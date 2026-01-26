package ram

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

// @API RAM GET /v1/resource-shares/quotas
func DataSourceQuotas() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceQuotasRead,
		Schema: map[string]*schema.Schema{
			"quotas": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     quotasSchema(),
			},
		},
	}
}

func quotasSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"resources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"quota": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"min": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"max": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"used": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}

	return &sc
}

func dataSourceQuotasRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/resource-shares/quotas"
		product = "ram"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RAM client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpts)
	if err != nil {
		return diag.Errorf("error retrieving RAM quotas, %s", err)
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

	mErr := multierror.Append(nil,
		d.Set("quotas", flattenQuotas(utils.PathSearch("quotas", respBody, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenQuotas(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	quotaMap, ok := resp.(map[string]interface{})
	if !ok {
		return nil
	}

	resources, ok := quotaMap["resources"].([]interface{})
	if !ok {
		return nil
	}

	resourcesRst := make([]interface{}, 0, len(resources))
	for _, v := range resources {
		resourcesRst = append(resourcesRst, map[string]interface{}{
			"type":  utils.PathSearch("type", v, nil),
			"quota": utils.PathSearch("quota", v, nil),
			"min":   utils.PathSearch("min", v, nil),
			"max":   utils.PathSearch("max", v, nil),
			"used":  utils.PathSearch("used", v, nil),
		})
	}

	return []interface{}{
		map[string]interface{}{
			"resources": resourcesRst,
		},
	}
}
