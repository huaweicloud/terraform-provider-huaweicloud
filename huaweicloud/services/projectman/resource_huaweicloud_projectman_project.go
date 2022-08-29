// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product ProjectMan
// ---------------------------------------------------------------

package projectman

import (
	"context"
	"strings"

	"github.com/chnsz/golangsdk"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/jmespath/go-jmespath"
)

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
				Type:         schema.TypeString,
				Required:     true,
				Description:  `The project name.`,
				ValidateFunc: validation.StringLenBetween(1, 128),
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
	config := meta.(*config.Config)
	region := config.GetRegion(d)

	// createProject: create a Project
	var (
		createProjectHttpUrl = "v4/project"
		createProjectProduct = "projectman"
	)
	createProjectClient, err := config.NewServiceClient(createProjectProduct, region)
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
	createProjectOpt.JSONBody = utils.RemoveNil(buildCreateProjectBodyParams(d, config))
	createProjectResp, err := createProjectClient.Request("POST", createProjectPath, &createProjectOpt)
	if err != nil {
		return diag.Errorf("error creating Project: %s", err)
	}

	createProjectRespBody, err := utils.FlattenResponse(createProjectResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := jmespath.Search("project_id", createProjectRespBody)
	if err != nil {
		return diag.Errorf("error creating Project: ID is not found in API response")
	}
	d.SetId(id.(string))

	return resourceProjectRead(ctx, d, meta)
}

func buildCreateProjectBodyParams(d *schema.ResourceData, config *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"project_name":  utils.ValueIngoreEmpty(d.Get("name")),
		"description":   utils.ValueIngoreEmpty(d.Get("description")),
		"enterprise_id": utils.ValueIngoreEmpty(common.GetEnterpriseProjectID(d, config)),
		"project_type":  utils.ValueIngoreEmpty(d.Get("type")),
		"source":        utils.ValueIngoreEmpty(d.Get("source")),
		"template_id":   utils.ValueIngoreEmpty(d.Get("template_id")),
	}
	return bodyParams
}

func resourceProjectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)

	var mErr *multierror.Error

	// getProject: Query the Project
	var (
		getProjectHttpUrl = "v4/projects/{id}"
		getProjectProduct = "projectman"
	)
	getProjectClient, err := config.NewServiceClient(getProjectProduct, region)
	if err != nil {
		return diag.Errorf("error creating Project Client: %s", err)
	}

	getProjectPath := getProjectClient.Endpoint + getProjectHttpUrl
	getProjectPath = strings.Replace(getProjectPath, "{id}", d.Id(), -1)

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
	config := meta.(*config.Config)
	region := config.GetRegion(d)

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
		updateProjectClient, err := config.NewServiceClient(updateProjectProduct, region)
		if err != nil {
			return diag.Errorf("error creating Project Client: %s", err)
		}

		updateProjectPath := updateProjectClient.Endpoint + updateProjectHttpUrl
		updateProjectPath = strings.Replace(updateProjectPath, "{id}", d.Id(), -1)

		updateProjectOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				204,
			},
		}
		updateProjectOpt.JSONBody = utils.RemoveNil(buildUpdateProjectBodyParams(d, config))
		_, err = updateProjectClient.Request("PUT", updateProjectPath, &updateProjectOpt)
		if err != nil {
			return diag.Errorf("error updating Project: %s", err)
		}
	}
	return resourceProjectRead(ctx, d, meta)
}

func buildUpdateProjectBodyParams(d *schema.ResourceData, config *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"project_name": utils.ValueIngoreEmpty(d.Get("name")),
		"description":  utils.ValueIngoreEmpty(d.Get("description")),
	}
	return bodyParams
}

func resourceProjectDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)

	// deleteProject: missing operation notes
	var (
		deleteProjectHttpUrl = "v4/projects/{id}"
		deleteProjectProduct = "projectman"
	)
	deleteProjectClient, err := config.NewServiceClient(deleteProjectProduct, region)
	if err != nil {
		return diag.Errorf("error creating Project Client: %s", err)
	}

	deleteProjectPath := deleteProjectClient.Endpoint + deleteProjectHttpUrl
	deleteProjectPath = strings.Replace(deleteProjectPath, "{id}", d.Id(), -1)

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
