-- ====================================================
-- TABLES
-- ====================================================

-- Users
CREATE TABLE IF NOT EXISTS users (
  id TEXT PRIMARY KEY NOT NULL,
  email TEXT NOT NULL UNIQUE,
  username TEXT NOT NULL,
  password TEXT NOT NULL,
  token TEXT,
  subscription_start BIGINT NOT NULL DEFAULT 0,
  subscription_end BIGINT NOT NULL DEFAULT 0
)

CREATE INDEX IF NOT EXISTS user_email_index ON users(email)

CREATE INDEX IF NOT EXISTS user_token_index ON users(token)

-- Recipes
CREATE TABLE IF NOT EXISTS recipes (
  id TEXT PRIMARY KEY NOT NULL,
  name TEXT NOT NULL,
  description TEXT NOT NULL,
  ingredients BLOB
)

-- ====================================================
-- RELATIONS
-- ====================================================

-- UserRecipe
CREATE TABLE IF NOT EXISTS user_recipe (
  user_id TEXT NOT NULL,
  recipe_id TEXT NOT NULL,
  PRIMARY KEY (user_id, recipe_id),
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (recipe_id) REFERENCES recipes(id)
)