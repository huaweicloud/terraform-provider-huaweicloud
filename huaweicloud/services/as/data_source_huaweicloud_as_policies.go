package as

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/autoscaling/v1/policies"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API AS GET /autoscaling-api/v1/{project_id}/scaling_policy/{scaling_group_id}/list
func DataSourceASPolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceASPoliciesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"scaling_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"scaling_policy_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scaling_policy_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scaling_policy_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scaling_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"alarm_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scheduled_policy": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"launch_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"recurrence_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"recurrence_value": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"start_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"end_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"action": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"operation": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"instance_number": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"instance_percentage": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"cool_down_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildDataSourcePolicyOpts(d *schema.ResourceData) policies.ListOpts {
	return policies.ListOpts{
		GroupID:  d.Get("scaling_group_id").(string),
		Name:     d.Get("scaling_policy_name").(string),
		Type:     d.Get("scaling_policy_type").(string),
		PolicyID: d.Get("scaling_policy_id").(string),
	}
}

func flattenDataSourcePolicies(policyResp []policies.Policy) []map[string]interface{} {
	policyList := make([]map[string]interface{}, 0, len(policyResp))
	for _, policy := range policyResp {
		policyMap := map[string]interface{}{
			"scaling_group_id": policy.ID,
			"id":               policy.PolicyID,
			"name":             policy.Name,
			"status":           policy.Status,
			"type":             policy.Type,
			"alarm_id":         policy.AlarmID,
			"scheduled_policy": flattenSchedulePolicy(policy.SchedulePolicy),
			"action":           flattenPolicyAction(policy.Action),
			"cool_down_time":   policy.CoolDownTime,
			"created_at":       policy.CreateTime,
		}

		policyList = append(policyList, policyMap)
	}
	return policyList
}

func dataSourceASPoliciesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		opts   = buildDataSourcePolicyOpts(d)
	)
	client, err := cfg.AutoscalingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating AS v1 client: %s", err)
	}

	policyResp, err := policies.List(client, opts).Extract()
	if err != nil {
		return diag.Errorf("error retrieving AS policies: %s", err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("policies", flattenDataSourcePolicies(policyResp)),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error saving AS policies data source fields: %s", mErr)
	}
	return nil
}
