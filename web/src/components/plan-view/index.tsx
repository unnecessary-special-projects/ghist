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
import {
  SortableContext,
  useSortable,
  arrayMove,
  verticalListSortingStrategy,
} from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import { DotsSixVertical } from "@phosphor-icons/react";
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
import { orderMilestones } from "../../utils/milestoneOrder";

const NONE_SENTINEL = "__none__";
const MILESTONE_PREFIX = "milestone:";

export interface IPlanView {
  tasks: Task[];
  onMilestoneChange: (id: number, milestone: string) => void;
  onCardClick?: (task: Task) => void;
  milestoneOrder: string[];
  onMilestoneOrderChange: (order: string[]) => void;
}

export const PlanView: React.FC<IPlanView> = (props) => {
  const [activeTask, setActiveTask] = useState<Task | null>(null);
  const [activeMilestone, setActiveMilestone] = useState<string | null>(null);

  const sensors = useSensors(
    useSensor(PointerSensor, { activationConstraint: { distance: 5 } })
  );

  const rawMilestones = useMemo(() => {
    const set = new Set<string>();
    for (const t of props.tasks) {
      set.add(t.milestone || "");
    }
    return [...set];
  }, [props.tasks]);

  const milestones = useMemo(
    () => orderMilestones(rawMilestones, props.milestoneOrder),
    [rawMilestones, props.milestoneOrder]
  );

  const sortableIds = useMemo(
    () => milestones.map((m) => `${MILESTONE_PREFIX}${m === "" ? NONE_SENTINEL : m}`),
    [milestones]
  );

  const handleDragStart = (event: DragStartEvent) => {
    const id = String(event.active.id);
    if (id.startsWith(MILESTONE_PREFIX)) {
      const ms = id.slice(MILESTONE_PREFIX.length);
      setActiveMilestone(ms === NONE_SENTINEL ? "" : ms);
    } else {
      const task = props.tasks.find((t) => t.id === event.active.id);
      if (task) setActiveTask(task);
    }
  };

  const handleDragEnd = (event: DragEndEvent) => {
    const { active, over } = event;
    setActiveTask(null);
    setActiveMilestone(null);
    if (!over) return;

    const activeId = String(active.id);

    if (activeId.startsWith(MILESTONE_PREFIX)) {
      // Milestone reorder
      const overId = String(over.id);
      if (activeId === overId) return;
      const oldIndex = sortableIds.indexOf(activeId);
      const newIndex = sortableIds.indexOf(overId);
      if (oldIndex === -1 || newIndex === -1) return;
      const newMilestones = arrayMove(milestones, oldIndex, newIndex);
      props.onMilestoneOrderChange(newMilestones);
    } else {
      // Card drag to milestone
      const taskId = active.id as number;
      const dropId = over.id as string;
      // Strip milestone prefix if dropping on a milestone header
      const rawDrop = dropId.startsWith(MILESTONE_PREFIX)
        ? dropId.slice(MILESTONE_PREFIX.length)
        : dropId;
      const newMilestone = rawDrop === NONE_SENTINEL ? "" : rawDrop;

      const task = props.tasks.find((t) => t.id === taskId);
      if (task && (task.milestone || "") !== newMilestone) {
        props.onMilestoneChange(taskId, newMilestone);
      }
    }
  };

  return (
    <DndContext
      sensors={sensors}
      onDragStart={handleDragStart}
      onDragEnd={handleDragEnd}
    >
      <div className={css.plan}>
        <SortableContext
          items={sortableIds}
          strategy={verticalListSortingStrategy}
        >
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
        </SortableContext>
      </div>

      <DragOverlay>
        {activeTask ? <PlanCard task={activeTask} isOverlay /> : null}
        {activeMilestone !== null ? (
          <div className={cx(css.group, css.groupOverlay)}>
            <div className={css.groupHeader}>
              <DotsSixVertical size={16} className={css.dragHandle} />
              <span className={css.dot} />
              <span className={css.groupLabel}>
                {activeMilestone === "" ? "No milestone" : activeMilestone}
              </span>
            </div>
          </div>
        ) : null}
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
  const sortableId = `${MILESTONE_PREFIX}${milestone === "" ? NONE_SENTINEL : milestone}`;
  const droppableId = milestone === "" ? NONE_SENTINEL : milestone;

  const {
    attributes,
    listeners,
    setNodeRef: setSortableNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({ id: sortableId });

  const { setNodeRef: setDroppableNodeRef, isOver } = useDroppable({
    id: droppableId,
  });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
  };

  return (
    <div
      ref={(node) => {
        setSortableNodeRef(node);
        setDroppableNodeRef(node);
      }}
      style={style}
      className={cx(css.group, {
        [css.groupOver]: isOver,
        [css.groupEmpty]: tasks.length === 0,
        [css.groupDragging]: isDragging,
      })}
    >
      <div className={css.groupHeader}>
        <span
          className={css.dragHandle}
          {...attributes}
          {...listeners}
        >
          <DotsSixVertical size={16} />
        </span>
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
