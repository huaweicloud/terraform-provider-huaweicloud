package rfs

import (
	"context"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var stackRollbackNonUpdatableParams = []string{
	"stack_name",
	"stack_id",
}

// @API RFS POST /v1/{project_id}/stacks/{stack_name}/rollbacks
func ResourceStackRollback() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceStackRollbackCreate,
		ReadContext:   resourceStackRollbackRead,
		UpdateContext: resourceStackRollbackUpdate,
		DeleteContext: resourceStackRollbackDelete,

		CustomizeDiff: config.FlexibleForceNew(stackRollbackNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"stack_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"stack_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"deployment_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildStackRollbackBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"stack_id": utils.ValueIgnoreEmpty(d.Get("stack_id")),
	}
}

func resourceStackRollbackCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		stackName = d.Get("stack_name").(string)
		httpUrl   = "v1/{project_id}/stacks/{stack_name}/rollbacks"
	)

	client, err := cfg.NewServiceClient("rfs", region)
	if err != nil {
		return diag.Errorf("error creating RFS client: %s", err)
	}

	requestId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate UUID: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{stack_name}", stackName)
	requestOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Client-Request-Id": requestId,
		},
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildStackRollbackBodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error triggering RFS stack rollback: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(stackName)

	deploymentId := utils.PathSearch("deployment_id", respBody, "").(string)
	if deploymentId == "" {
		return diag.Errorf("error triggering RFS stack rollback: Deployment ID is not found in API response")
	}

	if err := d.Set("deployment_id", deploymentId); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceStackRollbackRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in 'Read()' method because resource is a one-time action resource.
	return nil
}

func resourceStackRollbackUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in 'Update()' method because resource is a one-time action resource.
	return nil
}

func resourceStackRollbackDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to trigger stack rollback. Deleting this resource
    will not cancel the rollback operation, but will only remove resource information from
    tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
