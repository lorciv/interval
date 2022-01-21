package ranges

import "testing"

func TestOverlap(t *testing.T) {
	tests := []struct {
		descr string
		input []Range
		want  []Range
	}{
		{
			descr: "disjoint",
			input: []Range{
				{Start: 5, End: 10, Priority: 0},
				{Start: 20, End: 25, Priority: 2},
				{Start: 12, End: 16, Priority: 1},
			},
			want: []Range{
				{Start: 5, End: 10, Priority: 0},
				{Start: 12, End: 16, Priority: 1},
				{Start: 20, End: 25, Priority: 2},
			},
		},
		{
			descr: "adjacent",
			input: []Range{
				{Start: 5, End: 10, Priority: 0},
				{Start: 10, End: 20, Priority: 1},
				{Start: 20, End: 23, Priority: 2},
				{Start: 23, End: 26, Priority: 2},
			},
			want: []Range{
				{Start: 5, End: 10, Priority: 0},
				{Start: 10, End: 20, Priority: 1},
				{Start: 20, End: 26, Priority: 2},
			},
		},
		{
			descr: "overlap",
			input: []Range{
				{Start: 5, End: 13, Priority: 0},
				{Start: 10, End: 20, Priority: 1},
				{Start: 17, End: 25, Priority: 2},
			},
			want: []Range{
				{Start: 5, End: 13, Priority: 0},
				{Start: 13, End: 20, Priority: 1},
				{Start: 20, End: 25, Priority: 2},
			},
		},
		{
			descr: "overlap2",
			input: []Range{
				{Start: 5, End: 25, Priority: 3},
				{Start: 10, End: 20, Priority: 2},
				{Start: 13, End: 17, Priority: 1},
			},
			want: []Range{
				{Start: 5, End: 10, Priority: 3},
				{Start: 10, End: 13, Priority: 2},
				{Start: 13, End: 17, Priority: 1},
				{Start: 17, End: 20, Priority: 2},
				{Start: 20, End: 25, Priority: 3},
			},
		},
		{
			descr: "overlap3",
			input: []Range{
				{Start: 5, End: 10, Priority: 3},
				{Start: 10, End: 25, Priority: 2},
				{Start: 10, End: 20, Priority: 1},
				{Start: 10, End: 15, Priority: 0},
			},
			want: []Range{
				{Start: 5, End: 10, Priority: 3},
				{Start: 10, End: 15, Priority: 0},
				{Start: 15, End: 20, Priority: 1},
				{Start: 20, End: 25, Priority: 2},
			},
		},
	}

	for _, test := range tests {
		got := Overlap(test.input)
		if len(got) != len(test.want) {
			t.Fatalf("Overlap(%s): got %d items, want %d", test.descr, len(got), len(test.want))
		}
		for i := 0; i < len(got); i++ {
			if got[i] != test.want[i] {
				t.Errorf("Overlap(%s): range #%d: got %s, want %s", test.descr, i, got[i], test.want[i])
			}
		}
	}
}
