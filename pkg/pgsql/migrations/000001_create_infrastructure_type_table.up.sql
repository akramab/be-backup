CREATE TABLE IF NOT EXISTS infrastructure_type(
  id INT PRIMARY KEY,
  name TEXT NOT NULL,
  icon_url TEXT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO infrastructure_type (id, name, icon_url, created_at)
VALUES
  (1, 'Titik', 'titik.svg', '2023-06-29T12:54:18.610Z'),
  (2, 'Garis', 'garis.svg', '2023-06-29T12:54:18.610Z'),
  (3, 'Bidang', 'bidang.svg', '2023-06-29T12:54:18.610Z');