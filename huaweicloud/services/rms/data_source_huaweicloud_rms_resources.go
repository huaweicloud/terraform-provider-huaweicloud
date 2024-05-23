package rms

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CONFIG GET /v1/resource-manager/domains/{domain_id}/all-resources
// @API CONFIG GET /v1/resource-manager/domains/{domain_id}/tracked-resources
func DataSourceResources() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceResourcesRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tracked": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"resource_deleted": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"tags": common.TagsSchema(),
			"resources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
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
						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enterprise_project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enterprise_project_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"checksum": {
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
						"provisioning_state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": common.TagsComputedSchema(),
						"properties": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataSourceResourcesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("rms", region)
	if err != nil {
		return diag.Errorf("error creating RMS v1 client: %s", err)
	}

	var resources []interface{}
	if d.Get("tracked").(bool) || d.Get("resource_deleted").(bool) {
		getTrackedResourcesHttpUrl := "v1/resource-manager/domains/{domain_id}/tracked-resources"
		getTrackedResourcesPath := client.Endpoint + getTrackedResourcesHttpUrl
		getTrackedResourcesPath = strings.ReplaceAll(getTrackedResourcesPath, "{domain_id}", cfg.DomainID)

		resources, err = getResources(client, d, getTrackedResourcesPath)
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		getResourcesHttpUrl := "v1/resource-manager/domains/{domain_id}/all-resources"
		getResourcesPath := client.Endpoint + getResourcesHttpUrl
		getResourcesPath = strings.ReplaceAll(getResourcesPath, "{domain_id}", cfg.DomainID)

		resources, err = getResources(client, d, getResourcesPath)
		if err != nil {
			return diag.FromErr(err)
		}
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

func buildResourcesQueryParams(d *schema.ResourceData, marker string) string {
	res := ""

	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	if v, ok := d.GetOk("resource_id"); ok {
		res = fmt.Sprintf("%s&id=%v", res, v)
	}
	if v, ok := d.GetOk("type"); ok {
		res = fmt.Sprintf("%s&type=%v", res, v)
	}
	if v, ok := d.GetOk("region_id"); ok {
		res = fmt.Sprintf("%s&region_id=%v", res, v)
	}
	if v, ok := d.GetOk("enterprise_project_id"); ok {
		res = fmt.Sprintf("%s&ep_id=%v", res, v)
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := v.(map[string]interface{})
		for k, val := range tagsMap {
			tagsString := fmt.Sprintf(`%s%%3D%s`, k, val)
			res = fmt.Sprintf("%s&tags=%v", res, tagsString)
		}
	}
	if v, ok := d.GetOk("resource_deleted"); ok {
		res = fmt.Sprintf("%s&resource_deleted=%v", res, v)
	}

	if marker != "" {
		res = fmt.Sprintf("%s&marker=%v", res, marker)
	}

	if res != "" {
		res = "?" + res[1:]
	}

	return res
}

func getResources(client *golangsdk.ServiceClient, d *schema.ResourceData, getResourcesPath string) ([]interface{}, error) {
	var resources []interface{}
	var marker string
	getResourcesOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	for {
		requestPath := getResourcesPath + buildResourcesQueryParams(d, marker)
		resp, err := client.Request("GET", requestPath, &getResourcesOpt)
		if err != nil {
			return nil, fmt.Errorf("error retrieving RMS tracked resources: %s", err)
		}

		getTrackedResourcesRespBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		resourcesTemp, err := flattenResources(utils.PathSearch("resources", getTrackedResourcesRespBody, nil))
		if err != nil {
			return nil, err
		}
		resources = append(resources, resourcesTemp...)
		marker = utils.PathSearch("page_info.next_marker", getTrackedResourcesRespBody, "").(string)
		if marker == "" {
			break
		}
	}

	return resources, nil
}

func flattenResources(resourcesResourcesRaw interface{}) ([]interface{}, error) {
	if resourcesResourcesRaw == nil {
		return nil, nil
	}

	resourcesResources := resourcesResourcesRaw.([]interface{})
	res := make([]interface{}, len(resourcesResources))
	for i, v := range resourcesResources {
		resource := v.(map[string]interface{})

		properties, err := flattenResourceProperties(resource["properties"])
		if err != nil {
			return nil, err
		}
		res[i] = map[string]interface{}{
			"id":                      resource["id"],
			"name":                    resource["name"],
			"service":                 resource["provider"],
			"type":                    resource["type"],
			"region_id":               resource["region_id"],
			"project_id":              resource["project_id"],
			"project_name":            resource["project_name"],
			"enterprise_project_id":   resource["ep_id"],
			"enterprise_project_name": resource["ep_name"],
			"checksum":                resource["checksum"],
			"created_at":              resource["created"],
			"updated_at":              resource["updated"],
			"provisioning_state":      resource["provisioning_state"],
			"state":                   resource["state"],
			"tags":                    resource["tags"],
			"properties":              properties,
		}
	}

	return res, nil
}

func flattenResourceProperties(properties interface{}) (map[string]interface{}, error) {
	if properties == nil {
		return nil, nil
	}

	result := make(map[string]interface{})
	for k, v := range properties.(map[string]interface{}) {
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("generate json string of %s failed: %s", k, err)
		}
		result[k] = string(jsonBytes)
	}
	return result, nil
}
