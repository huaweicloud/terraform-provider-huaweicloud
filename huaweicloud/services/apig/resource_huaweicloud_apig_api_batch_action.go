package apig

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var apigApiBatchActionNonUpdatableParams = []string{"instance_id", "action", "env_id", "apis", "group_id", "remark"}

// @API APIG POST /v2/{project_id}/apigw/instances/{instance_id}/apis/publish
func ResourceApigApiBatchAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApigApiBatchActionCreate,
		ReadContext:   resourceApigApiBatchActionRead,
		UpdateContext: resourceApigApiBatchActionUpdate,
		DeleteContext: resourceApigApiBatchActionDelete,

		CustomizeDiff: config.FlexibleForceNew(apigApiBatchActionNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the APIG instance to which the API belongs is located`,
			},

			// Required parameters.
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the dedicated instance to which the APIs belong.`,
			},
			"action": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The action to perform on the APIs.`,
			},
			"env_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the environment where the action will be performed.`,
			},

			// Optional parameters.
			"apis": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1000,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of API IDs to perform the action on.`,
			},
			"group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the API group.`,
			},
			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The remark for the batch operation.`,
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

func buildApiBatchActionBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"env_id": d.Get("env_id"),
	}

	if v, ok := d.GetOk("apis"); ok {
		bodyParams["apis"] = v
	}
	if v, ok := d.GetOk("group_id"); ok {
		bodyParams["group_id"] = v
	}
	if v, ok := d.GetOk("remark"); ok {
		bodyParams["remark"] = v
	}
	return bodyParams
}

func buildBatchActionErrorMessages(failureList []interface{}) string {
	if len(failureList) < 1 {
		return ""
	}

	errorMsgs := make([]string, 0)
	for _, failure := range failureList {
		apiId := utils.PathSearch("api_id", failure, "unknown")
		errorCode := utils.PathSearch("error_code", failure, "unknown")
		errorMsg := utils.PathSearch("error_msg", failure, "unknown error")
		errorMsgs = append(errorMsgs, fmt.Sprintf("API %v failed: %v - %v", apiId, errorCode, errorMsg))
	}
	return strings.Join(errorMsgs, "\n")
}

func resourceApigApiBatchActionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
		action     = d.Get("action").(string)
		httpUrl    = "v2/{project_id}/apigw/instances/{instance_id}/apis/publish?action={action}"
	)

	client, err := cfg.ApigV2Client(region)
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)
	createPath = strings.ReplaceAll(createPath, "{action}", action)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildApiBatchActionBodyParams(d),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error performing batch action on APIs: %s", err)
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

	failureList := utils.PathSearch("failure", respBody, make([]interface{}, 0)).([]interface{})
	if len(failureList) > 0 {
		errorMessage := buildBatchActionErrorMessages(failureList)
		return diag.Errorf("batch action failed for some APIs:\n%s", errorMessage)
	}

	return resourceApigApiBatchActionRead(ctx, d, meta)
}

func resourceApigApiBatchActionRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceApigApiBatchActionUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceApigApiBatchActionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for performing an operation with the API list.
Deleting this resource will not clear the corresponding request record, but will only remove the resource information
from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
