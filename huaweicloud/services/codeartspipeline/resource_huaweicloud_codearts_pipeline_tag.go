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

var tagNonUpdatableParams = []string{
	"project_id",
}

// @API CodeArtsPipeline POST /v5/{project_id}/api/pipeline-tag/create
// @API CodeArtsPipeline GET /v5/{project_id}/api/pipeline-tag/list
// @API CodeArtsPipeline POST /v5/{project_id}/api/pipeline-tag/update
// @API CodeArtsPipeline DELETE /v5/{project_id}/api/pipeline-tag/delete
func ResourceCodeArtsPipelineTag() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePipelineTagCreateOrUpdate,
		ReadContext:   resourcePipelineTagRead,
		UpdateContext: resourcePipelineTagCreateOrUpdate,
		DeleteContext: resourcePipelineTagDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceImportStateFuncWithProjectIdAndId,
		},

		CustomizeDiff: config.FlexibleForceNew(tagNonUpdatableParams),

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
				Description: `Specifies the tag name.`,
			},
			"color": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the tag color.`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"project_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the CodeArts project name.`,
			},
		},
	}
}

func resourcePipelineTagCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	httpUrl := "v5/{project_id}/api/pipeline-tag/create"
	if !d.IsNewResource() {
		httpUrl = "v5/{project_id}/api/pipeline-tag/update"
	}
	createOrUpdatePath := client.Endpoint + httpUrl
	createOrUpdatePath = strings.ReplaceAll(createOrUpdatePath, "{project_id}", d.Get("project_id").(string))
	createOrUpdateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreatePipelineTagBodyParams(d)),
	}

	createOrUpdateResp, err := client.Request("POST", createOrUpdatePath, &createOrUpdateOpt)
	if err != nil {
		return diag.Errorf("error setting CodeArts Pipeline tag: %s", err)
	}
	createOrUpdateRespBody, err := utils.FlattenResponse(createOrUpdateResp)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := checkResponseError(createOrUpdateRespBody, ""); err != nil {
		return diag.Errorf("error setting CodeArts Pipeline tag: %s", err)
	}

	if d.IsNewResource() {
		getRespBody, err := GetPipelineTag(client, d.Get("project_id").(string))
		if err != nil {
			return diag.Errorf("error retrieving CodeArts Pipeline tags: %s", err)
		}

		searchPath := fmt.Sprintf("[?name=='%s']|[0].tag_id", d.Get("name").(string))
		id := utils.PathSearch(searchPath, getRespBody, "").(string)
		if id == "" {
			return diag.Errorf("unable to find the CodeArts Pipeline tag ID from the API response")
		}

		d.SetId(id)
	}

	return resourcePipelineTagRead(ctx, d, meta)
}

func buildCreatePipelineTagBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"tagId": utils.ValueIgnoreEmpty(d.Id()),
		"name":  d.Get("name"),
		"color": d.Get("color"),
	}

	return bodyParams
}

func resourcePipelineTagRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	getRespBody, err := GetPipelineTag(client, d.Get("project_id").(string))
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CodeArts Pipeline tags")
	}

	searchPath := fmt.Sprintf("[?tag_id=='%s']|[0]", d.Id())
	tag := utils.PathSearch(searchPath, getRespBody, nil)
	if tag == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving CodeArts Pipeline tag")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("project_id", utils.PathSearch("project_id", tag, nil)),
		d.Set("project_name", utils.PathSearch("project_name", tag, nil)),
		d.Set("name", utils.PathSearch("name", tag, nil)),
		d.Set("color", utils.PathSearch("color", tag, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func GetPipelineTag(client *golangsdk.ServiceClient, projectId string) (interface{}, error) {
	httpUrl := "v5/{project_id}/api/pipeline-tag/list?proj_id={project_id}"
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

func resourcePipelineTagDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	projectId := d.Get("project_id").(string)
	httpUrl := "v5/{project_id}/api/pipeline-tag/delete?tagId={id}"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", projectId)
	deletePath = strings.ReplaceAll(deletePath, "{id}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         map[string]interface{}{},
	}

	deleteResp, err := client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting CodeArts Pipeline tag")
	}

	deleteRespBody, err := utils.FlattenResponse(deleteResp)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := checkResponseError(deleteRespBody, projectNotFoundError2, tagNotFoundError); err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting CodeArts Pipeline tag")
	}

	return nil
}
