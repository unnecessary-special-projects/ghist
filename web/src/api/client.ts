import type { Task, Event, StatusSummary } from '../types';

const BASE = '/api';

async function request<T>(path: string, options?: RequestInit): Promise<T> {
  const res = await fetch(`${BASE}${path}`, {
    headers: { 'Content-Type': 'application/json' },
    ...options,
  });
  if (!res.ok) {
    const body = await res.json().catch(() => ({ error: res.statusText }));
    throw new Error(body.error || res.statusText);
  }
  return res.json();
}

export async function listTasks(
  status?: string,
  milestone?: string,
  priority?: string,
  type?: string,
): Promise<Task[]> {
  const params = new URLSearchParams();
  if (status) params.set('status', status);
  if (milestone) params.set('milestone', milestone);
  if (priority) params.set('priority', priority);
  if (type) params.set('type', type);
  const qs = params.toString();
  return request<Task[]>(`/tasks${qs ? `?${qs}` : ''}`);
}

export async function getTask(id: number): Promise<Task> {
  return request<Task>(`/tasks/${id}`);
}

export async function createTask(data: {
  title: string;
  description?: string;
  status?: string;
  milestone?: string;
  priority?: string;
  type?: string;
  legacy_id?: string;
}): Promise<Task> {
  return request<Task>('/tasks', {
    method: 'POST',
    body: JSON.stringify(data),
  });
}

export async function updateTask(
  id: number,
  data: Partial<Pick<Task, 'title' | 'description' | 'plan' | 'status' | 'milestone' | 'commit_hash' | 'priority' | 'type' | 'legacy_id'>>
): Promise<Task> {
  return request<Task>(`/tasks/${id}`, {
    method: 'PATCH',
    body: JSON.stringify(data),
  });
}

export async function deleteTask(id: number): Promise<void> {
  await request(`/tasks/${id}`, { method: 'DELETE' });
}

export async function listEvents(limit?: number): Promise<Event[]> {
  const qs = limit ? `?limit=${limit}` : '';
  return request<Event[]>(`/events${qs}`);
}

export async function listTaskEvents(taskId: number): Promise<Event[]> {
  return request<Event[]>(`/tasks/${taskId}/events`);
}

export async function createEvent(data: {
  type?: string;
  message: string;
  task_id?: number | null;
}): Promise<Event> {
  return request<Event>('/events', {
    method: 'POST',
    body: JSON.stringify(data),
  });
}

export async function getStatus(): Promise<StatusSummary> {
  return request<StatusSummary>('/status');
}

export async function getConfig(): Promise<{ github_repo_url: string }> {
  return request<{ github_repo_url: string }>('/config');
}

export async function getMilestoneOrder(): Promise<string[]> {
  return request<string[]>('/settings/milestone-order');
}

export async function setMilestoneOrder(order: string[]): Promise<string[]> {
  return request<string[]>('/settings/milestone-order', {
    method: 'PUT',
    body: JSON.stringify(order),
  });
}
