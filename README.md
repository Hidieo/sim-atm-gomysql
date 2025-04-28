# Simulasi Mesin ATM (CLI Go + MySQL)

Ini adalah proyek simulasi mesin ATM sederhana berbasis Command Line Interface (CLI) yang dibangun menggunakan bahasa Go (Golang) dan MySQL. Proyek ini dikembangkan sebagai bagian dari tugas Digital Skill Fair 38.0 di program pelatihan ibimbing.

## 📚 Fitur Aplikasi

- **Register**: Membuat akun baru dengan nama dan PIN.
- **Login**: Masuk ke akun menggunakan ID akun dan PIN.
- **Check Balance**: Melihat saldo rekening.
- **Deposit**: Menambahkan saldo ke akun.
- **Withdraw**: Menarik saldo dari akun (dengan pengecekan saldo cukup).
- **Transfer**: Mengirim saldo ke akun lain.
- **Transaction History**: Melihat riwayat transfer masuk dan keluar.

## 🛠️ Teknologi yang Digunakan

- **Bahasa Pemrograman**: Go (Golang)
- **Database**: MySQL
- **Library Tambahan**: [urfave/cli](https://github.com/urfave/cli) untuk pembuatan CLI interaktif

## 🧩 Struktur Database

Terdapat dua tabel utama:

- **accounts**: Menyimpan data akun pengguna (id, nama, pin, saldo, tanggal pembuatan).
- **transactions**: Menyimpan seluruh histori transaksi (deposit, withdraw, transfer).

## 🗂️ Instalasi & Setup

Ikuti langkah-langkah berikut untuk mengatur proyek ini di mesin lokal Anda:

1. Clone repository ini:

    ```bash
    git clone https://github.com/username/repo-ATM-simulation.git
    cd repo-ATM-simulation
    ```

2. Impor database MySQL menggunakan file `atm_simulation.sql`.

3. Atur koneksi database di file `main.go`. Sesuaikan kredensial dan pengaturan koneksi ke MySQL.

4. Jalankan aplikasi:

    ```bash
    go run main.go
    ```

## 🔥 Fitur Unggulan

- Sistem transaksi atomik menggunakan transaction di database (rollback saat error).
- Pencatatan histori transaksi untuk transfer masuk dan keluar.
- Validasi saldo sebelum withdraw atau transfer.
- Struktur database relasional untuk menjaga konsistensi data.

## 📄 License

Proyek ini dibuat untuk tujuan pembelajaran. Feel free to fork dan kembangkan lebih lanjut!
