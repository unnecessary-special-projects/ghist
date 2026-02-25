package store

import (
	"testing"
)

func newTestStore(t *testing.T) *Store {
	t.Helper()
	s, err := Open(t.TempDir())
	if err != nil {
		t.Fatalf("opening test store: %v", err)
	}
	t.Cleanup(func() { s.Close() })
	return s
}

// --- Task tests ---

func TestCreateAndGetTask(t *testing.T) {
	s := newTestStore(t)
	task, err := s.CreateTask(CreateTaskInput{Title: "Test task", Description: "A description", Milestone: "v1"})
	if err != nil {
		t.Fatalf("creating task: %v", err)
	}
	if task.ID != 1 {
		t.Errorf("expected id 1, got %d", task.ID)
	}
	if task.Title != "Test task" {
		t.Errorf("expected title 'Test task', got %q", task.Title)
	}
	if task.Status != "todo" {
		t.Errorf("expected status 'todo', got %q", task.Status)
	}
	if task.Milestone != "v1" {
		t.Errorf("expected milestone 'v1', got %q", task.Milestone)
	}
	if task.RefID != "GHST-1" {
		t.Errorf("expected ref_id 'GHST-1', got %q", task.RefID)
	}

	got, err := s.GetTask(task.ID)
	if err != nil {
		t.Fatalf("getting task: %v", err)
	}
	if got.Title != "Test task" {
		t.Errorf("expected title 'Test task', got %q", got.Title)
	}
}

func TestListTasks(t *testing.T) {
	s := newTestStore(t)
	s.CreateTask(CreateTaskInput{Title: "Task 1", Status: "todo"})
	s.CreateTask(CreateTaskInput{Title: "Task 2", Status: "in_progress", Milestone: "v1"})
	s.CreateTask(CreateTaskInput{Title: "Task 3", Status: "done", Milestone: "v1"})

	// List all
	tasks, err := s.ListTasks("", "", "", "")
	if err != nil {
		t.Fatalf("listing tasks: %v", err)
	}
	if len(tasks) != 3 {
		t.Errorf("expected 3 tasks, got %d", len(tasks))
	}

	// Filter by status
	tasks, err = s.ListTasks("in_progress", "", "", "")
	if err != nil {
		t.Fatalf("listing tasks: %v", err)
	}
	if len(tasks) != 1 {
		t.Errorf("expected 1 task, got %d", len(tasks))
	}

	// Filter by milestone
	tasks, err = s.ListTasks("", "v1", "", "")
	if err != nil {
		t.Fatalf("listing tasks: %v", err)
	}
	if len(tasks) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(tasks))
	}
}

func TestUpdateTask(t *testing.T) {
	s := newTestStore(t)
	s.CreateTask(CreateTaskInput{Title: "Original"})

	title := "Updated"
	status := "in_progress"
	task, err := s.UpdateTask(1, TaskUpdate{Title: &title, Status: &status})
	if err != nil {
		t.Fatalf("updating task: %v", err)
	}
	if task.Title != "Updated" {
		t.Errorf("expected title 'Updated', got %q", task.Title)
	}
	if task.Status != "in_progress" {
		t.Errorf("expected status 'in_progress', got %q", task.Status)
	}
}

func TestUpdateTaskPlan(t *testing.T) {
	s := newTestStore(t)
	s.CreateTask(CreateTaskInput{Title: "Plan test"})

	plan := "## Steps\n1. Do thing A\n2. Do thing B"
	task, err := s.UpdateTask(1, TaskUpdate{Plan: &plan})
	if err != nil {
		t.Fatalf("updating task plan: %v", err)
	}
	if task.Plan != plan {
		t.Errorf("expected plan %q, got %q", plan, task.Plan)
	}

	// Verify it persists via GetTask
	got, err := s.GetTask(1)
	if err != nil {
		t.Fatalf("getting task: %v", err)
	}
	if got.Plan != plan {
		t.Errorf("expected plan %q after get, got %q", plan, got.Plan)
	}

	// Verify plan shows in ListTasks
	tasks, err := s.ListTasks("", "", "", "")
	if err != nil {
		t.Fatalf("listing tasks: %v", err)
	}
	if len(tasks) != 1 || tasks[0].Plan != plan {
		t.Errorf("expected plan in list, got %q", tasks[0].Plan)
	}
}

func TestUpdateTaskNotFound(t *testing.T) {
	s := newTestStore(t)
	title := "Nope"
	_, err := s.UpdateTask(999, TaskUpdate{Title: &title})
	if err == nil {
		t.Fatal("expected error for non-existent task")
	}
}

func TestDeleteTask(t *testing.T) {
	s := newTestStore(t)
	s.CreateTask(CreateTaskInput{Title: "To delete"})

	err := s.DeleteTask(1)
	if err != nil {
		t.Fatalf("deleting task: %v", err)
	}

	_, err = s.GetTask(1)
	if err == nil {
		t.Fatal("expected error getting deleted task")
	}
}

func TestDeleteTaskNotFound(t *testing.T) {
	s := newTestStore(t)
	err := s.DeleteTask(999)
	if err == nil {
		t.Fatal("expected error for non-existent task")
	}
}

func TestTaskCountsByStatus(t *testing.T) {
	s := newTestStore(t)
	s.CreateTask(CreateTaskInput{Title: "T1", Status: "todo"})
	s.CreateTask(CreateTaskInput{Title: "T2", Status: "todo"})
	s.CreateTask(CreateTaskInput{Title: "T3", Status: "done"})

	counts, err := s.TaskCountsByStatus()
	if err != nil {
		t.Fatalf("counting: %v", err)
	}
	if counts["todo"] != 2 {
		t.Errorf("expected 2 todo, got %d", counts["todo"])
	}
	if counts["done"] != 1 {
		t.Errorf("expected 1 done, got %d", counts["done"])
	}
}

func TestMilestoneInfo(t *testing.T) {
	s := newTestStore(t)
	s.CreateTask(CreateTaskInput{Title: "T1", Status: "todo", Milestone: "v1"})
	s.CreateTask(CreateTaskInput{Title: "T2", Status: "done", Milestone: "v1"})
	s.CreateTask(CreateTaskInput{Title: "T3", Status: "todo", Milestone: "v2"})

	milestones, err := s.MilestoneInfo()
	if err != nil {
		t.Fatalf("querying milestones: %v", err)
	}
	if len(milestones) != 2 {
		t.Fatalf("expected 2 milestones, got %d", len(milestones))
	}
	if milestones[0].Name != "v1" || milestones[0].Total != 2 || milestones[0].Done != 1 {
		t.Errorf("unexpected v1 milestone: %+v", milestones[0])
	}
}

func TestTaskNewFields(t *testing.T) {
	s := newTestStore(t)
	task, err := s.CreateTask(CreateTaskInput{Title: "Fields test", Priority: "high", Type: "bug"})
	if err != nil {
		t.Fatalf("creating task with new fields: %v", err)
	}
	if task.Priority != "high" {
		t.Errorf("expected priority 'high', got %q", task.Priority)
	}
	if task.Type != "bug" {
		t.Errorf("expected type 'bug', got %q", task.Type)
	}
	if task.RefID != "GHST-1" {
		t.Errorf("expected ref_id 'GHST-1', got %q", task.RefID)
	}
}

func TestInPlanningStatus(t *testing.T) {
	s := newTestStore(t)
	task, err := s.CreateTask(CreateTaskInput{Title: "Planning task", Status: "in_planning"})
	if err != nil {
		t.Fatalf("creating in_planning task: %v", err)
	}
	if task.Status != "in_planning" {
		t.Errorf("expected status 'in_planning', got %q", task.Status)
	}
}

func TestPriorityAndTypeFiltering(t *testing.T) {
	s := newTestStore(t)
	s.CreateTask(CreateTaskInput{Title: "High bug", Priority: "high", Type: "bug"})
	s.CreateTask(CreateTaskInput{Title: "Low feature", Priority: "low", Type: "feature"})
	s.CreateTask(CreateTaskInput{Title: "High feature", Priority: "high", Type: "feature"})

	// Filter by priority
	tasks, err := s.ListTasks("", "", "high", "")
	if err != nil {
		t.Fatalf("listing tasks by priority: %v", err)
	}
	if len(tasks) != 2 {
		t.Errorf("expected 2 high-priority tasks, got %d", len(tasks))
	}

	// Filter by type
	tasks, err = s.ListTasks("", "", "", "feature")
	if err != nil {
		t.Fatalf("listing tasks by type: %v", err)
	}
	if len(tasks) != 2 {
		t.Errorf("expected 2 feature tasks, got %d", len(tasks))
	}

	// Filter by both
	tasks, err = s.ListTasks("", "", "high", "bug")
	if err != nil {
		t.Fatalf("listing tasks by priority+type: %v", err)
	}
	if len(tasks) != 1 {
		t.Errorf("expected 1 high-priority bug, got %d", len(tasks))
	}
}

func TestUpdateTaskPriorityAndType(t *testing.T) {
	s := newTestStore(t)
	s.CreateTask(CreateTaskInput{Title: "Update me"})

	priority := "urgent"
	taskType := "chore"
	legacyID := "JIRA-123"
	task, err := s.UpdateTask(1, TaskUpdate{Priority: &priority, Type: &taskType, LegacyID: &legacyID})
	if err != nil {
		t.Fatalf("updating task: %v", err)
	}
	if task.Priority != "urgent" {
		t.Errorf("expected priority 'urgent', got %q", task.Priority)
	}
	if task.Type != "chore" {
		t.Errorf("expected type 'chore', got %q", task.Type)
	}
	if task.LegacyID != "JIRA-123" {
		t.Errorf("expected legacy_id 'JIRA-123', got %q", task.LegacyID)
	}
}

func TestRefIDAutoGeneration(t *testing.T) {
	s := newTestStore(t)
	t1, _ := s.CreateTask(CreateTaskInput{Title: "First"})
	t2, _ := s.CreateTask(CreateTaskInput{Title: "Second"})
	t3, _ := s.CreateTask(CreateTaskInput{Title: "Third"})

	if t1.RefID != "GHST-1" {
		t.Errorf("expected GHST-1, got %q", t1.RefID)
	}
	if t2.RefID != "GHST-2" {
		t.Errorf("expected GHST-2, got %q", t2.RefID)
	}
	if t3.RefID != "GHST-3" {
		t.Errorf("expected GHST-3, got %q", t3.RefID)
	}
}

// --- Event tests ---

func TestCreateAndGetEvent(t *testing.T) {
	s := newTestStore(t)
	event, err := s.CreateEvent("log", "Something happened", "{}", nil)
	if err != nil {
		t.Fatalf("creating event: %v", err)
	}
	if event.ID != 1 {
		t.Errorf("expected id 1, got %d", event.ID)
	}
	if event.Message != "Something happened" {
		t.Errorf("unexpected message: %q", event.Message)
	}
	if event.TaskID != nil {
		t.Errorf("expected nil task_id, got %v", event.TaskID)
	}
}

func TestEventWithTask(t *testing.T) {
	s := newTestStore(t)
	task, _ := s.CreateTask(CreateTaskInput{Title: "A task"})
	taskID := task.ID

	event, err := s.CreateEvent("log", "Linked event", "{}", &taskID)
	if err != nil {
		t.Fatalf("creating event: %v", err)
	}
	if event.TaskID == nil || *event.TaskID != taskID {
		t.Errorf("expected task_id %d, got %v", taskID, event.TaskID)
	}
}

func TestListEvents(t *testing.T) {
	s := newTestStore(t)
	s.CreateEvent("log", "First", "{}", nil)
	s.CreateEvent("log", "Second", "{}", nil)
	s.CreateEvent("log", "Third", "{}", nil)

	events, err := s.ListEvents(2)
	if err != nil {
		t.Fatalf("listing events: %v", err)
	}
	if len(events) != 2 {
		t.Errorf("expected 2 events, got %d", len(events))
	}
}

func TestListEventsByTask(t *testing.T) {
	s := newTestStore(t)
	task, _ := s.CreateTask(CreateTaskInput{Title: "Task"})
	taskID := task.ID

	s.CreateEvent("log", "Linked", "{}", &taskID)
	s.CreateEvent("log", "Unlinked", "{}", nil)

	events, err := s.ListEventsByTask(taskID)
	if err != nil {
		t.Fatalf("listing events: %v", err)
	}
	if len(events) != 1 {
		t.Errorf("expected 1 event, got %d", len(events))
	}
}

// --- Opportunity tests ---

func TestCreateAndGetOpportunity(t *testing.T) {
	s := newTestStore(t)
	opp, err := s.CreateOpportunity("Feature idea", "Some notes")
	if err != nil {
		t.Fatalf("creating opportunity: %v", err)
	}
	if opp.Name != "Feature idea" {
		t.Errorf("unexpected name: %q", opp.Name)
	}

	got, err := s.GetOpportunity(opp.ID)
	if err != nil {
		t.Fatalf("getting opportunity: %v", err)
	}
	if got.Notes != "Some notes" {
		t.Errorf("unexpected notes: %q", got.Notes)
	}
}

func TestListOpportunities(t *testing.T) {
	s := newTestStore(t)
	s.CreateOpportunity("Opp 1", "")
	s.CreateOpportunity("Opp 2", "")

	opps, err := s.ListOpportunities()
	if err != nil {
		t.Fatalf("listing opportunities: %v", err)
	}
	if len(opps) != 2 {
		t.Errorf("expected 2 opportunities, got %d", len(opps))
	}
}

// --- Delete task cascades to events ---

func TestDeleteTaskSetsEventTaskNull(t *testing.T) {
	s := newTestStore(t)
	task, _ := s.CreateTask(CreateTaskInput{Title: "Task"})
	taskID := task.ID
	event, _ := s.CreateEvent("log", "Linked", "{}", &taskID)

	s.DeleteTask(taskID)

	got, err := s.GetEvent(event.ID)
	if err != nil {
		t.Fatalf("getting event after task delete: %v", err)
	}
	if got.TaskID != nil {
		t.Errorf("expected nil task_id after delete, got %v", got.TaskID)
	}
}
