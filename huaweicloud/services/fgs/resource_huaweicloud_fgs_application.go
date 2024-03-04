package fgs

import (
	"context"
	"encoding/json"
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
				Description:  `The template parameters.`,
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
		JSONBody:         utils.RemoveNil(buildCreateApplicationBodyParams(d)),
	}

	requestResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph application: %s", err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceId := utils.PathSearch("application_id", respBody, "")
	d.SetId(resourceId.(string))

	err = waitForApplicationStatusCompleted(ctx, client, d)
	if err != nil {
		return diag.Errorf("error waiting for the application (%s) status to become success: %s", resourceId, err)
	}

	return resourceApplicationRead(ctx, d, meta)
}

func buildCreateApplicationBodyParams(d *schema.ResourceData) map[string]interface{} {
	params := d.Get("params").(string)
	parseResult := make(map[string]interface{})
	err := json.Unmarshal([]byte(params), &parseResult)
	if err != nil {
		log.Printf("[ERROR] Invalid type of the params, not json format")
	}
	return map[string]interface{}{
		"name":        d.Get("name"),
		"template_id": d.Get("template_id"),
		"description": d.Get("description"),
		"agency_name": d.Get("agency_name"),
		"params":      parseResult,
	}
}

func waitForApplicationStatusCompleted(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      applicationStatusRefreshFunc(client, d, []string{"success"}),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func applicationStatusRefreshFunc(client *golangsdk.ServiceClient, d *schema.ResourceData, targets []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		applicationId := d.Id()
		respBody, err := getApplicationById(client, applicationId)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok && len(targets) < 1 {
				log.Printf("[DEBUG] The FunctionGraph application (%s) has been deleted", applicationId)
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

func getApplicationById(client *golangsdk.ServiceClient, resourceId string) (interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/fgs/applications/{id}"
	)

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{id}", resourceId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
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

	respBody, err := getApplicationById(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "FunctionGraph application")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", respBody, nil)),
		d.Set("description", utils.PathSearch("description", respBody, nil)),
		d.Set("stack_id", utils.PathSearch("stack_id", respBody, nil)),
		d.Set("stack_resources", flattenStackResource(respBody)),
		d.Set("repository", flattenRepository(respBody)),
		d.Set("status", utils.PathSearch("status", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenStackResource(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curJson := utils.PathSearch("stack_resources", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"physical_resource_id":   utils.PathSearch("physical_resource_id", v, nil),
			"physical_resource_name": utils.PathSearch("physical_resource_name", v, nil),
			"logical_resource_name":  utils.PathSearch("logical_resource_name", v, nil),
			"logical_resource_type":  utils.PathSearch("logical_resource_type", v, nil),
			"resource_status":        utils.PathSearch("resource_status", v, nil),
			"status_message":         utils.PathSearch("status_message", v, nil),
			"href":                   utils.PathSearch("href", v, nil),
			"display_name":           utils.PathSearch("display_name", v, nil),
		})
	}
	return rst
}

func flattenRepository(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curJson := utils.PathSearch("repo", resp, make(map[string]interface{}))
	repoElem := map[string]interface{}{
		"https_url":  utils.PathSearch("https_url", curJson, nil),
		"web_url":    utils.PathSearch("web_url", curJson, nil),
		"status":     utils.PathSearch("repo_status", curJson, nil),
		"project_id": utils.PathSearch("project_id", curJson, nil),
	}
	return []interface{}{repoElem}
}

func resourceApplicationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/fgs/applications/{id}"
	)

	client, err := cfg.NewServiceClient("fgs", region)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{id}", d.Id())
	// Due to API restrictions, the request body must pass in an empty JSON.
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error deleting FunctionGraph application: %s", err)
	}

	err = waitForApplicationDeleted(ctx, client, d)
	if err != nil {
		diag.Errorf("error waiting for the application (%s) status to become deleted: %s", d.Id(), err)
	}
	return nil
}

func waitForApplicationDeleted(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      applicationStatusRefreshFunc(client, d, nil),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}
