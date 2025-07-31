package secmaster

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var nonUpdatableParamsWorkspace = []string{
	"project_name",
	"enterprise_project_id",
	"enterprise_project_name",
	"is_view",
	"tags",
}

// @API SecMaster POST /v1/{project_id}/workspaces
// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}
// @API SecMaster PUT /v1/{project_id}/workspaces/{workspace_id}
// @API SecMaster DELETE /v1/{project_id}/workspaces/{workspace_id}
func ResourceWorkspace() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWorkspaceCreate,
		UpdateContext: resourceWorkspaceUpdate,
		ReadContext:   resourceWorkspaceRead,
		DeleteContext: resourceWorkspaceDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: customdiff.All(
			config.FlexibleForceNew(nonUpdatableParamsWorkspace),
			config.MergeDefaultTags(),
		),

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
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enterprise_project_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"view_bind_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"is_view": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"tags": common.TagsSchema(),
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creator_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creator_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modifier_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modifier_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"view_bind_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"workspace_agency_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     workspaceAgencyListSchema(),
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

func workspaceAgencyListSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"workspace_attribution": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"agency_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"iam_agency_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"iam_agency_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_spec_code": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"selected": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
	return &sc
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

	return resourceWorkspaceRead(ctx, d, meta)
}

func resourceWorkspaceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	getWorkspaceClient, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	ws, err := GetWorkspace(getWorkspaceClient, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "code", "SecMaster.20011003"), "error retrieving SecMaster workspace")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("create_time", utils.PathSearch("workspace.create_time", ws, nil)),
		d.Set("creator_id", utils.PathSearch("workspace.creator_id", ws, nil)),
		d.Set("creator_name", utils.PathSearch("workspace.creator_name", ws, nil)),
		d.Set("description", utils.PathSearch("workspace.description", ws, nil)),
		d.Set("domain_id", utils.PathSearch("workspace.domain_id", ws, nil)),
		d.Set("domain_name", utils.PathSearch("workspace.domain_name", ws, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("workspace.enterprise_project_id", ws, nil)),
		d.Set("enterprise_project_name", utils.PathSearch("workspace.enterprise_project_name", ws, nil)),
		d.Set("is_view", utils.PathSearch("workspace.is_view", ws, nil)),
		d.Set("modifier_id", utils.PathSearch("workspace.modifier_id", ws, nil)),
		d.Set("modifier_name", utils.PathSearch("workspace.modifier_name", ws, nil)),
		d.Set("name", utils.PathSearch("workspace.name", ws, nil)),
		d.Set("project_name", utils.PathSearch("workspace.project_name", ws, nil)),
		d.Set("update_time", utils.PathSearch("workspace.update_time", ws, nil)),
		d.Set("view_bind_id", utils.PathSearch("workspace.view_bind_id", ws, nil)),
		d.Set("view_bind_name", utils.PathSearch("workspace.view_bind_name", ws, nil)),
		d.Set("workspace_agency_list", flattenWorkspaceAgencyList(
			utils.PathSearch("workspace.workspace_agency_list", ws, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenWorkspaceAgencyList(workspaceAgencyList []interface{}) []interface{} {
	if len(workspaceAgencyList) == 0 {
		return nil
	}

	rst := make([]interface{}, len(workspaceAgencyList))
	for i, v := range workspaceAgencyList {
		rst[i] = map[string]interface{}{
			"project_id":            utils.PathSearch("project_id", v, nil),
			"id":                    utils.PathSearch("id", v, nil),
			"name":                  utils.PathSearch("name", v, nil),
			"region_id":             utils.PathSearch("region_id", v, nil),
			"workspace_attribution": utils.PathSearch("workspace_attribution", v, nil),
			"agency_version":        utils.PathSearch("agency_version", v, nil),
			"domain_id":             utils.PathSearch("domain_id", v, nil),
			"domain_name":           utils.PathSearch("domain_name", v, nil),
			"iam_agency_id":         utils.PathSearch("iam_agency_id", v, nil),
			"iam_agency_name":       utils.PathSearch("iam_agency_name", v, nil),
			"resource_spec_code":    utils.PathSearch("resource_spec_code", v, nil),
			"selected":              utils.PathSearch("selected", v, nil),
		}
	}

	return rst
}

func resourceWorkspaceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		updateWorkspaceHttpUrl = "v1/{project_id}/workspaces/{workspace_id}"
		updateWorkspaceProduct = "secmaster"
	)
	updateWorkspaceClient, err := cfg.NewServiceClient(updateWorkspaceProduct, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	updateWorkspacePath := updateWorkspaceClient.Endpoint + updateWorkspaceHttpUrl
	updateWorkspacePath = strings.ReplaceAll(updateWorkspacePath, "{project_id}", updateWorkspaceClient.ProjectID)
	updateWorkspacePath = strings.ReplaceAll(updateWorkspacePath, "{workspace_id}", d.Id())

	updateWorkspaceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=UTF-8",
			"X-Language":   "en-us",
		},
	}

	updateOpts := map[string]interface{}{
		"name":         d.Get("name"),
		"description":  utils.ValueIgnoreEmpty(d.Get("description")),
		"view_bind_id": utils.ValueIgnoreEmpty(d.Get("view_bind_id")),
	}

	updateWorkspaceOpt.JSONBody = utils.RemoveNil(updateOpts)

	_, err = updateWorkspaceClient.Request("PUT", updateWorkspacePath, &updateWorkspaceOpt)
	if err != nil {
		return diag.Errorf("error updating SecMaster workspace: %s", err)
	}

	return resourceWorkspaceRead(ctx, d, meta)
}

func resourceWorkspaceDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		deleteWorkspaceHttpUrl = "v1/{project_id}/workspaces/{workspace_id}"
		deleteWorkspaceProduct = "secmaster"
	)
	deleteWorkspaceClient, err := cfg.NewServiceClient(deleteWorkspaceProduct, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	deleteWorkspacePath := deleteWorkspaceClient.Endpoint + deleteWorkspaceHttpUrl
	deleteWorkspacePath = strings.ReplaceAll(deleteWorkspacePath, "{project_id}", deleteWorkspaceClient.ProjectID)
	deleteWorkspacePath = strings.ReplaceAll(deleteWorkspacePath, "{workspace_id}", d.Id())

	deleteWorkspaceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=UTF-8",
			"X-Language":   "en-us",
		},
	}

	_, err = deleteWorkspaceClient.Request("DELETE", deleteWorkspacePath, &deleteWorkspaceOpt)
	if err != nil {
		return diag.Errorf("error deleting SecMaster workspace: %s", err)
	}

	return nil
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

func GetWorkspace(client *golangsdk.ServiceClient, workspaceId string) (interface{}, error) {
	getWorkspaceHttpUrl := "v1/{project_id}/workspaces/{workspace_id}"
	getWorkspacePath := client.Endpoint + getWorkspaceHttpUrl
	getWorkspacePath = strings.ReplaceAll(getWorkspacePath, "{project_id}", client.ProjectID)
	getWorkspacePath = strings.ReplaceAll(getWorkspacePath, "{workspace_id}", workspaceId)

	getWorkspaceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=UTF-8",
			"X-Language":   "en-us",
		},
	}

	getWorkspaceResp, err := client.Request("GET", getWorkspacePath, &getWorkspaceOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getWorkspaceResp)
}
