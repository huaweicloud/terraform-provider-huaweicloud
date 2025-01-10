package live

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

// @API Live POST /v1/{project_id}/record/callbacks
// @API Live GET /v1/{project_id}/record/callbacks/{id}
// @API Live PUT /v1/{project_id}/record/callbacks/{id}
// @API Live DELETE /v1/{project_id}/record/callbacks/{id}
func ResourceRecordCallback() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRecordCallbackCreate,
		ReadContext:   resourceRecordCallbackRead,
		UpdateContext: resourceRecordCallbackUpdate,
		DeleteContext: resourceRecordCallbackDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
			"url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"types": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"sign_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"key": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				Computed:  true,
			},
		},
	}
}

func buildCreateOrUpdateCallbackBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"publish_domain":            d.Get("domain_name"),
		"app":                       "*",
		"notify_callback_url":       d.Get("url"),
		"notify_event_subscription": d.Get("types").([]interface{}),
		"sign_type":                 utils.ValueIgnoreEmpty(d.Get("sign_type")),
		"key":                       utils.ValueIgnoreEmpty(d.Get("key")),
	}
}

func resourceRecordCallbackCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/record/callbacks"
	)

	client, err := cfg.NewServiceClient("live", region)
	if err != nil {
		return diag.Errorf("error creating Live client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateOrUpdateCallbackBodyParams(d)),
	}

	requestResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating Live record callback configuration: %s", err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	ruleId := utils.PathSearch("id", respBody, "").(string)
	if ruleId == "" {
		return diag.Errorf("error creating Live record callback configuration: ID is not found in API response")
	}

	d.SetId(ruleId)

	return resourceRecordCallbackRead(ctx, d, meta)
}

func resourceRecordCallbackRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/record/callbacks/{id}"
	)

	client, err := cfg.NewServiceClient("live", region)
	if err != nil {
		return diag.Errorf("error creating Live client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{id}", d.Id())
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Live record callback configuration")
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("domain_name", utils.PathSearch("publish_domain", respBody, nil)),
		d.Set("types", utils.PathSearch("notify_event_subscription", respBody, make([]interface{}, 0))),
		d.Set("url", utils.PathSearch("notify_callback_url", respBody, nil)),
		d.Set("sign_type", utils.PathSearch("sign_type", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceRecordCallbackUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("live", region)
	if err != nil {
		return diag.Errorf("error creating Live client: %s", err)
	}

	if d.HasChanges("url", "types", "sign_type", "key") {
		httpUrl := "v1/{project_id}/record/callbacks/{id}"
		updatePath := client.Endpoint + httpUrl
		updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
		updatePath = strings.ReplaceAll(updatePath, "{id}", d.Id())
		updateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         utils.RemoveNil(buildCreateOrUpdateCallbackBodyParams(d)),
		}

		_, err := client.Request("PUT", updatePath, &updateOpt)
		if err != nil {
			return diag.Errorf("error updating Live record callback configuration: %s", err)
		}
	}

	return resourceRecordCallbackRead(ctx, d, meta)
}

func resourceRecordCallbackDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/record/callbacks/{id}"
	)

	client, err := cfg.NewServiceClient("live", region)
	if err != nil {
		return diag.Errorf("error creating Live client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{id}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting Live record callback configuration")
	}

	return nil
}
