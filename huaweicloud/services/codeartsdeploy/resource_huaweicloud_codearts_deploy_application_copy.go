package codeartsdeploy

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CodeArtsDeploy POST /v1/applications/{app_id}/duplicate
// @API CodeArtsDeploy GET /v1/applications/{app_id}/info
// @API CodeArtsDeploy PUT /v1/applications
// @API CodeArtsDeploy DELETE /v1/applications/{app_id}
// @API CodeArtsDeploy PUT /v1/applications/{app_id}/disable
// @API CodeArtsDeploy PUT /v3/applications/permission-level
// @API CodeArtsDeploy GET /v3/applications/permissions
// @API CodeArtsDeploy PUT /v1/projects/{project_id}/applications/groups/move
func ResourceDeployApplicationCopy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCodeArtsDeployApplicationCopyCreate,
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

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"source_app_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the application ID.`,
			},
			"project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the project ID for CodeArts service.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the application name.`,
			},
			"is_draft": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Specifies whether the application is in draft status.`,
			},
			"create_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the creation type.`,
			},
			"trigger_source": {
				Type:        schema.TypeString,
				Optional:    true,
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
			"operation_list": {
				Type:        schema.TypeList,
				Elem:        deployApplicationOperationSchema(),
				Optional:    true,
				Description: `Specifies the deployment orchestration list information.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the application description.`,
			},
			"resource_pool_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the custom slave resource pool ID.`,
			},
			"is_disable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Specifies whether to disable the application.`,
			},
			"group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the group ID.`,
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

func resourceCodeArtsDeployApplicationCopyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_deploy", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts deploy client: %s", err)
	}

	httpUrl := "v1/applications/{app_id}/duplicate"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{app_id}", d.Get("source_app_id").(string))
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error copying CodeArts deploy application: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	copyId := utils.PathSearch("result.id", createRespBody, "").(string)
	if copyId == "" {
		return diag.Errorf("unable to find the CodeArts deploy new application ID from the API response")
	}

	d.SetId(copyId)

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
		resultRespBody, err := getDeployApplication(client, d)
		if err != nil {
			return diag.FromErr(err)
		}

		taskId := utils.PathSearch("arrange_infos|[0].id", resultRespBody, "").(string)
		if taskId == "" {
			return diag.Errorf("unable to find deployment task ID")
		}

		d.Set("task_id", taskId)

		httpUrl := "v1/applications"
		updatePath := client.Endpoint + httpUrl
		updateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders: map[string]string{
				"Content-Type": "application/json;charset=utf-8",
			},
			JSONBody: utils.RemoveNil(buildUpdateDeployApplicationBodyParams(d)),
		}

		_, err = client.Request("PUT", updatePath, &updateOpt)
		if err != nil {
			return diag.Errorf("error updating CodeArts deploy application: %s", err)
		}
	}

	if _, ok := d.GetOk("group_id"); ok {
		if err := updateDeployApplicationGroupId(client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if _, ok := d.GetOk("permission_level"); ok {
		if err := updateDeployApplicationPermissionLevel(ctx, client, d, d.Timeout(schema.TimeoutCreate)); err != nil {
			return diag.FromErr(err)
		}
	}

	// it's ok to disable the app again if that is disabled already
	if err := updateDeployApplicationDisable(client, d); err != nil {
		return diag.FromErr(err)
	}

	return resourceDeployApplicationRead(ctx, d, meta)
}
