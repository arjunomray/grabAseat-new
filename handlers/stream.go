package handlers

import (
	"arjunomray/eventer/realtime"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// StreamEvents streams seat-availability updates to the client as
// Server-Sent Events. It is intentionally public: seat counts are not secret,
// and EventSource cannot send an Authorization header.
func StreamEvents(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("X-Accel-Buffering", "no")
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Flush()

	updates := realtime.Default.Subscribe()
	defer realtime.Default.Unsubscribe(updates)

	// Periodic comment keeps proxies from closing an idle connection.
	ping := time.NewTicker(25 * time.Second)
	defer ping.Stop()

	for {
		select {
		case <-c.Request.Context().Done():
			return
		case update, ok := <-updates:
			if !ok {
				return
			}
			data, err := json.Marshal(update)
			if err != nil {
				continue
			}
			if _, err := fmt.Fprintf(c.Writer, "data: %s\n\n", data); err != nil {
				return
			}
			c.Writer.Flush()
		case <-ping.C:
			if _, err := fmt.Fprint(c.Writer, ": ping\n\n"); err != nil {
				return
			}
			c.Writer.Flush()
		}
	}
}
