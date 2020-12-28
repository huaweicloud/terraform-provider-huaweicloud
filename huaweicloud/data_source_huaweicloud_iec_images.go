package huaweicloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/huaweicloud/golangsdk/openstack/iec/v1/images"
)

func dataSourceIecImages() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIecImagesV1Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"os_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Linux", "Windows", "Other",
				}, false),
			},
			"images": {
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
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"os_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIecImagesV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	iecClient, err := config.IECV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud IEC client: %s", err)
	}

	listOpts := images.ListOpts{
		Name:    d.Get("name").(string),
		OsType:  d.Get("os_type").(string),
		Status:  "active",
		SortKey: "name",
	}
	pages, err := images.List(iecClient, listOpts).AllPages()
	if err != nil {
		return fmt.Errorf("Unable to retrieve iec images: %s", err)
	}

	allImages, err := images.ExtractImages(pages)
	if err != nil {
		return fmt.Errorf("Unable to extract iec images: %s", err)
	}
	total := len(allImages.Images)
	if total < 1 {
		return fmt.Errorf("Your query returned no results of huaweicloud_iec_images. " +
			"Please change your search criteria and try again.")
	}

	log.Printf("[INFO] Retrieved [%d] IEC images using given filter", total)
	edgeImages := make([]map[string]interface{}, 0, total)
	for _, item := range allImages.Images {
		val := map[string]interface{}{
			"id":      item.ID,
			"name":    item.Name,
			"status":  item.Status,
			"os_type": item.OsType,
		}
		edgeImages = append(edgeImages, val)
	}
	if err := d.Set("images", edgeImages); err != nil {
		return fmt.Errorf("Error saving IEC iamges: %s", err)
	}

	d.SetId(allImages.Images[0].ID)
	return nil
}
