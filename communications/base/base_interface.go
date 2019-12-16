package base

import (
	"errors"
	"time"

<<<<<<< HEAD
	"github.com/idoall/gocryptotrader/config"
	"github.com/idoall/gocryptotrader/exchanges/orderbook"
	"github.com/idoall/gocryptotrader/exchanges/ticker"
	log "github.com/idoall/gocryptotrader/logger"
=======
	"github.com/thrasher-corp/gocryptotrader/config"
	log "github.com/thrasher-corp/gocryptotrader/logger"
>>>>>>> upstrem/master
)

// IComm is the main interface array across the communication packages
type IComm []ICommunicate

// ICommunicate enforces standard functions across communication packages
type ICommunicate interface {
	Setup(config *config.CommunicationsConfig)
	Connect() error
	PushEvent(Event) error
	IsEnabled() bool
	IsConnected() bool
	GetName() string
}

// Setup sets up communication variables and intiates a connection to the
// communication mediums
func (c IComm) Setup() {
	ServiceStarted = time.Now()
	for i := range c {
		if c[i].IsEnabled() && !c[i].IsConnected() {
			err := c[i].Connect()
			if err != nil {
				log.Errorf(log.CommunicationMgr, "Communications: %s failed to connect. Err: %s", c[i].GetName(), err)
				continue
			}
			log.Debugf(log.CommunicationMgr, "Communications: %v is enabled and online.", c[i].GetName())
		}
	}
}

// PushEvent pushes triggered events to all enabled communication links
func (c IComm) PushEvent(event Event) {
	for i := range c {
		if c[i].IsEnabled() && c[i].IsConnected() {
			err := c[i].PushEvent(event)
			if err != nil {
				log.Errorf(log.CommunicationMgr, "Communications error - PushEvent() in package %s with %v. Err %s",
					c[i].GetName(), event, err)
			}
		}
	}
}

// GetStatus returns the status of the comms relayers
func (c IComm) GetStatus() map[string]CommsStatus {
	result := make(map[string]CommsStatus)
	for x := range c {
		result[c[x].GetName()] = CommsStatus{
			Enabled:   c[x].IsEnabled(),
			Connected: c[x].IsConnected(),
		}
	}
	return result
}

// GetEnabledCommunicationMediums prints out enabled and connected communication
// packages
// (#debug output only)
func (c IComm) GetEnabledCommunicationMediums() error {
	var count int
	for i := range c {
		if c[i].IsEnabled() && c[i].IsConnected() {
			log.Debugf(log.CommunicationMgr, "Communications: Medium %s is enabled.", c[i].GetName())
			count++
		}
	}
	if count == 0 {
		return errors.New("no communication mediums are enabled")
	}
	return nil
}
