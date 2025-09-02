package apig

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

var apigApiActionNonUpdatableParams = []string{"instance_id", "api_id", "env_id", "action", "remark"}

// @API APIG POST /v2/{project_id}/apigw/instances/{instance_id}/apis/action
func ResourceApigApiAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApigApiActionCreate,
		ReadContext:   resourceApigApiActionRead,
		UpdateContext: resourceApigApiActionUpdate,
		DeleteContext: resourceApigApiActionDelete,

		CustomizeDiff: config.FlexibleForceNew(apigApiActionNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the APIG instance to which the API belongs is located.`,
			},

			// Required parameters.
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the dedicated instance to which the API belongs.`,
			},
			"api_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the API to be published.`,
			},
			"env_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the environment to which the API will be published.`,
			},
			"action": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The operation on the API will be performed.`,
			},

			// Optional parameters.
			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the publish action.`,
			},

			// Attributes.
			"publish_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the publish record.`,
			},
			"api_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the API.`,
			},
			"publish_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The time when the API was published, in UTC format.`,
			},
			"version_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The version ID of the online API.`,
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

func buildApiPublishBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"env_id": d.Get("env_id"),
		"api_id": d.Get("api_id"),
		"action": d.Get("action"),
		"remark": utils.ValueIgnoreEmpty(d.Get("remark")),
	}
}

func flattenApiPublishResponse(d *schema.ResourceData, respBody map[string]interface{}) error {
	if d.Get("action").(string) == "online" {
		mErr := multierror.Append(nil,
			d.Set("publish_id", utils.PathSearch("publish_id", respBody, "")),
			d.Set("api_name", utils.PathSearch("api_name", respBody, "")),
			d.Set("publish_time", utils.PathSearch("publish_time", respBody, "")),
			d.Set("version_id", utils.PathSearch("version_id", respBody, "")),
		)
		return mErr.ErrorOrNil()
	}
	return nil
}

func resourceApigApiActionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/apigw/instances/{instance_id}/apis/action"
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", d.Get("instance_id").(string))
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildApiPublishBodyParams(d),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error performing an operation with the API: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	err = flattenApiPublishResponse(d, respBody.(map[string]interface{}))
	if err != nil {
		return diag.Errorf("error setting API publish attributes: %s", err)
	}

	return resourceApigApiActionRead(ctx, d, meta)
}

func resourceApigApiActionRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceApigApiActionUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceApigApiActionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for performing an operation with the API. Deleting
this resource will not clear the corresponding request record, but will only remove the resource information from the
tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
