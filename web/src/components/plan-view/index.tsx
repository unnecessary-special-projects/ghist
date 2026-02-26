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
import { useDroppable } from "@dnd-kit/core";
import css from "./index.module.css";
import type { Task } from "../../types";
import { STATUS_COLORS } from "../../types";
import { TaskCard } from "../task-card";

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
        {activeTask ? <TaskCard task={activeTask} isOverlay /> : null}
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
            <div key={task.id} className={css.cardRow}>
              <span
                className={css.statusDot}
                style={{ backgroundColor: STATUS_COLORS[task.status] }}
                title={task.status}
              />
              <div className={css.cardWrapper}>
                <TaskCard task={task} onCardClick={onCardClick} />
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
