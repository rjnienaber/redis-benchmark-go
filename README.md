# redis-benchmark-go
This is an implementation of the [redis-benchmark] in Go. This was done as a proof-of-concept to see if Go could achieve 
similar performance to the C implementation. As a proof-of-concept, only a subset of the reference implementation's 
functionality was programmed (i.e. just the `PING` command) though it should be relatively easy to extend.

As a caveat to inferring too much from these results, please read [Pitfalls and misconceptions] about redis benchmarking.

## Benchmark
### Results

| Implementation | Requests/sec | Elapsed (s) | Mem (Kb) | CPU (%) | Notes                  |
| ---------------|-------------:| -----------:|---------:|--------:|:-----------------------|
| Go             | 247905.40    | 0:04.03     |   10520  |194      | GOMAXPROCS = 2, Go 1.9 |
| Go             | 219339.87    | 0:04.56     |   18052  |316      | GOMAXPROCS = 8, Go 1.9 |
| [Reference (C)]  | 214454.20    | 0:04.73     | 22428    | 99      |                        |
| Go             | 162336.18    | 0:06.16     |    8380  |99       | GOMAXPROCS = 1, Go 1.9 |
| [Node.js]        | 88983.80     | 0:11.33     |   72580  | 101     | Node.js 8.5.0   |

### Conditions
* Median value of 10 runs, lower elapsed time is better
* The benchmark implementation was run on the same machine as the server.
* Command:
```
$ time -f "Command: %C\nElapsed: %Es\nMEM (RSS): %Mkb\nCPU: %P\n"  redis-benchmark -n 1000000
```

## System info
**System** 
```
$ lscpu
Architecture:          x86_64
CPU op-mode(s):        32-bit, 64-bit
Byte Order:            Little Endian
CPU(s):                8
On-line CPU(s) list:   0-7
Thread(s) per core:    2
Core(s) per socket:    4
Socket(s):             1
NUMA node(s):          1
Vendor ID:             GenuineIntel
CPU family:            6
Model:                 94
Model name:            Intel(R) Core(TM) i7-6700HQ CPU @ 2.60GHz
Stepping:              3
CPU MHz:               1082.453
CPU max MHz:           3500.0000
CPU min MHz:           800.0000
BogoMIPS:              5183.91
Virtualisation:        VT-x
L1d cache:             32K
L1i cache:             32K
L2 cache:              256K
L3 cache:              6144K
```
**Memory**
```
$ free -m
              total        used        free      shared  buff/cache   available
Mem:          15901        6657        3525         217        5719        8494
Swap:         16239           0       16239
```

**Redis Server**: 4.0.2 64 bit

[redis-benchmark]: https://github.com/antirez/redis/blob/unstable/src/redis-benchmark.c
[Reference (C)]: https://github.com/antirez/redis/blob/unstable/src/redis-benchmark.c
[Pitfalls and misconceptions]: https://redis.io/topics/benchmarks#pitfalls-and-misconceptions
[Node.js]: https://github.com/rjnienaber/redis-benchmark-js
