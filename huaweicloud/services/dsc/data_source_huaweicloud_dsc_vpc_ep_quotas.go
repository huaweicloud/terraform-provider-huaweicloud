package dsc

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DSC GET /v1/{project_id}/vpcep/quotas
func DataSourceDscVpcEpQuotas() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDscVpcEpQuotasRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"quotas": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The VPC endpoint quota information list.",
				Elem:        dscVpcEpQuotaSchema(),
			},
		},
	}
}

func dscVpcEpQuotaSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"quota": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The quota.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The quota type.",
			},
			"used": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The used quota.",
			},
		},
	}
}

func buildDscVpcEpQuotasQueryParams(limit, offset int) string {
	return fmt.Sprintf("?limit=%d&offset=%d", limit, offset)
}

func dataSourceDscVpcEpQuotasRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "dsc"
		httpUrl = "v1/{project_id}/vpcep/quotas"
		offset  = 0
		limit   = 1000
		result  = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"content-type": "application/json;charset=UTF-8",
		},
	}

	for {
		currentPath := requestPath + buildDscVpcEpQuotasQueryParams(limit, offset)

		resp, err := client.Request("GET", currentPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving DSC VPC endpoint quotas: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		quotasResp := utils.PathSearch("quotas.resources", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, quotasResp...)

		if len(quotasResp) < limit {
			break
		}

		offset += len(quotasResp)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("quotas", flattenDscVpcEpQuotas(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDscVpcEpQuotas(quotas []interface{}) []interface{} {
	if len(quotas) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(quotas))
	for _, v := range quotas {
		rst = append(rst, map[string]interface{}{
			"quota": utils.PathSearch("quota", v, nil),
			"type":  utils.PathSearch("type", v, nil),
			"used":  utils.PathSearch("used", v, nil),
		})
	}

	return rst
}
