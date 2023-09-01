// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product CodeArts
// ---------------------------------------------------------------

package codearts

import (
	"context"
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
	applicationNotFound = "Deploy.00011020"
)

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
			"project_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the project name for CodeArts service.`,
			},
			// The value of this field is inconsistent between creation and query
			"template_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the deployment template ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the application name.`,
			},
			"resource_pool_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the resource pool ID.`,
			},
			// There is no value for this field in the details API
			"configs": {
				Type:        schema.TypeList,
				Elem:        deployApplicationConfigSchema(),
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the deployment parameters.`,
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
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The application state.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description.`,
			},
			"owner_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The user name of the application creator.`,
			},
			"owner_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the application creator.`,
			},
			"can_modify": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether the user has the modify permission.`,
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
			"role_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The role ID.`,
			},
		},
	}
}

func deployApplicationConfigSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the deployment parameter name, which can be customized.`,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Description: `Specifies the deployment parameter type, valid values are: **text**, **host_group**,
**enum** and **encrypt**.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the deployment parameter description.`,
			},
			"value": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the deployment parameter value.`,
			},
			"static_status": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default:  1,
				Description: `Specifies whether the parameter is a static parameter. If the value is **1**,
the parameter cannot be changed during deployment. If the value is **0**, the parameter can be changed and reported
to the pipeline.`,
			},
			"limits": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the enum values.`,
			},
		},
	}
	return &sc
}

func resourceDeployApplicationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/tasks/template-task"
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

	id, err := jmespath.Search("task_id", createRespBody)
	if err != nil || id == nil {
		return diag.Errorf("error creating CodeArts deploy application: ID is not found in API response")
	}
	d.SetId(id.(string))

	return resourceDeployApplicationRead(ctx, d, meta)
}

func buildCreateDeployApplicationBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"project_id":       d.Get("project_id"),
		"project_name":     d.Get("project_name"),
		"template_id":      d.Get("template_id"),
		"task_name":        d.Get("name"),
		"slave_cluster_id": d.Get("resource_pool_id"),
		"configs":          buildCreateApplicationConfigs(d.Get("configs")),
	}
	return bodyParams
}

func buildCreateApplicationConfigs(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]map[string]interface{}, 0, len(rawArray))
		for _, v := range rawArray {
			if raw, rawOk := v.(map[string]interface{}); rawOk {
				rst = append(rst, map[string]interface{}{
					"name":          utils.ValueIngoreEmpty(raw["name"]),
					"type":          utils.ValueIngoreEmpty(raw["type"]),
					"description":   utils.ValueIngoreEmpty(raw["description"]),
					"value":         utils.ValueIngoreEmpty(raw["value"]),
					"static_status": raw["static_status"],
					"limits":        buildApplicationConfigLimits(raw["limits"]),
				})
			}
		}
		return rst
	}
	return nil
}

func buildApplicationConfigLimits(rawLimits interface{}) []map[string]interface{} {
	if arr, ok := rawLimits.([]interface{}); ok {
		limitRaw := make([]map[string]interface{}, len(arr))
		for i, name := range arr {
			limitRaw[i] = map[string]interface{}{
				"name": name,
			}
		}
		return limitRaw
	}
	return nil
}

func resourceDeployApplicationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		mErr    *multierror.Error
		httpUrl = "v2/tasks/{task_id}"
		product = "codearts_deploy"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CodeArts deploy client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{task_id}", d.Id())
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
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

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("project_id", utils.PathSearch("project_id", getRespBody, nil)),
		d.Set("project_name", utils.PathSearch("project_name", getRespBody, nil)),
		d.Set("name", utils.PathSearch("name", getRespBody, nil)),
		d.Set("resource_pool_id", utils.PathSearch("slave_cluster_id", getRespBody, nil)),
		d.Set("created_at", utils.PathSearch("create_time", getRespBody, nil)),
		d.Set("updated_at", utils.PathSearch("update_time", getRespBody, nil)),
		d.Set("state", utils.PathSearch("state", getRespBody, nil)),
		d.Set("description", utils.PathSearch("description", getRespBody, nil)),
		d.Set("owner_name", utils.PathSearch("owner", getRespBody, nil)),
		d.Set("owner_id", utils.PathSearch("owner_id", getRespBody, nil)),
		d.Set("can_modify", utils.PathSearch("can_modify", getRespBody, nil)),
		d.Set("can_delete", utils.PathSearch("can_delete", getRespBody, nil)),
		d.Set("can_view", utils.PathSearch("can_view", getRespBody, nil)),
		d.Set("can_execute", utils.PathSearch("can_execute", getRespBody, nil)),
		d.Set("can_copy", utils.PathSearch("can_copy", getRespBody, nil)),
		d.Set("can_manage", utils.PathSearch("can_manage", getRespBody, nil)),
		d.Set("role_id", utils.PathSearch("role_id", getRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDeployApplicationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/tasks/{task_id}"
		product = "codearts_deploy"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CodeArts deploy client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{task_id}", d.Id())
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
