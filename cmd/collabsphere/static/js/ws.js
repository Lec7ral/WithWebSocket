import { state, resetRoomState } from "./state.js";
import { renderUsers, addChatMessage, setTypingIndicator, replayWhiteboardEvents } from "./ui.js";
import { applyRemoteDrawEvent } from "./canvas.js";
import { showToast } from "./ui.js";
import { getUser } from "./api.js";

let reconnecting = false;

export function connectToRoom(roomId) {
  if (!state.token || !roomId) return;
  if (state.ws) {
    state.ws.close();
    state.ws = null;
  }
  resetRoomState();
  state.currentRoomId = roomId;

  const scheme = location.protocol === "https:" ? "wss" : "ws";
  const url =
    `${scheme}://${location.host}/ws/` +
    encodeURIComponent(roomId) +
    `?token=${encodeURIComponent(state.token)}`;

  const ws = new WebSocket(url);
  state.ws = ws;

  ws.onopen = () => {
    reconnecting = false;
    showToast(`Conectado a ${roomId}`);
  };

  ws.onmessage = async (event) => {
    let msg;
    try {
      msg = JSON.parse(event.data);
    } catch {
      return;
    }
    const { type, payload, sender, room_id } = msg;
    if (room_id && !state.currentRoomId) {
      state.currentRoomId = room_id;
    }
    switch (type) {
      case "initial_state":
        if (payload.users) {
          payload.users.forEach((u) => {
            state.users.set(u.id, u);
          });
          renderUsers();
        }
        if (payload.whiteboard && Array.isArray(payload.whiteboard.events)) {
          replayWhiteboardEvents(payload.whiteboard.events);
        }
        break;

      case "user_list_update":
        state.users.clear();
        payload.forEach((u) => state.users.set(u.id, u));
        renderUsers();
        break;

      case "text_message": {
        const userId = sender;
        const content = payload;
        const fromSelf = userId === state.user?.id;
        const meta = await ensureUserMeta(userId);
        addChatMessage({
          text: content,
          fromSelf,
          dm: false,
          username: meta?.username || userId,
        });
        break;
      }

      case "direct_message": {
        const userId = sender;
        const content = payload;
        const fromSelf = userId === state.user?.id;
        const meta = await ensureUserMeta(userId);
        addChatMessage({
          text: content,
          fromSelf,
          dm: true,
          username: meta?.username || userId,
        });
        break;
      }

      case "draw_start":
      case "draw_move":
      case "draw_end":
      case "clear_board":
        if (sender && sender === state.user?.id) return;
        applyRemoteDrawEvent({ type, payload, sender });
        break;

      case "typing_start":
        if (sender && sender !== state.user?.id) {
          state.typingUsers.add(sender);
          updateTyping();
        }
        break;

      case "typing_stop":
        if (sender && sender !== state.user?.id) {
          state.typingUsers.delete(sender);
          updateTyping();
        }
        break;

      default:
        break;
    }
  };

  ws.onclose = () => {
    if (!reconnecting && state.currentRoomId) {
      reconnecting = true;
      showToast("Desconectado");
    }
  };

  ws.onerror = () => {
    showToast("Error de conexión");
  };
}

async function ensureUserMeta(userId) {
  if (!userId) return null;
  if (state.users.has(userId)) return state.users.get(userId);
  if (state.userNameByIdCache.has(userId)) {
    return state.userNameByIdCache.get(userId);
  }
  try {
    const data = await getUser(userId);
    state.userNameByIdCache.set(userId, data);
    return data;
  } catch {
    return null;
  }
}

function updateTyping() {
  const ids = Array.from(state.typingUsers);
  setTypingIndicator(ids);
}

export function sendMessage(obj) {
  if (!state.ws || state.ws.readyState !== WebSocket.OPEN) return;
  state.ws.send(JSON.stringify(obj));
}

