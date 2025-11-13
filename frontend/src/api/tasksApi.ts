import type { Task } from "../types";
import type { ApiResponse } from "../types";
import { api } from "./api.ts";

export async function fetchTasks(): Promise<Task[]> {
  try {
    const response = await api.get<ApiResponse<Task[]>>("/tasks");
    return response.data.data;
  } catch (error) {
    console.error("Error fetching tasks:", error);
    return [];
  }
}

export async function createTask(task: Task): Promise<Task> {
  const response = await api.post<ApiResponse<Task>>("/tasks", task);
  return response.data.data;
}

export async function updateTask(id: string, task: Partial<Task>) {
  const response = await api.put<ApiResponse<Task>>(`/tasks/${id}`, task);
  return response.data.data;
}

export async function deleteTask(id: string) {
  await api.delete(`/tasks/${id}`);
}
