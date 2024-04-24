import {
  getFirestore,
  doc,
  onSnapshot,
  query,
  collection,
  orderBy,
} from "https://www.gstatic.com/firebasejs/10.8.1/firebase-firestore.js";
import { app } from "./auth.js";

const db = getFirestore(app);
export { db };

const userId = "PxmNqnn3zmeb6ZC9cyTUG1NLfyF2";

const boardObj = {
  id: "",
  title: Element,
  description: Element,
  columns: {},
  groups: {},
};

document.addEventListener("DOMContentLoaded", () => {
  const board = document.querySelector(".board");
  if (board) {
    const board_title = document.querySelector(".board_title");
    const board_desc = document.querySelector(".board_description");
    boardObj.id = board.id;
    boardObj.title = board_title;
    boardObj.description = board_desc;
    const columns = board.querySelectorAll(".column");
    if (columns) {
      for (const col of columns) {
        const colId = col.id;
        boardObj.columns[colId] = col;
      }
    }
    const groups = board.querySelectorAll(".group");
    if (groups) {
      for (const group of groups) {
        group.name = group.querySelector(".group_name");
        group.tasks = {};
        const groupId = group.id;
        const tasks = group.querySelectorAll(".task");
        if (tasks) {
          for (const task of tasks) {
            const taskId = task.id;
            task.content = task.querySelector(".task_content");
            group.tasks[taskId] = task;
          }
        }
        boardObj.groups[groupId] = group;
      }
    }
  }
  console.log(boardObj);
});
document.addEventListener("DOMContentLoaded", () => {
  if (boardObj) {
    onSnapshot(doc(db, `users/${userId}/boards`, boardObj.id), (doc) => {
      mapDataToBoard(boardObj, doc.data());
    });
    if (boardObj.groups) {
      for (const groupId in boardObj.groups) {
        const group = boardObj.groups[groupId];
        onSnapshot(
          doc(db, `users/${userId}/boards/${boardObj.id}/groups`, groupId),
          (doc) => {
            mapDataToGroup(group, doc.data());
          },
        );

        if (group.tasks) {
          for (const taskId in group.tasks) {
            const task = group.tasks[taskId];
            onSnapshot(
              doc(
                db,
                `users/${userId}/boards/${boardObj.id}/groups/${groupId}/tasks`,
                taskId,
              ),
              (doc) => {
                mapDataToTask(task, doc.data());
              },
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
// Listen for changes to columns
document.addEventListener("DOMContentLoaded", () => {
  let isFirstSnapShot = true;
  const q = query(
    collection(db, `users/${userId}/boards/${boardObj.id}/columns`),
  );
  onSnapshot(q, (snapshot) => {
    if (!isFirstSnapShot) {
      snapshot.docChanges().forEach((change) => {
        if (change.type === "added") {
          createColumn(change.doc.id, change.doc.data());
        }
        if (change.type === "modified") {
          console.log("Modified column: ", change.doc.data());
          const column = boardObj.columns[change.doc.id];
          mapDataToColumn(column, change.doc.data());
        }
        if (change.type === "removed") {
          console.log("Removed column: ", change.doc.data());
          removeColumn(change.doc.id);
        }
      });
    }
    isFirstSnapShot = false;
  });
});

function mapDataToBoard(board, data) {
  board.title.value = data.title;
  board.description.innerText = data.description;
}
function mapDataToGroup(group, data) {
  group.name.value = data.name;
}
function mapDataToTask(task, data) {
  // task.content.value = data.content;
  // const column_values = data.column_values;
  // if (column_values) {
  //   const entries = Object.entries(column_values);
  //   entries.forEach(([key, value]) => {
  //     if (value["type"] === "Status") {
  //       console.log("status field", key);
  //     } else {
  //       const field = task.querySelector(`#${key}`);
  //       field.children[0].value = value["value"];
  //     }
  //   });
  // }
}
function mapDataToColumn(column, data) {
  column.children[0].value = data.name;
}
function createColumn(colId, data) {
  const column = document.createElement("th");
  column.id = colId;
  column.classList.add("column");
  if (data.type == "Status") {
    column.classList.add("status_column");
  }
  const input = document.createElement("input");
  input.value = data.name;
  input.setAttribute("data-type", data.type);
  input.type = "text";
  column.appendChild(input);

  const columns = document.querySelector(".columns");
  columns.insertBefore(column, columns.lastElementChild.previousElementSibling);
  boardObj.columns[colId] = column;
  // changes to column data will be handled by the column snapshot listener
  // populate tasks with empty fields for this column
  for (const group of Object.values(boardObj.groups)) {
    for (const task of Object.values(group.tasks)) {
      const field = document.createElement("td");
      field.id = colId;
      switch (data.type) {
        case "Text" || "Number":
          const input = document.createElement("input");
          input.setAttribute("data-type", data.type);
          input.name = "colVal";
          input.type = data.type.toLowerCase();
          input.setAttribute(
            "hx-put",
            `/task/update/colVal/${boardObj.id}/${group.id}/${task.id}/${colId}`,
          );
          field.appendChild(input);
          break;
        case "Status":
          const dropdown = document.createElement("sl-dropdown");
          dropdown.setAttribute("stay-open-on-select", "");
          const openButton = document.createElement("sl-button");
          openButton.slot = "trigger";
          openButton.setAttribute("caret", "");
          openButton.innerText = "Status";
          dropdown.appendChild(openButton);
          const menu = document.createElement("sl-menu");
          dropdown.appendChild(menu);
          field.appendChild(dropdown);
          // <sl-color-picker format="hex" no-format-toggle size="small" hoist></sl-color-picker>
          // <button hx-post="/testRoute">Add</button>
          const menu_item = document.createElement("sl-menu-item");
          const form = document.createElement("form");
          const formInput = document.createElement("input");
          formInput.name = "status_name";
          form.appendChild(formInput);
          const formColorPicker = document.createElement("sl-color-picker");
          formColorPicker.name = "status_color";
          formColorPicker.setAttribute("format", "hex");
          formColorPicker.setAttribute("no-format-toggle", "");
          formColorPicker.setAttribute("size", "small");
          formColorPicker.setAttribute("hoist", "");
          form.appendChild(formColorPicker);
          const formButton = document.createElement("button");
          formButton.innerText = "Add";
          formButton.setAttribute(
            "hx-post",
            `/status/create/${boardObj.id}/${group.id}/${task.id}`,
          );
          form.appendChild(formButton);

          menu_item.appendChild(form);
          menu.appendChild(menu_item);
          const divider = document.createElement("sl-divider");
          menu.appendChild(divider);
          task.statuses = {};
          break;
      }
      task.insertBefore(field, task.lastElementChild);
    }
  }
  if (data.type == "Status") {
    listenToStatus();
  }
}
function removeColumn(colId) {
  const column = boardObj.columns[colId];
  column.remove();
  delete boardObj.columns[colId];
  // remove fields from tasks
  for (const group of Object.values(boardObj.groups)) {
    for (const task of Object.values(group.tasks)) {
      const field = task.querySelector(`#${colId}`);
      field.remove();
    }
  }
}

// get and listen for changes to messages
document.addEventListener("DOMContentLoaded", () => {
  for (const group of Object.values(boardObj.groups)) {
    for (const task of Object.values(group.tasks)) {
      task.messages = {};
      const q = query(
        collection(
          db,
          `users/${userId}/boards/${boardObj.id}/groups/${group.id}/tasks/${task.id}/messages`,
        ),
      );
      onSnapshot(q, (snapshot) => {
        snapshot.docChanges().forEach((change) => {
          if (change.type === "added") {
            console.log("Added message: ", change.doc.data());
            createMessage(task, group, change.doc.id, change.doc.data());
          }
          if (change.type === "modified") {
            console.log("Modified message: ", change.doc.data());
            const message = task.messages[change.doc.id];
            mapDataToMessage(message, change.doc.data());
          }
          if (change.type === "removed") {
            console.log("Removed message: ", change.doc.data());
            removeMessage(task, change.doc.id);
          }
        });
      });
    }
  }
});
function createMessage(task, group, msgId, msgData) {
  const messages = task.querySelector(".messages");
  const message = document.createElement("div");
  message.id = msgId;
  message.classList.add("message");
  const content = document.createElement("p");
  content.classList.add("message_content");
  content.innerText = msgData.content;
  message.appendChild(content);
  message.content = content;
  const repliesDiv = document.createElement("div");
  repliesDiv.classList.add("replies");
  message.appendChild(repliesDiv);
  messages.appendChild(message);
  task.messages[msgId] = message;
  message.replies = {};
  // check if the message has replies
  const q = query(
    collection(
      db,
      `users/${userId}/boards/${boardObj.id}/groups/${group.id}/tasks/${task.id}/messages/${msgId}/replies`,
    ),
  );
  onSnapshot(q, (snapshot) => {
    if (snapshot.empty) {
      console.log("no replies");
    }
    snapshot.docChanges().forEach((change) => {
      if (change.type === "added") {
        console.log("Added reply: ", change.doc.data());
        createReply(message, change.doc.id, change.doc.data());
      }
      if (change.type === "modified") {
        console.log("Modified reply: ", change.doc.data());
        const reply = message.replies[change.doc.id];
        mapDataToReply(reply, change.doc.data());
      }
      if (change.type === "removed") {
        console.log("Removed reply: ", change.doc.data());
        removeReply(message, change.doc.id);
      }
    });
  });
}
function createReply(message, replyId, replyData) {
  console.log("creating reply");
  const replies = message.querySelector(".replies");
  const reply = document.createElement("div");
  reply.id = replyId;
  reply.classList.add("reply");
  const content = document.createElement("p");
  content.classList.add("reply_content");
  content.innerText = replyData.content;
  reply.appendChild(content);
  replies.appendChild(reply);
  reply.content = content;
  message.replies[replyId] = reply;
}
function mapDataToMessage(message, data) {
  message.content.innerText = data.content;
}
function mapDataToReply(reply, data) {
  reply.content.innerText = data.content;
}
function removeMessage(task, msgId) {
  const message = task.messages[msgId];
  message.remove();
  delete task.messages[msgId];
}
function removeReply(message, replyId) {
  const reply = message.replies[replyId];
  reply.remove();
  delete message.replies[replyId];
}

//status
function listenToStatus() {
  for (const group of Object.values(boardObj.groups)) {
    for (const task of Object.values(group.tasks)) {
      // status is a collection in a task where each doc is a status label.
      // the init statuses will be listened for when the column get created.
      // so the backend html template should just include dropdown, button, and menu. -- new input as wll
      // if the column already exist it will be listened for when the dom loads, then labels will render dianamically
      // as well as chosen label will be in the column values.
      const q = query(
        collection(
          db,
          `users/${userId}/boards/${boardObj.id}/groups/${group.id}/tasks/${task.id}/statuses`,
        ),
      );
      onSnapshot(q, (snapshot) => {
        snapshot.docChanges().forEach((change) => {
          if (change.type === "added") {
            console.log("status added", change.doc.data());
            createStatus();
          }
          if (change.type === "modified") {
            console.log("status modified", change.doc.data());
            mapDataToStatus();
          }
          if (change.type === "removed") {
            console.log("status removed", change.doc.data());
            removeStatus();
          }
        });
      });
    }
  }
  activateStatuses();
}
//if status already exist
document.addEventListener("DOMContentLoaded", () => {
  const status_column = document.querySelector(".status_column");
  if (status_column) {
    console.log("got status column");
    for (const group of Object.values(boardObj.groups)) {
      for (const task of Object.values(group.tasks)) {
        task.statuses = {};
        const q = query(
          collection(
            db,
            `users/${userId}/boards/${boardObj.id}/groups/${group.id}/tasks/${task.id}/statuses`,
          ),
        );
        onSnapshot(q, (snapshot) => {
          snapshot.docChanges().forEach((change) => {
            if (change.type === "added") {
              console.log("status added", change.doc.data());
              createStatus(task, change.doc.id, change.doc.data());
            }
            if (change.type === "modified") {
              const status = task.statuses[change.doc.id];
              console.log("status modified", change.doc.data());
              mapDataToStatus(status, change.doc.data());
            }
            if (change.type === "removed") {
              console.log("status removed", change.doc.data());
              removeStatus(task, change.doc.data());
            }
          });
        });
      }
    }
    activateStatuses();
  }
});
function createStatus(task, statusId, data) {
  const menu = task.querySelector("sl-menu");
  const status = document.createElement("sl-menu-item");
  status.id = statusId;
  const labelEl = document.createElement("label");
  labelEl.innerText = data.text;
  labelEl.style.backgroundColor = data.color;
  status.appendChild(labelEl);
  const colorPicker = document.createElement("sl-color-picker");
  colorPicker.classList.add("status_color_picker");
  colorPicker.setAttribute("format", "hex");
  colorPicker.setAttribute("no-format-toggle", "");
  colorPicker.setAttribute("size", "small");
  colorPicker.setAttribute("hoist", "");
  const editStatusTextButton = document.createElement("button");
  editStatusTextButton.classList.add("edit_status_text_button");
  editStatusTextButton.innerText = "E";
  status.appendChild(colorPicker);
  status.appendChild(editStatusTextButton);
  menu.appendChild(status);
  task.statuses[statusId] = status;
  activateStatus(status);
}
function mapDataToStatus(status, data) {
  const label = status.querySelector("label");
  label.innerText = data.text;
  label.style.backgroundColor = data.color;
}
function removeStatus(task, statusId) {
  const status = task.statuses[statusId];
  status.remove();
  delete task.statuses[statusId];
}
function activateStatuses() {
  for (const group of Object.values(boardObj.groups)) {
    for (const task of Object.values(group.tasks)) {
      for (const status of Object.values(task.statuses)) {
        const picker = status.querySelector("sl-color-picker");
        picker.addEventListener("sl-change", () => {
          const valid = picker.checkValidity();
          if (!valid) {
            return;
          }
          const new_color = picker.getFormattedValue("hex");
          htmx.ajax("PUT", "/testRoute", {
            swap: "none",
            values: { new_color: new_color },
          });
        });
      }
    }
  }
}
function activateStatus(status) {
  const picker = status.querySelector("sl-color-picker");
  picker.addEventListener("sl-change", () => {
    const valid = picker.checkValidity();
    if (!valid) {
      return;
    }
    const new_color = picker.getFormattedValue("hex");
    htmx.ajax("PUT", "/testRoute", {
      swap: "none",
      values: { new_color: new_color },
    });
  });
}
