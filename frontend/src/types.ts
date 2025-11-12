export interface Task {
  id: string;
  title: string;
  description?: string;
  status: "to_do" | "in_progress" | "done";
}
