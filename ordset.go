// Package ordset implements an ordered set data structure using generics.
//
// An OrderedSet is a combination of a linked list and a hash map. As values are appended/prepended to the linked list,
// their list element pointers are stored in a hash map. This way, set membership checks (is element 'x' in the set?)
// can be done in constant time. Comparatively equal (==) elements are not duplicated when they are added to the set.
//
// Think of an OrderedSet as queue of its own keys, where each key can be quickly looked up to find its position in the queue.
//
// Example usage of an ordered set:
//
//  set := ordset.New[int]()
//  set.Append(1)
//  set.Append(2)
//  set.Append(3)
//  set.Prepend(0)
//  set.Prepend(0)
//
//  set.Has(1)  // true
//  set.Len()   // 4
//  set.Slice() // []int{0, 1, 2, 3}
//
//  set.Pop()   // 3
//  set.Slice() // []int{0, 1, 2}
//
//
// Elements in the set can be removed in constant-time while preserving the order of the queue.
//
//  set.Slice() // []int{0, 10, 20}
//  set.Remove(10)
//  set.Slice() // []int{0, 20}
//
// The set can be iterated through using the Range method.
//
//  set.Slice() // []string{"zero", "one", "two", "three"}
//  err := set.Range(func(i int, str string) error {
//    fmt.Println(i, str) // prints "0 zero", "1 one", etc
//    return nil
//  })
package ordset

import (
	"container/list"
	"errors"
)

var ErrMarkNotFound = errors.New("reference value for insert is not in the OrderedSet")

// OrderedSet is an implementation of an ordered set of elements of type T. Appending, prepending
// or inserting into the set stores values in a linked list which encodes the order of the elements,
// and in a hash map which allows for quick lookup of elements.
type OrderedSet[T comparable] struct {
	list    *list.List
	mapping map[T]*list.Element
}

// New initializes an OrderedSet of element type T with an initial set of elements.
func New[T comparable](elems ...T) *OrderedSet[T] {
	set := &OrderedSet[T]{
		list:    list.New(),
		mapping: make(map[T]*list.Element),
	}
	for _, v := range elems {
		set.Append(v)
	}
	return set
}

// Len returns the number of elements in the OrderedSet.
func (o *OrderedSet[T]) Len() int {
	return o.list.Len()
}

// Has returns true if the given value v is a member of the OrderedSet.
func (o *OrderedSet[T]) Has(v T) bool {
	_, exists := o.mapping[v]
	return exists
}

// Front returns the element at the front of the OrderedSet.
func (o *OrderedSet[T]) Front() T {
	return o.list.Front().Value.(T)
}

// Back returns the element at the front of the OrderedSet.
func (o *OrderedSet[T]) Back() T {
	return o.list.Back().Value.(T)
}

// Append pushes a value to the back of the OrderedSet.
func (o *OrderedSet[T]) Append(v T) bool {
	if o.Has(v) {
		return false
	}
	o.mapping[v] = o.list.PushBack(v)
	return true
}

// Prepend pushes a value to the front of the OrderedSet.
func (o *OrderedSet[T]) Prepend(v T) bool {
	if o.Has(v) {
		return false
	}
	o.mapping[v] = o.list.PushFront(v)
	return true
}

// Pop extracts and removes a value from the right of the OrderedSet. Returns a boolean true value
// if an element was successfully popped. This will only ever be false if the OrderedSet is empty.
func (o *OrderedSet[T]) Pop() (v T, ok bool) {
	if o.Len() > 0 {
		elem := o.list.Back()
		o.list.Remove(elem)
		v, ok = elem.Value.(T)
		delete(o.mapping, v)
	}
	return
}

// Shift extracts and removes a value from the left of the OrderedSet. Returns a boolean true value
// if an element was successfully popped. This will only ever be false if the OrderedSet is empty.
func (o *OrderedSet[T]) Shift() (v T, ok bool) {
	if o.Len() > 0 {
		elem := o.list.Front()
		o.list.Remove(elem)
		v, ok = elem.Value.(T)
		delete(o.mapping, v)
	}
	return
}

// Insert inserts the given value v into the OrderedSet at a specific position relative to the given mark value.
// If the after parameter is true, the value v is inserted immediately behind mark. If after is false, v is inserted
// immediately in front of mark.
//
// If the value v is already a member of the set, Insert is a no-op. Use the Move method to reorder set elements.
func (o *OrderedSet[T]) Insert(v, mark T, after bool) (added bool, err error) {
	if !o.Has(mark) {
		return false, ErrMarkNotFound
	} else if o.Has(v) {
		// value already exists in set, no-op
		return false, nil
	}

	if after {
		o.mapping[v] = o.list.InsertAfter(v, o.mapping[mark])
	} else {
		o.mapping[v] = o.list.InsertBefore(v, o.mapping[mark])
	}
	return true, nil
}

// Move reorders repositions the set element value v relative to the given mark value.
// If the after parameter is true, the value v is moved to immediately behind mark. If
// after is false, v is moved to immediately in front of mark.
func (o *OrderedSet[T]) Move(v, mark T, after bool) (err error) {
	if !o.Has(mark) {
		return ErrMarkNotFound
	}
	if elem, ok := o.mapping[v]; ok {
		if after {
			o.list.MoveAfter(elem, o.mapping[mark])
		} else {
			o.list.MoveBefore(elem, o.mapping[mark])
		}
	}
	return nil
}

func (o *OrderedSet[T]) Remove(v T) bool {
	if elem, ok := o.mapping[v]; ok {
		o.list.Remove(elem)
		delete(o.mapping, v)
		return true
	}
	return false
}

func (o *OrderedSet[T]) Range(loop func(int, T) error) error {
	i := 0
	for elem := o.list.Front(); elem != nil; elem = elem.Next() {
		err := loop(i, elem.Value.(T))
		if err != nil {
			return err
		}
		i++
	}
	return nil
}

func (o *OrderedSet[T]) RangeReverse(loop func(int, T) error) error {
	i := 0
	for elem := o.list.Back(); elem != nil; elem = elem.Prev() {
		err := loop(i, elem.Value.(T))
		if err != nil {
			return err
		}
		i++
	}
	return nil
}

func (o *OrderedSet[T]) Slice() []T {
	slice := make([]T, o.Len())
	o.Range(func(i int, v T) error {
		slice[i] = v
		return nil
	})
	return slice
}
