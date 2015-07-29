package controllers

import (
	"container/list"
	"time"
	// "strings"
	"fmt"
)

const ARCHIVE_SIZE = 20

var (
	archive = list.New() // Event archives.
)

type EventType int

const (
	EVENT_JOIN    = iota //0
	EVENT_LEAVE          //1
	EVENT_MESSAGE        //2
)

type Event struct {
	Type      EventType // JOIN, LEAVE, MESSAGE
	User      string
	Timestamp int // Unix timestamp (secs)
	Content   string
}

func (e *Event) String() string {
	return fmt.Sprintf("type %d User: %s  Content: %s", e.Type, e.User, e.Content)
}

func newEvent(ep EventType, user, msg string) Event {
	return Event{ep, user, int(time.Now().Unix()), msg}
}

// NewArchive saves new event to archive list.
func NewArchive(event Event) {
	if archive.Len() >= ARCHIVE_SIZE {
		archive.Remove(archive.Front())
	}
	archive.PushBack(event)
}

// GetEvents returns all events after lastReceived.
func GetEvents(lastReceived int) []Event {
	events := make([]Event, 0, archive.Len())
	for event := archive.Front(); event != nil; event = event.Next() {
		e := event.Value.(Event)
		if e.Timestamp > int(lastReceived) {
			events = append(events, e)
		}
	}
	return events
}
