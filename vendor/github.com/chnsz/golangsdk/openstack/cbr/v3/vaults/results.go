package vaults

import (
	"fmt"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/pagination"
)

type commonResult struct {
	golangsdk.Result
}

type CreateResult struct {
	commonResult
}

type GetResult struct {
	commonResult
}

type UpdateResult struct {
	commonResult
}

type DeleteResult struct {
	golangsdk.ErrResult
}

type Vault struct {
	ID                  string             `json:"id"`
	Name                string             `json:"name"`
	Billing             Billing            `json:"billing"`
	Description         string             `json:"description"`
	ProjectID           string             `json:"project_id"`
	ProviderID          string             `json:"provider_id"`
	Resources           []ResourceResp     `json:"resources"`
	Tags                []tags.ResourceTag `json:"tags"`
	EnterpriseProjectID string             `json:"enterprise_project_id"`
	AutoBind            bool               `json:"auto_bind"`
	BindRules           VaultBindRules     `json:"bind_rules"`
	UserID              string             `json:"user_id"`
	CreatedAt           string             `json:"created_at"`
	AutoExpand          bool               `json:"auto_expand"`
}

type Billing struct {
	Allocated       int    `json:"allocated"`
	ChargingMode    string `json:"charging_mode"`
	CloudType       string `json:"cloud_type"`
	ConsistentLevel string `json:"consistent_level"`
	ObjectType      string `json:"object_type"`
	OrderID         string `json:"order_id"`
	ProductID       string `json:"product_id"`
	ProtectType     string `json:"protect_type"`
	Size            int    `json:"size"`
	SpecCode        string `json:"spec_code"`
	Status          string `json:"status"`
	StorageUnit     string `json:"storage_unit"`
	Used            int    `json:"used"`
	FrozenScene     string `json:"frozen_scene"`
}

type ResourceResp struct {
	ExtraInfo     ResourceExtraInfo `json:"extra_info"`
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	ProtectStatus string            `json:"protect_status"`
	Size          int               `json:"size"`
	Type          string            `json:"type"`
	BackupSize    int               `json:"backup_size"`
	BackupCount   int               `json:"backup_count"`
}

func (r commonResult) Extract() (*Vault, error) {
	var s struct {
		Vault *Vault `json:"vault"`
	}
	err := r.ExtractInto(&s)
	return s.Vault, err
}

type OrderResp struct {
	ErrText string  `json:"errText"`
	ErrCode string  `json:"error_code"`
	RetCode int     `json:"retCode"`
	Orders  []Order `json:"orders"`
}

type Order struct {
	CloudServiceId     string   `json:"cloudServiceId"`
	ID                 string   `json:"orderId"`
	ReserveInstanceIds []string `json:"reserveInstanceIds"`
	ResourceId         string   `json:"resourceId"`
	SubscribeResult    string   `json:"subscribeResult"`
}

func (r CreateResult) ExtractOrder() (*OrderResp, error) {
	var s OrderResp
	err := r.ExtractInto(&s)
	return &s, err
}

type AssociateResourcesResult struct {
	golangsdk.Result
}

func (r AssociateResourcesResult) Extract() ([]string, error) {
	var s struct {
		AddResourceIDs []string `json:"add_resource_ids"`
	}
	if r.Err != nil {
		return nil, r.Err
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return nil, fmt.Errorf("failed to extract Associated Resource IDs")
	}
	return s.AddResourceIDs, nil
}

type DissociateResourcesResult struct {
	golangsdk.Result
}

func (r DissociateResourcesResult) Extract() ([]string, error) {
	var s struct {
		RemoveResourceIDs []string `json:"remove_resource_ids"`
	}
	if r.Err != nil {
		return nil, r.Err
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return nil, fmt.Errorf("failed to extract Dissociated Resource IDs")
	}
	return s.RemoveResourceIDs, nil
}

type BindPolicyResult struct {
	golangsdk.Result
}

type PolicyBinding struct {
	VaultID  string `json:"vault_id"`
	PolicyID string `json:"policy_id"`
}

func (r BindPolicyResult) Extract() (*PolicyBinding, error) {
	var s struct {
		PolicyBinding *PolicyBinding `json:"associate_policy"`
	}
	err := r.ExtractInto(&s)
	return s.PolicyBinding, err
}

type UnbindPolicyResult struct {
	golangsdk.Result
}

func (r UnbindPolicyResult) Extract() (*PolicyBinding, error) {
	var s struct {
		PolicyBinding *PolicyBinding `json:"dissociate_policy"`
	}
	err := r.ExtractInto(&s)
	return s.PolicyBinding, err
}

type VaultPage struct {
	pagination.SinglePageBase
}

func ExtractVaults(r pagination.Page) (*[]Vault, error) {
	var s struct {
		Vaults []Vault `json:"vaults"`
	}
	err := (r.(VaultPage)).ExtractInto(&s)
	return &s.Vaults, err
}
