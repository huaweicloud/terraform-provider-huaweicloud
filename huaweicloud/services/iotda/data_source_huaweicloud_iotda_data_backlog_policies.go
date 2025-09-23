package iotda

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IoTDA GET /v5/iot/{project_id}/routing-rule/backlog-policy
func DataSourceDataBacklogPolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDataBacklogPoliciesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"policy_name": {
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
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"backlog_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"backlog_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildDataBacklogPoliciesQueryParams(d *schema.ResourceData) string {
	if v, ok := d.GetOk("policy_name"); ok {
		return fmt.Sprintf("&policy_name=%v", v)
	}

	return ""
}

func dataSourceDataBacklogPoliciesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		httpUrl     = "v5/iot/{project_id}/routing-rule/backlog-policy?limit=50"
		product     = "iotda"
		allPolicies []interface{}
		offset      = 0
	)

	isDerived := WithDerivedAuth(cfg, region)
	client, err := cfg.NewServiceClientWithDerivedAuth(product, region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildDataBacklogPoliciesQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		requestPathWithOffset := fmt.Sprintf("%s&offset=%d", requestPath, offset)
		resp, err := client.Request("GET", requestPathWithOffset, &requestOpt)
		if err != nil {
			return diag.Errorf("error querying IoTDA data backlog policies: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		policies := utils.PathSearch("backlog_policies", respBody, make([]interface{}, 0)).([]interface{})
		if len(policies) == 0 {
			break
		}

		allPolicies = append(allPolicies, policies...)
		offset += len(policies)
	}

	uuID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(uuID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("policies", flattenDataBacklogPolicies(allPolicies)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDataBacklogPolicies(policies []interface{}) []interface{} {
	if len(policies) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(policies))
	for _, v := range policies {
		rst = append(rst, map[string]interface{}{
			"id":           utils.PathSearch("policy_id", v, nil),
			"name":         utils.PathSearch("policy_name", v, nil),
			"description":  utils.PathSearch("description", v, nil),
			"backlog_size": utils.PathSearch("backlog_size", v, nil),
			"backlog_time": utils.PathSearch("backlog_time", v, nil),
		})
	}

	return rst
}
