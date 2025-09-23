package iotda

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

// @API IoTDA POST /v5/iot/{project_id}/apps
// @API IoTDA GET /v5/iot/{project_id}/apps/{app_id}
// @API IoTDA PUT /v5/iot/{project_id}/apps/{app_id}
// @API IoTDA DELETE /v5/iot/{project_id}/apps/{app_id}
func ResourceSpace() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSpaceCreate,
		ReadContext:   resourceSpaceRead,
		UpdateContext: resourceSpaceUpdate,
		DeleteContext: resourceSpaceDelete,
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
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"is_default": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceSpaceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		isDerived = WithDerivedAuth(cfg, region)
		httpUrl   = "v5/iot/{project_id}/apps"
	)

	client, err := cfg.NewServiceClientWithDerivedAuth("iotda", region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"app_name": d.Get("name"),
		},
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating IoTDA space: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	appId := utils.PathSearch("app_id", respBody, "").(string)
	if appId == "" {
		return diag.Errorf("error creating IoTDA space: ID is not found in API response")
	}

	d.SetId(appId)

	return resourceSpaceRead(ctx, d, meta)
}

func resourceSpaceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		isDerived = WithDerivedAuth(cfg, region)
		httpUrl   = "v5/iot/{project_id}/apps/{app_id}"
	)

	client, err := cfg.NewServiceClientWithDerivedAuth("iotda", region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{app_id}", d.Id())
	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		// When the resource does not exist, query API will return `403`, the error code is `IOTDA.001000`.
		return common.CheckDeletedDiag(d, common.ConvertExpected403ErrInto404Err(err, "error_code", "IOTDA.001000"),
			"error retrieving IoTDA space")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("app_name", getRespBody, nil)),
		d.Set("is_default", utils.PathSearch("default_app", getRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

// The **basic** edition instance not support update.
func resourceSpaceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		isDerived = WithDerivedAuth(cfg, region)
	)

	client, err := cfg.NewServiceClientWithDerivedAuth("iotda", region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	if d.HasChange("name") {
		httpUrl := "v5/iot/{project_id}/apps/{app_id}"
		updatePath := client.Endpoint + httpUrl
		updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
		updatePath = strings.ReplaceAll(updatePath, "{app_id}", d.Id())
		updateOpts := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody: map[string]interface{}{
				"app_name": d.Get("name"),
			},
		}

		_, err = client.Request("PUT", updatePath, &updateOpts)
		if err != nil {
			return diag.Errorf("error updating IoTDA space: %s", err)
		}
	}

	return resourceSpaceRead(ctx, d, meta)
}

func resourceSpaceDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		isDerived = WithDerivedAuth(cfg, region)
		httpUrl   = "v5/iot/{project_id}/apps/{app_id}"
	)

	client, err := cfg.NewServiceClientWithDerivedAuth("iotda", region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{app_id}", d.Id())
	deleteOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpts)
	if err != nil {
		// When the resource does not exist, delete API will return `403`, the error code is `IOTDA.001000`.
		return common.CheckDeletedDiag(d, common.ConvertExpected403ErrInto404Err(err, "error_code", "IOTDA.001000"),
			"error deleting IoTDA space")
	}

	return nil
}
