package live

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	livev1 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/live/v1"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/live/v1/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Live DELETE /v1/{project_id}/record/callbacks/{id}
// @API Live GET /v1/{project_id}/record/callbacks/{id}
// @API Live PUT /v1/{project_id}/record/callbacks/{id}
// @API Live GET /v1/{project_id}/record/callbacks
// @API Live POST /v1/{project_id}/record/callbacks
func ResourceRecordCallback() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRecordCallbackCreate,
		ReadContext:   resourceRecordCallbackRead,
		UpdateContext: resourceRecordCallbackUpdate,
		DeleteContext: resourceRecordCallbackDelete,
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
			"domain_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"types": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 4,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"sign_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"key": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				Computed:  true,
			},
		},
	}
}

func buildEventSubscriptionParams(d *schema.ResourceData) (*[]model.RecordCallbackConfigRequestNotifyEventSubscription, error) {
	v := d.Get("types").([]interface{})
	events := make([]model.RecordCallbackConfigRequestNotifyEventSubscription, len(v))
	for i, v := range v {
		var event model.RecordCallbackConfigRequestNotifyEventSubscription
		err := event.UnmarshalJSON([]byte(v.(string)))
		if err != nil {
			return nil, fmt.Errorf("error getting the argument types: %s", err)
		}
		events[i] = event
	}
	return &events, nil
}

func buildSignTypeParams(d *schema.ResourceData) (*model.RecordCallbackConfigRequestSignType, error) {
	if v, ok := d.GetOk("sign_type"); ok {
		var signType model.RecordCallbackConfigRequestSignType
		err := signType.UnmarshalJSON([]byte(v.(string)))
		if err != nil {
			return nil, fmt.Errorf("error getting the argument sign_type: %s", err)
		}

		return &signType, nil
	}

	return nil, nil
}

func buildCallbackConfigParams(d *schema.ResourceData) (*model.RecordCallbackConfigRequest, error) {
	events, err := buildEventSubscriptionParams(d)
	if err != nil {
		return nil, err
	}

	signType, err := buildSignTypeParams(d)
	if err != nil {
		return nil, err
	}

	req := model.RecordCallbackConfigRequest{
		PublishDomain:           d.Get("domain_name").(string),
		App:                     "*",
		NotifyEventSubscription: events,
		NotifyCallbackUrl:       utils.String(d.Get("url").(string)),
		SignType:                signType,
		Key:                     utils.StringIgnoreEmpty(d.Get("key").(string)),
	}

	return &req, nil
}

func getRecordCallBackId(client *livev1.LiveClient, publishDomain, appName string) (string, error) {
	resp, err := client.ListRecordCallbackConfigs(&model.ListRecordCallbackConfigsRequest{
		PublishDomain: &publishDomain,
		App:           &appName,
	})
	if err != nil {
		return "", fmt.Errorf("error retrieving Live callback configuration: %s", err)
	}

	if resp == nil || resp.CallbackConfig == nil || len(*resp.CallbackConfig) < 1 {
		return "", fmt.Errorf("error retrieving Live callback configuration: no data")
	}

	callbacks := *resp.CallbackConfig

	return *callbacks[0].Id, nil
}

func resourceRecordCallbackCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	region := c.GetRegion(d)
	client, err := c.HcLiveV1Client(region)
	if err != nil {
		return diag.Errorf("error creating Live v1 client: %s", err)
	}

	createBody, err := buildCallbackConfigParams(d)
	if err != nil {
		return diag.FromErr(err)
	}
	createOpts := &model.CreateRecordCallbackConfigRequest{
		Body: createBody,
	}

	_, err = client.CreateRecordCallbackConfig(createOpts)
	if err != nil {
		return diag.Errorf("error creating Live callback configuration: %s", err)
	}

	// ID is lost in CreateRecordCallbackConfig, so get from List API
	id, err := getRecordCallBackId(client, createOpts.Body.PublishDomain, createOpts.Body.App)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	return resourceRecordCallbackRead(ctx, d, meta)
}

func flattenEventSubscriptionAttribute(events *[]model.ShowRecordCallbackConfigResponseNotifyEventSubscription) []string {
	if events == nil {
		return nil
	}

	rst := make([]string, len(*events))
	for i, v := range *events {
		rst[i] = v.Value()
	}

	return rst
}

func flattenSignTypeAttribute(signType *model.ShowRecordCallbackConfigResponseSignType) string {
	if signType == nil {
		return ""
	}

	return signType.Value()
}

func resourceRecordCallbackRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	region := c.GetRegion(d)
	client, err := c.HcLiveV1Client(region)
	if err != nil {
		return diag.Errorf("error creating Live v1 client: %s", err)
	}

	response, err := client.ShowRecordCallbackConfig(&model.ShowRecordCallbackConfigRequest{Id: d.Id()})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Live callback configuration")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("domain_name", response.PublishDomain),
		d.Set("types", flattenEventSubscriptionAttribute(response.NotifyEventSubscription)),
		d.Set("url", response.NotifyCallbackUrl),
		d.Set("sign_type", flattenSignTypeAttribute(response.SignType)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceRecordCallbackUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	region := c.GetRegion(d)
	client, err := c.HcLiveV1Client(region)
	if err != nil {
		return diag.Errorf("error creating Live v1 client: %s", err)
	}

	updateBody, err := buildCallbackConfigParams(d)
	if err != nil {
		return diag.FromErr(err)
	}
	updateOpts := &model.UpdateRecordCallbackConfigRequest{
		Id:   d.Id(),
		Body: updateBody,
	}

	_, err = client.UpdateRecordCallbackConfig(updateOpts)
	if err != nil {
		return diag.Errorf("error updating Live callback configuration: %s", err)
	}

	return resourceRecordCallbackRead(ctx, d, meta)
}

func resourceRecordCallbackDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	region := c.GetRegion(d)
	client, err := c.HcLiveV1Client(region)
	if err != nil {
		return diag.Errorf("error creating Live v1 client: %s", err)
	}

	deleteOpts := &model.DeleteRecordCallbackConfigRequest{
		Id: d.Id(),
	}
	_, err = client.DeleteRecordCallbackConfig(deleteOpts)
	if err != nil {
		return diag.Errorf("error deleting Live callback configuration: %s", err)
	}

	return nil
}
