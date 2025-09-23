package codeartspipeline

import (
	"context"
	"log"
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

var microServiceNonUpdatableParams = []string{
	"project_id", "name", "type", "parent_id",
}

var (
	followMicroServiceHttpUrl      = "v2/{cloudProjectId}/component/{componentId}/follow"
	unfollowMicroServiceHttpUrl    = "v2/{cloudProjectId}/component/{componentId}/unfollow"
	updateMicroServiceHttpUrl      = "v2/{cloudProjectId}/component/{componentId}/update"
	updateMicroServiceReposHttpUrl = "v2/{cloudProjectId}/component/{componentId}/repo/update"
)

// @API CodeArtsPipeline POST /v2/{cloudProjectId}/component/create
// @API CodeArtsPipeline GET /v2/{cloudProjectId}/component/{componentId}/query
// @API CodeArtsPipeline PUT /v2/{cloudProjectId}/component/{componentId}/update
// @API CodeArtsPipeline DELETE /v2/{cloudProjectId}/component/{componentId}/delete
// @API CodeArtsPipeline PUT /v2/{cloudProjectId}/component/{componentId}/follow
// @API CodeArtsPipeline PUT /v2/{cloudProjectId}/component/{componentId}/unfollow
// @API CodeArtsPipeline GET /v2/{cloudProjectId}/component/{componentId}/follow/query
// @API CodeArtsPipeline PUT /v2/{cloudProjectId}/component/{componentId}/repo/update
func ResourceCodeArtsPipelineMicroService() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePipelineMicroServiceCreate,
		ReadContext:   resourcePipelineMicroServiceRead,
		UpdateContext: resourcePipelineMicroServiceUpdate,
		DeleteContext: resourcePipelineMicroServiceDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceImportStateFuncWithProjectIdAndId,
		},

		CustomizeDiff: config.FlexibleForceNew(microServiceNonUpdatableParams),

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
				Description: `Specifies the micro service type.`,
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the micro service name.`,
			},
			"repos": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: `Specifies the repository information.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Specifies the repository type.`,
						},
						"repo_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Specifies the repository ID.`,
						},
						"http_url": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Specifies the HTTP address of the Git repository.`,
						},
						"git_url": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Specifies the Git address of the Git repository.`,
						},
						"branch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Specifies the branch.`,
						},
						"language": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Specifies the language.`,
						},
						"endpoint_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Specifies the endpoint ID.`,
						},
					},
				},
			},
			"parent_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the micro service parent ID.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the micro service description.`,
			},
			"is_followed": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Specifies whether the micro service is followed.`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the micro service status.`,
			},
			"creator_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the creator ID.`,
			},
			"updater_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the updater ID.`,
			},
			"creator_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the creator name.`,
			},
			"updater_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the updater name.`,
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the create time.`,
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the update time.`,
			},
		},
	}
}

func resourcePipelineMicroServiceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	httpUrl := "v2/{cloudProjectId}/component/create"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{cloudProjectId}", d.Get("project_id").(string))
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreatePipelineMicroServiceBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline micro service: %s", err)
	}
	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := checkResponseError(createRespBody, ""); err != nil {
		return diag.Errorf("error creating CodeArts Pipeline micro service: %s", err)
	}

	id := utils.PathSearch("id", createRespBody, "").(string)
	if id == "" {
		return diag.Errorf("unable to find the CodeArts Pipeline micro service ID from the API response")
	}

	d.SetId(id)

	if d.Get("is_followed").(bool) {
		if err := updatePipelineMicroServiceField(client, d, followMicroServiceHttpUrl, nil); err != nil {
			return diag.Errorf("error updating micro service is_followed: %s", err)
		}
	}

	return resourcePipelineMicroServiceRead(ctx, d, meta)
}

func buildCreatePipelineMicroServiceBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":      d.Get("name"),
		"type":      d.Get("type"),
		"desc":      d.Get("description"),
		"repos":     buildPipelineMicroServiceRepos(d),
		"parent_id": utils.ValueIgnoreEmpty(d.Get("parent_id")),
	}

	return bodyParams
}

func buildPipelineMicroServiceRepos(d *schema.ResourceData) interface{} {
	rawRepos := d.Get("repos").(*schema.Set).List()
	if len(rawRepos) == 0 {
		return make([]interface{}, 0)
	}

	repos := make([]map[string]interface{}, 0, len(rawRepos))
	for _, v := range rawRepos {
		if repo, ok := v.(map[string]interface{}); ok {
			customVar := map[string]interface{}{
				"type":        repo["type"],
				"repo_id":     repo["repo_id"],
				"http_url":    repo["http_url"],
				"git_url":     repo["git_url"],
				"branch":      repo["branch"],
				"language":    repo["language"],
				"endpoint_id": utils.ValueIgnoreEmpty(repo["endpoint_id"]),
			}
			repos = append(repos, customVar)
		}
	}

	return repos
}

func updatePipelineMicroServiceField(client *golangsdk.ServiceClient, d *schema.ResourceData, httpUrl string,
	bodyParams interface{}) error {
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{cloudProjectId}", d.Get("project_id").(string))
	updatePath = strings.ReplaceAll(updatePath, "{componentId}", d.Id())
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	if bodyParams != nil {
		updateOpt.JSONBody = bodyParams
	}

	updateResp, err := client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return err
	}

	updateRespBody, err := utils.FlattenResponse(updateResp)
	if err != nil {
		// when update action is completely success, will return uuid in invalid json format
		log.Printf("[WARN] error flattening response: %s", err)
	} else {
		// when input params are not correct, will return 200 but error code in response body
		if err := checkResponseError(updateRespBody, ""); err != nil {
			return err
		}
	}

	return nil
}

func resourcePipelineMicroServiceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	getRespBody, err := GetPipelineMicroService(client, d.Get("project_id").(string), d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", projectNotFoundError3),
			"error retrieving CodeArts Pipeline micro service")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("project_id", utils.PathSearch("cloud_project_id", getRespBody, nil)),
		d.Set("name", utils.PathSearch("name", getRespBody, nil)),
		d.Set("description", utils.PathSearch("description", getRespBody, nil)),
		d.Set("type", utils.PathSearch("type", getRespBody, nil)),
		d.Set("parent_id", utils.PathSearch("parent_id", getRespBody, nil)),
		d.Set("repos", flattenPipelineMicroServiceRepos(getRespBody)),
		d.Set("creator_id", utils.PathSearch("creator_id", getRespBody, nil)),
		d.Set("updater_id", utils.PathSearch("updater_id", getRespBody, nil)),
		d.Set("creator_name", utils.PathSearch("creator_name", getRespBody, nil)),
		d.Set("updater_name", utils.PathSearch("updater_name", getRespBody, nil)),
		d.Set("create_time", utils.PathSearch("create_time", getRespBody, nil)),
		d.Set("update_time", utils.PathSearch("update_time", getRespBody, nil)),
		d.Set("status", utils.PathSearch("status", getRespBody, nil)),
	)

	if followed, err := getPipelineMicroServiceIsFollowed(client, d.Get("project_id").(string), d.Id()); err != nil {
		log.Println("[WARN] error retrieving micro service is_followed")
	} else {
		mErr = multierror.Append(mErr, d.Set("is_followed", followed))
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func GetPipelineMicroService(client *golangsdk.ServiceClient, projectId, id string) (interface{}, error) {
	httpUrl := "v2/{cloudProjectId}/component/{componentId}/query"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{cloudProjectId}", projectId)
	getPath = strings.ReplaceAll(getPath, "{componentId}", id)
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

	return getRespBody, nil
}

func flattenPipelineMicroServiceRepos(resp interface{}) []interface{} {
	reposList, ok := utils.PathSearch("repos", resp, make([]interface{}, 0)).([]interface{})
	if ok && len(reposList) > 0 {
		result := make([]interface{}, 0, len(reposList))
		for _, v := range reposList {
			repo := v.(map[string]interface{})
			m := map[string]interface{}{
				"type":        utils.PathSearch("type", repo, nil),
				"repo_id":     utils.PathSearch("repo_id", repo, nil),
				"http_url":    utils.PathSearch("http_url", repo, nil),
				"git_url":     utils.PathSearch("git_url", repo, nil),
				"branch":      utils.PathSearch("branch", repo, nil),
				"language":    utils.PathSearch("language", repo, nil),
				"endpoint_id": utils.PathSearch("endpoint_id", repo, nil),
			}
			result = append(result, m)
		}
		return result
	}

	return nil
}

func getPipelineMicroServiceIsFollowed(client *golangsdk.ServiceClient, projectId, id string) (interface{}, error) {
	httpUrl := "v2/{cloudProjectId}/component/{componentId}/follow/query"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{cloudProjectId}", projectId)
	getPath = strings.ReplaceAll(getPath, "{componentId}", id)
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

	return getRespBody, nil
}

func resourcePipelineMicroServiceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	if d.HasChanges("description") {
		err := updatePipelineMicroServiceField(client, d, updateMicroServiceHttpUrl, buildUpdatePipelineMicroServiceBodyParams(d))
		if err != nil {
			return diag.Errorf("error updating micro service: %s", err)
		}
	}

	if d.HasChanges("repos") {
		err := updatePipelineMicroServiceField(client, d, updateMicroServiceReposHttpUrl, buildPipelineMicroServiceRepos(d))
		if err != nil {
			return diag.Errorf("error updating micro service repos: %s", err)
		}
	}

	if d.HasChanges("is_followed") {
		httpUrl := followMicroServiceHttpUrl
		if !d.Get("is_followed").(bool) {
			httpUrl = unfollowMicroServiceHttpUrl
		}
		err := updatePipelineMicroServiceField(client, d, httpUrl, nil)
		if err != nil {
			return diag.Errorf("error updating micro service is_followed: %s", err)
		}
	}

	return resourcePipelineMicroServiceRead(ctx, d, meta)
}

func buildUpdatePipelineMicroServiceBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"desc": d.Get("description"),
	}

	return bodyParams
}

func resourcePipelineMicroServiceDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	projectId := d.Get("project_id").(string)
	httpUrl := "v2/{cloudProjectId}/component/{componentId}/delete"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{cloudProjectId}", projectId)
	deletePath = strings.ReplaceAll(deletePath, "{componentId}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", projectNotFoundError3),
			"error deleting CodeArts Pipeline micro service")
	}

	return nil
}
