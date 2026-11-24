import { action, decodeFrame } from "./wire.js";

const byId = (id) => document.getElementById(id);
const ui = {
  connection: byId("connection"),
  joinForm: byId("join-form"),
  name: byId("player-name"),
  lobby: byId("lobby"),
  lobbyTitle: byId("lobby-title"),
  lobbyHint: byId("lobby-hint"),
  start: byId("start"),
  rejoin: byId("rejoin"),
  board: byId("board"),
  controls: byId("controls"),
  reveal: byId("reveal"),
  flag: byId("flag"),
  flags: byId("flags"),
  flagCount: byId("flag-count"),
  players: byId("players"),
  playerCount: byId("player-count"),
  noticeTitle: byId("notice-title"),
  noticeText: byId("notice-text"),
  again: byId("again"),
  leave: byId("leave"),
};

const state = {
  socket: null,
  phase: "idle",
  playerId: null,
  hostId: null,
  players: [],
  startPending: false,
  board: null,
  tiles: [],
  cursors: [],
  tool: "reveal",
  outcome: null,
  notice: "",
};
let rematchTimer;

ui.name.value = (new URLSearchParams(location.search).get("player") || "玩家").slice(0, 15);

function send(bytes) {
  if (state.socket?.readyState !== WebSocket.OPEN) return;
  state.socket.send(bytes);
}

function connect(name) {
  const url = new URL("/ws", location.href);
  url.protocol = location.protocol === "https:" ? "wss:" : "ws:";
  url.searchParams.set("name", name);
  const socket = new WebSocket(url);
  socket.binaryType = "arraybuffer";
  state.socket = socket;
  state.phase = "connecting";
  state.notice = "";
  render();

  socket.addEventListener("open", render);
  socket.addEventListener("message", ({ data }) => {
    try {
      const frame = decodeFrame(data);
      if (frame) receive(frame);
    } catch (error) {
      state.notice = "棋盘消息读取失败：" + error.message;
      render();
    }
  });
  socket.addEventListener("error", () => {
    state.notice = "连接失败。确认试玩服务已启动，再试一次。";
  });
  socket.addEventListener("close", () => {
    const notice = state.notice || "已退出。你可以重新加入大厅。";
    clearTimeout(rematchTimer);
    Object.assign(state, {
      socket: null,
      phase: "idle",
      playerId: null,
      hostId: null,
      players: [],
      startPending: false,
      board: null,
      cursors: [],
      tool: "reveal",
      outcome: null,
      notice,
    });
    render();
  });
}

function receive(frame) {
  switch (frame.type) {
    case "welcome":
      state.playerId = frame.id;
      break;
    case "lobby":
      if (state.phase === "playing") return;
      if (!frame.players.some(({ id }) => id === state.playerId)) {
        if (state.phase !== "finished") {
          state.phase = "waiting";
          state.players = frame.players;
        }
        break;
      }
      clearTimeout(rematchTimer);
      state.phase = "lobby";
      state.hostId = frame.host;
      state.players = frame.players;
      state.outcome = null;
      state.notice = "";
      if (state.startPending && state.hostId === state.playerId && frame.mode === 11) {
        state.startPending = false;
        send(action.start);
      }
      break;
    case "board":
      if (state.phase !== "lobby" && state.phase !== "playing") return;
      if (state.phase === "lobby") {
        state.tool = "reveal";
        state.cursors = [];
      }
      state.phase = "playing";
      state.board = frame;
      paintBoard(frame);
      break;
    case "cursors":
      if (state.phase !== "playing") return;
      state.cursors = frame.cursors;
      paintCursors();
      return;
    case "won":
    case "lost":
      if (state.phase !== "playing") return;
      finish(frame);
      break;
    case "error":
      if (frame.message.startsWith("game in progress")) {
        state.phase = "waiting";
      } else {
        state.notice = frame.message;
      }
      break;
  }
  render();
}

function finish(frame) {
  state.phase = "finished";
  const won = frame.type === "won";
  const player = state.players.find(({ id }) => id === frame.id);
  state.outcome = {
    type: won ? "win" : "lose",
    heading: won ? "排雷成功！" : "踩到地雷了",
    text: won ? "两位搭档一起清空了棋盘。" : `${player?.name ?? "玩家"}在第 ${frame.y + 1} 行、第 ${frame.x + 1} 列踩中了雷。`,
    canRematch: false,
  };
  for (const tile of state.tiles) tile.disabled = true;
  rematchTimer = setTimeout(() => {
    state.outcome.canRematch = true;
    render();
  }, won ? 5100 : 3100);
}

function cellAppearance(value) {
  if (value === 0x00) return ["hidden", "未翻开", ""];
  if (value === 0x10) return ["flagged", "已插旗", "⚑"];
  if (value === 0x20) return ["exploded", "踩中地雷", "✹"];
  if (value === 0x39) return ["mine", "地雷", "✹"];
  const neighbors = value - 0x30;
  return ["revealed", `周围 ${neighbors} 雷`, neighbors || ""];
}

function paintBoard(frame) {
  if (ui.board.dataset.size !== `${frame.width}x${frame.height}`) {
    ui.board.style.setProperty("--columns", frame.width);
    ui.board.dataset.size = `${frame.width}x${frame.height}`;
    state.tiles = Array.from(frame.cells, (_, index) => {
      const x = index % frame.width;
      const y = Math.floor(index / frame.width);
      const tile = document.createElement("button");
      tile.className = "cell";
      tile.type = "button";
      tile.dataset.x = x;
      tile.dataset.y = y;
      tile.addEventListener("click", () => playTile(tile, state.tool, x, y));
      tile.addEventListener("contextmenu", (event) => {
        event.preventDefault();
        playTile(tile, "flag", x, y);
      });
      tile.addEventListener("pointerenter", () => {
        if (state.phase === "playing") send(action.cursor(state.playerId, x, y));
      });
      return tile;
    });
    ui.board.replaceChildren(...state.tiles);
  }
  frame.cells.forEach((value, index) => {
    const tile = state.tiles[index];
    tile.disabled = false;
    if (tile.dataset.value === String(value)) return;
    const [kind, label, content] = cellAppearance(value);
    if (tile.dataset.value !== undefined) {
      tile.classList.remove("changed");
      requestAnimationFrame(() => tile.classList.add("changed"));
    }
    tile.dataset.value = value;
    tile.dataset.state = kind;
    tile.dataset.number = kind === "revealed" ? content : "";
    tile.textContent = content;
    tile.setAttribute("aria-label", `第 ${Number(tile.dataset.y) + 1} 行第 ${Number(tile.dataset.x) + 1} 列，${label}`);
  });
  paintCursors();
}

function paintCursors() {
  const marks = new Set();
  for (const { id, x, y } of state.cursors) {
    if (id !== state.playerId) marks.add(y * state.board.width + x);
  }
  state.tiles.forEach((tile, index) => tile.classList.toggle("peer", marks.has(index)));
}

function playTile(tile, tool, x, y) {
  if (state.phase !== "playing") return;
  const cell = tile.dataset.state;
  if (cell === "hidden" || (cell === "flagged" && tool === "flag")) {
    send(action[tool](x, y));
  }
}

function renderPlayers(showScores) {
  ui.players.replaceChildren(...state.players.map((player, index) => {
    const row = document.createElement("li");
    row.className = "player";
    if (player.id === state.playerId) row.classList.add("self");
    const name = document.createElement("strong");
    name.textContent = player.name;
    const details = document.createElement("span");
    const tags = [];
    if (player.id === state.playerId) tags.push("你");
    if (player.id === state.hostId) tags.push("房主");
    if (showScores) tags.push(`翻开 ${state.board.scores[index] ?? 0} 格`);
    details.textContent = tags.join(" · ");
    row.append(name, details);
    return row;
  }));
  ui.playerCount.textContent = `${state.players.length} 人`;
}

function render() {
  const phase = state.phase;
  const playing = phase === "playing";
  const finished = phase === "finished";
  const inLobby = phase === "lobby";
  const isHost = inLobby && state.hostId === state.playerId;
  const connected = state.socket?.readyState === WebSocket.OPEN;
  const showBoard = playing || finished;

  ui.connection.textContent = connected ? "已连接" : phase === "connecting" ? "连接中…" : "未连接";
  ui.connection.classList.toggle("online", connected);
  ui.joinForm.hidden = phase !== "idle";
  ui.leave.hidden = !connected;
  ui.lobby.hidden = showBoard;
  ui.board.hidden = !showBoard;
  ui.flagCount.hidden = !showBoard;
  ui.flags.textContent = state.board?.flagsLeft ?? "";
  ui.controls.hidden = !playing;
  ui.start.hidden = !isHost;
  ui.start.disabled = state.players.length < 2 || state.startPending;
  ui.start.textContent = state.players.length < 2 ? "等待搭档" : "开始对局";
  ui.rejoin.hidden = phase !== "waiting";
  ui.again.hidden = !finished;
  ui.again.disabled = !state.outcome?.canRematch;
  ui.reveal.setAttribute("aria-pressed", state.tool === "reveal");
  ui.flag.setAttribute("aria-pressed", state.tool === "flag");
  renderPlayers(showBoard);

  if (phase === "connecting") {
    ui.lobbyTitle.textContent = "正在加入大厅";
    ui.lobbyHint.textContent = "等待服务端确认连接。";
  } else if (inLobby) {
    ui.lobbyTitle.textContent = "已进入大厅";
    ui.lobbyHint.textContent = isHost ? "搭档加入后，就可以开局。" : "等待房主开始对局。";
  } else if (phase === "waiting") {
    ui.lobbyTitle.textContent = "对局正在进行";
    ui.lobbyHint.textContent = "这局结束后，可以加入下一局。";
  } else {
    ui.lobbyTitle.textContent = "找个搭档，开一局";
    ui.lobbyHint.textContent = "两个人共享一张棋盘，一起排雷。";
  }

  ui.noticeTitle.parentElement.dataset.outcome = state.outcome?.type ?? "";
  ui.noticeTitle.textContent = state.outcome?.heading ?? (playing ? "对局进行中" : "等你开局");
  ui.noticeText.textContent = state.notice || state.outcome?.text || (playing
    ? "蓝框是搭档当前指向的格子。"
    : "双方加入大厅后，由房主开始对局。");
}

ui.joinForm.addEventListener("submit", (event) => {
  event.preventDefault();
  const name = ui.name.value.trim();
  if (!name || new TextEncoder().encode(name).length > 15) {
    ui.name.setCustomValidity("请填昵称，最多 15 个英文字母或 5 个汉字。");
    ui.name.reportValidity();
    return;
  }
  connect(name);
});
ui.name.addEventListener("input", () => ui.name.setCustomValidity(""));
ui.start.addEventListener("click", () => {
  state.startPending = true;
  send(action.cooperative);
  render();
});
const rejoin = () => send(action.rejoin);
ui.rejoin.addEventListener("click", rejoin);
ui.again.addEventListener("click", rejoin);
ui.leave.addEventListener("click", () => state.socket?.close());
for (const tool of ["reveal", "flag"]) {
  ui[tool].addEventListener("click", () => {
    state.tool = tool;
    render();
  });
}

render();
