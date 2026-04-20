import { useState } from "react";
import { useAuth } from "../context/AuthContext";
import s from "./LoginPage.module.css";

export default function LoginPage() {
  const { login, register, guestLogin } = useAuth();
  const [mode, setMode] = useState("login");
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [guestLoading, setGuestLoading] = useState(false);
  const [error, setError] = useState("");

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError("");
    setLoading(true);
    try {
      if (mode === "login") {
        await login(email, password);
      } else {
        await register(name, email, password);
      }
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const handleGuest = async () => {
    setError("");
    setGuestLoading(true);
    try {
      await guestLogin();
    } catch (err) {
      setError(err.message);
    } finally {
      setGuestLoading(false);
    }
  };

  return (
    <div className={s.wrapper}>
      <div className={s.card}>
        <h1 className={s.heading}>Expense Tracker</h1>
        <p className={s.sub}>Track where your money goes</p>

        <div className={s.tabs}>
          <button className={mode === "login" ? s.tabActive : s.tab}
            onClick={() => { setMode("login"); setError(""); }}>Login</button>
          <button className={mode === "register" ? s.tabActive : s.tab}
            onClick={() => { setMode("register"); setError(""); }}>Register</button>
        </div>

        <form onSubmit={handleSubmit}>
          {mode === "register" && (
            <input className={s.input} type="text" placeholder="Name"
              value={name} onChange={(e) => setName(e.target.value)} required />
          )}
          <input className={s.input} type="email" placeholder="Email"
            value={email} onChange={(e) => setEmail(e.target.value)} required />
          <input className={s.input} type="password" placeholder="Password (min 6 chars)"
            value={password} onChange={(e) => setPassword(e.target.value)} required minLength={6} />
          {error && <p className={s.error}>{error}</p>}
          <button className={s.btnPrimary} type="submit" disabled={loading}>
            {loading ? "Please wait..." : mode === "login" ? "Login" : "Create Account"}
          </button>
        </form>

        <div className={s.divider}><span className={s.dividerText}>or</span></div>

        <button className={s.btnGhost} onClick={handleGuest} disabled={guestLoading}>
          {guestLoading ? "Creating..." : "Continue as Guest"}
        </button>
      </div>
    </div>
  );
}
