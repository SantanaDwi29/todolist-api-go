# To-Do List API (Backend)

## 1. Tentang API Ini
Ini adalah layanan Backend berbasis RESTful API untuk aplikasi To-Do List. Dibangun menggunakan performa tinggi dari Golang, API ini bertanggung jawab untuk menangani semua logika bisnis, autentikasi pengguna secara aman, serta seluruh operasi database untuk entitas aplikasi seperti User, Todo, Category, FocusSession, dan Milestone.

## 2. Teknologi yang Digunakan
- **Bahasa Pemrograman:** Golang (v1.26)
- **Framework Web:** Gin Gonic
- **Database ORM:** GORM
- **Database Relasional:** MySQL
- **Autentikasi:** JSON Web Token (JWT) dengan hashing Bcrypt
- **Migrasi Database:** golang-migrate

## 3. Struktur Project Backend
Proyek ini mengadopsi struktur standar yang rapi (Mirip Hexagonal/Layered Architecture). Berikut adalah penjelasan setiap folder di dalam `todolist-api`:

- **`/cmd/server/`**
  - **Fungsi:** Titik masuk utama (entry point) untuk mengeksekusi aplikasi.
  - Berisi `main.go` yang memuat konfigurasi, menghubungkan database, menyambungkan rute, dan menjalankan HTTP server.
- **`/internal/config/`**
  - **Fungsi:** Memuat dan mengelola variabel lingkungan (environment variables) seperti credential DB dan JWT Secret.
- **`/internal/database/`**
  - **Fungsi:** Menyimpan skrip koneksi GORM ke database MySQL.
- **`/internal/models/`**
  - **Fungsi:** Menyimpan struktur data (structs) Golang yang mendefinisikan representasi tabel di MySQL (contoh: `user.go`, `todo.go`, `category.go`).
- **`/internal/repository/`**
  - **Fungsi:** Layer komunikasi langsung dengan Database. Disini kueri (Create, Read, Update, Delete) GORM dieksekusi agar logika kueri terpisah dari logika bisnis.
- **`/internal/service/`**
  - **Fungsi:** Layer logika bisnis (Business Logic). Menerima instruksi dari handler, melakukan validasi/algoritma, kemudian memanggil repository terkait.
- **`/internal/handler/`** (atau Controllers)
  - **Fungsi:** Layer penerima request HTTP. Mengambil JSON/Form payload, memvalidasi formatnya, meneruskannya ke Service, lalu mengirimkan respons JSON kembali ke client.
- **`/internal/routes/`**
  - **Fungsi:** Menyimpan `routes.go` yang memetakan URL API (misal `POST /api/v1/todos`) ke Handler yang spesifik.
- **`/internal/middleware/`**
  - **Fungsi:** Lapisan penengah HTTP (interceptors). Berisi fungsi validasi Token JWT (Auth Middleware) untuk mencegah akses tidak sah, dan konfigurasi CORS.
- **`/migrations/`**
  - **Fungsi:** File `.sql` (Up dan Down) untuk membuat dan memperbarui skema tabel secara otomatis (version control untuk database).
- **`/utils/`**
  - **Fungsi:** Berisi fungsi utilitas serbaguna, contoh: format response JSON standar, atau fungsi hashing password.

## 4. Aliran Data Backend
**Dari Mana ke Mana (Data Flow):**

1. **Client Request:** Frontend (browser) mengirimkan request menuju ke endpoint API tertentu (contoh: `POST /api/todos`).
2. **Routes & Middleware:** Request HTTP tiba di server Go. Rute (`/internal/routes/`) mencegatnya dan melewatkannya ke Middleware (`/internal/middleware/`) untuk dicek token JWT-nya.
3. **Handler Layer:** Jika otentikasi lolos, request dilempar ke Handler (`/internal/handler/`). Handler menerjemahkan JSON menjadi object Golang.
4. **Service Layer:** Handler mengirimkan data ke Service (`/internal/service/`) yang memastikan kebenaran proses bisnis (contoh: apakah kategori milik pengguna ini?).
5. **Repository Layer:** Jika logika valid, Service menginstruksikan Repository (`/internal/repository/`) untuk menyimpan data.
6. **Database:** Repository mengeksekusi sintaks MySQL melalui GORM.
7. **Response Flow:** Database membalikkan data tersimpan -> Repository -> Service -> Handler. Handler merangkum data tersebut menjadi response standar JSON sukses (beserta kode status HTTP 200/201) lalu melemparkannya kembali ke Frontend.

## 5. Format Response Standar
Semua endpoint di API ini diatur menggunakan utilitas response standar (dikelola di `utils/response.go` atau dikerjakan secara konsisten di layer Handler) agar bentuk JSON yang dikembalikan selalu terprediksi. Hal ini sangat mempermudah konsumsi data oleh Frontend.

**Contoh Response Sukses (HTTP 200 OK / 201 Created):**
```json
{
  "status": "success",
  "message": "Berhasil melakukan aksi (misal: data berhasil diambil)",
  "data": { 
     "id": 1,
     "title": "Tugas Contoh"
  }
}
```

**Contoh Response Error (HTTP 400 / 401 / 404 / 500):**
```json
{
  "status": "error",
  "message": "Pesan error spesifik (misal: Kredensial tidak valid)",
  "data": null
}
```

## 6. Langkah-Langkah Menjalankan & Migrasi Database
Untuk menjalankan Backend ini secara mandiri, ikuti langkah-langkah berikut:

**Langkah 1: Setup Database MySQL**
1. Buat database baru di MySQL server Anda dengan menggunakan perintah SQL:
   ```sql
   CREATE DATABASE todolist_db;
   ```

**Langkah 2: Konfigurasi Environment**
1. Salin file `.env.example` ke `.env` di dalam root folder `todolist-api`.
2. Buka `.env` dan sesuaikan kredensial koneksi agar sesuai dengan MySQL Anda:
   ```env
   DB_HOST=127.0.0.1
   DB_PORT=3306
   DB_USER=root
   DB_PASSWORD=
   DB_NAME=todolist_db
   PORT=8080
   JWT_SECRET=rahasia_jwt_anda
   ```

**Langkah 3: Menjalankan Server & Migrasi Otomatis**
- Buka terminal di direktori `todolist-api` dan jalankan:
  ```bash
  go mod tidy
  go run cmd/server/main.go
  ```
- **Proses Migrasi Database:** Saat file `main.go` dieksekusi, salah satu fungsi awal yang dipanggil adalah `database.MigrateDB()`. Fungsi ini secara otomatis mendeteksi dan mengeksekusi urutan file `*.up.sql` dari dalam folder `/migrations`. Ini berarti seluruh skema tabel aplikasi (Users, Todos, dsb.) akan langsung dibuat secara otomatis di database MySQL Anda tanpa perlu impor `.sql` secara manual!
- **Tanda Berhasil:** Di terminal akan muncul log:
  ```text
  Database connected successfully!
  Database SQL migration completed successfully!
  Server starting on port 8080
  ```

## 7. Command Line Interface (CLI)

Proyek ini juga dilengkapi dengan beberapa tool CLI (Command Line Interface) yang fungsionalitasnya mirip dengan Artisan di Laravel.

### Menjalankan Migrasi Secara Terpisah
Jika Anda sedang menjalankan server dan tidak ingin mematikannya untuk menjalankan file migrasi baru (mirip seperti `php artisan migrate`), Anda dapat membuka tab terminal baru dan mengeksekusi:
```bash
go run cmd/migrate/main.go
```

### Membuat OAuth Client (Untuk Frontend)
Untuk menghasilkan Client ID dan Client Secret yang akan digunakan oleh aplikasi Frontend Anda (mirip seperti `php artisan passport:client`), eksekusi perintah berikut di terminal:
```bash
go run cmd/passport/main.go --name="Nama Frontend App Anda"
```
Setelah dieksekusi, terminal akan menampilkan `Client ID` dan `Client Secret` yang bisa langsung Anda gunakan di sisi Frontend.
