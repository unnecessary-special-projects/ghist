Ghist (The Agentic PM)
Role: Senior Software Architect
Project: Ghist (Open-Source, Local-First AI Project Manager)

1. The Core Mission
Ghist isn't just a CRM; it is a Project Memory Layer. It solves the problem of AI agents (Claude, Cursor) losing track of the "Why" and "What Next" across long sessions. It adds a persistent, queryable state to any repo via a CLI that the agent handles.

2. The "Skills" System (Crucial)
The ghist init command must generate a GHIST_SKILLS.md (or append to CLAUDE.md / .cursorrules). These skills are not just "how to use the CLI," but Autonomous Behaviors:

A. The "Context Sync" Skill
Trigger: Every time a new chat session starts.

Action: The agent must run ghist status.

Outcome: Ghist outputs a high-level summary: "You are currently 60% through the 'Auth' milestone. 2 tasks are blocked. Last meeting notes suggest focusing on the Login bug."

B. The "Auto-Completion" Skill (Diff-Checking)
Trigger: After the agent performs a git commit or a file write.

Action: The agent runs ghist scan --diff.

Mechanism: The CLI analyzes the staged changes. If the code implements a feature described in an "In Progress" task, the agent should automatically move that task to "Done" and log the commit hash in the Ghist DB.

C. The "Lead/Meeting" Skill
Trigger: When a developer mentions a person or a requirement in chat (e.g., "John wants the button blue").

Action: The agent recognizes this as a "Lead Update" or "Meeting Note" and runs ghist lead sync --note "Requested blue button".

3. Technical Requirements (Refined)
.ghist/ Internal Structure
ghist.sqlite: The source of truth.

current_context.json: A "fast-cache" of the active task so the AI doesn't have to query the whole DB every time.

CLI Enhancements
ghist log "message": A way for the AI to "think out loud" into the project history.

ghist plan: The AI looks at the Leads (Requirements) and Tasks (Execution) and proposes the next 3 logical steps.

4. Specific Instructions for Claude Code
When you start the repo, tell Claude Code:

"I want you to build the CLI for Ghist.

Key Logic to Implement First:

Automated Skill Injection: When I run ghist init, you must write a set of instructions to CLAUDE.md that forces you (the agent) to check Ghist before starting any work.

Reflective Updates: Implement a feature where you, the agent, can run ghist reflect. This command should look at my last 5 terminal commands and suggest which tasks in the database should be updated or created.

The Schema: We need tables for tasks (id, title, status, commit_hash), leads (id, name, pain_points, value_prop), and events (a global timeline of everything that happened in the repo)."
