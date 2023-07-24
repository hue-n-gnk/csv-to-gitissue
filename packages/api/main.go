package main

import (
	"fmt"
	"hue-n-gnk/csv-to-gitisue/database/factory"
	"hue-n-gnk/csv-to-gitisue/helpers/env"
	"hue-n-gnk/csv-to-gitisue/service/backlogcvs"

	"github.com/pkg/errors"
)

func setup() error {
	// database setup
	e := env.Load()
	factory.NewConnection(factory.Config{
		User:     e.PostgrestUser,
		Password: e.PostgrestPassword,
		Host:     e.PostgrestHostname,
		Database: e.PostgrestDatabase,
		Port:     e.PostgrestPort,
	})

	return nil
}

func run() error {
	if err := setup(); err != nil {
		return errors.Wrap(err, "failed application setup")
	}
	// service start
	s := backlogcvs.NewService()
	return s.ToGithub()
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
	fmt.Println("DONE")
	return
}
