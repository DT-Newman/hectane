package util

import (
	"io"
	"net"
	"net/mail"
	"strings"
)

// Attempt to find the mail servers for the specified host. MX records are
// checked first. If one or more were found, the records are converted into an
// array of strings (sorted by priority). If none were found, the original host
// is returned.
func FindMailServers(host string) []string {
	if mx, err := net.LookupMX(host); err == nil {
		servers := make([]string, len(mx))
		for i, r := range mx {
			servers[i] = strings.TrimSuffix(r.Host, ".")
		}
		return servers
	} else {
		return []string{host}
	}
}

// Group a list of email addresses by their host. An error will be returned if
// any of the addresses are invalid.
func GroupAddressesByHost(addrs []string) (map[string][]string, error) {
	m := make(map[string][]string)
	for _, a := range addrs {
		if addr, err := mail.ParseAddress(a); err == nil {
			parts := strings.Split(addr.Address, "@")
			if m[parts[1]] == nil {
				m[parts[1]] = make([]string, 0, 1)
			}
			m[parts[1]] = append(m[parts[1]], addr.Address)
		} else {
			return nil, err
		}
	}
	return m, nil
}

// Given a reader for a MIME message, extract the address that the message is
// being sent from and the addresses that it is being delivered to.
func ExtractAddresses(r io.Reader) (string, []string, error) {
	if m, err := mail.ReadMessage(r); err == nil {
		addrs := make([]string, 0, 1)
		for _, h := range []string{"To", "Cc", "Bcc"} {
			if addrList, err := m.Header.AddressList(h); err == nil {
				for _, a := range addrList {
					addrs = append(addrs, a.Address)
				}
			}
		}
		return m.Header.Get("From"), addrs, nil
	} else {
		return "", nil, err
	}
}
