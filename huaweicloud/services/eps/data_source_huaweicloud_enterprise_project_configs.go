package eps

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

// @API EPS GET /v1/enterprise-projects/configs
func DataSourceConfigs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceConfigsRead,

		Schema: map[string]*schema.Schema{
			// Attributes.
			"support_item": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeBool},
			},
		},
	}
}

func dataSourceConfigsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("eps", region)
	if err != nil {
		return diag.Errorf("error creating EPS client: %s", err)
	}

	return ShowEpConfigs(client, d)
}

func ShowEpConfigs(client *golangsdk.ServiceClient, d *schema.ResourceData) diag.Diagnostics {
	httpUrl := "v1/enterprise-projects/configs"
	showPath := client.Endpoint + httpUrl

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", showPath, &opt)
	if err != nil {
		return diag.FromErr(err)
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr := multierror.Append(nil,
		d.Set("support_item", map[string]interface{}{
			"delete_ep_support": utils.PathSearch("support_item.delete_ep_support", respBody, nil).(bool),
		}),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
