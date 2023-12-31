CREATE TABLE IF NOT EXISTS infrastructure(
  id TEXT PRIMARY KEY,
  sub_type_id INT NOT NULL REFERENCES infrastructure_sub_type(id),
  name TEXT NOT NULL,
  details TEXT not null,
  status TEXT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);