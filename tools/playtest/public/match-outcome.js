export function matchOutcome(match) {
  if (match.Cancel) return { label: "已取消", state: "cancel" };
  if (match.Won) return { label: "排雷成功", state: "win" };
  return { label: "挑战失败", state: "lose" };
}
