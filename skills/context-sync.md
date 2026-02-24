# Context Sync

Synchronize project context at the start and end of every AI agent session.

## Session Start Protocol

1. Run `ghist status` to get a snapshot of the current project state.
2. Review the task list — identify tasks that are `in_planning`, `in_progress`, or `blocked`.
3. Check recent events for decisions or notes from previous sessions.
4. If resuming work on a task, update its status to `in_progress`:
   ```
   ghist task update <id> --status in_progress
   ```

## Session End Protocol

1. Log a summary of what was accomplished:
   ```
   ghist log "Completed X, started Y, blocked on Z"
   ```
2. Update task statuses to reflect current state:
   ```
   ghist task update <id> --status done
   ghist task update <id> --status blocked
   ```
3. If new tasks were discovered during the session, add them:
   ```
   ghist task add "New task title" --description "Details"
   ```
4. Run `ghist status` one final time to verify the project state is accurate.

## Why This Matters

Without context sync, the next session starts from scratch. By consistently
running these commands, you ensure continuity across sessions — every agent
that picks up the project knows exactly where things stand.
