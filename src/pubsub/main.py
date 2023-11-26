import base64
import functools
import hashlib
import json
import os
import subprocess
from contextlib import contextmanager
from http import HTTPStatus
from tempfile import TemporaryDirectory

from flask import Flask
from flask import abort
from flask import request
from google.cloud.storage import Client as StorageClient

app = Flask(__name__)

storage_client = StorageClient()
bucket = storage_client.bucket(os.environ["BUCKET"])


@contextmanager
def directory(path):
    original_dir = os.getcwd()
    try:
        os.chdir(path)
        yield
    finally:
        os.chdir(original_dir)


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

            message = envelope["message"]

            data = None

            if isinstance(message, dict) and "data" in message:
                data = base64.b64decode(message["data"]).decode("utf-8").strip()

            if not data:
                message = "received an empty message from Pub/Sub"
                print(message)
                abort(HTTPStatus.BAD_REQUEST, description=f"Bad Request: {message}")

            return func(json.loads(data))

        return wrapper

    return decorator


def run(source: str) -> str:
    with TemporaryDirectory() as path:
        with directory(path):
            with open("main.cpp", "w+t") as main:
                main.write(source)
                main.flush()

                command = [
                    "emcc",
                    "-O3",
                    "-flto",
                    "-s",
                    "ENVIRONMENT=node",
                    # "-s",
                    # "PURE_WASI=1",
                    "-s",
                    "WASM=1",
                    "main.cpp",
                ]

                result = subprocess.run(command, capture_output=True, text=True)

                if result.returncode != 0:
                    raise Exception(result.stderr)

                command = [
                    "node",
                    "a.out.js",
                ]

                result = subprocess.run(command, capture_output=True, text=True)

                if result.returncode != 0:
                    raise Exception(result.stderr)

                return result.stdout


@app.post("/")
@unenvelop()
def index(data):
    if "source" not in data:
        abort(HTTPStatus.NO_CONTENT)

    result = run(data["source"])

    # if len(result) > 1024:
    #    # upload to storage

    filename = f"{hashlib.sha256(str(data).encode()).hexdigest()}.txt"

    blob = bucket.blob(filename)
    blob.upload_from_string(result)
    blob.make_public()
    print(blob.public_url)
