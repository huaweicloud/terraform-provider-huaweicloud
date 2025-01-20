package fgs

import (
	"context"
	"fmt"
	"log"
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

// @API FunctionGraph POST /v2/{project_id}/fgs/applications
// @API FunctionGraph GET /v2/{project_id}/fgs/applications/{id}
// @API FunctionGraph DELETE /v2/{project_id}/fgs/applications/{id}
func ResourceApplication() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApplicationCreate,
		ReadContext:   resourceApplicationRead,
		DeleteContext: resourceApplicationDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the application is located.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The application name`,
			},
			"template_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Description: utils.SchemaDesc(
					`The ID of the template used by the application.`,
					utils.SchemaDescInput{
						Required: true,
					}),
			},
			"agency_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `The agency name used by the application.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `The description of the application.`,
			},
			"params": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsJSON,
				Description:  `The template parameters, in JSON format.`,
			},
			"stack_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the stack where the application is deployed.`,
			},
			"stack_resources": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        stackResourceSchema(),
				Description: `The list of the stack resources information.`,
			},
			"repository": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        repositorySchema(),
				Description: `The repository information.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The application status.`,
			},
		},
	}
}

func stackResourceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"physical_resource_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The physical resource ID.`,
			},
			"physical_resource_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The physical resource name.`,
			},
			"logical_resource_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The logical resource name.`,
			},
			"logical_resource_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The logical resource type.`,
			},
			"resource_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of resource.`,
			},
			"status_message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status information.`,
			},
			"href": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The hyperlink.`,
			},
			"display_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The cloud service name.`,
			},
		},
	}
}

func repositorySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"https_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The HTTP address of the repository.`,
			},
			"web_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The repository link.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The repository status.`,
			},
			"project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: `The project ID of the repository.`,
			},
		},
	}
}

func buildCreateApplicationBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":        d.Get("name"),
		"template_id": d.Get("template_id"),
		"description": d.Get("description"),
		"agency_name": d.Get("agency_name"),
		"params":      utils.StringToJson(d.Get("params").(string)),
	}
}

func applicationStatusRefreshFunc(client *golangsdk.ServiceClient, appId string, targets []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		respBody, err := GetApplicationById(client, appId)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok && len(targets) < 1 {
				log.Printf("[DEBUG] The FunctionGraph application (%s) has been deleted", appId)
				// When the error code is 404, the value of respBody is nil, and a non-null value is returned to avoid continuing the loop check.
				return "Resource Not Found", "COMPLETED", nil
			}
			return respBody, "ERROR", err
		}

		status := utils.PathSearch("status", respBody, "").(string)
		unexpectedStatuses := []string{
			"CreateFail", "InitingFailed", "RegisterFailed", "InstallFailed",
			"UpdateFailed", "RollbackFailed", "UnRegisterFailed", "DeleteFailed",
		}
		if utils.StrSliceContains(unexpectedStatuses, status) {
			return respBody, "ERROR", fmt.Errorf("unexpect status (%s)", status)
		}

		if utils.StrSliceContains(targets, status) {
			return respBody, "COMPLETED", nil
		}
		return respBody, "PENDING", nil
	}
}

func waitForApplicationCreateCompleted(ctx context.Context, client *golangsdk.ServiceClient, appId string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      applicationStatusRefreshFunc(client, appId, []string{"success"}),
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceApplicationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/fgs/applications"
	)

	client, err := cfg.NewServiceClient("fgs", region)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildCreateApplicationBodyParams(d)),
	}

	requestResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph application: %s", err)
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	appId := utils.PathSearch("application_id", respBody, "").(string)
	if appId == "" {
		return diag.Errorf("unable to find the application ID from the API response")
	}
	d.SetId(appId)

	err = waitForApplicationCreateCompleted(ctx, client, appId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the application (%s) status to become success: %s", appId, err)
	}

	return resourceApplicationRead(ctx, d, meta)
}

func GetApplicationById(client *golangsdk.ServiceClient, resourceId string) (interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/fgs/applications/{id}"
	)

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{id}", resourceId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(requestResp)
}

func resourceApplicationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("fgs", region)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph client: %s", err)
	}

	respBody, err := GetApplicationById(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "FunctionGraph application")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", respBody, nil)),
		d.Set("description", utils.PathSearch("description", respBody, nil)),
		d.Set("stack_id", utils.PathSearch("stack_id", respBody, nil)),
		d.Set("stack_resources", flattenStackResources(utils.PathSearch("stack_resources",
			respBody, make([]interface{}, 0)).([]interface{}))),
		d.Set("repository", flattenRepository(utils.PathSearch("repo", respBody, nil))),
		d.Set("status", utils.PathSearch("status", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenStackResources(stackResources []interface{}) []map[string]interface{} {
	if len(stackResources) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(stackResources))
	for _, stackResource := range stackResources {
		result = append(result, map[string]interface{}{
			"physical_resource_id":   utils.PathSearch("physical_resource_id", stackResource, nil),
			"physical_resource_name": utils.PathSearch("physical_resource_name", stackResource, nil),
			"logical_resource_name":  utils.PathSearch("logical_resource_name", stackResource, nil),
			"logical_resource_type":  utils.PathSearch("logical_resource_type", stackResource, nil),
			"resource_status":        utils.PathSearch("resource_status", stackResource, nil),
			"status_message":         utils.PathSearch("status_message", stackResource, nil),
			"href":                   utils.PathSearch("href", stackResource, nil),
			"display_name":           utils.PathSearch("display_name", stackResource, nil),
		})
	}
	return result
}

func flattenRepository(repository interface{}) []map[string]interface{} {
	if repository == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"https_url":  utils.PathSearch("https_url", repository, nil),
			"web_url":    utils.PathSearch("web_url", repository, nil),
			"status":     utils.PathSearch("repo_status", repository, nil),
			"project_id": utils.PathSearch("project_id", repository, nil),
		},
	}
}

func resourceApplicationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/fgs/applications/{id}"
		appId   = d.Id()
	)

	client, err := cfg.NewServiceClient("fgs", region)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{id}", appId)
	// Due to API restrictions, the request body must pass in an empty JSON.
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error deleting FunctionGraph application (%s): %s", appId, err)
	}

	err = waitForApplicationDeleteCompleted(ctx, client, appId, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		diag.Errorf("error waiting for the application (%s) status to become deleted: %s", appId, err)
	}
	return nil
}

func waitForApplicationDeleteCompleted(ctx context.Context, client *golangsdk.ServiceClient, appId string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      applicationStatusRefreshFunc(client, appId, nil),
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}
