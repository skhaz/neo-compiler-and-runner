import asyncio
import hashlib
import os
import subprocess
import traceback
from tempfile import TemporaryDirectory

from google.cloud.storage import Client as StorageClient
from starlette.applications import Starlette
from starlette.requests import Request
from starlette.responses import Response
from starlette.routing import Route
from telegram import Update
from telegram.ext import Application
from telegram.ext import CommandHandler
from telegram.ext import ContextTypes
from wasmtime import Config
from wasmtime import Engine
from wasmtime import ExitTrap
from wasmtime import Func
from wasmtime import Linker
from wasmtime import Module
from wasmtime import Store
from wasmtime import WasiConfig

application = (
    Application.builder().token(os.environ["TELEGRAM_TOKEN"]).updater(None).build()
)

storage_client = StorageClient()
bucket = storage_client.bucket(os.environ["BUCKET"])


def run(source: str) -> str:
    with TemporaryDirectory() as path:
        os.chdir(path)

        with open("main.cpp", "w+t") as main:
            main.write(source)
            main.flush()

            try:
                result = subprocess.run(
                    [
                        "em++",
                        "-s",
                        "ENVIRONMENT=node",
                        "-s",
                        "WASM=1",
                        "-s",
                        "PURE_WASI=1",
                        "main.cpp",
                    ],
                    capture_output=True,
                    text=True,
                    check=True,
                )

                if result.returncode != 0:
                    return result.stderr

                with open("a.out.wasm", "rb") as binary:
                    wasi = WasiConfig()
                    wasi.stdout_file = "a.out.stdout"
                    wasi.stderr_file = "a.out.stderr"

                    config = Config()
                    # config.consume_fuel = True
                    engine = Engine(config)
                    store = Store(engine)
                    store.set_wasi(wasi)
                    # store.set_limits(16 * 1024 * 1024)
                    # store.set_fuel(10_000_000_000)

                    linker = Linker(engine)
                    linker.define_wasi()
                    module = Module(store.engine, binary.read())
                    instance = linker.instantiate(store, module)
                    start = instance.exports(store)["_start"]
                    assert isinstance(start, Func)

                    try:
                        start(store)
                    except ExitTrap as e:
                        if e.code != 0:
                            with open("a.out.stderr", "rt") as stderr:
                                return stderr.read()

                    with open("a.out.stdout", "rt") as stdout:
                        return stdout.read()
            except subprocess.CalledProcessError as e:
                return e.stderr
            except Exception as e:  # noqa
                return str(e)


async def on_run(update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
    message = update.message
    if not message:
        return

    text = message.text
    if not text:
        return

    text = text.lstrip("/run")

    if not text:
        await message.reply_text("Luke, I need the code for the Death Star's system.")
        return

    try:
        coro = asyncio.to_thread(run, text)
        result = await asyncio.wait_for(coro, timeout=30)
        if len(result) > 1:
            blob = bucket.blob(hashlib.sha256(str(text).encode()).hexdigest())
            blob.upload_from_string(result)
            blob.make_public()
            result = blob.public_url

        await message.reply_text(result)
    except asyncio.TimeoutError:
        await message.reply_text("⏰😮‍💨")
        return
    except Exception as e:
        await message.reply_text(f"{e}\n{traceback.format_exc()}")


def equals(left: str | None, right: str | None) -> bool:
    if not left or not right:
        return False

    if len(left) != len(right):
        return False

    for c1, c2 in zip(left, right):
        if c1 != c2:
            return False

    return True


async def webhook(request: Request):
    if not equals(
        request.headers.get("X-Telegram-Bot-Api-Secret-Token"),
        os.environ["SECRET"],
    ):
        return Response(content="Unauthorized", status_code=401)

    payload = await request.json()

    async with application:
        await application.process_update(Update.de_json(payload, application.bot))

    return Response(status_code=200)


application.add_handler(CommandHandler("run", on_run))

app = Starlette(
    debug=True,
    routes=[
        Route("/", webhook, methods=["POST"]),
    ],
)
