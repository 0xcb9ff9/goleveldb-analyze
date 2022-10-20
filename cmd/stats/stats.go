package stats

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/0xcb9ff9/goleveldb-analyze/v2/flags"
	"github.com/spf13/cobra"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

func GetCommand() *cobra.Command {
	kvsize := &cobra.Command{
		Use:   "stats",
		Short: "print stats about the leveldb",
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
		ReadOnly:           true,
		BlockCacheCapacity: 32 * opt.MiB,
	})

	if err != nil {
		panic(err)
	}
	defer db.Close()

	iter := db.NewIterator(nil, nil)
	defer iter.Release()

	closeCh := make(chan struct{})

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		for iter.Next() {
			select {
			case <-closeCh:
				return
			default: // pass
			}

			_ = iter.Key()
			_ = iter.Value()
		}
	}()

	fmt.Println("Counting... wait 1 minute")

	time.Sleep(1 * time.Minute) // read 1 minute
	close(closeCh)

	wg.Wait()

	stats, err := db.GetProperty("leveldb.stats")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(stats)
}
