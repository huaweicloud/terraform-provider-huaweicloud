package css

import (
	"context"
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

var aiOpsSettingNonUpdatableParams = []string{"cluster_id"}

// @API CSS POST /v1.0/{project_id}/clusters/{cluster_id}/ai-ops/setting
// @API CSS GET /v1.0/{project_id}/clusters/{cluster_id}/ai-ops/setting
// @API CSS PUT /v1.0/{project_id}/clusters/{cluster_id}/ai-ops/close
func ResourceAiOpsSetting() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAiOpsSettingCreate,
		ReadContext:   resourceAiOpsSettingRead,
		UpdateContext: resourceAiOpsSettingUpdate,
		DeleteContext: resourceAiOpsSettingDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(aiOpsSettingNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"check_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"period": {
				Type:     schema.TypeString,
				Required: true,
			},
			"check_items": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func buildAiOpsSettingBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"check_type":  d.Get("check_type"),
		"period":      d.Get("period"),
		"check_items": utils.ValueIgnoreEmpty(utils.ExpandToStringListBySet(d.Get("check_items").(*schema.Set))),
	}

	return bodyParams
}

func resourceAiOpsSettingCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		createHttpUrl = "v1.0/{project_id}/clusters/{cluster_id}/ai-ops/setting"
		clusterId     = d.Get("cluster_id").(string)
	)

	client, err := cfg.NewServiceClient("css", region)
	if err != nil {
		return diag.Errorf("error creating CSS client: %s", err)
	}

	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{cluster_id}", clusterId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildAiOpsSettingBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error opening the auto ai-ops setting: %s", err)
	}

	_, err = utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(clusterId)

	return resourceAiOpsSettingRead(ctx, d, meta)
}

func resourceAiOpsSettingRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("css", region)
	if err != nil {
		return diag.Errorf("error creating CSS client: %s", err)
	}

	checkSettings, err := GetAiOpsSettingInfo(client, d.Id())
	if err != nil {
		// When the cluster does not exist, the response HTTP status code of the query API is `403`
		return common.CheckDeletedDiag(d, common.ConvertExpected403ErrInto404Err(err, "errCode", "CSS.0015"),
			"error retrieving the auto ai-ops setting")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("cluster_id", d.Id()),
		d.Set("check_type", utils.PathSearch("check_type", checkSettings, nil)),
		d.Set("period", utils.PathSearch("period", checkSettings, nil)),
		d.Set("check_items", utils.PathSearch("check_items", checkSettings, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func GetAiOpsSettingInfo(client *golangsdk.ServiceClient, clusterId string) (interface{}, error) {
	httpUrl := "v1.0/{project_id}/clusters/{cluster_id}/ai-ops/setting"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{cluster_id}", clusterId)
	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	if respBody == nil || len(respBody.(map[string]interface{})) == 0 {
		// When the auto ai-ops setting is not opening, the response HTTP status code of the query API is `200`,
		// and return an empty object.
		return nil, golangsdk.ErrDefault404{}
	}

	return respBody, nil
}

func resourceAiOpsSettingUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1.0/{project_id}/clusters/{cluster_id}/ai-ops/setting"
	)

	client, err := cfg.NewServiceClient("css", region)
	if err != nil {
		return diag.Errorf("error creating CSS client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{cluster_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildAiOpsSettingBodyParams(d)),
	}

	_, err = client.Request("POST", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating the auto ai-ops setting: %s", err)
	}

	return resourceAiOpsSettingRead(ctx, d, meta)
}

func resourceAiOpsSettingDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1.0/{project_id}/clusters/{cluster_id}/ai-ops/close"
	)

	client, err := cfg.NewServiceClient("css", region)
	if err != nil {
		return diag.Errorf("error creating CSS client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{cluster_id}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("PUT", deletePath, &deleteOpt)
	if err != nil {
		// If the auto ai-ops setting is not opening, the response HTTP status code of the deletion API is `400`.
		// If the cluster does not exist, the response HTTP status code of the deletion API is `403`.
		return common.CheckDeletedDiag(
			d,
			common.ConvertExpected400ErrInto404Err(
				common.ConvertExpected403ErrInto404Err(err, "errCode", "CSS.0015"),
				"errCode",
				"CSS.0001",
			),
			"error closing the ai-ops setting",
		)
	}

	return nil
}
