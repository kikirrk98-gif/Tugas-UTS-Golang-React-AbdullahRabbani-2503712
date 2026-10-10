const tgl = (s) =>
  new Date(s).toLocaleDateString("id-ID", { day: "numeric", month: "short", year: "numeric" });

export default function NoteCard({ note, onEdit, onPin, onDelete, confirming, onConfirm, onCancel }) {
  return (
    <article className={"card" + (note.pinned ? " pinned" : "")}>
      <header>
        <span className="chip">{note.category}</span>
        <button className={"icon" + (note.pinned ? " on" : "")} onClick={() => onPin(note)}
          aria-label={note.pinned ? "Lepas sematan" : "Sematkan"} title={note.pinned ? "Lepas sematan" : "Sematkan"}>
          {note.pinned ? "★" : "☆"}
        </button>
      </header>
      <h3>{note.title}</h3>
      <p>{note.content || <em>Tanpa isi</em>}</p>
      <footer>
        <time dateTime={note.updated_at}>{tgl(note.updated_at)}</time>
        {confirming ? (
          <span className="row">
            <button className="link danger" onClick={() => onDelete(note)}>Ya, hapus</button>
            <button className="link" onClick={onCancel}>Batal</button>
          </span>
        ) : (
          <span className="row">
            <button className="link" onClick={() => onEdit(note)}>Ubah</button>
            <button className="link danger" onClick={() => onConfirm(note.id)}>Hapus</button>
          </span>
        )}
      </footer>
    </article>
  );
}
