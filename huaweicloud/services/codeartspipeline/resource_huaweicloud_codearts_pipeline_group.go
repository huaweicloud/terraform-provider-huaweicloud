package codeartspipeline

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var pipelineGroupNonUpdatableParams = []string{
	"project_id", "parent_id",
}

// @API CodeArtsPipeline POST /v5/{project_id}/api/pipeline-group/create
// @API CodeArtsPipeline GET /v5/{project_id}/api/pipeline-group/tree
// @API CodeArtsPipeline POST /v5/{project_id}/api/pipeline-group/update
// @API CodeArtsPipeline DELETE /v5/{project_id}/api/pipeline-group/delete
func ResourceCodeArtsPipelineGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePipelineGroupCreate,
		ReadContext:   resourcePipelineGroupRead,
		UpdateContext: resourcePipelineGroupUpdate,
		DeleteContext: resourcePipelineGroupDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceImportStateFuncWithProjectIdAndId,
		},

		CustomizeDiff: config.FlexibleForceNew(pipelineGroupNonUpdatableParams),

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
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the group name.`,
			},
			"parent_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the group parent ID.`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"path_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the path ID.`,
			},
			"ordinal": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the group ordinal`,
			},
			"creator": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the creator ID.`,
			},
			"updater": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the updater ID.`,
			},
			"create_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the create time.`,
			},
			"update_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the update time.`,
			},
			"children": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Indicates the child group name list.`,
			},
		},
	}
}

func resourcePipelineGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	httpUrl := "v5/{project_id}/api/pipeline-group/create"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", d.Get("project_id").(string))
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreatePipelineGroupBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline group: %s", err)
	}
	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := checkResponseError(createRespBody, ""); err != nil {
		return diag.Errorf("error creating CodeArts Pipeline group: %s", err)
	}

	id := utils.PathSearch("id", createRespBody, "").(string)
	if id == "" {
		return diag.Errorf("unable to find the CodeArts Pipeline group ID from the API response")
	}

	d.SetId(id)

	return resourcePipelineGroupRead(ctx, d, meta)
}

func buildCreatePipelineGroupBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":       d.Get("name"),
		"project_id": d.Get("project_id"),
		"parent_id":  utils.ValueIgnoreEmpty(d.Get("parent_id")),
	}

	return bodyParams
}

func resourcePipelineGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	groups, err := GetPipelineGroups(client, d.Get("project_id").(string))
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CodeArts Pipeline groups")
	}

	// filter group by group ID
	group := filterGroupByGroupId(groups.([]interface{}), d.Id())
	if group == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving pipeline group")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("project_id", utils.PathSearch("project_id", group, nil)),
		d.Set("name", utils.PathSearch("name", group, nil)),
		d.Set("parent_id", utils.PathSearch("parent_id", group, nil)),
		d.Set("path_id", utils.PathSearch("path_id", group, nil)),
		d.Set("ordinal", utils.PathSearch("ordinal", group, nil)),
		d.Set("creator", utils.PathSearch("creator", group, nil)),
		d.Set("updater", utils.PathSearch("updater", group, nil)),
		d.Set("create_time", utils.PathSearch("create_time", group, nil)),
		d.Set("update_time", utils.PathSearch("update_time", group, nil)),
		d.Set("children", utils.PathSearch("children[*].name", group, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func GetPipelineGroups(client *golangsdk.ServiceClient, projectId string) (interface{}, error) {
	httpUrl := "v5/{project_id}/api/pipeline-group/tree"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", projectId)
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

	if err := checkResponseError(getRespBody, projectNotFoundError2); err != nil {
		return nil, err
	}

	return getRespBody, nil
}

func filterGroupByGroupId(currentList []interface{}, id string) interface{} {
	// search first level
	searchPath := fmt.Sprintf(`[?id=='%s']|[0]`, id)
	group := utils.PathSearch(searchPath, currentList, nil)

	if group == nil {
		// search next level
		children := utils.PathSearch("[*].children[]", currentList, make([]interface{}, 0)).([]interface{})
		if len(children) != 0 {
			group = filterGroupByGroupId(children, id)
		}
	}

	return group
}

func resourcePipelineGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	httpUrl := "v5/{project_id}/api/pipeline-group/update"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", d.Get("project_id").(string))
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"name": d.Get("name"),
			"id":   d.Id(),
		},
	}

	updateResp, err := client.Request("POST", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating pipeline group: %s", err)
	}

	updateRespBody, err := utils.FlattenResponse(updateResp)
	if err != nil {
		return diag.Errorf("error flattening pipeline group response: %s", err)
	}

	if err := checkResponseError(updateRespBody, ""); err != nil {
		return diag.Errorf("error updating pipeline group: %s", err)
	}

	return resourcePipelineGroupRead(ctx, d, meta)
}

func resourcePipelineGroupDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	projectId := d.Get("project_id").(string)
	httpUrl := "v5/{project_id}/api/pipeline-group/delete?id={id}"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", projectId)
	deletePath = strings.ReplaceAll(deletePath, "{id}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	deleteResp, err := client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting CodeArts Pipeline group")
	}

	deleteRespBody, err := utils.FlattenResponse(deleteResp)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := checkResponseError(deleteRespBody, projectNotFoundError2, groupNotFoundError); err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting CodeArts Pipeline group")
	}

	return nil
}
