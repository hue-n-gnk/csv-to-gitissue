# migrate

## Services and libraries

- [migrate](https://github.com/golang-migrate/migrate) contains migration for Postgresql

### Install

    make install

### Migrate

    make migrate

### Drop (only use for development env)

    make drop

### How to add new migration

1. add sql file to path [ddl](packages/migrate/ddl) with name format ${increase number}_${name}.up.sql
   - The increase number must be uniqued.
2. run `make migrate` to apply to DB.
