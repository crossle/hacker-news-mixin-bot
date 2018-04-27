CREATE TABLE IF NOT EXISTS subscribers (
  user_id          VARCHAR(36) PRIMARY KEY,
  created_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX subscribers_user_id ON subscribers (user_id);
