import { useCallback, useEffect, useState } from "react";
import api, { errMsg } from "../api";
import NoteCard from "../components/NoteCard";
import NoteModal from "../components/NoteModal";

export default function Notes() {
    const [notes, setNotes] = useState([]);
  const [meta, setMeta] = useState({ page: 1, pages: 0, total: 0 });
  const [stats, setStats] = useState({ total: 0, pinned: 0, categories: [] });
  const [q, setQ] = useState("");
  const [search, setSearch] = useState("");
  const [category, setCategory] = useState("");
  const [page, setPage] = useState(1);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [modal, setModal] = useState(null); // null | {note?: Note}
  const [saving, setSaving] = useState(false);
  const [formErr, setFormErr] = useState("");
  const [confirmId, setConfirmId] = useState(null);

  // Debounce pencarian agar tidak memanggil API di setiap ketikan.
  useEffect(() => {
    const t = setTimeout(() => { setSearch(q); setPage(1); }, 300);
    return () => clearTimeout(t);
  }, [q]);

  const load = useCallback(async () => {
    setLoading(true);
    try {
      const [n, s] = await Promise.all([
        api.get("/notes", { params: { q: search, category, page, limit: 9 } }),
        api.get("/stats"),
      ]);
      setNotes(n.data.data);
      setMeta(n.data.meta);
      setStats(s.data.data);
      setError("");
    } catch (err) {
      setError(errMsg(err));
    } finally {
      setLoading(false);
    }
  }, [search, category, page]);

  useEffect(() => { load(); }, [load]);

  const save = async (form) => {
    setSaving(true);
    setFormErr("");
    try {
      if (modal.note) await api.put(`/notes/${modal.note.id}`, form);
      else await api.post("/notes", form);
      setModal(null);
      load();
    } catch (err) {
      setFormErr(errMsg(err));
    } finally {
      setSaving(false);
    }
  };

  const pin = async (note) => { await api.patch(`/notes/${note.id}/pin`); load(); };
  const remove = async (note) => {
    await api.delete(`/notes/${note.id}`);
    setConfirmId(null);
    // Jika halaman jadi kosong, mundur satu halaman.
    if (notes.length === 1 && page > 1) setPage(page - 1);
    else load();
  };

  return (
    <section className="section-notes">
      <h2>Catatan</h2>
      <p className="statline">
        total: <b>{stats.total}</b> · disematkan: <b>{stats.pinned}</b>
      </p>

      <div className="filters" role="group" aria-label="Filter kategori">
        <button className={"chipbtn" + (category === "" ? " on" : "")}
          onClick={() => { setCategory(""); setPage(1); }}>
          Semua <small>{stats.total}</small>
        </button>
        {stats.categories.map((c) => (
          <button key={c.category} className={"chipbtn" + (category === c.category ? " on" : "")}
            onClick={() => { setCategory(c.category); setPage(1); }}>
            {c.category} <small>{c.total}</small>
          </button>
        ))}
      </div>

      <div className="tools">
        <input type="search" id="cari" aria-label="Cari catatan" placeholder="cari judul atau isi..."
          value={q} onChange={(e) => setQ(e.target.value)} />
        <button className="btn primary" onClick={() => { setFormErr(""); setModal({}); }}>+ Catatan baru</button>
      </div>

      {error && <p className="error" role="alert">{error}</p>}

      {!loading && !error && notes.length === 0 && (
        <div className="empty">
          <h3>{search || category ? "Tidak ada catatan yang cocok" : "Belum ada catatan"}</h3>
          <p>{search || category ? "Coba kata kunci atau kategori lain." : "Klik “Catatan baru” untuk menulis yang pertama."}</p>
        </div>
      )}

      <div className={"grid" + (loading ? " loading" : "")}>
        {notes.map((n) => (
          <NoteCard key={n.id} note={n}
            onEdit={(note) => { setFormErr(""); setModal({ note }); }}
            onPin={pin} onDelete={remove}
            confirming={confirmId === n.id}
            onConfirm={setConfirmId} onCancel={() => setConfirmId(null)} />
        ))}
      </div>

      {meta.pages > 1 && (
        <nav className="pager" aria-label="Halaman">
          <button className="btn" disabled={page <= 1} onClick={() => setPage(page - 1)}>← Sebelumnya</button>
          <span>Halaman {meta.page} dari {meta.pages}</span>
          <button className="btn" disabled={page >= meta.pages} onClick={() => setPage(page + 1)}>Berikutnya →</button>
        </nav>
      )}

      {modal && (
        <NoteModal note={modal.note} onSave={save} onClose={() => setModal(null)} error={formErr} busy={saving} />
      )}
    </section>
  );
}
