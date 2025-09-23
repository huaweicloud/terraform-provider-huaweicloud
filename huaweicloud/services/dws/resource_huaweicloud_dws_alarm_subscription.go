// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product DWS
// ---------------------------------------------------------------

package dws

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DWS GET /v2/{project_id}/alarm-subs
// @API DWS POST /v2/{project_id}/alarm-subs
// @API DWS PUT /v2/{project_id}/alarm-subs/{alarm_sub_id}
// @API DWS DELETE /v2/{project_id}/alarm-subs/{alarm_sub_id}
func ResourceDwsAlarmSubs() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDwsAlarmSubsCreate,
		UpdateContext: resourceDwsAlarmSubsUpdate,
		ReadContext:   resourceDwsAlarmSubsRead,
		DeleteContext: resourceDwsAlarmSubsDelete,
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
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the alarm subscription.`,
			},
			"enable": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `whether the alarm subscription is enabled.`,
			},
			"notification_target": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The notification target.`,
			},
			"notification_target_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of notification target.`,
			},
			"notification_target_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The type of notification target. Currently only **SMN** is supported.`,
			},
			"time_zone": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The time_zone of the alarm subscription.`,
			},
			"alarm_level": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: `The level of alarm. separate multiple alarm levels with commas (,).
`,
			},
		},
	}
}

func resourceDwsAlarmSubsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createDwsAlarmSubs: create a DWS alarm subscription.
	var (
		createDwsAlarmSubsHttpUrl = "v2/{project_id}/alarm-subs"
		createDwsAlarmSubsProduct = "dws"
	)
	createDwsAlarmSubsClient, err := cfg.NewServiceClient(createDwsAlarmSubsProduct, region)
	if err != nil {
		return diag.Errorf("error creating DWS Client: %s", err)
	}

	createDwsAlarmSubsPath := createDwsAlarmSubsClient.Endpoint + createDwsAlarmSubsHttpUrl
	createDwsAlarmSubsPath = strings.ReplaceAll(createDwsAlarmSubsPath, "{project_id}", createDwsAlarmSubsClient.ProjectID)

	createDwsAlarmSubsOpt := golangsdk.RequestOpts{
		MoreHeaders:      requestOpts.MoreHeaders,
		KeepResponseBody: true,
	}
	createDwsAlarmSubsOpt.JSONBody = utils.RemoveNil(buildCreateDwsAlarmSubsBodyParams(d))
	createDwsAlarmSubsResp, err := createDwsAlarmSubsClient.Request("POST", createDwsAlarmSubsPath, &createDwsAlarmSubsOpt)
	if err != nil {
		return diag.Errorf("error creating DWS alarm subscription: %s", err)
	}

	createDwsAlarmSubsRespBody, err := utils.FlattenResponse(createDwsAlarmSubsResp)
	if err != nil {
		return diag.FromErr(err)
	}

	subscriptionId := utils.PathSearch("id", createDwsAlarmSubsRespBody, "").(string)
	if subscriptionId == "" {
		return diag.Errorf("unable to find the DWS alarm subscription ID from the API response")
	}
	d.SetId(subscriptionId)

	return resourceDwsAlarmSubsRead(ctx, d, meta)
}

func buildCreateDwsAlarmSubsBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":                     utils.ValueIgnoreEmpty(d.Get("name")),
		"alarm_level":              utils.ValueIgnoreEmpty(d.Get("alarm_level")),
		"enable":                   d.Get("enable"),
		"notification_target":      utils.ValueIgnoreEmpty(d.Get("notification_target")),
		"notification_target_name": utils.ValueIgnoreEmpty(d.Get("notification_target_name")),
		"notification_target_type": utils.ValueIgnoreEmpty(d.Get("notification_target_type")),
		"time_zone":                utils.ValueIgnoreEmpty(d.Get("time_zone")),
	}
	return bodyParams
}

func resourceDwsAlarmSubsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getDwsAlarmSubs: Query the DWS alarm subscription.
	var (
		getDwsAlarmSubsHttpUrl = "v2/{project_id}/alarm-subs"
		getDwsAlarmSubsProduct = "dws"
	)
	getDwsAlarmSubsClient, err := cfg.NewServiceClient(getDwsAlarmSubsProduct, region)
	if err != nil {
		return diag.Errorf("error creating DWS Client: %s", err)
	}

	getDwsAlarmSubsPath := getDwsAlarmSubsClient.Endpoint + getDwsAlarmSubsHttpUrl
	getDwsAlarmSubsPath = strings.ReplaceAll(getDwsAlarmSubsPath, "{project_id}", getDwsAlarmSubsClient.ProjectID)

	getDwsAlarmSubsResp, err := pagination.ListAllItems(
		getDwsAlarmSubsClient,
		"offset",
		getDwsAlarmSubsPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DWS alarm subscription")
	}

	getDwsAlarmSubsRespJson, err := json.Marshal(getDwsAlarmSubsResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var getDwsAlarmSubsRespBody interface{}
	err = json.Unmarshal(getDwsAlarmSubsRespJson, &getDwsAlarmSubsRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonPath := fmt.Sprintf("alarm_subscriptions[?id=='%s']|[0]", d.Id())
	rawData := utils.PathSearch(jsonPath, getDwsAlarmSubsRespBody, nil)
	if rawData == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "no data found")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", rawData, nil)),
		d.Set("alarm_level", utils.PathSearch("alarm_level", rawData, nil)),
		d.Set("enable", utils.PathSearch("enable", rawData, nil)),
		d.Set("notification_target", utils.PathSearch("notification_target", rawData, nil)),
		d.Set("notification_target_name", utils.PathSearch("notification_target_name", rawData, nil)),
		d.Set("notification_target_type", utils.PathSearch("notification_target_type", rawData, nil)),
		d.Set("time_zone", utils.PathSearch("time_zone", rawData, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDwsAlarmSubsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateDwsAlarmSubsChanges := []string{
		"name",
		"alarm_level",
		"enable",
		"notification_target",
		"notification_target_name",
		"notification_target_type",
	}

	if d.HasChanges(updateDwsAlarmSubsChanges...) {
		// updateDwsAlarmSubs: update the DWS alarm subscription.
		var (
			updateDwsAlarmSubsHttpUrl = "v2/{project_id}/alarm-subs/{alarm_sub_id}"
			updateDwsAlarmSubsProduct = "dws"
		)
		updateDwsAlarmSubsClient, err := cfg.NewServiceClient(updateDwsAlarmSubsProduct, region)
		if err != nil {
			return diag.Errorf("error creating DWS Client: %s", err)
		}

		updateDwsAlarmSubsPath := updateDwsAlarmSubsClient.Endpoint + updateDwsAlarmSubsHttpUrl
		updateDwsAlarmSubsPath = strings.ReplaceAll(updateDwsAlarmSubsPath, "{project_id}", updateDwsAlarmSubsClient.ProjectID)
		updateDwsAlarmSubsPath = strings.ReplaceAll(updateDwsAlarmSubsPath, "{alarm_sub_id}", d.Id())

		updateDwsAlarmSubsOpt := golangsdk.RequestOpts{
			MoreHeaders:      requestOpts.MoreHeaders,
			KeepResponseBody: true,
		}
		updateDwsAlarmSubsOpt.JSONBody = utils.RemoveNil(buildUpdateDwsAlarmSubsBodyParams(d))
		_, err = updateDwsAlarmSubsClient.Request("PUT", updateDwsAlarmSubsPath, &updateDwsAlarmSubsOpt)
		if err != nil {
			return diag.Errorf("error updating DWS alarm subscription: %s", err)
		}
	}
	return resourceDwsAlarmSubsRead(ctx, d, meta)
}

func buildUpdateDwsAlarmSubsBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":                     utils.ValueIgnoreEmpty(d.Get("name")),
		"alarm_level":              utils.ValueIgnoreEmpty(d.Get("alarm_level")),
		"enable":                   d.Get("enable"),
		"notification_target":      utils.ValueIgnoreEmpty(d.Get("notification_target")),
		"notification_target_name": utils.ValueIgnoreEmpty(d.Get("notification_target_name")),
		"notification_target_type": utils.ValueIgnoreEmpty(d.Get("notification_target_type")),
	}
	return bodyParams
}

func resourceDwsAlarmSubsDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteDwsAlarmSubs: delete DWS alarm subscription
	var (
		deleteDwsAlarmSubsHttpUrl = "v2/{project_id}/alarm-subs/{alarm_sub_id}"
		deleteDwsAlarmSubsProduct = "dws"
	)
	deleteDwsAlarmSubsClient, err := cfg.NewServiceClient(deleteDwsAlarmSubsProduct, region)
	if err != nil {
		return diag.Errorf("error creating DWS Client: %s", err)
	}

	deleteDwsAlarmSubsPath := deleteDwsAlarmSubsClient.Endpoint + deleteDwsAlarmSubsHttpUrl
	deleteDwsAlarmSubsPath = strings.ReplaceAll(deleteDwsAlarmSubsPath, "{project_id}", deleteDwsAlarmSubsClient.ProjectID)
	deleteDwsAlarmSubsPath = strings.ReplaceAll(deleteDwsAlarmSubsPath, "{alarm_sub_id}", d.Id())

	deleteDwsAlarmSubsOpt := golangsdk.RequestOpts{
		MoreHeaders:      requestOpts.MoreHeaders,
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	_, err = deleteDwsAlarmSubsClient.Request("DELETE", deleteDwsAlarmSubsPath, &deleteDwsAlarmSubsOpt)
	if err != nil {
		return diag.Errorf("error deleting DWS alarm subscription: %s", err)
	}

	return nil
}
