import { Draggable } from "@hello-pangea/dnd";
import type { Task } from "../types";
import classNames from "classnames";

interface TaskCardProps {
  index: number;
  task: Task;
}

const TaskCard = ({ index, task }: TaskCardProps) => {
  return (
    <Draggable draggableId={task.id} index={index}>
      {(provided, snapshot) => (
        <div
          className={classNames("card mb-2 task-card", {
            dragging: snapshot.isDragging,
          })}
          ref={provided.innerRef}
          {...provided.draggableProps}
          {...provided.dragHandleProps}
        >
          <div className="card-body p-2">
            <h6 className="mb-1">{task.title}</h6>
            {task.description && (
              <p className="mb-0 text-muted">{task.description}</p>
            )}
            <small className="text-uppercase badge bg-secondary mt-1">
              {task.status}
            </small>
          </div>
        </div>
      )}
    </Draggable>
  );
};

export default TaskCard;
