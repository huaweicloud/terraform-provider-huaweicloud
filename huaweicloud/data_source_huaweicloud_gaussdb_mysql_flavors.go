package huaweicloud

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk"
)

func dataSourceGaussdbMysqlFlavors() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGaussdbMysqlFlavorsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"engine": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "gaussdb-mysql",
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "8.0",
			},
			"availability_zone_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "single",
			},
			"flavors": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vcpus": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"memory": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceGaussdbMysqlFlavorsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	client, err := config.gaussdbV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud GaussDB client: %s", err)
	}

	link := fmt.Sprintf("flavors/%s?version_name=%s&availability_zone_mode=%s", d.Get("engine").(string), d.Get("version").(string), d.Get("availability_zone_mode").(string))
	url := client.ServiceURL(link)

	r, err := sendGaussdbMysqlFlavorsListRequest(client, url)
	if err != nil {
		return err
	}

	flavors := make([]interface{}, 0, len(r.([]interface{})))
	for _, item := range r.([]interface{}) {
		val := item.(map[string]interface{})

		flavors = append(flavors, map[string]interface{}{
			"vcpus":   val["vcpus"],
			"memory":  val["ram"],
			"name":    val["spec_code"],
			"mode":    val["instance_mode"],
			"version": val["version_name"],
		})
	}

	d.SetId("flavors")
	return d.Set("flavors", flavors)
}

func sendGaussdbMysqlFlavorsListRequest(client *golangsdk.ServiceClient, url string) (interface{}, error) {
	r := golangsdk.Result{}
	_, r.Err = client.Get(url, &r.Body, &golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"X-Language":   "en-us",
		}})
	if r.Err != nil {
		return nil, fmt.Errorf("Error fetching flavors for gaussdb mysql, error: %s", r.Err)
	}

	v, err := navigateValue(r.Body, []string{"flavors"}, nil)
	if err != nil {
		return nil, err
	}
	return v, nil
}
