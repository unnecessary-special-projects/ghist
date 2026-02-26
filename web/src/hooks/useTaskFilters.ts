import { useCallback, useState } from 'react';
import type { Task, TaskPriority, TaskType } from '../types';

export type ViewMode = 'list' | 'board' | 'plan';
export type SortOption = 'newest' | 'updated' | 'priority' | 'title';

const PRIORITY_ORDER: Record<TaskPriority, number> = {
  urgent: 0,
  high: 1,
  medium: 2,
  low: 3,
  '': 4,
};

export function useTaskFilters() {
  const [viewMode, setViewMode] = useState<ViewMode>('list');
  const [priorityFilter, setPriorityFilter] = useState<TaskPriority | 'all'>('all');
  const [typeFilter, setTypeFilter] = useState<TaskType | 'all'>('all');
  const [searchQuery, setSearchQuery] = useState('');
  const [sortBy, setSortBy] = useState<SortOption>('newest');
  const [milestoneFilter, setMilestoneFilter] = useState<Set<string>>(new Set());

  const toggleMilestone = useCallback((milestone: string) => {
    setMilestoneFilter((prev) => {
      const next = new Set(prev);
      if (next.has(milestone)) {
        next.delete(milestone);
      } else {
        next.add(milestone);
      }
      return next;
    });
  }, []);

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
      if (milestoneFilter.size > 0) {
        result = result.filter((t) => {
          const val = t.milestone || '';
          return milestoneFilter.has(val);
        });
      }

      result = [...result].sort((a, b) => {
        switch (sortBy) {
          case 'newest':
            return b.created_at.localeCompare(a.created_at);
          case 'updated':
            return b.updated_at.localeCompare(a.updated_at);
          case 'priority':
            return (PRIORITY_ORDER[a.priority] ?? 4) - (PRIORITY_ORDER[b.priority] ?? 4);
          case 'title':
            return a.title.localeCompare(b.title);
          default:
            return 0;
        }
      });

      return result;
    },
    [priorityFilter, typeFilter, searchQuery, sortBy, milestoneFilter]
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
    sortBy,
    setSortBy,
    milestoneFilter,
    toggleMilestone,
    filterTasks,
  };
}
