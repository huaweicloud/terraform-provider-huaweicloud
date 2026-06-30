package eip

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

// @API EIP GET /v3/{domain_id}/global-eips/support-instances/{access_site}
func DataSourceGlobalEipTenantSupportRegions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGlobalEipTenantSupportRegionsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"access_site": {
				Type:     schema.TypeString,
				Required: true,
			},
			"fields": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"support_regions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_border_group": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"access_site": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildGlobalEipTenantSupportRegionsQueryParams(d *schema.ResourceData) string {
	req := ""

	if fields := d.Get("fields").([]interface{}); len(fields) > 0 {
		for _, v := range fields {
			req += fmt.Sprintf("&fields=%v", v)
		}
	}

	if req != "" {
		req = "?" + req[1:]
	}

	return req
}

func dataSourceGlobalEipTenantSupportRegionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "geip"
		httpUrl = "v3/{domain_id}/global-eips/support-instances/{access_site}"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GEIP client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{domain_id}", cfg.DomainID)
	requestPath = strings.ReplaceAll(requestPath, "{access_site}", d.Get("access_site").(string))
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	requestPathWithQuery := requestPath + buildGlobalEipTenantSupportRegionsQueryParams(d)
	resp, err := client.Request("GET", requestPathWithQuery, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving global EIP tenant support regions: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("support_regions", flattenGlobalEipTenantSupportRegions(
			utils.PathSearch("support_regions", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGlobalEipTenantSupportRegions(regions []interface{}) []interface{} {
	if len(regions) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(regions))
	for _, r := range regions {
		rst = append(rst, map[string]interface{}{
			"id":                  utils.PathSearch("id", r, nil),
			"instance_type":       utils.PathSearch("instance_type", r, nil),
			"public_border_group": utils.PathSearch("public_border_group", r, nil),
			"region_id":           utils.PathSearch("region_id", r, nil),
			"access_site":         utils.PathSearch("access_site", r, nil),
			"status":              utils.PathSearch("status", r, nil),
			"created_at":          utils.PathSearch("created_at", r, nil),
			"updated_at":          utils.PathSearch("updated_at", r, nil),
		})
	}

	return rst
}
