import asyncio
import json
import os

from nats.aio.client import Client as NATS


async def main():
    nc = NATS()

    await nc.connect(
        servers=["ws://localhost:8080"],
        user="app",
        password="app",
    )

    subject = os.getenv("NATS_SUBJECT", "requests")

    payload = b"hello"

    await nc.publish(subject, payload)
    await nc.flush()

    print(f"Published to {subject}: {payload}")

    await nc.close()


if __name__ == "__main__":
    asyncio.run(main())