CREATE TABLE IF NOT EXISTS infrastructure_sub_type(
  id INT PRIMARY KEY,
  type_id INT NOT NULL REFERENCES infrastructure_type(id),
  name TEXT NOT NULL,
  icon_url TEXT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO infrastructure_sub_type (id, type_id, name, icon_url, created_at)
VALUES
  (1, 1, 'Air Bersih', 'air_bersih.svg', '2023-06-29T12:54:18.610Z'),
  (2, 1, 'Air Kotor', 'air_kotor.svg', '2023-06-29T12:54:18.610Z'),
  (3, 1, 'Titik Persampahan', 'titik_persampahan.svg', '2023-06-29T12:54:18.610Z'),
  (4, 2, 'Jalan', 'jalan.svg', '2023-06-29T12:54:18.610Z'),
  (5, 2, 'Drainase', 'drainase.svg', '2023-06-29T12:54:18.610Z'),
  (6, 3, 'Lahan Parkir', 'lahan_parkir.svg', '2023-06-29T12:54:18.610Z');