// Package ranges implements operations to combine a set of overlapping intervals (e.g. time intervals)
// into a sequence that does not overlap. Priorities can be specified so that certain ranges prevail over others.
package ranges

import (
	"fmt"
	"math"
	"sort"
)

// Range represents an interval (e.g. a time interval) from start to end. Start is included,
// end is excluded. Ranges can be assigned a non-negative integer that represents its priority.
// Increasing numbers (0, 1, 2, 3...) represent decreasing priority, with 0 being the highest.
type Range struct {
	Start, End int
	Priority   int
}

func (p Range) String() string {
	return fmt.Sprintf("{%d -> %d prio %d}", p.Start, p.End, p.Priority)
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

func eventify(ranges []Range) []event {
	var events []event

	for _, r := range ranges {
		if r.Start > r.End {
			r.Start, r.End = r.End, r.Start
		}

		events = append(events, event{
			typ:      eventStart,
			time:     r.Start,
			priority: r.Priority,
		}, event{
			typ:      eventEnd,
			time:     r.End,
			priority: r.Priority,
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

// Overlap combines a list of potentially overlapping ranges into a list of sequential ranges, which are guaranteed not to overlap.
// Generated ranges are assigned the highest priority (i.e. lowest value) computed from the input.
func Overlap(ranges []Range) []Range {
	var sequence []Range

	curPrio := math.MaxInt // current priority
	var count multiCounter
	for _, e := range eventify(ranges) {
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
				sequence = append(sequence, Range{
					Start:    e.time,
					Priority: prio,
				})
			}
		}
		curPrio = prio
	}

	return sequence
}
