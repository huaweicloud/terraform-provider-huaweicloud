package swr

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var temporaryLoginCommandNonUpdatableParams = []string{
	"enhanced",
}

// @API SWR POST /v2/manage/utils/secret
// @API SWR POST /v2/manage/utils/authorizationtoken
func ResourceSwrTemporaryLoginCommand() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSwrTemporaryLoginCommandCreate,
		ReadContext:   resourceSwrTemporaryLoginCommandRead,
		UpdateContext: resourceSwrTemporaryLoginCommandUpdate,
		DeleteContext: resourceSwrTemporaryLoginCommandDelete,

		CustomizeDiff: config.FlexibleForceNew(temporaryLoginCommandNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"enhanced": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether to create enhanced login command.`,
			},
			"x_swr_docker_login": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the docker login command.`,
			},
			"auths": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the authentication information.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the authentication information key.`,
						},
						"auth": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the base64-encoded authentication information.`,
						},
					},
				},
			},
			"x_expire_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the expiration time of the login command.`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func resourceSwrTemporaryLoginCommandCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	createHttpUrl := "v2/manage/utils/secret"
	if d.Get("enhanced").(bool) {
		createHttpUrl = "v2/manage/utils/authorizationtoken"
	}
	createPath := client.Endpoint + createHttpUrl
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating SWR image auto sync: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	xSwrDockerLogin := createResp.Header.Get("X-Swr-Dockerlogin")
	xSwrExpireAt := createResp.Header.Get("X-Swr-Expireat")

	mErr := multierror.Append(nil,
		d.Set("x_swr_docker_login", xSwrDockerLogin),
		d.Set("x_expire_at", xSwrExpireAt),
		d.Set("auths", flattenSwrTemporaryLoginCommandResponse(createRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSwrTemporaryLoginCommandResponse(resp interface{}) interface{} {
	if auths, ok := utils.PathSearch("auths", resp, nil).(map[string]interface{}); ok && len(auths) > 0 {
		result := make([]interface{}, 0, len(auths))
		for k, v := range auths {
			auth := v.(map[string]interface{})
			m := map[string]interface{}{
				"key":  k,
				"auth": utils.PathSearch("auth", auth, nil),
			}
			result = append(result, m)
		}

		return result
	}
	return nil
}

func resourceSwrTemporaryLoginCommandRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceSwrTemporaryLoginCommandUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceSwrTemporaryLoginCommandDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting temporary login command resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
