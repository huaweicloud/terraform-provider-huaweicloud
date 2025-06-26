package codearts

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CodeArts POST /v1/repositories
// @API CodeArts GET /v1/repositories/{repository_uuid}/status
// @API CodeArts GET /v2/repositories/{repository_uuid}
// @API CodeArts DELETE /v1/repositories/{repository_uuid}
func ResourceRepository() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRepositoryCreate,
		ReadContext:   resourceRepositoryRead,
		DeleteContext: resourceRepositoryDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
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
				ForceNew:    true,
				Description: `The repository name.`,
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The project ID for CodeHub service.`,
			},
			"visibility_level": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  `The visibility level.`,
				ValidateFunc: validation.IntInSlice([]int{0, 20}),
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The repository description.`,
			},
			"import_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The HTTPS address of the template repository encrypted using Base64.`,
			},
			"gitignore_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The program language type for generating .gitignore files.`,
			},
			"license_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Default:     1,
				Description: `The license ID for public repository.`,
			},
			"enable_readme": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Description:  `Whether to generate the README.md file.`,
				Default:      1,
				ValidateFunc: validation.IntInSlice([]int{0, 1}),
			},
			"import_members": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Description:  `Whether to import the project members.`,
				Default:      1,
				ValidateFunc: validation.IntInSlice([]int{0, 1}),
			},
			"https_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The HTTPS URL that used to the fork repository.`,
			},
			"ssh_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The SSH URL that used to the fork repository.`,
			},
			"web_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The Web URL, accessing this URL will redirect to the repository detail page.`,
			},
			"lfs_size": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The LFS capacity, in MB. If the capacity is greater than 1024M, the unit is GB.`,
			},
			"capacity": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The total size of the repository, in MB. If the capacity is greater than 1024M, the unit is GB.`,
			},
			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The repository status.`,
			},
			"create_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time.`,
			},
			"update_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The last update time.`,
			},
			"repository_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The repository primart key ID.`,
			},
		},
	}
}

func waitForRepositoryActive(ctx context.Context, cfg *config.Config, d *schema.ResourceData) error {
	var (
		getRepositoryHttpUrl = "v1/repositories/{repository_uuid}/status"
		getRepositoryProduct = "codehub"
	)

	region := cfg.GetRegion(d)
	getRepositoryClient, err := cfg.NewServiceClient(getRepositoryProduct, region)
	if err != nil {
		return fmt.Errorf("error creating repository client: %s", err)
	}

	getRepositoryPath := getRepositoryClient.Endpoint + getRepositoryHttpUrl
	getRepositoryPath = strings.ReplaceAll(getRepositoryPath, "{repository_uuid}", d.Id())

	createRepositoryOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	stateConf := &resource.StateChangeConf{
		Target:       []string{"ACTIVE"},
		Refresh:      repositoryRefreshFunc(getRepositoryClient, getRepositoryPath, createRepositoryOpt),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
		// We can't query the repository after it becomes ACTIVE immediately
		ContinuousTargetOccurence: 2,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	return err
}

func repositoryRefreshFunc(client *golangsdk.ServiceClient, path string,
	opts golangsdk.RequestOpts) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.Request("GET", path, &opts)
		if err != nil {
			return nil, "ERROR", err
		}
		respBody, err := utils.FlattenResponse(resp)
		if err == nil && respBody != nil {
			status := utils.PathSearch("status", respBody, "").(string)
			if status == "success" {
				return resp, "ACTIVE", nil
			}
		}

		return nil, "ERROR", err
	}
}

func resourceRepositoryCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// createRepository: Create a CodeHub repository
	var (
		createRepositoryHttpUrl = "v1/repositories"
		createRepositoryProduct = "codehub"
	)

	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	createRepositoryClient, err := cfg.NewServiceClient(createRepositoryProduct, region)
	if err != nil {
		return diag.Errorf("error creating repository client: %s", err)
	}

	createRepositoryPath := createRepositoryClient.Endpoint + createRepositoryHttpUrl

	createRepositoryOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	createRepositoryOpt.JSONBody = buildCreateRepositoryBodyParams(d)
	createRepositoryResp, err := createRepositoryClient.Request("POST", createRepositoryPath, &createRepositoryOpt)
	if err != nil {
		return diag.Errorf("error creating CodeHub repository: %s", err)
	}

	createRepositoryRespBody, err := utils.FlattenResponse(createRepositoryResp)
	if err != nil {
		return diag.FromErr(err)
	}
	repositoryId := utils.PathSearch("result.repository_uuid", createRepositoryRespBody, "").(string)
	if repositoryId == "" {
		return diag.Errorf("unable to find the CodeHub repository ID from the API response")
	}
	d.SetId(repositoryId)

	if err = waitForRepositoryActive(ctx, cfg, d); err != nil {
		return diag.Errorf("timout waiting for CodeHub repository to become active: %s", err)
	}
	return resourceRepositoryRead(ctx, d, meta)
}

func buildCreateRepositoryBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":             utils.ValueIgnoreEmpty(d.Get("name")),
		"project_uuid":     utils.ValueIgnoreEmpty(d.Get("project_id")),
		"visibility_level": d.Get("visibility_level"),
		"description":      utils.ValueIgnoreEmpty(d.Get("description")),
		"import_url":       utils.ValueIgnoreEmpty(d.Get("import_url")),
		"gitignore_id":     utils.ValueIgnoreEmpty(d.Get("gitignore_id")),
		"license_id":       utils.ValueIgnoreEmpty(d.Get("license_id")),
		"enable_readme":    d.Get("enable_readme"),
		"import_members":   d.Get("import_members"),
	}
	return bodyParams
}

func resourceRepositoryRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// getRepository: Query the resource detail of the CodeHub repository
	var (
		getRepositoryHttpUrl = "v2/repositories/{repository_uuid}"
		getRepositoryProduct = "codehub"
		mErr                 *multierror.Error
	)

	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	getRepositoryClient, err := conf.NewServiceClient(getRepositoryProduct, region)
	if err != nil {
		return diag.Errorf("error creating repository client: %s", err)
	}

	getRepositoryPath := getRepositoryClient.Endpoint + getRepositoryHttpUrl
	getRepositoryPath = strings.ReplaceAll(getRepositoryPath, "{repository_uuid}", d.Id())

	getRepositoryOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getRepositoryResp, err := getRepositoryClient.Request("GET", getRepositoryPath, &getRepositoryOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, parseRepositoryRequestError(err), "error retrieving CodeHub repository")
	}

	getRepositoryRespBody, err := utils.FlattenResponse(getRepositoryResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(mErr,
		d.Set("region", region),
		d.Set("https_url", utils.PathSearch("result.https_url", getRepositoryRespBody, nil)),
		d.Set("ssh_url", utils.PathSearch("result.ssh_url", getRepositoryRespBody, nil)),
		d.Set("web_url", utils.PathSearch("result.web_url", getRepositoryRespBody, nil)),
		d.Set("lfs_size", utils.PathSearch("result.lfs_size", getRepositoryRespBody, nil)),
		d.Set("project_id", utils.PathSearch("result.project_uuid", getRepositoryRespBody, nil)),
		d.Set("capacity", utils.PathSearch("result.repository_size", getRepositoryRespBody, nil)),
		d.Set("status", utils.PathSearch("result.status", getRepositoryRespBody, nil)),
		d.Set("create_at", utils.PathSearch("result.create_at", getRepositoryRespBody, nil)),
		d.Set("update_at", utils.PathSearch("result.update_at", getRepositoryRespBody, nil)),
		d.Set("visibility_level", utils.PathSearch("result.visibility_level", getRepositoryRespBody, nil)),
		d.Set("repository_id", utils.PathSearch("result.repository_id", getRepositoryRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

// The error code "CH.080401" means that the repository may have been deleted.
func parseRepositoryRequestError(respErr error) error {
	if _, ok := respErr.(golangsdk.ErrDefault401); ok {
		return golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Body: []byte("the member has been removed from the repository or the repository has been removed."),
			},
		}
	}
	return respErr
}

func resourceRepositoryDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// deleteRepository: Remove an existing CodeHub repository
	var (
		deleteRepositoryHttpUrl = "v1/repositories/{repository_uuid}"
		deleteRepositoryProduct = "codehub"
	)

	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	deleteRepositoryClient, err := conf.NewServiceClient(deleteRepositoryProduct, region)
	if err != nil {
		return diag.Errorf("error creating repository client: %s", err)
	}

	deleteRepositoryPath := deleteRepositoryClient.Endpoint + deleteRepositoryHttpUrl
	deleteRepositoryPath = strings.ReplaceAll(deleteRepositoryPath, "{repository_uuid}", d.Id())

	deleteRepositoryOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	_, err = deleteRepositoryClient.Request("DELETE", deleteRepositoryPath, &deleteRepositoryOpt)
	if err != nil {
		return diag.Errorf("error deleting CodeHub repository: %s", err)
	}

	return nil
}
