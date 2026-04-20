import { useState, useEffect, useCallback } from "react";
import { api } from "../lib/api";

export function useExpenses() {
  const [expenses, setExpenses] = useState([]);
  const [totalPaisa, setTotalPaisa] = useState(0);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [category, setCategory] = useState("");
  const [sort, setSort] = useState("date_desc");

  const fetch_ = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const params = new URLSearchParams();
      if (category) params.set("category", category);
      params.set("sort", sort);
      const res = await api.get(`/expenses/?${params}`);
      setExpenses(res.data.expenses);
      setTotalPaisa(res.data.total_paisa);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }, [category, sort]);

  useEffect(() => { fetch_(); }, [fetch_]);

  const toggleSort = useCallback(() => {
    setSort((s) => (s === "date_desc" ? "date_asc" : "date_desc"));
  }, []);

  return { expenses, totalPaisa, loading, error, category, setCategory, sort, toggleSort, refetch: fetch_ };
}
