package duplicate

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

var (
	params = &duplicateParams{}
)

type duplicateParams struct {
	src string
	out string

	// option
	tableSize      int
	tableTotalSize int
	bloomKeyBit    int
	blockSize      int
	filterBaseLg   int
	blockCache     int
	fileFds        int
}

func GetCommand() *cobra.Command {
	duplicate := &cobra.Command{
		Use:   "duplicate",
		Short: "duplicate goleveldb",
		Run:   runCommand,
	}

	setFlags(duplicate)

	return duplicate
}
func setFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(
		&params.src,
		"src",
		"",
		"duplicate database source",
	)

	cmd.Flags().StringVar(
		&params.out,
		"out",
		"",
		"duplicate database output",
	)

	cmd.Flags().IntVar(
		&params.tableSize,
		"table-size",
		4,
		"leveldb table size (unit mib)",
	)

	cmd.Flags().IntVar(
		&params.tableTotalSize,
		"table-total-size",
		100,
		"leveldb table total size (unit mib)",
	)

	cmd.Flags().IntVar(
		&params.bloomKeyBit,
		"bloom-bits",
		2048,
		"leveldb bloom key bit",
	)

	cmd.Flags().IntVar(
		&params.blockSize,
		"block-size",
		512,
		"leveldb bloom key bit",
	)

	cmd.Flags().IntVar(
		&params.filterBaseLg,
		"filter-base",
		16,
		"set FilterBaseLg, is the log size for filter block to create a bloom filter.",
	)

	cmd.Flags().IntVar(
		&params.blockCache,
		"cache",
		1024,
		"sorted table block caching",
	)

	cmd.Flags().IntVar(
		&params.fileFds,
		"handles",
		2048,
		"open files max handle",
	)
}

func runCommand(cmd *cobra.Command, _ []string) {

	if _, err := os.Stat(params.src); os.IsNotExist(err) {
		log.Fatalf("leveldb source path \"%s\" does not exist", params.src)
	}

	srcdb, err := leveldb.OpenFile(params.src, &opt.Options{
		ReadOnly:               true,
		NoSync:                 true,
		BlockCacheCapacity:     params.blockCache * opt.MiB,
		OpenFilesCacheCapacity: params.fileFds,
		DisableSeeksCompaction: true,
	})

	if err != nil {
		panic(err)
	}
	defer srcdb.Close()

	outdb, err := leveldb.OpenFile(params.out, &opt.Options{
		BlockCacheCapacity:            params.blockCache * opt.MiB,
		OpenFilesCacheCapacity:        params.fileFds,
		CompactionTableSize:           params.tableSize * opt.MiB,
		CompactionTotalSize:           params.tableTotalSize * opt.MiB,
		WriteBuffer:                   (params.tableSize * 2) * opt.MiB,
		CompactionTableSizeMultiplier: 1.1,
		Filter:                        filter.NewBloomFilter(params.bloomKeyBit),
		BlockSize:                     params.blockSize * opt.KiB,
		FilterBaseLg:                  params.filterBaseLg,
		DisableSeeksCompaction:        true,
		NoSync:                        true,
	})

	if err != nil {
		panic(err)
	}
	defer srcdb.Close()

	iter := srcdb.NewIterator(nil, &opt.ReadOptions{
		DontFillCache: true,
		Strict:        opt.StrictJournalChecksum | opt.StrictBlockChecksum | opt.StrictReader,
	})
	defer iter.Release()

	for iter.Next() {
		key := iter.Key()
		value := iter.Value()

		outdb.Put(key, value, nil)

		log.Printf("write key: %s", EncodeToHex(key))
	}

	stats, err := outdb.GetProperty("leveldb.stats")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(stats)
}

func EncodeToHex(str []byte) string {
	return hex.EncodeToString(str)
}
