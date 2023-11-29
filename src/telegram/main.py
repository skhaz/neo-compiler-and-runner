import json
import os
import traceback

from google.cloud.pubsublite_v1.types.publisher import PublishRequest
from google.cloud.pubsublite_v1.types import PubSubMessage
from google.cloud.pubsublite_v1 import PublisherServiceAsyncClient
from starlette.applications import Starlette
from starlette.requests import Request
from starlette.responses import Response
from starlette.routing import Route
from google.cloud.pubsublite_v1 import PublisherServiceAsyncClient, PublisherOptions
from google.cloud.pubsublite_v1.types import PubSubMessage, PublishRequest

from telegram import Update
from telegram.ext import Application
from telegram.ext import CommandHandler
from telegram.ext import ContextTypes

application = (
    Application.builder().token(os.environ["TELEGRAM_TOKEN"]).updater(None).build()
)

pubsub = PublisherServiceAsyncClient()


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
        payload = {
            "message": {
                "chat": {
                    "id": message.chat.id,
                },
                "message": {
                    "id": message.message_id,
                },
                "source": text,
            }
        }

        client_options = PublisherOptions(topic_path=os.environ["PUBSUB_TOPIC_PATH"])
        pubsub = PublisherServiceAsyncClient(client_options=client_options)

        message = PubSubMessage(data=json.dumps(payload).encode("utf-8"))

        request = PublishRequest(messages=[message])

        async def request_generator():
            yield request

        try:
            await pubsub.publish(requests=request_generator())
        except Exception as e:
            await message.reply_text(f"{e}\n{traceback.format_exc()}")
            return

        await message.reply_text("Ok")

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


application.add_handler(CommandHandler("run2", on_run))

app = Starlette(
    debug=True,
    routes=[
        Route("/", webhook, methods=["POST"]),
    ],
)
