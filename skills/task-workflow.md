# Task Workflow

A structured workflow for picking up, planning, executing, and completing tasks.

## The Loop

### 1. Find the task

```
ghist task list
ghist task show <id>
```

Review the task details, description, and any existing plan.

### 2. Start planning

```
ghist task update <id> --status in_planning
```

This signals to other agents (and future sessions) that planning is underway.

### 3. Write and save the plan

Write an implementation plan and **save it to the task as soon as it's ready**:

```
cat <<'EOF' | ghist task update <id> --plan-stdin
## Approach
- Step 1: ...
- Step 2: ...

## Files to change
- path/to/file.go — description of change
EOF
```

**This is mandatory.** The plan must live on the task — not just in conversation context — so that any session can pick it up. If the user's agent (e.g. Claude Code) produces a plan through its own planning workflow, save that plan to the task using `--plan-stdin` as soon as it's finalized.

### 4. Execute

Move the task to in_progress:

```
ghist task update <id> --status in_progress
```

Work through the plan step by step. **Keep the plan current on the task** — update it on every meaningful change (new files, different approach, scope change, completed steps) so it always reflects the real state of the work:

```
cat <<'EOF' | ghist task update <id> --plan-stdin
## Approach (revised)
- Step 1: ... (done)
- Step 2: ... (done)
- Step 3: ... (added — discovered during implementation)

## Files changed
- path/to/file.go — what changed
- path/to/new_file.go — why it was added
EOF
```

### 5. Complete

When the work is done, **append implementation notes to the plan** on the task. This keeps everything about the task — what was planned, what was actually done, and why — in one place.

```
cat <<'EOF' | ghist task update <id> --plan-stdin
## Approach
- Step 1: Added migration for new table
- Step 2: Implemented API endpoint
- Step 3: Wrote tests

## Files changed
- internal/store/store.go — added migration
- internal/api/tasks.go — new endpoint
- internal/api/tasks_test.go — test coverage

## Implementation Notes
Went with JWT over sessions to keep the API stateless. Added tests for
token expiry and invalid signatures. Skipped refresh tokens for now —
can add in a follow-up if needed.
EOF
```

Then link the commit and ask the user before closing:

```
ghist task update <id> --commit-hash abc1234
```

> "I've finished the implementation and updated the task with notes. Should I close it?"

Only mark as done once confirmed:

```
ghist task update <id> --status done
```

## Rules

1. **Save plans to the task as soon as they're ready.** Never keep the plan only in conversation. Use `--plan-stdin`.
2. **Keep the plan current.** Update it on every meaningful change as work progresses.
3. **Write implementation notes on the task** when done — not just as a log event.
4. **Ask before closing.** The user decides when a task is done.

## Why Persist Plans?

- **Session continuity**: If the session ends, the next agent reads the plan instead of starting over.
- **Auditability**: Plans are stored alongside the task, so you can review what was intended vs. what was done.
- **Parallelism**: Multiple agents can read each other's plans to avoid conflicting work.
