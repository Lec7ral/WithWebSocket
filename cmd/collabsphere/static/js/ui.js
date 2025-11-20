import { state } from "./state.js";
import { listRooms } from "./api.js";

let usersListEl,
  userCountEl,
  roomListEl,
  currentRoomLabelEl,
  chatMessagesEl,
  typingIndicatorEl,
  dmSelectEl,
  toastEl;

export function initUI() {
  usersListEl = document.getElementById("users-list");
  userCountEl = document.getElementById("user-count");
  roomListEl = document.getElementById("room-list");
  currentRoomLabelEl = document.getElementById("current-room-label");
  chatMessagesEl = document.getElementById("chat-messages");
  typingIndicatorEl = document.getElementById("typing-indicator");
  dmSelectEl = document.getElementById("dm-recipient-select");
  toastEl = document.getElementById("toast");
}

export function renderUsers() {
  if (!usersListEl) return;
  usersListEl.innerHTML = "";
  const users = Array.from(state.users.values());
  userCountEl.textContent = users.length.toString();
  dmSelectEl.innerHTML = `<option value="">Sala</option>`;
  users.forEach((u) => {
    const item = document.createElement("div");
    item.className = "user-item";
    if (u.id === state.user?.id) item.classList.add("self");

    const info = document.createElement("div");
    const nameSpan = document.createElement("div");
    nameSpan.className = "user-name";
    nameSpan.textContent = u.username || "usuario";
    const idSpan = document.createElement("div");
    idSpan.className = "user-id";
    idSpan.textContent = u.id.slice(0, 8);
    info.appendChild(nameSpan);
    info.appendChild(idSpan);

    const dmBtn = document.createElement("button");
    dmBtn.className = "user-dm-btn";
    dmBtn.textContent = "";
    dmBtn.addEventListener("click", () => {
      dmSelectEl.value = u.id;
      dmSelectEl.dispatchEvent(new Event("change"));
    });

    item.appendChild(info);
    item.appendChild(dmBtn);
    usersListEl.appendChild(item);

    if (u.id !== state.user?.id) {
      const opt = document.createElement("option");
      opt.value = u.id;
      opt.textContent = `DM: ${u.username}`;
      dmSelectEl.appendChild(opt);
    }
  });
}

export function addChatMessage({ text, fromSelf, dm, username }) {
  if (!chatMessagesEl) return;
  const wrapper = document.createElement("div");
  wrapper.className = "chat-message";
  if (fromSelf) wrapper.classList.add("self");
  if (dm) wrapper.classList.add("dm");

  const meta = document.createElement("div");
  meta.className = "chat-meta";
  const senderSpan = document.createElement("span");
  senderSpan.textContent = username || (fromSelf ? "Tú" : "Usuario");
  const tagSpan = document.createElement("span");
  tagSpan.textContent = dm ? "DM" : "Sala";
  meta.appendChild(senderSpan);
  meta.appendChild(tagSpan);

  const textDiv = document.createElement("div");
  textDiv.className = "chat-text";
  textDiv.textContent = text;

  wrapper.appendChild(meta);
  wrapper.appendChild(textDiv);
  chatMessagesEl.appendChild(wrapper);
  chatMessagesEl.scrollTop = chatMessagesEl.scrollHeight;
}

export function setTypingIndicator(userIds) {
  if (!typingIndicatorEl) return;
  if (!userIds || userIds.length === 0) {
    typingIndicatorEl.textContent = "";
    return;
  }
  const names = userIds
    .map((id) => state.users.get(id)?.username || "Alguien")
    .slice(0, 2);
  let text = "";
  if (names.length === 1) {
    text = `${names[0]} está escribiendo...`;
  } else {
    text = `${names[0]} y ${names[1]} están escribiendo...`;
  }
  typingIndicatorEl.textContent = text;
}

export function setCurrentRoomLabel(roomId) {
  if (!currentRoomLabelEl) return;
  currentRoomLabelEl.textContent = roomId ? `Sala: ${roomId}` : "";
}

let toastTimeout;
export function showToast(msg) {
  if (!toastEl) return;
  toastEl.textContent = msg;
  toastEl.classList.remove("hidden");
  clearTimeout(toastTimeout);
  toastTimeout = setTimeout(() => {
    toastEl.classList.add("hidden");
  }, 2500);
}

export async function refreshRoomList() {
  if (!roomListEl) return;
  try {
    const rooms = await listRooms();
    roomListEl.innerHTML = "";
    rooms.forEach((room) => {
      const chip = document.createElement("button");
      chip.type = "button";
      chip.className = "room-chip";
      const nameSpan = document.createElement("span");
      nameSpan.textContent = room.id;
      const countSpan = document.createElement("span");
      countSpan.className = "badge";
      countSpan.textContent = room.clientCount.toString();
      chip.appendChild(nameSpan);
      chip.appendChild(countSpan);
      chip.addEventListener("click", () => {
        const input = document.getElementById("room-id-input");
        if (input) {
          input.value = room.id;
        }
        document.getElementById("join-room-btn")?.click();
      });
      roomListEl.appendChild(chip);
    });
  } catch (e) {
    showToast("No se pudieron cargar las salas");
  }
}