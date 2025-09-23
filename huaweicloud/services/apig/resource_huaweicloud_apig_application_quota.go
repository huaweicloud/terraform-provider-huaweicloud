package apig

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

// @API APIG POST /v2/{project_id}/apigw/instances/{instance_id}/app-quotas
// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/app-quotas/{app_quota_id}
// @API APIG PUT /v2/{project_id}/apigw/instances/{instance_id}/app-quotas/{app_quota_id}
// @API APIG DELETE /v2/{project_id}/apigw/instances/{instance_id}/app-quotas/{app_quota_id}
func ResourceApplicationQuota() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApplicationQuotaCreate,
		ReadContext:   resourceApplicationQuotaRead,
		UpdateContext: resourceApplicationQuotaUpdate,
		DeleteContext: resourceApplicationQuotaDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceApplicationQuotaImport,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Specifies the ID of the dedicated instance to which the application quota belongs.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the name of the application quota.",
			},
			"time_unit": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the limited time unit of the application quota.",
			},
			"call_limits": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Specifies the access limit of the application quota.",
			},
			"time_interval": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Specifies the limited time value for flow control of the application quota.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the description of the application quota.",
			},
			"bind_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of bound APPs.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the application quota, in RFC3339 format.`,
			},
		},
	}
}

func buildQuotaParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":          d.Get("name"),
		"time_unit":     d.Get("time_unit"),
		"call_limits":   d.Get("call_limits"),
		"time_interval": d.Get("time_interval"),
		"remark":        d.Get("description"),
	}
	return bodyParams
}

func resourceApplicationQuotaCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                = meta.(*config.Config)
		region             = cfg.GetRegion(d)
		createQuotaHttpUrl = "v2/{project_id}/apigw/instances/{instance_id}/app-quotas"
		createQuotaProduct = "apig"
	)

	quotaClient, err := cfg.NewServiceClient(createQuotaProduct, region)
	if err != nil {
		return diag.Errorf("error creating APIG Client: %s", err)
	}
	createQuotaPath := quotaClient.Endpoint + createQuotaHttpUrl
	createQuotaPath = strings.ReplaceAll(createQuotaPath, "{project_id}", quotaClient.ProjectID)
	createQuotaPath = strings.ReplaceAll(createQuotaPath, "{instance_id}", d.Get("instance_id").(string))
	createQuotaOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createQuotaOpt.JSONBody = utils.RemoveNil(buildQuotaParams(d))
	createQuotaResp, err := quotaClient.Request("POST", createQuotaPath, &createQuotaOpt)
	if err != nil {
		return diag.Errorf("error creating APIG application quota: %s", err)
	}

	createQuotaBody, err := utils.FlattenResponse(createQuotaResp)
	if err != nil {
		return diag.FromErr(err)
	}

	quotaId := utils.PathSearch("app_quota_id", createQuotaBody, "").(string)
	if quotaId == "" {
		return diag.Errorf("unable to find the APIG application quota ID from the API response")
	}
	d.SetId(quotaId)

	return resourceApplicationQuotaRead(ctx, d, meta)
}

func resourceApplicationQuotaRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	getApplicationQuotaHttpUrl := "v2/{project_id}/apigw/instances/{instance_id}/app-quotas/{app_quota_id}"
	getApplicationQuotaProduct := "apig"

	getApplicationQuotaClient, err := cfg.NewServiceClient(getApplicationQuotaProduct, region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	getApplicationQuotaPath := getApplicationQuotaClient.Endpoint + getApplicationQuotaHttpUrl
	getApplicationQuotaPath = strings.ReplaceAll(getApplicationQuotaPath, "{project_id}", getApplicationQuotaClient.ProjectID)
	getApplicationQuotaPath = strings.ReplaceAll(getApplicationQuotaPath, "{instance_id}", d.Get("instance_id").(string))
	getApplicationQuotaPath = strings.ReplaceAll(getApplicationQuotaPath, "{app_quota_id}", d.Id())

	getApplicationQuotaOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getApplicationQuotaResp, err := getApplicationQuotaClient.Request("GET", getApplicationQuotaPath, &getApplicationQuotaOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "APIG application quota")
	}

	respBody, err := utils.FlattenResponse(getApplicationQuotaResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", respBody, nil)),
		d.Set("time_unit", utils.PathSearch("time_unit", respBody, nil)),
		d.Set("call_limits", utils.PathSearch("call_limits", respBody, nil)),
		d.Set("time_interval", utils.PathSearch("time_interval", respBody, nil)),
		d.Set("description", utils.PathSearch("remark", respBody, nil)),
		d.Set("bind_num", utils.PathSearch("bound_app_num", respBody, nil)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("create_time",
			respBody, "").(string))/1000, false)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting APIG application quota fields: %s", err)
	}

	return nil
}

func resourceApplicationQuotaUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateQuotaHttpUrl := "v2/{project_id}/apigw/instances/{instance_id}/app-quotas/{app_quota_id}"
	updateQuotaProduct := "apig"

	updateQuotaClient, err := cfg.NewServiceClient(updateQuotaProduct, region)
	if err != nil {
		return diag.Errorf("error creating APIG Client: %s", err)
	}
	updateQuotaPath := updateQuotaClient.Endpoint + updateQuotaHttpUrl
	updateQuotaPath = strings.ReplaceAll(updateQuotaPath, "{project_id}", updateQuotaClient.ProjectID)
	updateQuotaPath = strings.ReplaceAll(updateQuotaPath, "{instance_id}", d.Get("instance_id").(string))
	updateQuotaPath = strings.ReplaceAll(updateQuotaPath, "{app_quota_id}", d.Id())

	updateQuotaOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildQuotaParams(d),
	}

	_, err = updateQuotaClient.Request("PUT", updateQuotaPath, &updateQuotaOpt)
	if err != nil {
		return diag.Errorf("error updating APIG application quota: %s", err)
	}
	return resourceApplicationQuotaRead(ctx, d, meta)
}

func resourceApplicationQuotaDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	quotaProduct := "apig"
	quotaClient, err := cfg.NewServiceClient(quotaProduct, region)
	if err != nil {
		return diag.Errorf("error creating APIG Client: %s", err)
	}
	delQuotaHttpUrl := "v2/{project_id}/apigw/instances/{instance_id}/app-quotas/{app_quota_id}"
	delQuotaPath := quotaClient.Endpoint + delQuotaHttpUrl
	delQuotaPath = strings.ReplaceAll(delQuotaPath, "{project_id}", quotaClient.ProjectID)
	delQuotaPath = strings.ReplaceAll(delQuotaPath, "{instance_id}", d.Get("instance_id").(string))
	delQuotaPath = strings.ReplaceAll(delQuotaPath, "{app_quota_id}", d.Id())

	delQuotaOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = quotaClient.Request("DELETE", delQuotaPath, &delQuotaOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting APIG application quota")
	}
	return nil
}

func resourceApplicationQuotaImport(_ context.Context, d *schema.ResourceData, _ interface{}) (
	[]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format of import ID, want to '<instance_id>/<id>', but got '%s'", d.Id())
	}

	d.SetId(parts[1])
	mErr := multierror.Append(nil,
		d.Set("instance_id", parts[0]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
