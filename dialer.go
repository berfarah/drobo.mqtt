package main

import (
	"bufio"
	"bytes"
	"net"
	"strings"
)

const (
	nasdXMLOpen  = "<ESATMUpdate>"
	nasdXMLClose = "</ESATMUpdate>"
)

// Credit to github.com/cosmouser/drobo_exporter
func getDroboInfo(target string) (drobo Drobo, err error) {
	conn, err := net.Dial("tcp", target)
	if err != nil {
		return drobo, err
	}
	defer conn.Close()

	var data bytes.Buffer
	var scanning bool
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		tmp := scanner.Text()

		if strings.Contains(tmp, nasdXMLOpen) && !scanning {
			scanning = true
		}

		if scanning {
			if n, err := data.WriteString(tmp); err != nil || n != len(tmp) {
				return drobo, err
			}
		}

		if strings.Contains(tmp, nasdXMLClose) && scanning {
			break
		}
	}

	return ReadXML(data.Bytes())
}
