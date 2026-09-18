// Formats an integer amount of minor units (cents) as a price string.
export function formatPrice(cents) {
  return `Rs.${((cents ?? 0) / 100).toFixed(2)}`;
}
