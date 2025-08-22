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

var strSliceParamKeys = []string{"api_ids"}

// @API APIG DELETE /v2/{project_id}/apigw/instances/{instance_id}/app-auths/{app_auth_id}
// @API APIG POST /v2/{project_id}/apigw/instances/{instance_id}/app-auths
// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/app-auths/binded-apis
func ResourceApplicationAuthorization() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApplicationAuthorizationCreate,
		ReadContext:   resourceApplicationAuthorizationRead,
		UpdateContext: resourceApplicationAuthorizationUpdate,
		DeleteContext: resourceApplicationAuthorizationDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceApplicationAuthorizationImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the application and APPCODEs are located.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the dedicated instance to which the application and APIs belong.",
			},
			"application_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the application authorized to access the APIs.",
			},
			"env_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The environment ID where the APIs were published.",
			},
			"api_ids": {
				Type:             schema.TypeSet,
				Required:         true,
				Description:      `The list of API IDs to be authorized for the application.`,
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

func createApplicationAuthorizationForApis(client *golangsdk.ServiceClient, instanceId, envId, appId string, apiIds []interface{}) error {
	httpUrl := "v2/{project_id}/apigw/instances/{instance_id}/app-auths"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: map[string]interface{}{
			"app_ids": []string{appId},
			"env_id":  envId,
			"api_ids": apiIds,
		},
	}

	requestBody, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return fmt.Errorf("error authorizing APIs to application (%s) under dedicated instance (%s): %s", appId, instanceId, err)
	}
	respBody, err := utils.FlattenResponse(requestBody)
	if err != nil {
		return err
	}

	failedApiRecords := utils.PathSearch("auths[?auth_result.status=='FAILED'].api_id", respBody, make([]interface{}, 0)).([]interface{})
	if len(failedApiRecords) > 0 {
		return fmt.Errorf("error authorizing APIs to application (%s) under dedicated instance (%s), the failed API IDs are: %s",
			appId, instanceId, failedApiRecords)
	}
	return nil
}

func resourceApplicationAuthorizationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	var (
		instanceId = d.Get("instance_id").(string)
		appId      = d.Get("application_id").(string)
		envId      = d.Get("env_id").(string)
		apiIds     = d.Get("api_ids").(*schema.Set).List()
		resourceId = fmt.Sprintf("%s/%s/%s", instanceId, envId, appId)
	)
	// Lock the resource to prevent concurrent updates (error APIG.3500 will be returned if the etcd data synchronize
	// failed)
	config.MutexKV.Lock(resourceId)
	defer config.MutexKV.Unlock(resourceId)

	err = createApplicationAuthorizationForApis(client, instanceId, envId, appId, apiIds)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(resourceId)

	// If the request is successful, obtain the values of all slice parameters first and save them to the corresponding
	// '_origin' attributes for subsequent determination and construction of the request body during next updates.
	// And whether corresponding parameters are changed, the origin values must be refreshed.
	err = utils.RefreshSliceParamOriginValues(d, strSliceParamKeys)
	if err != nil {
		// Don't fail the creation if origin refresh fails
		log.Printf("[WARN] Unable to refresh the origin values: %s", err)
	}

	return resourceApplicationAuthorizationRead(ctx, d, meta)
}

func listApplicationAuthorizedApisUnderEnv(client *golangsdk.ServiceClient, instanceId, envId, appId string) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/apigw/instances/{instance_id}/app-auths/binded-apis?app_id={app_id}&env_id={env_id}&limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", instanceId)
	listPath = strings.ReplaceAll(listPath, "{app_id}", appId)
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

		authorizedApis := utils.PathSearch("auths", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, authorizedApis...)
		if len(authorizedApis) < limit {
			break
		}
		offset += limit
	}
	if len(result) < 1 {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v2/{project_id}/apigw/instances/{instance_id}/app-auths/binded-apis",
				RequestId: "NONE",
				Body:      []byte("all APIs are not authorized for the application"),
			},
		}
	}
	return result, nil
}

func GetLocalAuthorizedApiIds(client *golangsdk.ServiceClient, instanceId, envId, appId string, originApiIds []interface{}) ([]interface{}, error) {
	authorizedApis, err := listApplicationAuthorizedApisUnderEnv(client, instanceId, envId, appId)
	if err != nil {
		return nil, err
	}

	authorizedApiIds := utils.PathSearch("[*].api_id", authorizedApis, make([]interface{}, 0)).([]interface{})
	if len(originApiIds) > 0 && len(utils.FildSliceIntersection(authorizedApiIds, originApiIds)) < 1 {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v2/{project_id}/apigw/instances/{instance_id}/app-auths/binded-apis",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("all locally managed APIs have been deauthorized: %v", originApiIds)),
			},
		}
	}
	return authorizedApiIds, nil
}

func resourceApplicationAuthorizationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	var (
		instanceId   = d.Get("instance_id").(string)
		appId        = d.Get("application_id").(string)
		envId        = d.Get("env_id").(string)
		originApiIds = d.Get("api_ids_origin").([]interface{})
	)

	authorizedApiIds, err := GetLocalAuthorizedApiIds(client, instanceId, envId, appId, originApiIds)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error querying authorized APIs from application (%s) under specified environment (%s)",
			appId, envId))
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("api_ids", authorizedApiIds),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving authorization fields for specified application (%s): %s", appId, err)
	}

	return nil
}

func resourceApplicationAuthorizationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	var (
		resourceId = d.Id()
		instanceId = d.Get("instance_id").(string)
		appId      = d.Get("application_id").(string)
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
		log.Printf("[DEBUG] Prepare to delete the authorization for specified API IDs: %v", rmApiIds)
		err := deleteApplicationAuthorizationForApis(client, instanceId, envId, appId, rmApiIds)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if len(newApiIds) > 0 {
		log.Printf("[DEBUG] Prepare to create the authorization for specified API IDs: %v", newApiIds)
		err = createApplicationAuthorizationForApis(client, instanceId, envId, appId, newApiIds)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// If the request is successful, obtain the values of all slice parameters first and save them to the corresponding
	// '_origin' attributes for subsequent determination and construction of the request body during next updates.
	// And whether corresponding parameters are changed, the origin values must be refreshed.
	err = utils.RefreshSliceParamOriginValues(d, strSliceParamKeys)
	if err != nil {
		// Don't fail the update if origin refresh fails
		log.Printf("[WARN] Unable to refresh the origin values: %s", err)
	}

	return resourceApplicationAuthorizationRead(ctx, d, meta)
}

func resourceApplicationAuthorizationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("apig", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	var (
		resourceId = d.Id()
		instanceId = d.Get("instance_id").(string)
		appId      = d.Get("application_id").(string)
		envId      = d.Get("env_id").(string)

		rmApiIds = getConfiguredApiIds(d)
	)

	// Lock the resource to prevent concurrent updates (error APIG.3500 will be returned if the etcd data synchronize
	// failed)
	config.MutexKV.Lock(resourceId)
	defer config.MutexKV.Unlock(resourceId)

	if err := deleteApplicationAuthorizationForApis(client, instanceId, envId, appId, rmApiIds); err != nil {
		return diag.Errorf("error deleting configured API bindings: %s", err)
	}

	return nil
}

// getConfiguredApiIds retrieves API IDs from configuration or origin
func getConfiguredApiIds(d *schema.ResourceData) []interface{} {
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

func deleteApplicationAuthorizationForApi(client *golangsdk.ServiceClient, instanceId, envId, appId, authId string) error {
	httpUrl := "v2/{project_id}/apigw/instances/{instance_id}/app-auths/{app_auth_id}?app_id={app_id}&env_id={env_id}"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", instanceId)
	deletePath = strings.ReplaceAll(deletePath, "{app_auth_id}", authId)
	deletePath = strings.ReplaceAll(deletePath, "{app_id}", appId)
	deletePath = strings.ReplaceAll(deletePath, "{env_id}", envId)

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err := client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return fmt.Errorf("error deleting application authorization (%s) under specified environment (%s): %s", authId, envId, err)
	}
	return nil
}

func deleteApplicationAuthorizationForApis(client *golangsdk.ServiceClient, instanceId, envId, appId string, apiIds []interface{}) error {
	notFoundErr := fmt.Sprintf("[DEBUG] All APIs have been unauthorized form application (%s) under dedicated instance (%s)", appId, instanceId)

	authorizedApis, err := listApplicationAuthorizedApisUnderEnv(client, instanceId, envId, appId)
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			log.Println(notFoundErr)
			return nil
		}
		return fmt.Errorf("error querying authorized APIs for application (%s) under dedicated instance (%s)", appId, instanceId)
	}

	var mErr *multierror.Error
	for _, apiId := range apiIds {
		authId := utils.PathSearch(fmt.Sprintf("[?api_id=='%s'].id|[0]", apiId), authorizedApis, "").(string)
		if authId == "" {
			log.Printf("[DEBUG] Unable to find the authorization ID via API ID (%s), so skip this auth deletion", apiId)
			continue
		}
		log.Printf("[DEBUG] Prepare to delete the authorization ID (%s) via API ID (%s)", authId, apiId)
		err = deleteApplicationAuthorizationForApi(client, instanceId, envId, appId, authId)
		if err != nil {
			mErr = multierror.Append(mErr, err)
		}
	}

	return mErr.ErrorOrNil()
}

func resourceApplicationAuthorizationImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<instance_id>/<env_id>/<application_id>', "+
			"but got '%s'", importedId)
	}

	mErr := multierror.Append(nil,
		d.Set("instance_id", parts[0]),
		d.Set("env_id", parts[1]),
		d.Set("application_id", parts[2]),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return []*schema.ResourceData{d},
			fmt.Errorf("error saving application authorization resource fields during import: %s", err)
	}
	return []*schema.ResourceData{d}, nil
}
