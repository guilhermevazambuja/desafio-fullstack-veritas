import type { Task } from "../types";
import type { ApiResponse } from "../types";
import { api } from "./api.ts";

export async function fetchTasks(): Promise<Task[]> {
  try {
    const response = await api.get<ApiResponse<Task[]>>("/tasks");
    return response.data.data;
  } catch (error) {
    console.error("Erro ao buscar tasks:", error);
    return [];
  }
}
