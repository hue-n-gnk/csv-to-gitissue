// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package entities

const TableNameSchemaMigration = "schema_migrations"

// SchemaMigration mapped from table <schema_migrations>
type SchemaMigration struct {
	Version int64 `gorm:"column:version;type:bigint;primaryKey" json:"version"`
	Dirty   bool  `gorm:"column:dirty;type:boolean;not null" json:"dirty"`
}

// TableName SchemaMigration's table name
func (*SchemaMigration) TableName() string {
	return TableNameSchemaMigration
}