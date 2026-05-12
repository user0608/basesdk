export const toStringArray = (value: unknown): string[] => {
  if (value == null) return [];
  if (Array.isArray(value)) {
    return value.flatMap(toStringArray);
  }
  if (value instanceof Set) {
    return [...value].flatMap(toStringArray);
  }
  if (typeof value === "string") {
    const normalizedValue = value.trim();

    return normalizedValue ? [normalizedValue] : [];
  }
  if (typeof value === "number") {
    return Number.isFinite(value) ? [String(value)] : [];
  }
  if (typeof value === "boolean" || typeof value === "bigint") {
    return [String(value)];
  }
  return [];
};
