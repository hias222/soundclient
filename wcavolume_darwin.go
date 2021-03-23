package main

import (
	"fmt"
)

var version = "latest"
var revision = "latest"

type VolumeFlag struct {
	Value float32
	IsSet bool
}

func endpointVolume(volumeFlag VolumeFlag) (err error) {
	fmt.Println(volumeFlag)
	fmt.Println("------not implemented --")
	return
}
