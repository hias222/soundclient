package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type MessageSocket struct {
	serverURL string

	stopChannel chan bool
	connected   bool

	lastKnownNumSliders        int
	currentSliderPercentValues []float32

	sliderMoveConsumers []chan SliderMoveEvent
}

type SliderMoveEvent struct {
	SliderID     int
	PercentValue float32
}

func AllUsers() {
	fmt.Println("All Users")
}

func NewWebsocket() (*MessageSocket, error) {

	newSocket := &MessageSocket{
		serverURL:           "ws://localhost:8080" + "/ws",
		stopChannel:         make(chan bool),
		connected:           false,
		sliderMoveConsumers: []chan SliderMoveEvent{},
	}

	log.Println("Created socket instance")

	return newSocket, nil
}

func (newSocket *MessageSocket) receiveHandler(connection *websocket.Conn) {
	defer close(done)
	for {
		_, msg, err := connection.ReadMessage()
		if err != nil {
			log.Println("Error in receive:", err)
			return
		}
		newSocket.testit(msg)
	}
}

func (newSocket *MessageSocket) testit(msg []byte) {
	log.Printf("For Max: %s\n", msg)

	moveEvents := []SliderMoveEvent{}

	moveEvents = append(moveEvents, SliderMoveEvent{
		SliderID:     1,
		PercentValue: 0.1,
	})

	for _, consumer := range newSocket.sliderMoveConsumers {
		for _, moveEvent := range moveEvents {
			consumer <- moveEvent
		}
	}

}

func (newSocket *MessageSocket) Start() error {

	if newSocket.connected {
		log.Println("Already connected, can't start another without closing first")
		return errors.New("serial: connection already active")
	}

	conn, _, err := websocket.DefaultDialer.Dial(newSocket.serverURL, nil)
	if err != nil {
		log.Fatal("webodels: Error connecting to Websocket Server:", err)
	}
	//defer conn.Close()
	newSocket.connected = true
	go newSocket.receiveHandler(conn)

	log.Println("webmodel start ....")

	return nil

}

func (newSocket *MessageSocket) Stop() {
	if newSocket.connected {
		log.Println("Shutting down socket connection")
		//newSocket.stopChannel <- true
	} else {
		log.Println("Not currently connected, nothing to stop")
	}
}

func (newSocket *MessageSocket) SubscribeToSliderMoveEvents() chan SliderMoveEvent {
	ch := make(chan SliderMoveEvent)
	newSocket.sliderMoveConsumers = append(newSocket.sliderMoveConsumers, ch)

	return ch
}
