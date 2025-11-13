import { Droppable } from "@hello-pangea/dnd";
import TaskCard from "./TaskCard";
import type { Task } from "../types.ts";
import { useState } from "react";
import AddTaskModal from "./AddTaskModal.tsx";

interface ColumnProps {
  id: string;
  tasks: Task[];
  title: string;
  onTaskCreated: (task: Task) => void;
}

const Column = ({ id, tasks, title, onTaskCreated }: ColumnProps) => {
  const [showModal, setShowModal] = useState(false);

  return (
    <Droppable droppableId={id}>
      {(provided) => (
        <div
          ref={provided.innerRef}
          {...provided.droppableProps}
          className="bg-light p-3 rounded shadow-sm"
        >
          <div className="d-flex justify-content-between align-items-center mb-3">
            <h5 className="mb-0">{title}</h5>
            <button
              className="btn btn-sm btn-outline-dark"
              onClick={() => setShowModal(true)}
            >
              +
            </button>
          </div>

          {tasks.map((task, index) => (
            <TaskCard key={task.id} task={task} index={index} />
          ))}
          {provided.placeholder}

          {showModal && (
            <AddTaskModal
              defaultStatus={id as "to_do" | "in_progress" | "done"}
              onClose={() => setShowModal(false)}
              onTaskCreated={onTaskCreated}
            />
          )}
        </div>
      )}
    </Droppable>
  );
};

export default Column;
