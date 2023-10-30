package compact

import (
	"fmt"
	"log"
	"os"

	"github.com/0xcb9ff9/goleveldb-analyze/v2/flags"
	"github.com/spf13/cobra"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

func GetCommand() *cobra.Command {
	kvsize := &cobra.Command{
		Use:   "compact",
		Short: "Compact leveldb database. WARNING: May take a very long time",
		Run:   runCommand,
	}

	return kvsize
}

func runCommand(cmd *cobra.Command, _ []string) {
	leveldbPath := cmd.Flag(flags.LeveldbPathFlag).Value.String()

	if _, err := os.Stat(leveldbPath); os.IsNotExist(err) {
		log.Fatalf("leveldb path \"%s\" does not exist", leveldbPath)
	}

	db, err := leveldb.OpenFile(leveldbPath, &opt.Options{
		BlockCacheCapacity: 32 * opt.MiB,
	})

	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.CompactRange(util.Range{Start: nil, Limit: nil})
	if err != nil {
		log.Fatal(err)
	}

	stats, err := db.GetProperty("leveldb.stats")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(stats)
}
