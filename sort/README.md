# sort

| Sort          | Best Case         | Average Case      | Worst Case        | Memory | Comment                              |
|---------------|:-----------------:|:-----------------:|:-----------------:|:------:|--------------------------------------|
| InsertionSort | N<sup>2</sup> / 2 | N<sup>2</sup> / 2 | N<sup>2</sup> / 2 | 1      | N exchanges, Suitable for small N    |
| ShellSort     | N                 | ?                 | ?                 | 1      | Tight code, Subquadratic             |
| MergeSort     | NlgN              | NlgN              | NlgN              | N      | Stable, NlgN guarantee, Extra memory |
| MergeSortRec  | NlgN              | NlgN              | NlgN              | N      | Stable, NlgN guarantee, Extra memory |
| HeapSort      | NlgN              | NlgN              | NlgN              | 1      | NlgN guarantee                       |
| QuickSort     | NlgN              | 2NlnN             | N<sup>2</sup> / 2 | clogN  | NlgN probabilistic guarantee         |
| QuickSort3Way | N                 | 2NlnN             | N<sup>2</sup> / 2 | clogN  | Faster in presence of duplicate keys |

By running benchmarks, you can compare the performance of these algorithms with each other and also with built-in sort algorithm of go.
