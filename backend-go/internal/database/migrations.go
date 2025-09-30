package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Migration struct {
	Version   int
	Name      string
	UpSQL     string
	DownSQL   string
	AppliedAt *time.Time
}

type MigrationManager struct {
	db     *sql.DB
	logger *log.Logger
}

func NewMigrationManager(db *sql.DB) *MigrationManager {
	return &MigrationManager{
		db:     db,
		logger: log.New(os.Stdout, "[MIGRATION] ", log.LstdFlags),
	}
}

func (mm *MigrationManager) CreateMigrationsTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS migrations (
			version INTEGER PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`

	_, err := mm.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	mm.logger.Println("Migrations table created/verified")
	return nil
}

func (mm *MigrationManager) GetAppliedMigrations() (map[int]bool, error) {
	applied := make(map[int]bool)

	query := "SELECT version FROM migrations ORDER BY version"
	rows, err := mm.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query applied migrations: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var version int
		if err := rows.Scan(&version); err != nil {
			return nil, fmt.Errorf("failed to scan migration version: %w", err)
		}
		applied[version] = true
	}

	return applied, nil
}

func (mm *MigrationManager) LoadMigrationsFromDir(dir string) ([]Migration, error) {
	var migrations []Migration

	files, err := filepath.Glob(filepath.Join(dir, "*.sql"))
	if err != nil {
		return nil, fmt.Errorf("failed to read migration files: %w", err)
	}

	for _, file := range files {
		filename := filepath.Base(file)
		if !strings.HasSuffix(filename, ".up.sql") {
			continue
		}

		versionStr := strings.Split(filename, "_")[0]
		version, err := strconv.Atoi(versionStr)
		if err != nil {
			mm.logger.Printf("Skipping invalid migration file: %s", filename)
			continue
		}

		name := strings.TrimSuffix(filename, ".up.sql")
		name = strings.TrimPrefix(name, versionStr+"_")

		upSQL, err := os.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("failed to read up migration file %s: %w", file, err)
		}

		downFile := strings.Replace(file, ".up.sql", ".down.sql", 1)
		var downSQL []byte
		if _, err := os.Stat(downFile); err == nil {
			downSQL, err = os.ReadFile(downFile)
			if err != nil {
				mm.logger.Printf("Warning: failed to read down migration file %s: %v", downFile, err)
			}
		}

		migrations = append(migrations, Migration{
			Version: version,
			Name:    name,
			UpSQL:   string(upSQL),
			DownSQL: string(downSQL),
		})
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

func (mm *MigrationManager) LoadBuiltinMigrations() []Migration {
	return []Migration{
		{
			Version: 1,
			Name:    "initial_schema",
			UpSQL: `
				CREATE TABLE IF NOT EXISTS users (
					id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
					email VARCHAR(255) UNIQUE NOT NULL,
					name VARCHAR(255),
					password VARCHAR(255) NOT NULL,
					role VARCHAR(50) DEFAULT 'user',
					image VARCHAR(500),
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
				);

				CREATE TABLE IF NOT EXISTS categories (
					id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
					name VARCHAR(255) NOT NULL,
					slug VARCHAR(255) UNIQUE NOT NULL,
					description TEXT,
					image VARCHAR(500),
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
				);

				CREATE TABLE IF NOT EXISTS products (
					id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
					name VARCHAR(255) NOT NULL,
					slug VARCHAR(255) UNIQUE NOT NULL,
					description TEXT,
					price DECIMAL(10,2) NOT NULL,
					compare_price DECIMAL(10,2),
					images TEXT[],
					in_stock BOOLEAN DEFAULT true,
					stock INTEGER DEFAULT 0,
					featured BOOLEAN DEFAULT false,
					category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
				);

				CREATE TABLE IF NOT EXISTS orders (
					id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
					user_id UUID REFERENCES users(id) ON DELETE CASCADE,
					total DECIMAL(10,2) NOT NULL,
					status VARCHAR(50) DEFAULT 'pending',
					shipping_address JSONB,
					billing_address JSONB,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
				);

				CREATE TABLE IF NOT EXISTS order_items (
					id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
					order_id UUID REFERENCES orders(id) ON DELETE CASCADE,
					product_id UUID REFERENCES products(id) ON DELETE CASCADE,
					quantity INTEGER NOT NULL,
					price DECIMAL(10,2) NOT NULL,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
				);

				CREATE TABLE IF NOT EXISTS cart_items (
					id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
					user_id UUID REFERENCES users(id) ON DELETE CASCADE,
					product_id UUID REFERENCES products(id) ON DELETE CASCADE,
					quantity INTEGER NOT NULL,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					UNIQUE(user_id, product_id)
				);

				CREATE TABLE IF NOT EXISTS wishlist_items (
					id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
					user_id UUID REFERENCES users(id) ON DELETE CASCADE,
					product_id UUID REFERENCES products(id) ON DELETE CASCADE,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					UNIQUE(user_id, product_id)
				);

				CREATE TABLE IF NOT EXISTS reviews (
					id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
					user_id UUID REFERENCES users(id) ON DELETE CASCADE,
					product_id UUID REFERENCES products(id) ON DELETE CASCADE,
					rating INTEGER CHECK (rating >= 1 AND rating <= 5),
					comment TEXT,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
				);

				CREATE TABLE IF NOT EXISTS payments (
					id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
					order_id UUID REFERENCES orders(id) ON DELETE CASCADE,
					amount DECIMAL(10,2) NOT NULL,
					currency VARCHAR(3) DEFAULT 'USD',
					status VARCHAR(50) DEFAULT 'pending',
					payment_method VARCHAR(50),
					transaction_id VARCHAR(255),
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
				);

				CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
				CREATE INDEX IF NOT EXISTS idx_products_category_id ON products(category_id);
				CREATE INDEX IF NOT EXISTS idx_products_featured ON products(featured);
				CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);
				CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);
				CREATE INDEX IF NOT EXISTS idx_cart_items_user_id ON cart_items(user_id);
				CREATE INDEX IF NOT EXISTS idx_wishlist_items_user_id ON wishlist_items(user_id);
				CREATE INDEX IF NOT EXISTS idx_reviews_product_id ON reviews(product_id);
				CREATE INDEX IF NOT EXISTS idx_payments_order_id ON payments(order_id);
			`,
			DownSQL: `
				DROP TABLE IF EXISTS payments;
				DROP TABLE IF EXISTS reviews;
				DROP TABLE IF EXISTS wishlist_items;
				DROP TABLE IF EXISTS cart_items;
				DROP TABLE IF EXISTS order_items;
				DROP TABLE IF EXISTS orders;
				DROP TABLE IF EXISTS products;
				DROP TABLE IF EXISTS categories;
				DROP TABLE IF EXISTS users;
			`,
		},
		{
			Version: 2,
			Name:    "add_indexes",
			UpSQL: `
				CREATE INDEX IF NOT EXISTS idx_products_name ON products(name);
				CREATE INDEX IF NOT EXISTS idx_products_price ON products(price);
				CREATE INDEX IF NOT EXISTS idx_products_created_at ON products(created_at);
				CREATE INDEX IF NOT EXISTS idx_orders_created_at ON orders(created_at);
				CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at);
			`,
			DownSQL: `
				DROP INDEX IF EXISTS idx_products_name;
				DROP INDEX IF EXISTS idx_products_price;
				DROP INDEX IF EXISTS idx_products_created_at;
				DROP INDEX IF EXISTS idx_orders_created_at;
				DROP INDEX IF EXISTS idx_users_created_at;
			`,
		},
		{
			Version: 3,
			Name:    "add_soft_delete",
			UpSQL: `
				ALTER TABLE users ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;
				ALTER TABLE products ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;
				ALTER TABLE categories ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;
				ALTER TABLE orders ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;
				
				CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);
				CREATE INDEX IF NOT EXISTS idx_products_deleted_at ON products(deleted_at);
				CREATE INDEX IF NOT EXISTS idx_categories_deleted_at ON categories(deleted_at);
				CREATE INDEX IF NOT EXISTS idx_orders_deleted_at ON orders(deleted_at);
			`,
			DownSQL: `
				DROP INDEX IF EXISTS idx_users_deleted_at;
				DROP INDEX IF EXISTS idx_products_deleted_at;
				DROP INDEX IF EXISTS idx_categories_deleted_at;
				DROP INDEX IF EXISTS idx_orders_deleted_at;
				
				ALTER TABLE users DROP COLUMN IF EXISTS deleted_at;
				ALTER TABLE products DROP COLUMN IF EXISTS deleted_at;
				ALTER TABLE categories DROP COLUMN IF EXISTS deleted_at;
				ALTER TABLE orders DROP COLUMN IF EXISTS deleted_at;
			`,
		},
	}
}

func (mm *MigrationManager) Up(ctx context.Context) error {
	if err := mm.CreateMigrationsTable(); err != nil {
		return err
	}

	applied, err := mm.GetAppliedMigrations()
	if err != nil {
		return err
	}

	migrations := mm.LoadBuiltinMigrations()

	for _, migration := range migrations {
		if applied[migration.Version] {
			mm.logger.Printf("Migration %d_%s already applied, skipping", migration.Version, migration.Name)
			continue
		}

		mm.logger.Printf("Applying migration %d_%s", migration.Version, migration.Name)

		tx, err := mm.db.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("failed to begin transaction for migration %d: %w", migration.Version, err)
		}

		if _, err := tx.ExecContext(ctx, migration.UpSQL); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute migration %d_%s: %w", migration.Version, migration.Name, err)
		}

		_, err = tx.ExecContext(ctx, "INSERT INTO migrations (version, name) VALUES ($1, $2)", migration.Version, migration.Name)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %d_%s: %w", migration.Version, migration.Name, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration %d_%s: %w", migration.Version, migration.Name, err)
		}

		mm.logger.Printf("Successfully applied migration %d_%s", migration.Version, migration.Name)
	}

	mm.logger.Println("All migrations applied successfully")
	return nil
}

func (mm *MigrationManager) Down(ctx context.Context, targetVersion int) error {
	if err := mm.CreateMigrationsTable(); err != nil {
		return err
	}

	applied, err := mm.GetAppliedMigrations()
	if err != nil {
		return err
	}

	migrations := mm.LoadBuiltinMigrations()

	for i := len(migrations) - 1; i >= 0; i-- {
		migration := migrations[i]

		if migration.Version <= targetVersion {
			break
		}

		if !applied[migration.Version] {
			continue
		}

		if migration.DownSQL == "" {
			mm.logger.Printf("No down migration for %d_%s, skipping", migration.Version, migration.Name)
			continue
		}

		mm.logger.Printf("Rolling back migration %d_%s", migration.Version, migration.Name)

		tx, err := mm.db.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("failed to begin transaction for rollback %d: %w", migration.Version, err)
		}

		if _, err := tx.ExecContext(ctx, migration.DownSQL); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute rollback %d_%s: %w", migration.Version, migration.Name, err)
		}

		_, err = tx.ExecContext(ctx, "DELETE FROM migrations WHERE version = $1", migration.Version)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to remove migration record %d_%s: %w", migration.Version, migration.Name, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit rollback %d_%s: %w", migration.Version, migration.Name, err)
		}

		mm.logger.Printf("Successfully rolled back migration %d_%s", migration.Version, migration.Name)
	}

	mm.logger.Printf("Rolled back to version %d", targetVersion)
	return nil
}

func (mm *MigrationManager) Status() ([]Migration, error) {
	if err := mm.CreateMigrationsTable(); err != nil {
		return nil, err
	}

	applied, err := mm.GetAppliedMigrations()
	if err != nil {
		return nil, err
	}

	migrations := mm.LoadBuiltinMigrations()

	for i := range migrations {
		if applied[migrations[i].Version] {
			now := time.Now()
			migrations[i].AppliedAt = &now
		}
	}

	return migrations, nil
}

func (mm *MigrationManager) CreateMigration(name string) error {
	migrations := mm.LoadBuiltinMigrations()
	nextVersion := 1

	if len(migrations) > 0 {
		nextVersion = migrations[len(migrations)-1].Version + 1
	}

	upFile := fmt.Sprintf("%03d_%s.up.sql", nextVersion, name)
	downFile := fmt.Sprintf("%03d_%s.down.sql", nextVersion, name)

	upContent := fmt.Sprintf("-- Migration: %s\n-- Version: %d\n-- Up\n\n", name, nextVersion)
	downContent := fmt.Sprintf("-- Migration: %s\n-- Version: %d\n-- Down\n\n", name, nextVersion)

	if err := os.WriteFile(upFile, []byte(upContent), 0644); err != nil {
		return fmt.Errorf("failed to create up migration file: %w", err)
	}

	if err := os.WriteFile(downFile, []byte(downContent), 0644); err != nil {
		return fmt.Errorf("failed to create down migration file: %w", err)
	}

	mm.logger.Printf("Created migration files: %s, %s", upFile, downFile)
	return nil
}
