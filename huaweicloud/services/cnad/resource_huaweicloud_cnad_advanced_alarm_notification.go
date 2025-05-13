package cnad

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API AAD POST /v1/cnad/alarm-config
// @API AAD GET /v1/cnad/alarm-config
// @API AAD DELETE /v1/cnad/alarm-config
func ResourceAlarmNotification() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAlarmNotificationCreate,
		UpdateContext: resourceAlarmNotificationUpdate,
		ReadContext:   resourceAlarmNotificationRead,
		DeleteContext: resourceAlarmNotificationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"topic_urn": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the topic urn of SMN.`,
			},
			"is_close_attack_source_flag": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether to enable the alarm content to shield the attack source information.`,
			},
		},
	}
}

func configAlarmNotification(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	requestPath := client.Endpoint + "v1/cnad/alarm-config"
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		JSONBody: map[string]interface{}{
			"topic_urn": d.Get("topic_urn"),
		},
	}

	_, err := client.Request("POST", requestPath, &requestOpt)
	return err
}

func resourceAlarmNotificationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		product  = "aad"
		topicUrn = d.Get("topic_urn").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CNAD client: %s", err)
	}

	if err := configAlarmNotification(client, d); err != nil {
		return diag.Errorf("error creating CNAD advanced alarm notification: %s", err)
	}
	d.SetId(topicUrn)

	return resourceAlarmNotificationRead(ctx, d, meta)
}

func GetAlarmNotification(client *golangsdk.ServiceClient) (interface{}, error) {
	requestPath := client.Endpoint + "v1/cnad/alarm-config"
	requestOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	// The empty topic urn indicating that the resource does not exist.
	if utils.PathSearch("topic_urn", respBody, "").(string) == "" {
		return nil, golangsdk.ErrDefault404{}
	}

	return respBody, err
}

func resourceAlarmNotificationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "aad"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CNAD client: %s", err)
	}

	respBody, err := GetAlarmNotification(client)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CNAD advanced alarm notification")
	}

	mErr := multierror.Append(
		d.Set("topic_urn", utils.PathSearch("topic_urn", respBody, nil)),
		d.Set("is_close_attack_source_flag", utils.PathSearch("is_close_attack_source_flag", respBody, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceAlarmNotificationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "aad"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CNAD client: %s", err)
	}

	if err := configAlarmNotification(client, d); err != nil {
		return diag.Errorf("error updating CNAD advanced alarm notification: %s", err)
	}

	return resourceAlarmNotificationRead(ctx, d, meta)
}

func resourceAlarmNotificationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "aad"
		httpUrl = "v1/cnad/alarm-config"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CNAD client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error deleting CNAD advanced alarm notification: %s", err)
	}
	return nil
}
