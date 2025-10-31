package dns

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var ZoneAuthorizationNonUpdatableParams = []string{
	"zone_name",
}

// @API DNS POST /v2/authorize-txtrecord
func ResourceZoneAuthorization() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceZoneAuthorizationCreate,
		ReadContext:   resourceZoneAuthorizationRead,
		UpdateContext: resourceZoneAuthorizationUpdate,
		DeleteContext: resourceZoneAuthorizationDelete,

		CustomizeDiff: config.FlexibleForceNew(ZoneAuthorizationNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			// Required parameters.
			"zone_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the sub-domain to be authorized.`,
				DiffSuppressFunc: func(_, oldVal, newVal string, _ *schema.ResourceData) bool {
					return strings.TrimSuffix(oldVal, ".") == strings.TrimSuffix(newVal, ".")
				},
			},

			// Attributes.
			"second_level_zone_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The second-level domain name to which the sub-domain belongs.`,
			},
			"record": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The TXT record information.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The host record of the TXT record.`,
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The record value of the TXT record.`,
						},
					},
				},
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The authorization status.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the authorization, in RFC3339 format.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The latest update time of the authorization, in RFC3339 format.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func buildZoneAuthorizationBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"zone_name": d.Get("zone_name"),
	}
}

func flattenZoneAuthorizeRecord(record map[string]interface{}) []map[string]interface{} {
	if len(record) < 1 {
		return nil
	}

	return []map[string]interface{}{
		{
			"host":  utils.PathSearch("host", record, nil),
			"value": utils.PathSearch("value", record, nil),
		},
	}
}

func resourceZoneAuthorizationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)

	client, err := cfg.NewServiceClient("dns", "")
	if err != nil {
		return diag.Errorf("error creating DNS client: %s", err)
	}

	httpUrl := "v2/authorize-txtrecord"
	createPath := client.Endpoint + httpUrl

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: buildZoneAuthorizationBodyParams(d),
	}

	resp, err := client.Request("POST", createPath, &opt)
	if err != nil {
		return diag.Errorf("error creating sub-domain authorization: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.Errorf("error parsing response of sub-domain authorization: %s", err)
	}

	authorizationId := utils.PathSearch("id", respBody, "").(string)
	if authorizationId == "" {
		return diag.Errorf("unable to find the authorization ID from the API response")
	}
	d.SetId(authorizationId)

	mErr := multierror.Append(nil,
		d.Set("zone_name", utils.PathSearch("zone_name", respBody, nil)),
		d.Set("second_level_zone_name", utils.PathSearch("second_level_zone_name", respBody, nil)),
		d.Set("record", flattenZoneAuthorizeRecord(utils.PathSearch("record",
			respBody, make(map[string]interface{})).(map[string]interface{}))),
		d.Set("status", utils.PathSearch("status", respBody, nil)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("created_at",
			respBody, "").(string), "2006-01-02T15:04:05.000")/1000, false)),
		d.Set("updated_at", utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("updated_at",
			respBody, "").(string), "2006-01-02T15:04:05.000")/1000, false)),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return resourceZoneAuthorizationRead(ctx, d, meta)
}

func resourceZoneAuthorizationRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceZoneAuthorizationUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceZoneAuthorizationDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to request a sub-domain authorization when creating a
sub-domain prompts this following error:
'domain conflicts with other tenants, you need to add TXT authorization verification'.
Deleting this resource will not clear the corresponding authorization record, but will only remove the resource
information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
