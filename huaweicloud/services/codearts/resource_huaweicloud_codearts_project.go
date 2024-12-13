package codearts

import (
	"context"
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

// @API CodeArts POST /v4/project
// @API CodeArts PUT /v4/projects/{id}
// @API CodeArts GET /v4/projects/{id}
// @API CodeArts DELETE /v4/projects/{id}
func ResourceProject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProjectCreate,
		UpdateContext: resourceProjectUpdate,
		ReadContext:   resourceProjectRead,
		DeleteContext: resourceProjectDelete,
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
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The project name.`,
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The type of project.`,
				ValidateFunc: validation.StringInSlice([]string{
					"scrum", "xboard", "basic", "phoenix",
				}, false),
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The description about the project.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The enterprise project ID of the project.`,
			},
			"source": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The source of project.`,
			},
			"template_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The template id which used to create project.`,
			},
			"archive": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Whether the project is archived.`,
			},
			"project_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The project code.`,
			},
			"project_num_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number id of project.`,
			},
		},
	}
}

func resourceProjectCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createProject: create a Project
	var (
		createProjectHttpUrl = "v4/project"
		createProjectProduct = "projectman"
	)
	createProjectClient, err := cfg.NewServiceClient(createProjectProduct, region)
	if err != nil {
		return diag.Errorf("error creating Project Client: %s", err)
	}

	createProjectPath := createProjectClient.Endpoint + createProjectHttpUrl

	createProjectOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	createProjectOpt.JSONBody = utils.RemoveNil(buildCreateProjectBodyParams(d, cfg))
	createProjectResp, err := createProjectClient.Request("POST", createProjectPath, &createProjectOpt)
	if err != nil {
		return diag.Errorf("error creating Project: %s", err)
	}

	createProjectRespBody, err := utils.FlattenResponse(createProjectResp)
	if err != nil {
		return diag.FromErr(err)
	}

	projectId := utils.PathSearch("project_id", createProjectRespBody, "").(string)
	if projectId == "" {
		return diag.Errorf("unable to find the CodeArts project ID from the API response")
	}
	d.SetId(projectId)

	return resourceProjectRead(ctx, d, meta)
}

func buildCreateProjectBodyParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"project_name":  utils.ValueIgnoreEmpty(d.Get("name")),
		"description":   utils.ValueIgnoreEmpty(d.Get("description")),
		"enterprise_id": utils.ValueIgnoreEmpty(cfg.GetEnterpriseProjectID(d)),
		"project_type":  utils.ValueIgnoreEmpty(d.Get("type")),
		"source":        utils.ValueIgnoreEmpty(d.Get("source")),
		"template_id":   utils.ValueIgnoreEmpty(d.Get("template_id")),
	}
	return bodyParams
}

func resourceProjectRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getProject: Query the Project
	var (
		getProjectHttpUrl = "v4/projects/{id}"
		getProjectProduct = "projectman"
	)
	getProjectClient, err := cfg.NewServiceClient(getProjectProduct, region)
	if err != nil {
		return diag.Errorf("error creating Project Client: %s", err)
	}

	getProjectPath := getProjectClient.Endpoint + getProjectHttpUrl
	getProjectPath = strings.ReplaceAll(getProjectPath, "{id}", d.Id())

	getProjectOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getProjectResp, err := getProjectClient.Request("GET", getProjectPath, &getProjectOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Project")
	}

	getProjectRespBody, err := utils.FlattenResponse(getProjectResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("archive", utils.PathSearch("project.archive", getProjectRespBody, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("project.enterprise_id", getProjectRespBody, nil)),
		d.Set("name", utils.PathSearch("project.name", getProjectRespBody, nil)),
		d.Set("project_code", utils.PathSearch("project.project_code", getProjectRespBody, nil)),
		d.Set("project_num_id", utils.PathSearch("project.project_num_id", getProjectRespBody, nil)),
		d.Set("type", utils.PathSearch("project.project_type", getProjectRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceProjectUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateProjecthasChanges := []string{
		"name",
		"description",
	}

	if d.HasChanges(updateProjecthasChanges...) {
		// updateProject: update the Project
		var (
			updateProjectHttpUrl = "v4/projects/{id}"
			updateProjectProduct = "projectman"
		)
		updateProjectClient, err := cfg.NewServiceClient(updateProjectProduct, region)
		if err != nil {
			return diag.Errorf("error creating Project Client: %s", err)
		}

		updateProjectPath := updateProjectClient.Endpoint + updateProjectHttpUrl
		updateProjectPath = strings.ReplaceAll(updateProjectPath, "{id}", d.Id())

		updateProjectOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				204,
			},
		}
		updateProjectOpt.JSONBody = utils.RemoveNil(buildUpdateProjectBodyParams(d))
		_, err = updateProjectClient.Request("PUT", updateProjectPath, &updateProjectOpt)
		if err != nil {
			return diag.Errorf("error updating Project: %s", err)
		}
	}
	return resourceProjectRead(ctx, d, meta)
}

func buildUpdateProjectBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"project_name": utils.ValueIgnoreEmpty(d.Get("name")),
		"description":  utils.ValueIgnoreEmpty(d.Get("description")),
	}
	return bodyParams
}

func resourceProjectDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteProject: missing operation notes
	var (
		deleteProjectHttpUrl = "v4/projects/{id}"
		deleteProjectProduct = "projectman"
	)
	deleteProjectClient, err := cfg.NewServiceClient(deleteProjectProduct, region)
	if err != nil {
		return diag.Errorf("error creating Project Client: %s", err)
	}

	deleteProjectPath := deleteProjectClient.Endpoint + deleteProjectHttpUrl
	deleteProjectPath = strings.ReplaceAll(deleteProjectPath, "{id}", d.Id())

	deleteProjectOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}
	_, err = deleteProjectClient.Request("DELETE", deleteProjectPath, &deleteProjectOpt)
	if err != nil {
		return diag.Errorf("error deleting Project: %s", err)
	}

	return nil
}
