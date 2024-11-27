package geminidb

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/geminidb/v3/instances"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API GaussDBforNoSQL GET /v3/{project_id}/dedicated-resources
func DataSourceGeminiDBDehResource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGeminiDBDehResourceRead,

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
				Type:     schema.TypeString,
				Computed: true,
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

func dataSourceGeminiDBDehResourceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.GeminiDBV3Client(region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	pages, err := instances.ListDeh(client).AllPages()
	if err != nil {
		return diag.Errorf("error getting GaussDB cassandra dedicated list: %s", err)
	}

	allResources, err := instances.ExtractDehResources(pages)
	if err != nil {
		return diag.Errorf("unable to retrieve dedicated resources: %s", err)
	}

	resourceName := d.Get("resource_name").(string)
	refinedResources := []instances.DehResource{}
	for _, refResource := range allResources.Resources {
		if refResource.EngineName != "cassandra" {
			continue
		}
		if resourceName != "" && refResource.ResourceName != resourceName {
			continue
		}
		refinedResources = append(refinedResources, refResource)
	}

	if len(refinedResources) < 1 {
		return diag.Errorf("your query returned no results. " +
			"please change your search criteria and try again.")
	}

	if len(refinedResources) > 1 {
		return diag.Errorf("your query returned more than one result." +
			" please try a more specific search criteria")
	}

	resource := refinedResources[0]

	log.Printf("[DEBUG] Retrieved Resource %s: %+v", resource.Id, resource)
	d.SetId(resource.Id)

	mErr := multierror.Append(
		d.Set("resource_name", resource.ResourceName),
		d.Set("availability_zone", resource.AvailabilityZone),
		d.Set("architecture", resource.Architecture),
		d.Set("vcpus", resource.Capacity.Vcpus),
		d.Set("ram", resource.Capacity.Ram),
		d.Set("volume", resource.Capacity.Volume),
		d.Set("status", resource.Status),
		d.Set("region", region),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
