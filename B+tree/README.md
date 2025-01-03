Data structure:

A node consists of:
1. A fixed-sized header containing the type of the node (leaf node or internal node)
2. The number of keys.
3. A list of pointers to the child nodes. (Used by internal nodes).
4. A list of offsets pointing to each key-value pair.
5. Packed KV pairs.

```
| type | nkeys | pointers   | offsets    | key-values
| 2B   | 2B    | nkeys * 8B | nkeys * 2B | ...
```

This is the format of the KV pair. Lengths followed by data.
```
| klen | vlen | key | val |
| 2B   | 2B   | ... | ... |
```