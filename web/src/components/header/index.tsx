import css from "./index.module.css";

export interface IHeader {
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
          <button className={`${css.navItem} ${css.navItemActive}`}>Tasks</button>
        </div>
      </div>
      <div className={css.right}>
        <button className={css.addButton} onClick={props.onNewTask}>
          + New Task
        </button>
      </div>
    </header>
  );
};
