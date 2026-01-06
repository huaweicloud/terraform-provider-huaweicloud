package lts

import (
	"context"
	"encoding/json"
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

// @API LTS DELETE /v2/{project_id}/lts/struct/template
// @API LTS GET /v2/{project_id}/lts/struct/template
// @API LTS POST /v3/{project_id}/lts/struct/template
// @API LTS PUT /v3/{project_id}/lts/struct/template
func ResourceStructConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceStructConfigCreate,
		UpdateContext: resourceStructConfigUpdate,
		ReadContext:   resourceStructConfigRead,
		DeleteContext: resourceStructConfigDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceLtsStructConfigImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"log_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the log group ID.`,
			},
			"log_stream_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the log stream ID.`,
			},
			"template_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the type of the template.`,
			},
			"template_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the template name.`,
			},
			"template_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the template ID.`,
			},
			"demo_fields": {
				Type:        schema.TypeList,
				Elem:        structConfigFieldSchema(),
				Optional:    true,
				Description: `Specifies the example field array.`,
			},
			"tag_fields": {
				Type:        schema.TypeList,
				Elem:        structConfigFieldSchema(),
				Optional:    true,
				Description: `Specifies the tag field array.`,
			},
			"quick_analysis": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether to enable demo_fields and tag_fields quick analysis.`,
			},
			"demo_log": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The sample log event.`,
			},
		},
	}
}

func structConfigFieldSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"field_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the field name.`,
			},
			"is_analysis": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether quick analysis is enabled.`,
			},
		},
	}
	return &sc
}

func resourceStructConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/lts/struct/template"
		product = "lts"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateOrUpdateStructConfigBodyParams(d)),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating LTS structuring configuration: %s", err)
	}

	// Since the creation API response body is empty, the resource ID is obtained from the detailed API.
	detailRespBody, err := queryStructConfigDetail(client, d)
	if err != nil {
		return diag.Errorf("error creating LTS structuring configuration (failed to query detail API): %s", err)
	}

	templateId := utils.PathSearch("id", detailRespBody, "").(string)
	if templateId == "" {
		return diag.Errorf("unable to find the LTS structuring configuration ID from the API response")
	}
	d.SetId(templateId)

	return resourceStructConfigRead(ctx, d, meta)
}

func buildCreateOrUpdateStructConfigBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"log_group_id":   d.Get("log_group_id"),
		"log_stream_id":  d.Get("log_stream_id"),
		"template_type":  d.Get("template_type"),
		"template_id":    d.Get("template_id"),
		"template_name":  d.Get("template_name"),
		"quick_analysis": d.Get("quick_analysis"),
		"demo_fields":    buildStructConfigFields(d.Get("demo_fields")),
		"tag_fields":     buildStructConfigFields(d.Get("tag_fields")),
	}
	return bodyParams
}

func buildStructConfigFields(rawParams interface{}) []map[string]interface{} {
	rawArray, ok := rawParams.([]interface{})
	if !ok || len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(rawArray))
	for _, v := range rawArray {
		if raw, isMap := v.(map[string]interface{}); isMap {
			rst = append(rst, map[string]interface{}{
				"field_name":  raw["field_name"],
				"is_analysis": raw["is_analysis"],
			})
		}
	}
	return rst
}

func resourceStructConfigRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		mErr    *multierror.Error
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "lts"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	detailRespBody, err := queryStructConfigDetail(client, d)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving LTS structuring configuration")
	}

	// update the resource ID for import operation
	id := utils.PathSearch("id", detailRespBody, "")
	if id != "" {
		d.SetId(id.(string))
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("log_group_id", utils.PathSearch("logGroupId", detailRespBody, nil)),
		d.Set("log_stream_id", utils.PathSearch("logStreamId", detailRespBody, nil)),
		d.Set("template_name", utils.PathSearch("templateName", detailRespBody, nil)),
		d.Set("demo_log", utils.PathSearch("demoLog", detailRespBody, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

// queryStructConfigDetail use to query structuring configuration detail.
// The response of detail API is string value.
// If the response string value is empty, means the structuring configuration is not exist.
func queryStructConfigDetail(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	getPath := client.Endpoint + "v2/{project_id}/lts/struct/template"
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += buildGetStructConfigQueryParams(d)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	rawString, isString := getRespBody.(string)
	if !isString {
		return nil, err
	}

	if rawString == "" {
		// the structuring configuration is not exist
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v2/{project_id}/lts/struct/template",
				RequestId: "NONE",
				Body:      []byte(`the structuring configuration is not exist`),
			},
		}
	}

	var rst map[string]interface{}
	if err := json.Unmarshal([]byte(rawString), &rst); err != nil {
		return nil, err
	}
	return rst, nil
}

func buildGetStructConfigQueryParams(d *schema.ResourceData) string {
	logGroupId := d.Get("log_group_id").(string)
	logStreamId := d.Get("log_stream_id").(string)
	return fmt.Sprintf("?logGroupId=%s&logStreamId=%s", logGroupId, logStreamId)
}

func resourceStructConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/lts/struct/template"
		product = "lts"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateOrUpdateStructConfigBodyParams(d)),
	}

	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating LTS structuring configuration: %s", err)
	}
	return resourceStructConfigRead(ctx, d, meta)
}

func resourceStructConfigDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/lts/struct/template"
		product = "lts"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildDeleteStructConfigBodyParams(d),
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error deleting LTS structuring configuration: %s", err)
	}

	// The response to deletion API is always successful.
	// Determine whether the deletion is successful through the details API
	_, err = queryStructConfigDetail(client, d)
	if err == nil {
		return diag.Errorf("error deleting LTS structuring configuration: detailed information can still be found")
	}
	return common.CheckDeletedDiag(d, err, "error deleting LTS structuring configuration")
}

func buildDeleteStructConfigBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"id": d.Id(),
	}
	return bodyParams
}

func resourceLtsStructConfigImportState(_ context.Context, d *schema.ResourceData, _ interface{}) (
	[]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid resource ID format for LTS structuring configuration,"+
			" want '<log_group_id>/<log_stream_id>', but got '%s'", d.Id())
	}
	mErr := multierror.Append(nil,
		d.Set("log_group_id", parts[0]),
		d.Set("log_stream_id", parts[1]),
	)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
