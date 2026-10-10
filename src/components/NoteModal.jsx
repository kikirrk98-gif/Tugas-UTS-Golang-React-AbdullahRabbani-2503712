import { useEffect, useRef, useState } from "react";

const SARAN = ["Umum", "Kuliah", "Kerja", "Ide", "Belanja"];

export default function NoteModal({ note, onSave, onClose, error, busy }) {
  const [form, setForm] = useState({
    title: note?.title ?? "",
    content: note?.content ?? "",
    category: note?.category ?? "Umum",
    pinned: note?.pinned ?? false,
  });
  const first = useRef(null);

  useEffect(() => {
    first.current?.focus();
    const onKey = (e) => e.key === "Escape" && onClose();
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  }, [onClose]);

  const submit = (e) => {
    e.preventDefault();
    onSave(form);
  };

  return (
    <div className="overlay" onMouseDown={(e) => e.target === e.currentTarget && onClose()}>
      <form className="modal" onSubmit={submit} role="dialog" aria-modal="true" aria-labelledby="mt">
        <h2 id="mt">{note ? "Ubah catatan" : "Catatan baru"}</h2>
        <label htmlFor="n-title">Judul
          <input id="n-title" ref={first} value={form.title} maxLength={200}
            onChange={(e) => setForm({ ...form, title: e.target.value })} required />
        </label>
        <label htmlFor="n-content">Isi
          <textarea id="n-content" rows={7} value={form.content}
            onChange={(e) => setForm({ ...form, content: e.target.value })} />
        </label>
        <div className="two">
          <label htmlFor="n-cat">Kategori
            <input id="n-cat" list="saran" value={form.category} maxLength={50}
              onChange={(e) => setForm({ ...form, category: e.target.value })} />
            <datalist id="saran">{SARAN.map((s) => <option key={s} value={s} />)}</datalist>
          </label>
          <label className="check" htmlFor="n-pin">
            <input id="n-pin" type="checkbox" checked={form.pinned}
              onChange={(e) => setForm({ ...form, pinned: e.target.checked })} />
            Sematkan di atas
          </label>
        </div>
        {error && <p className="error" role="alert">{error}</p>}
        <div className="actions">
          <button type="button" className="btn" onClick={onClose}>Batal</button>
          <button className="btn primary" disabled={busy}>{busy ? "Menyimpan..." : "Simpan"}</button>
        </div>
      </form>
    </div>
  );
}
