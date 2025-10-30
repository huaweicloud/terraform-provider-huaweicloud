package coc

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var cloudVendorUserResourcesSyncNonUpdatableParams = []string{"vendor", "account_id"}

// @API RDS POST /v1/multicloud-resources/sync
func ResourceCloudVendorUserResourcesSync() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCloudVendorUserResourcesSyncCreate,
		ReadContext:   resourceCloudVendorUserResourcesSyncRead,
		UpdateContext: resourceCloudVendorUserResourcesSyncUpdate,
		DeleteContext: resourceCloudVendorUserResourcesSyncDelete,

		CustomizeDiff: config.FlexibleForceNew(cloudVendorUserResourcesSyncNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"vendor": {
				Type:     schema.TypeString,
				Required: true,
			},
			"account_id": {
				Type:     schema.TypeString,
				Optional: true,
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

func resourceCloudVendorUserResourcesSyncCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v1/multicloud-resources/sync"
		product = "coc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateCloudVendorUserResourcesSyncBodyParams(d))

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating COC cloud vendor user resources sync: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("data", createRespBody, "").(string)
	if id == "" {
		return diag.Errorf("unable to find data from the response")
	}

	d.SetId(id)

	return nil
}

func buildCreateCloudVendorUserResourcesSyncBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"vendor":     d.Get("vendor"),
		"account_id": utils.ValueIgnoreEmpty(d.Get("account_id")),
	}

	return bodyParams
}

func resourceCloudVendorUserResourcesSyncRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceCloudVendorUserResourcesSyncUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceCloudVendorUserResourcesSyncDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting COC cloud vendor user resources sync is not supported. The restoration record is only " +
		"removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
