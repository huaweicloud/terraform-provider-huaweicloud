package apig

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
	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/appauths"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API APIG DELETE /v2/{project_id}/apigw/instances/{instance_id}/app-auths/{app_auth_id}
// @API APIG POST /v2/{project_id}/apigw/instances/{instance_id}/app-auths
// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/app-auths/binded-apis
// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/app-auths/unbinded-apis
func ResourceAppAuth() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppAuthCreate,
		ReadContext:   resourceAppAuthRead,
		UpdateContext: resourceAppAuthUpdate,
		DeleteContext: resourceAppAuthDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceAppAuthImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Update: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
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
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The authorized API IDs",
			},
		},
	}
}

func resourceAppAuthCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.ApigV2Client(region)
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	var (
		appId  = d.Get("application_id").(string)
		envId  = d.Get("env_id").(string)
		apiIds = utils.ExpandToStringListBySet(d.Get("api_ids").(*schema.Set))
	)
	err = createAppAuthForApis(ctx, client, d, apiIds, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprintf("%s/%s", envId, appId))

	return resourceAppAuthRead(ctx, d, meta)
}

func flattenAuthorizedApis(apiInfos []appauths.ApiAuthInfo) []string {
	result := make([]string, len(apiInfos))
	for i, val := range apiInfos {
		result[i] = val.ApiId
	}
	return result
}

func resourceAppAuthRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		appId      = d.Get("application_id").(string)
		opts       = appauths.ListOpts{
			InstanceId: instanceId,
			AppId:      appId,
		}
	)
	resp, err := appauths.ListAuthorized(client, opts)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error querying authorized APIs from application (%s) under dedicated instance (%s)",
			appId, instanceId))
	}
	if len(resp) < 1 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "Application Authorization")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("api_ids", flattenAuthorizedApis(resp)),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving authorization fields for specified application (%s): %s", appId, err)
	}
	return nil
}

func resourceAppAuthUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	var (
		oldRaw, newRaw = d.GetChange("api_ids")

		addSet = newRaw.(*schema.Set).Difference(oldRaw.(*schema.Set))
		rmSet  = oldRaw.(*schema.Set).Difference(newRaw.(*schema.Set))
	)
	if rmSet.Len() > 0 {
		apiIds := utils.ExpandToStringListBySet(rmSet)
		err := deleteAppAuthFromApis(ctx, client, d, apiIds, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if addSet.Len() > 0 {
		apiIds := utils.ExpandToStringListBySet(addSet)
		err = createAppAuthForApis(ctx, client, d, apiIds, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceAppAuthRead(ctx, d, meta)
}

func resourceAppAuthDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.ApigV2Client(region)
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	apiIds := d.Get("api_ids").(*schema.Set)
	if err = deleteAppAuthFromApis(ctx, client, d, utils.ExpandToStringListBySet(apiIds), d.Timeout(schema.TimeoutDelete)); err != nil {
		return common.CheckDeletedDiag(d, err, "error unbinding APIs from application")
	}
	return nil
}

func createAppAuthForApis(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData, apiIds []string,
	timeout time.Duration) error {
	var (
		instanceId = d.Get("instance_id").(string)
		appId      = d.Get("application_id").(string)
		envId      = d.Get("env_id").(string)
		opts       = appauths.CreateOpts{
			InstanceId: instanceId,
			AppIds:     []string{appId},
			EnvId:      envId,
			ApiIds:     apiIds,
		}
	)

	_, err := appauths.Create(client, opts)
	if err != nil {
		return fmt.Errorf("error authorizing APIs to application (%s) under dedicated instance (%s): %s", appId, instanceId, err)
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: authApisStateRefreshFunc(client, instanceId, appId, envId, apiIds),
		Timeout: timeout,
		// In most cases, the unbind operation will be completed immediately, but in a few cases, it needs to wait
		// for a short period of time, and the polling is performed by incrementing the time here.
		MinTimeout: 2 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for API authorize operations completed: %s", err)
	}

	return nil
}

func deleteAppAuthFromApis(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData, apiIds []string,
	timeout time.Duration) error {
	var (
		instanceId = d.Get("instance_id").(string)
		appId      = d.Get("application_id").(string)
		envId      = d.Get("env_id").(string)
		opts       = appauths.ListOpts{
			InstanceId: instanceId,
			AppId:      appId,
			EnvId:      envId,
		}
		notFoundErr = fmt.Sprintf("[DEBUG] All APIs have been unauthorized form application (%s) under dedicated instance (%s)", appId, instanceId)
	)

	resp, err := appauths.ListAuthorized(client, opts)
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			log.Println(notFoundErr)
			return err
		}
		return fmt.Errorf("error querying authorized APIs for application (%s) under dedicated instance (%s)", appId, instanceId)
	}
	if len(resp) < 1 {
		log.Println(notFoundErr)
		return nil
	}

	for _, val := range resp {
		if !utils.StrSliceContains(apiIds, val.ApiId) {
			continue
		}

		authId := val.ID
		err = appauths.Delete(client, instanceId, authId)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] All APIs has been unauthorized from the application (%s)", appId)
				continue
			}
			return fmt.Errorf("error unauthorizing APIs form application (%s) under dedicated instance (%s): %s", appId, instanceId, err)
		}
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: unauthApisStateRefreshFunc(client, instanceId, appId, envId, apiIds),
		Timeout: timeout,
		// In most cases, the unbind operation will be completed immediately, but in a few cases, it needs to wait
		// for a short period of time, and the polling is performed by incrementing the time here.
		MinTimeout: 2 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for API unauthorize operations completed: %s", err)
	}
	return nil
}

func authApisStateRefreshFunc(client *golangsdk.ServiceClient, instanceId, appId, envId string, apiIds []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		opts := appauths.ListOpts{
			InstanceId: instanceId,
			AppId:      appId,
			EnvId:      envId,
		}
		resp, err := appauths.ListUnauthorized(client, opts)
		if err != nil {
			return resp, "", err
		}

		for _, val := range resp {
			if utils.StrSliceContains(apiIds, val.ID) {
				return resp, "PENDING", nil
			}
		}
		return resp, "COMPLETED", nil
	}
}

func unauthApisStateRefreshFunc(client *golangsdk.ServiceClient, instanceId, appId, envId string, apiIds []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		opts := appauths.ListOpts{
			InstanceId: instanceId,
			AppId:      appId,
			EnvId:      envId,
		}
		resp, err := appauths.ListAuthorized(client, opts)
		if err != nil {
			// The API returns a 404 error, which means that the instance or application has been deleted.
			// In this case, there's no need to disassociate API, also this action has been completed.
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return "instance_or_application_not_exist", "COMPLETED", nil
			}
			return resp, "", err
		}

		for _, val := range resp {
			if utils.StrSliceContains(apiIds, val.ApiId) {
				return resp, "PENDING", nil
			}
		}
		return resp, "COMPLETED", nil
	}
}

func resourceAppAuthImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<instance_id>/<id>' (the format of resource ID is "+
			"'<env_id>/<application_id>'), but got '%s'", importedId)
	}

	d.SetId(fmt.Sprintf("%s/%s", parts[1], parts[2]))
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
