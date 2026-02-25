import { useCallback, useEffect, useState } from "react";
import css from "./index.module.css";
import type { Event, EventType, Task } from "../../types";
import { EVENT_TYPES, EVENT_TYPE_LABELS, EVENT_TYPE_COLORS } from "../../types";
import * as api from "../../api/client";
import { useSSE } from "../../hooks/useSSE";

interface IActivityFeed {
  tasks: Task[];
}

export const ActivityFeed: React.FC<IActivityFeed> = ({ tasks }) => {
  const [events, setEvents] = useState<Event[]>([]);
  const [message, setMessage] = useState("");
  const [type, setType] = useState<EventType>("log");
  const [submitting, setSubmitting] = useState(false);

  const loadEvents = useCallback(async () => {
    const data = await api.listEvents(100);
    setEvents(data);
  }, []);

  useEffect(() => {
    loadEvents();
  }, [loadEvents]);

  useSSE(loadEvents);

  const handleSubmit = async () => {
    if (!message.trim()) return;
    setSubmitting(true);
    try {
      await api.createEvent({ type, message: message.trim() });
      setMessage("");
      await loadEvents();
    } finally {
      setSubmitting(false);
    }
  };

  const taskMap = Object.fromEntries(tasks.map((t) => [t.id, t]));

  return (
    <div className={css.container}>
      <div className={css.inner}>
        <div className={css.form}>
          <textarea
            className={css.textarea}
            value={message}
            onChange={(e) => setMessage(e.target.value)}
            placeholder="Log a decision, note, or update..."
            rows={3}
          />
          <div className={css.formActions}>
            <select
              className={css.select}
              value={type}
              onChange={(e) => setType(e.target.value as EventType)}
            >
              {EVENT_TYPES.map((t) => (
                <option key={t} value={t}>{EVENT_TYPE_LABELS[t]}</option>
              ))}
            </select>
            <button
              className={css.addBtn}
              onClick={handleSubmit}
              disabled={submitting || !message.trim()}
            >
              Add Entry
            </button>
          </div>
        </div>

        {events.length === 0 ? (
          <div className={css.empty}>
            <p className={css.emptyText}>No activity yet</p>
            <p className={css.emptyHint}>Decisions, notes, and task updates will appear here</p>
          </div>
        ) : (
          <div className={css.timeline}>
            {events.map((e) => {
              const color = EVENT_TYPE_COLORS[e.type as EventType] ?? "#768390";
              const linkedTask = e.task_id ? taskMap[e.task_id] : null;
              return (
                <div key={e.id} className={css.event}>
                  <div className={css.eventLeft}>
                    <span
                      className={css.typeBadge}
                      style={{ color, backgroundColor: `${color}1a` }}
                    >
                      {EVENT_TYPE_LABELS[e.type as EventType] ?? e.type}
                    </span>
                  </div>
                  <div className={css.eventBody}>
                    <p className={css.eventMessage}>{e.message}</p>
                    <div className={css.eventMeta}>
                      <span className={css.eventTime}>{formatDate(e.created_at)}</span>
                      {linkedTask && (
                        <span className={css.taskChip}>{linkedTask.ref_id} {linkedTask.title}</span>
                      )}
                    </div>
                  </div>
                </div>
              );
            })}
          </div>
        )}
      </div>
    </div>
  );
};

function formatDate(iso: string): string {
  if (!iso) return "";
  const d = new Date(iso);
  return d.toLocaleDateString(undefined, { year: "numeric", month: "short", day: "numeric" })
    + " " + d.toLocaleTimeString(undefined, { hour: "2-digit", minute: "2-digit" });
}
