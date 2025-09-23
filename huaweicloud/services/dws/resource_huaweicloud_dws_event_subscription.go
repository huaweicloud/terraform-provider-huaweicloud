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

// @API DWS POST /v2/{project_id}/event-subs
// @API DWS GET /v2/{project_id}/event-subs
// @API AWS PUT /v2/{project_id}/event-subs/{event_sub_id}
// @API DWS DELETE /v2/{project_id}/event-subs/{event_sub_id}
func ResourceDwsEventSubs() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDwsEventSubsCreate,
		UpdateContext: resourceDwsEventSubsUpdate,
		ReadContext:   resourceDwsEventSubsRead,
		DeleteContext: resourceDwsEventSubsDelete,
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
				Description: `The name of the event subscription.`,
			},
			"enable": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `whether the event subscription is enabled.`,
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
			"source_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `ID of source event.`,
			},
			"source_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The type of source event.`,
			},
			"category": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The category of source event.`,
			},
			"severity": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The severity of source event.`,
			},
			"time_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The time_zone of the event subscription.`,
			},
		},
	}
}

func resourceDwsEventSubsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createDwsEventSubs: create a DWS event subscription.
	var (
		createDwsEventSubsHttpUrl = "v2/{project_id}/event-subs"
		createDwsEventSubsProduct = "dws"
	)
	createDwsEventSubsClient, err := cfg.NewServiceClient(createDwsEventSubsProduct, region)
	if err != nil {
		return diag.Errorf("error creating DWS Client: %s", err)
	}

	createDwsEventSubsPath := createDwsEventSubsClient.Endpoint + createDwsEventSubsHttpUrl
	createDwsEventSubsPath = strings.ReplaceAll(createDwsEventSubsPath, "{project_id}", createDwsEventSubsClient.ProjectID)

	createDwsEventSubsOpt := golangsdk.RequestOpts{
		MoreHeaders:      requestOpts.MoreHeaders,
		KeepResponseBody: true,
	}
	createDwsEventSubsOpt.JSONBody = utils.RemoveNil(buildCreateDwsEventSubsBodyParams(d))
	createDwsEventSubsResp, err := createDwsEventSubsClient.Request("POST", createDwsEventSubsPath, &createDwsEventSubsOpt)
	if err != nil {
		return diag.Errorf("error creating DWS event subscription: %s", err)
	}

	createDwsEventSubsRespBody, err := utils.FlattenResponse(createDwsEventSubsResp)
	if err != nil {
		return diag.FromErr(err)
	}

	subscriptionId := utils.PathSearch("id", createDwsEventSubsRespBody, "").(string)
	if subscriptionId == "" {
		return diag.Errorf("unable to find the DWS event subscription ID from the API response")
	}
	d.SetId(subscriptionId)

	return resourceDwsEventSubsRead(ctx, d, meta)
}

func buildCreateDwsEventSubsBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":                     utils.ValueIgnoreEmpty(d.Get("name")),
		"source_id":                utils.ValueIgnoreEmpty(d.Get("source_id")),
		"source_type":              utils.ValueIgnoreEmpty(d.Get("source_type")),
		"category":                 utils.ValueIgnoreEmpty(d.Get("category")),
		"severity":                 utils.ValueIgnoreEmpty(d.Get("severity")),
		"enable":                   utils.StringToInt(utils.String(d.Get("enable").(string))),
		"notification_target":      utils.ValueIgnoreEmpty(d.Get("notification_target")),
		"notification_target_name": utils.ValueIgnoreEmpty(d.Get("notification_target_name")),
		"notification_target_type": utils.ValueIgnoreEmpty(d.Get("notification_target_type")),
		"time_zone":                utils.ValueIgnoreEmpty(d.Get("time_zone")),
	}
	return bodyParams
}

func resourceDwsEventSubsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getDwsEventSubs: Query the DWS event subscription.
	var (
		getDwsEventSubsHttpUrl = "v2/{project_id}/event-subs"
		getDwsEventSubsProduct = "dws"
	)
	getDwsEventSubsClient, err := cfg.NewServiceClient(getDwsEventSubsProduct, region)
	if err != nil {
		return diag.Errorf("error creating DWS Client: %s", err)
	}

	getDwsEventSubsPath := getDwsEventSubsClient.Endpoint + getDwsEventSubsHttpUrl
	getDwsEventSubsPath = strings.ReplaceAll(getDwsEventSubsPath, "{project_id}", getDwsEventSubsClient.ProjectID)

	getDwsEventSubsResp, err := pagination.ListAllItems(
		getDwsEventSubsClient,
		"offset",
		getDwsEventSubsPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DwsEventSubs")
	}

	getDwsEventSubsRespJson, err := json.Marshal(getDwsEventSubsResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var getDwsEventSubsRespBody interface{}
	err = json.Unmarshal(getDwsEventSubsRespJson, &getDwsEventSubsRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonPath := fmt.Sprintf("event_subscriptions[?id=='%s']|[0]", d.Id())
	rawData := utils.PathSearch(jsonPath, getDwsEventSubsRespBody, nil)
	if rawData == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "no data found")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", rawData, nil)),
		d.Set("source_id", utils.PathSearch("source_id", rawData, nil)),
		d.Set("source_type", utils.PathSearch("source_type", rawData, nil)),
		d.Set("category", utils.PathSearch("category", rawData, nil)),
		d.Set("severity", utils.PathSearch("severity", rawData, nil)),
		d.Set("enable", fmt.Sprint(utils.PathSearch("enable", rawData, nil))),
		d.Set("notification_target", utils.PathSearch("notification_target", rawData, nil)),
		d.Set("notification_target_name", utils.PathSearch("notification_target_name", rawData, nil)),
		d.Set("notification_target_type", utils.PathSearch("notification_target_type", rawData, nil)),
		d.Set("time_zone", utils.PathSearch("time_zone", rawData, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDwsEventSubsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateDwsEventSubsChanges := []string{
		"name",
		"source_id",
		"source_type",
		"category",
		"severity",
		"enable",
		"notification_target",
		"notification_target_name",
		"notification_target_type",
	}

	if d.HasChanges(updateDwsEventSubsChanges...) {
		// updateDwsEventSubs: update the DWS event subscription.
		var (
			updateDwsEventSubsHttpUrl = "v2/{project_id}/event-subs/{event_sub_id}"
			updateDwsEventSubsProduct = "dws"
		)
		updateDwsEventSubsClient, err := cfg.NewServiceClient(updateDwsEventSubsProduct, region)
		if err != nil {
			return diag.Errorf("error creating AS Client: %s", err)
		}

		updateDwsEventSubsPath := updateDwsEventSubsClient.Endpoint + updateDwsEventSubsHttpUrl
		updateDwsEventSubsPath = strings.ReplaceAll(updateDwsEventSubsPath, "{project_id}", updateDwsEventSubsClient.ProjectID)
		updateDwsEventSubsPath = strings.ReplaceAll(updateDwsEventSubsPath, "{event_sub_id}", d.Id())

		updateDwsEventSubsOpt := golangsdk.RequestOpts{
			MoreHeaders:      requestOpts.MoreHeaders,
			KeepResponseBody: true,
		}
		updateDwsEventSubsOpt.JSONBody = utils.RemoveNil(buildUpdateDwsEventSubsBodyParams(d))
		_, err = updateDwsEventSubsClient.Request("PUT", updateDwsEventSubsPath, &updateDwsEventSubsOpt)
		if err != nil {
			return diag.Errorf("error updating DWS event subscription: %s", err)
		}
	}
	return resourceDwsEventSubsRead(ctx, d, meta)
}

func buildUpdateDwsEventSubsBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":                     utils.ValueIgnoreEmpty(d.Get("name")),
		"source_id":                utils.ValueIgnoreEmpty(d.Get("source_id")),
		"source_type":              utils.ValueIgnoreEmpty(d.Get("source_type")),
		"category":                 utils.ValueIgnoreEmpty(d.Get("category")),
		"severity":                 utils.ValueIgnoreEmpty(d.Get("severity")),
		"enable":                   utils.StringToInt(utils.String(d.Get("enable").(string))),
		"notification_target":      utils.ValueIgnoreEmpty(d.Get("notification_target")),
		"notification_target_name": utils.ValueIgnoreEmpty(d.Get("notification_target_name")),
		"notification_target_type": utils.ValueIgnoreEmpty(d.Get("notification_target_type")),
	}
	return bodyParams
}

func resourceDwsEventSubsDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteDwsEventSubs: delete DWS event subscription
	var (
		deleteDwsEventSubsHttpUrl = "v2/{project_id}/event-subs/{event_sub_id}"
		deleteDwsEventSubsProduct = "dws"
	)
	deleteDwsEventSubsClient, err := cfg.NewServiceClient(deleteDwsEventSubsProduct, region)
	if err != nil {
		return diag.Errorf("error creating DWS Client: %s", err)
	}

	deleteDwsEventSubsPath := deleteDwsEventSubsClient.Endpoint + deleteDwsEventSubsHttpUrl
	deleteDwsEventSubsPath = strings.ReplaceAll(deleteDwsEventSubsPath, "{project_id}", deleteDwsEventSubsClient.ProjectID)
	deleteDwsEventSubsPath = strings.ReplaceAll(deleteDwsEventSubsPath, "{event_sub_id}", d.Id())

	deleteDwsEventSubsOpt := golangsdk.RequestOpts{
		MoreHeaders:      requestOpts.MoreHeaders,
		KeepResponseBody: true,
	}
	_, err = deleteDwsEventSubsClient.Request("DELETE", deleteDwsEventSubsPath, &deleteDwsEventSubsOpt)
	if err != nil {
		return diag.Errorf("error deleting DWS event subscription: %s", err)
	}

	return nil
}
