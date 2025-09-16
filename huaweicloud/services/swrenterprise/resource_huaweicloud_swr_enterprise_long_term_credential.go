package swrenterprise

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var enterpriseLongTermCredentialNonUpdatableParams = []string{
	"instance_id", "name",
}

// @API SWR POST /v2/{project_id}/instances/{instance_id}/long-term-credential
// @API SWR GET /v2/{project_id}/instances/{instance_id}/long-term-credentials
// @API SWR PUT /v2/{project_id}/instances/{instance_id}/long-term-credentials/{credential_id}
// @API SWR DELETE /v2/{project_id}/instances/{instance_id}/long-term-credentials/{credential_id}
func ResourceSwrEnterpriseLongTermCredential() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSwrEnterpriseLongTermCredentialCreate,
		UpdateContext: resourceSwrEnterpriseLongTermCredentialUpdate,
		ReadContext:   resourceSwrEnterpriseLongTermCredentialRead,
		DeleteContext: resourceSwrEnterpriseLongTermCredentialDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceSwrEnterpriseLongTermCredentialImportStateFunc,
		},

		CustomizeDiff: config.FlexibleForceNew(enterpriseLongTermCredentialNonUpdatableParams),

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
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the credential name.`,
			},
			"enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Specifies whether to enable the credential.`,
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
			"user_profile": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the user profile.`,
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

func resourceSwrEnterpriseLongTermCredentialCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	createHttpUrl := "v2/{project_id}/instances/{instance_id}/long-term-credential"
	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", d.Get("instance_id").(string))
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateSwrEnterpriseLongTermCredentialBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating SWR enterprise instance long term credential: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("token_id", createRespBody, "").(string)
	if id == "" {
		return diag.Errorf("unable to find SWR enterprise instance long term credential token ID from the API response")
	}

	d.SetId(id)

	authToken := utils.PathSearch("auth_token", createRespBody, "").(string)
	if authToken == "" {
		return diag.Errorf("unable to find SWR enterprise instance long term credential auth token from the API response")
	}
	if err := d.Set("auth_token", authToken); err != nil {
		return diag.Errorf("error saving auth token: %s", err)
	}

	if !d.Get("enable").(bool) {
		if err := updateSwrEnterpriseLongTermCredentialEnable(client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceSwrEnterpriseLongTermCredentialRead(ctx, d, meta)
}

func buildCreateSwrEnterpriseLongTermCredentialBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name": d.Get("name"),
	}

	return bodyParams
}

func updateSwrEnterpriseLongTermCredentialEnable(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	updateHttpUrl := "v2/{project_id}/instances/{instance_id}/long-term-credentials/{credential_id}"
	updatePath := client.Endpoint + updateHttpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", d.Get("instance_id").(string))
	updatePath = strings.ReplaceAll(updatePath, "{credential_id}", d.Id())
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildUpdateSwrEnterpriseLongTermCredentialEnableBodyParams(d),
	}

	_, err := client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("error updating SWR enterprise instance long term credential: %s", err)
	}

	return nil
}

func buildUpdateSwrEnterpriseLongTermCredentialEnableBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"enable": d.Get("enable"),
	}

	return bodyParams
}

func resourceSwrEnterpriseLongTermCredentialRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	getHttpUrl := "v2/{project_id}/instances/{instance_id}/long-term-credentials"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Get("instance_id").(string))
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving SWR enterprise instance long term credential")
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	searchPath := fmt.Sprintf("auth_tokens[?token_id=='%s']|[0]", d.Id())
	token := utils.PathSearch(searchPath, getRespBody, nil)
	if token == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving SWR enterprise instance long term credential")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", token, nil)),
		d.Set("enable", utils.PathSearch("enable", token, nil)),
		d.Set("user_id", utils.PathSearch("user_id", token, nil)),
		d.Set("user_profile", utils.PathSearch("user_profile", token, nil)),
		d.Set("created_at", utils.PathSearch("created_at", token, nil)),
		d.Set("expire_date", utils.PathSearch("expire_date", token, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceSwrEnterpriseLongTermCredentialUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	if d.HasChanges("enable") {
		if err := updateSwrEnterpriseLongTermCredentialEnable(client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceSwrEnterpriseLongTermCredentialRead(ctx, d, meta)
}

func resourceSwrEnterpriseLongTermCredentialDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	deleteHttpUrl := "v2/{project_id}/instances/{instance_id}/long-term-credentials/{credential_id}"
	deletePath := client.Endpoint + deleteHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", d.Get("instance_id").(string))
	deletePath = strings.ReplaceAll(deletePath, "{credential_id}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting SWR enterprise instance long term credential")
	}

	return nil
}

func resourceSwrEnterpriseLongTermCredentialImportStateFunc(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<instance_id>/<id>', but got '%s'", d.Id())
	}

	d.SetId(parts[1])
	if err := d.Set("instance_id", parts[0]); err != nil {
		return nil, fmt.Errorf("error saving instance ID: %s", err)
	}

	return []*schema.ResourceData{d}, nil
}
