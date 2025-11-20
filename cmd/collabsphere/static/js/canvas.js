import { sendMessage } from "./ws.js";

let canvas, ctx;
let isDrawing = false;
let color = "#2f8f5b";
let lineWidth = 3;
let lastPoint = null;

export function initCanvas() {
  canvas = document.getElementById("whiteboard");
  if (!canvas) return;
  ctx = canvas.getContext("2d");
  resizeCanvas();
  window.addEventListener("resize", resizeCanvas);

  // Mouse
  canvas.addEventListener("mousedown", handlePointerDown);
  canvas.addEventListener("mousemove", handlePointerMove);
  window.addEventListener("mouseup", handlePointerUp);

  // Touch
  canvas.addEventListener("touchstart", handlePointerDown, { passive: false });
  canvas.addEventListener("touchmove", handlePointerMove, { passive: false });
  window.addEventListener("touchend", handlePointerUp);
  window.addEventListener("touchcancel", handlePointerUp);

  const colorInput = document.getElementById("color-input");
  const widthInput = document.getElementById("line-width-input");
  const clearBtn = document.getElementById("clear-board-btn");

  colorInput?.addEventListener("input", (e) => {
    color = e.target.value || color;
  });
  widthInput?.addEventListener("input", (e) => {
    lineWidth = Number(e.target.value) || lineWidth;
  });
  clearBtn?.addEventListener("click", () => {
    clearCanvas();
    sendMessage({ type: "clear_board", payload: null });
  });
}

function resizeCanvas() {
  if (!canvas) return;
  const rect = canvas.getBoundingClientRect();
  const ratio = window.devicePixelRatio || 1;
  canvas.width = rect.width * ratio;
  canvas.height = rect.height * ratio;
  ctx.setTransform(ratio, 0, 0, ratio, 0, 0);
}

function getPointFromEvent(ev) {
  let x, y;
  if (ev.touches && ev.touches[0]) {
    const t = ev.touches[0];
    const rect = canvas.getBoundingClientRect();
    x = t.clientX - rect.left;
    y = t.clientY - rect.top;
  } else if (ev.changedTouches && ev.changedTouches[0]) {
    const t = ev.changedTouches[0];
    const rect = canvas.getBoundingClientRect();
    x = t.clientX - rect.left;
    y = t.clientY - rect.top;
  } else {
    const rect = canvas.getBoundingClientRect();
    x = ev.clientX - rect.left;
    y = ev.clientY - rect.top;
  }
  return { x, y };
}

function handlePointerDown(ev) {
  ev.preventDefault();
  if (!ctx) return;
  const pt = getPointFromEvent(ev);
  isDrawing = true;
  lastPoint = pt;
  drawPoint(pt, color, lineWidth, true);
  sendMessage({
    type: "draw_start",
    payload: { x: pt.x, y: pt.y, color, lineWidth },
  });
}

function handlePointerMove(ev) {
  if (!isDrawing) return;
  ev.preventDefault();
  const pt = getPointFromEvent(ev);
  drawLine(lastPoint, pt, color, lineWidth, true);
  lastPoint = pt;
  sendMessage({
    type: "draw_move",
    payload: { x: pt.x, y: pt.y, color, lineWidth },
  });
}

function handlePointerUp(ev) {
  if (!isDrawing) return;
  ev.preventDefault();
  isDrawing = false;
  lastPoint = null;
  sendMessage({ type: "draw_end", payload: null });
}

function drawPoint(pt, color, width, localOnly) {
  if (!ctx) return;
  ctx.save();
  ctx.fillStyle = color;
  ctx.beginPath();
  ctx.arc(pt.x, pt.y, width / 2, 0, Math.PI * 2);
  ctx.fill();
  ctx.restore();
}

function drawLine(from, to, color, width, localOnly) {
  if (!ctx || !from || !to) return;
  ctx.save();
  ctx.strokeStyle = color;
  ctx.lineWidth = width;
  ctx.lineCap = "round";
  ctx.beginPath();
  ctx.moveTo(from.x, from.y);
  ctx.lineTo(to.x, to.y);
  ctx.stroke();
  ctx.restore();
}

export function clearCanvas() {
  if (!ctx || !canvas) return;
  ctx.clearRect(0, 0, canvas.width, canvas.height);
}

export function applyRemoteDrawEvent(evt) {
  const { type, payload } = evt;
  if (!ctx) return;
  if (type === "clear_board") {
    clearCanvas();
    return;
  }
  if (type === "draw_start") {
    const pt = { x: payload.x, y: payload.y };
    drawPoint(pt, payload.color, payload.lineWidth, false);
    lastPoint = pt;
  } else if (type === "draw_move") {
    const pt = { x: payload.x, y: payload.y };
    drawLine(lastPoint, pt, payload.color, payload.lineWidth, false);
    lastPoint = pt;
  } else if (type === "draw_end") {
    lastPoint = null;
  }
}

export function replayWhiteboardEvents(events) {
  clearCanvas();
  lastPoint = null;
  events.forEach((evt) => {
    applyRemoteDrawEvent(evt);
  });
}