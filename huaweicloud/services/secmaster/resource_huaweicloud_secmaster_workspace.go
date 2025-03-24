package secmaster

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var nonUpdatableParamsWorkspace = []string{
	"name",
	"project_name",
	"description",
	"enterprise_project_id",
	"enterprise_project_name",
	"view_bind_id",
	"is_view",
	"tags",
}

// @API SecMaster POST /v1/{project_id}/workspaces
func ResourceWorkspace() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWorkspaceCreate,
		UpdateContext: resourceWorkspaceUpdate,
		ReadContext:   resourceWorkspaceRead,
		DeleteContext: resourceWorkspaceDelete,

		CustomizeDiff: config.FlexibleForceNew(nonUpdatableParamsWorkspace),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"view_bind_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_view": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"tags": common.TagsSchema(),
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func resourceWorkspaceCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		createWorkspaceHttpUrl = "v1/{project_id}/workspaces"
		createWorkspaceProduct = "secmaster"
	)
	createWorkspaceClient, err := cfg.NewServiceClient(createWorkspaceProduct, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	createWorkspacePath := createWorkspaceClient.Endpoint + createWorkspaceHttpUrl
	createWorkspacePath = strings.ReplaceAll(createWorkspacePath, "{project_id}", createWorkspaceClient.ProjectID)

	createWorkspaceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}

	createOpts := map[string]interface{}{
		"region_id":               region,
		"name":                    d.Get("name"),
		"project_name":            d.Get("project_name"),
		"description":             utils.ValueIgnoreEmpty(d.Get("description")),
		"enterprise_project_id":   cfg.GetEnterpriseProjectID(d),
		"enterprise_project_name": utils.ValueIgnoreEmpty(d.Get("enterprise_project_name")),
		"view_bind_id":            utils.ValueIgnoreEmpty(d.Get("view_bind_id")),
		"is_view":                 d.Get("is_view"),
		"tags":                    utils.ExpandResourceTagsMap(d.Get("tags").(map[string]interface{})),
	}

	createWorkspaceOpt.JSONBody = utils.RemoveNil(createOpts)

	resp, err := createWorkspaceClient.Request("POST", createWorkspacePath, &createWorkspaceOpt)
	if err != nil {
		return diag.Errorf("error creating SecMaster workspace: %s", err)
	}

	reponseBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("id", reponseBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating SecMaster workspace: ID is not found in API response")
	}

	d.SetId(id)

	return nil
}

func resourceWorkspaceRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceWorkspaceUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceWorkspaceDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for SecMaster workspace resource. Deleting this resource will
		not change the status of the currently SecMaster workspace resource, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
