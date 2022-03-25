// Package ranges implements operations to combine a set of overlapping intervals (e.g. time intervals)
// into a sequence that does not overlap. Priorities can be specified so that certain ranges prevail over others.
package ranges

import (
	"fmt"
	"math"
	"sort"
)

// Interval represents an interval of values from start to end (e.g. a time interval). Start is included,
// end is excluded. Intervals can be assigned a non-negative integer that represents its priority.
// Increasing numbers (0, 1, 2, 3...) represent decreasing priority, with 0 being the highest.
type Interval struct {
	Start, End int
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
	time     int
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
			d = int(events[i].typ) - int(events[j].typ)
		}
		if d == 0 {
			d = events[i].priority - events[j].priority
		}
		return d < 0
	})

	return events
}

type multiCounter []int

func (m *multiCounter) incr(priority int) {
	for len(*m) < priority+1 {
		*m = append(*m, 0)
	}
	(*m)[priority]++
}

func (m *multiCounter) decr(priority int) {
	if len(*m) < priority+1 || (*m)[priority] == 0 {
		panic("illegal counter decrement")
	}
	(*m)[priority]--
}

// Overlap combines a list of potentially overlapping intervals into a list of sequential intervals,
// which are guaranteed not to overlap. Generated intervals are assigned the highest priority (i.e. lowest value)
// computed from the input.
func Overlap(intervals []Interval) []Interval {
	var sequence []Interval

	curPrio := math.MaxInt // current priority
	var count multiCounter
	for _, e := range eventify(intervals) {
		switch e.typ {
		case eventStart:
			count.incr(e.priority)
		case eventEnd:
			count.decr(e.priority)
		}

		prio := math.MaxInt // priority after last event
		for i := 0; i < len(count); i++ {
			if count[i] > 0 {
				prio = i
				break
			}
		}

		if prio != curPrio {
			last := len(sequence) - 1
			if last >= 0 && curPrio < math.MaxInt {
				sequence[last].End = e.time
			}
			if prio < math.MaxInt {
				sequence = append(sequence, Interval{
					Start:    e.time,
					Priority: prio,
				})
			}
		}
		curPrio = prio
	}

	return sequence
}
