export function formatRupees(paisa) {
  const n = Number(paisa);
  if (isNaN(n)) return "₹0.00";
  const rupees = Math.floor(Math.abs(n) / 100);
  const p = Math.abs(n) % 100;
  return `₹${rupees.toLocaleString("en-IN")}.${String(p).padStart(2, "0")}`;
}

export function generateId() {
  return "xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx".replace(/[xy]/g, (c) => {
    const r = (Math.random() * 16) | 0;
    return (c === "x" ? r : (r & 0x3) | 0x8).toString(16);
  });
}

export function todayISO() {
  return new Date().toISOString().slice(0, 10);
}

export const CATEGORIES = ["Food", "Transport", "Shopping", "Bills", "Entertainment", "Health", "Education", "Other"];
