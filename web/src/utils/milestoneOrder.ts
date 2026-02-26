export function orderMilestones(milestones: string[], order: string[]): string[] {
  const ordered: string[] = [];
  // First: milestones that appear in the saved order
  for (const m of order) {
    if (milestones.includes(m)) ordered.push(m);
  }
  // Then: milestones not in the saved order (alphabetical), empty string last
  const remaining = milestones.filter((m) => !order.includes(m) && m !== "");
  remaining.sort((a, b) => a.localeCompare(b));
  ordered.push(...remaining);
  // Empty string always last
  if (milestones.includes("")) ordered.push("");
  return ordered;
}
