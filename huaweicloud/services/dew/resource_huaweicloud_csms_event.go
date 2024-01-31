// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product DEW
// ---------------------------------------------------------------

package dew

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DEW POST /v1/{project_id}/csms/events
// @API DEW GET /v1/{project_id}/csms/events/{event_name}
// @API DEW PUT /v1/{project_id}/csms/events/{event_name}
// @API DEW DELETE /v1/{project_id}/csms/events/{event_name}
func ResourceCsmsEvent() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCsmsEventCreate,
		UpdateContext: resourceCsmsEventUpdate,
		ReadContext:   resourceCsmsEventRead,
		DeleteContext: resourceCsmsEventDelete,
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
				ForceNew:    true,
				Description: `Specifies the name of CSMS event.`,
			},
			"event_types": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				Description: `Specifies the event list.`,
			},
			"status": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the event status.`,
			},
			"notification_target_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the notification target type.`,
			},
			"notification_target_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the notification target ID.`,
			},
			"notification_target_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the notification target name.`,
			},
			"event_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The event ID.`,
			},
		},
	}
}

func resourceCsmsEventCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/csms/events"
		product = "kms"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating KMS client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildCreateCsmsEventBodyParams(d),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating CSMS event: %s", err)
	}

	name := d.Get("name").(string)
	d.SetId(name)

	return resourceCsmsEventRead(ctx, d, meta)
}

func buildCreateCsmsEventBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":         d.Get("name"),
		"event_types":  d.Get("event_types"),
		"state":        d.Get("status"),
		"notification": buildNotificationBodyParam(d),
	}
	return bodyParams
}

func buildNotificationBodyParam(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"target_type": d.Get("notification_target_type"),
		"target_id":   d.Get("notification_target_id"),
		"target_name": d.Get("notification_target_name"),
	}
}

func resourceCsmsEventRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		mErr    *multierror.Error
		httpUrl = "v1/{project_id}/csms/events/{event_name}"
		product = "kms"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating KMS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{event_name}", d.Id())
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CSMS event")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("event.name", getRespBody, nil)),
		d.Set("event_types", utils.PathSearch("event.event_types", getRespBody, nil)),
		d.Set("status", utils.PathSearch("event.state", getRespBody, nil)),
		d.Set("notification_target_type", utils.PathSearch("event.notification.target_type", getRespBody, nil)),
		d.Set("notification_target_id", utils.PathSearch("event.notification.target_id", getRespBody, nil)),
		d.Set("notification_target_name", utils.PathSearch("event.notification.target_name", getRespBody, nil)),
		d.Set("event_id", utils.PathSearch("event.event_id", getRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceCsmsEventUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/csms/events/{event_name}"
		product = "kms"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating KMS client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{event_name}", d.Id())
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildUpdateCsmsEventBodyParams(d),
	}

	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating CSMS event: %s", err)
	}
	return resourceCsmsEventRead(ctx, d, meta)
}

func buildUpdateCsmsEventBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"event_types":  d.Get("event_types"),
		"state":        d.Get("status"),
		"notification": buildNotificationBodyParam(d),
	}
}

func resourceCsmsEventDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/csms/events/{event_name}"
		product = "kms"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating KMS client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{event_name}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error deleting CSMS event: %s", err)
	}

	return nil
}
