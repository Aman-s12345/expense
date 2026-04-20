import { useState } from "react";
import { useCreateExpense } from "../hooks/useCreateExpense";
import { CATEGORIES, todayISO } from "../lib/utils";
import s from "./AddExpenseForm.module.css";

export default function AddExpenseForm({ onAdded }) {
  const [amount, setAmount] = useState("");
  const [category, setCategory] = useState("Food");
  const [description, setDescription] = useState("");
  const [date, setDate] = useState(todayISO());
  const [success, setSuccess] = useState("");

  const { create, loading, error } = useCreateExpense(() => {
    setAmount("");
    setDescription("");
    setDate(todayISO());
    setSuccess("Expense added!");
    setTimeout(() => setSuccess(""), 2000);
    onAdded();
  });

  const handleSubmit = (e) => {
    e.preventDefault();
    create({ amount, category, description, date });
  };

  return (
    <div className={s.card}>
      <h2 className={s.title}>Add Expense</h2>
      <form onSubmit={handleSubmit}>
        <div className={s.row}>
          <label className={s.field}>
            <span className={s.label}>Amount (₹)</span>
            <input className={s.input} type="number" step="0.01" min="0.01" placeholder="0.00"
              value={amount} onChange={(e) => setAmount(e.target.value)} required />
          </label>
          <label className={s.field}>
            <span className={s.label}>Category</span>
            <select className={s.input} value={category} onChange={(e) => setCategory(e.target.value)}>
              {CATEGORIES.map((c) => <option key={c} value={c}>{c}</option>)}
            </select>
          </label>
        </div>
        <div className={s.row}>
          <label className={s.field}>
            <span className={s.label}>Description</span>
            <input className={s.input} type="text" placeholder="What was this for?"
              value={description} onChange={(e) => setDescription(e.target.value)} />
          </label>
          <label className={s.field}>
            <span className={s.label}>Date</span>
            <input className={s.input} type="date" value={date}
              onChange={(e) => setDate(e.target.value)} required />
          </label>
        </div>
        {error && <p className={s.error}>{error}</p>}
        {success && <p className={s.success}>{success}</p>}
        <button className={s.submit} type="submit" disabled={loading}>
          {loading ? "Saving..." : "Add Expense"}
        </button>
      </form>
    </div>
  );
}
