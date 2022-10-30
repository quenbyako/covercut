package sequence_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	. "github.com/quenbyako/covercut/cmd/sequence"
)

func TestCut(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name string
		cut  Sequence[int]
		obj  Sequence[int]
		want []Sequence[int]
	}{{
		cut:  Sequence[int]{Start: 0, Stop: 100},
		obj:  Sequence[int]{Start: 10, Stop: 15},
		want: []Sequence[int]{},
	}, {
		cut:  Sequence[int]{Start: 0, Stop: 100},
		obj:  Sequence[int]{Start: 101, Stop: 150},
		want: []Sequence[int]{{Start: 101, Stop: 150}},
	}, {
		cut:  Sequence[int]{Start: 500, Stop: 1000},
		obj:  Sequence[int]{Start: 101, Stop: 150},
		want: []Sequence[int]{{Start: 101, Stop: 150}},
	}, {
		cut:  Sequence[int]{Start: 10, Stop: 12},
		obj:  Sequence[int]{Start: 0, Stop: 30},
		want: []Sequence[int]{{Start: 0, Stop: 10}, {Start: 12, Stop: 30}},
	}, {
		cut:  Sequence[int]{Start: 0, Stop: 12},
		obj:  Sequence[int]{Start: 0, Stop: 30},
		want: []Sequence[int]{{Start: 12, Stop: 30}},
	}, {
		cut:  Sequence[int]{Start: 13, Stop: 30},
		obj:  Sequence[int]{Start: 0, Stop: 30},
		want: []Sequence[int]{{Start: 0, Stop: 13}},
	}, {
		cut:  Sequence[int]{Start: 5, Stop: 5},
		obj:  Sequence[int]{Start: 0, Stop: 10},
		want: []Sequence[int]{{Start: 0, Stop: 5}, {Start: 5, Stop: 10}},
	}} {
		tt := tt // for parallel tests

		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.obj.Cut(tt.cut))
		})
	}
}

func TestCmpLiberal(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name string
		obj  Sequence[int]
		i    int
		want int
	}{{
		obj:  Sequence[int]{Start: 0, Stop: 10},
		i:    9,
		want: 0,
	}, {
		obj:  Sequence[int]{Start: 0, Stop: 10},
		i:    0,
		want: -1,
	}, {
		obj:  Sequence[int]{Start: 0, Stop: 10},
		i:    10,
		want: 1,
	}, {
		obj:  Sequence[int]{Start: 0, Stop: 10},
		i:    100,
		want: 1,
	}, {
		obj:  Sequence[int]{Start: 0, Stop: 10},
		i:    -100,
		want: -1,
	}} {
		tt := tt // for parallel tests

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tt.want, tt.obj.CmpLiberal(tt.i))
		})
	}
}

func TestCmpGreedy(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name string
		obj  Sequence[int]
		i    int
		want int
	}{{
		obj:  Sequence[int]{Start: 0, Stop: 10},
		i:    9,
		want: 0,
	}, {
		obj:  Sequence[int]{Start: 0, Stop: 10},
		i:    0,
		want: 0,
	}, {
		obj:  Sequence[int]{Start: 0, Stop: 10},
		i:    10,
		want: 0,
	}, {
		obj:  Sequence[int]{Start: 0, Stop: 10},
		i:    100,
		want: 1,
	}, {
		obj:  Sequence[int]{Start: 0, Stop: 10},
		i:    -100,
		want: -1,
	}} {
		tt := tt // for parallel tests

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tt.want, tt.obj.CmpGreedy(tt.i))
		})
	}
}
