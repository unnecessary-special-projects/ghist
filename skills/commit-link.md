# Commit Linking

When ghist notifies you that a git commit was detected, link it to the relevant in-progress task.

## On Commit Detection

When you receive a system message like:

> "ghist: commit abc1234 was just made..."

Follow these steps:

1. Check both in-progress and recently completed tasks:
   ```
   ghist task list --status in_progress
   ghist task list --status done
   ```
   A task may have already been closed during this session before the commit happened — check both.

2. Consider the commit in context — what was just implemented, and which task does it relate to? Use the `updated_at` timestamps on done tasks to identify ones closed recently.

3. If there is a clear match, link the commit:
   ```
   ghist task update <id> --commit-hash abc1234
   ```

4. If multiple tasks are covered by the commit, link all of them.

5. If no tasks clearly match, skip it. Do not guess.

## When Not to Link

- The commit is a minor or unrelated change (formatting, typo fix, dependency bump)
- You cannot confidently identify which task it belongs to
