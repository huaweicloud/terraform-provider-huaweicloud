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
			"support_item_attribute": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"delete_ep_support_attribute": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether enterprise projects can be deleted.`,
						},
					},
				},
				Description: `The list of configurations.`,
			},
			// The structure of attribute field `support_item` is flawed; it should be deprecated and the field
			// `support_item_attribute` should be used instead.
			"support_item": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeBool},
				Description: utils.SchemaDesc(
					"Refer to the attribute field delete_ep_support_attribute.",
					utils.SchemaDescInput{
						Deprecated: true,
					},
				),
			},
		},
	}
}

func dataSourceConfigsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/enterprise-projects/configs"
	)
	client, err := cfg.NewServiceClient("eps", region)
	if err != nil {
		return diag.Errorf("error creating EPS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving EPS configs: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr := multierror.Append(
		d.Set("support_item_attribute", flattenSupportItemAttribute(utils.PathSearch("support_item", respBody, nil))),
		d.Set("support_item", flattenSupportItem(utils.PathSearch("support_item.delete_ep_support", respBody, false).(bool))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSupportItem(respBool bool) map[string]interface{} {
	return map[string]interface{}{
		"delete_ep_support": respBool,
	}
}

func flattenSupportItemAttribute(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"delete_ep_support_attribute": utils.PathSearch("delete_ep_support", respBody, nil),
		},
	}
}
