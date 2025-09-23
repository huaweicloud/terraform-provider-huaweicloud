package live

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

// @API Live POST /v1/{project_id}/auth/chain
func ResourceUrlAuthentication() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUrlAuthenticationCreate,
		ReadContext:   resourceUrlAuthenticationRead,
		DeleteContext: resourceUrlAuthenticationDelete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the domain name to which the URL validation belongs.`,
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the type of the domain name.`,
			},
			"stream_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the stream name.`,
			},
			"app_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the application name.`,
			},
			"check_level": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the check level.`,
			},
			"start_time": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the start time of the valid access time defined by the user.`,
			},
			"key_chain": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The generated signed URLs.`,
			},
		},
	}
}

func buildSignUrlBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"domain":      d.Get("domain_name"),
		"domain_type": d.Get("type"),
		"stream":      d.Get("stream_name"),
		"app":         d.Get("app_name"),
		"check_level": utils.ValueIgnoreEmpty(d.Get("check_level")),
		"start_time":  utils.ValueIgnoreEmpty(d.Get("start_time")),
	}
}

func resourceUrlAuthenticationCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/auth/chain"
	)

	client, err := cfg.NewServiceClient("live", region)
	if err != nil {
		return diag.Errorf("error creating Live client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildSignUrlBodyParams(d)),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating URL authentication: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	keyChain := utils.PathSearch("keychain", respBody, make([]interface{}, 0)).([]interface{})
	if len(keyChain) == 0 {
		return diag.Errorf("err creating URL authentication, the 'key_chain' not found in API response")
	}

	resourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(resourceId)

	mErr := multierror.Append(
		d.Set("key_chain", keyChain),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceUrlAuthenticationRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// This resource is only a one-time action resource for generating a signer URL
	return nil
}

func resourceUrlAuthenticationDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// This resource is only a one-time action resource for generating a signer URL
	return nil
}
