package cbr

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

// @API CBR POST /v3/{project_id}/agent/check
func DataSourceAgentChecks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAgentChecksRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"agent_status": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"resource_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"resource_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"agent_status_attr": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"installed": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_old": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"message": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"code": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildAgentChecksBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"agent_status": buildResourceInfosBodyParams(d.Get("agent_status").([]interface{})),
	}

	return bodyParams
}

func buildResourceInfosBodyParams(resourceInfos []interface{}) []map[string]interface{} {
	if len(resourceInfos) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(resourceInfos))
	for _, v := range resourceInfos {
		resourceInfo, ok := v.(map[string]interface{})
		if !ok {
			return nil
		}

		params := map[string]interface{}{
			"resource_id":   resourceInfo["resource_id"],
			"resource_type": resourceInfo["resource_type"],
			"resource_name": utils.ValueIgnoreEmpty(resourceInfo["resource_name"]),
		}
		result = append(result, params)
	}

	return result
}

func dataSourceAgentChecksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/agent/check"
	)

	client, err := cfg.NewServiceClient("cbr", region)
	if err != nil {
		return diag.Errorf("error creating CBR client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildAgentChecksBodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving agent status: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	agentStatus := utils.PathSearch("agent_status", respBody, make([]interface{}, 0)).([]interface{})

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("agent_status_attr", flattenAgentChecks(agentStatus)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAgentChecks(agentData []interface{}) []interface{} {
	if len(agentData) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(agentData))
	for _, v := range agentData {
		rst = append(rst, map[string]interface{}{
			"resource_id": utils.PathSearch("resource_id", v, nil),
			"version":     utils.PathSearch("version", v, nil),
			"installed":   utils.PathSearch("installed", v, nil),
			"is_old":      utils.PathSearch("is_old", v, nil),
			"message":     utils.PathSearch("message", v, nil),
			"code":        utils.PathSearch("code", v, nil),
		})
	}

	return rst
}
