package secmaster

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API SecMaster POST /v1/{project_id}/workspaces/{workspace_id}/soc/policys/search
func DataSourceSocPolicysSearch() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSocPolicysSearchRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"condition": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     policysConditionSchema(),
			},
			"sort": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     policysSortSchema(),
			},
			"group_by": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     policysGroupBySchema(),
			},
			"data": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func policysConditionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"conditions": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     policysConditionConditionsSchema(),
			},
			"logics": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func policysConditionConditionsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func policysSortSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"sort_by": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"order": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func policysGroupBySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"group_by_fields": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"group_by_hit": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     policysGroupByHitSchema(),
			},
		},
	}
}

func policysGroupByHitSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"source": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dest": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func buildSocPolicysSearchBodyParams(d *schema.ResourceData, offset int) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"offset":    utils.ValueIgnoreEmpty(offset),
		"condition": buildPolicysConditionBodyParams(d.Get("condition").([]interface{})),
		"sort":      buildPolicysSortBodyParams(d.Get("sort").([]interface{})),
		"group_by":  buildPolicysGroupByBodyParams(d.Get("group_by").([]interface{})),
	}

	return bodyParams
}

func buildPolicysConditionBodyParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"conditions": buildPolicysConditionConditionsBodyParams(rawMap["conditions"].([]interface{})),
		"logics":     utils.ValueIgnoreEmpty(rawMap["logics"]),
	}
}

func buildPolicysConditionConditionsBodyParams(rawArray []interface{}) []map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(rawArray))
	for _, v := range rawArray {
		rawMap, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		rst = append(rst, map[string]interface{}{
			"name": utils.ValueIgnoreEmpty(rawMap["name"]),
			"data": utils.ValueIgnoreEmpty(rawMap["data"]),
		})
	}

	return rst
}

func buildPolicysSortBodyParams(rawArray []interface{}) []map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(rawArray))
	for _, v := range rawArray {
		rawMap, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		rst = append(rst, map[string]interface{}{
			"sort_by": utils.ValueIgnoreEmpty(rawMap["sort_by"]),
			"order":   utils.ValueIgnoreEmpty(rawMap["order"]),
		})
	}

	return rst
}

func buildPolicysGroupByBodyParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"group_by_fields": utils.ValueIgnoreEmpty(rawMap["group_by_fields"]),
		"group_by_hit":    buildPolicysGroupByHitBodyParams(rawMap["group_by_hit"].([]interface{})),
	}
}

func buildPolicysGroupByHitBodyParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"source": utils.ValueIgnoreEmpty(rawMap["source"]),
		"dest":   utils.ValueIgnoreEmpty(rawMap["dest"]),
	}
}

func dataSourceSocPolicysSearchRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/policys/search"
		offset  = 0
		allData = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", d.Get("workspace_id").(string))

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		requestOpt.JSONBody = utils.RemoveNil(buildSocPolicysSearchBodyParams(d, offset))
		resp, err := client.Request("POST", requestPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving SOC policys: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		dataRaw := utils.PathSearch("data", respBody, make([]interface{}, 0)).([]interface{})
		if len(dataRaw) == 0 {
			break
		}

		allData = append(allData, dataRaw...)
		offset += len(dataRaw)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("data", utils.JsonToString(allData)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
