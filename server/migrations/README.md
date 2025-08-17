# Database Migrations

This directory contains your database migration files.

## Naming Convention

Use timestamps for ordering: `YYYYMMDDHHMMSS_description.sql`

Examples:

- `20240817120000_create_users_table.sql`
- `20240817120100_create_quizzes_table.sql`
- `20240817120200_add_index_to_users_email.sql`

## Migration Tools

Consider using:

- golang-migrate/migrate
- pressly/goose
- Custom migration system
