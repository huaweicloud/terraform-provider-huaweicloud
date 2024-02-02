// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product COC
// ---------------------------------------------------------------

package coc

import (
	"context"
	"encoding/json"
	"log"
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

const defaultSensitiveValue = "*****************"

// @API COC POST /v1/job/scripts
// @API COC GET /v1/job/scripts/{script_uuid}
// @API COC PUT /v1/job/scripts/{script_uuid}
// @API COC DELETE /v1/job/scripts/{script_uuid}
func ResourceScript() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceScriptCreate,
		ReadContext:   resourceScriptRead,
		UpdateContext: resourceScriptUpdate,
		DeleteContext: resourceScriptDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

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
			"risk_level": {
				Type:     schema.TypeString,
				Required: true,
			},
			"version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"content": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: suppressDosOrUnix,
			},
			"parameters": {
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
							DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
								return oldValue == defaultSensitiveValue
							},
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
			"updated_at": {
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
		"name":          d.Get("name"),
		"description":   d.Get("description"),
		"type":          d.Get("type"),
		"content":       d.Get("content"),
		"script_params": buildScriptParamsBody(d.Get("parameters")),
		"properties": map[string]interface{}{
			"risk_level": d.Get("risk_level"),
			"version":    d.Get("version"),
		},
	}
	return bodyParams
}

func resourceScriptCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	return resourceScriptRead(ctx, d, meta)
}

func resourceScriptRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		// error_msg: invalid script uuid.
		if hasErrorCode(err, "COC.00040601") {
			err = golangsdk.ErrDefault404{}
		}

		return common.CheckDeletedDiag(d, err, "error retrieving COC script")
	}

	getScriptRespBody, err := utils.FlattenResponse(getScriptResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(nil,
		d.Set("name", utils.PathSearch("data.name", getScriptRespBody, nil)),
		d.Set("description", utils.PathSearch("data.description", getScriptRespBody, nil)),
		d.Set("type", utils.PathSearch("data.type", getScriptRespBody, nil)),
		d.Set("content", utils.PathSearch("data.content", getScriptRespBody, nil)),
		d.Set("status", utils.PathSearch("data.status", getScriptRespBody, nil)),
		d.Set("risk_level", utils.PathSearch("data.properties.risk_level", getScriptRespBody, nil)),
		d.Set("version", utils.PathSearch("data.properties.version", getScriptRespBody, nil)),
		d.Set("parameters", flattenScriptParams(getScriptRespBody)),
		d.Set("created_at", flattenScriptTimeStamp(getScriptRespBody, "data.gmt_created")),
		d.Set("updated_at", flattenScriptTimeStamp(getScriptRespBody, "data.gmt_modified")),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting COC script fields: %s", err)
	}

	return nil
}

func flattenScriptTimeStamp(resp interface{}, path string) interface{} {
	timeStamp := utils.PathSearch(path, resp, nil)
	if timeStamp == nil {
		return nil
	}

	return utils.FormatTimeStampRFC3339(int64(timeStamp.(float64))/1000, false)
}

func flattenScriptParams(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curJson := utils.PathSearch("data.script_params", resp, make([]interface{}, 0))
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

func buildUpdateScriptBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"description":   d.Get("description"),
		"type":          d.Get("type"),
		"content":       d.Get("content"),
		"script_params": buildScriptParamsBody(d.Get("parameters")),
		"properties": map[string]interface{}{
			"risk_level": d.Get("risk_level"),
			"version":    d.Get("version"),
		},
	}
	return bodyParams
}

func resourceScriptUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	updateScriptOpt.JSONBody = utils.RemoveNil(buildUpdateScriptBodyParams(d))
	_, err = client.Request("PUT", updateScriptPath, &updateScriptOpt)
	if err != nil {
		return diag.Errorf("error updating COC script: %s", err)
	}

	return resourceScriptRead(ctx, d, meta)
}

func resourceScriptDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

func suppressDosOrUnix(_, old, new string, _ *schema.ResourceData) bool {
	return strings.ReplaceAll(old, "\r\n", "\n") == strings.ReplaceAll(new, "\r\n", "\n")
}

func hasErrorCode(err error, expectCode string) bool {
	if errCode, ok := err.(golangsdk.ErrDefault500); ok {
		var response interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &response); jsonErr == nil {
			errorCode, parseErr := jmespath.Search("error_code", response)
			if parseErr != nil {
				log.Printf("[WARN] failed to parse error_code from response body: %s", parseErr)
			}

			if errorCode == expectCode {
				return true
			}
		}
	}

	return false
}
