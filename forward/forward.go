package forward

import (
	"SrvCat/config"
	"SrvCat/storage"
	"github.com/fsnotify/fsnotify"
	"github.com/kataras/golog"
	"github.com/spf13/viper"
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

type forwardPort struct {
	Port   string
	Listen net.Listener
}

type forward struct {
	HandledPorts []*forwardPort
	Lock         sync.Mutex
}

var Forward = new(forward)

func init() {
	for _, portPair := range config.Config.Forwards {
		appendForward(portPair)
	}
	viper.OnConfigChange(func(e fsnotify.Event) {
		golog.Infof("Config file changed")
		Forward.RefreshPorts()
	})
}

func (f *forward) RefreshPorts() {
	f.Lock.Lock()
	var forwardPorts []string
	forwardPorts = append(forwardPorts, config.Config.Forwards...)
OuterLoop:
	for i := len(f.HandledPorts) - 1; i >= 0; {
		for j := len(forwardPorts) - 1; j >= 0; {
			if f.HandledPorts[i].Port == forwardPorts[j] {
				forwardPorts = append(forwardPorts[:j], forwardPorts[j+1:]...)
				continue OuterLoop
			} else {
				j--
			}
		}
		f.HandledPorts[i].Listen.Close()
		f.HandledPorts = append(f.HandledPorts[:i], f.HandledPorts[i+1:]...)
		i--
	}
	for _, portPair := range forwardPorts {
		appendForward(portPair)
	}
	f.Lock.Unlock()
}

func appendForward(portPair string) {
	ports := strings.Split(portPair, ":")
	fromAddr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:"+ports[0])
	if err != nil {
		golog.Errorf("Port pair %v resolve err: %v\n", portPair, err)
		return
	}
	toAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:"+ports[1])
	if err != nil {
		golog.Errorf("Port pair %v resolve err: %v\n", portPair, err)
		return
	}
	listener, err := net.ListenTCP("tcp", fromAddr)
	if err != nil {
		golog.Errorf("Port pair %v listen err: %v\n", portPair, err)
		return
	}
	go handleListener(listener, toAddr)
	port := &forwardPort{Port: portPair, Listen: listener}
	golog.Infof("Proxy listen on port pair: %v", portPair)
	Forward.HandledPorts = append(Forward.HandledPorts, port)
}

func handleListener(listener net.Listener, to *net.TCPAddr) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			golog.Errorf("Addr %v accept err: %v\n", listener.Addr().String(), err)
			continue
		}
		ip, _, err := net.SplitHostPort(conn.RemoteAddr().String())
		if err != nil {
			golog.Errorf("Remote %v parse err: %v\n", conn.RemoteAddr().String(), err)
			continue
		}
		log.Println("Connect request from ip: ", ip)
		verified, err := storage.Sqlite.GetVerified(ip, time.Now().Add(-time.Duration(config.Config.Settings.VerifyPeriod)*time.Minute).UnixMilli())
		if verified {
			golog.Infof("Accept connect from ip: ", ip)
			go handleConnect(conn, to)
			continue
		}
		conn.Close()
	}
}

func handleConnect(conn net.Conn, to *net.TCPAddr) {
	defer conn.Close()
	dial, err := net.DialTCP("tcp", nil, to)
	defer dial.Close()
	if err != nil {
		golog.Errorf("Port %v dial err: %v\n", to, err)
		return
	}
	go func() {
		defer conn.Close()
		defer dial.Close()
		io.Copy(conn, dial)
	}()
	io.Copy(dial, conn)
}
