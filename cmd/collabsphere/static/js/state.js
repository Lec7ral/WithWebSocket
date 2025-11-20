export const state = {
  token: null,
  user: null, // { id, username, sessionId }
  currentRoomId: null,
  ws: null,
  users: new Map(), // id -> { id, username }
  userNameByIdCache: new Map(),
  typingUsers: new Set(),
};

export function setToken(token) {
  state.token = token;
  if (token) {
    localStorage.setItem("collab_token", token);
  } else {
    localStorage.removeItem("collab_token");
  }
}

export function decodeJwt(token) {
  try {
    const [, payloadBase64] = token.split(".");
    if (!payloadBase64) return null;
    const json = atob(payloadBase64.replace(/-/g, "+").replace(/_/g, "/"));
    return JSON.parse(json);
  } catch {
    return null;
  }
}

export function setUserFromToken(token) {
  if (!token) {
    state.user = null;
    return;
  }
  const payload = decodeJwt(token);
  if (!payload) {
    state.user = null;
    return;
  }
  state.user = {
    id: payload.userId,
    username: payload.username,
    sessionId: payload.sessionId,
  };
}

export function resetRoomState() {
  state.currentRoomId = null;
  state.users.clear();
  state.typingUsers.clear();
}

