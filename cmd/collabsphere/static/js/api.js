import { setToken, setUserFromToken } from "./state.js";

const BASE = "";

async function jsonFetch(path, opts = {}) {
  const res = await fetch(BASE + path, {
    headers: {
      "Content-Type": "application/json",
      ...(opts.headers || {}),
    },
    ...opts,
  });
  if (!res.ok) {
    const text = await res.text();
    const err = new Error(text || res.statusText);
    err.status = res.status;
    throw err;
  }
  const ct = res.headers.get("content-type") || "";
  if (ct.includes("application/json")) {
    return res.json();
  }
  return res.text();
}

export async function login(username) {
  const body = JSON.stringify({ username });
  const data = await jsonFetch("/login", { method: "POST", body });
  if (!data.token) {
    throw new Error("Token no recibido");
  }
  setToken(data.token);
  setUserFromToken(data.token);
  return data.token;
}

export async function listRooms() {
  return jsonFetch("/api/rooms", { method: "GET" });
}

export async function getUser(userId) {
  return jsonFetch(`/api/users/${encodeURIComponent(userId)}`, {
    method: "GET",
  });
}

