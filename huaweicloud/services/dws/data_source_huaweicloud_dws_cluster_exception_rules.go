package dws

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DWS GET /v1/{project_id}/clusters/{cluster_id}/workload/rules
func DataSourceClusterExceptionRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceClusterExceptionRulesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the exception rules are located.`,
			},

			// Required parameters.
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the cluster to which the exception rules belong.`,
			},

			// Optional parameters.
			"rule_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the exception rule to be queried, which is the fuzzy query.`,
			},

			// Attributes.
			"rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the exception rule.`,
						},
						"configurations": {
							Type:        schema.TypeMap,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The configuration items of the exception rule.`,
						},
					},
				},
				Description: `The list of the exception rules that matched filter parameters.`,
			},
		},
	}
}

func buildClusterExceptionRulesQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("rule_name"); ok {
		res = fmt.Sprintf("%s&rule_name=%v", res, v)
	}

	return res
}

func flattenClusterExceptionRules(exceptionRules []interface{}) []map[string]interface{} {
	if len(exceptionRules) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(exceptionRules))
	for _, rule := range exceptionRules {
		result = append(result, map[string]interface{}{
			"name":           utils.PathSearch("name", rule, nil),
			"configurations": utils.PathSearch("except_rules", rule, make(map[string]interface{})),
		})
	}

	return result
}

func dataSourceClusterExceptionRulesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		clusterId = d.Get("cluster_id").(string)
	)

	client, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	items, err := listClusterExceptionRules(client, clusterId, buildClusterExceptionRulesQueryParams(d))
	if err != nil {
		return diag.Errorf("error querying cluster exception rules: %s", err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("rules", flattenClusterExceptionRules(items)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
