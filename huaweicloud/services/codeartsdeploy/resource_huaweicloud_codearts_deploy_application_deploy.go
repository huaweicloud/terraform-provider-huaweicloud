package codeartsdeploy

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CodeArtsDeploy POST /v2/tasks/{task_id}/start
// @API CodeArtsDeploy GET /v2/history/tasks/{task_id}/params
func ResourceDeployApplicationDeploy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDeployApplicationDeployCreate,
		ReadContext:   resourceDeployApplicationDeployRead,
		DeleteContext: resourceDeployApplicationDeployDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceDeployApplicationDeploymentRecordImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"task_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the deployment task ID.`,
			},
			"record_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the deployment record ID of an application.`,
			},
			"trigger_source": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the trigger source.`,
			},
			"params": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the parameters transferred during application deployment.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: `Specifies the parameter name transferred when deploying application.`,
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: `Specifies the parameter value transferred during application deployment.`,
						},
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: `Specifies the parameter type. If a dynamic parameter is set, the type is mandatory.`,
						},
					},
				},
			},
		},
	}
}

func resourceDeployApplicationDeployCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("codearts_deploy", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CodeArts deploy client: %s", err)
	}

	httpUrl := "v2/tasks/{task_id}/start"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{task_id}", d.Get("task_id").(string))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf-8",
		},
		JSONBody: utils.RemoveNil(buildDeployApplicationDeployBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error deploying CodeArts deploy application: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("id", createRespBody, "").(string)
	if id == "" {
		return diag.Errorf("unable to find the deployment record ID from the API response")
	}

	d.SetId(id)

	return resourceDeployApplicationDeployRead(ctx, d, meta)
}

func buildDeployApplicationDeployBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"record_id":      utils.ValueIgnoreEmpty(d.Get("record_id")),
		"trigger_source": utils.ValueIgnoreEmpty(d.Get("trigger_source")),
		"params":         buildDeployApplicationDeployBodyParamsParams(d),
	}
}

func buildDeployApplicationDeployBodyParamsParams(d *schema.ResourceData) []map[string]interface{} {
	rawArray := d.Get("params").([]interface{})
	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(rawArray))
	for _, v := range rawArray {
		if raw, isMap := v.(map[string]interface{}); isMap {
			rst = append(rst, map[string]interface{}{
				"key":   raw["name"],
				"value": raw["value"],
				"type":  raw["type"],
			})
		}
	}

	return rst
}

func resourceDeployApplicationDeployRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_deploy", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts deploy client: %s", err)
	}

	httpUrl := "v2/history/tasks/{task_id}/params?record_id={record_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{task_id}", d.Get("task_id").(string))
	getPath = strings.ReplaceAll(getPath, "{record_id}", d.Id())

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf-8",
		},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", "Deploy.00011303"), "")
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("params", flattenDeploymentRecordExecutionParameters(getRespBody, d)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDeploymentRecordExecutionParameters(respBody interface{}, d *schema.ResourceData) []interface{} {
	rawParams := d.Get("params").([]interface{})
	names := make(map[string]bool)
	for _, v := range rawParams {
		if raw, isMap := v.(map[string]interface{}); isMap {
			names[raw["name"].(string)] = true
		}
	}

	if resp, isList := respBody.([]interface{}); isList {
		rst := make([]interface{}, 0, len(names))
		for _, v := range resp {
			name := utils.PathSearch("name", v, "").(string)
			if !names[name] {
				continue
			}
			rst = append(rst, map[string]interface{}{
				"name":  name,
				"type":  utils.PathSearch("type", v, nil),
				"value": utils.PathSearch("value", v, nil),
			})
		}

		return rst
	}
	return nil
}

func resourceDeployApplicationDeployDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting application deploy resource is not supported. The resource is only removed from the state," +
		" the deployment record remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

func resourceDeployApplicationDeploymentRecordImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<task_id>/<record_id>', but got '%s'",
			d.Id())
	}

	if err := d.Set("task_id", parts[0]); err != nil {
		return nil, fmt.Errorf("error saving task ID: %s", err)
	}
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}
