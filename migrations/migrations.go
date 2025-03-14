// migrations/migrations.go

package migrations

import (
    "GoShort/internal/db"
    "log"

    "gorm.io/gorm"
)

type Migration struct {
    ID      string
    Migrate func(tx *gorm.DB) error
}

var migrations = []Migration{
    {
        ID: "20240101_create_urls_table",
        Migrate: func(tx *gorm.DB) error {
            return tx.Exec(`
                CREATE TABLE IF NOT EXISTS urls (
                    id SERIAL PRIMARY KEY,
                    long_url TEXT NOT NULL,
                    short_url TEXT UNIQUE NOT NULL,
                    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                    expiry TIMESTAMP,
                    clicks INT NOT NULL DEFAULT 0
                );
            `).Error
        },
    },
    {
        ID: "20240102_create_users_table",
        Migrate: func(tx *gorm.DB) error {
            return tx.Exec(`
                CREATE TABLE IF NOT EXISTS users (
                    id SERIAL PRIMARY KEY,
                    username TEXT UNIQUE NOT NULL,
                    password TEXT NOT NULL,
                    role TEXT NOT NULL,
                    last_login TIMESTAMP,
                    created_at TIMESTAMP NOT NULL DEFAULT NOW()
                );
            `).Error
        },
    },
    {
        ID: "20250314_add_userid_to_urls",
        Migrate: func(tx *gorm.DB) error {
            return tx.Exec(`
                ALTER TABLE urls
                ADD COLUMN IF NOT EXISTS user_id INT REFERENCES users(id) ON DELETE CASCADE;
            `).Error
        },
    },
}

func RunMigrations() {
    for _, migration := range migrations {
        if err := applyMigration(migration); err != nil {
            log.Fatalf("Migration failed (ID: %s): %v", migration.ID, err)
        }
        log.Printf("Migration applied successfully: %s", migration.ID)
    }
}

func applyMigration(migration Migration) error {
    log.Printf("Applying migration: %s", migration.ID)
    return db.DB.Transaction(migration.Migrate)
}
