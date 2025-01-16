package live

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Live POST /v1/{project_id}/template/transcodings
// @API Live GET /v1/{project_id}/template/transcodings
// @API Live PUT /v1/{project_id}/template/transcodings
// @API Live DELETE /v1/{project_id}/template/transcodings
func ResourceTranscoding() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTranscodingCreate,
		ReadContext:   resourceTranscodingRead,
		UpdateContext: resourceTranscodingUpdate,
		DeleteContext: resourceTranscodingDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceTranscodingImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"domain_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"app_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"video_encoding": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"H264", "H265"}, false),
			},
			"templates": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 4,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"width": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"height": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"bitrate": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"frame_rate": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"i_frame_interval": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"gop": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"bitrate_adaptive": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"i_frame_policy": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			// This field is not returned by API, so the Computed attribute is not added.
			"trans_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"low_bitrate_hd": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func buildCreateOrUpdateTranscodingBodyParams(d *schema.ResourceData) map[string]interface{} {
	transcodingParams := map[string]interface{}{
		"domain":       d.Get("domain_name"),
		"app_name":     d.Get("app_name"),
		"trans_type":   utils.ValueIgnoreEmpty(d.Get("trans_type")),
		"quality_info": buildTemplatesBodyParams(d),
	}
	return transcodingParams
}

func buildTemplatesBodyParams(d *schema.ResourceData) []map[string]interface{} {
	rawParams := d.Get("templates").([]interface{})

	if len(rawParams) == 0 {
		return nil
	}

	templates := make([]map[string]interface{}, 0, len(rawParams))
	for _, v := range rawParams {
		raw := v.(map[string]interface{})
		params := map[string]interface{}{
			"templateName":     raw["name"],
			"quality":          "userdefine",
			"hdlb":             buildTemplatesHdlb(d),
			"codec":            d.Get("video_encoding"),
			"width":            raw["width"],
			"height":           raw["height"],
			"bitrate":          raw["bitrate"],
			"video_frame_rate": utils.ValueIgnoreEmpty(raw["frame_rate"]),
			"protocol":         utils.ValueIgnoreEmpty(raw["protocol"]),
			"iFrameInterval":   convertStrToInt(raw["i_frame_interval"].(string)),
			"gop":              convertStrToInt(raw["gop"].(string)),
			"bitrate_adaptive": utils.ValueIgnoreEmpty(raw["bitrate_adaptive"]),
			"i_frame_policy":   utils.ValueIgnoreEmpty(raw["i_frame_policy"]),
		}
		templates = append(templates, params)
	}

	return templates
}

func buildTemplatesHdlb(d *schema.ResourceData) string {
	var bitrateHdlb string
	if _, ok := d.GetOk("low_bitrate_hd"); ok {
		bitrateHdlb = "on"
	} else {
		bitrateHdlb = "off"
	}
	return bitrateHdlb
}

func resourceTranscodingCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		domainName = d.Get("domain_name").(string)
		appName    = d.Get("app_name").(string)
	)

	client, err := cfg.NewServiceClient("live", region)
	if err != nil {
		return diag.Errorf("error creating Live client: %s", err)
	}

	createHttpUrl := "v1/{project_id}/template/transcodings"
	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateOrUpdateTranscodingBodyParams(d)),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating Live transcoding: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/%s", domainName, appName))

	return resourceTranscodingRead(ctx, d, meta)
}

func resourceTranscodingRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		domainName = d.Get("domain_name").(string)
		appName    = d.Get("app_name").(string)
	)

	client, err := cfg.NewServiceClient("live", region)
	if err != nil {
		return diag.Errorf("error creating Live client: %s", err)
	}

	getResp, err := GetTranscodingTemplates(client, domainName, appName)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Live transcoding")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("domain_name", domainName),
		d.Set("app_name", utils.PathSearch("app_name", getResp, nil)),
		d.Set("templates", flattentranscodingTemplates(utils.PathSearch("quality_info", getResp, make([]interface{}, 0)))),
		d.Set("video_encoding", utils.PathSearch("quality_info[0].codec", getResp, nil)),
		d.Set("low_bitrate_hd", flattenTemplatesHdlb(utils.PathSearch("quality_info[0].hdlb", getResp, "").(string))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattentranscodingTemplates(resp interface{}) []map[string]interface{} {
	rawArray := resp.([]interface{})
	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, len(rawArray))
	for i, v := range rawArray {
		params := map[string]interface{}{
			"name":             utils.PathSearch("templateName", v, nil),
			"width":            utils.PathSearch("width", v, nil),
			"height":           utils.PathSearch("height", v, nil),
			"bitrate":          utils.PathSearch("bitrate", v, nil),
			"frame_rate":       utils.PathSearch("video_frame_rate", v, nil),
			"protocol":         utils.PathSearch("protocol", v, nil),
			"i_frame_interval": convertIntToStr(int(utils.PathSearch("iFrameInterval", v, float64(0)).(float64))),
			"gop":              convertIntToStr(int(utils.PathSearch("gop", v, float64(0)).(float64))),
			"bitrate_adaptive": utils.PathSearch("bitrate_adaptive", v, nil),
			"i_frame_policy":   utils.PathSearch("i_frame_policy", v, nil),
		}
		rst[i] = params
	}

	return rst
}

func flattenTemplatesHdlb(hdlb string) bool {
	return hdlb == "on"
}

func GetTranscodingTemplates(client *golangsdk.ServiceClient, domainName, appName string) (interface{}, error) {
	getHttpUrl := "v1/{project_id}/template/transcodings"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = fmt.Sprintf("%s?domain=%s&app_name=%s", getPath, domainName, appName)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, common.ConvertExpected400ErrInto404Err(err, "error_code", domainNameNotExistsCode)
	}

	getRespBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	template := utils.PathSearch("templates|[0]", getRespBody, nil)
	if template == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return template, nil
}

func resourceTranscodingUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("live", region)
	if err != nil {
		return diag.Errorf("error creating Live client: %s", err)
	}

	updateHttpUrl := "v1/{project_id}/template/transcodings"
	updatePath := client.Endpoint + updateHttpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200, 201, 204,
		},
		JSONBody: utils.RemoveNil(buildCreateOrUpdateTranscodingBodyParams(d)),
	}

	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating Live transcoding: %s", err)
	}

	return resourceTranscodingRead(ctx, d, meta)
}

func resourceTranscodingDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		domainName = d.Get("domain_name").(string)
		appName    = d.Get("app_name").(string)
	)

	client, err := cfg.NewServiceClient("live", region)
	if err != nil {
		return diag.Errorf("error creating Live client: %s", err)
	}

	deleteHttpUrl := "v1/{project_id}/template/transcodings"
	deletePath := client.Endpoint + deleteHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = fmt.Sprintf("%s?domain=%s&app_name=%s", deletePath, domainName, appName)

	deleteOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpts)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", domainNameNotExistsCode),
			"error deleting Live transcoding")
	}

	return nil
}

func resourceTranscodingImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<domain_name>/<app_name>', but got '%s'",
			importedId)
	}

	mErr := multierror.Append(nil,
		d.Set("domain_name", parts[0]),
		d.Set("app_name", parts[1]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}

func convertIntToStr(v interface{}) string {
	return fmt.Sprintf("%v", v)
}

func convertStrToInt(str string) interface{} {
	if str == "" {
		return nil
	}

	resp, err := strconv.Atoi(str)
	if err != nil {
		log.Printf("[WARN] Failed to convert the string (%s) to int", str)
		return 0
	}

	return resp
}
