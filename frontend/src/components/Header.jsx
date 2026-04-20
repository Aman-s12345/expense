import { useAuth } from "../context/AuthContext";
import s from "./Header.module.css";

export default function Header() {
  const { user, logout } = useAuth();

  return (
    <header className={s.bar}>
      <span className={s.logo}>Expense Tracker</span>
      <div className={s.right}>
        <span className={s.name}>{user.is_guest ? "Guest" : user.name}</span>
        <button className={s.btn} onClick={logout}>Logout</button>
      </div>
    </header>
  );
}
