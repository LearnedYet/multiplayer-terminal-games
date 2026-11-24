import { matchOutcome } from "./match-outcome.js";

const pageSize = 8;
const requestedPage = Number(new URLSearchParams(location.search).get("page"));
const validPage = Number.isSafeInteger(requestedPage) && requestedPage > 0 && Number.isSafeInteger(requestedPage * pageSize);
const page = validPage ? requestedPage : 1;
const status = document.getElementById("list-status");
const list = document.getElementById("match-list");
const navigation = document.getElementById("page-nav");
const previous = document.getElementById("prev-page");
const next = document.getElementById("next-page");

function matchRow(match) {
  const item = document.createElement("li");
  const link = document.createElement("a");
  link.className = "match-link";
  link.href = "/result.html?" + new URLSearchParams({ id: match.GameID, page });

  const outcome = matchOutcome(match);
  const badge = document.createElement("strong");
  badge.className = "record-outcome";
  badge.dataset.outcome = outcome.state;
  badge.textContent = outcome.label;

  const id = document.createElement("span");
  id.className = "match-id";
  id.textContent = match.GameID;

  const arrow = document.createElement("span");
  arrow.className = "match-arrow";
  arrow.setAttribute("aria-hidden", "true");
  arrow.textContent = "›";
  link.append(badge, id, arrow);
  item.append(link);
  return item;
}

async function loadPage() {
  // 多取一条来确定是否还有下一页，列表仍只显示本页的条数。
  const query = new URLSearchParams({ limit: pageSize + 1, offset: (page - 1) * pageSize });
  try {
    const response = await fetch("/matches?" + query);
    if (response.status === 404 || response.status === 501) {
      status.textContent = "当前试玩服务还未开放战绩列表。";
      return;
    }
    if (!response.ok) throw new Error(`HTTP ${response.status}`);

    const matches = await response.json();
    list.replaceChildren(...matches.slice(0, pageSize).map(matchRow));
    status.textContent = matches.length ? "" : page === 1 ? "还没有已保存的对局。" : "这一页没有对局。";

    if (page > 1) previous.href = "?page=" + (page - 1);
    if (matches.length > pageSize) next.href = "?page=" + (page + 1);
    document.getElementById("page-number").textContent = `第 ${page} 页`;
    navigation.hidden = !matches.length && page === 1;
  } catch (error) {
    console.error("读取战绩列表失败", error);
    status.textContent = "战绩读取失败，请刷新重试。";
  }
}

loadPage();
