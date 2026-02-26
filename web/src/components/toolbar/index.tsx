import { useMemo } from "react";
import cx from "classnames";
import css from "./index.module.css";
import type { Task, TaskPriority, TaskType } from "../../types";
import { PRIORITIES, PRIORITY_LABELS, TASK_TYPES, TYPE_LABELS } from "../../types";
import type { ViewMode, SortOption } from "../../hooks/useTaskFilters";

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
}

const SORT_LABELS: Record<SortOption, string> = {
  newest: "Newest First",
  updated: "Recently Updated",
  priority: "Priority",
  title: "Title A\u2013Z",
};

export const Toolbar: React.FC<IToolbar> = (props) => {
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

  const showMilestones = milestones.length > 1 || (milestones.length === 1 && milestones[0] !== "");

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
        <div className={css.milestoneBar}>
          {milestones.map((m) => {
            const active = props.milestoneFilter.has(m);
            return (
              <button
                key={m === "" ? "__none__" : m}
                className={cx(css.chip, { [css.chipActive]: active })}
                onClick={() => props.onToggleMilestone(m)}
              >
                {m === "" ? "No milestone" : m}
              </button>
            );
          })}
        </div>
      )}
    </div>
  );
};
