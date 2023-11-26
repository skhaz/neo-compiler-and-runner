import base64
import functools
import os
from http import HTTPStatus

from flask import Flask
from flask import abort
from flask import request
from google.cloud.storage import Client as StorageClient

app = Flask(__name__)

storage_client = StorageClient()
bucket = storage_client.bucket(os.environ["BUCKET"])


def unenvelop():
    def decorator(func):
        @functools.wraps(func)
        def wrapper(*args, **kwargs):
            envelope = request.get_json()

            if not envelope:
                message = "no Pub/Sub message received"
                print(message)
                abort(HTTPStatus.BAD_REQUEST, description=f"Bad Request: {message}")

            if not isinstance(envelope, dict) or "message" not in envelope:
                message = "invalid Pub/Sub message format"
                print(message)
                abort(HTTPStatus.BAD_REQUEST, description=f"Bad Request: {message}")

            pubsub_message = envelope["message"]

            data = None

            if isinstance(pubsub_message, dict) and "data" in pubsub_message:
                data = base64.b64decode(pubsub_message["data"]).decode("utf-8").strip()

            if not data:
                message = "received an empty message from Pub/Sub"
                print(message)
                abort(HTTPStatus.BAD_REQUEST, description=f"Bad Request: {message}")

            return func(data)

        return wrapper

    return decorator


@app.post("/")
@unenvelop()
def index(data):
    print("Received data: ", data)
    blob = bucket.blob("data.txt")
    blob.upload_from_string(data)
    blob.make_public()
    print("Uploaded to: ", blob.public_url)
