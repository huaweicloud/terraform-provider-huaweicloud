package huaweicloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/huaweicloud/golangsdk/openstack/rds/v1/datastores"
	"github.com/huaweicloud/golangsdk/openstack/rds/v1/flavors"
)

func dataSourceRdsFlavorV1() *schema.Resource {
	return &schema.Resource{
		Read:               dataSourcedataSourceRdsFlavorV1Read,
		DeprecationMessage: "use huaweicloud_rds_flavors data source instead",
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"datastore_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"datastore_version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ram": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"speccode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourcedataSourceRdsFlavorV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	rdsClient, err := config.RdsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud rds client: %s", err)
	}

	datastoresList, err := datastores.List(rdsClient, d.Get("datastore_name").(string)).Extract()
	if err != nil {
		return fmt.Errorf("Unable to retrieve datastores: %s ", err)
	}

	if len(datastoresList) < 1 {
		return fmt.Errorf("Returned no datastore result. ")
	}
	var datastoreId string
	for _, datastore := range datastoresList {
		if datastore.Name == d.Get("datastore_version").(string) {
			datastoreId = datastore.ID
		}
	}
	if datastoreId == "" {
		return fmt.Errorf("Returned no datastore ID. ")
	}
	log.Printf("[DEBUG] Received datastore Id: %s", datastoreId)

	flavorsList, err := flavors.List(rdsClient, datastoreId, d.Get("region").(string)).Extract()
	if err != nil {
		return fmt.Errorf("Unable to retrieve flavors: %s", err)
	}
	if len(flavorsList) < 1 {
		return fmt.Errorf("Returned no flavor result. ")
	}

	var rdsFlavor flavors.Flavor
	if d.Get("speccode").(string) == "" {
		rdsFlavor = flavorsList[0]
	} else {
		for _, flavor := range flavorsList {
			if flavor.SpecCode == d.Get("speccode").(string) {
				rdsFlavor = flavor
			}
		}
	}
	log.Printf("[DEBUG] Retrieved flavorId %s: %+v ", rdsFlavor.ID, rdsFlavor)
	if rdsFlavor.ID == "" {
		return fmt.Errorf("Returned no flavor Id. ")
	}

	d.SetId(rdsFlavor.ID)

	d.Set("id", rdsFlavor.ID)
	d.Set("name", rdsFlavor.Name)
	d.Set("ram", rdsFlavor.Ram)
	d.Set("speccode", rdsFlavor.SpecCode)
	d.Set("region", GetRegion(d, config))

	return nil
}
