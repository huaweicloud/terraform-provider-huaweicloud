package codeartspipeline

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var (
	pipelinePermissionNonUpdatableParams = []string{
		"project_id", "pipeline_id", "role_id", "user_id",
	}

	updateRolePermissionHttpUrl = "v5/{project_id}/api/pipeline-permissions/{pipeline_id}/update-role-permission"
	updateUserPermissionHttpUrl = "v5/{project_id}/api/pipeline-permissions/{pipeline_id}/update-user-permission"
)

// @API CodeArtsPipeline POST /v5/{project_id}/api/pipeline-permissions/{pipeline_id}/update-role-permission
// @API CodeArtsPipeline POST /v5/{project_id}/api/pipeline-permissions/{pipeline_id}/update-user-permission
// @API CodeArtsPipeline GET /v5/{project_id}/api/pipeline-permissions/{pipeline_id}/role-permission
// @API CodeArtsPipeline GET /v5/{project_id}/api/pipeline-permissions/{pipeline_id}/user-permission
func ResourceCodeArtsPipelinePermission() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePipelinePermissionCreateOrUpdate,
		ReadContext:   resourcePipelinePermissionRead,
		UpdateContext: resourcePipelinePermissionCreateOrUpdate,
		DeleteContext: resourcePipelinePermissionDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourcePipelinePermissionImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(pipelinePermissionNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the CodeArts project ID.`,
			},
			"pipeline_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the pipeline ID.`,
			},
			"role_id": {
				Type:         schema.TypeInt,
				Optional:     true,
				ExactlyOneOf: []string{"role_id", "user_id"},
				Description:  `Specifies the role ID.`,
			},
			"user_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the user ID.`,
			},
			"operation_query": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether the role has the permission to query.`,
			},
			"operation_execute": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether the role has the permission to execute.`,
			},
			"operation_update": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether the role has the permission to update.`,
			},
			"operation_delete": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether the role has the permission to delete.`,
			},
			"operation_authorize": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether the role has the permission to authorize.`,
			},
			"role_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the role name.`,
			},
			"user_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the user name.`,
			},
		},
	}
}

func resourcePipelinePermissionCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("codearts_pipeline", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CodeArts pipeline client: %s", err)
	}

	projectId := d.Get("project_id").(string)
	pipelineId := d.Get("pipeline_id").(string)
	userId := d.Get("user_id").(string)

	httpUrl := updateRolePermissionHttpUrl
	if userId != "" {
		httpUrl = updateUserPermissionHttpUrl
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", projectId)
	updatePath = strings.ReplaceAll(updatePath, "{pipeline_id}", pipelineId)
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildPipelinePermissionBodyParams(d)),
	}

	updateResp, err := client.Request("POST", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating pipeline permission: %s", err)
	}
	updateRespBody, err := utils.FlattenResponse(updateResp)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := checkResponseError(updateRespBody, ""); err != nil {
		return diag.Errorf("error updating pipeline permission: %s", err)
	}

	if d.IsNewResource() {
		if userId != "" {
			d.SetId(projectId + "/" + pipelineId + "/user/" + userId)
		} else {
			d.SetId(projectId + "/" + pipelineId + "/role/" + strconv.Itoa(d.Get("role_id").(int)))
		}
	}

	return resourcePipelinePermissionRead(ctx, d, meta)
}

func buildPipelinePermissionBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"role_id":             utils.ValueIgnoreEmpty(d.Get("role_id")),
		"user_id":             utils.ValueIgnoreEmpty(d.Get("user_id")),
		"operation_query":     d.Get("operation_query"),
		"operation_execute":   d.Get("operation_execute"),
		"operation_update":    d.Get("operation_update"),
		"operation_delete":    d.Get("operation_delete"),
		"operation_authorize": d.Get("operation_authorize"),
	}
}

func resourcePipelinePermissionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts pipeline client: %s", err)
	}

	projectId := d.Get("project_id").(string)
	pipelineId := d.Get("pipeline_id").(string)
	userId := d.Get("user_id").(string)

	var rst interface{}
	if userId != "" {
		rst, err = GetPipelineUesrPermissions(client, projectId, pipelineId, userId)
		if err != nil {
			return common.CheckDeletedDiag(d, err, "error retrieving CodeArts pipeline user permissions")
		}
	} else {
		rst, err = GetPipelineRolePermissions(client, projectId, pipelineId, strconv.Itoa(d.Get("role_id").(int)))
		if err != nil {
			return common.CheckDeletedDiag(d, err, "error retrieving CodeArts pipeline role permissions")
		}
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("operation_query", utils.PathSearch("operation_query", rst, nil)),
		d.Set("operation_execute", utils.PathSearch("operation_execute", rst, nil)),
		d.Set("operation_update", utils.PathSearch("operation_update", rst, nil)),
		d.Set("operation_delete", utils.PathSearch("operation_delete", rst, nil)),
		d.Set("operation_authorize", utils.PathSearch("operation_authorize", rst, nil)),
		d.Set("role_name", utils.PathSearch("role_name", rst, nil)),
		d.Set("user_name", utils.PathSearch("user_name", rst, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func GetPipelineUesrPermissions(client *golangsdk.ServiceClient, projectId, pipelineId, userId string) (interface{}, error) {
	getHttpUrl := "v5/{project_id}/api/pipeline-permissions/{pipeline_id}/user-permission?limit=10"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", projectId)
	getPath = strings.ReplaceAll(getPath, "{pipeline_id}", pipelineId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	// have to send
	getPath += "&subject="

	offset := 0
	for {
		currentPath := getPath + fmt.Sprintf("&offset=%d", offset)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return nil, err
		}
		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return nil, fmt.Errorf("error flatten response: %s", err)
		}
		if err := checkResponseError(getRespBody, projectNotFoundError2); err != nil {
			return nil, err
		}

		expression := fmt.Sprintf("users[?user_id=='%v']|[0]", userId)
		user := utils.PathSearch(expression, getRespBody, nil)
		if user != nil {
			return user, nil
		}

		users := utils.PathSearch("users", getRespBody, make([]interface{}, 0)).([]interface{})
		if len(users) == 0 {
			return nil, golangsdk.ErrDefault404{}
		}

		offset += 10
	}
}

func GetPipelineRolePermissions(client *golangsdk.ServiceClient, projectId, pipelineId, roleId string) (interface{}, error) {
	getHttpUrl := "v5/{project_id}/api/pipeline-permissions/{pipeline_id}/role-permission"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", projectId)
	getPath = strings.ReplaceAll(getPath, "{pipeline_id}", pipelineId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, fmt.Errorf("error flatten response: %s", err)
	}
	if err := checkResponseError(getRespBody, projectNotFoundError2); err != nil {
		return nil, err
	}

	expression := fmt.Sprintf("roles[?role_id==`%v`]|[0]", roleId)
	role := utils.PathSearch(expression, getRespBody, nil)
	if role == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return role, nil
}

func resourcePipelinePermissionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting permission resource is not supported. The resource is only removed from the state," +
		" the pipeline permission remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

func resourcePipelinePermissionImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 4 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<project_id>/<pipeline_id>/role/<role_id>'"+
			" or '<project_id>/<pipeline_id>/user/<user_id>', but got '%s'", d.Id())
	}

	mErr := multierror.Append(nil,
		d.Set("project_id", parts[0]),
		d.Set("pipeline_id", parts[1]),
	)

	if parts[2] == "role" {
		roleId, err := strconv.Atoi(parts[3])
		if err != nil {
			return nil, fmt.Errorf("invalid role_id', got '%s'", parts[3])
		}
		mErr = multierror.Append(mErr,
			d.Set("role_id", roleId),
		)
	} else {
		mErr = multierror.Append(mErr,
			d.Set("user_id", parts[3]),
		)
	}

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
