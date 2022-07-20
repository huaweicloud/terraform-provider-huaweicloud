package huaweicloud

import (
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/taurusdb/v3/instances"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func dataSourceGaussDBMysqlDehResource() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGaussDBMysqlDehResourceRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"resource_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"availability_zone": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"architecture": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vcpus": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ram": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"volume": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGaussDBMysqlDehResourceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	region := GetRegion(d, config)
	client, err := config.GaussdbV3Client(region)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud GaussDB client: %s", err)
	}

	pages, err := instances.ListDeh(client).AllPages()
	if err != nil {
		return err
	}

	allResources, err := instances.ExtractDehResources(pages)
	if err != nil {
		return fmtp.Errorf("Unable to retrieve dedicated resources: %s", err)
	}

	resource_name := d.Get("resource_name").(string)
	refinedResources := []instances.DehResource{}
	for _, refResource := range allResources.Resources {
		if refResource.EngineName != "taurus" {
			continue
		}
		if resource_name != "" && refResource.ResourceName != resource_name {
			continue
		}
		refinedResources = append(refinedResources, refResource)
	}

	if len(refinedResources) < 1 {
		return fmtp.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedResources) > 1 {
		return fmtp.Errorf("Your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	resource := refinedResources[0]

	logp.Printf("[DEBUG] Retrieved Resource %s: %+v", resource.Id, resource)
	d.SetId(resource.Id)

	d.Set("resource_name", resource.ResourceName)
	d.Set("availability_zone", resource.AvailabilityZone)
	d.Set("architecture", resource.Architecture)
	d.Set("vcpus", resource.Capacity.Vcpus)
	d.Set("ram", resource.Capacity.Ram)
	d.Set("volume", resource.Capacity.Volume)
	d.Set("status", resource.Status)
	d.Set("region", region)

	return nil
}
