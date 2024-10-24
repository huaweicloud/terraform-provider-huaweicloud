package iotda

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iotda/v5/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IoTDA POST /v5/iot/{project_id}/auth/accesscode
func ResourceAccessCredential() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAccessCredentialCreate,
		ReadContext:   resourceAccessCredentialRead,
		DeleteContext: resourceAccessCredentialDelete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"force_disconnect": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"access_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildAccessCredentialBodyParams(d *schema.ResourceData) *model.CreateAccessCodeRequest {
	bodyParams := model.CreateAccessCodeRequest{
		Body: &model.CreateAccessCodeRequestBody{
			Type:            utils.StringIgnoreEmpty(d.Get("type").(string)),
			ForceDisconnect: utils.Bool(d.Get("force_disconnect").(bool)),
		},
	}

	return &bodyParams
}

func resourceAccessCredentialCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	isDerived := WithDerivedAuth(cfg, region)
	client, err := cfg.HcIoTdaV5Client(region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	createOpts := buildAccessCredentialBodyParams(d)
	respBody, err := client.CreateAccessCode(createOpts)
	if err != nil || respBody == nil {
		return diag.Errorf("error creating access credential: %s", err)
	}

	resourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(resourceId)

	accessKey := utils.StringValue(respBody.AccessCode)
	accessCode := utils.StringValue(respBody.AccessKey)
	if accessKey == "" || accessCode == "" {
		return diag.Errorf("error creating access credential: 'access_key' or 'access_code' is empty in API response")
	}

	mErr := multierror.Append(
		d.Set("access_key", accessKey),
		d.Set("access_code", accessCode),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceAccessCredentialRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a action resource.
	return nil
}

func resourceAccessCredentialDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Delete()' method because the resource is a action resource.
	return nil
}
