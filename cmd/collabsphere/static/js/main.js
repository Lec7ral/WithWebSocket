import { initCanvas } from "./canvas.js";
import { initUI, renderUsers, setCurrentRoomLabel, showToast, refreshRoomList, addChatMessage } from "./ui.js";
import { state, setToken, setUserFromToken } from "./state.js";
import { login } from "./api.js";
import { connectToRoom, sendMessage } from "./ws.js";

let typingTimeout = null;
let lastTypingSent = false;

function setupAuth() {
  const loginBtn = document.getElementById("login-btn");
  const usernameInput = document.getElementById("username-input");
  const logoutBtn = document.getElementById("logout-btn");
  const authOut = document.getElementById("auth-logged-out");
  const authIn = document.getElementById("auth-logged-in");
  const currentUserSpan = document.getElementById("current-user");

  async function doLogin() {
    const username = (usernameInput.value || "").trim();
    if (!username) return;
    try {
      await login(username);
      authOut.classList.add("hidden");
      authIn.classList.remove("hidden");
      currentUserSpan.textContent = `Conectado como ${state.user.username}`;
      await refreshRoomList();
      showToast("Autenticado");
    } catch (e) {
      showToast("Error de login");
    }
  }

  loginBtn.addEventListener("click", doLogin);
  usernameInput.addEventListener("keydown", (e) => {
    if (e.key === "Enter") doLogin();
  });

  logoutBtn.addEventListener("click", () => {
    setToken(null);
    setUserFromToken(null);
    authIn.classList.add("hidden");
    authOut.classList.remove("hidden");
    currentUserSpan.textContent = "";
    renderUsers();
    setCurrentRoomLabel("");
    showToast("Sesión cerrada");
  });

  // Restore token
  const savedToken = localStorage.getItem("collab_token");
  if (savedToken) {
    setToken(savedToken);
    setUserFromToken(savedToken);
    if (state.user) {
      authOut.classList.add("hidden");
      authIn.classList.remove("hidden");
      currentUserSpan.textContent = `Conectado como ${state.user.username}`;
      refreshRoomList();
    } else {
      setToken(null);
    }
  }
}

function setupRooms() {
  const joinBtn = document.getElementById("join-room-btn");
  const roomInput = document.getElementById("room-id-input");

  async function joinRoom() {
    const roomId = (roomInput.value || "").trim();
    if (!roomId) return;
    if (!state.token) {
      showToast("Primero inicia sesión");
      return;
    }
    connectToRoom(roomId);
    setCurrentRoomLabel(roomId);
  }

  joinBtn.addEventListener("click", joinRoom);
  roomInput.addEventListener("keydown", (e) => {
    if (e.key === "Enter") joinRoom();
  });
}

function setupChat() {
  const input = document.getElementById("chat-input");
  const sendBtn = document.getElementById("send-chat-btn");
  const dmSelect = document.getElementById("dm-recipient-select");

  function sendChat() {
    const text = (input.value || "").trim();
    if (!text) return;
    const recipientId = dmSelect.value;
    if (recipientId) {
      sendMessage({
        type: "direct_message",
        payload: {
          recipient_id: recipientId,
          content: text,
        },
      });
      addChatMessage({
        text,
        fromSelf: true,
        dm: true,
        username: state.user?.username || "Tú",
      });
    } else {
      sendMessage({
        type: "text_message",
        payload: text,
      });
      addChatMessage({
        text,
        fromSelf: true,
        dm: false,
        username: state.user?.username || "Tú",
      });
    }
    input.value = "";
    stopTyping();
  }

  sendBtn.addEventListener("click", sendChat);
  input.addEventListener("keydown", (e) => {
    if (e.key === "Enter") {
      e.preventDefault();
      sendChat();
    }
  });

  input.addEventListener("input", () => {
    handleTypingChange(input.value);
  });
}

function handleTypingChange(text) {
  const isTyping = text.trim().length > 0;
  if (isTyping && !lastTypingSent) {
    sendMessage({ type: "typing_start", payload: null });
    lastTypingSent = true;
  }
  clearTimeout(typingTimeout);
  typingTimeout = setTimeout(() => {
    stopTyping();
  }, 1500);
}

function stopTyping() {
  if (lastTypingSent) {
    sendMessage({ type: "typing_stop", payload: null });
    lastTypingSent = false;
  }
  clearTimeout(typingTimeout);
}

window.addEventListener("DOMContentLoaded", () => {
  initUI();
  initCanvas();
  setupAuth();
  setupRooms();
  setupChat();
});