# Auto-Completion

Detect when a task is finished, write implementation notes on the task, and ask the user before closing it.

## When to Consider a Task Done

Look for these signals that a task might be complete:

1. **All acceptance criteria met** — if the task description lists specific criteria, verify each one.
2. **Tests passing** — if the task involves code, confirm tests pass.
3. **No remaining blockers** — verify nothing is left unresolved.

Never assume a task is done. The user decides when to close it.

## The Completion Flow

### 1. Update the task plan with implementation notes

Before asking the user, **update the plan on the task** to include a summary of what was done. This keeps the full history — plan and outcome — on the task itself, not just in the event log.

Read the current plan, append an "Implementation Notes" section, and save it back:

```
cat <<'EOF' | ghist task update <id> --plan-stdin
<existing plan content here>

## Implementation Notes
- Implemented JWT middleware in internal/auth
- Chose stateless tokens over sessions to keep the API simple
- Added tests for expiry and invalid signature cases
- Skipped refresh tokens — can add later if needed
EOF
```

**This is mandatory.** Every completed task must have implementation notes saved to the task's plan field.

### 2. Link the commit

```
ghist task update <id> --commit-hash abc1234
```

### 3. Ask the user

Do not mark the task as done automatically. Instead, ask:

> "I've finished the implementation and updated the task with notes. Should I close it?"

Wait for confirmation before proceeding.

### 4. Close only on confirmation

Once the user confirms:

```
ghist task update <id> --status done
```

Optionally log a brief event for the timeline:

```
ghist log "Completed: <one-line summary>" --task <id>
```

## Partial Completion

If a task is only partially done, update the plan to reflect progress and keep it in progress:

```
cat <<'EOF' | ghist task update <id> --plan-stdin
<existing plan, with completed steps marked>

## Progress Notes
- Completed steps 1-3
- Step 4 remaining: needs API review before proceeding
EOF
```

Do not ask the user to close a partially complete task.
