package codeartsdeploy

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CodeArtsDeploy POST /v1/applications
// @API CodeArtsDeploy GET /v1/applications/{app_id}/info
// @API CodeArtsDeploy PUT /v1/applications
// @API CodeArtsDeploy DELETE /v1/applications/{app_id}
// @API CodeArtsDeploy PUT /v1/applications/{app_id}/disable
// @API CodeArtsDeploy PUT /v3/applications/permission-level
// @API CodeArtsDeploy GET /v3/applications/permissions
// @API CodeArtsDeploy PUT /v1/projects/{project_id}/applications/groups/move
func ResourceDeployApplication() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDeployApplicationCreate,
		ReadContext:   resourceDeployApplicationRead,
		UpdateContext: resourceDeployApplicationUpdate,
		DeleteContext: resourceDeployApplicationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
		},

		// precheck values for `is_draft` and `operation_list`,
		// when `operation_list` is empty and `is_draft` is true, the application is actually created,
		// and it's only shown in list API, which means 404 will return when get single application.
		CustomizeDiff: func(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
			operationList := d.Get("operation_list").([]interface{})
			isDraft := d.Get("is_draft").(bool)

			if len(operationList) == 0 && !isDraft {
				return fmt.Errorf("the argument (operation_list) is required when application is not in draft status")
			}

			if oldValue, newValue := d.GetChange("group_id"); oldValue != "" && newValue == "" {
				return fmt.Errorf("group_id can not be updated to empty")
			}

			return nil
		},

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
				ForceNew:    true,
				Description: `Specifies the project ID for CodeArts service.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the application name.`,
			},
			"is_draft": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: `Specifies whether the application is in draft status.`,
			},
			"create_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the creation type.`,
			},
			"trigger_source": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies where a deployment task can be executed.`,
			},
			"artifact_source_system": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the source information transferred by the pipeline.`,
			},
			"artifact_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the artifact type for the pipeline source.`,
			},
			"template_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "schema: Deprecated; Currently, the field is useless for creating API.",
			},
			// `operation_list` has order
			"operation_list": {
				Type:        schema.TypeList,
				Elem:        deployApplicationOperationSchema(),
				Optional:    true,
				Description: `Specifies the deployment orchestration list information.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the application description.`,
			},
			"resource_pool_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the custom slave resource pool ID.`,
			},
			"group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the group ID.`,
			},
			"is_disable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether to disable the application.`,
			},
			"permission_level": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the permission level.`,
			},
			"is_care": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether the user has favorited the application.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The create time.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The update time.`,
			},
			"project_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The project name.`,
			},
			"can_modify": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether the user has the editing permission.`,
			},
			"can_disable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether the user has the permission to disable application.`,
			},
			"can_delete": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether the user has the deletion permission.`,
			},
			"can_view": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether the user has the view permission.`,
			},
			"can_execute": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether the user has the deployment permission`,
			},
			"can_copy": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether the user has the copy permission.`,
			},
			"can_manage": {
				Type:     schema.TypeBool,
				Computed: true,
				Description: `Indicates whether the user has the management permission, including adding, deleting,
modifying, querying deployment and permission modification.`,
			},
			"can_create_env": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether the user has the permission to create an environment.`,
			},
			"task_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The deployment task ID.`,
			},
			"task_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The deployment task name.`,
			},
			"steps": {
				Type:        schema.TypeMap,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: `The deployment steps.`,
			},
			"permission_matrix": {
				Type:        schema.TypeList,
				Elem:        deployApplicationPermissionMatrixSchema(),
				Computed:    true,
				Description: `Indicates the permission matrix.`,
			},
		},
	}
}

func deployApplicationOperationSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the step name.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the step description.`,
			},
			"code": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the download URL.`,
			},
			"params": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the parameter.`,
			},
			"entrance": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the entry function.`,
			},
			"version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the version.`,
			},
			"module_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the module ID.`,
			},
		},
	}
	return &sc
}

func deployApplicationPermissionMatrixSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"role_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the role ID.`,
			},
			"role_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the role name.`,
			},
			"role_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the role type.`,
			},
			"can_modify": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether the role has the editing permission.`,
			},
			"can_disable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether the role has the permission to disable application.`,
			},
			"can_delete": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether the role has the deletion permission.`,
			},
			"can_view": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether the role has the view permission.`,
			},
			"can_execute": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether the role has the deployment permission.`,
			},
			"can_copy": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether the role has the copy permission.`,
			},
			"can_manage": {
				Type:     schema.TypeBool,
				Computed: true,
				Description: `Indicates whether the role has the management permission, including adding, deleting,
modifying, querying deployment and permission modification.`,
			},
			"can_create_env": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether the role has the permission to create an environment.`,
			},
		},
	}
	return &sc
}

func resourceDeployApplicationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/applications"
		product = "codearts_deploy"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CodeArts deploy client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf-8",
		},
		JSONBody: utils.RemoveNil(buildCreateDeployApplicationBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating CodeArts deploy application: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	appId := utils.PathSearch("result.id", createRespBody, "").(string)
	if appId == "" {
		return diag.Errorf("unable to find the deploy application ID from the API response")
	}
	d.SetId(appId)

	// `is_disable` defaults to false when the application is created
	if d.Get("is_disable").(bool) {
		if err := updateDeployApplicationDisable(client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	// `permission_level` defaults to project level when the application is created
	if v := d.Get("permission_level").(string); v == "instance" {
		if err := updateDeployApplicationPermissionLevel(ctx, client, d, d.Timeout(schema.TimeoutCreate)); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceDeployApplicationRead(ctx, d, meta)
}

func updateDeployApplicationDisable(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	httpUrl := "v1/applications/{app_id}/disable"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{app_id}", d.Id())
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf-8",
		},
		JSONBody: map[string]interface{}{
			"is_disable": d.Get("is_disable"),
		},
	}

	_, err := client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("error updating CodeArts deploy application: %s", err)
	}

	return nil
}

func updateDeployApplicationPermissionLevel(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	timeout time.Duration) error {
	httpUrl := "v3/applications/permission-level"
	updatePath := client.Endpoint + httpUrl
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"project_id":       d.Get("project_id"),
			"permission_level": d.Get("permission_level"),
			"application_ids":  []string{d.Id()},
		},
	}

	err := resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		_, err := client.Request("PUT", updatePath, &updateOpt)
		isRetry, err := handleDeployApplicationPermissionLevelOperationError(err)
		if isRetry {
			// lintignore:R018
			time.Sleep(10 * time.Second)
			return resource.RetryableError(err)
		}
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error updating CodeArts deploy application permission level: %s", err)
	}

	return nil
}

func buildCreateDeployApplicationBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"project_id":       d.Get("project_id"),
		"name":             d.Get("name"),
		"description":      d.Get("description"),
		"is_draft":         d.Get("is_draft"),
		"create_type":      d.Get("create_type"),
		"slave_cluster_id": d.Get("resource_pool_id"),
		"group_id":         utils.ValueIgnoreEmpty(d.Get("group_id")),
		"trigger":          buildTriggerBodyParam(d),
		"arrange_infos":    buildArrangeInfoBodyParam(d),
	}
	return bodyParams
}

func buildTriggerBodyParam(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"trigger_source":         d.Get("trigger_source"),
		"artifact_source_system": d.Get("artifact_source_system"),
		"artifact_type":          d.Get("artifact_type"),
	}
}

func buildArrangeInfoBodyParam(d *schema.ResourceData) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"id":             utils.ValueIgnoreEmpty(d.Get("task_id")),
			"template_id":    d.Get("template_id"),
			"operation_list": buildArrangeInfoOperationListBodyParam(d),
		},
	}
}

func buildArrangeInfoOperationListBodyParam(d *schema.ResourceData) []map[string]interface{} {
	if rawArray, ok := d.Get("operation_list").([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]map[string]interface{}, 0, len(rawArray))
		for _, v := range rawArray {
			if raw, isMap := v.(map[string]interface{}); isMap {
				rst = append(rst, map[string]interface{}{
					"name":        raw["name"],
					"description": raw["description"],
					"code":        raw["code"],
					"params":      raw["params"],
					"entrance":    raw["entrance"],
					"version":     raw["version"],
					"module_id":   raw["module_id"],
				})
			}
		}
		return rst
	}
	return nil
}

func resourceDeployApplicationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		mErr   *multierror.Error
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("codearts_deploy", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts deploy client: %s", err)
	}

	resultRespBody, err := getDeployApplication(client, d)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CodeArts deploy application")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("project_id", utils.PathSearch("project_id", resultRespBody, nil)),
		d.Set("name", utils.PathSearch("name", resultRespBody, nil)),
		d.Set("description", utils.PathSearch("description", resultRespBody, nil)),
		d.Set("create_type", utils.PathSearch("create_type", resultRespBody, nil)),
		d.Set("resource_pool_id", utils.PathSearch("slave_cluster_id", resultRespBody, nil)),
		d.Set("created_at", utils.PathSearch("create_time", resultRespBody, nil)),
		d.Set("updated_at", utils.PathSearch("update_time", resultRespBody, nil)),
		d.Set("project_name", utils.PathSearch("project_name", resultRespBody, nil)),
		d.Set("can_disable", utils.PathSearch("can_disable", resultRespBody, nil)),
		d.Set("can_modify", utils.PathSearch("can_modify", resultRespBody, nil)),
		d.Set("can_delete", utils.PathSearch("can_delete", resultRespBody, nil)),
		d.Set("can_view", utils.PathSearch("can_view", resultRespBody, nil)),
		d.Set("can_execute", utils.PathSearch("can_execute", resultRespBody, nil)),
		d.Set("can_copy", utils.PathSearch("can_copy", resultRespBody, nil)),
		d.Set("can_manage", utils.PathSearch("can_manage", resultRespBody, nil)),
		d.Set("can_create_env", utils.PathSearch("can_create_env", resultRespBody, nil)),
		d.Set("task_id", utils.PathSearch("arrange_infos|[0].id", resultRespBody, nil)),
		d.Set("task_name", utils.PathSearch("arrange_infos|[0].name", resultRespBody, nil)),
		d.Set("steps", flattenDeployApplicationSteps(resultRespBody)),
		d.Set("is_disable", utils.PathSearch("is_disable", resultRespBody, nil)),
		d.Set("is_care", utils.PathSearch("is_care", resultRespBody, nil)),
		d.Set("permission_level", utils.PathSearch("permission_level", resultRespBody, nil)),
	)

	permissionMatrix, err := getDeployApplicationPermissionMatrix(client, d)
	if err != nil {
		log.Printf("[WARN] failed to retrieve application permission matrix: %s", err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("permission_matrix", flattenDeployApplicationPermissionMatrix(permissionMatrix)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getDeployApplication(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	httpUrl := "v1/applications/{app_id}/info"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{app_id}", d.Id())
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf-8",
		},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	resultRespBody := utils.PathSearch("result", getRespBody, nil)
	if resultRespBody == nil {
		return nil, fmt.Errorf("error retrieving CodeArts deploy application: result is not found in API response")
	}

	return resultRespBody, nil
}

func getDeployApplicationPermissionMatrix(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	permissionLevel := d.Get("permission_level").(string)

	httpUrl := "v3/applications/permissions"
	getPath := client.Endpoint + httpUrl
	if permissionLevel == "instance" {
		getPath += fmt.Sprintf("?app_id=%s", d.Id())
	} else {
		getPath += fmt.Sprintf("?project_id=%s", d.Get("project_id").(string))
	}
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	permissionMatrix := utils.PathSearch("result", getRespBody, make([]interface{}, 0)).([]interface{})
	if len(permissionMatrix) == 0 {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Body: []byte("error retrieving CodeArts deploy application permission matrix, empty list"),
			},
		}
	}

	return permissionMatrix, nil
}

// flattenDeployApplicationSteps use to flatten deployment steps.
// An example of the return value of this function is as follows: '{"step1":"XXX", "step2":"XXX"}'
func flattenDeployApplicationSteps(resp interface{}) interface{} {
	steps := utils.PathSearch("arrange_infos|[0].steps", resp, nil)
	if steps == nil {
		return nil
	}
	rst := make(map[string]interface{})
	if stepMap, ok := steps.(map[string]interface{}); ok {
		for key, val := range stepMap {
			rst[key] = utils.PathSearch("name", val, "")
		}
	}
	return rst
}

func flattenDeployApplicationPermissionMatrix(respBody interface{}) []interface{} {
	if resp, isList := respBody.([]interface{}); isList {
		rst := make([]interface{}, 0, len(resp))
		for _, v := range resp {
			rst = append(rst, map[string]interface{}{
				"role_id":        utils.PathSearch("role_id", v, nil),
				"role_name":      utils.PathSearch("name", v, nil),
				"role_type":      utils.PathSearch("role_type", v, nil),
				"can_modify":     utils.PathSearch("can_modify", v, nil),
				"can_disable":    utils.PathSearch("can_disable", v, nil),
				"can_delete":     utils.PathSearch("can_delete", v, nil),
				"can_view":       utils.PathSearch("can_view", v, nil),
				"can_execute":    utils.PathSearch("can_execute", v, nil),
				"can_copy":       utils.PathSearch("can_copy", v, nil),
				"can_manage":     utils.PathSearch("can_manage", v, nil),
				"can_create_env": utils.PathSearch("can_create_env", v, nil),
			})
		}
		return rst
	}

	return nil
}

func resourceDeployApplicationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_deploy", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts deploy client: %s", err)
	}

	changes := []string{
		"name",
		"is_draft",
		"trigger_source",
		"artifact_source_system",
		"artifact_type",
		"template_id",
		"operation_list",
		"description",
		"resource_pool_id",
	}

	if d.HasChanges(changes...) {
		httpUrl := "v1/applications"
		updatePath := client.Endpoint + httpUrl
		updateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders: map[string]string{
				"Content-Type": "application/json;charset=utf-8",
			},
			JSONBody: utils.RemoveNil(buildUpdateDeployApplicationBodyParams(d)),
		}

		_, err := client.Request("PUT", updatePath, &updateOpt)
		if err != nil {
			return diag.Errorf("error updating CodeArts deploy application: %s", err)
		}
	}

	if d.HasChange("group_id") {
		if err := updateDeployApplicationGroupId(client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("permission_level") {
		if err := updateDeployApplicationPermissionLevel(ctx, client, d, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("is_disable") {
		if err := updateDeployApplicationDisable(client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceDeployApplicationRead(ctx, d, meta)
}

func buildUpdateDeployApplicationBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"id":               d.Id(),
		"project_id":       d.Get("project_id"),
		"name":             d.Get("name"),
		"description":      d.Get("description"),
		"is_draft":         d.Get("is_draft"),
		"create_type":      d.Get("create_type"),
		"slave_cluster_id": d.Get("resource_pool_id"),
		"trigger":          buildTriggerBodyParam(d),
		"arrange_infos":    buildArrangeInfoBodyParam(d),
	}
	return bodyParams
}

func updateDeployApplicationGroupId(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	httpUrl := "v1/projects/{project_id}/applications/groups/move"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", d.Get("project_id").(string))
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"group_id":        d.Get("group_id"),
			"application_ids": []string{d.Id()},
		},
	}

	updateResp, err := client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("error updating group ID for application: %s", err)
	}
	updateRespBody, err := utils.FlattenResponse(updateResp)
	if err != nil {
		return err
	}

	resultRespBody := utils.PathSearch("result", updateRespBody, make([]interface{}, 0)).([]interface{})
	if len(resultRespBody) > 0 {
		return fmt.Errorf("error updating group ID for application: %v", resultRespBody)
	}

	return nil
}

func resourceDeployApplicationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/applications/{app_id}"
		product = "codearts_deploy"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CodeArts deploy client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{app_id}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf-8",
		},
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting CodeArts deploy application")
	}

	return nil
}

// Deploy.00060222: Some application permissions in this project are being updated. Try again later.
// Error will occur when two or more applications using instance level permission are creating.
func handleDeployApplicationPermissionLevelOperationError(err error) (bool, error) {
	if err == nil {
		return false, nil
	}

	if errCode, ok := err.(golangsdk.ErrDefault400); ok {
		var apiError interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &apiError); jsonErr != nil {
			return false, jsonErr
		}

		errorCode, errorCodeErr := jmespath.Search("error_code", apiError)
		if errorCodeErr != nil {
			return false, fmt.Errorf("error parse error code from response body: %s", errorCodeErr)
		}
		if errorCode == "Deploy.00060222" {
			return true, err
		}
	}
	return false, err
}
