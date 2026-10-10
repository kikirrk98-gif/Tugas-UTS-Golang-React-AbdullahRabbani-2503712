# Catatan Go

Aplikasi catatan fullstack. Backend Go (Gin, GORM, SQLite, JWT), frontend React (Vite, React Router, Axios).

## Fitur
- Registrasi dan login (password di-hash bcrypt, sesi memakai JWT 24 jam)
- CRUD catatan milik masing-masing user, sematkan catatan, kategori
- Pencarian, filter kategori, paginasi, statistik
- Validasi di server, penanganan 401 otomatis di frontend

## Struktur
```
backend/    config, database, models, middleware, controllers, routes, utils, main.go
frontend/   src/{api.js, auth.jsx, App.jsx, pages, components, styles.css}
```

## Menjalankan (development)
Butuh Go 1.22+ dan Node 18+.

```sh
# terminal 1: API di :8080
cd backend
cp .env.example .env     # ubah JWT_SECRET
go mod tidy              # unduh dependensi dan buat go.sum
go run .

# terminal 2: React di :5173 (proxy /api ke :8080)
cd frontend
npm install
npm run dev
```
Buka http://localhost:5173

## Produksi (satu biner)
```sh
cd frontend && npm install && npm run build
cd ../backend && go build -o catatan . && ./catatan
```
Go menyajikan `frontend/dist` sendiri, jadi cukup buka http://localhost:8080.
Lokasi build bisa diubah dengan `STATIC_DIR`.

## Endpoint API
| Method | Path | Keterangan |
|---|---|---|
| POST | /api/auth/register | daftar, mengembalikan token |
| POST | /api/auth/login | login, mengembalikan token |
| GET | /api/me | profil (butuh token) |
| GET | /api/stats | total, disematkan, jumlah per kategori |
| GET | /api/notes?q=&category=&page=&limit= | daftar catatan |
| POST | /api/notes | buat catatan |
| GET/PUT/DELETE | /api/notes/:id | detail, ubah, hapus |
| PATCH | /api/notes/:id/pin | sematkan atau lepas |

Header untuk endpoint terlindungi: `Authorization: Bearer <token>`.
