import cx from "classnames";
import { DndContext, DragEndEvent, DragOverlay, DragStartEvent, PointerSensor, useSensor, useSensors } from "@dnd-kit/core";
import { useDroppable } from "@dnd-kit/core";
import { useState } from "react";
import css from "./index.module.css";
import type { Task, TaskStatus } from "../../types";
import { BOARD_STATUSES, STATUS_LABELS, STATUS_COLORS } from "../../types";
import { TaskCard } from "../task-card";

export interface IBoard {
  tasks: Task[];
  onStatusChange: (id: number, status: TaskStatus) => void;
  onCardClick?: (task: Task) => void;
}

export const Board: React.FC<IBoard> = (props) => {
  const [activeTask, setActiveTask] = useState<Task | null>(null);

  const sensors = useSensors(
    useSensor(PointerSensor, { activationConstraint: { distance: 5 } })
  );

  const handleDragStart = (event: DragStartEvent) => {
    const task = props.tasks.find((t) => t.id === event.active.id);
    if (task) setActiveTask(task);
  };

  const handleDragEnd = (event: DragEndEvent) => {
    setActiveTask(null);
    const { active, over } = event;
    if (!over) return;

    const taskId = active.id as number;
    const newStatus = over.id as TaskStatus;

    const task = props.tasks.find((t) => t.id === taskId);
    if (task && task.status !== newStatus) {
      props.onStatusChange(taskId, newStatus);
    }
  };

  return (
    <DndContext sensors={sensors} onDragStart={handleDragStart} onDragEnd={handleDragEnd}>
      <div className={css.columns}>
        {BOARD_STATUSES.map((status) => (
          <Column
            key={status}
            status={status}
            tasks={props.tasks.filter((t) => t.status === status)}
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

function Column({ status, tasks, onCardClick }: { status: TaskStatus; tasks: Task[]; onCardClick?: (task: Task) => void }) {
  const { setNodeRef, isOver } = useDroppable({ id: status });

  return (
    <div
      ref={setNodeRef}
      className={cx(css.column, { [css.columnOver]: isOver, [css.columnEmpty]: tasks.length === 0 })}
      style={isOver ? { borderColor: STATUS_COLORS[status] } : undefined}
    >
      <div className={css.header}>
        <span className={css.dot} style={{ backgroundColor: STATUS_COLORS[status] }} />
        <span className={css.label}>{STATUS_LABELS[status]}</span>
        <span className={css.count}>{tasks.length}</span>
      </div>
      <div className={css.cards}>
        {tasks.map((task) => (
          <TaskCard key={task.id} task={task} onCardClick={onCardClick} />
        ))}
      </div>
    </div>
  );
}
