package apig

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var strSliceParamKeysForPluginBatchApisAssociate = []string{"api_ids"}

// ResourcePluginBatchApisAssociate defines the provider resource of the APIG plugin binding.
// @API APIG PUT /v2/{project_id}/apigw/instances/{instance_id}/plugins/{plugin_id}/detach
// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/plugins/{plugin_id}/attached-apis
// @API APIG POST /v2/{project_id}/apigw/instances/{instance_id}/plugins/{plugin_id}/attach
func ResourcePluginBatchApisAssociate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePluginBatchApisAssociateCreate,
		ReadContext:   resourcePluginBatchApisAssociateRead,
		UpdateContext: resourcePluginBatchApisAssociateUpdate,
		DeleteContext: resourcePluginBatchApisAssociateDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourcePluginBatchApisAssociateImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the plugin is located.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the dedicated instance to which the plugin belongs.",
			},
			"plugin_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The plugin ID.",
			},
			"env_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The environment ID where the API was published.",
			},
			"api_ids": {
				Type:             schema.TypeSet,
				Required:         true,
				Description:      `The list of API IDs to be bound by the plugin.`,
				Elem:             &schema.Schema{Type: schema.TypeString},
				DiffSuppressFunc: utils.SuppressStrSliceDiffs(),
			},
			"api_ids_origin": {
				Type:             schema.TypeList,
				Optional:         true,
				Computed:         true,
				Elem:             &schema.Schema{Type: schema.TypeString},
				DiffSuppressFunc: utils.SuppressDiffAll,
				Description: utils.SchemaDesc(
					`The script configuration value of this change is also the original value used for comparison with
the new value next time the change is made. The corresponding parameter name is 'api_ids'.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
		},
	}
}

func bindPluginToApis(client *golangsdk.ServiceClient, instanceId, pluginId, envId string, apiIds []interface{}) error {
	httpUrl := "v2/{project_id}/apigw/instances/{instance_id}/plugins/{plugin_id}/attach"
	bindPath := client.Endpoint + httpUrl
	bindPath = strings.ReplaceAll(bindPath, "{project_id}", client.ProjectID)
	bindPath = strings.ReplaceAll(bindPath, "{instance_id}", instanceId)
	bindPath = strings.ReplaceAll(bindPath, "{plugin_id}", pluginId)

	bindOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: map[string]interface{}{
			"env_id":  envId,
			"api_ids": apiIds,
		},
	}

	requestBody, err := client.Request("POST", bindPath, &bindOpt)
	if err != nil {
		return fmt.Errorf("error binding APIs to plugin (%s) under dedicated instance (%s): %s", pluginId, instanceId, err)
	}
	respBody, err := utils.FlattenResponse(requestBody)
	if err != nil {
		return err
	}

	// Check for failed bindings
	failedApiRecords := utils.PathSearch("bindings[?binding_result.status=='FAILED'].api_id", respBody, make([]interface{}, 0)).([]interface{})
	if len(failedApiRecords) > 0 {
		return fmt.Errorf("error binding APIs to plugin (%s) under dedicated instance (%s), the failed API IDs are: %s",
			pluginId, instanceId, failedApiRecords)
	}
	return nil
}

func resourcePluginBatchApisAssociateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.ApigV2Client(region)
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	var (
		instanceId = d.Get("instance_id").(string)
		pluginId   = d.Get("plugin_id").(string)
		envId      = d.Get("env_id").(string)
		apiIds     = d.Get("api_ids").(*schema.Set).List()
		resourceId = fmt.Sprintf("%s/%s/%s", instanceId, pluginId, envId)
	)
	// Lock the resource to prevent concurrent updates (error APIG.3500 will be returned if the etcd data synchronize
	// failed)
	config.MutexKV.Lock(resourceId)
	defer config.MutexKV.Unlock(resourceId)

	err = bindPluginToApis(client, instanceId, pluginId, envId, apiIds)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(resourceId)

	// If the request is successful, obtain the values of all slice parameters first and save them to the corresponding
	// '_origin' attributes for subsequent determination and construction of the request body during next updates.
	// And whether corresponding parameters are changed, the origin values must be refreshed.
	err = utils.RefreshSliceParamOriginValues(d, strSliceParamKeysForPluginBatchApisAssociate)
	if err != nil {
		// Don't fail the creation if origin refresh fails
		log.Printf("[WARN] Unable to refresh the origin values: %s", err)
	}

	return resourcePluginBatchApisAssociateRead(ctx, d, meta)
}

func listPluginBoundApisUnderEnv(client *golangsdk.ServiceClient, instanceId, pluginId, envId string) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/apigw/instances/{instance_id}/plugins/{plugin_id}/attached-apis?env_id={env_id}&limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", instanceId)
	listPath = strings.ReplaceAll(listPath, "{plugin_id}", pluginId)
	listPath = strings.ReplaceAll(listPath, "{env_id}", envId)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestBody, err := client.Request("GET", listPathWithOffset, &getOpt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestBody)
		if err != nil {
			return nil, err
		}

		boundApis := utils.PathSearch("apis", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, boundApis...)
		if len(boundApis) < limit {
			break
		}
		offset += limit
	}
	if len(result) < 1 {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v2/{project_id}/apigw/instances/{instance_id}/plugins/{plugin_id}/attached-apis",
				RequestId: "NONE",
				Body:      []byte("all APIs are not bound to the plugin"),
			},
		}
	}
	return result, nil
}

func GetLocalBoundApiIdsForPlugin(client *golangsdk.ServiceClient, instanceId, pluginId, envId string,
	originApiIds []interface{}) ([]interface{}, error) {
	boundApis, err := listPluginBoundApisUnderEnv(client, instanceId, pluginId, envId)
	if err != nil {
		return nil, err
	}

	// Extract API IDs from the bound APIs
	boundApiIds := utils.PathSearch("[*].api_id", boundApis, make([]interface{}, 0)).([]interface{})
	if len(originApiIds) > 0 && len(utils.FildSliceIntersection(boundApiIds, originApiIds)) < 1 {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v2/{project_id}/apigw/instances/{instance_id}/plugins/{plugin_id}/attached-apis",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("all locally managed APIs have been unbound: %v", originApiIds)),
			},
		}
	}
	return boundApiIds, nil
}

func resourcePluginBatchApisAssociateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.ApigV2Client(region)
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	var (
		instanceId   = d.Get("instance_id").(string)
		pluginId     = d.Get("plugin_id").(string)
		envId        = d.Get("env_id").(string)
		originApiIds = d.Get("api_ids_origin").([]interface{})
	)

	boundApiIds, err := GetLocalBoundApiIdsForPlugin(client, instanceId, pluginId, envId, originApiIds)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error querying bound APIs from plugin (%s) under specified environment (%s)",
			pluginId, envId))
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("api_ids", boundApiIds),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving binding fields for specified plugin (%s): %s", pluginId, err)
	}

	return nil
}

func unbindPluginFromApi(client *golangsdk.ServiceClient, instanceId, pluginId, envId string, apiIds []string) error {
	httpUrl := "v2/{project_id}/apigw/instances/{instance_id}/plugins/{plugin_id}/detach"
	unbindPath := client.Endpoint + httpUrl
	unbindPath = strings.ReplaceAll(unbindPath, "{project_id}", client.ProjectID)
	unbindPath = strings.ReplaceAll(unbindPath, "{instance_id}", instanceId)
	unbindPath = strings.ReplaceAll(unbindPath, "{plugin_id}", pluginId)

	unbindOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: map[string]interface{}{
			"env_id":  envId,
			"api_ids": apiIds,
		},
		OkCodes: []int{204},
	}

	_, err := client.Request("PUT", unbindPath, &unbindOpt)
	if err != nil {
		return fmt.Errorf("error unbinding APIs (%v) from plugin (%s) under specified environment (%s): %s", apiIds, pluginId, envId, err)
	}
	return nil
}

func unbindPluginFromApis(client *golangsdk.ServiceClient, instanceId, pluginId, envId string, apiIds []interface{}) error {
	notFoundErr := fmt.Sprintf("[DEBUG] All APIs have been unbound from plugin (%s) under dedicated instance (%s)", pluginId, instanceId)

	boundApis, err := listPluginBoundApisUnderEnv(client, instanceId, pluginId, envId)
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			log.Println(notFoundErr)
			return nil
		}
		return fmt.Errorf("error querying bound APIs for plugin (%s) under dedicated instance (%s): %s", pluginId, instanceId, err)
	}

	var (
		mErr                  *multierror.Error
		prepareToUnbindApiIds = make([]string, 0)
	)
	for _, apiId := range apiIds {
		apiId := utils.PathSearch(fmt.Sprintf("[?api_id=='%s'].api_id|[0]", apiId), boundApis, "").(string)
		if apiId == "" {
			log.Printf("[DEBUG] Unable to find the bound API ID (%s), so skip this unbinding", apiId)
			continue
		}
		prepareToUnbindApiIds = append(prepareToUnbindApiIds, apiId)
	}

	if len(prepareToUnbindApiIds) < 1 {
		log.Println(notFoundErr)
		return nil
	}

	log.Printf("[DEBUG] Prepare to unbind API IDs (%v) from plugin (%s)", prepareToUnbindApiIds, pluginId)
	err = unbindPluginFromApi(client, instanceId, pluginId, envId, prepareToUnbindApiIds)
	if err != nil {
		mErr = multierror.Append(mErr, err)
	}

	return mErr.ErrorOrNil()
}

func resourcePluginBatchApisAssociateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	var (
		resourceId = d.Id()
		instanceId = d.Get("instance_id").(string)
		pluginId   = d.Get("plugin_id").(string)
		envId      = d.Get("env_id").(string)

		consoleApiIds, scriptApiIds = d.GetChange("api_ids")

		consoleApiIdsList = consoleApiIds.(*schema.Set).List()
		scriptApiIdsList  = scriptApiIds.(*schema.Set).List()
		originApiIdsList  = d.Get("api_ids_origin").([]interface{})
	)

	// Lock the resource to prevent concurrent updates (error APIG.3500 will be returned if the etcd data synchronize
	// failed)
	config.MutexKV.Lock(resourceId)
	defer config.MutexKV.Unlock(resourceId)

	newApiIds := utils.FindSliceElementsNotInAnother(scriptApiIdsList, consoleApiIdsList)
	rmApiIds := utils.FindSliceElementsNotInAnother(originApiIdsList, scriptApiIdsList)

	if len(rmApiIds) > 0 {
		log.Printf("[DEBUG] Prepare to unbind the specified API IDs: %v", rmApiIds)
		err := unbindPluginFromApis(client, instanceId, pluginId, envId, rmApiIds)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if len(newApiIds) > 0 {
		log.Printf("[DEBUG] Prepare to bind the specified API IDs: %v", newApiIds)
		err = bindPluginToApis(client, instanceId, pluginId, envId, newApiIds)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// If the request is successful, obtain the values of all slice parameters first and save them to the corresponding
	// '_origin' attributes for subsequent determination and construction of the request body during next updates.
	// And whether corresponding parameters are changed, the origin values must be refreshed.
	err = utils.RefreshSliceParamOriginValues(d, strSliceParamKeysForPluginBatchApisAssociate)
	if err != nil {
		// Don't fail the update if origin refresh fails
		log.Printf("[WARN] Unable to refresh the origin values: %s", err)
	}

	return resourcePluginBatchApisAssociateRead(ctx, d, meta)
}

func resourcePluginBatchApisAssociateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	var (
		resourceId = d.Id()
		instanceId = d.Get("instance_id").(string)
		pluginId   = d.Get("plugin_id").(string)
		envId      = d.Get("env_id").(string)

		rmApiIds = getConfiguredApiIdsForPlugin(d)
	)

	// Lock the resource to prevent concurrent updates (error APIG.3500 will be returned if the etcd data synchronize
	// failed)
	config.MutexKV.Lock(resourceId)
	defer config.MutexKV.Unlock(resourceId)

	if err := unbindPluginFromApis(client, instanceId, pluginId, envId, rmApiIds); err != nil {
		return diag.Errorf("error deleting configured API bindings: %s", err)
	}

	return nil
}

func resourcePluginBatchApisAssociateImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<instance_id>/<plugin_id>/<env_id>', "+
			"but got '%s'", importedId)
	}

	mErr := multierror.Append(nil,
		d.Set("instance_id", parts[0]),
		d.Set("plugin_id", parts[1]),
		d.Set("env_id", parts[2]),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("error saving associate resource field: %s", err)
	}
	return []*schema.ResourceData{d}, nil
}

// getConfiguredApiIdsForPlugin retrieves API IDs from configuration or origin
func getConfiguredApiIdsForPlugin(d *schema.ResourceData) []interface{} {
	// Fallback to origin (last known configuration)
	if origin, ok := d.Get("api_ids_origin").([]interface{}); ok && len(origin) > 0 {
		log.Printf("[DEBUG] Found %d API ID(s) from the origin attribute: %v", len(origin), origin)
		return origin
	}

	log.Printf("[DEBUG] Unable to find the API IDs from the origin attribute, so try to get from current state")
	// After resource imported, the origin attribute is not set, so try to get from current state
	current := d.Get("api_ids").(*schema.Set).List()
	log.Printf("[DEBUG] Found %d API ID(s) from the current state: %v", len(current), current)

	return current
}
