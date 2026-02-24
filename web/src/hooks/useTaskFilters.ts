import { useCallback, useState } from 'react';
import type { Task, TaskPriority, TaskType } from '../types';

export type ViewMode = 'list' | 'board';

export function useTaskFilters() {
  const [viewMode, setViewMode] = useState<ViewMode>('list');
  const [priorityFilter, setPriorityFilter] = useState<TaskPriority | 'all'>('all');
  const [typeFilter, setTypeFilter] = useState<TaskType | 'all'>('all');
  const [searchQuery, setSearchQuery] = useState('');

  const filterTasks = useCallback(
    (tasks: Task[]): Task[] => {
      let result = tasks;

      if (priorityFilter !== 'all') {
        result = result.filter((t) => t.priority === priorityFilter);
      }
      if (typeFilter !== 'all') {
        result = result.filter((t) => t.type === typeFilter);
      }
      if (searchQuery.trim()) {
        const q = searchQuery.toLowerCase();
        result = result.filter(
          (t) =>
            t.title.toLowerCase().includes(q) ||
            t.description.toLowerCase().includes(q) ||
            t.plan.toLowerCase().includes(q)
        );
      }

      return result;
    },
    [priorityFilter, typeFilter, searchQuery]
  );

  return {
    viewMode,
    setViewMode,
    priorityFilter,
    setPriorityFilter,
    typeFilter,
    setTypeFilter,
    searchQuery,
    setSearchQuery,
    filterTasks,
  };
}
