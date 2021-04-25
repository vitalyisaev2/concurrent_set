# Concurrent set data structure 
This repository contains benchmarks for different implementations of linked list based concurrent set. 
The algorithms are taken from The Art of Multiprocessor Programming, 2nd edition, 2021, by Herlihy, Shavit, Luchangco and Spear.

## Concurrent write
- Each thread tries to insert 1024 items to the set.
![](report/insert_ascending_array.svg)
![](report/insert_shuffled_array.svg)

## Concurrent read 
- Each thread tries to seek the items in the pre-prepared set consisting of 1024 items.
![](report/contains_ascending_array.svg)
![](report/contains_shuffled_array.svg)
  
## Concurrent write and read
- Half of the threads are inserting 1024 items to the set, while the other half is seeking for the items.
![](report/insert_and_contains_ascending_array.svg)
![](report/insert_and_contains_shuffled_array.svg)
