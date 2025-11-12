import { Droppable } from "@hello-pangea/dnd";
import TaskCard from "./TaskCard";
import type { Task } from "../types.ts";

interface ColumnProps {
  id: string;
  tasks: Task[];
  title: string;
}

const Column = ({ id, tasks, title }: ColumnProps) => {
  return (
    <Droppable droppableId={id}>
      {(provided) => (
        <div
          className="bg-light p-3 rounded shadow-sm"
          ref={provided.innerRef}
          {...provided.droppableProps}
        >
          <h5 className="text-center mb-3">{title}</h5>
          {tasks.map((task, index) => (
            <TaskCard key={task.id} task={task} index={index} />
          ))}

          {provided.placeholder}
        </div>
      )}
    </Droppable>
  );
};

export default Column;
