package app

import (
	"errors"
	"flag"
	"strconv"
	"strings"
)

var CmdConfig cmdConfig

func (c *cmdConfig) Parse() error {
	sv := new(ServHost)
	_ = flag.Value(sv)
	sh := new(ShorLink)
	_ = flag.Value(sh)

	flag.Var(sv, "a", "Net address host:port")
	flag.Var(sh, "b", "short link server")
	flag.Parse()

	CmdConfig.ServHost = *sv
	CmdConfig.ShorLink = *sh
	return nil
}

type cmdConfig struct {
	ServHost ServHost
	ShorLink ShorLink
}
type ServHost struct {
	Host string
	Port int
}

func (a *ServHost) String() string {
	if a.Port == 0 {
		a.Port = 8080
	}
	if a.Host == "" {
		a.Host = "localhost"
	}
	return a.Host + ":" + strconv.Itoa(a.Port)
}

func (a *ServHost) Set(s string) error {
	hp := strings.Split(s, ":")
	if len(hp) != 2 {
		return errors.New("need address in a form host:port")
	}
	port, err := strconv.Atoi(hp[1])
	if err != nil {
		return err
	}

	if port == 0 {
		a.Port = 8080
	} else {
		a.Port = port
	}

	if hp[0] == "" {
		a.Host = "localhost"
	} else {
		a.Host = hp[0]
	}
	return nil
}

type ShorLink struct {
	Addr string
}

func (a *ShorLink) String() string {
	if a.Addr == "" {
		a.Addr = urlPrefix + delimiter
	}
	return a.Addr
}

func (a *ShorLink) Set(s string) error {
	hp := strings.Split(s, ":")
	if hp[0] == "" {
		a.Addr = urlPrefix + delimiter
	} else {
		a.Addr = hp[0]
	}
	return nil
}
