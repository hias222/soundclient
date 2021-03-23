package main

import (
	"log"
)

type Mapping struct {
	sound     *Soundcontrol
	connected bool
}

func NewMapping(s *Soundcontrol) (*Mapping, error) {

	newMapping := &Mapping{
		sound:     s,
		connected: false,
	}

	log.Println("Created socket instance")

	return newMapping, nil
}

func (m *Mapping) initialize() error {

	m.setupOnSliderMove()

	return nil
}

func (m *Mapping) handleSliderMoveEvent(event SliderMoveEvent) {
	// {"id": 1, "percent": 0.2}
	log.Println("todo move")
	log.Printf("slider %d percent %g", event.SliderID, event.PercentValue)
}

func (m *Mapping) setupOnSliderMove() {
	sliderEventsChannel := m.sound.socket.SubscribeToSliderMoveEvents()

	go func() {
		for {
			select {
			case event := <-sliderEventsChannel:
				m.handleSliderMoveEvent(event)
			}
		}
	}()
}
