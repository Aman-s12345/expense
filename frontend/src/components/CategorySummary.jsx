import { formatRupees } from "../lib/utils";
import s from "./CategorySummary.module.css";

export default function CategorySummary({ expenses }) {
  const totals = {};
  expenses.forEach((e) => {
    totals[e.category] = (totals[e.category] || 0) + e.amount_paisa;
  });

  const sorted = Object.entries(totals).sort((a, b) => b[1] - a[1]);
  if (sorted.length === 0) return null;

  const max = sorted[0][1];

  return (
    <div className={s.card}>
      <h2 className={s.title}>By Category</h2>
      {sorted.map(([cat, total]) => (
        <div key={cat} className={s.item}>
          <div className={s.info}>
            <span>{cat}</span>
            <span className={s.amount}>{formatRupees(total)}</span>
          </div>
          <div className={s.barBg}>
            <div className={s.barFill} style={{ width: `${(total / max) * 100}%` }} />
          </div>
        </div>
      ))}
    </div>
  );
}
