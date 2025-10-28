package coc

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var cloudVendorAccountNonUpdatableParams = []string{"vendor", "account_id"}

// @API COC POST /v1/vendor-account
// @API COC GET /v1/vendor-account
// @API COC PUT /v1/vendor-account/{id}
// @API COC DELETE /v1/vendor-account/{id}
func ResourceCloudVendorAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCloudVendorAccountCreate,
		ReadContext:   resourceCloudVendorAccountRead,
		UpdateContext: resourceCloudVendorAccountUpdate,
		DeleteContext: resourceCloudVendorAccountDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(cloudVendorAccountNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"vendor": {
				Type:     schema.TypeString,
				Required: true,
			},
			"account_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"account_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ak": {
				Type:     schema.TypeString,
				Required: true,
			},
			"sk": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sync_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"failure_msg": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sync_date": {
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

func resourceCloudVendorAccountCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v1/vendor-account"
		product = "coc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildCreateCloudVendorAccountBodyParams(d),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating COC cloud vendor account: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.Errorf("error flattening cloud vendor account: %s", err)
	}

	id := utils.PathSearch("data", createRespBody, "").(string)
	if id == "" {
		return diag.Errorf("unable to find the COC cloud vendor account ID from the API response")
	}

	d.SetId(id)

	return resourceCloudVendorAccountRead(ctx, d, meta)
}

func buildCreateCloudVendorAccountBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"vendor":       d.Get("vendor"),
		"account_id":   d.Get("account_id"),
		"account_name": d.Get("account_name"),
		"ak":           d.Get("ak"),
		"sk":           d.Get("sk"),
	}

	return bodyParams
}

func resourceCloudVendorAccountRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v1/vendor-account?limit=100"
		product = "coc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)

	getResp, err := pagination.ListAllItems(
		client,
		"offset",
		getPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving COC cloud vendor account")
	}

	getRespJson, err := json.Marshal(getResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var getRespBody interface{}
	err = json.Unmarshal(getRespJson, &getRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	cloudVendorAccount := utils.PathSearch(fmt.Sprintf("data[?id=='%s']|[0]", d.Id()), getRespBody, nil)
	if cloudVendorAccount == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("vendor", utils.PathSearch("vendor", cloudVendorAccount, nil)),
		d.Set("account_id", utils.PathSearch("account_id", cloudVendorAccount, nil)),
		d.Set("account_name", utils.PathSearch("account_name", cloudVendorAccount, nil)),
		d.Set("ak", utils.PathSearch("ak", cloudVendorAccount, nil)),
		d.Set("domain_id", utils.PathSearch("domain_id", cloudVendorAccount, nil)),
		d.Set("sync_status", utils.PathSearch("sync_status", cloudVendorAccount, nil)),
		d.Set("failure_msg", utils.PathSearch("failure_msg", cloudVendorAccount, nil)),
		d.Set("sync_date", utils.PathSearch("sync_date", cloudVendorAccount, nil)),
		d.Set("create_time", utils.PathSearch("create_time", cloudVendorAccount, nil)),
		d.Set("update_time", utils.PathSearch("update_time", cloudVendorAccount, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceCloudVendorAccountUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v1/vendor-account/{id}"
		product = "coc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	if d.HasChanges("account_name", "ak", "sk") {
		updatePath := client.Endpoint + httpUrl
		updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
		updatePath = strings.ReplaceAll(updatePath, "{id}", d.Id())

		updateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		updateOpt.JSONBody = buildUpdateCloudVendorAccountBodyParams(d)

		_, err = client.Request("PUT", updatePath, &updateOpt)
		if err != nil {
			return diag.Errorf("error updating COC cloud vendor account: %s", err)
		}
	}

	return resourceCloudVendorAccountRead(ctx, d, meta)
}

func buildUpdateCloudVendorAccountBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := make(map[string]interface{})
	if d.HasChange("account_name") {
		bodyParams["account_name"] = d.Get("account_name")
	}
	if d.HasChange("ak") {
		bodyParams["ak"] = d.Get("ak")
	}
	if d.HasChange("sk") {
		bodyParams["sk"] = d.Get("sk")
	}

	return bodyParams
}

func resourceCloudVendorAccountDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v1/vendor-account/{id}"
		product = "coc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{id}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code",
			"COC.00101102"), "error deleting COC cloud vendor account")
	}

	return nil
}
