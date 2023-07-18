package nameservers

// Zone represents a DNS zone.
type NameServer struct {
	// Type of the name server.Value options:
	// public: indicates a public name server.
	// private: indicates a private name server.
	Type string `json:"type"`

	// Region ID. When you query a public name server, leave this parameter blank.
	Region string `json:"region"`

	// Array of name server record objects
	Records []Record `json:"ns_records"`
}

type Record struct {
	// Host name. This parameter is left blank when a private name server is used.
	HostName string `json:"hostname"`
	// Address of the name server. When the server is a public name server, this parameter is left blank.
	Address string `json:"address"`
	// the priority. If the value of priority is 1, the DNS server is the first one to resolve domain names.
	Priority int `json:"priority"`
}
