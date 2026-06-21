-- Drop database dan buat ulang (untuk development)
-- PERINGATAN: Ini akan menghapus semua data!

DROP DATABASE IF EXISTS bruvela_db;
CREATE DATABASE bruvela_db;

-- Setelah ini, jalankan:
-- psql -U yogatama.egiantoro -d bruvela_db -f migrations/schema.sql
-- psql -U yogatama.egiantoro -d bruvela_db -f migrations/seed.sql
