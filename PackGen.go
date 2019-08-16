package main

import (
	"fmt"
	"net"
	"syscall"
	"log"
	"time"
	"flag"
	"sync"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

type connDetails struct {
	srcMac	string	
	srcIp	string
	handle 	*pcap.Handle
}

var srcIPs []string
var srcMACs []string
var dstIp, dstMac string
var wg sync.WaitGroup

func createWorkers(count int, jobs chan connDetails) {
	for index := 0; index < count; index++ {
		wg.Add(1)
		go packetGen(jobs)
	}
}

//Original - http://play.golang.org/p/m8TNTtygK0
//Modified - https://play.golang.org/p/IALp9qRmq-e
func incIP(ip net.IP, index *int) {
	(*index)++
	for j:= len(ip)-1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func incMAC(mac net.HardwareAddr, index *int) {
	(*index)++
	for j:= len(mac)-1; j >= 0; j-- {
		mac[j]++
		if mac[j] > 0 {
			break
		}
	}
}


func ipMacGenerator(count int, srcIp string, srcMac string) {
	startIP := net.ParseIP(srcIp)
	startMAC, _ := net.ParseMAC(srcMac)
	//fmt.Println(startIP, startMAC)

	for index:=0;index<count;incIP(startIP, &index) {
		srcIPs = append(srcIPs, startIP.String())
	}

	for index:=0; index<count;incMAC(startMAC, &index) {
		srcMACs = append(srcMACs, startMAC.String())
	}
}

func packetGen(connChannel chan connDetails) {

	defer wg.Done()
	for connDetail := range connChannel {

		//fmt.Println(connDetail)
		payLoad := gopacket.Payload([]byte{1,2,3,4})
		l:= uint16(len(payLoad))

		udp := layers.UDP {
			SrcPort: 12347,
			DstPort: 12034,
			Length:   l + 8,
        		Checksum: 0,
		}

		ip := layers.IPv4 {
			Version:    0x4,
        		IHL:        5,
        		Length:     20 + l,
        		TTL:        255,
        		Flags:      0x40,
        		FragOffset: 0,
        		Checksum:   0,
			SrcIP : net.ParseIP(connDetail.srcIp),
			DstIP : net.ParseIP(dstIp),
			Protocol: syscall.IPPROTO_UDP,
		}

		dstMAC, _ := net.ParseMAC(dstMac)
		srcMAC, _ := net.ParseMAC(connDetail.srcMac)

		eth := layers.Ethernet{
			SrcMAC:       srcMAC, 
			DstMAC:       dstMAC,
			EthernetType: layers.EthernetTypeIPv4,
		}
		
		buf := gopacket.NewSerializeBuffer()
		opts := gopacket.SerializeOptions{
			FixLengths:       true,
			ComputeChecksums: true,
		}

		//gopacket.SerializeLayers(buf, opts, &eth, &ip, &udp, &payLoad)
		payLoad.SerializeTo(buf, opts)
		udp.SerializeTo(buf, opts)
		ip.SerializeTo(buf, opts)
		eth.SerializeTo(buf, opts)
		//fmt.Println(buf.Bytes())


		err := connDetail.handle.WritePacketData(buf.Bytes())
		if err != nil {
			log.Fatal(err)
		}
	}
}


func main() {

	flag_intf := flag.String("i", "lo", "device name")
	flag_dstIp := flag.String("di", "127.0.0.1", "destination ip")
	flag_dstMac := flag.String("dm", "FF:FF:FF:FF:FF:FF", "destination MAC address")
	flag_srcMac := flag.String("sm", "00:01:02:00:00:00", "Src MAC address")
	flag_srcIp := flag.String("si", "20.0.0.1", "source ip")
	flag_count := flag.Int("c", 10, "Number of connections to be generated with auto-generated MAC and IP")
	flag_rate := flag.Int("r", 30, "Rate of packet generation needed in seconds")
	
	flag.Parse()

	ifName := *flag_intf
	dstIp = *flag_dstIp
	dstMac = *flag_dstMac
	srcIp := *flag_srcIp
	srcMac := *flag_srcMac
	count := *flag_count
	rate := *flag_rate
	//fmt.Println(ifName, dstIp, count)
	handle, err := pcap.OpenLive(ifName, 1024, false, 30*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	//Autogenerate the IP and MAC for packet generation
	ipMacGenerator(count, srcIp, srcMac)

	//fmt.Println(srcIPs, srcMACs)
	jobsChannel := make(chan connDetails, rate)

	//Spawn many goroutines based on the rate per second and give them the 
	createWorkers(rate, jobsChannel)

	//Divide the requests equally across the requested rate
	interval := time.Second / time.Duration(rate)

	//Start the timer with the interval
	timerChan := time.NewTicker(interval).C
	oneSecChan := time.NewTimer(time.Second).C

	currCount := 0
	sentRate := 0
	for {
		select {
		case <-timerChan:
			if (currCount >= count) {
				close(jobsChannel)
				wg.Wait()
				fmt.Println("Done generating packets.\nGood Bye!!!")
				return	
			}
			conn := connDetails{srcIp: srcIPs[currCount], srcMac: srcMACs[currCount], handle: handle}
			jobsChannel <- conn
			currCount++
			sentRate = sentRate + 1
		case <-oneSecChan:
			fmt.Println("CurrentRate is", sentRate)
			sentRate = 0
			oneSecChan = time.NewTimer(time.Second).C
		}
	}
}
