import { useEffect } from 'react';

/**
 * Subscribes to the server-sent events stream and calls onUpdate whenever
 * the server broadcasts a data change. The EventSource reconnects automatically
 * on connection loss.
 */
export function useSSE(onUpdate: () => void) {
  useEffect(() => {
    const es = new EventSource('/api/events/stream');
    es.onmessage = onUpdate;
    es.onerror = () => {}; // browser handles reconnection automatically
    return () => es.close();
  }, [onUpdate]);
}
