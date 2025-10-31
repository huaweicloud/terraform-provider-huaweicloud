package identitycenter

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var identityCenterRegisteredRegionParams = []string{"region_id"}

// @API IdentityCenter POST /v1/register-regions
// @API IdentityCenter GET  /v1/registered-regions
func ResourceIdentityCenterRegisteredRegion() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityCenterRegisteredRegionCreate,
		ReadContext:   resourceIdentityCenterRegisteredRegionRead,
		DeleteContext: resourceIdentityCenterRegisteredRegionDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(identityCenterRegisteredRegionParams),

		Schema: map[string]*schema.Schema{
			"region_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceIdentityCenterRegisteredRegionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	regionId := d.Get("region_id").(string)

	var (
		httpUrl = "v1/register-regions"
		product = "identitycenter"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating IdentityCenter client: %s", err)
	}

	registerPath := client.Endpoint + httpUrl

	registerOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildRegisterRegionBodyParams(d)),
	}

	_, err = client.Request("POST", registerPath, &registerOpt)
	if err != nil {
		return diag.Errorf("error registering IdentityCenter region: %s", err)
	}

	d.SetId(regionId)

	return resourceIdentityCenterRegisteredRegionRead(ctx, d, meta)
}

func buildRegisterRegionBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"region_id": utils.ValueIgnoreEmpty(d.Get("region_id")),
	}
	return bodyParams
}

func resourceIdentityCenterRegisteredRegionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		getHttpUrl = "v1/registered-regions"
		product    = "identitycenter"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center Client: %s", err)
	}

	getPath := client.Endpoint + getHttpUrl

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getResp, err := client.Request("GET", getPath, &getOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Identity Center registered region")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	regions := utils.PathSearch(fmt.Sprintf("regions|[?region_id =='%s']|[0]", d.Id()), getRespBody, nil)
	if regions == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "no region found.")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region_id", utils.PathSearch("region_id", regions, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceIdentityCenterRegisteredRegionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting IdentityCenter region resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
