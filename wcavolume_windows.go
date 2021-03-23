package main

import (
	"fmt"

	"github.com/go-ole/go-ole"
	"github.com/moutend/go-wca/pkg/wca"
)

var version = "latest"
var revision = "latest"

type VolumeFlag struct {
	Value float32
	IsSet bool
}

func endpointVolume(volumeFlag VolumeFlag) (err error) {
	fmt.Println(volumeFlag)
	if err = ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED); err != nil {
		return
	}
	defer ole.CoUninitialize()

	var mmde *wca.IMMDeviceEnumerator
	if err = wca.CoCreateInstance(wca.CLSID_MMDeviceEnumerator, 0, wca.CLSCTX_ALL, wca.IID_IMMDeviceEnumerator, &mmde); err != nil {
		return
	}
	defer mmde.Release()

	var mmd *wca.IMMDevice
	if err = mmde.GetDefaultAudioEndpoint(wca.ERender, wca.EConsole, &mmd); err != nil {
		return
	}
	defer mmd.Release()

	var ps *wca.IPropertyStore
	if err = mmd.OpenPropertyStore(wca.STGM_READ, &ps); err != nil {
		return
	}
	defer ps.Release()

	var pv wca.PROPVARIANT
	if err = ps.GetValue(&wca.PKEY_Device_FriendlyName, &pv); err != nil {
		return
	}
	fmt.Printf("%s\n", pv.String())

	var aev *wca.IAudioEndpointVolume
	if err = mmd.Activate(wca.IID_IAudioEndpointVolume, wca.CLSCTX_ALL, nil, &aev); err != nil {
		return
	}
	defer aev.Release()

	if volumeFlag.IsSet {
		if err = aev.SetMasterVolumeLevelScalar(volumeFlag.Value, nil); err != nil {
			return
		}
	}

	var channels uint32
	if err = aev.GetChannelCount(&channels); err != nil {
		return
	}

	var mute bool
	if err = aev.GetMute(&mute); err != nil {
		return
	}

	var masterVolumeLevel float32
	if err = aev.GetMasterVolumeLevel(&masterVolumeLevel); err != nil {
		return
	}

	var masterVolumeLevelScalar float32
	if err = aev.GetMasterVolumeLevelScalar(&masterVolumeLevelScalar); err != nil {
		return
	}

	fmt.Println("--------")
	fmt.Printf("Channels: %d\n", channels)
	fmt.Printf("Mute state: %v\n", mute)
	fmt.Println("Master volume level:")
	fmt.Printf("  %v [dB]\n", masterVolumeLevel)
	fmt.Printf("  %v [scalar]\n", masterVolumeLevelScalar)
	fmt.Println("--------")

	return
}
