package events_service

import (
	"fmt"
	"github.com/timhugh/digitalvenue/core"
)

func main() {
	event := core.WebhookEvent{
		EventId: "123",
		Type:    "user.created",
	}
	fmt.Printf("Event: %+v\n", event)
}
