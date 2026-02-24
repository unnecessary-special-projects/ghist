import cx from "classnames";
import { List, Columns } from "@phosphor-icons/react";
import css from "./index.module.css";
import type { TaskPriority, TaskType } from "../../types";
import { PRIORITIES, PRIORITY_LABELS, TASK_TYPES, TYPE_LABELS } from "../../types";
import type { ViewMode } from "../../hooks/useTaskFilters";

export interface IToolbar {
  viewMode: ViewMode;
  onViewModeChange: (mode: ViewMode) => void;
  priorityFilter: TaskPriority | "all";
  onPriorityFilterChange: (priority: TaskPriority | "all") => void;
  typeFilter: TaskType | "all";
  onTypeFilterChange: (type: TaskType | "all") => void;
  searchQuery: string;
  onSearchQueryChange: (query: string) => void;
}

export const Toolbar: React.FC<IToolbar> = (props) => {
  return (
    <div className={css.toolbar}>
      <div className={css.left}>
        <div className={css.segmented}>
          <button
            className={cx(css.segBtn, { [css.segBtnActive]: props.viewMode === "list" })}
            onClick={() => props.onViewModeChange("list")}
            title="List view"
          >
            <List size={18} weight="duotone" />
          </button>
          <button
            className={cx(css.segBtn, { [css.segBtnActive]: props.viewMode === "board" })}
            onClick={() => props.onViewModeChange("board")}
            title="Board view"
          >
            <Columns size={18} weight="duotone" />
          </button>
        </div>
      </div>

      <div className={css.right}>
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
  );
};
