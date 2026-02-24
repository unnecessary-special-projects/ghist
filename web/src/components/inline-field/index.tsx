import cx from "classnames";
import { useEffect, useRef, useState } from "react";
import css from "./index.module.css";

interface IOption {
  value: string;
  label: string;
  color?: string;
}

export interface IInlineField {
  label: string;
  value: string;
  onSave: (value: string) => void;
  type?: "text" | "textarea" | "select";
  options?: IOption[];
  renderValue?: (value: string) => React.ReactNode;
  readOnly?: boolean;
  className?: string;
}

export const InlineField: React.FC<IInlineField> = (props) => {
  const { label, value, onSave, type = "text", options, renderValue, readOnly, className } = props;
  const [editing, setEditing] = useState(false);
  const [draft, setDraft] = useState(value);
  const inputRef = useRef<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>(null);

  useEffect(() => {
    setDraft(value);
  }, [value]);

  useEffect(() => {
    if (editing && inputRef.current) {
      inputRef.current.focus();
    }
  }, [editing]);

  const save = () => {
    setEditing(false);
    if (draft !== value) {
      onSave(draft);
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (type === "text" && e.key === "Enter") {
      save();
    }
    if (e.key === "Escape") {
      setDraft(value);
      setEditing(false);
    }
  };

  if (readOnly) {
    return (
      <div className={cx(css.field, className)}>
        <span className={css.label}>{label}</span>
        <div className={css.value}>
          {renderValue ? renderValue(value) : value || <span className={css.empty}>&mdash;</span>}
        </div>
      </div>
    );
  }

  if (editing) {
    return (
      <div className={cx(css.field, className)}>
        <span className={css.label}>{label}</span>
        {type === "select" && options ? (
          <select
            ref={inputRef as React.RefObject<HTMLSelectElement>}
            className={css.input}
            value={draft}
            onChange={(e) => {
              setDraft(e.target.value);
              setEditing(false);
              if (e.target.value !== value) {
                onSave(e.target.value);
              }
            }}
            onBlur={save}
          >
            {options.map((o) => (
              <option key={o.value} value={o.value}>{o.label}</option>
            ))}
          </select>
        ) : type === "textarea" ? (
          <textarea
            ref={inputRef as React.RefObject<HTMLTextAreaElement>}
            className={`${css.input} ${css.inputTextarea}`}
            value={draft}
            onChange={(e) => setDraft(e.target.value)}
            onBlur={save}
            onKeyDown={handleKeyDown}
          />
        ) : (
          <input
            ref={inputRef as React.RefObject<HTMLInputElement>}
            className={css.input}
            value={draft}
            onChange={(e) => setDraft(e.target.value)}
            onBlur={save}
            onKeyDown={handleKeyDown}
          />
        )}
      </div>
    );
  }

  return (
    <div className={cx(css.field, className)}>
      <span className={css.label}>{label}</span>
      <div className={css.valueEditable} onClick={() => setEditing(true)}>
        {renderValue ? renderValue(value) : value || <span className={css.empty}>&mdash;</span>}
      </div>
    </div>
  );
};
