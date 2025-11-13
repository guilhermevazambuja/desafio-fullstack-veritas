export interface Task {
  id: string;
  title: string;
  description?: string;
  status: "to_do" | "in_progress" | "done";
}

export interface ApiResponse<T> {
  data: T;
}
