package sequence

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

type Sequence[T constraints.Ordered] struct {
	Start T
	Stop  T
}

func (r Sequence[T]) String() string {
	return fmt.Sprintf("Range{%v:%v}", r.Start, r.Stop)
}

// TODO: optimize it!!!
func (r Sequence[T]) CutMany(toCut ...Sequence[T]) []Sequence[T] {
	res := []Sequence[T]{r}
	for _, cutter := range toCut {
		more := []Sequence[T]{}
		for _, item := range res {
			more = append(more, item.Cut(cutter)...)
		}
		res = more
	}
	return res
}

func (r Sequence[T]) Cut(toCut Sequence[T]) []Sequence[T] {
	switch {
	case toCut.CmpGreedy(r.Start) == 0 && toCut.CmpGreedy(r.Stop) == 0:
		return []Sequence[T]{}
	case toCut.CmpGreedy(r.Start) > 0 || toCut.CmpGreedy(r.Stop) < 0:
		return []Sequence[T]{r}
	case toCut.Start > r.Start && toCut.Stop < r.Stop:
		return []Sequence[T]{{Start: r.Start, Stop: toCut.Start}, {Start: toCut.Stop, Stop: r.Stop}}
	case toCut.Start > r.Start && toCut.Stop >= r.Stop:
		return []Sequence[T]{{Start: r.Start, Stop: toCut.Start}}
	case toCut.Start <= r.Start && toCut.Stop < r.Stop:
		return []Sequence[T]{{Start: toCut.Stop, Stop: r.Stop}}
	default:
		panic(fmt.Sprintf("whaaat? %v %v", r, toCut))
	}
}

func (r Sequence[T]) CmpGreedy(t T) int {
	switch {
	case r.Start > t:
		return -1
	case r.Stop >= t:
		return 0
	default:
		return 1
	}
}

func (r Sequence[T]) CmpLiberal(t T) int {
	switch {
	case r.Start >= t:
		return -1
	case r.Stop > t:
		return 0
	default:
		return 1
	}
}
