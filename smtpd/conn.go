package smtpd

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

type Conn struct {
	net.Conn
	IsTLS  bool
	Errors []error
}

// write communicates back to the connected client
func (c *Conn) write(format string, vars ...interface{}) error {
	c.SetDeadline(time.Now().Add(10 * time.Second))
	_, err := c.Write([]byte(fmt.Sprintf(format, vars...) + "\r\n"))
	return err
}

// writeOK is a convenience function for sending the default OK response
func (c *Conn) writeOK() error {
	return c.write("250 OK")
}

// read handles brokering incoming SMTP protocol
func (c *Conn) read(max int) (string, error) {
	msg, err := c.rawRead("\r\n", max)
	log.Println(strings.TrimSpace(msg))
	return msg, err
}

// readData brokers the special case of SMTP data messages
func (c *Conn) readData(max int) (string, error) {
	return c.rawRead("\r\n.\r\n", max)
}

// rawRead performs the actual read from the connection
func (c *Conn) rawRead(suffix string, max int) (input string, err error) {
	var reply string
	reader := bufio.NewReader(c)
	for err == nil {
		c.SetDeadline(time.Now().Add(10 * time.Second))
		reply, err = reader.ReadString('\n')
		if reply != "" {
			input = input + reply
			if len(input) > max && max > 0 {
				return strings.Trim(input, " \n\r"), fmt.Errorf("Maximum DATA size exceeded (%v)", max)
			}
		}
		if err != nil {
			break
		}
		if strings.HasSuffix(input, suffix) {
			break
		}
	}
	return strings.Trim(input, " \n\r"), err
}
