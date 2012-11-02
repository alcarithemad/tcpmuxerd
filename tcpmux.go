package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"strings"
	"strconv"
)

type Service struct {
	Type	bool
	Port	int
	Path	string
}

func (srv *Service) Run(s *net.TCPConn) {
	if (srv.Type == true) {
		cmd := exec.Command(srv.Path)
		stdin, err := cmd.StdinPipe()
		if (err != nil) {
                        s.Write([]byte("- service not available at this time\n"))
                        s.Close()
                        return
                }
		stdout, err := cmd.StdoutPipe()
		if (err != nil) {
                        s.Write([]byte("- service not available at this time\n"))
                        s.Close()
                        return
                }
		err = cmd.Start()
		if (err != nil) {
                        s.Write([]byte("- service not available at this time\n"))
                        s.Close()
                        return
                }
		s.Write([]byte("+ connecting\n"))
		io.Copy(s, stdout)
		io.Copy(stdin, s)
	} else {
		d, err := net.DialTCP("tcp", nil, &net.TCPAddr{net.IPv4(127,0,0,1), srv.Port})
		if (err != nil) {
			s.Write([]byte("- service not available at this time\n"))
			s.Close()
			return
		}
		s.Write([]byte("+ connecting\n"))
		go io.Copy(s, d)
		go io.Copy(d, s)
	}
}

func readconf() *map[string]Service {
	cf_file, err := os.Open("tcpmux.conf")
	if (err != nil) { panic(err) }
	stat, err := cf_file.Stat()
	if (err != nil) { panic(err) }
	size := stat.Size()
	rawconf := make([]byte, size)
	cf_file.Read(rawconf)
	conf := make(map[string]Service, 1)
	lines := strings.Split(string(rawconf), "\n")
	for _, line := range lines {
		if (len(line) < 3) {continue}
		l := strings.Split(line, " ")
		srv_entry := Service{}
		if (l[1] == "!") {
			srv_entry.Type = true
			srv_entry.Path = l[2]
		} else {
			port, err := strconv.ParseInt(l[2], 10, 16)
			if (err != nil) { panic(err) } // TODO: print a better error
			srv_entry.Port = int(port)
		}
		conf[strings.ToUpper(l[0])] = srv_entry
	}
	return &conf
}

func process(s *net.TCPConn, conf *map[string]Service) {
	srv := make([]byte, 100)
	_, err := s.Read(srv)
	if (err != nil) {
			panic(err)
	}
	service := strings.Split(strings.ToUpper(string(srv)), "\n")[0]
	//fmt.Printf("'%s'", service)
	if (service == "HELP") {
		for name, _ := range *conf {
			s.Write([]byte(fmt.Sprintf("%s\n", name)))
		}
		s.Close()
		return
	}
	srv_entry, exists := (*conf)[string(service)]
	if (exists) {
		srv_entry.Run(s)
	} else { 
		s.Write([]byte("- service not found\n"))
		s.Close()
	}
}

func main() {
	fmt.Println("TCPMUX starting up")
	conf := readconf()
	sock, err := net.ListenTCP("tcp", &net.TCPAddr{net.IPv4(0,0,0,0), 1})
	if (err != nil) { panic(err) }
	defer func(){ sock.Close() }()
	for {
		conn, err := sock.AcceptTCP()
		if (err != nil) { continue }
		go process(conn, conf)
	}
}
