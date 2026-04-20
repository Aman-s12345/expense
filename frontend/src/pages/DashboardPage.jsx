import Header from "../components/Header";
import AddExpenseForm from "../components/AddExpenseForm";
import ExpenseList from "../components/ExpenseList";
import CategorySummary from "../components/CategorySummary";
import { useExpenses } from "../hooks/useExpenses";
import s from "./DashboardPage.module.css";

export default function DashboardPage() {
  const {
    expenses, totalPaisa, loading, error,
    category, setCategory, sort, toggleSort, refetch,
  } = useExpenses();

  return (
    <div>
      <Header />
      <main className={s.main}>
        <AddExpenseForm onAdded={refetch} />
        <div className={s.columns}>
          <div className={s.left}>
            <ExpenseList
              expenses={expenses}
              totalPaisa={totalPaisa}
              loading={loading}
              error={error}
              category={category}
              setCategory={setCategory}
              sort={sort}
              toggleSort={toggleSort}
              onRetry={refetch}
            />
          </div>
          <div className={s.right}>
            <CategorySummary expenses={expenses} />
          </div>
        </div>
      </main>
    </div>
  );
}
