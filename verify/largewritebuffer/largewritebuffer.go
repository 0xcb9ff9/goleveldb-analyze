package largewritebuffer

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/0xcb9ff9/goleveldb-analyze/v2/flags"
	"github.com/spf13/cobra"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

func GetCommand() *cobra.Command {
	kvsize := &cobra.Command{
		Use:   "verify-large-writebuffer",
		Short: "verify large write buffer defect",
		Run:   runCommand,
	}

	return kvsize
}

func runCommand(cmd *cobra.Command, _ []string) {
	leveldbPath := cmd.Flag(flags.LeveldbPathFlag).Value.String()

	// WriteBuffer > (CompactionTableSize * 8)
	// CompactionTotalSize = CompactionTableSize * 4
	// random write 32 byte key and 500 byte value
	// random read 32 byte key (exist key)
	// total write > 4GB

	// result:
	//   read slow,ldb exist size > 40mb file, all data in level 0

	BlockCacheCapacity := 64 * opt.MiB
	CompactionTableSize := 8 * opt.MiB
	CompactionTotalSize := CompactionTableSize * 10

	// WriteBuffer >= CompactionTotalSize
	// read time > 1s

	defectOption := &opt.Options{
		OpenFilesCacheCapacity:        512,
		CompactionTableSize:           CompactionTableSize,
		CompactionTotalSize:           CompactionTotalSize,
		CompactionTableSizeMultiplier: 1.1,
		BlockCacheCapacity:            BlockCacheCapacity,
		WriteBuffer:                   CompactionTotalSize * 4,
		Filter:                        filter.NewBloomFilter(2048),
		Compression:                   opt.DefaultCompression,
		NoSync:                        false,
		DisableSeeksCompaction:        true,
		BlockSize:                     256 * opt.KiB,
		FilterBaseLg:                  19, // 512 KiB
	}

	db, err := leveldb.OpenFile(leveldbPath, defectOption)

	// no defect Option
	// but read time ~ 0.01 second
	// WriteBuffer < CompactionTotalSize

	// noDefectOption := &opt.Options{
	// 	OpenFilesCacheCapacity:        512,
	// 	CompactionTableSize:           CompactionTableSize,
	// 	CompactionTotalSize:           CompactionTotalSize,
	// 	CompactionTableSizeMultiplier: 1.1,
	// 	BlockCacheCapacity:            BlockCacheCapacity,
	// 	WriteBuffer:                   CompactionTableSize * 2,
	// 	Filter:                        filter.NewBloomFilter(2048),
	// 	Compression:                   opt.DefaultCompression,
	// 	NoSync:                        false,
	// 	DisableSeeksCompaction:        true,
	// 	BlockSize:                     256 * opt.KiB,
	// 	FilterBaseLg:                  19, // 512 KiB
	// }

	// db, err := leveldb.OpenFile(leveldbPath, noDefectOption)

	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	var wg sync.WaitGroup
	wg.Add(1)

	closeCh := make(chan struct{})

	go func() {
		defer wg.Done()

		key := make([]byte, 32)
		val := make([]byte, 500)

		for {
			select {
			case <-closeCh:
				return
			default: // pass
			}

			_, _ = rand.Read(key)
			_, _ = rand.Read(val)

			err := db.Put(key, val, nil)
			if err != nil {
				log.Println(err)
			}

		}
	}()

	wg.Add(1)

	go func() {
		defer wg.Done()

		key := make([]byte, 32)

		log.Println("wait write (30 seconds)")
		time.Sleep(30 * time.Second)

		count := 0
		t := big.NewInt(0)

		for {
			select {
			case <-closeCh:
				return
			default: // pass
			}

			_, _ = rand.Read(key)

			startTime := time.Now()
			db.Get(key, nil)
			elapsed := time.Since(startTime)

			count++
			t = t.Add(t, big.NewInt(elapsed.Milliseconds()))

			if count%100 == 0 {
				t = t.Div(t, big.NewInt(100))
				fmt.Printf("[%d]:read time: %d ms\n", count, t)
				fmt.Println(db.GetProperty("leveldb.stats"))
				t = big.NewInt(0)
			}
		}
	}()

	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)

	<-cancelChan

	close(closeCh)

	wg.Wait()
}
