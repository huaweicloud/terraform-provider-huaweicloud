// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product COC
// ---------------------------------------------------------------

package ecs

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func ResourceComputeScript() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceComputeScriptCreate,
		ReadContext:   resourceComputeScriptRead,
		UpdateContext: resourceComputeScriptUpdate,
		DeleteContext: resourceComputeScriptDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Description: "schema: Internal",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"content": {
				Type:     schema.TypeString,
				Required: true,
			},
			"risk_level": {
				Type:     schema.TypeString,
				Required: true,
			},
			"version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"params": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
						"description": {
							Type:     schema.TypeString,
							Required: true,
						},
						"sensitive": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},

			// attributes
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildScriptParamsBody(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		params := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			params[i] = map[string]interface{}{
				"param_order":       i + 1, // param_order starts counting from 1
				"param_name":        raw["name"],
				"param_value":       raw["value"],
				"param_description": raw["description"],
				"sensitive":         raw["sensitive"],
			}
		}
		return params
	}

	return nil
}

func buildCreateScriptBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        d.Get("name"),
		"description": d.Get("description"),
		"type":        d.Get("type"),
		"content":     d.Get("content"),
		"params":      buildScriptParamsBody(d.Get("params")),
		"properties": map[string]interface{}{
			"risk_level": d.Get("risk_level"),
			"version":    d.Get("version"),
		},
	}
	return bodyParams
}

func resourceComputeScriptCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	product := "coc"

	client, err := cfg.NewServiceClient(product, cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	createScriptHttpUrl := "v1/job/scripts"
	createScriptPath := client.Endpoint + createScriptHttpUrl

	createScriptOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	createScriptOpt.JSONBody = utils.RemoveNil(buildCreateScriptBodyParams(d))
	createScriptResp, err := client.Request("POST", createScriptPath, &createScriptOpt)
	if err != nil {
		return diag.Errorf("error creating COC script: %s", err)
	}
	createScriptRespBody, err := utils.FlattenResponse(createScriptResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := jmespath.Search("data", createScriptRespBody)
	if err != nil {
		return diag.Errorf("error creating COC script: ID is not found in API response")
	}

	d.SetId(id.(string))
	return resourceComputeScriptRead(ctx, d, meta)
}

func resourceComputeScriptRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	product := "coc"

	client, err := cfg.NewServiceClient(product, cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	getScriptHttpUrl := "v1/job/scripts/{id}"
	getScriptPath := client.Endpoint + getScriptHttpUrl
	getScriptPath = strings.ReplaceAll(getScriptPath, "{id}", d.Id())

	getScriptOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getScriptResp, err := client.Request("GET", getScriptPath, &getScriptOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving COC script")
	}

	getScriptRespBody, err := utils.FlattenResponse(getScriptResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(nil,
		d.Set("name", utils.PathSearch("name", getScriptRespBody, nil)),
		d.Set("description", utils.PathSearch("description", getScriptRespBody, nil)),
		d.Set("type", utils.PathSearch("type", getScriptRespBody, nil)),
		d.Set("content", utils.PathSearch("content", getScriptRespBody, nil)),
		d.Set("status", utils.PathSearch("status", getScriptRespBody, nil)),
		d.Set("risk_level", utils.PathSearch("propertites.risk_level", getScriptRespBody, nil)),
		d.Set("version", utils.PathSearch("propertites.version", getScriptRespBody, nil)),
		d.Set("params", flattenScriptParams(getScriptRespBody)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting COC script fields: %s", err)
	}

	return nil
}

func flattenScriptParams(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curJson := utils.PathSearch("script_params", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, len(curArray))
	for i, v := range curArray {
		rst[i] = map[string]interface{}{
			"name":        utils.PathSearch("param_name", v, nil),
			"value":       utils.PathSearch("param_value", v, nil),
			"description": utils.PathSearch("param_description", v, nil),
			"sensitive":   utils.PathSearch("sensitive", v, false),
		}
	}
	return rst
}

func buildUpdateScriptBodyParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"description": d.Get("description"),
		"type":        d.Get("type"),
		"content":     d.Get("content"),
		"params":      buildScriptParamsBody(d.Get("params")),
		"properties": map[string]interface{}{
			"risk_level": d.Get("risk_level"),
			"version":    d.Get("version"),
		},
	}
	return bodyParams
}

func resourceComputeScriptUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	product := "coc"

	client, err := cfg.NewServiceClient(product, cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	updateScriptHttpUrl := "v1/job/scripts/{id}"
	updateScriptPath := client.Endpoint + updateScriptHttpUrl
	updateScriptPath = strings.ReplaceAll(updateScriptPath, "{id}", d.Id())

	updateScriptOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	updateScriptOpt.JSONBody = utils.RemoveNil(buildUpdateScriptBodyParams(d, cfg))
	_, err = client.Request("PUT", updateScriptPath, &updateScriptOpt)
	if err != nil {
		return diag.Errorf("error updating COC script: %s", err)
	}

	return resourceComputeScriptRead(ctx, d, meta)
}

func resourceComputeScriptDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	product := "coc"

	client, err := cfg.NewServiceClient(product, cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	deleteScriptHttpUrl := "v1/job/scripts/{id}"
	deleteScriptPath := client.Endpoint + deleteScriptHttpUrl
	deleteScriptPath = strings.ReplaceAll(deleteScriptPath, "{id}", d.Id())

	deleteScriptOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	_, err = client.Request("DELETE", deleteScriptPath, &deleteScriptOpt)
	if err != nil {
		return diag.Errorf("error deleting COC script: %s", err)
	}

	return nil
}
