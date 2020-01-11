package main

import (
	"errors"
	"github.com/armon/go-socks5"
	"github.com/mroth/weightedrand"
	"golang.org/x/net/context"
	"log"
	"math/rand"
	"net"
	"time"
)

const (
	ConfigFilename = "config.json"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	cfg, err := loadConfig(ConfigFilename)
	if err != nil {
		log.Fatalln("Unable to load config:", err)
	}

	chooser, err := getLocalIPAddrChooser(cfg.Interfaces)
	if err != nil {
		log.Fatalln("Unable to load interfaces:", err)
	}

	server, err := socks5.New(&socks5.Config{
		Dial: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
			lIPAddr := chooser.Pick().(*net.IPAddr)
			rAddr, err := net.ResolveTCPAddr("tcp", addr)
			if err != nil {
				return nil, err
			}
			return net.DialTCP(network, &net.TCPAddr{
				IP:   lIPAddr.IP,
				Zone: lIPAddr.Zone,
			}, rAddr)
		},
	})
	if err != nil {
		log.Fatalln("Unable to initialize SOCKS5 server:", err)
	}

	err = server.ListenAndServe("tcp", cfg.SOCKS5ListenAddr)
	log.Fatalln("Unable to listen and serve:", err)
}

func getLocalIPAddrChooser(cfgIfs []interfaceConfig) (*weightedrand.Chooser, error) {
	if len(cfgIfs) == 0 {
		return nil, errors.New("no interface in config")
	}

	localIfs, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	localIfsMap := make(map[string]*net.IPAddr, len(localIfs))
	for _, i := range localIfs {
		if addrs, err := i.Addrs(); err == nil && len(addrs) > 0 {
			// TODO: IPv6 support
			if addr := getIPv4IPAddr(addrs); addr != nil {
				localIfsMap[i.Name] = addr
			}
		}
	}
	log.Printf("Found %d available local interfaces.", len(localIfsMap))

	choices := make([]weightedrand.Choice, 0, len(cfgIfs))
	for _, i := range cfgIfs {
		if addr := localIfsMap[i.Name]; addr != nil {
			choices = append(choices, weightedrand.Choice{
				Item:   addr,
				Weight: i.Weight,
			})
			log.Printf("Using %s (%s), weight %d", i.Name, addr.String(), i.Weight)
		}
	}
	if len(choices) == 0 {
		return nil, errors.New("no matching interface found")
	}

	chooser := weightedrand.NewChooser(choices...)
	return &chooser, nil
}

func getIPv4IPAddr(addrs []net.Addr) *net.IPAddr {
	for _, addr := range addrs {
		switch addrT := addr.(type) {
		case *net.IPAddr:
			if ip4 := addrT.IP.To4(); ip4 != nil {
				return &net.IPAddr{
					IP:   ip4,
					Zone: addrT.Zone,
				}
			}
		case *net.IPNet:
			if ip4 := addrT.IP.To4(); ip4 != nil {
				return &net.IPAddr{
					IP: ip4,
				}
			}
		}
	}
	return nil
}
