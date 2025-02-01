package app

import (
	"errors"
	"flag"
	"strconv"
	"strings"
)

const (
	urlPrefix   = "http://"
	delimiter   = "/"
	contentType = "text/plain; charset=utf-8"
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
	a.normalize()
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
	a.Host = hp[0]
	a.Port = port
	a.normalize()
	return nil
}

func (a *ServHost) normalize() {
	if a.Port == 0 {
		a.Port = 8080
	}
	if a.Host == "" {
		a.Host = "localhost"
	}
}

type ShorLink struct {
	Addr string
}

func (a *ShorLink) String() string {
	return a.Addr
}

func (a *ShorLink) Set(s string) error {
	hp := strings.Split(s, ":")
	a.Addr = hp[0]
	return nil
}
