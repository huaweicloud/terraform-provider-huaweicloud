package dsc

import (
	"context"
	"fmt"
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

var nonUpdatableParamsDscAssetAuthorization = []string{"type"}

// @API DSC PUT /v1/{project_id}/sdg/asset/authorization
// @API DSC GET /v1/{project_id}/sdg/asset/authorization
func ResourceDscAssetAuthorization() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDscAssetAuthorizationCreate,
		ReadContext:   resourceDscAssetAuthorizationRead,
		UpdateContext: resourceDscAssetAuthorizationUpdate,
		DeleteContext: resourceDscAssetAuthorizationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDscAssetAuthorizationImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(nonUpdatableParamsDscAssetAuthorization),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"authorization_status": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"bigdata_authorization": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"dashboard_authorization": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"database_authorization": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"lts_authorization": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"mrs_authorization": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"obs_authorization": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func buildDscAssetAuthorizationBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"authorization_status": d.Get("authorization_status"),
	}
}

func updateDscAssetAuthorization(client *golangsdk.ServiceClient, assetType string, requestBody map[string]interface{}) error {
	httpUrl := "v1/{project_id}/sdg/asset/authorization"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath += fmt.Sprintf("?type=%s", assetType)

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         requestBody,
	}

	_, err := client.Request("PUT", updatePath, &updateOpt)
	return err
}

func resourceDscAssetAuthorizationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "dsc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	assetType := d.Get("type").(string)
	requestBody := buildDscAssetAuthorizationBodyParams(d)
	if err := updateDscAssetAuthorization(client, assetType, requestBody); err != nil {
		return diag.Errorf("error authorizing DSC asset in creation operation: %s", err)
	}

	d.SetId(assetType)

	return resourceDscAssetAuthorizationRead(ctx, d, meta)
}

func resourceDscAssetAuthorizationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "dsc"
		httpUrl = "v1/{project_id}/sdg/asset/authorization"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += fmt.Sprintf("?type=%s", d.Id())
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	resp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DSC asset authorization")
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("bigdata_authorization", utils.PathSearch("bigdata_authorization", respBody, nil)),
		d.Set("dashboard_authorization", utils.PathSearch("dashboard_authorization", respBody, nil)),
		d.Set("database_authorization", utils.PathSearch("database_authorization", respBody, nil)),
		d.Set("lts_authorization", utils.PathSearch("lts_authorization", respBody, nil)),
		d.Set("mrs_authorization", utils.PathSearch("mrs_authorization", respBody, nil)),
		d.Set("obs_authorization", utils.PathSearch("obs_authorization", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDscAssetAuthorizationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "dsc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	if d.HasChange("authorization_status") {
		requestBody := buildDscAssetAuthorizationBodyParams(d)
		if err := updateDscAssetAuthorization(client, d.Get("type").(string), requestBody); err != nil {
			return diag.Errorf("error updating DSC asset authorization: %s", err)
		}
	}

	return resourceDscAssetAuthorizationRead(ctx, d, meta)
}

func resourceDscAssetAuthorizationDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDscAssetAuthorizationImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, d.Set("type", d.Id())
}
