package swrenterprise

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

var enterpriseTemporaryCredentialNonUpdatableParams = []string{
	"instance_id",
}

// @API SWR POST /v2/{project_id}/instances/{instance_id}/temp-credential
func ResourceSwrEnterpriseTemporaryCredential() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSwrEnterpriseTemporaryCredentialCreate,
		UpdateContext: resourceSwrEnterpriseTemporaryCredentialUpdate,
		ReadContext:   resourceSwrEnterpriseTemporaryCredentialRead,
		DeleteContext: resourceSwrEnterpriseTemporaryCredentialDelete,

		CustomizeDiff: config.FlexibleForceNew(enterpriseTemporaryCredentialNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the enterprise instance ID.`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"user_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the user ID.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the creation time.`,
			},
			"expire_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the expired time.`,
			},
			"auth_token": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: `Indicates the auth token.`,
			},
		},
	}
}

func resourceSwrEnterpriseTemporaryCredentialCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	createHttpUrl := "v2/{project_id}/instances/{instance_id}/temp-credential"
	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", d.Get("instance_id").(string))
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating SWR enterprise instance temporary credential: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("token_id", createRespBody, "").(string)
	if id == "" {
		return diag.Errorf("unable to find SWR enterprise instance temporary credential token ID from the API response")
	}

	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("auth_token", utils.PathSearch("auth_token", createRespBody, nil)),
		d.Set("user_id", utils.PathSearch("user_id", createRespBody, nil)),
		d.Set("created_at", utils.PathSearch("created_at", createRespBody, nil)),
		d.Set("expire_date", utils.PathSearch("expire_date", createRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceSwrEnterpriseTemporaryCredentialRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceSwrEnterpriseTemporaryCredentialUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceSwrEnterpriseTemporaryCredentialDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting SWR enterprise instance temporary credential resource is not supported. The resource is only " +
		"removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
