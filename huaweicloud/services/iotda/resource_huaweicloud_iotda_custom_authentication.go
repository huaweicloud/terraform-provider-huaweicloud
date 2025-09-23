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

// @API IoTDA POST /v5/iot/{project_id}/device-authorizers
// @API IoTDA GET /v5/iot/{project_id}/device-authorizers/{authorizer_id}
// @API IoTDA PUT /v5/iot/{project_id}/device-authorizers/{authorizer_id}
// @API IoTDA DELETE /v5/iot/{project_id}/device-authorizers/{authorizer_id}
func ResourceCustomAuthentication() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCustomAuthenticationCreate,
		ReadContext:   resourceCustomAuthenticationRead,
		UpdateContext: resourceCustomAuthenticationUpdate,
		DeleteContext: resourceCustomAuthenticationDelete,
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
			"authorizer_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"func_urn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"signing_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"signing_token": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				Computed:  true,
			},
			"signing_public_key": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"default_authorizer": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cache_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"func_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildCustomAuthenticationBodyParams(d *schema.ResourceData) map[string]interface{} {
	authenticationParams := map[string]interface{}{
		"authorizer_name":    d.Get("authorizer_name"),
		"func_urn":           d.Get("func_urn"),
		"signing_enable":     d.Get("signing_enable"),
		"signing_token":      utils.ValueIgnoreEmpty(d.Get("signing_token")),
		"signing_public_key": utils.ValueIgnoreEmpty(d.Get("signing_public_key")),
		"default_authorizer": d.Get("default_authorizer"),
		"status":             utils.ValueIgnoreEmpty(d.Get("status")),
		"cache_enable":       d.Get("cache_enable"),
	}

	return authenticationParams
}

func resourceCustomAuthenticationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		isDerived = WithDerivedAuth(cfg, region)
		httpUrl   = "v5/iot/{project_id}/device-authorizers"
	)

	client, err := cfg.NewServiceClientWithDerivedAuth("iotda", region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCustomAuthenticationBodyParams(d)),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating IoTDA custom authentication: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	authorizerId := utils.PathSearch("authorizer_id", respBody, "").(string)
	if authorizerId == "" {
		return diag.Errorf("error creating IoTDA custom authentication: ID is not found in API response")
	}

	d.SetId(authorizerId)

	return resourceCustomAuthenticationRead(ctx, d, meta)
}

func resourceCustomAuthenticationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		isDerived = WithDerivedAuth(cfg, region)
		httpUrl   = "v5/iot/{project_id}/device-authorizers/{authorizer_id}"
	)

	client, err := cfg.NewServiceClientWithDerivedAuth("iotda", region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{authorizer_id}", d.Id())
	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		// When the resource does not exist, query API will return `404`.
		return common.CheckDeletedDiag(d, err, "error retrieving IoTDA custom authentication")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("authorizer_name", utils.PathSearch("authorizer_name", getRespBody, nil)),
		d.Set("func_name", utils.PathSearch("func_name", getRespBody, nil)),
		d.Set("func_urn", utils.PathSearch("func_urn", getRespBody, nil)),
		d.Set("signing_enable", utils.PathSearch("signing_enable", getRespBody, nil)),
		d.Set("signing_token", utils.PathSearch("signing_token", getRespBody, nil)),
		d.Set("signing_public_key", utils.PathSearch("signing_public_key", getRespBody, nil)),
		d.Set("default_authorizer", utils.PathSearch("default_authorizer", getRespBody, nil)),
		d.Set("status", utils.PathSearch("status", getRespBody, nil)),
		d.Set("cache_enable", utils.PathSearch("cache_enable", getRespBody, nil)),
		d.Set("create_time", utils.PathSearch("create_time", getRespBody, nil)),
		d.Set("update_time", utils.PathSearch("update_time", getRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceCustomAuthenticationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		isDerived = WithDerivedAuth(cfg, region)
		httpUrl   = "v5/iot/{project_id}/device-authorizers/{authorizer_id}"
	)

	client, err := cfg.NewServiceClientWithDerivedAuth("iotda", region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{authorizer_id}", d.Id())
	updateOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCustomAuthenticationBodyParams(d)),
	}

	_, err = client.Request("PUT", updatePath, &updateOpts)
	if err != nil {
		return diag.Errorf("error updating IoTDA custom authentication: %s", err)
	}

	return resourceCustomAuthenticationRead(ctx, d, meta)
}

func resourceCustomAuthenticationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		isDerived = WithDerivedAuth(cfg, region)
		httpUrl   = "v5/iot/{project_id}/device-authorizers/{authorizer_id}"
	)

	client, err := cfg.NewServiceClientWithDerivedAuth("iotda", region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{authorizer_id}", d.Id())
	deleteOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpts)
	if err != nil {
		// When the resource does not exist, delete API will return `404`.
		return common.CheckDeletedDiag(d, err, "error deleting IoTDA custom authentication")
	}

	return nil
}
