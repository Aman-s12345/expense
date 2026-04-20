import { CATEGORIES, formatRupees } from "../lib/utils";
import s from "./ExpenseList.module.css";

export default function ExpenseList({
  expenses, totalPaisa, loading, error,
  category, setCategory, sort, toggleSort, onRetry,
}) {
  return (
    <div className={s.card}>
      <div className={s.header}>
        <h2 className={s.title}>Expenses</h2>
        <span className={s.total}>Total: {formatRupees(totalPaisa)}</span>
      </div>

      <div className={s.filters}>
        <select className={s.select} value={category} onChange={(e) => setCategory(e.target.value)}>
          <option value="">All Categories</option>
          {CATEGORIES.map((c) => <option key={c} value={c}>{c}</option>)}
        </select>
        <button className={s.sortBtn} onClick={toggleSort}>
          Date {sort === "date_desc" ? "↓" : "↑"}
        </button>
      </div>

      {error && (
        <div className={s.errorWrap}>
          <p className={s.error}>{error}</p>
          <button className={s.retryBtn} onClick={onRetry}>Retry</button>
        </div>
      )}

      {loading ? (
        <p className={s.empty}>Loading expenses...</p>
      ) : expenses.length === 0 ? (
        <p className={s.empty}>
          {category ? `No expenses in "${category}".` : "No expenses yet. Add one above!"}
        </p>
      ) : (
        <div className={s.table}>
          <div className={s.thead}>
            <span className={s.cDate}>Date</span>
            <span className={s.cDesc}>Description</span>
            <span className={s.cCat}>Category</span>
            <span className={s.cAmt}>Amount</span>
          </div>
          {expenses.map((e) => (
            <div key={e.id} className={s.row}>
              <span className={s.cDate}>{e.date}</span>
              <span className={s.cDesc}>{e.description || "—"}</span>
              <span className={s.cCat}><span className={s.pill}>{e.category}</span></span>
              <span className={s.cAmt}>{formatRupees(e.amount_paisa)}</span>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
