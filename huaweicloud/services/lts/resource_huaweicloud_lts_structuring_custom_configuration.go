// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product LTS
// ---------------------------------------------------------------

package lts

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const (
	regexParseType = "custom_regex"
	jsonParseType  = "json"
	splitParseType = "split"
	nginxParseType = "nginx"
)

// @API LTS GET /v2/{project_id}/lts/struct/template
// @API LTS POST /v2/{project_id}/lts/struct/template
// @API LTS PUT /v2/{project_id}/lts/struct/template
// @API LTS DELETE /v2/{project_id}/lts/struct/template
func ResourceStructCustomConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceStructCustomConfigCreate,
		UpdateContext: resourceStructCustomConfigUpdate,
		ReadContext:   resourceStructCustomConfigRead,
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
			"content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies a sample log event.`,
			},
			"demo_fields": {
				Type:        schema.TypeList,
				Elem:        structCustomConfigDemoFieldSchema(),
				Required:    true,
				Description: `Specifies the list of example fields.`,
			},
			"regex_rules": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the regular expression.`,
				ExactlyOneOf: []string{
					"layers", "tokenizer", "log_format",
				},
			},
			"layers": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `Specifies the maximum parsing layers.`,
			},
			"tokenizer": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the delimiter.`,
			},
			"log_format": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the nginx configuration.`,
			},
			"tag_fields": {
				Type:        schema.TypeList,
				Elem:        structCustomConfigTagFieldSchema(),
				Optional:    true,
				Description: `Specifies the tag field list.`,
			},
		},
	}
}

func structCustomConfigTagFieldSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"field_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the field name.`,
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the field data type.`,
			},
			"is_analysis": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether quick analysis is enabled.`,
			},
			"content": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the field content.`,
			},
		},
	}
	return &sc
}

func structCustomConfigDemoFieldSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"field_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the field name.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the field data type.`,
			},
			"is_analysis": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether quick analysis is enabled.`,
			},
			"content": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the content.`,
			},
		},
	}
	return &sc
}

func resourceStructCustomConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateOrUpdateStructCustomConfigBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating LTS structuring custom configuration: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, isString := createRespBody.(string)
	if !isString {
		return diag.Errorf("error creating LTS structuring custom configuration: the API response is not string")
	}

	if id == "" {
		return diag.Errorf("error creating LTS structuring custom configuration:" +
			" ID is not found in API response")
	}
	d.SetId(id)
	return resourceStructCustomConfigRead(ctx, d, meta)
}

func buildCreateOrUpdateStructCustomConfigBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"log_group_id":  d.Get("log_group_id"),
		"log_stream_id": d.Get("log_stream_id"),
		"content":       d.Get("content"),
		"parse_type":    buildStructCustomConfigParseType(d),
		"regex_rules":   utils.ValueIgnoreEmpty(d.Get("regex_rules")),
		"layers":        utils.ValueIgnoreEmpty(d.Get("layers")),
		"tokenizer":     utils.ValueIgnoreEmpty(d.Get("tokenizer")),
		"log_format":    utils.ValueIgnoreEmpty(d.Get("log_format")),
		"tag_fields":    buildStructCustomConfigFields(d.Get("tag_fields")),
		"demo_fields":   buildStructCustomConfigFields(d.Get("demo_fields")),
	}
	return bodyParams
}

func buildStructCustomConfigParseType(d *schema.ResourceData) string {
	if _, ok := d.GetOk("regex_rules"); ok {
		return regexParseType
	}
	if _, ok := d.GetOk("layers"); ok {
		return jsonParseType
	}
	if _, ok := d.GetOk("tokenizer"); ok {
		return splitParseType
	}
	return nginxParseType
}

func buildStructCustomConfigFields(rawParams interface{}) []map[string]interface{} {
	rawArray, ok := rawParams.([]interface{})
	if !ok || len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(rawArray))
	for _, v := range rawArray {
		if raw, isMap := v.(map[string]interface{}); isMap {
			rst = append(rst, map[string]interface{}{
				"fieldName":  utils.ValueIgnoreEmpty(raw["field_name"]),
				"isAnalysis": utils.ValueIgnoreEmpty(raw["is_analysis"]),
				"content":    utils.ValueIgnoreEmpty(raw["content"]),
				"type":       utils.ValueIgnoreEmpty(raw["type"]),
			})
		}
	}
	return rst
}

func resourceStructCustomConfigRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		return common.CheckDeletedDiag(d, err, "error retrieving LTS structuring custom configuration")
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
		d.Set("content", utils.PathSearch("demoLog", detailRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceStructCustomConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateOrUpdateStructCustomConfigBodyParams(d)),
	}

	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating LTS structuring custom configuration: %s", err)
	}
	return resourceStructCustomConfigRead(ctx, d, meta)
}
