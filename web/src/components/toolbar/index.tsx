import { useMemo } from "react";
import cx from "classnames";
import {
  DndContext,
  DragEndEvent,
  PointerSensor,
  useSensor,
  useSensors,
} from "@dnd-kit/core";
import {
  SortableContext,
  useSortable,
  arrayMove,
  horizontalListSortingStrategy,
} from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import css from "./index.module.css";
import type { Task, TaskPriority, TaskType } from "../../types";
import { PRIORITIES, PRIORITY_LABELS, TASK_TYPES, TYPE_LABELS } from "../../types";
import type { ViewMode, SortOption } from "../../hooks/useTaskFilters";
import { orderMilestones } from "../../utils/milestoneOrder";

const NONE_SENTINEL = "__none__";

export interface IToolbar {
  viewMode: ViewMode;
  onViewModeChange: (mode: ViewMode) => void;
  priorityFilter: TaskPriority | "all";
  onPriorityFilterChange: (priority: TaskPriority | "all") => void;
  typeFilter: TaskType | "all";
  onTypeFilterChange: (type: TaskType | "all") => void;
  searchQuery: string;
  onSearchQueryChange: (query: string) => void;
  sortBy: SortOption;
  onSortChange: (sort: SortOption) => void;
  tasks: Task[];
  milestoneFilter: Set<string>;
  onToggleMilestone: (milestone: string) => void;
  milestoneOrder: string[];
  onMilestoneOrderChange: (order: string[]) => void;
}

const SORT_LABELS: Record<SortOption, string> = {
  newest: "Newest First",
  updated: "Recently Updated",
  priority: "Priority",
  title: "Title A\u2013Z",
};

export const Toolbar: React.FC<IToolbar> = (props) => {
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
    () => milestones.map((m) => (m === "" ? NONE_SENTINEL : m)),
    [milestones]
  );

  const showMilestones = milestones.length > 1 || (milestones.length === 1 && milestones[0] !== "");

  const handleDragEnd = (event: DragEndEvent) => {
    const { active, over } = event;
    if (!over || active.id === over.id) return;
    const oldIndex = sortableIds.indexOf(String(active.id));
    const newIndex = sortableIds.indexOf(String(over.id));
    if (oldIndex === -1 || newIndex === -1) return;
    const newMilestones = arrayMove(milestones, oldIndex, newIndex);
    props.onMilestoneOrderChange(newMilestones);
  };

  return (
    <div>
      <div className={css.toolbar}>
        <div className={css.left}>
          <div className={css.segmented}>
            <button
              className={cx(css.segBtn, { [css.segBtnActive]: props.viewMode === "plan" })}
              onClick={() => props.onViewModeChange("plan")}
            >
              Plan
            </button>
            <button
              className={cx(css.segBtn, { [css.segBtnActive]: props.viewMode === "list" })}
              onClick={() => props.onViewModeChange("list")}
            >
              List
            </button>
            <button
              className={cx(css.segBtn, { [css.segBtnActive]: props.viewMode === "board" })}
              onClick={() => props.onViewModeChange("board")}
            >
              Board
            </button>
          </div>
        </div>

        <div className={css.right}>
          <select
            className={css.select}
            value={props.sortBy}
            onChange={(e) => props.onSortChange(e.target.value as SortOption)}
          >
            {(Object.keys(SORT_LABELS) as SortOption[]).map((key) => (
              <option key={key} value={key}>{SORT_LABELS[key]}</option>
            ))}
          </select>

          <select
            className={css.select}
            value={props.priorityFilter}
            onChange={(e) => props.onPriorityFilterChange(e.target.value as TaskPriority | "all")}
          >
            <option value="all">All Priorities</option>
            {PRIORITIES.filter((p) => p !== "").map((p) => (
              <option key={p} value={p}>{PRIORITY_LABELS[p]}</option>
            ))}
          </select>

          <select
            className={css.select}
            value={props.typeFilter}
            onChange={(e) => props.onTypeFilterChange(e.target.value as TaskType | "all")}
          >
            <option value="all">All Types</option>
            {TASK_TYPES.filter((t) => t !== "").map((t) => (
              <option key={t} value={t}>{TYPE_LABELS[t]}</option>
            ))}
          </select>

          <input
            className={css.search}
            type="text"
            placeholder="Search tasks..."
            value={props.searchQuery}
            onChange={(e) => props.onSearchQueryChange(e.target.value)}
          />
        </div>
      </div>

      {showMilestones && (
        <DndContext sensors={sensors} onDragEnd={handleDragEnd}>
          <SortableContext
            items={sortableIds}
            strategy={horizontalListSortingStrategy}
          >
            <div className={css.milestoneBar}>
              {milestones.map((m) => (
                <SortableChip
                  key={m === "" ? NONE_SENTINEL : m}
                  id={m === "" ? NONE_SENTINEL : m}
                  label={m === "" ? "No milestone" : m}
                  active={props.milestoneFilter.has(m)}
                  onClick={() => props.onToggleMilestone(m)}
                />
              ))}
            </div>
          </SortableContext>
        </DndContext>
      )}
    </div>
  );
};

function SortableChip({
  id,
  label,
  active,
  onClick,
}: {
  id: string;
  label: string;
  active: boolean;
  onClick: () => void;
}) {
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({ id });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging ? 0.4 : undefined,
  };

  return (
    <button
      ref={setNodeRef}
      style={style}
      {...attributes}
      {...listeners}
      className={cx(css.chip, { [css.chipActive]: active })}
      onClick={onClick}
    >
      {label}
    </button>
  );
}
