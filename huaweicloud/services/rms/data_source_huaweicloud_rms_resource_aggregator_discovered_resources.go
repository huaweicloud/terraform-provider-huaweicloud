package rms

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CONFIG POST /v1/resource-manager/domains/{domain_id}/aggregators/aggregate-data/aggregate-discovered-resources
func DataSourceAggregatorDiscoveredResources() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAggregatorDiscoveredResourcesRead,

		Schema: map[string]*schema.Schema{
			"aggregator_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the resource aggregator ID.`,
			},
			"service_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"filter": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"resource_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"resource_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"resources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_account_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAggregatorDiscoveredResourcesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("rms", region)
	if err != nil {
		return diag.Errorf("error creating RMS v1 client: %s", err)
	}

	getAggregatorDiscoveredResourcesHttpUrl := "v1/resource-manager/domains/{domain_id}/aggregators/aggregate-data/aggregate-discovered-resources"
	getAggregatorDiscoveredResourcesPath := client.Endpoint + getAggregatorDiscoveredResourcesHttpUrl
	getAggregatorDiscoveredResourcesPath = strings.ReplaceAll(getAggregatorDiscoveredResourcesPath, "{domain_id}", cfg.DomainID)

	resources, err := getAggregatorDiscoveredResources(client, d, getAggregatorDiscoveredResourcesPath)
	if err != nil {
		return diag.FromErr(err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(
		nil,
		d.Set("resources", resources),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildAggregatorDiscoveredResourcesQueryParams(marker string) string {
	res := "?limit=200"

	if marker != "" {
		res += fmt.Sprintf("&marker=%v", marker)
	}

	return res
}

func getAggregatorDiscoveredResources(client *golangsdk.ServiceClient, d *schema.ResourceData,
	getAggregatorDiscoveredResourcesPath string) ([]interface{}, error) {
	var resources []interface{}
	var marker string
	getAggregatorDiscoveredResourcesOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildAggregatorDiscoveredResourcesBodyParams(d)),
	}
	for {
		requestPath := getAggregatorDiscoveredResourcesPath + buildAggregatorDiscoveredResourcesQueryParams(marker)
		resp, err := client.Request("POST", requestPath, &getAggregatorDiscoveredResourcesOpt)
		if err != nil {
			return nil, fmt.Errorf("error retrieving aggregator discovered resources: %s", err)
		}

		getTrackedAggregatorDiscoveredResourcesRespBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		resourcesTemp := flattenAggregatorDiscoveredResources(
			utils.PathSearch("resource_identifiers", getTrackedAggregatorDiscoveredResourcesRespBody, nil))
		if err != nil {
			return nil, err
		}
		resources = append(resources, resourcesTemp...)
		marker = utils.PathSearch("page_info.next_marker", getTrackedAggregatorDiscoveredResourcesRespBody, "").(string)
		if marker == "" {
			break
		}
	}

	return resources, nil
}

func buildAggregatorDiscoveredResourcesBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"aggregator_id": d.Get("aggregator_id"),
		"provider":      utils.ValueIgnoreEmpty(d.Get("service_type")),
		"resource_type": utils.ValueIgnoreEmpty(d.Get("resource_type")),
		"filter":        buildAggregatorDiscoveredResourcesFilterBodyParams(d),
	}
	return bodyParams
}

func buildAggregatorDiscoveredResourcesFilterBodyParams(d *schema.ResourceData) map[string]interface{} {
	filterRaw := d.Get("filter").([]interface{})
	if len(filterRaw) == 0 {
		return nil
	}

	if filter, ok := filterRaw[0].(map[string]interface{}); ok {
		bodyParams := map[string]interface{}{
			"account_id":    utils.ValueIgnoreEmpty(filter["account_id"]),
			"region_id":     utils.ValueIgnoreEmpty(filter["region_id"]),
			"resource_id":   utils.ValueIgnoreEmpty(filter["resource_id"]),
			"resource_name": utils.ValueIgnoreEmpty(filter["resource_name"]),
		}

		return bodyParams
	}
	return nil
}

func flattenAggregatorDiscoveredResources(resourcesAggregatorDiscoveredResourcesRaw interface{}) []interface{} {
	if resourcesAggregatorDiscoveredResourcesRaw == nil {
		return nil
	}

	resourcesAggregatorDiscoveredResources := resourcesAggregatorDiscoveredResourcesRaw.([]interface{})
	res := make([]interface{}, len(resourcesAggregatorDiscoveredResources))
	for i, v := range resourcesAggregatorDiscoveredResources {
		resource := v.(map[string]interface{})
		res[i] = map[string]interface{}{
			"resource_id":       resource["resource_id"],
			"resource_name":     resource["resource_name"],
			"service":           resource["provider"],
			"type":              resource["type"],
			"source_account_id": resource["source_account_id"],
			"region_id":         resource["region_id"],
		}
	}

	return res
}
