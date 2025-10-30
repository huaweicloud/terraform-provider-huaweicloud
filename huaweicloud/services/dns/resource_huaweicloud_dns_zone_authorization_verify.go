package dns

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var zoneAuthorizationVerifyNonUpdatableParams = []string{"authorization_id"}

// @API DNS POST /v2/authorize-txtrecord/{id}/verify
func ResourceZoneAuthorizationVerify() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceZoneAuthorizationVerifyCreate,
		ReadContext:   resourceZoneAuthorizationVerifyRead,
		UpdateContext: resourceZoneAuthorizationVerifyUpdate,
		DeleteContext: resourceZoneAuthorizationVerifyDelete,

		CustomizeDiff: config.FlexibleForceNew(zoneAuthorizationVerifyNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			// Required parameters.
			"authorization_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The request ID of the sub-domain authorization to be verified.`,
			},

			// Attributes.
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The authorization status after verification.`,
			},
		},
	}
}

func resourceZoneAuthorizationVerifyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)

	client, err := cfg.NewServiceClient("dns", "")
	if err != nil {
		return diag.Errorf("error creating DNS client: %s", err)
	}

	httpUrl := "v2/authorize-txtrecord/{id}/verify"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{id}", d.Get("authorization_id").(string))

	opt := golangsdk.RequestOpts{KeepResponseBody: true}
	resp, err := client.Request("POST", createPath, &opt)
	if err != nil {
		return diag.Errorf("error verifying sub-domain authorization: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.Errorf("error parsing verification response of sub-domain authorization: %s", err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("status", utils.PathSearch("status", respBody, nil)),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return resourceZoneAuthorizationVerifyRead(ctx, d, meta)
}

func resourceZoneAuthorizationVerifyRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceZoneAuthorizationVerifyUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceZoneAuthorizationVerifyDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	warning := `This resource is a one-time action for verify sub-domain authorization. Deleting this resource will not
clear the verification result in the service, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  warning,
		},
	}
}
