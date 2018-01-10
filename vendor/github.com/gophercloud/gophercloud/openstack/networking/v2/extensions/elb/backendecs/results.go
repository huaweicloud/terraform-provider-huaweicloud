package backendecs

import (
	"fmt"

	"github.com/gophercloud/gophercloud"
)

type Listener struct {
	ID string `json:"id"`
}

type backend struct {
	ServerAddress string     `json:"server_address"`
	ID            string     `json:"id"`
	Address       string     `json:"address"`
	Status        string     `json:"status"`
	HealthStatus  string     `json:"health_status"`
	UpdateTime    string     `json:"update_time"`
	CreateTime    string     `json:"create_time"`
	ServerName    string     `json:"server_name"`
	ServerID      string     `json:"server_id"`
	Listeners     []Listener `json:"listeners"`
}

type Backend struct {
	PublicAddress  string     `json:"public_address"`
	ID             string     `json:"id"`
	PrivateAddress string     `json:"private_address"`
	Status         string     `json:"status"`
	HealthStatus   string     `json:"health_status"`
	UpdateTime     string     `json:"update_time"`
	CreateTime     string     `json:"create_time"`
	ServerName     string     `json:"server_name"`
	ServerID       string     `json:"server_id"`
	Listeners      []Listener `json:"listeners"`
}

// GetResult represents the result of a get operation.
type GetResult struct {
	gophercloud.Result
}

func (r GetResult) Extract() (*Backend, error) {
	var s []backend
	err := r.ExtractInto(&s)
	if err == nil {
		if len(s) != 1 {
			return nil, fmt.Errorf("get %d backends: %#v", len(s), s)
		}
		b := s[0]
		be := Backend{
			PublicAddress:  b.Address,
			ID:             b.ID,
			PrivateAddress: b.ServerAddress,
			Status:         b.Status,
			HealthStatus:   b.HealthStatus,
			UpdateTime:     b.UpdateTime,
			CreateTime:     b.CreateTime,
			ServerName:     b.ServerName,
			ServerID:       b.ServerID,
		}
		l := make([]Listener, len(b.Listeners))
		for i, v := range b.Listeners {
			l[i] = Listener{ID: v.ID}
		}
		be.Listeners = l

		return &be, nil
	}
	return nil, err
}
