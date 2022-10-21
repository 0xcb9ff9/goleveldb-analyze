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