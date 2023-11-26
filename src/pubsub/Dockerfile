FROM python:3.12-slim AS base

ENV PATH /opt/venv/bin:$PATH
ENV PYTHONUNBUFFERED 1
ENV PYTHONDONTWRITEBYTECODE 1

FROM base AS builder
RUN python -m venv /opt/venv
COPY requirements.txt .
RUN pip install --no-cache-dir --requirement requirements.txt

FROM base
WORKDIR /app
COPY --from=builder /opt/venv /opt/venv
COPY . .

RUN useradd -r user
USER user

CMD exec gunicorn --bind :$PORT --workers 1 --threads 8 --timeout 0 main:app