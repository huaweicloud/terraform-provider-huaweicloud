package secmaster

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
		},

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

func resourceWorkspaceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=UTF-8",
			"X-Language":   "en-us",
		},
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

	err = createWorkspaceWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the SecMaster workspace (%s) creation to complete: %s", d.Id(), err)
	}

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

func createWorkspaceWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			cfg := meta.(*config.Config)
			region := cfg.GetRegion(d)

			var (
				getWorkspaceHttpUrl = "v1/{project_id}/workspaces/{workspace_id}/sa/guides?version=1"
				getWorkspaceProduct = "secmaster"
			)
			getWorkspaceClient, err := cfg.NewServiceClient(getWorkspaceProduct, region)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error creating SecMaster client: %s", err)
			}

			getWorkspacePath := getWorkspaceClient.Endpoint + getWorkspaceHttpUrl
			getWorkspacePath = strings.ReplaceAll(getWorkspacePath, "{project_id}", getWorkspaceClient.ProjectID)
			getWorkspacePath = strings.ReplaceAll(getWorkspacePath, "{workspace_id}", d.Id())

			getWorkspaceOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
				MoreHeaders: map[string]string{
					"Content-Type": "application/json",
					"X-Language":   "en-us",
				},
			}

			getWorkspaceResp, err := getWorkspaceClient.Request("GET", getWorkspacePath, &getWorkspaceOpt)
			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault400); ok {
					return getWorkspaceResp, "PENDING", nil
				}
				return nil, "ERROR", err
			}
			return getWorkspaceResp, "COMPLETED", nil
		},
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}
