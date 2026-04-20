import { useState, useRef, useCallback } from "react";
import { api } from "../lib/api";
import { generateId } from "../lib/utils";

export function useCreateExpense(onSuccess) {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const key = useRef(generateId());

  const create = useCallback(async ({ amount, category, description, date }) => {
    setError(null);

    const paisa = Math.round(parseFloat(amount) * 100);
    if (isNaN(paisa) || paisa <= 0) { setError("Amount must be positive"); return false; }
    if (!date) { setError("Date is required"); return false; }
    if (!category) { setError("Category is required"); return false; }

    setLoading(true);
    try {
      await api.post("/expenses/", {
        amount: paisa,
        category,
        description: description.trim(),
        date,
        idempotency_key: key.current,
      });
      key.current = generateId();
      onSuccess?.();
      return true;
    } catch (err) {
      setError(err.message);
      return false;
    } finally {
      setLoading(false);
    }
  }, [onSuccess]);

  return { create, loading, error };
}
