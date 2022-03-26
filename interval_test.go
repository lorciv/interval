package interval

import "testing"

func TestOverlap(t *testing.T) {
	tests := []struct {
		descr string
		input []Interval
		want  []Interval
	}{
		{
			descr: "disjoint",
			input: []Interval{
				{Start: 5, End: 10, Priority: 0},
				{Start: 20, End: 25, Priority: 2},
				{Start: 12, End: 16, Priority: 1},
			},
			want: []Interval{
				{Start: 5, End: 10, Priority: 0},
				{Start: 12, End: 16, Priority: 1},
				{Start: 20, End: 25, Priority: 2},
			},
		},
		{
			descr: "disjoint inverted",
			input: []Interval{
				{Start: 10, End: 5, Priority: 0},
				{Start: 25, End: 20, Priority: 2},
				{Start: 12, End: 16, Priority: 1},
			},
			want: []Interval{
				{Start: 5, End: 10, Priority: 0},
				{Start: 12, End: 16, Priority: 1},
				{Start: 20, End: 25, Priority: 2},
			},
		},
		{
			descr: "adjacent",
			input: []Interval{
				{Start: 5, End: 10, Priority: 0},
				{Start: 10, End: 20, Priority: 1},
				{Start: 20, End: 23, Priority: 2},
				{Start: 23, End: 26, Priority: 2},
			},
			want: []Interval{
				{Start: 5, End: 10, Priority: 0},
				{Start: 10, End: 20, Priority: 1},
				{Start: 20, End: 26, Priority: 2},
			},
		},
		{
			descr: "overlap",
			input: []Interval{
				{Start: 5, End: 13, Priority: 0},
				{Start: 10, End: 20, Priority: 1},
				{Start: 17, End: 25, Priority: 2},
			},
			want: []Interval{
				{Start: 5, End: 13, Priority: 0},
				{Start: 13, End: 20, Priority: 1},
				{Start: 20, End: 25, Priority: 2},
			},
		},
		{
			descr: "overlap2",
			input: []Interval{
				{Start: 5, End: 25, Priority: 3},
				{Start: 10, End: 20, Priority: 2},
				{Start: 13, End: 17, Priority: 1},
			},
			want: []Interval{
				{Start: 5, End: 10, Priority: 3},
				{Start: 10, End: 13, Priority: 2},
				{Start: 13, End: 17, Priority: 1},
				{Start: 17, End: 20, Priority: 2},
				{Start: 20, End: 25, Priority: 3},
			},
		},
		{
			descr: "overlap3",
			input: []Interval{
				{Start: 5, End: 10, Priority: 3},
				{Start: 10, End: 25, Priority: 2},
				{Start: 10, End: 20, Priority: 1},
				{Start: 10, End: 15, Priority: 0},
			},
			want: []Interval{
				{Start: 5, End: 10, Priority: 3},
				{Start: 10, End: 15, Priority: 0},
				{Start: 15, End: 20, Priority: 1},
				{Start: 20, End: 25, Priority: 2},
			},
		},
		{
			descr: "single range",
			input: []Interval{
				{Start: 5, End: 23, Priority: 0},
			},
			want: []Interval{
				{Start: 5, End: 23, Priority: 0},
			},
		},
		{
			descr: "negative start",
			input: []Interval{
				{Start: -5, End: 23, Priority: 0},
			},
			want: []Interval{
				{Start: -5, End: 23, Priority: 0},
			},
		},
		{
			descr: "non-consecutive priorities",
			input: []Interval{
				{Start: 0, End: 10, Priority: 4},
				{Start: 3, End: 6, Priority: 2},
			},
			want: []Interval{
				{Start: 0, End: 3, Priority: 4},
				{Start: 3, End: 6, Priority: 2},
				{Start: 6, End: 10, Priority: 4},
			},
		},
		{
			descr: "stacked",
			input: []Interval{
				{Start: 3, End: 10, Priority: 4},
				{Start: 3, End: 10, Priority: 2},
				{Start: 3, End: 10, Priority: 5},
			},
			want: []Interval{
				{Start: 3, End: 10, Priority: 2},
			},
		},
	}

	for _, test := range tests {
		got := Sequence(test.input)
		if len(got) != len(test.want) {
			t.Fatalf("Overlap(%s): got %d items, want %d", test.descr, len(got), len(test.want))
		}
		for i := 0; i < len(got); i++ {
			if got[i] != test.want[i] {
				t.Errorf("Overlap(%s): interval #%d: got %s, want %s", test.descr, i, got[i], test.want[i])
			}
		}
	}
}
