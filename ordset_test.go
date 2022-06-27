package ordset_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/kklash/ordset"
)

type T struct{ x int }

func TestOrderedSet(t *testing.T) {
	t.Run("Append/Prepend/Front/Back/Has/Length", func(t *testing.T) {
		set := ordset.New[int]()
		if !set.Append(3) {
			t.Errorf("expected set.Append(3) to return true")
			return
		}
		if !set.Append(5) {
			t.Errorf("expected set.Append(5) to return true")
			return
		}
		if !set.Prepend(1) {
			t.Errorf("expected set.Prepend(1) to return true")
			return
		}

		if set.Append(5) {
			t.Errorf("expected set.Append(5) again to return false")
			return
		}
		if set.Prepend(1) {
			t.Errorf("expected set.Append(1) to return false")
			return
		}

		if !set.Has(3) {
			t.Errorf("expected set to contain 3")
			return
		}
		if !set.Has(5) {
			t.Errorf("expected set to contain 5")
			return
		}
		if !set.Has(1) {
			t.Errorf("expected set to contain 1")
			return
		}
		if set.Has(12) {
			t.Errorf("expected set NOT to contain 12")
			return
		}
		if set.Len() != 3 {
			t.Errorf("expected set length to be 3")
			return
		}

		if v := set.Front(); v != 1 {
			t.Errorf("expected Front to return 1, got %d", v)
			return
		}
		if v := set.Back(); v != 5 {
			t.Errorf("expected Back to return 5, got %d", v)
			return
		}
	})

	t.Run("set of struct values", func(t *testing.T) {
		structs := []T{
			{0}, {3}, {4},
			{1}, {19}, {21},

			// duplicate struct values, should not be stored in final set
			{4}, {3}, {0},
		}

		set := ordset.New[T](structs...)

		resultStructs := set.Slice()
		if !reflect.DeepEqual(resultStructs, structs[:6]) {
			t.Errorf("failed to produce expected set of structs")
			return
		}
	})

	t.Run("set of struct pointers", func(t *testing.T) {
		structs := []*T{
			{0}, {3}, {4},
			{1}, {19}, {21},

			// duplicate struct values, different pointers
			{4}, {3}, {0},
		}

		set := ordset.New[*T](structs...)

		resultStructs := set.Slice()
		if !reflect.DeepEqual(resultStructs, structs) {
			t.Errorf("failed to produce expected set of pointers")
			return
		}
	})

	t.Run("Insert", func(t *testing.T) {
		set := ordset.New[int](1, 2, 3, 5)

		added, err := set.Insert(4, 3, true)
		if err != nil {
			t.Errorf("failed to insert value 4: %s", err)
			return
		} else if !added {
			t.Errorf("insert of value 4 did not complete")
			return
		}

		result := set.Slice()
		expected := []int{1, 2, 3, 4, 5}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("failed to produce expected set of ints after insertion")
			return
		}

		added, err = set.Insert(4, 3, false)
		if err != nil {
			t.Errorf("failed to re-insert value 4: %s", err)
			return
		} else if added {
			t.Errorf("expected re-insertion of value 4 to return added=false")
			return
		}

		result = set.Slice()
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("set should not have changed after re-insertion")
			return
		}

		_, err = set.Insert(4, 99999, true)
		if !errors.Is(err, ordset.ErrMarkNotFound) {
			t.Errorf("expected ErrMarkNotFound when inserting using mark not in set, got %v", err)
			return
		}
	})

	t.Run("Move", func(t *testing.T) {
		set := ordset.New[int](1, 2, 4, 3, 5)

		if err := set.Move(4, 3, true); err != nil {
			t.Errorf("failed to move value around in set: %s", err)
			return
		}

		result := set.Slice()
		expected := []int{1, 2, 3, 4, 5}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("failed to produce expected set of ints after move")
			return
		}

		if err := set.Move(4, 999, true); !errors.Is(err, ordset.ErrMarkNotFound) {
			t.Errorf("expected ErrMarkNotFound when moving using mark not in set, got %v", err)
			return
		}
	})

	t.Run("Pop", func(t *testing.T) {
		values := []int{1, 2, 3, 4, 5}
		set := ordset.New[int](values...)

		for i := len(values) - 1; i >= 0; i-- {
			v, ok := set.Pop()
			if !ok {
				t.Errorf("expected pop to return ok=true")
				return
			}
			if v != values[i] {
				t.Errorf("popped value did not match: %d != %d", v, values[i])
				return
			}
		}

		if v, ok := set.Pop(); ok || v != 0 {
			t.Errorf("expected pop of empty set to return 0, false, got (%d, %v)", v, ok)
			return
		}
	})

	t.Run("Shift", func(t *testing.T) {
		values := []int{1, 2, 3, 4, 5}
		set := ordset.New[int](values...)

		for i := 0; i < len(values); i++ {
			v, ok := set.Shift()
			if !ok {
				t.Errorf("expected shift to return ok=true")
				return
			}
			if v != values[i] {
				t.Errorf("popped value did not match: %d != %d", v, values[i])
				return
			}
		}
		if v, ok := set.Shift(); ok || v != 0 {
			t.Errorf("expected shift of empty set to return 0, false, got (%d, %v)", v, ok)
			return
		}
	})
}
