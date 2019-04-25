package sniffer

import (
	"sync"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

// Sniffer is a container of all the sniffer-related data
type Sniffer struct {
	XboxAdapter *pcap.Handle
	Packets     chan Packet
	readBytes   bool
	waitGroup   *sync.WaitGroup
}

type Packet struct {
	Data        []byte
	CaptureInfo gopacket.CaptureInfo
}

/* Public Functions */

// Start attempts to start the sniffer and return a pointer to the struct
func Start() (*Sniffer, error) {
	sniffer, err := createSniffer()
	if err != nil {
		return nil, err
	}
	sniffer.readBytes = true
	sniffer.waitGroup.Add(1)
	go sniffer.beginRead()
	return sniffer, nil
}

// Stop attempts to stop the sniffer
func (s *Sniffer) Stop() {
	s.readBytes = false
	s.waitGroup.Wait()
	return
}

/* Private Functions */

func createSniffer() (*Sniffer, error) {
	handle, err := getXboxAdapterHandle()
	if err != nil {
		return nil, err
	}
	sniffer := &Sniffer{
		XboxAdapter: handle,
		Packets:     make(chan Packet),
		waitGroup:   &sync.WaitGroup{},
	}
	return sniffer, nil
}

func (s *Sniffer) beginRead() {
	for s.readBytes {
		data, captureInfo, err := s.XboxAdapter.ReadPacketData()
		if err != nil {
			continue
		}
		s.Packets <- Packet{
			Data:        data,
			CaptureInfo: captureInfo,
		}
	}
	s.waitGroup.Done()
}

func getXboxAdapterInterface() (*pcap.Interface, error) {
	interfaces, err := pcap.FindAllDevs()
	if err != nil {
		return nil, err
	}
	for _, dev := range interfaces {
		if len(dev.Addresses) == 0 {
			return &dev, nil
		}
	}
	return nil, errorNotFound
}

func getXboxAdapterHandle() (*pcap.Handle, error) {
	adapterIf, err := getXboxAdapterInterface()
	if err != nil {
		return nil, err
	}
	handle, err := pcap.OpenLive(adapterIf.Name, 45, true, 50)
	if err != nil {
		return nil, err
	}
	return handle, nil
}
