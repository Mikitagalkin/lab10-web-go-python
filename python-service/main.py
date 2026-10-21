"""
FastAPI-сервис для сравнения с Go.

Задание 7 средней: graceful shutdown через lifespan.
"""

from contextlib import asynccontextmanager
from datetime import datetime

from fastapi import FastAPI


@asynccontextmanager
async def lifespan(app: FastAPI):
    """Graceful shutdown (задание 7 средней)."""
    print(f"[{datetime.now().isoformat()}] FastAPI starting...")
    yield
    print(f"[{datetime.now().isoformat()}] FastAPI shutting down...")


app = FastAPI(
    title="Python Service",
    description="FastAPI для сравнения с Go (ЛР10)",
    version="1.0.0",
    lifespan=lifespan,
)


@app.get("/health")
async def health():
    return {"status": "ok", "timestamp": datetime.now().isoformat()}


@app.get("/ping")
async def ping():
    return {
        "message": "pong from Python",
        "timestamp": datetime.now().timestamp(),
    }
