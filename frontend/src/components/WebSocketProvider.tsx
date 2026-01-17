import { useEffect, useRef, useState } from "react";

export function useWebSocket(url: string) {
  const ws = useRef<WebSocket | null>(null);
  const [messages, setMessages] = useState<string[]>([]);

  useEffect(() => {
    ws.current = new WebSocket(url);

    ws.current.onopen = () => console.log("WebSocket connected");
    ws.current.onmessage = (event) => {
      setMessages((prev) => [...prev, event.data]);
    };
    ws.current.onerror = (err) => console.error("WebSocket error", err);
    ws.current.onclose = () => console.log("WebSocket disconnected");

    return () => {
      ws.current?.close();
    };
  }, [url]);

  const sendMessage = (msg: string) => {
    if (ws.current?.readyState === WebSocket.OPEN) {
      ws.current.send(msg);
    }
  };

  return { messages, sendMessage };
}
