import cx from "classnames";
import { useDraggable } from "@dnd-kit/core";
import css from "./index.module.css";
import type { Task } from "../../types";
import { PRIORITY_COLORS, PRIORITY_LABELS, TYPE_COLORS, TYPE_LABELS } from "../../types";

export interface ITaskCard {
  task: Task;
  isOverlay?: boolean;
  onCardClick?: (task: Task) => void;
  compact?: boolean;
}

export const TaskCard: React.FC<ITaskCard> = (props) => {

  const { attributes, listeners, setNodeRef, isDragging } = useDraggable({
    id: props.task.id,
  });

  return (
    <div
      ref={ setNodeRef }
      { ...listeners }
      { ...attributes }
      onClick={ () => props.onCardClick?.(props.task) }
      className={ cx(
        css.card, 
        {
          [css.cardCompact]: props.compact,
          [css.dragging]: isDragging && !props.isOverlay,
          [css.overlay]: props.isOverlay,
        }
      )}>
      
      <div className={css.titleRow}>
        <span className={css.title}>{props.task.title}</span>
      </div>

      {!props.compact && props.task.description && (
        <p className={css.description}>{props.task.description}</p>
      )}

      {!props.compact && (
        <div className={css.pills}>
          <span className={css.refId}>{props.task.ref_id}</span>
          <div className={css.pillGroup}>
            {props.task.priority && (
              <span
                className={css.priorityPill}
                style={{ color: PRIORITY_COLORS[props.task.priority], backgroundColor: `${PRIORITY_COLORS[props.task.priority]}1a` }}
              >
                <span className={css.pillLabel}>Priority:</span> {PRIORITY_LABELS[props.task.priority]}
              </span>
            )}
            {props.task.type && (
              <span
                className={css.typeBadge}
                style={{ color: TYPE_COLORS[props.task.type], backgroundColor: `${TYPE_COLORS[props.task.type]}1a` }}
              >
                <span className={css.pillLabel}>Type:</span> {TYPE_LABELS[props.task.type]}
              </span>
            )}
            {props.task.milestone && (
              <span className={css.milestone}>
                <span className={css.pillLabel}>Milestone:</span> {props.task.milestone}
              </span>
            )}
            {props.task.plan && (
              <span className={css.planBadge}>Plan</span>
            )}
          </div>
        </div>
      )}
    </div>
  );
};

