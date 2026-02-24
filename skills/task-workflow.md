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

### 3. Write a plan

Before writing code, write an implementation plan and persist it:

```
cat <<'EOF' | ghist task update <id> --plan-stdin
## Approach
- Step 1: ...
- Step 2: ...

## Files to change
- path/to/file.go — description of change
EOF
```

The plan is stored on the task and survives session boundaries. If a session
ends mid-task, the next agent can read the plan and pick up where you left off.

If you update the plan later, use `--plan-stdin` again — it overwrites the existing plan on the same task.

### 4. Execute

First, move the task to in_progress:

```
ghist task update <id> --status in_progress
```

Then work through the plan step by step. Log progress as you go:

```
ghist log "Completed step 1: added migration" --task <id>
ghist log "Step 2 blocked: need API review" --task <id>
```

### 5. Complete

Do not mark the task as done automatically. Instead:

1. Write a short implementation note summarising what was done:
   ```
   ghist log "Brief summary of what was implemented and any key decisions made" --type decision --task <id>
   ```

2. Link the commit if applicable:
   ```
   ghist task update <id> --commit-hash abc1234
   ```

3. Ask the user for confirmation:
   > "I think this is done — should I close the task?"

4. Only mark as done once confirmed:
   ```
   ghist task update <id> --status done
   ```

## Why Persist Plans?

- **Session continuity**: If the session ends, the next agent reads the plan instead of starting over.
- **Auditability**: Plans are stored alongside the task, so you can review what was intended vs. what was done.
- **Parallelism**: Multiple agents can read each other's plans to avoid conflicting work.
