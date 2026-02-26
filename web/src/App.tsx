import { useCallback, useEffect, useState } from 'react';
import css from './App.module.css';
import { List } from './components/list';
import { Board } from './components/board';
import { PlanView } from './components/plan-view';
import { TaskDrawer } from './components/task-drawer';
import { Header } from './components/header';
import type { AppView } from './components/header';
import { Toolbar } from './components/toolbar';
import { ActivityFeed } from './components/activity-feed';
import { useTaskFilters } from './hooks/useTaskFilters';
import { useSSE } from './hooks/useSSE';
import type { Task, TaskStatus, TaskPriority, TaskType } from './types';
import * as api from './api/client';

export function App() {
  const [appView, setAppView] = useState<AppView>('tasks');
  const [tasks, setTasks] = useState<Task[]>([]);
  const [drawerTask, setDrawerTask] = useState<Task | null>(null);
  const [drawerMode, setDrawerMode] = useState<'view' | 'create' | null>(null);
  const [repoURL, setRepoURL] = useState<string>('');
  const [milestoneOrder, setMilestoneOrder] = useState<string[]>([]);

  const {
    viewMode, setViewMode,
    priorityFilter, setPriorityFilter,
    typeFilter, setTypeFilter,
    searchQuery, setSearchQuery,
    sortBy, setSortBy,
    milestoneFilter, toggleMilestone,
    filterTasks,
  } = useTaskFilters();

  const loadTasks = useCallback(async () => {
    const data = await api.listTasks();
    setTasks(data);
  }, []);

  const loadMilestoneOrder = useCallback(async () => {
    const order = await api.getMilestoneOrder();
    setMilestoneOrder(order);
  }, []);

  const handleRefresh = useCallback(async () => {
    await Promise.all([loadTasks(), loadMilestoneOrder()]);
  }, [loadTasks, loadMilestoneOrder]);

  useEffect(() => {
    handleRefresh();
    api.getConfig().then((c) => setRepoURL(c.github_repo_url)).catch(() => {});
  }, [handleRefresh]);

  useSSE(handleRefresh);

  const handleMilestoneOrderChange = async (order: string[]) => {
    setMilestoneOrder(order); // optimistic
    try {
      await api.setMilestoneOrder(order);
    } catch {
      await loadMilestoneOrder(); // revert
    }
  };

  const handleStatusChange = async (id: number, status: TaskStatus) => {
    // Optimistic update
    setTasks((prev) => prev.map((t) => (t.id === id ? { ...t, status } : t)));
    if (drawerTask?.id === id) {
      setDrawerTask((prev) => prev ? { ...prev, status } : null);
    }
    try {
      await api.updateTask(id, { status });
      await loadTasks();
    } catch {
      await loadTasks(); // revert
    }
  };

  const handleMilestoneChange = async (id: number, milestone: string) => {
    setTasks((prev) => prev.map((t) => (t.id === id ? { ...t, milestone } : t)));
    if (drawerTask?.id === id) {
      setDrawerTask((prev) => prev ? { ...prev, milestone } : null);
    }
    try {
      await api.updateTask(id, { milestone });
      await loadTasks();
    } catch {
      await loadTasks();
    }
  };

  const handleCardClick = (task: Task) => {
    setDrawerTask(task);
    setDrawerMode('view');
  };

  const handleNewTask = () => {
    setDrawerTask(null);
    setDrawerMode('create');
  };

  const handleCloseDrawer = () => {
    setDrawerMode(null);
    setDrawerTask(null);
  };

  const handleFieldSave = async (id: number, data: Record<string, string>) => {
    // Optimistic update
    setTasks((prev) => prev.map((t) => (t.id === id ? { ...t, ...data } : t)));
    if (drawerTask?.id === id) {
      setDrawerTask((prev) => prev ? { ...prev, ...data } : null);
    }
    try {
      const updated = await api.updateTask(id, data);
      setTasks((prev) => prev.map((t) => (t.id === id ? updated : t)));
      if (drawerTask?.id === id) {
        setDrawerTask(updated);
      }
    } catch {
      await loadTasks(); // revert
    }
  };

  const handleCreateTask = async (data: {
    title: string;
    description: string;
    status: TaskStatus;
    milestone: string;
    priority: TaskPriority;
    type: TaskType;
  }) => {
    const newTask = await api.createTask(data);
    await loadTasks();
    setDrawerTask(newTask);
    setDrawerMode('view');
  };

  const handleDeleteTask = async (id: number) => {
    await api.deleteTask(id);
    if (drawerTask?.id === id) {
      setDrawerMode(null);
      setDrawerTask(null);
    }
    await loadTasks();
  };

  const filteredTasks = filterTasks(tasks);

  return (
    <div className={css.app}>
      <Header view={appView} onViewChange={setAppView} onNewTask={handleNewTask} />
      {appView === 'activity' ? (
        <ActivityFeed tasks={tasks} />
      ) : (
        <>
          <Toolbar
            viewMode={viewMode}
            onViewModeChange={setViewMode}
            priorityFilter={priorityFilter}
            onPriorityFilterChange={setPriorityFilter}
            typeFilter={typeFilter}
            onTypeFilterChange={setTypeFilter}
            searchQuery={searchQuery}
            onSearchQueryChange={setSearchQuery}
            sortBy={sortBy}
            onSortChange={setSortBy}
            tasks={tasks}
            milestoneFilter={milestoneFilter}
            onToggleMilestone={toggleMilestone}
            milestoneOrder={milestoneOrder}
            onMilestoneOrderChange={handleMilestoneOrderChange}
          />
          {viewMode === 'list' ? (
            <List
              tasks={filteredTasks}
              onStatusChange={handleStatusChange}
              onCardClick={handleCardClick}
            />
          ) : viewMode === 'board' ? (
            <Board
              tasks={filteredTasks}
              onStatusChange={handleStatusChange}
              onCardClick={handleCardClick}
            />
          ) : (
            <PlanView
              tasks={filteredTasks}
              onMilestoneChange={handleMilestoneChange}
              onCardClick={handleCardClick}
              milestoneOrder={milestoneOrder}
              onMilestoneOrderChange={handleMilestoneOrderChange}
            />
          )}
        </>
      )}
      <TaskDrawer
        task={drawerTask}
        mode={drawerMode}
        onClose={handleCloseDrawer}
        onUpdateTask={handleFieldSave}
        onCreateTask={handleCreateTask}
        onDeleteTask={handleDeleteTask}
        repoURL={repoURL}
      />
    </div>
  );
}

