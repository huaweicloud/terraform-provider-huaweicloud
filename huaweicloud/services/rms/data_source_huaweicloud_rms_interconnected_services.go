package rms

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CONFIG GET /v1/resource-manager/domains/{domain_id}/all-providers
func DataSourceRmsInterconnectedServices() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRmsInterconnectedServicesRead,
		Schema: map[string]*schema.Schema{
			"resource_providers": {
				Type:     schema.TypeList,
				Elem:     interconnectedServicesResourceProvider(),
				Computed: true,
			},
		},
	}
}

func interconnectedServicesResourceProvider() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"provider": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"category_display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_types": {
				Type:     schema.TypeList,
				Elem:     interconnectedServicesResourceProviderResourceType(),
				Computed: true,
			},
		},
	}
	return &sc
}

func interconnectedServicesResourceProviderResourceType() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"global": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"regions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"console_endpoint_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"console_detail_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"console_list_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"track": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func dataSourceRmsInterconnectedServicesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v1/resource-manager/domains/{domain_id}/all-providers"
		product = "rms"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Config client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{domain_id}", cfg.DomainID)

	listResp, err := pagination.ListAllItems(
		client,
		"offset",
		listPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return diag.Errorf("error retrieving Config interconnected services: %s", err)
	}

	listRespJson, err := json.Marshal(listResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var listRespBody interface{}
	err = json.Unmarshal(listRespJson, &listRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("resource_providers", flattenInterconnectedServicesResourceProviders(listRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenInterconnectedServicesResourceProviders(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("resource_providers", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"provider":              utils.PathSearch("provider", v, nil),
			"display_name":          utils.PathSearch("display_name", v, nil),
			"category_display_name": utils.PathSearch("category_display_name", v, nil),
			"resource_types":        flattenInterconnectedServicesResourceProvidersResourceTypes(v),
		})
	}
	return rst
}

func flattenInterconnectedServicesResourceProvidersResourceTypes(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("resource_types", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"name":                utils.PathSearch("name", v, nil),
			"display_name":        utils.PathSearch("display_name", v, nil),
			"global":              utils.PathSearch("global", v, nil),
			"regions":             utils.PathSearch("regions", v, nil),
			"console_endpoint_id": utils.PathSearch("console_endpoint_id", v, nil),
			"console_list_url":    utils.PathSearch("console_list_url", v, nil),
			"console_detail_url":  utils.PathSearch("console_detail_url", v, nil),
			"track":               utils.PathSearch("track", v, nil),
		})
	}
	return rst
}
