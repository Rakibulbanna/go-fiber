-- Create users table
CREATE TABLE IF NOT EXISTS "users" (
    "id" SERIAL PRIMARY KEY,
    "email" VARCHAR(255) UNIQUE NOT NULL,
    "password" VARCHAR(255) NOT NULL,
    "name" VARCHAR(255) NOT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMP NULL
);

CREATE INDEX IF NOT EXISTS "idx_users_deleted_at" ON "users" ("deleted_at");
CREATE UNIQUE INDEX IF NOT EXISTS "idx_users_email" ON "users" ("email");

-- Create books table
CREATE TABLE IF NOT EXISTS "books" (
    "id" SERIAL PRIMARY KEY,
    "user_id" INTEGER NOT NULL,
    "author" VARCHAR(255),
    "title" VARCHAR(255) NOT NULL,
    "publisher" VARCHAR(255) NOT NULL,
    "year" INTEGER NOT NULL,
    CONSTRAINT "fk_books_user" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS "idx_books_user_id" ON "books" ("user_id");

