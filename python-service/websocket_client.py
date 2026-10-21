"""
WebSocket-клиент для подключения к Go-чату.

Задание 4 повышенной сложности: WebSocket чат.
"""

import asyncio
import json

import websockets


async def receive_messages(ws):
    """Получает сообщения от сервера и печатает."""
    try:
        async for message in ws:
            data = json.loads(message)
            print(f"[RECEIVED] {data['username']}: {data['text']}")
    except websockets.exceptions.ConnectionClosed:
        print("[INFO] Соединение закрыто")
    except asyncio.CancelledError:
        raise


async def main():
    """Подключается к Go-чату, отправляет сообщения, получает ответы."""
    uri = "ws://localhost:8080/ws"
    print(f"Подключение к {uri}...")

    async with websockets.connect(uri) as ws:
        print("[INFO] Соединение установлено\n")

        # ✅ Сохраняем ссылку на задачу (важно!)
        receive_task = asyncio.create_task(receive_messages(ws))

        messages = [
            {"username": "python_user", "text": "Привет из Python!"},
            {"username": "python_user", "text": "Как дела, Go?"},
            {"username": "python_user", "text": "До связи!"},
        ]

        for msg in messages:
            await ws.send(json.dumps(msg))
            print(f"[SENT] {msg['username']}: {msg['text']}")
            await asyncio.sleep(1)

        # Дать время на получение
        await asyncio.sleep(2)

        # Корректно отменяем задачу
        receive_task.cancel()
        try:
            await receive_task
        except asyncio.CancelledError:
            pass

    print("\n[INFO] Клиент завершил работу")


if __name__ == "__main__":
    try:
        asyncio.run(main())
    except KeyboardInterrupt:
        print("\n[INFO] Прервано пользователем")
