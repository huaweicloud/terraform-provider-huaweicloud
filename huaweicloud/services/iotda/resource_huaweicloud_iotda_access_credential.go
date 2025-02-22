package iotda

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

func buildCreateAccessCredentialBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"type":             utils.ValueIgnoreEmpty(d.Get("type")),
		"force_disconnect": d.Get("force_disconnect"),
	}
}

func resourceAccessCredentialCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/iot/{project_id}/auth/accesscode"
		product = "iotda"
	)

	isDerived := WithDerivedAuth(cfg, region)
	client, err := cfg.NewServiceClientWithDerivedAuth(product, region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateAccessCredentialBodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating access credential: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	accessKey := utils.PathSearch("access_key", respBody, "").(string)
	accessCode := utils.PathSearch("access_code", respBody, "").(string)
	if accessKey == "" || accessCode == "" {
		return diag.Errorf("error creating access credential: 'access_key' or 'access_code' is empty in API response")
	}

	resourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(resourceId)

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
