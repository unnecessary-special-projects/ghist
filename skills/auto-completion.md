# Auto-Completion

Automatically detect when tasks are completed and update their status.

## How It Works

After completing a piece of work, check if any tasks match what was just done.
If a task's requirements are fully met, mark it as done:

```
ghist task update <id> --status done
```

## Detection Signals

Look for these signals that a task might be complete:

1. **All acceptance criteria met** — if the task description lists specific criteria, verify each one.
2. **Tests passing** — if the task involves code, confirm tests pass.
3. **Code committed** — link the commit hash to the task:
   ```
   ghist task update <id> --status done --commit-hash abc1234
   ```
4. **No remaining blockers** — if the task was blocked, verify the blocker is resolved.

## Partial Completion

If a task is partially done, don't mark it as complete. Instead:

1. Log what was accomplished:
   ```
   ghist log "Completed steps 1-3 of task" --task <id>
   ```
2. Keep the task as `in_progress`.
3. Optionally update the description with remaining work:
   ```
   ghist task update <id> --description "Remaining: step 4, step 5"
   ```

## Best Practices

- Always verify completion before marking done — false completions create confusion.
- Link commits to tasks whenever possible for traceability.
- Log completion events so future sessions can see what was accomplished and when.
