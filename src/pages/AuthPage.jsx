import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../auth";
import { errMsg } from "../api";

export default function AuthPage({ mode }) {
  const isLogin = mode === "login";
  const { login, register } = useAuth();
  const nav = useNavigate();
  const [form, setForm] = useState({ name: "", email: "", password: "" });
  const [error, setError] = useState("");
  const [busy, setBusy] = useState(false);

  const set = (k) => (e) => setForm({ ...form, [k]: e.target.value });

  const submit = async (e) => {
    e.preventDefault();
    setError("");
    setBusy(true);
    try {
      if (isLogin) await login(form.email, form.password);
      else await register(form.name, form.email, form.password);
      nav("/", { replace: true });
    } catch (err) {
      setError(errMsg(err));
    } finally {
      setBusy(false);
    }
  };

  return (
    <section className="section-auth">
      <h2>{isLogin ? "Masuk" : "Buat akun"}</h2>
      <p className="lead">
        {isLogin
          ? "Masuk untuk membuka catatan Anda. Sesi memakai token JWT berlaku 24 jam."
          : "Daftar untuk mulai menulis catatan. Password disimpan sebagai hash bcrypt."}
      </p>
      <form className="form" onSubmit={submit} noValidate>
        {!isLogin && (
          <label htmlFor="name">nama
            <input id="name" value={form.name} onChange={set("name")} autoComplete="name" required />
          </label>
        )}
        <label htmlFor="email">email
          <input id="email" type="email" value={form.email} onChange={set("email")} autoComplete="email" required />
        </label>
        <label htmlFor="password">password
          <input id="password" type="password" value={form.password} onChange={set("password")}
            autoComplete={isLogin ? "current-password" : "new-password"} minLength={6} required />
        </label>
        {error && <p className="error" role="alert">{error}</p>}
        <button className="btn primary" disabled={busy}>
          {busy ? "memproses..." : isLogin ? "Masuk" : "Daftar"}
        </button>
      </form>
    </section>
  );
}
