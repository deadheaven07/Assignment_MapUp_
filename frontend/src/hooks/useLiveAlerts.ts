import { useEffect, useState } from 'react';
import type { AlertEvent } from '../types';

export function useLiveAlerts() {
  const [alerts, setAlerts] = useState<AlertEvent[]>([]);
  const [connected, setConnected] = useState(false);

  useEffect(() => {
    let socket: WebSocket | null = null;
    let retry: number | undefined;
    let closed = false;

    const connect = () => {
      socket = new WebSocket(import.meta.env.VITE_WS_URL ?? 'ws://localhost:8080/ws/alerts');
      socket.onopen = () => setConnected(true);
      socket.onclose = () => {
        setConnected(false);
        if (!closed) retry = window.setTimeout(connect, 1500);
      };
      socket.onmessage = (message) => {
        const event = JSON.parse(message.data) as AlertEvent;
        setAlerts((current) => [event, ...current].slice(0, 30));
      };
    };

    connect();
    return () => {
      closed = true;
      window.clearTimeout(retry);
      socket?.close();
    };
  }, []);

  return { alerts, connected };
}
