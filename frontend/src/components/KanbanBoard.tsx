import { DragDropContext, type DropResult } from "@hello-pangea/dnd";
import { useState } from "react";
import Column from "./Column";
import type { Task } from "../types";

type TaskMap = {
  to_do: Task[];
  in_progress: Task[];
  done: Task[];
};

const KanbanBoard = () => {
  const [tasks, setTasks] = useState<TaskMap>({
    to_do: [
      { id: "0001", title: "AAA", description: "Desc A", status: "to_do" },
      { id: "0002", title: "BBB", description: "Desc B", status: "to_do" },
    ],
    in_progress: [
      {
        id: "0003",
        title: "CCC",
        description: "Desc C",
        status: "in_progress",
      },
    ],
    done: [{ id: "0004", title: "DDD", description: "Desc D", status: "done" }],
  });

  const onDragEnd = (result: DropResult) => {
    const { source, destination } = result;
    if (!destination) return;

    const sourceCol = source.droppableId as keyof TaskMap;
    const destCol = destination.droppableId as keyof TaskMap;
    const isSameCol = sourceCol === destCol;

    const sourceTasks = [...tasks[sourceCol]];
    const [movedTask] = sourceTasks.splice(source.index, 1);

    // Atualiza o status da tarefa se mudou de coluna
    const updatedTask = {
      ...movedTask,
      status: destCol,
    };

    const destTasks = isSameCol ? sourceTasks : [...tasks[destCol]];
    destTasks.splice(destination.index, 0, updatedTask);

    setTasks({
      ...tasks,
      [sourceCol]: sourceTasks,
      [destCol]: destTasks,
    });
  };

  return (
    <div className="container mt-4">
      <h2 className="text-center mb-4">📌 Progress Board</h2>
      <DragDropContext onDragEnd={onDragEnd}>
        <div className="row">
          <div className="col-md-4">
            <Column title="To Do" tasks={tasks.to_do} id="to_do" />
          </div>
          <div className="col-md-4">
            <Column
              title="In Progress"
              tasks={tasks.in_progress}
              id="in_progress"
            />
          </div>
          <div className="col-md-4">
            <Column title="Done" tasks={tasks.done} id="done" />
          </div>
        </div>
      </DragDropContext>
    </div>
  );
};

export default KanbanBoard;
