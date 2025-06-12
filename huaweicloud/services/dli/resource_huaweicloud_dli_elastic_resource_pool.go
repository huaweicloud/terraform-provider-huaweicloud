package dli

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

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DLI POST /v3/{project_id}/elastic-resource-pools
// @API DLI GET /v3/{project_id}/elastic-resource-pools
// @API DLI PUT /v3/{project_id}/elastic-resource-pools/{elastic_resource_pool_name}
// @API DLI DELETE /v3/{project_id}/elastic-resource-pools/{elastic_resource_pool_name}
func ResourceElasticResourcePool() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceElasticResourcePoolCreate,
		ReadContext:   resourceElasticResourcePoolRead,
		UpdateContext: resourceElasticResourcePoolUpdate,
		DeleteContext: resourceElasticResourcePoolDelete,

		Timeouts: &schema.ResourceTimeout{
			// Create and Update method both using custom default timeout as general timeout of the corresponding
			// refresh function.
			Default: schema.DefaultTimeout(30 * time.Minute),
			Delete:  schema.DefaultTimeout(10 * time.Minute),
		},

		Importer: &schema.ResourceImporter{
			StateContext: resourceAssociationImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the elastic resource pool is located.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The name of the elastic resource pool.`,
				// The server will automatically convert uppercase letters to lowercase letters.
				DiffSuppressFunc: utils.SuppressCaseDiffs(),
			},
			"max_cu": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The maximum number of CUs for elastic resource pool scaling.`,
			},
			"min_cu": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The minimum number of CUs for elastic resource pool scaling.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the elastic resource pool.`,
			},
			"cidr": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The CIDR block of network to associate with the elastic resource pool.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the enterprise project to which the elastic resource pool belongs.`,
			},
			"tags": common.TagsForceNewSchema(),
			"label": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The attribute fields of the elastic resource pool.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The current status of the elastic resource pool.`,
			},
			"current_cu": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The current CU number of the elastic resource pool.`,
			},
			"actual_cu": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The actual CU number of the elastic resource pool.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the elastic resource pool.`,
			},
		},
	}
}

func resourceElasticResourcePoolCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/elastic-resource-pools"
	)

	client, err := cfg.NewServiceClient("dli", region)
	if err != nil {
		return diag.Errorf("error creating DLI client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateElasticResourcePoolBodyParams(cfg, d)),
	}

	requestResp, err := client.Request("POST", createPath, &createOpts)
	if err != nil {
		return diag.Errorf("error creating DLI elastic resource pool: %s", err)
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}
	if !utils.PathSearch("is_success", respBody, true).(bool) {
		return diag.Errorf("unable to create the elastic resource pool: %s",
			utils.PathSearch("message", respBody, "Message Not Found"))
	}

	resourceName := utils.PathSearch("elastic_resource_pool_name", respBody, "").(string)
	respBody, err = GetElasticResourcePoolByName(client, resourceName)
	if err != nil {
		return diag.FromErr(err)
	}
	resourceId := utils.PathSearch("resource_id", respBody, "").(string)
	d.SetId(resourceId)

	err = waitForElasticResourcePoolStatusCompleted(ctx, client, d)
	if err != nil {
		diag.Errorf("error waiting for the elastic resource pool (%s) status to become success: %s", resourceId, err)
	}

	return resourceElasticResourcePoolRead(ctx, d, meta)
}

func buildCreateElasticResourcePoolBodyParams(cfg *config.Config, d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"elastic_resource_pool_name": d.Get("name"),
		"description":                d.Get("description"),
		"max_cu":                     d.Get("max_cu"),
		"min_cu":                     d.Get("min_cu"),
		"cidr_in_vpc":                utils.StringIgnoreEmpty(d.Get("cidr").(string)),
		"enterprise_project_id":      cfg.GetEnterpriseProjectID(d),
		"tags":                       utils.ExpandResourceTagsMap(d.Get("tags").(map[string]interface{})),
		"label":                      utils.ValueIgnoreEmpty(d.Get("label")),
	}
}

func waitForElasticResourcePoolStatusCompleted(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      elasticResourcePoolStatusRefreshFunc(client, d, []string{"AVAILABLE"}),
		Timeout:      d.Timeout(schema.TimeoutDefault),
		Delay:        10 * time.Second,
		PollInterval: 20 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func elasticResourcePoolStatusRefreshFunc(client *golangsdk.ServiceClient, d *schema.ResourceData, targets []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resourceName := d.Get("name").(string)
		respBody, err := GetElasticResourcePoolByName(client, resourceName)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok && len(targets) < 1 {
				log.Printf("[DEBUG] The DLI elastic resource pool (%s) has been deleted", resourceName)
				// When the error code is 404, the value of respBody is nil, and a non-null value is returned to avoid continuing the loop check.
				return "Resource Not Found", "COMPLETED", nil
			}
			return respBody, "ERROR", err
		}

		status := utils.PathSearch("status", respBody, "").(string)
		unexpectedStatuses := []string{
			"FAILED",
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

func resourceElasticResourcePoolRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("dli", region)
	if err != nil {
		return diag.Errorf("error creating DLI client: %s", err)
	}

	respBody, err := GetElasticResourcePoolByName(client, d.Get("name").(string))
	if err != nil {
		return common.CheckDeletedDiag(d, err, "DLI elastic resource pool")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("elastic_resource_pool_name", respBody, nil)),
		d.Set("max_cu", int(utils.PathSearch("max_cu", respBody, float64(0)).(float64))),
		d.Set("min_cu", int(utils.PathSearch("min_cu", respBody, float64(0)).(float64))),
		d.Set("description", utils.PathSearch("description", respBody, nil)),
		d.Set("cidr", utils.PathSearch("cidr_in_vpc", respBody, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", respBody, nil)),
		d.Set("status", utils.PathSearch("status", respBody, nil)),
		d.Set("current_cu", utils.PathSearch("current_cu", respBody, nil)),
		d.Set("actual_cu", utils.PathSearch("actual_cu", respBody, nil)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(
			int64(utils.PathSearch("create_time", respBody, float64(0)).(float64))/1000, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

// GetElasticResourcePools is a method used to query all elastic resource pools in a specified region.
func GetElasticResourcePools(client *golangsdk.ServiceClient) ([]interface{}, error) {
	var (
		httpUrl  = "v3/{project_id}/elastic-resource-pools"
		listOpts = golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}
		offset = 0
		limit  = 100
		result = make([]interface{}, 0)
	)

	basePath := client.Endpoint + httpUrl
	basePath = strings.ReplaceAll(basePath, "{project_id}", client.ProjectID)

	for {
		listPath := basePath + fmt.Sprintf("?offset=%d&limit=%d", offset, limit)
		requestResp, err := client.Request("GET", listPath, &listOpts)
		if err != nil {
			return nil, fmt.Errorf("error query DLI elastic resource pool list: %s", err)
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		if !utils.PathSearch("is_success", respBody, true).(bool) {
			return nil, fmt.Errorf("unable to query the elastic resource pools: %s",
				utils.PathSearch("message", respBody, "Message Not Found"))
		}
		pools := utils.PathSearch("elastic_resource_pools", respBody, make([]interface{}, 0)).([]interface{})
		if len(pools) < 1 {
			break
		}
		result = append(result, pools...)
		offset += len(pools)
	}

	return result, nil
}

// GetElasticResourcePoolByName is the method used to query the elastic resource pool matching the name.
func GetElasticResourcePoolByName(client *golangsdk.ServiceClient, resourceName string) (interface{}, error) {
	pools, err := GetElasticResourcePools(client)
	if err != nil {
		return nil, err
	}

	// The elastic resource pool name in the API response does not contain uppercase letters (will be automatically converted).
	lowercaseName := strings.ToLower(resourceName)
	if pool := utils.PathSearch(fmt.Sprintf("[*]|[?elastic_resource_pool_name=='%s']|[0]", lowercaseName), pools, nil); pool != nil {
		return pool, nil
	}
	return nil, golangsdk.ErrDefault404{
		ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
			Body: []byte(fmt.Sprintf("unable to find the elastic resource pool using its name (%s)", lowercaseName)),
		},
	}
}

func updateElasticResourcePoolMetadata(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	resourceName := d.Get("name").(string)
	httpUrl := "v3/{project_id}/elastic-resource-pools/{elastic_resource_pool_name}"

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{elastic_resource_pool_name}", resourceName)
	updateOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildUpdateElasticResourcePoolBodyParams(d),
	}

	requestResp, err := client.Request("PUT", updatePath, &updateOpts)
	if err != nil {
		return fmt.Errorf("error updating DLI elastic resource pool: %s", err)
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return err
	}
	if !utils.PathSearch("is_success", respBody, true).(bool) {
		return fmt.Errorf("unable to update the elastic resource pool: %s",
			utils.PathSearch("message", respBody, "Message Not Found"))
	}

	err = waitForElasticResourcePoolStatusCompleted(ctx, client, d)
	if err != nil {
		return fmt.Errorf("error waiting for the elastic resource pool (%s) status to become success: %s", resourceName, err)
	}
	return nil
}

func resourceElasticResourcePoolUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		resourceId = d.Id()
	)

	client, err := cfg.NewServiceClient("dli", region)
	if err != nil {
		return diag.Errorf("error creating DLI client: %s", err)
	}

	if d.HasChanges("description", "min_cu", "max_cu") {
		if err := updateElasticResourcePoolMetadata(ctx, client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("enterprise_project_id") {
		migrateOpts := config.MigrateResourceOpts{
			ResourceId:   resourceId,
			ResourceType: "dli-elastic-resource-pool",
			RegionId:     region,
			ProjectId:    client.ProjectID,
		}
		if err := cfg.MigrateEnterpriseProject(ctx, d, migrateOpts); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceElasticResourcePoolRead(ctx, d, meta)
}

func buildUpdateElasticResourcePoolBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"description": d.Get("description"),
		"max_cu":      utils.IntIgnoreEmpty(d.Get("max_cu").(int)),
		"min_cu":      utils.IntIgnoreEmpty(d.Get("min_cu").(int)),
	}
}

func resourceElasticResourcePoolDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		httpUrl      = "v3/{project_id}/elastic-resource-pools/{elastic_resource_pool_name}"
		resourceName = d.Get("name").(string)
	)

	client, err := cfg.NewServiceClient("dli", region)
	if err != nil {
		return diag.Errorf("error creating DLI client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{elastic_resource_pool_name}", resourceName)
	// Due to API restrictions, the request body must pass in an empty JSON.
	deleteOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	requestResp, err := client.Request("DELETE", deletePath, &deleteOpts)
	if err != nil {
		return diag.Errorf("error deleting DLI elastic resource pool (%s): %s", resourceName, err)
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}
	if !utils.PathSearch("is_success", respBody, true).(bool) {
		return diag.Errorf("unable to delete the elastic resource pool: %s",
			utils.PathSearch("message", respBody, "Message Not Found"))
	}

	err = waitForElasticResourcePoolDeleted(ctx, client, d)
	if err != nil {
		diag.Errorf("error waiting for the elastic resource pool (%s) status to become deleted: %s", d.Id(), err)
	}
	return nil
}

func waitForElasticResourcePoolDeleted(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      elasticResourcePoolStatusRefreshFunc(client, d, nil),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceAssociationImportState(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData,
	error) {
	var (
		cfg          = meta.(*config.Config)
		resourceName = d.Id()
	)

	client, err := cfg.NewServiceClient("dli", cfg.GetRegion(d))
	if err != nil {
		return nil, fmt.Errorf("error creating DLI client: %s", err)
	}

	respBody, err := GetElasticResourcePoolByName(client, resourceName)
	if err != nil {
		return nil, err
	}

	resourceId := utils.PathSearch("resource_id", respBody, "").(string)
	d.SetId(resourceId)
	return []*schema.ResourceData{d}, d.Set("name", resourceName)
}
