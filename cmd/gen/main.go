package main

import (
	"gamesync/internal/dbx"
	"log"
	"os"

	"gorm.io/gen"
	"gorm.io/gorm/logger"
)


func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath: "internal/query",
		Mode: gen.WithDefaultQuery|gen.WithQueryInterface,
		FieldNullable: true,
	})

	db, err := dbx.ConnectDb(os.Getenv("GAMESYNC_DB_TYPE"), os.Getenv("GAMESYNC_DB_URL"), logger.Error)
	if err != nil {
		log.Fatal(err)
	}
	g.UseDB(db)
	g.ApplyBasic(g.GenerateAllTable()...)
	g.Execute()
}

