# 🧩 SudokuCLI

**go-Sudoku** adalah permainan **Sudoku interaktif berbasis terminal (CLI)** dengan fitur visualisasi papan sebagai **gambar**. Cocok untuk belajar logika dan mengasah otak, langsung dari terminal!

> Dibuat oleh: © 2025 HyHy  
> Bahasa: Go (Golang)  
> Output: PNG (menggunakan [fogleman/gg](https://github.com/fogleman/gg))

---

## 📦 Fitur Utama

- Bermain Sudoku langsung dari terminal
- Dukungan tingkat kesulitan: `easy`, `normal`, `hard`, `extreme`
- Output visual sebagai `sudoku.png`
- Peringatan kesalahan dengan **penanda merah**
- Bantuan (hint) terbatas 3x
- Timeout otomatis (5 menit idle)

---

## 🛠️ Cara Install

### 1. Clone repository

```bash
git clone https://github.com/Terror-Machine/go-sudoku
cd go-sudoku
````

### 2. Install dependensi

Pastikan Go sudah terpasang (Go 1.18 atau lebih baru):

```bash
go mod tidy
```

Jika belum, pasang Go terlebih dahulu: [https://go.dev/dl/](https://go.dev/dl/)

---

## ▶️ Cara Menjalankan

```bash
go run main.go
```

Setelah dijalankan, Anda akan melihat prompt:

```text
Selamat datang di Game Sudoku CLI!
Papan Sudoku (extreme) telah dibuat dan disimpan sebagai 'sudoku.png'.
```

Setiap perubahan akan memperbarui file gambar `sudoku.png`.

---

## 🎮 Cara Bermain

### 📌 Format perintah

| Perintah          | Deskripsi                                      |
| ----------------- | ---------------------------------------------- |
| `a1 5`            | Mengisi kotak kolom A baris 1 dengan angka 5   |
| `a1 0`            | Menghapus angka di kotak A1                    |
| `hint`            | Mendapatkan bantuan otomatis (maks 3 kali)     |
| `cek`             | Mengecek apakah jawaban sudah benar atau belum |
| `new easy`        | Memulai game baru dengan tingkat easy          |
| `new extreme`     | Game baru tingkat paling sulit                 |
| `exit` / `keluar` | Keluar dari permainan                          |

### 📸 Output visual

Setiap aksi akan memperbarui file `sudoku.png`, contoh tampilannya:

* 🔢 Angka asli dari puzzle: **hitam**
* 🔷 Angka yang Anda isi: **biru**
* ❌ Angka salah: **merah**

---

## 💡 Tips

* Jika ingin mengganti font atau ukuran gambar, modifikasi fungsi `generateSudokuImage()` di `main.go`
* Gambar akan tertimpa setiap kali Anda melakukan aksi
* Cocok untuk disandingkan dengan `watch` di Linux/macOS:

  ```bash
  watch -n 1 feh sudoku.png
  ```

---

## 📄 Lisensi

MIT © 2025 HyHy

---

## 🤝 Kontribusi

Pull request & issue dipersilahkan. 
Mie Sepesial Pake Telor boleh juga. 