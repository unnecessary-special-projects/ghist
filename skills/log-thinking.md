# Log Thinking

Record decisions, reasoning, and architectural notes to the event timeline.

## When to Log

Log whenever you make a decision that a future session should know about:

- **Architectural decisions** — "Chose SQLite over PostgreSQL because local-first"
- **Trade-offs** — "Skipped pagination for MVP, will add when task count > 100"
- **Blocked paths** — "Tried approach X, didn't work because Y"
- **Dependencies discovered** — "Task A depends on task B being done first"

## How to Log

Basic log entry:
```
ghist log "Chose REST over GraphQL for simplicity"
```

Log linked to a task:
```
ghist log "Decided to use @dnd-kit for drag-and-drop" --task 5
```

Log with a type for categorization:
```
ghist log "API rate limiting needed before public release" --type decision
ghist log "Bug: drag-and-drop doesn't work on mobile Safari" --type note
```

## Log Types

- `log` (default) — general notes
- `decision` — architectural or design decisions
- `note` — observations, bugs found, things to remember

## Why This Matters

AI agents don't remember previous sessions. Without logged decisions, the next
session might revisit the same trade-offs or make conflicting choices. Logging
your thinking creates a decision trail that any future agent can follow.

## Best Practices

- Be specific — "Chose X because Y" is better than "Made a decision"
- Link to tasks when the decision affects a specific task
- Log blocking discoveries immediately so they're not lost
- Keep entries concise — one decision per log entry
