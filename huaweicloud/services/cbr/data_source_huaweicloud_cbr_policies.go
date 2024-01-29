package cbr

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/cbr/v3/policies"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API CBR GET /v3/{project_id}/policies
func DataSourcePolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePoliciesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The region in which to query the CBR policies.",
			},
			"policy_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of CBR policies to query.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of CBR policies to query.",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The type of CBR policies to query.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to enable the CBR policy.",
			},
			"vault_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of vault to which the CBR policy resource belongs.",
			},
			"policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The policy ID.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The policy name.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The protection type of the CBR policy.",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to enable the CBR policy.",
						},
						"backup_cycle": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"interval": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The number of days between each backup.",
									},
									"days": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The weekly backup time.",
									},
									"execution_times": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The execution time of the policy.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
								},
							},
							Description: "The scheduling rule for the CBR policy backup execution.",
						},
						"destination_region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the replication destination region.",
						},
						"enable_acceleration": {
							Type:     schema.TypeBool,
							Computed: true,
							Description: "Whether to enable the acceleration function to shorten the replication time for " +
								"cross-region",
						},
						"destination_project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the replication destination project.",
						},
						"backup_quantity": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The maximum number of retained backups.",
						},
						"time_period": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The duration (in days) for retained backups.",
						},
						"long_term_retention": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"daily": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The latest backup of each day is saved in the long term.",
									},
									"weekly": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The latest backup of each week is saved in the long term.",
									},
									"monthly": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The latest backup of each month is saved in the long term.",
									},
									"yearly": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The latest backup of each year is saved in the long term.",
									},
									"full_backup_interval": {
										Type:     schema.TypeInt,
										Computed: true,
										Description: "How often (after how many incremental backups) a full backup is " +
											"performed.",
									},
								},
							},
							Description: "The long-term retention rules.",
						},
						"time_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The UTC time zone.",
						},
						"associated_vaults": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vault_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The vault ID of the associated CBR policy",
									},
									"destination_vault_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The destination vault ID associated with CBR policy",
									},
								},
							},
							Description: "The vault associated with the CBR policy",
						},
					},
				},
				Description: "All CBR policies that match the filter parameters.",
			},
		},
	}
}

func flattenAllAssociatedVaults(associatedVaults []policies.PolicyAssociateVault) []map[string]interface{} {
	if len(associatedVaults) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(associatedVaults))
	for i, associatedVault := range associatedVaults {
		result[i] = map[string]interface{}{
			"vault_id":             associatedVault.VaultID,
			"destination_vault_id": associatedVault.DestinationVaultID,
		}
	}

	return result
}

func flattenAllPolicies(d *schema.ResourceData, policyList []policies.Policy) ([]map[string]interface{}, error) {
	if len(policyList) < 1 {
		return nil, nil
	}

	rawID, idExist := d.GetOk("policy_id")
	rawName, nameExist := d.GetOk("name")
	rawEnabled, enabledExist := d.GetOk("enabled")

	result := make([]map[string]interface{}, 0)
	for _, policy := range policyList {
		if idExist && rawID.(string) != policy.ID {
			continue
		}
		if nameExist && rawName.(string) != policy.Name {
			continue
		}
		if enabledExist && rawEnabled.(bool) != policy.Enabled {
			continue
		}

		operationDefinition := policy.OperationDefinition
		policyMap := map[string]interface{}{
			"id":                     policy.ID,
			"name":                   policy.Name,
			"type":                   policy.OperationType,
			"enabled":                policy.Enabled,
			"destination_region":     operationDefinition.DestinationRegion,
			"destination_project_id": operationDefinition.DestinationProjectID,
			"associated_vaults":      flattenAllAssociatedVaults(policy.AssociatedVaults),
		}
		backupCycle, err := flattenPolicyBackupCycle(policy.Trigger.Properties.Pattern)
		if err != nil {
			return append(result, policyMap), err
		}

		policyMap["backup_cycle"] = backupCycle
		if operationDefinition.MaxBackups != -1 {
			policyMap["backup_quantity"] = operationDefinition.MaxBackups
			policyMap["long_term_retention"] = flattenLongTermRetention(operationDefinition)
			if operationDefinition.Timezone != "" {
				policyMap["time_zone"] = operationDefinition.Timezone
			}
		}

		if operationDefinition.RetentionDurationDays != -1 {
			policyMap["time_period"] = operationDefinition.RetentionDurationDays
		}

		result = append(result, policyMap)
	}

	return result, nil
}

func dataSourcePoliciesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.CbrV3Client(region)
	if err != nil {
		return diag.Errorf("error creating CBR v3 client: %s", err)
	}

	listOpts := policies.ListOpts{
		OperationType: d.Get("type").(string),
		VaultID:       d.Get("vault_id").(string),
	}
	allPages, err := policies.List(client, listOpts).AllPages()
	if err != nil {
		return diag.FromErr(err)
	}

	policyList, err := policies.ExtractPolicies(allPages)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "CBR policies")
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	result, err := flattenAllPolicies(d, policyList)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("policies", result),
	)
	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error saving policies data source fields: %s", mErr)
	}

	return nil
}
