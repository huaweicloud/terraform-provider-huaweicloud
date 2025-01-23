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

var nonUpdatableParamsApproval = []string{"workspace_id", "version_id", "result", "content"}

// @API SecMaster POST /v1/{project_id}/workspaces/{workspace_id}/soc/playbooks/versions/{version_id}/approval
func ResourcePlaybookApproval() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePlaybookApprovalCreate,
		UpdateContext: resourcePlaybookApprovalUpdate,
		ReadContext:   resourcePlaybookApprovalRead,
		DeleteContext: resourcePlaybookApprovalDelete,

		CustomizeDiff: config.FlexibleForceNew(nonUpdatableParamsApproval),

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
				Description: `Specifies the ID of the workspace to which the playbook belongs.`,
			},
			"version_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the version ID of the playbook.`,
			},
			"result": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the result of playbook approval.`,
			},
			"content": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the content of the playbook approval.`,
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

func resourcePlaybookApprovalCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		createPlaybookApprovalHttpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/playbooks/versions/{version_id}/approval"
		createPlaybookApprovalProduct = "secmaster"
	)
	createPlaybookApprovalClient, err := cfg.NewServiceClient(createPlaybookApprovalProduct, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	createPlaybookApprovalPath := createPlaybookApprovalClient.Endpoint + createPlaybookApprovalHttpUrl
	createPlaybookApprovalPath = strings.ReplaceAll(createPlaybookApprovalPath, "{project_id}", createPlaybookApprovalClient.ProjectID)
	createPlaybookApprovalPath = strings.ReplaceAll(createPlaybookApprovalPath, "{workspace_id}", d.Get("workspace_id").(string))
	createPlaybookApprovalPath = strings.ReplaceAll(createPlaybookApprovalPath, "{version_id}", d.Get("version_id").(string))

	createPlaybookApprovalOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	createPlaybookApprovalOpt.JSONBody = map[string]interface{}{
		"result":  utils.ValueIgnoreEmpty(d.Get("result")),
		"content": utils.ValueIgnoreEmpty(d.Get("content")),
	}

	_, err = createPlaybookApprovalClient.Request("POST", createPlaybookApprovalPath, &createPlaybookApprovalOpt)
	if err != nil {
		return diag.Errorf("error creating playbook approval: %s", err)
	}

	d.SetId(d.Get("version_id").(string))

	return nil
}

func resourcePlaybookApprovalRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourcePlaybookApprovalUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourcePlaybookApprovalDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for playbook approval resource. Deleting this resource will
		not change the status of the currently playbook approval resource, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
