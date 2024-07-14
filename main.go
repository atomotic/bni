package main

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log/slog"
	"os"

	_ "github.com/marcboeker/go-duckdb"
)

var migrations = `SET autoinstall_known_extensions=true;
SET autoload_known_extensions=true;
CREATE TABLE IF NOT EXISTS bni(id VARCHAR, isbn VARCHAR, title VARCHAR, "data" JSON, source VARCHAR);`

func main() {
	db, err := sql.Open("duckdb", "bni.ddb")
	if err != nil {
		slog.Error("open db", "msg", err)
		os.Exit(1)
	}
	defer db.Close()

	_, err = db.Exec(migrations)
	if err != nil {
		slog.Error("apply migrations", "msg", err)
		os.Exit(1)
	}

	xmlFile, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println(err)
	}
	defer xmlFile.Close()
	byteValue, _ := io.ReadAll(xmlFile)

	var c Collection
	xml.Unmarshal(byteValue, &c)

	for _, rec := range c.Rec {
		fmt.Printf("%s - %s\n", rec.ID(), rec.ISBN())
		j, err := json.Marshal(rec)
		if err != nil {
			fmt.Println(err)
		}
		_, err = db.Exec(`INSERT INTO bni (id, isbn, title, data, source) values (?,?,?,?,?)`,
			rec.ID(), rec.ISBN(), rec.Title(), string(j), os.Args[1])
		if err != nil {
			fmt.Println(err)
		}
	}

	// _, err = db.Exec("COPY bni TO 'bni.parquet' (FORMAT PARQUET);")
	// if err != nil {
	// 	fmt.Println(err)
	// }

}
