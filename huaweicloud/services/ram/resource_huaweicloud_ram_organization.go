// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product RAM
// ---------------------------------------------------------------

package ram

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RAM POST /v1/organization-share/enable
// @API RAM POST /v1/organization-share/disable
// @API RAM GET /v1/organization-share
func ResourceRAMOrganization() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRAMOrganizationCreate,
		UpdateContext: resourceRAMOrganizationUpdate,
		ReadContext:   resourceRAMOrganizationRead,
		DeleteContext: resourceRAMOrganizationDelete,

		Schema: map[string]*schema.Schema{
			"enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: `Specifies whether sharing with organizations is enabled.`,
			},
		},
	}
}

func resourceRAMOrganizationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		product  = "ram"
		enabled  = d.Get("enabled").(bool)
		ramError error
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RAM client: %s", err)
	}

	if enabled {
		ramError = enableRAMOrganization(client)
	} else {
		ramError = disableRAMOrganization(client)
	}

	if ramError != nil {
		return diag.FromErr(ramError)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(generateUUID)
	return resourceRAMOrganizationRead(ctx, d, meta)
}

func resourceRAMOrganizationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		product  = "ram"
		enabled  = d.Get("enabled").(bool)
		ramError error
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RAM client: %s", err)
	}

	if enabled {
		ramError = enableRAMOrganization(client)
	} else {
		ramError = disableRAMOrganization(client)
	}

	if ramError != nil {
		return diag.FromErr(ramError)
	}

	return resourceRAMOrganizationRead(ctx, d, meta)
}

func enableRAMOrganization(client *golangsdk.ServiceClient) error {
	enableOrganizationPath := client.Endpoint + "v1/organization-share/enable"
	enableOrganizationOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err := client.Request("POST", enableOrganizationPath, &enableOrganizationOpt)
	if err != nil {
		return fmt.Errorf("error enabling RAM organization: %s", err)
	}
	return nil
}

func disableRAMOrganization(client *golangsdk.ServiceClient) error {
	disableOrganizationPath := client.Endpoint + "v1/organization-share/disable"
	disableOrganizationOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err := client.Request("POST", disableOrganizationPath, &disableOrganizationOpt)
	if err != nil {
		return fmt.Errorf("error disabling RAM organization: %s", err)
	}
	return nil
}

func resourceRAMOrganizationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "ram"
		httpUrl = "v1/organization-share"
		mErr    *multierror.Error
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RAM client: %s", err)
	}

	getOrganizationPath := client.Endpoint + httpUrl
	getOrganizationOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getOrganizationResp, err := client.Request("GET", getOrganizationPath, &getOrganizationOpt)
	if err != nil {
		return diag.Errorf("error retrieving RAM organization: %s", err)
	}

	getOrganizationRespBody, err := utils.FlattenResponse(getOrganizationResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("enabled", utils.PathSearch("enabled", getOrganizationRespBody, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceRAMOrganizationDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}
