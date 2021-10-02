package dns

import (
	"fmt"

	"github.com/miekg/dns"
)

type Client interface {
	Query(hostname string, dnsType uint16) (*dns.Msg, error)
}

type dnsClient struct {
	client     *dns.Client
	serverAddr string
}

func NewClient(serverAddr string) Client {
	return &dnsClient{
		client:     new(dns.Client),
		serverAddr: serverAddr,
	}
}

func (c *dnsClient) Query(hostname string, dnsType uint16) (*dns.Msg, error) {
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(hostname), dnsType)
	r, _, err := c.client.Exchange(m, c.serverAddr)
	if err != nil {
		return nil, fmt.Errorf("could not get %s record - %w", dns.TypeToString[dnsType], err)
	}
	return r, nil
}

type fallbackClient struct {
	clients []Client
}

func NewFallbackClient(serverAddrs []string) Client {
	clients := make([]Client, 0)
	for _, s := range serverAddrs {
		clients = append(clients, NewClient(s))
	}
	return &fallbackClient{clients}
}

func (c *fallbackClient) Query(hostname string, dnsType uint16) (*dns.Msg, error) {
	var lastError error
	for _, cl := range c.clients {
		if r, err := cl.Query(hostname, dnsType); err == nil {
			return r, err
		} else {
			lastError = err
		}
	}
	return nil, lastError
}
