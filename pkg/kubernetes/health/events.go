package health

import (
	"time"

	"github.com/ghostsquad/go-timejumper"
	corev1 "k8s.io/api/core/v1"
)

const (
	// EventTypeWarning is for Warning events
	EventTypeWarning = "Warning"
)

// TrailingWarningEvent returns all "Warning" type events from the end of the list
// events should be sorted by "lastTimestamp"
func TrailingWarningEvent(events []corev1.Event) []corev1.Event {
	var warningEvents []corev1.Event

	for i := len(events) - 1; i >= 0; i-- {
		e := events[i]
		if e.Type == EventTypeWarning {
			warningEvents = append(warningEvents, e)
		} else {
			break
		}
	}

	return warningEvents
}

// WarningEventAgeOptions provide options for finding events
type WarningEventAgeOptions struct {
	clock timejumper.Clock
}

// WarningEventAgeOption provides for a functional API for func WarningEventsLessThanAge
type WarningEventAgeOption func(options *WarningEventAgeOptions)

// WarningEventsLessThanAge returns all "Warning" type events that have happened since the given time (1 second minimum resolution)
func WarningEventsLessThanAge(events []corev1.Event, dur time.Duration, options ...WarningEventAgeOption) []corev1.Event {
	opts := &WarningEventAgeOptions{
		clock: timejumper.RealClock{},
	}

	for _, ofn := range options {
		ofn(opts)
	}

	now := opts.clock.Now()

	var warningEvents []corev1.Event

	for i := range events {
		e := events[i]
		if e.Type == EventTypeWarning {
			diff := now.Sub(e.LastTimestamp.Time)

			if diff.Seconds() > dur.Seconds() {
				warningEvents = append(warningEvents, e)
			}
		}
	}

	return warningEvents
}
