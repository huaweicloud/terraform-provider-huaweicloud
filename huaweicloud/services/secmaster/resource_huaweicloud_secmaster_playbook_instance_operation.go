package secmaster

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var playbookInstanceOptNonUpdatableParams = []string{"workspace_id", "instance_id", "operation"}

// @API SecMaster POST /v1/{project_id}/workspaces/{workspace_id}/soc/playbooks/instances/{instance_id}/operation
func ResourcePlaybookInstanceOperation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePlaybookInstanceOperationCreate,
		UpdateContext: resourcePlaybookInstanceOperationUpdate,
		ReadContext:   resourcePlaybookInstanceOperationRead,
		DeleteContext: resourcePlaybookInstanceOperationDelete,

		CustomizeDiff: config.FlexibleForceNew(playbookInstanceOptNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the workspace to which the playbook instance belongs.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the playbook instance ID.`,
			},
			"operation": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the operation of the playbook instance.`,
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

func resourcePlaybookInstanceOperationCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	instanceId := d.Get("instance_id").(string)

	var (
		createHttpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/playbooks/instances/{instance_id}/operation"
		createProduct = "secmaster"
	)
	createClient, err := cfg.NewServiceClient(createProduct, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	createPath := createClient.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", createClient.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{workspace_id}", d.Get("workspace_id").(string))
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	createOpt.JSONBody = map[string]interface{}{
		"operation": d.Get("operation"),
	}

	_, err = createClient.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating playbook instance operation: %s", err)
	}

	d.SetId(instanceId)

	return nil
}

func resourcePlaybookInstanceOperationRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourcePlaybookInstanceOperationUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourcePlaybookInstanceOperationDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for playbook instance operation resource. Deleting this resource will not change
		the status of the currently playbook instance operation resource, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
