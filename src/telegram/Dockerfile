FROM python:3.12-slim-bookworm AS base

ENV PATH=/opt/venv/bin:$PATH
ENV PYTHONUNBUFFERED 1
ENV PYTHONDONTWRITEBYTECODE 1

FROM base AS builder
RUN python -m venv /opt/venv
COPY requirements.txt .
RUN pip install --no-cache-dir --requirement requirements.txt

FROM base
WORKDIR /opt/app
COPY --from=builder /opt/venv /opt/venv
COPY . .

RUN useradd -r user
USER user

CMD exec uvicorn main:app --host 0.0.0.0 --port $PORT