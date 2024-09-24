-- Hapus kolom description
ALTER TABLE notifications
DROP COLUMN description;

-- Tambahkan kembali kolom description dengan tipe data jsonb
ALTER TABLE notifications
ADD COLUMN description jsonb;