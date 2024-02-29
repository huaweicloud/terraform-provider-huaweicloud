// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product CodeArts
// ---------------------------------------------------------------

package codearts

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const (
	applicationNotFound = "Deploy.00011021"
)

// @API CodeArtsDeploy POST /v1/applications
// @API CodeArtsDeploy GET /v1/applications/{app_id}/info
// @API CodeArtsDeploy DELETE /v1/applications/{app_id}
func ResourceDeployApplication() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDeployApplicationCreate,
		ReadContext:   resourceDeployApplicationRead,
		DeleteContext: resourceDeployApplicationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
				ForceNew:    true,
				Description: `Specifies the application name.`,
			},
			"is_draft": {
				Type:        schema.TypeBool,
				Required:    true,
				ForceNew:    true,
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
				ForceNew:    true,
				Description: `Specifies where a deployment task can be executed.`,
			},
			"artifact_source_system": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the source information transferred by the pipeline.`,
			},
			"artifact_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the artifact type for the pipeline source.`,
			},
			"template_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "schema: Deprecated; Currently, the field is useless for creating API.",
			},
			"operation_list": {
				Type:        schema.TypeList,
				Elem:        deployApplicationOperationSchema(),
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the deployment orchestration list information.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the application description.`,
			},
			"resource_pool_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the custom slave resource pool ID.`,
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
				Description: `Check whether the user has the management permission, including adding, deleting,
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
		},
	}
}

func deployApplicationOperationSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the step name.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the step description.`,
			},
			"code": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the download URL.`,
			},
			"params": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the parameter.`,
			},
			"entrance": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the entry function.`,
			},
			"version": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the version.`,
			},
			"module_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the module ID.`,
			},
		},
	}
	return &sc
}

func preCheckBeforeCreateApplication(d *schema.ResourceData) error {
	rawArray, isArray := d.Get("operation_list").([]interface{})
	if isArray && len(rawArray) == 0 {
		if !d.Get("is_draft").(bool) {
			return fmt.Errorf("the argument (operation_list) is required when application is not in draft status")
		}
	}
	return nil
}

func resourceDeployApplicationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if err := preCheckBeforeCreateApplication(d); err != nil {
		return diag.FromErr(err)
	}

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

	if err := checkResponseError(createRespBody, applicationNotFound); err != nil {
		return diag.Errorf("error creating CodeArts deploy application: %s", err)
	}

	id, err := jmespath.Search("result.id", createRespBody)
	if err != nil || id == nil {
		return diag.Errorf("error creating DeployApplication: ID is not found in API response")
	}
	d.SetId(id.(string))

	return resourceDeployApplicationRead(ctx, d, meta)
}

func buildCreateDeployApplicationBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
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
		mErr    *multierror.Error
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/applications/{app_id}/info"
		product = "codearts_deploy"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CodeArts deploy client: %s", err)
	}

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
		return diag.Errorf("error retrieving CodeArts deploy application: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := checkResponseError(getRespBody, applicationNotFound); err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CodeArts deploy application")
	}

	resultRespBody := utils.PathSearch("result", getRespBody, nil)
	if resultRespBody == nil {
		return diag.Errorf("error retrieving CodeArts deploy application: result is not found in API response")
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
	)

	return diag.FromErr(mErr.ErrorOrNil())
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

	deleteResp, err := client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error deleting CodeArts deploy application: %s", err)
	}
	deleteRespBody, err := utils.FlattenResponse(deleteResp)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := checkResponseError(deleteRespBody, applicationNotFound); err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting CodeArts deploy application")
	}

	return nil
}
