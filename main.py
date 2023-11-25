import asyncio
import os
import subprocess
from io import BufferedReader
from tempfile import NamedTemporaryFile
from tempfile import TemporaryDirectory
from typing import Tuple

from starlette.applications import Starlette
from starlette.requests import Request
from starlette.responses import Response
from starlette.routing import Route
from telegram import Update
from telegram.ext import Application
from telegram.ext import CommandHandler
from telegram.ext import ContextTypes
from wasmtime import Engine
from wasmtime import ExitTrap
from wasmtime import Func
from wasmtime import Linker
from wasmtime import Module
from wasmtime import Store
from wasmtime import WasiConfig


def compile(source: str) -> Tuple[bool, str, str]:
    with open("main.cpp", "w+t") as fh:
        fh.write(source)
        fh.flush()

        command = [
            "emcc",
            "-O3",
            "-flto",
            "-s",
            "ENVIRONMENT=node",
            "-s",
            "PURE_WASI=1",
            "-s",
            "WASM=1",
            "main.cpp",
        ]

        result = subprocess.run(command, capture_output=True, text=True)

        return result.returncode == 0, result.stdout, f"compile error {result.stderr}"


def execute(binary: BufferedReader) -> Tuple[bool, str, str]:
    with NamedTemporaryFile(mode="w+t", delete=False) as stdout:
        with NamedTemporaryFile(mode="w+t", delete=False) as stderr:
            wasi = WasiConfig()
            wasi.stdout_file = stdout.name
            wasi.stderr_file = stderr.name

            engine = Engine()
            store = Store(engine)
            store.set_wasi(wasi)
            linker = Linker(engine)
            linker.define_wasi()
            module = Module(store.engine, binary.read())
            instance = linker.instantiate(store, module)
            start = instance.exports(store)["_start"]
            assert isinstance(start, Func)

            try:
                start(store)
            except Exception as exc:
                stdout.seek(0)
                stderr.seek(0)
                return False, stdout.read(), f"{exc}: {stderr.read()}"
            stdout.seek(0)
            stderr.seek(0)
            return False, stdout.read(), stderr.read()


def run(source: str) -> Tuple[bool, str, str]:
    with TemporaryDirectory() as path:
        os.chdir(path)

        result, stdout, stderr = compile(source)
        if not result:
            return result, stdout, stderr

        with open("a.out.wasm", "rb") as binary:
            return execute(binary)


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


async def on_run(update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
    message = update.message
    if not message:
        return

    text = message.text
    if not text:
        await message.reply_text("not text.")
        return

    text = text.lstrip("/run")

    if not text:
        await message.reply_text("Please provide a source code.")
        return

    loop = asyncio.get_event_loop()

    await message.reply_text("Running...")

    result, stdout, stderr = await loop.run_in_executor(None, run, text)

    await message.reply_text("Done.")

    if result:
        await message.reply_text(f"Stdout {stdout}")
        return

    if not result:
        await message.reply_text(f"Stderr {stderr}")
        return


application = (
    Application.builder().token(os.environ["TELEGRAM_TOKEN"]).updater(None).build()
)

application.add_handler(CommandHandler("run", on_run))

app = Starlette(
    debug=True,
    routes=[
        Route("/", webhook, methods=["POST"]),
    ],
)
