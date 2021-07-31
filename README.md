# liquigo
Liquibase inspired database migration tool written in Go.

Supports SQL updates to PortgreSQL database.

Rollback is not yet implemented.

## Usage

DB change log is a `yaml` file, listing all SQL files. Files processed in the order as listed.

Each SQL file contains changesets. A changeset starts with a `-- changeset id-of-a-changeset` and may contain one or multiple SQL statements. Each SQL must end with a `;`. Each changeset applied in one separate transaction.

Two tables are created: `dbchangelog` - to keep the history of applied changesets and `dbchangelock` - to allow only one liquigo process at a time.

To modify an SQL that has been applied and avoid the error of mismatched MD5, add a comment to the changeset: `-- md5 new-md5-value`.

By default SQL statements split with `;`. If you have a statement with `;`, stored procedure for example, that you don't want to split, add a comment to changeset: `-- splitStatements false`.
