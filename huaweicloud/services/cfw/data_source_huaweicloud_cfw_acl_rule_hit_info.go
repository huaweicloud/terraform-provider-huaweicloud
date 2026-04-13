package cfw

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CFW POST /v1/{project_id}/cfw/{fw_instance_id}/acl-rule/hit-info/batch-query
func DataSourceAclRuleHitInfo() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAclRuleHitInfoRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"fw_instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rule_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"records": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rule_hit_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"rule_last_hit_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildAclRuleHitInfoBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"rule_ids": utils.ExpandToStringList(d.Get("rule_ids").([]interface{})),
	}
}

func dataSourceAclRuleHitInfoRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "cfw"
		httpUrl = "v1/{project_id}/cfw/{fw_instance_id}/acl-rule/hit-info/batch-query"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{fw_instance_id}", d.Get("fw_instance_id").(string))
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildAclRuleHitInfoBodyParams(d),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving CFW ACL rule hit info: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("records", flattenAclRuleHitRecords(
			utils.PathSearch("data.records", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAclRuleHitRecords(recordsResp []interface{}) []interface{} {
	if len(recordsResp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(recordsResp))
	for _, record := range recordsResp {
		result = append(result, map[string]interface{}{
			"rule_id":            utils.PathSearch("rule_id", record, nil),
			"rule_hit_count":     utils.PathSearch("rule_hit_count", record, nil),
			"rule_last_hit_time": utils.PathSearch("rule_last_hit_time", record, nil),
		})
	}

	return result
}
