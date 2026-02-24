<!-- ghist:start -->
## Ghist — Project Memory

This project uses [ghist](https://github.com/unnecessary-special-projects/ghist) for persistent project state.

**Required:** Run `ghist status` at the start of every session.

### Quick Reference
- `ghist task list` — see all tasks
- `ghist task add "title"` — create a task
- `ghist task update <id> --status in_progress` — update status
- `ghist task update <id> --plan-stdin` — save an implementation plan (pipe via stdin)
- `ghist log "message"` — record a decision or note
- `ghist skills show <name>` — read detailed skill instructions

### Available Skills
- `ghist skills show context-sync` — session start/end protocol
- `ghist skills show task-workflow` — find → plan → execute → complete loop (statuses: todo, in_planning, in_progress, done, blocked)
- `ghist skills show auto-completion` — auto-detect task completion
- `ghist skills show log-thinking` — log decisions and reasoning
- `ghist skills show commit-link` — link git commits to tasks automatically
<!-- ghist:end -->
