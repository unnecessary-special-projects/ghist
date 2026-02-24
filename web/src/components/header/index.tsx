import cx from "classnames";
import css from "./index.module.css";

export type AppView = "tasks" | "activity";

export interface IHeader {
  view: AppView;
  onViewChange: (view: AppView) => void;
  onNewTask: () => void;
}

export const Header: React.FC<IHeader> = (props) => {
  return (
    <header className={css.header}>
      <div className={css.left}>
        <img src="/ghist-logo.png" alt="ghist" className={css.logo} />
      </div>
      <div className={css.center}>
        <div className={css.nav}>
          <button
            className={cx(css.navItem, { [css.navItemActive]: props.view === "tasks" })}
            onClick={() => props.onViewChange("tasks")}
          >
            Tasks
          </button>
          <button
            className={cx(css.navItem, { [css.navItemActive]: props.view === "activity" })}
            onClick={() => props.onViewChange("activity")}
          >
            Activity
          </button>
        </div>
      </div>
      <div className={css.right}>
        {props.view === "tasks" && (
          <button className={css.addButton} onClick={props.onNewTask}>
            + New Task
          </button>
        )}
      </div>
    </header>
  );
};
