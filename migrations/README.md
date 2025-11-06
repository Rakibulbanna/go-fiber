# Database Migrations

This project uses **Atlas** to automatically generate migrations from GORM models.

## How It Works

1. **Define your models** in the `models/` directory (e.g., `user.go`, `books.go`)
2. **Generate migrations** using `make migrate-diff` - this compares your GORM models with the current database schema
3. **Apply migrations** using `make migrate-apply`

## Commands

### Generate Migrations from GORM Models
```bash
make migrate-diff
```
This command will:
- Read your GORM models from `models/` directory
- Compare them with the current database schema
- Generate SQL migration files in `migrations/` directory

### Apply Migrations
```bash
make migrate-apply
```
Applies all pending migrations to your database.

### Rollback Last Migration
```bash
make migrate-down
```
Rolls back the last applied migration.

### Check Migration Status
```bash
make migrate-status
```
Shows which migrations have been applied and which are pending.

## Workflow

1. **Modify your GORM models** (add fields, change types, etc.)
2. **Run `make migrate-diff`** to generate migration files
3. **Review the generated migration files** in `migrations/` directory
4. **Run `make migrate-apply`** to apply the changes

## Important Notes

- **Never write raw SQL manually** - all migrations are auto-generated from GORM models
- **Always review generated migrations** before applying them
- **Add new models** to `cmd/atlas/main.go` when creating new model files
- The schema loader program (`cmd/atlas/main.go`) extracts schema from GORM models

## Adding New Models

When you create a new model file, add it to the `modelsList` in `cmd/atlas/main.go`:

```go
modelsList := []interface{}{
    &models.User{},
    &models.Book{},
    &models.YourNewModel{}, // Add your new model here
}
```

