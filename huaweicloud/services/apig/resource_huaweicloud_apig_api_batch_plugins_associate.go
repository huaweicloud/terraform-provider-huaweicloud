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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var (
	strSliceParamKeysForApiBatchPluginsAssociate = []string{
		"plugin_ids",
	}
	apiBatchPluginsAssociateNonUpdatableParams = []string{
		"instance_id",
		"api_id",
		"env_id",
	}
)

// @API APIG POST /v2/{project_id}/apigw/instances/{instance_id}/apis/{api_id}/plugins/attach
// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/apis/{api_id}/attached-plugins
// @API APIG PUT /v2/{project_id}/apigw/instances/{instance_id}/apis/{api_id}/plugins/detach
func ResourceApiBatchPluginsAssociate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApiBatchPluginsAssociateCreate,
		ReadContext:   resourceApiBatchPluginsAssociateRead,
		UpdateContext: resourceApiBatchPluginsAssociateUpdate,
		DeleteContext: resourceApiBatchPluginsAssociateDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceApiBatchPluginsAssociateImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(apiBatchPluginsAssociateNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the API and plugins are located.",
			},

			// Required parameters.
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the dedicated instance to which the API and plugins belong.",
			},
			"api_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the API to be bound with plugins.",
			},
			"env_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the environment where the API is published.",
			},
			"plugin_ids": {
				Type:             schema.TypeSet,
				Required:         true,
				Description:      `The list of plugin IDs to be bound to the API.`,
				Elem:             &schema.Schema{Type: schema.TypeString},
				DiffSuppressFunc: utils.SuppressStrSliceDiffs(),
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},

			// Internal attributes.
			"plugin_ids_origin": {
				Type:             schema.TypeList,
				Optional:         true,
				Computed:         true,
				Elem:             &schema.Schema{Type: schema.TypeString},
				DiffSuppressFunc: utils.SuppressDiffAll,
				Description: utils.SchemaDesc(
					`The script configuration value of this change is also the original value used for comparison with
the new value next time the change is made. The corresponding parameter name is 'plugin_ids'.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
		},
	}
}

func bindPluginsToApi(client *golangsdk.ServiceClient, instanceId, apiId, envId string, pluginIds []interface{}) error {
	httpUrl := "v2/{project_id}/apigw/instances/{instance_id}/apis/{api_id}/plugins/attach"
	bindPath := client.Endpoint + httpUrl
	bindPath = strings.ReplaceAll(bindPath, "{project_id}", client.ProjectID)
	bindPath = strings.ReplaceAll(bindPath, "{instance_id}", instanceId)
	bindPath = strings.ReplaceAll(bindPath, "{api_id}", apiId)

	bindOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: map[string]interface{}{
			"env_id":     envId,
			"plugin_ids": pluginIds,
		},
	}

	requestBody, err := client.Request("POST", bindPath, &bindOpt)
	if err != nil {
		return fmt.Errorf("error binding plugins to API (%s) under dedicated instance (%s): %s", apiId, instanceId, err)
	}
	respBody, err := utils.FlattenResponse(requestBody)
	if err != nil {
		return err
	}

	// Check for failed bindings
	failedPluginRecords := utils.PathSearch("attached_plugins[?attached_plugins.status=='FAILED'].plugin_id", respBody,
		make([]interface{}, 0)).([]interface{})
	if len(failedPluginRecords) > 0 {
		return fmt.Errorf("error binding plugins to API (%s) under dedicated instance (%s), the failed plugin IDs are: %s",
			apiId, instanceId, failedPluginRecords)
	}
	return nil
}

func resourceApiBatchPluginsAssociateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient(region, "apig")
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	var (
		instanceId = d.Get("instance_id").(string)
		apiId      = d.Get("api_id").(string)
		envId      = d.Get("env_id").(string)
		pluginIds  = d.Get("plugin_ids").(*schema.Set).List()
		resourceId = fmt.Sprintf("%s/%s/%s", instanceId, apiId, envId)
	)
	// Lock the resource to prevent concurrent updates (error APIG.3500 will be returned if the etcd data synchronize
	// failed)
	config.MutexKV.Lock(resourceId)
	defer config.MutexKV.Unlock(resourceId)

	err = bindPluginsToApi(client, instanceId, apiId, envId, pluginIds)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(resourceId)

	// If the request is successful, obtain the values of all slice parameters first and save them to the corresponding
	// '_origin' attributes for subsequent determination and construction of the request body during next updates.
	// And whether corresponding parameters are changed, the origin values must be refreshed.
	err = utils.RefreshSliceParamOriginValues(d, strSliceParamKeysForApiBatchPluginsAssociate)
	if err != nil {
		// Don't fail the creation if origin refresh fails
		log.Printf("[WARN] Unable to refresh the origin values: %s", err)
	}

	return resourceApiBatchPluginsAssociateRead(ctx, d, meta)
}

func listApiBoundPluginsUnderEnv(client *golangsdk.ServiceClient, instanceId, apiId, envId string) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/apigw/instances/{instance_id}/apis/{api_id}/attached-plugins?env_id={env_id}&limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", instanceId)
	listPath = strings.ReplaceAll(listPath, "{api_id}", apiId)
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

		boundPlugins := utils.PathSearch("plugins", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, boundPlugins...)
		if len(boundPlugins) < limit {
			break
		}
		offset += limit
	}
	if len(result) < 1 {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v2/{project_id}/apigw/instances/{instance_id}/apis/{api_id}/attached-plugins",
				RequestId: "NONE",
				Body:      []byte("all plugins are not bound to the API"),
			},
		}
	}
	return result, nil
}

func GetLocalBoundPluginIdsForApi(client *golangsdk.ServiceClient, instanceId, apiId, envId string,
	originPluginIds []interface{}) ([]interface{}, error) {
	boundPlugins, err := listApiBoundPluginsUnderEnv(client, instanceId, apiId, envId)
	if err != nil {
		return nil, err
	}

	// Extract plugin IDs from the bound plugins
	boundPluginIds := utils.PathSearch("[*].plugin_id", boundPlugins, make([]interface{}, 0)).([]interface{})
	if len(originPluginIds) > 0 && len(utils.FildSliceIntersection(boundPluginIds, originPluginIds)) < 1 {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v2/{project_id}/apigw/instances/{instance_id}/apis/{api_id}/attached-plugins",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("all locally managed plugins have been unbound: %v", originPluginIds)),
			},
		}
	}
	return boundPluginIds, nil
}

func resourceApiBatchPluginsAssociateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient(region, "apig")
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	var (
		instanceId      = d.Get("instance_id").(string)
		apiId           = d.Get("api_id").(string)
		envId           = d.Get("env_id").(string)
		originPluginIds = d.Get("plugin_ids_origin").([]interface{})
	)

	boundPluginIds, err := GetLocalBoundPluginIdsForApi(client, instanceId, apiId, envId, originPluginIds)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error querying bound plugins from API (%s) under specified environment (%s)",
			apiId, envId))
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("plugin_ids", boundPluginIds),
	)

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving binding fields for specified API (%s): %s", apiId, err)
	}

	return nil
}

func unbindPluginsFromApi(client *golangsdk.ServiceClient, instanceId, apiId, envId string, pluginIds []string) error {
	httpUrl := "v2/{project_id}/apigw/instances/{instance_id}/apis/{api_id}/plugins/detach"
	unbindPath := client.Endpoint + httpUrl
	unbindPath = strings.ReplaceAll(unbindPath, "{project_id}", client.ProjectID)
	unbindPath = strings.ReplaceAll(unbindPath, "{instance_id}", instanceId)
	unbindPath = strings.ReplaceAll(unbindPath, "{api_id}", apiId)

	unbindOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: map[string]interface{}{
			"env_id":     envId,
			"plugin_ids": pluginIds,
		},
		OkCodes: []int{204},
	}

	_, err := client.Request("PUT", unbindPath, &unbindOpt)
	if err != nil {
		return fmt.Errorf("error unbinding plugins (%v) from API (%s) under specified environment (%s): %s", pluginIds, apiId, envId, err)
	}
	return nil
}

func unbindPluginsFromApis(client *golangsdk.ServiceClient, instanceId, apiId, envId string, pluginIds []interface{}) error {
	notFoundErr := fmt.Sprintf("[DEBUG] All plugins have been unbound from API (%s) under dedicated instance (%s)", apiId, instanceId)

	boundPlugins, err := listApiBoundPluginsUnderEnv(client, instanceId, apiId, envId)
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			log.Println(notFoundErr)
			return nil
		}
		return fmt.Errorf("error querying bound plugins for API (%s) under dedicated instance (%s): %s", apiId, instanceId, err)
	}

	var (
		mErr                     *multierror.Error
		prepareToUnbindPluginIds = make([]string, 0)
	)
	for _, pluginId := range pluginIds {
		boundPluginId := utils.PathSearch(fmt.Sprintf("[?plugin_id=='%s'].plugin_id|[0]", pluginId), boundPlugins, "").(string)
		if boundPluginId == "" {
			log.Printf("[DEBUG] Unable to find the bound plugin ID (%s), so skip this unbinding", pluginId)
			continue
		}
		prepareToUnbindPluginIds = append(prepareToUnbindPluginIds, boundPluginId)
	}

	if len(prepareToUnbindPluginIds) < 1 {
		log.Println(notFoundErr)
		return nil
	}

	log.Printf("[DEBUG] Prepare to unbind plugin IDs (%v) from API (%s)", prepareToUnbindPluginIds, apiId)
	err = unbindPluginsFromApi(client, instanceId, apiId, envId, prepareToUnbindPluginIds)
	if err != nil {
		mErr = multierror.Append(mErr, err)
	}

	return mErr.ErrorOrNil()
}

func resourceApiBatchPluginsAssociateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient(region, "apig")
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	var (
		resourceId = d.Id()
		instanceId = d.Get("instance_id").(string)
		apiId      = d.Get("api_id").(string)
		envId      = d.Get("env_id").(string)

		consolePluginIds, scriptPluginIds = d.GetChange("plugin_ids")

		consolePluginIdsList = consolePluginIds.(*schema.Set).List()
		scriptPluginIdsList  = scriptPluginIds.(*schema.Set).List()
		originPluginIdsList  = d.Get("plugin_ids_origin").([]interface{})
	)

	// Lock the resource to prevent concurrent updates (error APIG.3500 will be returned if the etcd data synchronize
	// failed)
	config.MutexKV.Lock(resourceId)
	defer config.MutexKV.Unlock(resourceId)

	newPluginIds := utils.FindSliceElementsNotInAnother(scriptPluginIdsList, consolePluginIdsList)
	rmPluginIds := utils.FindSliceElementsNotInAnother(originPluginIdsList, scriptPluginIdsList)

	if len(rmPluginIds) > 0 {
		log.Printf("[DEBUG] Prepare to unbind the specified plugin IDs: %v", rmPluginIds)
		err := unbindPluginsFromApis(client, instanceId, apiId, envId, rmPluginIds)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if len(newPluginIds) > 0 {
		log.Printf("[DEBUG] Prepare to bind the specified plugin IDs: %v", newPluginIds)
		err = bindPluginsToApi(client, instanceId, apiId, envId, newPluginIds)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// If the request is successful, obtain the values of all slice parameters first and save them to the corresponding
	// '_origin' attributes for subsequent determination and construction of the request body during next updates.
	// And whether corresponding parameters are changed, the origin values must be refreshed.
	err = utils.RefreshSliceParamOriginValues(d, strSliceParamKeysForApiBatchPluginsAssociate)
	if err != nil {
		// Don't fail the update if origin refresh fails
		log.Printf("[WARN] Unable to refresh the origin values: %s", err)
	}

	return resourceApiBatchPluginsAssociateRead(ctx, d, meta)
}

// getConfiguredPluginIdsForApi retrieves plugin IDs from configuration or origin
func getConfiguredPluginIdsForApi(d *schema.ResourceData) []interface{} {
	// Fallback to origin (last known configuration)
	if origin, ok := d.Get("plugin_ids_origin").([]interface{}); ok && len(origin) > 0 {
		log.Printf("[DEBUG] Found %d plugin ID(s) from the origin attribute: %v", len(origin), origin)
		return origin
	}

	log.Printf("[DEBUG] Unable to find the plugin IDs from the origin attribute, so try to get from current state")
	// After resource imported, the origin attribute is not set, so try to get from current state
	current := d.Get("plugin_ids").(*schema.Set).List()
	log.Printf("[DEBUG] Found %d plugin ID(s) from the current state: %v", len(current), current)

	return current
}

func resourceApiBatchPluginsAssociateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient(region, "apig")
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	var (
		resourceId = d.Id()
		instanceId = d.Get("instance_id").(string)
		apiId      = d.Get("api_id").(string)
		envId      = d.Get("env_id").(string)

		rmPluginIds = getConfiguredPluginIdsForApi(d)
	)

	// Lock the resource to prevent concurrent updates (error APIG.3500 will be returned if the etcd data synchronize
	// failed)
	config.MutexKV.Lock(resourceId)
	defer config.MutexKV.Unlock(resourceId)

	if err := unbindPluginsFromApis(client, instanceId, apiId, envId, rmPluginIds); err != nil {
		return diag.Errorf("error deleting configured plugin bindings: %s", err)
	}

	return nil
}

func resourceApiBatchPluginsAssociateImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<instance_id>/<api_id>/<env_id>', "+
			"but got '%s'", importedId)
	}

	mErr := multierror.Append(nil,
		d.Set("instance_id", parts[0]),
		d.Set("api_id", parts[1]),
		d.Set("env_id", parts[2]),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("error saving associate resource field: %s", err)
	}
	return []*schema.ResourceData{d}, nil
}
