package main

import (
	"bitbucket.org/gdamore/mangos"
	"bitbucket.org/gdamore/mangos/protocol/rep"
	"bitbucket.org/gdamore/mangos/transport/all"
	"fmt"
	"github.com/ugorji/go/codec"
	"reflect"
	"strconv"
	"time"
)

var (
	mh codec.MsgpackHandle
	b  []byte
)

func boomerangMetrics(map[string][]string) {}

func jsMetrics(map[string]string) {}

func delta(start string, end string) (error, int) {
	s, err := strconv.Atoi(start)
	if err != nil {
		return err, -1
	}
	e, err := strconv.Atoi(end)
	if err != nil {
		return err, -1
	}
	return nil, e - s
}

func decode(buf []byte) (error, map[string][]string) {

	doc := map[string][]string(nil)
	dec := codec.NewDecoderBytes(buf, &mh)
	err := dec.Decode(&doc)

	if err != nil {
		return err, nil
	}
	return nil, doc
}
func string2int(in map[string][]string) (error, map[string]int) {
	out := make(map[string]int)
	return nil, out
}

func main() {
	// consumer --type boomerang --remote tcp://127.0.0.1:8000 --statsd 192.168.33.20:8125

	url := "tcp://127.0.0.1:8000"

	responseServerReady := make(chan struct{})
	responseServer, err := rep.NewSocket()
	defer responseServer.Close()

	all.AddTransports(responseServer)
	if err != nil {
		fmt.Println("Error connecting: ", err)
		return
	}

	fmt.Println("Consumer ready")

	go func() {
		var err error
		var serverMsg *mangos.Message

		if err = responseServer.Listen(url); err != nil {
			fmt.Printf("\nServer listen failed: %v", err)
			return
		}

		close(responseServerReady)
		mh.MapType = reflect.TypeOf(map[string][]string(nil))

		for {
			if serverMsg, err = responseServer.RecvMsg(); err != nil {
				fmt.Printf("\nServer receive failed: %v", err)
			}

			err, d := decode(serverMsg.Body)
			fmt.Println("------ server msg ------ ")
			nt_dns, _ := delta(d["nt_dns_st"][0], d["nt_dns_end"][0])
			nt_con, _ := delta(d["nt_con_st"][0], d["nt_con_end"][0])
			nt_domcontloaded, _ := delta(d["nt_domcontloaded_st"][0], d["nt_domcontloaded_end"][0])
			//			roundtrip, _ := delta(d["rt_start"][0], d["rt_end"][0])
			fmt.Println("Navigation timing DNS: ", nt_dns)
			fmt.Println("Navigation timing Connection: ", nt_con)
			fmt.Println("Navigation timing DOM content loaded: ", nt_domcontloaded)
			//			fmt.Println("Roundtrip: ", roundtrip)
			//			for k, v := range d {
			//				fmt.Println(k, v)
			//			}
			fmt.Println("------ server msg ------ ")

			serverMsg.Body = []byte("OK")

			if err = responseServer.SendMsg(serverMsg); err != nil {
				fmt.Printf("\nServer send failed: %v", err)
				return
			}
		}
		fmt.Println("Listening")
	}()

	for {
		time.Sleep(10 * time.Second)
	}
}
