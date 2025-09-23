package dsc

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DSC PUT /v1/{project_id}/sdg/smn/topic
// @API DSC GET /v1/{project_id}/sdg/smn/topics
func ResourceAlarmNotification() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAlarmNotificationCreate,
		ReadContext:   resourceAlarmNotificationRead,
		UpdateContext: resourceAlarmNotificationUpdate,
		DeleteContext: resourceAlarmNotificationDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
				Description: `The region in which to create the resource. If omitted, the provider-level region will be used.`,
			},
			"alarm_topic_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The alarm topic ID.`,
			},
			"topic_urn": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The unique resource identifier of an SMN topic.`,
			},
			"status": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The alarm notification status.`,
			},
		},
	}
}

func buildConfigAlarmNotificationBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"id":        d.Get("alarm_topic_id"),
		"topic_urn": d.Get("topic_urn"),
		"status":    d.Get("status"),
	}
}

func updateAlarmNotification(client *golangsdk.ServiceClient, requestBody map[string]interface{}) error {
	requestPath := client.Endpoint + "v1/{project_id}/sdg/smn/topic"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         requestBody,
	}

	resp, err := client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return err
	}
	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return err
	}

	status := utils.PathSearch("status", respBody, "").(string)
	if status != "success" {
		return fmt.Errorf("got an unexcept status (%s) in API response", status)
	}
	return nil
}

func resourceAlarmNotificationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "dsc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	requestBody := buildConfigAlarmNotificationBodyParams(d)
	if err := updateAlarmNotification(client, requestBody); err != nil {
		return diag.Errorf("error configuring DSC alarm notification in creation operation: %s", err)
	}

	d.SetId(d.Get("alarm_topic_id").(string))

	return resourceAlarmNotificationRead(ctx, d, meta)
}

func resourceAlarmNotificationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		mErr    *multierror.Error
		product = "dsc"
		httpUrl = "v1/{project_id}/sdg/smn/topics"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving DSC alarm notification: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	alarmTopicID := utils.PathSearch("id", respBody, "").(string)
	if alarmTopicID == "" {
		// Normally, this logic cannot be tested.
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("alarm_topic_id", alarmTopicID),
		d.Set("topic_urn", utils.PathSearch("default_topic_urn", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceAlarmNotificationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "dsc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	requestBody := buildConfigAlarmNotificationBodyParams(d)
	if err := updateAlarmNotification(client, requestBody); err != nil {
		return diag.Errorf("error configuring DSC alarm notification in update operation: %s", err)
	}

	return resourceAlarmNotificationRead(ctx, d, meta)
}

func buildDeleteAlarmNotificationBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"id":        d.Get("alarm_topic_id"),
		"topic_urn": d.Get("topic_urn"),
		"status":    0,
	}
}

func resourceAlarmNotificationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "dsc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	requestBody := buildDeleteAlarmNotificationBodyParams(d)
	if err := updateAlarmNotification(client, requestBody); err != nil {
		return diag.Errorf("error disabling DSC alarm notification in deletion operation: %s", err)
	}

	return nil
}
