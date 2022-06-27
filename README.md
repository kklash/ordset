# Ordered Set

Package `ordset` implements an ordered set data structure using generics.

An `OrderedSet` is a combination of a linked list and a hash map. As values are appended/prepended to the linked list,
their list element pointers are stored in a hash map. This way, set membership checks (is element 'x' in the set?)
can be done in constant time. Comparatively equal (==) elements are not duplicated when they are added to the set.

Think of an `OrderedSet` as queue of its own keys, where each key can be quickly looked up to find its position in the queue.

Example usage of an ordered set:

```go
set := ordset.New[int]()
set.Append(1)
set.Append(2)
set.Append(3)
set.Prepend(0)
set.Prepend(0)

set.Has(1)  // true
set.Len()   // 4
set.Slice() // []int{0, 1, 2, 3}

set.Pop()   // 3
set.Slice() // []int{0, 1, 2}
```

Elements in the set can be removed in constant-time while preserving the order of the queue.

```go
set.Slice() // []int{0, 10, 20}
set.Remove(10)
set.Slice() // []int{0, 20}
```

An `OrderedSet` can be iterated through using the Range method.

```go
set.Slice() // []string{"zero", "one", "two", "three"}
err := set.Range(func(i int, str string) error {
  fmt.Println(i, str) // prints "0 zero", "1 one", etc
  return nil
})
```

### [See here for full documentation](https://pkg.go.dev/github.com/kklash/ordset)
