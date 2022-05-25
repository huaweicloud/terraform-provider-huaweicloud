package live

import (
	"context"
	"fmt"
	"log"
	"regexp"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	v1 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/live/v1"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/live/v1/model"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const (
	eventRecordNewFileStart = "RECORD_NEW_FILE_START"
	eventRecordFileComplete = "RECORD_FILE_COMPLETE"
	eventRecordOver         = "RECORD_OVER"
	eventRecordFailed       = "RECORD_FAILED"

	allAppName = "*"
)

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
				ValidateFunc: validation.All(
					validation.StringMatch(regexp.MustCompile(`^http://|^https://`),
						"The URL must start with http:// or https://."),
					validation.StringDoesNotMatch(regexp.MustCompile(`\?{1}`), "The URL cannot contain parameters."),
				),
			},

			"types": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 4,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{eventRecordNewFileStart, eventRecordFileComplete,
						eventRecordOver, eventRecordFailed}, false),
				},
			},
		},
	}
}

func resourceRecordCallbackCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	region := c.GetRegion(d)
	client, err := c.HcLiveV1Client(region)
	if err != nil {
		return diag.Errorf("error creating Live v1 client: %s", err)
	}

	createBody, err := buildCallbackConfigParams(d, region)
	if err != nil {
		return diag.FromErr(err)
	}
	createOpts := &model.CreateRecordCallbackConfigRequest{
		Body: createBody,
	}

	log.Printf("[DEBUG] Create Live callback configuration params : %#v", createOpts)

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

func resourceRecordCallbackRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		setTypeToState(d, response.NotifyEventSubscription),
		d.Set("url", response.NotifyCallbackUrl),
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

	updateBody, err := buildCallbackConfigParams(d, region)
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

func resourceRecordCallbackDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

func buildCallbackConfigParams(d *schema.ResourceData, region string) (*model.RecordCallbackConfigRequest, error) {
	v := d.Get("types").([]interface{})
	events := make([]model.RecordCallbackConfigRequestNotifyEventSubscription, len(v))
	for i, v := range v {
		var event model.RecordCallbackConfigRequestNotifyEventSubscription
		err := event.UnmarshalJSON([]byte(v.(string)))
		if err != nil {
			return nil, fmt.Errorf("error getting the argument %q: %s", "types", err)
		}
		events[i] = event
	}

	req := model.RecordCallbackConfigRequest{
		PublishDomain:           d.Get("domain_name").(string),
		App:                     allAppName,
		NotifyEventSubscription: &events,
		NotifyCallbackUrl:       utils.String(d.Get("url").(string)),
	}

	return &req, nil

}

func getRecordCallBackId(client *v1.LiveClient, publishDomain, appName string) (string, error) {
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

func setTypeToState(d *schema.ResourceData, event *[]model.ShowRecordCallbackConfigResponseNotifyEventSubscription) error {
	if event != nil {
		types := make([]string, len(*event))
		for i, v := range *event {
			event := utils.MarshalValue(v)
			types[i] = event
		}
		return d.Set("types", types)
	}
	return nil
}
