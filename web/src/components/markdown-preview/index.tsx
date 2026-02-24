import Markdown from "react-markdown";
import remarkGfm from "remark-gfm";
import type { Components } from "react-markdown";
import css from "./index.module.css";

export interface IMarkdownPreview {
  content: string;
}

const components: Components = {
  h1: ({ children }) => <h1 className={css.h1}>{children}</h1>,
  h2: ({ children }) => <h2 className={css.h2}>{children}</h2>,
  h3: ({ children }) => <h3 className={css.h3}>{children}</h3>,
  p: ({ children }) => <p className={css.p}>{children}</p>,
  a: ({ href, children }) => <a href={href} className={css.a} target="_blank" rel="noopener noreferrer">{children}</a>,
  code: ({ className, children }) => {
    const isBlock = className?.includes("language-");
    if (isBlock) {
      return (
        <pre className={css.codeBlock}>
          <code className={css.code}>{children}</code>
        </pre>
      );
    }
    return <code className={css.codeInline}>{children}</code>;
  },
  pre: ({ children }) => <>{children}</>,
  ul: ({ children }) => <ul className={css.ul}>{children}</ul>,
  ol: ({ children }) => <ol className={css.ol}>{children}</ol>,
  li: ({ children }) => <li className={css.li}>{children}</li>,
  table: ({ children }) => (
    <div className={css.tableWrap}>
      <table className={css.table}>{children}</table>
    </div>
  ),
  th: ({ children }) => <th className={css.th}>{children}</th>,
  td: ({ children }) => <td className={css.td}>{children}</td>,
  blockquote: ({ children }) => <blockquote className={css.blockquote}>{children}</blockquote>,
  hr: () => <hr className={css.hr} />,
};

export const MarkdownPreview: React.FC<IMarkdownPreview> = (props) => {
  return (
    <div className={css.root}>
      <Markdown remarkPlugins={[remarkGfm]} components={components}>
        {props.content}
      </Markdown>
    </div>
  );
};
