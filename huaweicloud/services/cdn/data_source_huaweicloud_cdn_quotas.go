package cdn

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

// @API CDN GET /v1.0/cdn/quota
func DataSourceQuotas() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceQuotasRead,

		Schema: map[string]*schema.Schema{
			// Attributes.
			"quotas": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"limit": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The limit of the resource quota.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the resource quota.`,
						},
						"used": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The used capacity of the resource quota.`,
						},
						"user_domain_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The domain ID of the user to which the resource quota belong.`,
						},
					},
				},
				Description: `The list of the resource quotas that matched filter parameters.`,
			},
		},
	}
}

func getQuotas(client *golangsdk.ServiceClient) ([]interface{}, error) {
	httpUrl := "v1.0/cdn/quota"
	getPath := client.Endpoint + httpUrl

	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch("quotas", respBody, make([]interface{}, 0)).([]interface{}), nil
}

func flattenQuotas(quotas []interface{}) []map[string]interface{} {
	if len(quotas) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(quotas))
	for _, quota := range quotas {
		result = append(result, map[string]interface{}{
			"limit":          utils.PathSearch("quota_limit", quota, nil),
			"type":           utils.PathSearch("type", quota, nil),
			"used":           utils.PathSearch("used", quota, nil),
			"user_domain_id": utils.PathSearch("user_domain_id", quota, nil),
		})
	}

	return result
}

func dataSourceQuotasRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("cdn", "")
	if err != nil {
		return diag.Errorf("error creating CDN client: %s", err)
	}

	resp, err := getQuotas(client)
	if err != nil {
		return diag.Errorf("error querying resource quotas: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(
		d.Set("quotas", flattenQuotas(resp)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
