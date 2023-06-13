package cbr

import (
	"context"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/cbr/v3/policies"
	"github.com/chnsz/golangsdk/openstack/cbr/v3/vaults"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func DataSourceCbrVaultsV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCbrVaultsV3Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The region in which to query the vaults.",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
				Description:  "The vault name.",
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				//If the validation content has changed, please update the resource type map.
				ValidateFunc: validation.StringInSlice([]string{
					VaultTypeServer, VaultTypeDisk, VaultTypeTurbo,
				}, false),
				Description: "The object type of the vault.",
			},
			"consistent_level": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"crash_consistent", "app_consistent",
				}, false),
				Description: "The consistent level (specification) of the vault.",
			},
			"protection_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"backup", "replication",
				}, false),
				Description: "The protection type of the vault.",
			},
			"size": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 10485760),
				Description:  "The vault sapacity, in GB.",
			},
			"auto_expand_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to enable automatic expansion of the backup protection type vault.",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the enterprise project to which the vault belongs.",
			},
			"policy_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the policy associated with the vault.",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The vault status.",
			},
			"vaults": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The vault ID in UUID format.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The vault name.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The object type of the vault.",
						},
						"consistent_level": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The consistent level (specification) of the vault.",
						},
						"protection_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The protection type of the vault.",
						},
						"size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The vault capacity, in GB.",
						},
						"auto_expand_enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to enable automatic expansion of the backup protection type vault.",
						},
						"enterprise_project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The enterprise project ID.",
						},
						"policy_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the policy associated with the vault.",
						},
						"allocated": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "The allocated capacity of the vault, in GB.",
						},
						"used": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "The used capacity, in GB.",
						},
						"spec_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The specification code.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The vault status.",
						},
						"storage": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the bucket for the vault.",
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "The key/value pairs to associate with the vault.",
						},
						"resources": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"server_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The ID of the ECS instance to be backed up.",
									},
									"excludes": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "The array of disk IDs which will be excluded in the backup.",
									},
									"includes": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "The array of disk or SFS file system IDs which will be included in the backup.",
									},
								},
							},
							Description: "The array of one or more resources to attach to the vault.",
						},
					},
				},
			},
		},
	}
}

func buildVaultListOpts(d *schema.ResourceData, config *config.Config) vaults.ListOpts {
	return vaults.ListOpts{
		Limit:               100,
		CloudType:           "public",
		Name:                d.Get("name").(string),
		ObjectType:          d.Get("type").(string),
		ProtectType:         d.Get("protection_type").(string),
		PolicyID:            d.Get("policy_id").(string),
		EnterpriseProjectID: common.GetEnterpriseProjectID(d, config),
		Status:              d.Get("status").(string),
	}
}

func filterCbrVaults(d *schema.ResourceData, vaultList []vaults.Vault) ([]interface{}, error) {
	return utils.FilterSliceWithField(vaultList, map[string]interface{}{
		"Billing.ConsistentLevel": d.Get("consistent_level").(string),
		"Billing.Size":            d.Get("size").(int),
		"AutoExpand":              d.Get("auto_expand_enabled").(bool),
	})
}

func getCbrPolicyOfSpecificVault(client *golangsdk.ServiceClient, vaultId string) (*policies.Policy, error) {
	opt := policies.ListOpts{
		VaultID: vaultId,
	}
	pages, err := policies.List(client, opt).AllPages()
	if err != nil {
		return nil, err
	}
	resp, err := policies.ExtractPolicies(pages)
	if err != nil {
		return nil, err
	}

	if len(resp) > 0 {
		return &resp[0], nil
	}
	return nil, fmtp.Errorf("No policies are bound to the vault.")
}

func setCbrAllVaultParameters(client *golangsdk.ServiceClient, d *schema.ResourceData,
	vaultList []interface{}) error {
	result := make([]map[string]interface{}, len(vaultList))
	ids := make([]string, len(vaultList))
	for i, val := range vaultList {
		vault := val.(vaults.Vault)
		vMap := map[string]interface{}{
			"id":                    vault.ID,
			"name":                  vault.Name,
			"enterprise_project_id": vault.EnterpriseProjectID,
			"type":                  vault.Billing.ObjectType,
			"protection_type":       vault.Billing.ProtectType,
			"status":                vault.Billing.Status,
			"consistent_level":      vault.Billing.ConsistentLevel,
			"size":                  vault.Billing.Size,
			"allocated":             vault.Billing.Allocated,
			"used":                  vault.Billing.Used,
			"spec_code":             vault.Billing.SpecCode,
			"storage":               vault.Billing.StorageUnit,
			"auto_expand_enabled":   vault.AutoExpand,
			"tags":                  utils.TagsToMap(vault.Tags),
			"resources":             makeVaultResources(vault.Billing.ObjectType, vault.Resources),
		}

		// Query the CBR policy which bound to the vault by ID
		if policy, err := getCbrPolicyOfSpecificVault(client, vault.ID); err != nil {
			logp.Printf("[DEBUG] Unable to find the policy for specific vault: %s", err)
		} else {
			vMap["policy_id"] = policy.ID
		}
		result[i] = vMap
		ids[i] = vault.ID
	}

	d.SetId(hashcode.Strings(ids))

	return d.Set("vaults", result)
}

func dataSourceCbrVaultsV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.CbrV3Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating Huaweicloud CBR v3 client: %s", err)
	}

	opt := buildVaultListOpts(d, config)

	pages, err := vaults.List(client, opt).AllPages()
	if err != nil {
		return fmtp.DiagErrorf("Error getting vault details: %s", err)
	}

	resp, err := vaults.ExtractVaults(pages)
	if err != nil {
		return fmtp.DiagErrorf("error getting vault details: %s", err)
	}
	// Use the following parameters to filter the result of the List method return: consistent_level, size and
	// auto_expand_enabled.
	vaultList, err := filterCbrVaults(d, *resp)
	if err != nil {
		return fmtp.DiagErrorf("Error filting vaults by consistent_level, size and auto_expand_enabled: %s", err)
	}

	// Set the ID and other parameters.
	err = setCbrAllVaultParameters(client, d, vaultList)
	if err != nil {
		return fmtp.DiagErrorf("Error setting vaults parameter: %s", err)
	}
	return nil
}
