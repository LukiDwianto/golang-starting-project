# Golang Starter Project

Ini adalah proyek starter untuk aplikasi Golang menggunakan berbagai teknologi modern untuk pengembangan dan deployment.

## Teknologi yang Digunakan

- **Gin**: Framework web minimalis dan sangat cepat untuk Golang.
- **PostgreSQL**: Database relasional yang andal dan kaya fitur.
- **Docker Compose**: Mengelola dan mengorkestrasi layanan dalam container Docker.
- **Air**: Alat untuk live-reloading dalam pengembangan Golang.
- **Makefile**: Mengotomatisasi tugas umum seperti build, run, dan pengambilan logs.

## Struktur Proyek

```plaintext
.
|-- .data                 # Data volume untuk postgres
|-- constants             # Folder untuk variable constants
|-- controllers           # Controllers untuk route
|-- helpers               # Utilitas umum
|-- middleware            # Middleware untuk route
|-- models                # Definisi model untuk database
|-- router                # Definisi routes API
├-- air.toml              # Konfigurasi untuk Air (live-reloading)
├-- docker-compose.yml    # Konfigurasi untuk Docker Compose
├-- Dockerfile.dev        # Dockerfile development untuk aplikasi Golang
├-- Dockerfile.prod       # Dockerfile production untuk aplikasi Golang
├── go.mod                # Dependency manager untuk Golang
├── go.sum                # Checksum dependencies
├── main.go               # Entry point aplikasi
└-- Makefile              # Makefile untuk mengotomatisasi tugas
```

## Cara Menggunakan

### Prasyarat

- **Docker** dan **Docker Compose** harus terinstal di sistem Anda.
- **Golang** minimal versi 1.17 terinstal di sistem Anda.
- **Make** terinstal di sistem Anda untuk menjalankan perintah `Makefile`.

### Langkah-langkah

1. **Clone Repository**

   ```bash
   git clone <repository-url>
   cd <project-directory>
   ```

2. **Menyiapkan Environment**

   - Buat file `.env` berdasarkan file `.env.example` dan sesuaikan konfigurasi sesuai kebutuhan Anda.

3. **Menjalankan Aplikasi**

   Anda dapat menjalankan aplikasi dengan menggunakan `Makefile`:

   - **up_build**: Menghentikan container Docker yang berjalan, membangun, dan memulai container Docker dalam mode detached.

     ```bash
     make up_build
     ```

   - **logs**: Mengambil logs dari semua layanan.

     ```bash
     make logs
     ```

   - **init**: Membangun dan memulai container Docker dalam mode detached.

     ```bash
     make init
     ```

4. **Pengembangan dengan Live Reloading**

   Aplikasi ini dikonfigurasi untuk menggunakan Air selama pengembangan untuk live-reloading. Jalankan perintah berikut:

   ```bash
   air
   ```

5. **Docker Setup**

   - **docker-compose.yml**: Mengelola dua layanan utama:
     - `db`: Layanan database PostgreSQL yang menjalankan PostgreSQL 13.
     - `web`: Layanan aplikasi Golang yang dibangun dari `Dockerfile.dev` untuk pengembangan.
   - **Dockerfile.dev**: Digunakan untuk pengembangan dengan Air.
   - **Dockerfile.prod**: Digunakan untuk membangun aplikasi untuk produksi.

## Kontribusi

Kontribusi sangat diterima! Silakan buat pull request atau laporkan masalah.
