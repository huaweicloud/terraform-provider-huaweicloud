package iec

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk/openstack/iec/v1/images"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API IEC GET /v1/images
func DataSourceImages() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceImagesRead,

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

func dataSourceImagesRead(d *schema.ResourceData, meta interface{}) error {
	cfg := meta.(*config.Config)

	iecClient, err := cfg.IECV1Client(cfg.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating IEC client: %s", err)
	}

	listOpts := images.ListOpts{
		Name:    d.Get("name").(string),
		OsType:  d.Get("os_type").(string),
		Status:  "active",
		SortKey: "name",
	}
	pages, err := images.List(iecClient, listOpts).AllPages()
	if err != nil {
		return fmt.Errorf("unable to retrieve IEC images: %s", err)
	}

	allImages, err := images.ExtractImages(pages)
	if err != nil {
		return fmt.Errorf("unable to extract IEC images: %s", err)
	}
	total := len(allImages.Images)
	if total < 1 {
		return fmt.Errorf("your query returned no results of iec_images, " +
			"please change your search criteria and try again")
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
	mErr := multierror.Append(d.Set("images", edgeImages))
	if err := mErr.ErrorOrNil(); err != nil {
		return fmt.Errorf("error saving IEC iamges: %s", err)
	}

	d.SetId(allImages.Images[0].ID)
	return nil
}
