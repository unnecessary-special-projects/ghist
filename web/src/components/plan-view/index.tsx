import { useMemo, useState } from "react";
import cx from "classnames";
import {
  DndContext,
  DragEndEvent,
  DragOverlay,
  DragStartEvent,
  PointerSensor,
  useSensor,
  useSensors,
} from "@dnd-kit/core";
import { useDroppable, useDraggable } from "@dnd-kit/core";
import css from "./index.module.css";
import type { Task } from "../../types";
import {
  STATUS_COLORS,
  STATUS_LABELS,
  PRIORITY_COLORS,
  PRIORITY_LABELS,
  TYPE_COLORS,
  TYPE_LABELS,
} from "../../types";

const NONE_SENTINEL = "__none__";

export interface IPlanView {
  tasks: Task[];
  onMilestoneChange: (id: number, milestone: string) => void;
  onCardClick?: (task: Task) => void;
}

export const PlanView: React.FC<IPlanView> = (props) => {
  const [activeTask, setActiveTask] = useState<Task | null>(null);

  const sensors = useSensors(
    useSensor(PointerSensor, { activationConstraint: { distance: 5 } })
  );

  const milestones = useMemo(() => {
    const set = new Set<string>();
    for (const t of props.tasks) {
      set.add(t.milestone || "");
    }
    const sorted = [...set].sort((a, b) => {
      if (a === "") return 1;
      if (b === "") return -1;
      return a.localeCompare(b);
    });
    return sorted;
  }, [props.tasks]);

  const handleDragStart = (event: DragStartEvent) => {
    const task = props.tasks.find((t) => t.id === event.active.id);
    if (task) setActiveTask(task);
  };

  const handleDragEnd = (event: DragEndEvent) => {
    setActiveTask(null);
    const { active, over } = event;
    if (!over) return;

    const taskId = active.id as number;
    const dropId = over.id as string;
    const newMilestone = dropId === NONE_SENTINEL ? "" : dropId;

    const task = props.tasks.find((t) => t.id === taskId);
    if (task && (task.milestone || "") !== newMilestone) {
      props.onMilestoneChange(taskId, newMilestone);
    }
  };

  return (
    <DndContext
      sensors={sensors}
      onDragStart={handleDragStart}
      onDragEnd={handleDragEnd}
    >
      <div className={css.plan}>
        {milestones.map((milestone) => (
          <MilestoneGroup
            key={milestone === "" ? NONE_SENTINEL : milestone}
            milestone={milestone}
            tasks={props.tasks.filter(
              (t) => (t.milestone || "") === milestone
            )}
            onCardClick={props.onCardClick}
          />
        ))}
      </div>

      <DragOverlay>
        {activeTask ? <PlanCard task={activeTask} isOverlay /> : null}
      </DragOverlay>
    </DndContext>
  );
};

function MilestoneGroup({
  milestone,
  tasks,
  onCardClick,
}: {
  milestone: string;
  tasks: Task[];
  onCardClick?: (task: Task) => void;
}) {
  const droppableId = milestone === "" ? NONE_SENTINEL : milestone;
  const { setNodeRef, isOver } = useDroppable({ id: droppableId });

  return (
    <div
      ref={setNodeRef}
      className={cx(css.group, {
        [css.groupOver]: isOver,
        [css.groupEmpty]: tasks.length === 0,
      })}
    >
      <div className={css.groupHeader}>
        <span className={css.dot} />
        <span className={css.groupLabel}>
          {milestone === "" ? "No milestone" : milestone}
        </span>
        <span className={css.count}>{tasks.length}</span>
      </div>

      {tasks.length > 0 && (
        <div className={css.cards}>
          {tasks.map((task) => (
            <PlanCard key={task.id} task={task} onCardClick={onCardClick} />
          ))}
        </div>
      )}
    </div>
  );
}

function PlanCard({
  task,
  onCardClick,
  isOverlay,
}: {
  task: Task;
  onCardClick?: (task: Task) => void;
  isOverlay?: boolean;
}) {
  const { attributes, listeners, setNodeRef, isDragging } = useDraggable({
    id: task.id,
  });

  return (
    <div
      ref={setNodeRef}
      {...listeners}
      {...attributes}
      onClick={() => onCardClick?.(task)}
      className={cx(css.row, {
        [css.rowDragging]: isDragging && !isOverlay,
        [css.rowOverlay]: isOverlay,
      })}
    >
      <span className={css.refId}>{task.ref_id}</span>
      <span className={css.title}>{task.title}</span>
      <div className={css.pills}>
        <span
          className={css.pill}
          style={{
            color: STATUS_COLORS[task.status],
            backgroundColor: `${STATUS_COLORS[task.status]}1a`,
          }}
        >
          {STATUS_LABELS[task.status]}
        </span>
        {task.priority && (
          <span
            className={css.pill}
            style={{
              color: PRIORITY_COLORS[task.priority],
              backgroundColor: `${PRIORITY_COLORS[task.priority]}1a`,
            }}
          >
            {PRIORITY_LABELS[task.priority]}
          </span>
        )}
        {task.type && (
          <span
            className={css.pill}
            style={{
              color: TYPE_COLORS[task.type],
              backgroundColor: `${TYPE_COLORS[task.type]}1a`,
            }}
          >
            {TYPE_LABELS[task.type]}
          </span>
        )}
      </div>
    </div>
  );
}
