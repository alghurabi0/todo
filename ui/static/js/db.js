import {
  getFirestore,
  doc,
  onSnapshot,
} from "https://www.gstatic.com/firebasejs/10.8.1/firebase-firestore.js";
import { app } from "./auth.js";

const db = getFirestore(app);
export { db };

document.addEventListener("DOMContentLoaded", () => {
  const userId = "PxmNqnn3zmeb6ZC9cyTUG1NLfyF2";
  const board = document.querySelector(".board");

  if (board) {
    const boardId = board.id;
    onSnapshot(doc(db, `users/${userId}/boards`, boardId), (doc) => {
      mapDataToBoard(board, doc.data());
    });

    const columns = board.querySelectorAll(".column");
    if (columns) {
      for (const column of columns) {
        const columnId = column.id;
        onSnapshot(
          doc(db, `users/${userId}/boards/${boardId}/columns`, columnId),
          (doc) => {
            mapDataToColumn(column, doc.data());
          }
        );
      }
    } else {
      console.log("no columns");
    }

    const groups = board.querySelectorAll(".group");
    if (groups) {
      for (const group of groups) {
        const groupId = group.id;
        onSnapshot(
          doc(db, `users/${userId}/boards/${boardId}/groups`, groupId),
          (doc) => {
            mapDataToGroup(group, doc.data());
          }
        );

        const tasks = group.querySelectorAll(".task");
        if (tasks) {
          for (const task of tasks) {
            const taskId = task.id;
            onSnapshot(
              doc(
                db,
                `users/${userId}/boards/${boardId}/groups/${groupId}/tasks`,
                taskId
              ),
              (doc) => {
                mapDataToTask(task, doc.data());
              }
            );
          }
        } else {
          console.log("no tasks");
        }
      }
    } else {
      console.log("no groups");
    }
  } else {
    console.log("no board");
  }
});

function mapDataToBoard(board, data) {
  const board_title = document.querySelector(".board_title");
  const board_desc = document.querySelector(".board_description");
  board_title.innerText = data.title;
  board_desc.innerText = data.description;
}
function mapDataToGroup(group, data) {
  const group_name = group.querySelector(".group_name");
  group_name.innerText = data.name;
}
function mapDataToTask(task, data) {
  const task_content = task.querySelector(".task_content");
  task_content.innerText = data.content;

  const column_values = data.column_values;
  if (column_values) {
    const entries = Object.entries(column_values);
    entries.forEach(([key, value]) => {
      const field = task.querySelector(`#${key}`);
      field.innerText = value;
    });
  }
}
function mapDataToColumn(column, data) {
  column.innerText = data.name;
}
