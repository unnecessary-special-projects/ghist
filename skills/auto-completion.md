# Auto-Completion

Detect when a task is finished, write an implementation note, and ask the user before closing it.

## When to Consider a Task Done

Look for these signals that a task might be complete:

1. **All acceptance criteria met** — if the task description lists specific criteria, verify each one.
2. **Tests passing** — if the task involves code, confirm tests pass.
3. **No remaining blockers** — verify nothing is left unresolved.

Never assume a task is done. The user decides when to close it.

## The Completion Flow

### 1. Write an implementation note

Before asking the user, log a short summary of what was done — what changed, what approach was taken, and any relevant trade-offs. Keep it brief.

```
ghist log "Implemented JWT middleware in internal/auth. Chose stateless tokens over sessions to keep the API simple. Added tests for expiry and invalid signature cases." --type decision --task <id>
```

### 2. Link the commit if applicable

```
ghist task update <id> --commit-hash abc1234
```

### 3. Ask the user

Do not mark the task as done automatically. Instead, ask:

> "I think this is done — the implementation is complete and tests are passing. Should I close the task?"

Wait for confirmation before proceeding.

### 4. Close only on confirmation

Once the user confirms:

```
ghist task update <id> --status done
```

## Partial Completion

If a task is only partially done, log what was accomplished and keep it in progress:

```
ghist log "Completed steps 1-3, remaining: step 4 (needs API review)" --task <id>
```

Do not ask the user to close a partially complete task.
