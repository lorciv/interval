// Package intervals implements operations to combine a set of overlapping intervals (e.g. time intervals)
// into a sequence that does not overlap. Priorities can be specified so that certain ranges prevail over others.
package interval

import (
	"fmt"
	"sort"
)

// Interval represents an interval of values from start to end (e.g. a time interval). Start is included,
// end is excluded. Intervals can be assigned a non-negative integer that represents its priority.
// Increasing numbers (0, 1, 2, 3...) represent decreasing priority, with 0 being the highest.
type Interval struct {
	Start, End int64
	Priority   int
}

func (v Interval) String() string {
	return fmt.Sprintf("{%d -> %d prio %d}", v.Start, v.End, v.Priority)
}

type eventType int

const (
	eventStart eventType = iota
	eventEnd
)

type event struct {
	time     int64
	typ      eventType
	priority int
}

func eventify(intervals []Interval) []event {
	var events []event

	for _, v := range intervals {
		if v.Start > v.End {
			v.Start, v.End = v.End, v.Start
		}

		events = append(events, event{
			typ:      eventStart,
			time:     v.Start,
			priority: v.Priority,
		}, event{
			typ:      eventEnd,
			time:     v.End,
			priority: v.Priority,
		})
	}

	sort.Slice(events, func(i, j int) bool {
		d := events[i].time - events[j].time
		if d == 0 {
			d = int64(events[i].typ) - int64(events[j].typ)
		}
		if d == 0 {
			d = int64(events[i].priority) - int64(events[j].priority)
		}
		return d < 0
	})

	return events
}

// Sequence combines a list of potentially overlapping intervals into a list of sequential intervals,
// which are guaranteed not to overlap. Generated intervals are assigned the highest priority (i.e. lowest value)
// computed from the input.
func Sequence(intervals []Interval) []Interval {
	var seq []Interval

	count := make(map[int]int)

	var current *Interval = nil

	for _, e := range eventify(intervals) {
		switch e.typ {
		case eventStart:
			count[e.priority]++
			if current != nil && current.Priority <= e.priority {
				break
			}
			// Current interval ends (if any)
			if current != nil {
				current.End = e.time
				seq = append(seq, *current)
			}
			// New interval starts at higher priority
			current = &Interval{
				Start:    e.time,
				Priority: e.priority,
			}
		case eventEnd:
			count[e.priority]--
			if count[current.Priority] > 0 {
				break
			}
			// Current interval ends
			current.End = e.time
			seq = append(seq, *current)
			// New interval starts at lower priority (if any)
			started := false
			prios := make([]int, 0, len(count))
			for p := range count {
				prios = append(prios, p)
			}
			sort.Ints(prios)
			for _, p := range prios {
				if count[p] > 0 {
					current = &Interval{
						Start:    e.time,
						Priority: p,
					}
					started = true
					break
				}
			}
			if !started {
				current = nil
			}
		}
	}

	return seq
}
