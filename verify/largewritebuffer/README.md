## goleveldb large WriteBuffer defect

Test environment:
  * OS: Fedora Linux 36.20221008.0 (Silverblue)
  * CPU: Intel Core i7 10700K
  * Memory: 8GB RAM
  * Disk: INTEL SS DPEKNW512G8 (1.00)
  * Go Version: 1.19.2 linux/amd64
  * Partition format: brtfs

Condition:

 * WriteBuffer > CompactionTotalSize * 4
 * keep writing
 * read key (not exist in db)

Result:

1. ldb file size > CompactionTableSize
2. read key slow (> 100ms)
3. level 0 is large (> 1 GB)

```
[100]:read time: 99 ms
Compactions
 Level |   Tables   |    Size(MB)   |    Time(sec)  |    Read(MB)   |   Write(MB)
-------+------------+---------------+---------------+---------------+---------------
   0   |          2 |     946.45611 |      26.40013 |       0.00000 |     946.45611
-------+------------+---------------+---------------+---------------+---------------
 Total |          2 |     946.45611 |      26.40013 |       0.00000 |     946.45611
 <nil>
[200]:read time: 119 ms
Compactions
 Level |   Tables   |    Size(MB)   |    Time(sec)  |    Read(MB)   |   Write(MB)
-------+------------+---------------+---------------+---------------+---------------
   0   |          3 |    1419.68426 |      39.63862 |       0.00000 |    1419.68426
-------+------------+---------------+---------------+---------------+---------------
 Total |          3 |    1419.68426 |      39.63862 |       0.00000 |    1419.68426
 <nil>
[300]:read time: 217 ms
Compactions
 Level |   Tables   |    Size(MB)   |    Time(sec)  |    Read(MB)   |   Write(MB)
-------+------------+---------------+---------------+---------------+---------------
   0   |          5 |    2366.14155 |      57.16061 |       0.00000 |    2366.14155
-------+------------+---------------+---------------+---------------+---------------
 Total |          5 |    2366.14155 |      57.16061 |       0.00000 |    2366.14155
 <nil>
[400]:read time: 259 ms
Compactions
 Level |   Tables   |    Size(MB)   |    Time(sec)  |    Read(MB)   |   Write(MB)
-------+------------+---------------+---------------+---------------+---------------
   0   |          6 |    2839.37023 |      83.36314 |       0.00000 |    2839.37023
-------+------------+---------------+---------------+---------------+---------------
 Total |          6 |    2839.37023 |      83.36314 |       0.00000 |    2839.37023
 <nil>
```

No defect Option:

[largewritebuffer.go#L71](./largewritebuffer.go#L71)

```go
BlockCacheCapacity := 64 * opt.MiB
CompactionTableSize := 8 * opt.MiB
CompactionTotalSize := CompactionTableSize * 10

option := &opt.Options{
	OpenFilesCacheCapacity:        512,
	CompactionTableSize:           CompactionTableSize,
	CompactionTotalSize:           CompactionTotalSize,
	CompactionTableSizeMultiplier: 1.1,
	BlockCacheCapacity:            BlockCacheCapacity,
	WriteBuffer:                   CompactionTableSize * 2, // limit WriteBuffer < CompactionTotalSize
	Filter:                        filter.NewBloomFilter(2048),
	Compression:                   opt.DefaultCompression,
	NoSync:                        false,
	DisableSeeksCompaction:        true,
	BlockSize:                     256 * opt.KiB,
	FilterBaseLg:                  19, // 512 KiB
}
```

Result:

```
[6800]:read time: 0 ms
Compactions
 Level |   Tables   |    Size(MB)   |    Time(sec)  |    Read(MB)   |   Write(MB)
-------+------------+---------------+---------------+---------------+---------------
   0   |          5 |     118.40486 |      13.18796 |       0.00000 |    1349.82141
   1   |         76 |     991.41655 |     113.58148 |    7307.58506 |    7306.42303
   2   |         18 |     238.83577 |       0.32603 |      53.07663 |      53.07443
-------+------------+---------------+---------------+---------------+---------------
 Total |         99 |    1348.65718 |     127.09548 |    7360.66169 |    8709.31887
 <nil>
[6900]:read time: 1 ms
Compactions
 Level |   Tables   |    Size(MB)   |    Time(sec)  |    Read(MB)   |   Write(MB)
-------+------------+---------------+---------------+---------------+---------------
   0   |          5 |     118.40486 |      13.18796 |       0.00000 |    1349.82141
   1   |         76 |     991.41655 |     113.58148 |    7307.58506 |    7306.42303
   2   |         18 |     238.83577 |       0.32603 |      53.07663 |      53.07443
-------+------------+---------------+---------------+---------------+---------------
 Total |         99 |    1348.65718 |     127.09548 |    7360.66169 |    8709.31887
 <nil>
[7000]:read time: 2 ms
Compactions
 Level |   Tables   |    Size(MB)   |    Time(sec)  |    Read(MB)   |   Write(MB)
-------+------------+---------------+---------------+---------------+---------------
   0   |          6 |     142.08577 |      13.63124 |       0.00000 |    1373.50232
   1   |         76 |     991.41655 |     113.58148 |    7307.58506 |    7306.42303
   2   |         18 |     238.83577 |       0.32603 |      53.07663 |      53.07443
-------+------------+---------------+---------------+---------------+---------------
 Total |        100 |    1372.33809 |     127.53876 |    7360.66169 |    8732.99978
 <nil>
[7100]:read time: 3 ms
Compactions
 Level |   Tables   |    Size(MB)   |    Time(sec)  |    Read(MB)   |   Write(MB)
-------+------------+---------------+---------------+---------------+---------------
   0   |          6 |     142.08577 |      13.63124 |       0.00000 |    1373.50232
   1   |         76 |     991.41655 |     113.58148 |    7307.58506 |    7306.42303
   2   |         18 |     238.83577 |       0.32603 |      53.07663 |      53.07443
-------+------------+---------------+---------------+---------------+---------------
 Total |        100 |    1372.33809 |     127.53876 |    7360.66169 |    8732.99978
 <nil>
```