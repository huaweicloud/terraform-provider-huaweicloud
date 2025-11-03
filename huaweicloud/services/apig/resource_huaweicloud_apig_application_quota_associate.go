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

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API APIG POST /v2/{project_id}/apigw/instances/{instance_id}/app-quotas/{app_quota_id}/binding-apps
// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/app-quotas/{app_quota_id}/bindable-apps
// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/app-quotas/{app_quota_id}/bound-apps
// @API APIG DELETE /v2/{project_id}/apigw/instances/{instance_id}/app-quotas/{app_quota_id}/bound-apps/{app_id}
func ResourceApplicationQuotaAssociate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApplicationQuotaAssociateCreate,
		ReadContext:   resourceApplicationQuotaAssociateRead,
		UpdateContext: resourceApplicationQuotaAssociateUpdate,
		DeleteContext: resourceApplicationQuotaAssociateDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceApplicationQuotaAssociateImportState,
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
				Description: "The region where the application quota (policy) is located.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the dedicated instance to which the application quota (policy) belongs.",
			},
			"quota_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the application quota (policy).",
			},
			"applications": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The application ID bound to the application quota.",
						},
						"bind_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The binding time, in RFC3339 format.",
						},
					},
				},
				Description: "The configuration of applications bound to the quota.",
			},
		},
	}
}

func parseAssociatedAppIds(applications *schema.Set) []string {
	result := make([]string, 0, applications.Len())
	for _, val := range applications.List() {
		app := val.(map[string]interface{})
		result = append(result, app["id"].(string))
	}
	return result
}

func bindableAppRefreshFunc(client *golangsdk.ServiceClient, instanceId, quotaId string, appIds []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var (
			httpUrl = "v2/{project_id}/apigw/instances/{instance_id}/app-quotas/{app_quota_id}/bindable-apps"
		)

		queryPath := client.Endpoint + httpUrl
		queryPath = strings.ReplaceAll(queryPath, "{project_id}", client.ProjectID)
		queryPath = strings.ReplaceAll(queryPath, "{instance_id}", instanceId)
		queryPath = strings.ReplaceAll(queryPath, "{app_quota_id}", quotaId)
		opt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		requestResp, err := client.Request("GET", queryPath, &opt)
		if err != nil {
			return requestResp, "ERROR", fmt.Errorf("error querying bindable applications by specified quota ID (%s): %s", quotaId, err)
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return respBody, "ERROR", err
		}

		if utils.StrSliceContainsAnother(utils.ExpandToStringList(utils.PathSearch("apps[*].app_id", respBody,
			make([]interface{}, 0)).([]interface{})), appIds) {
			return "matched", "COMPLETED", nil
		}
		return "continue", "PENDING", nil
	}
}

func precheckAllAppIdsAreAvailable(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData, appIds []string,
	timeout time.Duration) error {
	var (
		instanceId = d.Get("instance_id").(string)
		quotaId    = d.Get("quota_id").(string)
	)
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: bindableAppRefreshFunc(client, instanceId, quotaId, appIds),
		Timeout: timeout,
		// In most cases, the bind operation will be completed immediately, but in a few cases, it needs to wait
		// for a short period of time, and the polling is performed by incrementing the time here.
		MinTimeout: 2 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for the application pre-check completed: %s", err)
	}
	return nil
}

func QueryQuotaAssociatedApplications(client *golangsdk.ServiceClient, instanceId, quotaId string) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/apigw/instances/{instance_id}/app-quotas/{app_quota_id}/bound-apps"
		offset  = 0
		result  = make([]interface{}, 0)
	)
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", instanceId)
	listPath = strings.ReplaceAll(listPath, "{app_quota_id}", quotaId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s?limit=100&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		unbindPublishIds := utils.PathSearch("apps", respBody, make([]interface{}, 0)).([]interface{})
		if len(unbindPublishIds) < 1 {
			break
		}
		result = append(result, unbindPublishIds...)
		offset += len(unbindPublishIds)
	}
	return result, nil
}

func appsBindingRefreshFunc(client *golangsdk.ServiceClient, instanceId, quota string,
	appIds []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		respBody, err := QueryQuotaAssociatedApplications(client, instanceId, quota)
		if err != nil {
			return nil, "ERROR", err
		}

		if utils.StrSliceContainsAnother(utils.ExpandToStringList(utils.PathSearch("[*].app_id",
			respBody, make([]interface{}, 0)).([]interface{})), appIds) {
			return respBody, "COMPLETED", nil
		}
		return "continue", "PENDING", nil
	}
}

func associateAppsToQuota(ctx context.Context, client *golangsdk.ServiceClient, instanceId, quotaId string,
	appIds []string, timeout time.Duration) error {
	httpUrl := "v2/{project_id}/apigw/instances/{instance_id}/app-quotas/{app_quota_id}/binding-apps"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)
	createPath = strings.ReplaceAll(createPath, "{app_quota_id}", quotaId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"app_ids": appIds,
		},
	}

	_, err := client.Request("POST", createPath, &opt)
	if err != nil {
		return fmt.Errorf("failed to associate application(s) to the application quota (%s): %s", quotaId, err)
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: appsBindingRefreshFunc(client, instanceId, quotaId, appIds),
		Timeout: timeout,
		// In most cases, the bind operation will be completed immediately, but in a few cases, it needs to wait
		// for a short period of time, and the polling is performed by incrementing the time here.
		MinTimeout: 2 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for the binding completed: %s", err)
	}
	return nil
}

func resourceApplicationQuotaAssociateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
		quotaId    = d.Get("quota_id").(string)
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	err = precheckAllAppIdsAreAvailable(ctx, client, d, parseAssociatedAppIds(d.Get("applications").(*schema.Set)),
		d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}
	err = associateAppsToQuota(ctx, client, instanceId, quotaId, parseAssociatedAppIds(d.Get("applications").(*schema.Set)),
		d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(quotaId)

	return resourceApplicationQuotaAssociateRead(ctx, d, meta)
}

func flattenQuotaAssociatedApplications(associatedApps []interface{}) []interface{} {
	if len(associatedApps) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(associatedApps))
	for _, associatedApp := range associatedApps {
		result = append(result, map[string]interface{}{
			"id": utils.PathSearch("app_id", associatedApp, nil),
			// Use the time zone configured by the execution machine as the basis for parsing its value.
			"bind_time": utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("bound_time",
				associatedApp, "").(string))/1000, false),
		})
	}
	return result
}

func resourceApplicationQuotaAssociateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
		quotaId    = d.Get("quota_id").(string)
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	respBody, err := QueryQuotaAssociatedApplications(client, instanceId, quotaId)
	if err != nil {
		return common.CheckDeletedDiag(d, err,
			fmt.Sprintf("error retrieving the associated application(s) from the application quota (%s)", quotaId))
	}
	associatedApps := flattenQuotaAssociatedApplications(utils.PathSearch("[*]", respBody, make([]interface{}, 0)).([]interface{}))
	if len(associatedApps) < 1 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "all APPs have been unbound")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("applications", associatedApps),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving the fields of the associated application(s): %s", err)
	}
	return nil
}

func disassociateAppsFromQuota(ctx context.Context, client *golangsdk.ServiceClient, instanceId, quotaId string,
	appIds []string, timeout time.Duration) error {
	var (
		httpUrl  = "v2/{project_id}/apigw/instances/{instance_id}/app-quotas/{app_quota_id}/bound-apps/{app_id}"
		basePath = client.Endpoint + httpUrl
	)
	basePath = strings.ReplaceAll(basePath, "{project_id}", client.ProjectID)
	basePath = strings.ReplaceAll(basePath, "{instance_id}", instanceId)
	basePath = strings.ReplaceAll(basePath, "{app_quota_id}", quotaId)

	for _, appId := range appIds {
		deletePath := strings.ReplaceAll(basePath, "{app_id}", appId)
		opt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		_, err := client.Request("DELETE", deletePath, &opt)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] The application (%s) has been disassociated from the application quota", appId)
				continue
			}
			return fmt.Errorf("failed to disassociate application from the application quota (%s): %s", quotaId, err)
		}
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: appsUnbindingRefreshFunc(client, instanceId, quotaId, appIds),
		Timeout: timeout,
		// In most cases, the bind operation will be completed immediately, but in a few cases, it needs to wait
		// for a short period of time, and the polling is performed by incrementing the time here.
		MinTimeout: 2 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for the unbinding completed: %s", err)
	}
	return nil
}

func appsUnbindingRefreshFunc(client *golangsdk.ServiceClient, instanceId, quota string,
	appIds []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		respBody, err := QueryQuotaAssociatedApplications(client, instanceId, quota)
		if err != nil {
			// The API returns a 404 error, which means that the instance has been deleted.
			// In this case, there's no need to disassociate quota and application, also this action has been completed.
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return "instance_not_exist", "COMPLETED", nil
			}
			return nil, "ERROR", err
		}

		if !utils.IsSliceContainsAnyAnotherSliceElement(utils.ExpandToStringList(utils.PathSearch("[*].app_id",
			respBody, make([]interface{}, 0)).([]interface{})), appIds, false, false) {
			return respBody, "COMPLETED", nil
		}
		return "continue", "PENDING", nil
	}
}

func resourceApplicationQuotaAssociateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	var (
		instanceId     = d.Get("instance_id").(string)
		quotaId        = d.Get("quota_id").(string)
		oldVal, newVal = d.GetChange("applications")
		rmRaw          = oldVal.(*schema.Set).Difference(newVal.(*schema.Set))
		addRaw         = newVal.(*schema.Set).Difference(oldVal.(*schema.Set))
	)

	if rmRaw.Len() > 0 {
		if err = disassociateAppsFromQuota(ctx, client, instanceId, quotaId, parseAssociatedAppIds(rmRaw),
			d.Timeout(schema.TimeoutUpdate)); err != nil {
			return diag.FromErr(err)
		}
	}

	if addRaw.Len() > 0 {
		err = precheckAllAppIdsAreAvailable(ctx, client, d, parseAssociatedAppIds(addRaw), d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.FromErr(err)
		}
		if err = associateAppsToQuota(ctx, client, instanceId, quotaId, parseAssociatedAppIds(addRaw),
			d.Timeout(schema.TimeoutUpdate)); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceApplicationQuotaAssociateRead(ctx, d, meta)
}

func resourceApplicationQuotaAssociateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
		quotaId    = d.Get("quota_id").(string)
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	oldVal, _ := d.GetChange("applications")
	if err = disassociateAppsFromQuota(ctx, client, instanceId, quotaId, parseAssociatedAppIds(oldVal.(*schema.Set)),
		d.Timeout(schema.TimeoutDelete)); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceApplicationQuotaAssociateImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<instance_id>/<quota_id>', but got '%s'", importedId)
	}

	d.SetId(parts[1])
	mErr := multierror.Append(nil,
		d.Set("instance_id", parts[0]),
		d.Set("quota_id", parts[1]),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return []*schema.ResourceData{d},
			fmt.Errorf("error saving the fields of the associate configuration during import: %s", err)
	}
	return []*schema.ResourceData{d}, nil
}
