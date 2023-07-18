package nameservers

import (
	"github.com/chnsz/golangsdk"
)

type ListOpts struct {
	// Type of the name server. Value options:
	// public: indicates a public name server.
	// private: indicates a private name server.
	// It is left blank by default.
	Type string `q:"type"`
	// Region ID. When you query a public name server, leave this parameter blank.
	// Exact matching will work. It is left blank by default.
	Region string `q:"region"`
}

// ToNameServersQuery formats a ListOpts into a query string.
func (opts ListOpts) ToNameServersQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	return q.String(), err
}

// List implements a nameserver List request.
func List(client *golangsdk.ServiceClient, opts *ListOpts) ([]NameServer, error) {
	url := baseURL(client)
	if opts != nil {
		query, err := opts.ToNameServersQuery()
		if err != nil {
			return nil, err
		}
		url += query
	}

	var s struct {
		NameServers []NameServer `json:"nameservers"`
	}
	_, err := client.Get(url, &s, nil)
	return s.NameServers, err
}
