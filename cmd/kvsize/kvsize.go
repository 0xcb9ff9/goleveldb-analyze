package kvsize

import (
	"log"
	"os"

	"github.com/0xcb9ff9/goleveldb-analyze/v2/flags"
	"github.com/spf13/cobra"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var (
	params = &kvsizeParams{}
)

type kvsizeParams struct {
	out string
}

func GetCommand() *cobra.Command {
	kvsize := &cobra.Command{
		Use:   "kvsize",
		Short: "Get the size of the key-value pairs in the database",
		Run:   runCommand,
	}

	setFlags(kvsize)

	return kvsize
}

func setFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(
		&params.out,
		"out",
		"./kvsize.sqlite3",
		"out sqlite database path",
	)
}

func runCommand(cmd *cobra.Command, _ []string) {
	leveldbPath := cmd.Flag(flags.LeveldbPathFlag).Value.String()

	if _, err := os.Stat(leveldbPath); os.IsNotExist(err) {
		log.Fatalf("leveldb path \"%s\" does not exist", leveldbPath)
	}

	db, err := leveldb.OpenFile(leveldbPath, &opt.Options{
		ReadOnly:           true,
		BlockCacheCapacity: 32 * opt.MiB,
	})

	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqliteDB, err := sql.Open("sqlite3", params.out)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqliteDB.Exec("CREATE TABLE IF NOT EXISTS kvsize (id INTEGER NOT NULL PRIMARY KEY, key_size INTEGER, value_size INTEGER, sum_size INTEGER)")

	iter := db.NewIterator(nil, nil)
	defer iter.Release()

	tx, err := sqliteDB.Begin()
	if err != nil {
		log.Fatal(err)
	}

	count := 0

	for iter.Next() {
		key := iter.Key()
		value := iter.Value()

		if count > 1024 {
			err := tx.Commit()
			if err != nil {
				log.Fatal(err)
			}

			tx, err = sqliteDB.Begin()
			if err != nil {
				log.Fatal(err)
			}

			count = 0
		}

		_, err = tx.Exec("INSERT INTO kvsize (key_size, value_size, sum_size) VALUES (?, ?, ?)", len(key), len(value), len(key)+len(value))
		if err != nil {
			log.Fatal(err)
		}

		count++
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

}
