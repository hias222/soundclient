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

func (m *Mapping) handleSliderMoveEvent(event SliderMoveEvent) {
	log.Println("todo move")
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