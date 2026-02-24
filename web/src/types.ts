export type TaskStatus = 'todo' | 'in_planning' | 'in_progress' | 'done' | 'blocked';
export type TaskPriority = '' | 'low' | 'medium' | 'high' | 'urgent';
export type TaskType = '' | 'bug' | 'feature' | 'improvement' | 'chore';

export interface Task {
  id: number;
  title: string;
  description: string;
  plan: string;
  status: TaskStatus;
  milestone: string;
  commit_hash: string;
  priority: TaskPriority;
  type: TaskType;
  ref_id: string;
  legacy_id: string;
  created_at: string;
  updated_at: string;
}

export interface Event {
  id: number;
  type: string;
  message: string;
  metadata: string;
  task_id: number | null;
  created_at: string;
}

export interface StatusSummary {
  total_tasks: number;
  tasks_by_status: Record<string, number>;
  milestones: MilestoneInfo[];
  recent_events: Event[];
}

export interface MilestoneInfo {
  name: string;
  total: number;
  done: number;
}

export const STATUSES: TaskStatus[] = ['todo', 'in_planning', 'in_progress', 'done', 'blocked'];
export const LIST_STATUSES: TaskStatus[] = [ 'in_progress', 'in_planning', 'todo', 'done', 'blocked'];
export const BOARD_STATUSES: TaskStatus[] = ['todo', 'in_planning', 'in_progress', 'done', 'blocked'];

export const STATUS_LABELS: Record<TaskStatus, string> = {
  todo: 'To Do',
  in_planning: 'In Planning',
  in_progress: 'In Progress',
  done: 'Done',
  blocked: 'Blocked',
};

export const STATUS_COLORS: Record<TaskStatus, string> = {
  todo: '#768390',
  in_planning: '#a371f7',
  in_progress: '#d29922',
  done: '#238636',
  blocked: '#f85149',
};

export const PRIORITIES: TaskPriority[] = ['', 'low', 'medium', 'high', 'urgent'];

export const PRIORITY_LABELS: Record<TaskPriority, string> = {
  '': 'None',
  low: 'Low',
  medium: 'Medium',
  high: 'High',
  urgent: 'Urgent',
};

export const PRIORITY_COLORS: Record<TaskPriority, string> = {
  '': '#484f58',
  low: '#3fb950',
  medium: '#d29922',
  high: '#f85149',
  urgent: '#da3633',
};

export const TASK_TYPES: TaskType[] = ['', 'bug', 'feature', 'improvement', 'chore'];

export const TYPE_LABELS: Record<TaskType, string> = {
  '': 'None',
  bug: 'Bug',
  feature: 'Feature',
  improvement: 'Improvement',
  chore: 'Chore',
};

export const TYPE_COLORS: Record<TaskType, string> = {
  '': '#484f58',
  bug: '#f85149',
  feature: '#a371f7',
  improvement: '#58a6ff',
  chore: '#768390',
};
