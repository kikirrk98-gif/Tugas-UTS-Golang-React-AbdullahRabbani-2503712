import { NavLink, Outlet, Link } from "react-router-dom";
import { useAuth } from "../auth";

export function Layout() {
  const { user, logout } = useAuth();
  const items = user
    ? [{ to: "/", label: "Catatan", end: true }]
    : [{ to: "/masuk", label: "Masuk" }, { to: "/daftar", label: "Daftar" }];

  return (
    <div className="app-container">
      <header className="site-header">
        <h1>Catatan Go</h1>
        <p className="role">
          Aplikasi catatan pribadi dengan backend Go (Gin, GORM, JWT) dan frontend React.
        </p>
        <p className="tagline">fokus: crud.catatan &amp; autentikasi.jwt</p>
      </header>

      <nav className="site-nav" aria-label="Navigasi utama">
        <ul>
          {items.map((it) => (
            <li key={it.to}>
              <NavLink to={it.to} end={it.end}>
                <span className="port-dot" aria-hidden="true" />
                {it.label}
              </NavLink>
            </li>
          ))}
          {user && (
            <li className="push">
              <button className="nav-btn" onClick={logout}>
                <span className="port-dot" aria-hidden="true" />
                Keluar ({user.name})
              </button>
            </li>
          )}
        </ul>
      </nav>

      <main className="site-content">
        <Outlet />
      </main>

      <footer className="site-footer">
        <p>© 2026 Catatan Go — dibangun dengan Go &amp; React</p>
      </footer>
    </div>
  );
}

export function NotFound() {
  return (
    <section className="section-404">
      <h2>404</h2>
      <p>Halaman tidak ditemukan.</p>
      <Link to="/">Kembali ke Beranda</Link>
    </section>
  );
}
