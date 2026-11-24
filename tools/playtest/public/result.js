import { matchOutcome } from "./match-outcome.js";

const form = document.getElementById("query-form");
const input = document.getElementById("match-id");
const status = document.getElementById("query-status");
const record = document.getElementById("record");
const lookup = document.getElementById("lookup");
const stats = document.getElementById("stats");
const submit = form.querySelector("button");

const query = new URLSearchParams(location.search);
const page = query.get("page");
input.value = (query.get("id") ?? "").trim();
if (page) {
  document.getElementById("history-link").href = "/history.html?" + new URLSearchParams({ page });
  lookup.hidden = Boolean(input.value);
}

function displayResult(match) {
  const outcome = matchOutcome(match);
  const badge = document.getElementById("outcome");
  badge.textContent = outcome.label;
  badge.dataset.outcome = outcome.state;
  document.getElementById("record-id").textContent = match.GameID;
  const rows = (match.PlayerResult ?? []).map((player) => {
    const row = document.createElement("tr");
    for (const value of [player.Name, player.Reveals, player.Flags, player.Left ? "中途离开" : "留到结算"]) {
      const cell = row.insertCell();
      cell.textContent = value;
    }
    return row;
  });
  stats.replaceChildren(...rows);
  record.hidden = false;
}

async function queryMatch() {
  const id = input.value.trim();
  if (!id) {
    status.textContent = "请输入对局编号。";
    return;
  }

  record.hidden = true;
  submit.disabled = true;
  status.textContent = "正在读取战绩…";
  try {
    const response = await fetch("/matches/" + encodeURIComponent(id));
    if (response.status === 404) {
      status.textContent = "没查到这局。可以检查对局编号。";
      return;
    }
    if (!response.ok) {
      status.textContent = "查询失败（HTTP " + response.status + "）。";
      return;
    }
    displayResult(await response.json());
    status.textContent = "";
    const url = new URL(location.href);
    url.searchParams.set("id", id);
    history.replaceState(null, "", url);
  } catch (error) {
    console.error("读取单局战绩失败", error);
    status.textContent = "读取失败，请确认试玩服务正在运行。";
  } finally {
    submit.disabled = false;
    if (record.hidden) lookup.hidden = false;
  }
}

form.addEventListener("submit", (event) => {
  event.preventDefault();
  queryMatch();
});

if (input.value) queryMatch();
