package dew

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var verifyMacNonUpdatableParams = []string{
	"key_id",
	"mac_algorithm",
	"message",
	"mac",
}

// @API DEW POST /v1.0/{project_id}/kms/verify-mac
func ResourceVerifyMac() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVerifyMacCreate,
		ReadContext:   resourceVerifyMacRead,
		UpdateContext: resourceVerifyMacUpdate,
		DeleteContext: resourceVerifyMacDelete,

		CustomizeDiff: config.FlexibleForceNew(verifyMacNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"key_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"mac_algorithm": {
				Type:     schema.TypeString,
				Required: true,
			},
			"message": {
				Type:     schema.TypeString,
				Required: true,
			},
			"mac": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"mac_valid": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func buildVerifyMacBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"key_id":        d.Get("key_id"),
		"mac_algorithm": d.Get("mac_algorithm"),
		"message":       d.Get("message"),
		"mac":           d.Get("mac"),
	}

	return bodyParams
}

func resourceVerifyMacCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1.0/{project_id}/kms/verify-mac"
	)

	client, err := cfg.NewServiceClient("kms", region)
	if err != nil {
		return diag.Errorf("error creating KMS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildVerifyMacBodyParams(d),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error verifying MAC: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resourceId)

	mErr := multierror.Append(nil,
		d.Set("mac_valid", utils.PathSearch("mac_valid", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceVerifyMacRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceVerifyMacUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceVerifyMacDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Delete()' method because the resource is a one-time action resource.
	return nil
}
