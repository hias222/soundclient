package main

import (
	"encoding/json"
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

type JsonMessage struct {
	Messagetype int    `json:"type"`
	MessageBody string `json:"body"`
}

type SoundMessage struct {
	Type    int `json:"type"`
	Message SliderMessage
}

type SliderMessage struct {
	SliderID     int     `json:"id,omitempty"`
	PercentValue float32 `json:"percent,omitempty"`
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
		newSocket.mapToJsonMsg(msg)
	}
}

func (newSocket *MessageSocket) mapToJsonMsg(msg []byte) {
	var soundMessage SoundMessage
	sliderError := json.Unmarshal(msg, &soundMessage)

	if sliderError != nil {
		return
	}

	if soundMessage.Type != 2 {
		return
	}

	moveEvents := []SliderMoveEvent{}

	moveEvents = append(moveEvents, SliderMoveEvent{
		SliderID:     soundMessage.Message.SliderID,
		PercentValue: soundMessage.Message.PercentValue,
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
	log.Println("SubscribeToSliderMoveEvents")

	return ch
}
