package ordset_test

import (
	"errors"
	"fmt"

	"github.com/kklash/ordset"
)

func ExampleOrderedSet_Range() {
	set := ordset.New[string]("zero", "one", "two", "three")

	err := set.Range(func(i int, str string) error {
		fmt.Println(i, str)
		if str == "two" {
			return errors.New("reached two, giving up now")
		}
		return nil
	})

	fmt.Println(err)

	// output:
	//
	// 0 zero
	// 1 one
	// 2 two
	// reached two, giving up now
}

func ExampleOrderedSet_Insert() {
	set := ordset.New[int](10, 20, 30, 50, 60)

	// Inserts 40 after 50
	added, err := set.Insert(40, 50, false)
	if err != nil {
		panic(err)
	}
	fmt.Println("1 - added 40:", added)
	fmt.Println("1 - set:", set.Slice())

	// No-op as 40 is already a member of the set
	added, err = set.Insert(40, 10, true)
	if err != nil {
		panic(err)
	}
	fmt.Println("2 - added 40:", added)
	fmt.Println("2 - set:", set.Slice())

	_, err = set.Insert(40, 9999, true)
	fmt.Println("3 - error:", err)

	// output:
	//
	// 1 - added 40: true
	// 1 - set: [10 20 30 40 50 60]
	// 2 - added 40: false
	// 2 - set: [10 20 30 40 50 60]
	// 3 - error: reference value for insert is not in the OrderedSet
}

func ExampleOrderedSet_Move() {
	set := ordset.New[int](100, 200, 150, 250)
	if err := set.Move(150, 200, false); err != nil {
		panic(err)
	}

	fmt.Println(set.Slice())

	// output:
	//
	// [100 150 200 250]
}
