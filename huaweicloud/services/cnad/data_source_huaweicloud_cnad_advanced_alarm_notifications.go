package cnad

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API AAD GET /v1/cnad/alarm-config
func DataSourceAlarmNotifications() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAlarmNotificationsRead,
		Schema: map[string]*schema.Schema{
			"topic_urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The topic urn of SMN.`,
			},
			"is_close_attack_source_flag": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether to enable the alarm content to shield the attack source information.`,
			},
		},
	}
}

func dataSourceAlarmNotificationsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "aad"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CNAD client: %s", err)
	}

	requestPath := client.Endpoint + "v1/cnad/alarm-config"
	requestOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving CNAD advanced alarm notifications: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(generateUUID)

	mErr := multierror.Append(
		d.Set("topic_urn", utils.PathSearch("topic_urn", respBody, nil)),
		d.Set("is_close_attack_source_flag", utils.PathSearch("is_close_attack_source_flag", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
