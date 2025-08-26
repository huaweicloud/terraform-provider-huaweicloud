package rms

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CONFIG POST /v1/resource-manager/domains/{domain_id}/aggregators/aggregate-resource-config
func DataSourceResourceAggregatorResourceDetail() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceResourceAggregatorResourceDetailRead,

		Schema: map[string]*schema.Schema{
			"aggregator_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"resource_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"service_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"source_account_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"resource_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"aggregator_domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ep_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"properties": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceResourceAggregatorResourceDetailRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v1/resource-manager/domains/{domain_id}/aggregators/aggregate-resource-config"
		product = "rms"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Config client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{domain_id}", cfg.DomainID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildResourceAggregatorResourceDetailQueryParams(d),
	}

	getResp, err := client.Request("POST", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving Config resource aggregator resource detail, %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
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
		d.Set("resource_id", utils.PathSearch("resource_id", getRespBody, nil)),
		d.Set("aggregator_id", utils.PathSearch("aggregator_id", getRespBody, nil)),
		d.Set("service_type", utils.PathSearch("provider", getRespBody, nil)),
		d.Set("type", utils.PathSearch("type", getRespBody, nil)),
		d.Set("region_id", utils.PathSearch("region_id", getRespBody, nil)),
		d.Set("resource_name", utils.PathSearch("name", getRespBody, nil)),
		d.Set("aggregator_domain_id", utils.PathSearch("aggregator_domain_id", getRespBody, nil)),
		d.Set("ep_id", utils.PathSearch("ep_id", getRespBody, nil)),
		d.Set("created", utils.PathSearch("created", getRespBody, nil)),
		d.Set("updated", utils.PathSearch("updated", getRespBody, nil)),
		d.Set("tags", utils.PathSearch("tags", getRespBody, nil)),
		d.Set("properties", flattenResourceAggregatorResourceDetailProperties(getRespBody)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildResourceAggregatorResourceDetailQueryParams(d *schema.ResourceData) map[string]interface{} {
	resourceIdentifierParams := map[string]interface{}{
		"resource_id":       d.Get("resource_id"),
		"provider":          d.Get("service_type"),
		"type":              d.Get("type"),
		"source_account_id": d.Get("source_account_id"),
		"region_id":         d.Get("region_id"),
		"resource_name":     utils.ValueIgnoreEmpty(d.Get("resource_name")),
	}
	bodyParams := map[string]interface{}{
		"aggregator_id":       d.Get("aggregator_id"),
		"resource_identifier": resourceIdentifierParams,
	}
	return bodyParams
}

func flattenResourceAggregatorResourceDetailProperties(resp interface{}) map[string]interface{} {
	curJson := utils.PathSearch("properties", resp, nil)
	if curJson == nil {
		return nil
	}
	rst := make(map[string]interface{})
	for k, v := range curJson.(map[string]interface{}) {
		jsonBytes, _ := json.Marshal(v)
		rst[k] = string(jsonBytes)
	}
	return rst
}
