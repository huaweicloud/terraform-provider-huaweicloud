package notebook

import (
	"strconv"

	"github.com/chnsz/golangsdk/pagination"
)

const (
	StatusInit         = "INIT"
	StatusCreating     = "CREATING"
	StatusStarting     = "STARTING"
	StatusStopping     = "STOPPING"
	StatusDeleting     = "DELETING"
	StatusRunning      = "RUNNING"
	StatusStopped      = "STOPPED"
	StatusSnapshotting = "SNAPSHOTTING"
	StatusCreateFailed = "CREATE_FAILED"
	StatusStartFailed  = "START_FAILED"
	StatusDeleteFailed = "DELETE_FAILED"
	StatusError        = "ERROR"
	StatusDeleted      = "DELETED"
	StatusFrozen       = "FROZEN"
)

type Notebook struct {
	ActionProgress []JobProgress `json:"action_progress"`
	Description    string        `json:"description"`
	Endpoints      []Endpoints   `json:"endpoints"`
	FailReason     string        `json:"fail_reason"`
	Feature        string        `json:"feature"`
	Flavor         string        `json:"flavor"`
	Id             string        `json:"id"`
	Image          Image         `json:"image"`
	Lease          Lease         `json:"lease"`
	Name           string        `json:"name"`
	Pool           PoolRes       `json:"pool"`
	Status         string        `json:"status"`
	Token          string        `json:"token"`
	Url            string        `json:"url"`
	Volume         VolumeRes     `json:"volume"`
	WorkspaceId    string        `json:"workspace_id"`
	CreateAt       int           `json:"create_at"`
	UpdateAt       int           `json:"update_at"`
}

type JobProgress struct {
	NotebookId      string `json:"notebook_id"`
	Status          string `json:"status"`
	Step            int    `json:"step"`
	StepDescription string `json:"step_description"`
}

type Endpoints struct {
	AllowedAccessIps []string `json:"allowed_access_ips"`
	Service          string   `json:"service"`
	KeyPairNames     []string `json:"key_pair_names"`
	Uri              string   `json:"uri"`
}

type Image struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	SwrPath string `json:"swr_path"`
	Type    string `json:"type"`
}

type Lease struct {
	CreateAt int  `json:"create_at"`
	Duration int  `json:"duration"`
	Enable   bool `json:"enable"`
	UpdateAt int  `json:"update_at"`
}

type PoolRes struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type VolumeRes struct {
	Capacity  int    `json:"capacity"`
	Category  string `json:"category"`
	MountPath string `json:"mount_path"`
	Ownership string `json:"ownership"`
	Status    string `json:"status"`
	URI       string `json:"uri"`
	ID        string `json:"id"`
}

type ListNotebooks struct {
	Current int        `json:"current"`
	Data    []Notebook `json:"data"`
	Pages   int        `json:"pages"`
	Size    int        `json:"size"`
	Total   int        `json:"total"`
}

type ImagesResp struct {
	Current int           `json:"current"`
	Data    []ImageDetail `json:"data"`
	Pages   int           `json:"pages"`
	Size    int           `json:"size"`
	Total   int           `json:"total"`
}

type ImageDetail struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	SwrPath     string `json:"swr_path"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Arch        string `json:"arch"`
}

type ImagePage struct {
	pagination.OffsetPageBase
}

// IsEmpty checks whether a ImagePage struct is empty.
func (current ImagePage) IsEmpty() (bool, error) {
	resp, err := extractImagesResp(current)
	if err != nil {
		return false, err
	}
	return len(resp.Data) == 0, nil
}

func extractImagesResp(r pagination.Page) (ImagesResp, error) {
	var s ImagesResp
	err := (r.(ImagePage)).ExtractInto(&s)
	return s, err
}

func ExtractImages(r pagination.Page) ([]ImageDetail, error) {
	var s ImagesResp
	err := (r.(ImagePage)).ExtractInto(&s)
	return s.Data, err
}

// NextPageURL generates the URL for the page of results after this one.
func (current ImagePage) NextPageURL() (string, error) {
	next := current.NextOffset()
	if next == 0 {
		return "", nil
	}
	resp, _ := extractImagesResp(current)
	currentURL := current.URL
	q := currentURL.Query()

	if resp.Current+1 >= resp.Pages {
		// This page is the last page.
		return "", nil
	}

	q.Set("offset", strconv.Itoa(next))
	currentURL.RawQuery = q.Encode()
	return currentURL.String(), nil
}

type flavorResp struct {
	Flavors []Flavor `json:"flavors"`
}

type Flavor struct {
	Arch        string      `json:"arch"`
	Ascend      AscendInfo  `json:"ascend"`
	Billing     BillingInfo `json:"billing"`
	Category    string      `json:"category"`
	Description string      `json:"description"`
	Feature     string      `json:"feature"`
	Free        bool        `json:"free"`
	Gpu         GPUInfo     `json:"gpu"`
	Id          string      `json:"id"`
	Memory      int         `json:"memory"`
	Name        string      `json:"name"`
	SoldOut     bool        `json:"sold_out"`
	Storages    []string    `json:"storages"`
	Vcpus       int         `json:"vcpus"`
}

type AscendInfo struct {
	Npu       int    `json:"npu"`
	NpuMemory string `json:"npu_memory"`
	Type      string `json:"type"`
}

type BillingInfo struct {
	Code    string `json:"code"`
	UnitNum int    `json:"unit_num"`
}

type GPUInfo struct {
	Gpu       int    `json:"gpu"`
	GpuMemory string `json:"gpu_memory"`
	Type      string `json:"type"`
}

type ModelartsError struct {
	ErrorCode string `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
}

type MountStorage struct {
	Category  string `json:"category"`
	Id        string `json:"id"`
	MountPath string `json:"mount_path"`
	Status    string `json:"status"`
	Uri       string `json:"uri"`
}

type MountStorageListResp struct {
	Current int            `json:"current"`
	Data    []MountStorage `json:"data"`
	Pages   int            `json:"pages"`
	Size    int            `json:"size"`
	Total   int            `json:"total"`
}
