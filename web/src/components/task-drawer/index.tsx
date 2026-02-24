import cx from "classnames";
import { useCallback, useEffect, useState } from "react";
import css from "./index.module.css";
import type { Task, TaskStatus, TaskPriority, TaskType, EventType, Event } from "../../types";
import {
  STATUSES, STATUS_LABELS, STATUS_COLORS,
  PRIORITIES, PRIORITY_LABELS, PRIORITY_COLORS,
  TASK_TYPES, TYPE_LABELS, TYPE_COLORS,
  EVENT_TYPES, EVENT_TYPE_LABELS, EVENT_TYPE_COLORS,
} from "../../types";
import { InlineField } from "../inline-field";
import { MarkdownPreview } from "../markdown-preview";
import * as api from "../../api/client";

export interface ITaskDrawer {
  task: Task | null;
  mode: "view" | "create" | null;
  onClose: () => void;
  onUpdateTask: (id: number, data: Record<string, string>) => void;
  onCreateTask: (data: { title: string; description: string; status: TaskStatus; milestone: string; priority: TaskPriority; type: TaskType }) => void;
  onDeleteTask: (id: number) => void;
}

type Tab = "details" | "plan" | "activity";

export const TaskDrawer: React.FC<ITaskDrawer> = (props) => {
  const [tab, setTab] = useState<Tab>("details");
  const open = props.mode !== null;

  const [createTitle, setCreateTitle] = useState("");
  const [createDescription, setCreateDescription] = useState("");
  const [createStatus, setCreateStatus] = useState<TaskStatus>("todo");
  const [createMilestone, setCreateMilestone] = useState("");
  const [createPriority, setCreatePriority] = useState<TaskPriority>("");
  const [createType, setCreateType] = useState<TaskType>("");

  useEffect(() => {
    if (open) setTab("details");
  }, [open, props.task?.id, props.mode]);

  useEffect(() => {
    if (props.mode === "create") {
      setCreateTitle("");
      setCreateDescription("");
      setCreateStatus("todo");
      setCreateMilestone("");
      setCreatePriority("");
      setCreateType("");
    }
  }, [props.mode]);

  useEffect(() => {
    const handleKey = (e: KeyboardEvent) => {
      if (e.key === "Escape" && open) props.onClose();
    };
    document.addEventListener("keydown", handleKey);
    return () => document.removeEventListener("keydown", handleKey);
  }, [open, props.onClose]);

  const handleCreate = () => {
    if (!createTitle.trim()) return;
    props.onCreateTask({
      title: createTitle.trim(),
      description: createDescription,
      status: createStatus,
      milestone: createMilestone,
      priority: createPriority,
      type: createType,
    });
  };

  return (
    <>
      <div
        className={cx(css.backdrop, open ? css.backdropVisible : css.backdropHidden)}
        onClick={props.onClose}
      />
      <div className={cx(css.drawer, open ? css.drawerOpen : css.drawerClosed)}>
        {props.mode === "create" ? (
          <CreateContent
            title={createTitle} onTitleChange={setCreateTitle}
            description={createDescription} onDescriptionChange={setCreateDescription}
            status={createStatus} onStatusChange={setCreateStatus}
            milestone={createMilestone} onMilestoneChange={setCreateMilestone}
            priority={createPriority} onPriorityChange={setCreatePriority}
            type={createType} onTypeChange={setCreateType}
            onClose={props.onClose} onCreate={handleCreate}
          />
        ) : props.task ? (
          <>
            <div className={css.header}>
              <div className={css.headerTitle}>
                <span className={css.taskRefId}>{props.task.ref_id}</span>
                <InlineField
                  label=""
                  value={props.task.title}
                  onSave={(v) => props.onUpdateTask(props.task!.id, { title: v })}
                  renderValue={(v) => <span className={css.taskTitleText}>{v}</span>}
                  className={css.titleField}
                />
              </div>
              <button className={css.closeBtn} onClick={props.onClose}>&times;</button>
            </div>

            <div className={css.tabs}>
              <div className={css.tabGroup}>
                <button
                  className={cx(css.tab, { [css.tabActive]: tab === "details" })}
                  onClick={() => setTab("details")}
                >
                  Details
                </button>
                <button
                  className={cx(css.tab, { [css.tabActive]: tab === "plan" })}
                  onClick={() => setTab("plan")}
                >
                  Plan
                </button>
                <button
                  className={cx(css.tab, { [css.tabActive]: tab === "activity" })}
                  onClick={() => setTab("activity")}
                >
                  Activity
                </button>
              </div>
            </div>

            <div className={css.content}>
              {tab === "details" ? (
                <DetailsTab task={props.task} onUpdate={props.onUpdateTask} onDelete={props.onDeleteTask} />
              ) : tab === "plan" ? (
                <PlanTab task={props.task} onUpdate={props.onUpdateTask} />
              ) : (
                <ActivityTab task={props.task} />
              )}
            </div>
          </>
        ) : null}
      </div>
    </>
  );
};

// ---------- Create Content ----------

function CreateContent({
  title, onTitleChange, description, onDescriptionChange,
  status, onStatusChange, milestone, onMilestoneChange,
  priority, onPriorityChange, type, onTypeChange,
  onClose, onCreate,
}: {
  title: string; onTitleChange: (v: string) => void;
  description: string; onDescriptionChange: (v: string) => void;
  status: TaskStatus; onStatusChange: (v: TaskStatus) => void;
  milestone: string; onMilestoneChange: (v: string) => void;
  priority: TaskPriority; onPriorityChange: (v: TaskPriority) => void;
  type: TaskType; onTypeChange: (v: TaskType) => void;
  onClose: () => void;
  onCreate: () => void;
}) {
  return (
    <>
      <div className={css.header}>
        <span className={css.taskTitleText}>New Task</span>
        <button className={css.closeBtn} onClick={onClose}>&times;</button>
      </div>
      <div className={css.content}>
        <div className={css.details}>
          <div className={css.formField}>
            <span className={css.fieldLabel}>Title</span>
            <input className={css.input} value={title} onChange={(e) => onTitleChange(e.target.value)} placeholder="Task title" autoFocus />
          </div>
          <div className={css.formField}>
            <span className={css.fieldLabel}>Status</span>
            <select className={css.input} value={status} onChange={(e) => onStatusChange(e.target.value as TaskStatus)}>
              {STATUSES.map((s) => <option key={s} value={s}>{STATUS_LABELS[s]}</option>)}
            </select>
          </div>
          <div className={css.formField}>
            <span className={css.fieldLabel}>Priority</span>
            <select className={css.input} value={priority} onChange={(e) => onPriorityChange(e.target.value as TaskPriority)}>
              {PRIORITIES.map((p) => <option key={p} value={p}>{PRIORITY_LABELS[p]}</option>)}
            </select>
          </div>
          <div className={css.formField}>
            <span className={css.fieldLabel}>Type</span>
            <select className={css.input} value={type} onChange={(e) => onTypeChange(e.target.value as TaskType)}>
              {TASK_TYPES.map((t) => <option key={t} value={t}>{TYPE_LABELS[t]}</option>)}
            </select>
          </div>
          <div className={css.formField}>
            <span className={css.fieldLabel}>Milestone</span>
            <input className={css.input} value={milestone} onChange={(e) => onMilestoneChange(e.target.value)} placeholder="e.g. v1.0" />
          </div>
          <div className={css.formField}>
            <span className={css.fieldLabel}>Description</span>
            <textarea className={`${css.input} ${css.inputTextarea}`} value={description} onChange={(e) => onDescriptionChange(e.target.value)} placeholder="Optional description" />
          </div>
          <button className={css.createBtn} onClick={onCreate}>Create Task</button>
        </div>
      </div>
    </>
  );
}

// ---------- Details Tab ----------

function DetailsTab({ task, onUpdate, onDelete }: { task: Task; onUpdate: (id: number, data: Record<string, string>) => void; onDelete: (id: number) => void }) {
  const save = (field: string) => (value: string) => onUpdate(task.id, { [field]: value });

  return (
    <div className={css.details}>
      <InlineField label="Ref ID" value={task.ref_id} onSave={() => {}} readOnly />
      <InlineField
        label="Status" value={task.status} onSave={save("status")} type="select"
        options={STATUSES.map((s) => ({ value: s, label: STATUS_LABELS[s], color: STATUS_COLORS[s] }))}
        renderValue={(v) => (
          <span className={css.statusBadge}>
            <span className={css.statusDot} style={{ backgroundColor: STATUS_COLORS[v as TaskStatus] }} />
            {STATUS_LABELS[v as TaskStatus]}
          </span>
        )}
      />
      <InlineField
        label="Priority" value={task.priority} onSave={save("priority")} type="select"
        options={PRIORITIES.map((p) => ({ value: p, label: PRIORITY_LABELS[p], color: PRIORITY_COLORS[p] }))}
        renderValue={(v) => {
          const p = v as TaskPriority;
          if (!p) return <span className={css.noneText}>None</span>;
          return (
            <span className={css.statusBadge}>
              <span className={css.statusDot} style={{ backgroundColor: PRIORITY_COLORS[p] }} />
              {PRIORITY_LABELS[p]}
            </span>
          );
        }}
      />
      <InlineField
        label="Type" value={task.type} onSave={save("type")} type="select"
        options={TASK_TYPES.map((t) => ({ value: t, label: TYPE_LABELS[t], color: TYPE_COLORS[t] }))}
        renderValue={(v) => {
          const t = v as TaskType;
          if (!t) return <span className={css.noneText}>None</span>;
          return (
            <span className={css.typeBadge} style={{ color: TYPE_COLORS[t], backgroundColor: `${TYPE_COLORS[t]}1a` }}>
              {TYPE_LABELS[t]}
            </span>
          );
        }}
      />
      <InlineField label="Milestone" value={task.milestone} onSave={save("milestone")} />
      <InlineField label="Description" value={task.description} onSave={save("description")} type="textarea" />
      {task.commit_hash && (
        <InlineField label="Commit" value={task.commit_hash} onSave={() => {}} readOnly
          renderValue={(v) => <code className={css.commitHash}>{v.slice(0, 8)}</code>}
        />
      )}
      {task.legacy_id && (
        <InlineField label="Legacy ID" value={task.legacy_id} onSave={save("legacy_id")} />
      )}
      <InlineField label="Created" value={formatDate(task.created_at)} onSave={() => {}} readOnly />
      <InlineField label="Updated" value={formatDate(task.updated_at)} onSave={() => {}} readOnly />

      <div className={css.drawerActions}>
        <button className={css.deleteBtn} onClick={() => onDelete(task.id)}>Delete Task</button>
      </div>
    </div>
  );
}

// ---------- Plan Tab ----------

function PlanTab({ task, onUpdate }: { task: Task; onUpdate: (id: number, data: Record<string, string>) => void }) {
  const [editMode, setEditMode] = useState(false);
  const [draft, setDraft] = useState(task.plan);

  useEffect(() => {
    setDraft(task.plan);
    setEditMode(false);
  }, [task.id, task.plan]);

  const toggleEdit = () => {
    if (editMode) {
      if (draft !== task.plan) {
        onUpdate(task.id, { plan: draft });
      }
      setEditMode(false);
    } else {
      setEditMode(true);
    }
  };

  return (
    <div className={css.planContainer}>
      <div className={css.planHeader}>
        <button className={css.planToggle} onClick={toggleEdit}>
          {editMode ? "Preview" : "Edit"}
        </button>
      </div>
      {editMode ? (
        <textarea
          className={css.planEditor}
          value={draft}
          onChange={(e) => setDraft(e.target.value)}
          onBlur={() => { if (draft !== task.plan) onUpdate(task.id, { plan: draft }); }}
          placeholder="Write your plan in markdown..."
          autoFocus
        />
      ) : task.plan ? (
        <MarkdownPreview content={task.plan} />
      ) : (
        <div className={css.emptyPlan}>
          <p className={css.emptyText}>No plan yet</p>
          <p className={css.emptyHint}>Click Edit to write a plan manually, or ask your AI coding agent to plan this task</p>
        </div>
      )}
    </div>
  );
}

// ---------- Activity Tab ----------

function ActivityTab({ task }: { task: Task }) {
  const [events, setEvents] = useState<Event[]>([]);
  const [message, setMessage] = useState('');
  const [type, setType] = useState<EventType>('log');
  const [submitting, setSubmitting] = useState(false);

  const loadEvents = useCallback(async () => {
    const data = await api.listTaskEvents(task.id);
    setEvents(data);
  }, [task.id]);

  useEffect(() => {
    loadEvents();
  }, [loadEvents]);

  const handleSubmit = async () => {
    if (!message.trim()) return;
    setSubmitting(true);
    try {
      await api.createEvent({ type, message: message.trim(), task_id: task.id });
      setMessage('');
      await loadEvents();
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className={css.activityContainer}>
      <div className={css.eventList}>
        {events.length === 0 ? (
          <div className={css.emptyActivity}>
            <p className={css.emptyText}>No activity yet</p>
            <p className={css.emptyHint}>Log decisions and notes as you work on this task</p>
          </div>
        ) : (
          events.map((e) => (
            <div key={e.id} className={css.eventItem}>
              <div className={css.eventItemHeader}>
                <span
                  className={css.eventTypeBadge}
                  style={{
                    color: EVENT_TYPE_COLORS[e.type as EventType] ?? '#768390',
                    backgroundColor: `${EVENT_TYPE_COLORS[e.type as EventType] ?? '#768390'}1a`,
                  }}
                >
                  {EVENT_TYPE_LABELS[e.type as EventType] ?? e.type}
                </span>
                <span className={css.eventTime}>{formatDate(e.created_at)}</span>
              </div>
              <p className={css.eventMessage}>{e.message}</p>
            </div>
          ))
        )}
      </div>
      <div className={css.eventForm}>
        <textarea
          className={`${css.input} ${css.inputTextarea}`}
          value={message}
          onChange={(e) => setMessage(e.target.value)}
          placeholder="Add a note or decision..."
          rows={3}
        />
        <div className={css.eventFormActions}>
          <select
            className={css.input}
            value={type}
            onChange={(e) => setType(e.target.value as EventType)}
            style={{ width: 'auto' }}
          >
            {EVENT_TYPES.map((t) => (
              <option key={t} value={t}>{EVENT_TYPE_LABELS[t]}</option>
            ))}
          </select>
          <button
            className={css.createBtn}
            onClick={handleSubmit}
            disabled={submitting || !message.trim()}
          >
            Add
          </button>
        </div>
      </div>
    </div>
  );
}

function formatDate(iso: string): string {
  if (!iso) return "\u2014";
  const d = new Date(iso);
  return d.toLocaleDateString(undefined, { year: "numeric", month: "short", day: "numeric" })
    + " " + d.toLocaleTimeString(undefined, { hour: "2-digit", minute: "2-digit" });
}
