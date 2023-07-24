package main

import (
	"fmt"
	"hue-n-gnk/csv-to-gitisue/helpers/env"

	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func main() {

	cf := gen.Config{
		OutPath:      "../../database/dao",
		ModelPkgPath: "entities",
		// Ignore generate unit test due to
		// https://github.com/go-gorm/gen/issues/638
		// WithUnitTest:   true,
		FieldNullable:  true,
		FieldCoverable: true,
		FieldSignable:  true,
		// FieldWithIndexTag: true,
		FieldWithTypeTag: true,
		Mode:             gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode

	}
	g := gen.NewGenerator(cf)

	db, err := connectDatabase()
	if err != nil {
		panic(errors.Wrap(err, "Connect database failed"))
	}
	g.UseDB(db)

	g.ApplyBasic(g.GenerateAllTable()...)
	// Generate the code
	g.Execute()
}

func connectDatabase() (*gorm.DB, error) {
	e := env.Load()

	var dsn string
	gormConfig := &gorm.Config{}
	dsn = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", e.PostgrestHostname, e.PostgrestPort, e.PostgrestUser, e.PostgrestDatabase, e.PostgrestPassword)
	return gorm.Open(postgres.Open(dsn), gormConfig)
}
