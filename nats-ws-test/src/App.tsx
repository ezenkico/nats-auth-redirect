import { useEffect, useState } from "react";
import { connect, type NatsConnection, StringCodec } from "nats.ws";

const sc = StringCodec();

export default function App() {
  const [status, setStatus] = useState("disconnected");
  const [messages, setMessages] = useState<string[]>([]);

  useEffect(() => {
    let nc: NatsConnection | undefined;
    let stopped = false;

    async function start() {
      try {
        console.log("Start")
        setStatus("connecting");

        nc = await connect({
          servers: "ws://localhost:8080",
          token: "test-token",
        });

        console.log("Test")

        setStatus("connected");

        const sub = nc.subscribe("requests");

        for await (const msg of sub) {
          console.log("Message Recieved.  Stopped: ", stopped);
          if (stopped) break;

          const text = sc.decode(msg.data);

          setMessages((prev) => [
            `${msg.subject}: ${text}`,
            ...prev.slice(0, 49),
          ]);
        }
      } catch (err) {
        console.error(err);
        setStatus("error");
      }
    }

    start();

    return () => {
      stopped = true;
      nc?.close();
    };
  }, []);

  return (
    <main style={{ padding: 24, fontFamily: "sans-serif" }}>
      <h1>NATS WebSocket Listener</h1>

      <p>
        Status: <strong>{status}</strong>
      </p>

      <section>
        <h2>Messages</h2>

        {messages.length === 0 ? (
          <p>No messages yet.</p>
        ) : (
          <ul>
            {messages.map((message, index) => (
              <li key={index}>
                <pre>{message}</pre>
              </li>
            ))}
          </ul>
        )}
      </section>
    </main>
  );
}