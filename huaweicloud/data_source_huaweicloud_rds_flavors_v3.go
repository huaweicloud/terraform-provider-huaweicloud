package huaweicloud

import (
	"fmt"
	"strconv"

	"github.com/chnsz/golangsdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func DataSourceRdsFlavorV3() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRdsFlavorV3Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"db_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"MySQL", "PostgreSQL", "SQLServer",
				}, true),
			},
			"db_version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_mode": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ha", "single", "replica",
				}, false),
			},
			"vcpus": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"memory": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"flavors": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vcpus": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"memory": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceRdsFlavorV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)

	client, err := config.RdsV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud rds client: %s", err)
	}

	link := fmt.Sprintf("flavors/%s?version_name=%s",
		d.Get("db_type").(string), d.Get("db_version").(string))
	url := client.ServiceURL(link)

	r, err := sendRdsFlavorV3ListRequest(client, url)
	if err != nil {
		return err
	}

	mode := d.Get("instance_mode").(string)
	cpu := d.Get("vcpus").(int)
	mem := d.Get("memory").(int)
	flavors := make([]interface{}, 0, len(r.([]interface{})))
	for _, item := range r.([]interface{}) {
		val := item.(map[string]interface{})
		vcpu, _ := strconv.Atoi(val["vcpus"].(string))
		if cpu > 0 && vcpu != cpu {
			continue
		}

		if mem > 0 && int(val["ram"].(float64)) != mem {
			continue
		}

		if mode == val["instance_mode"].(string) {
			flavors = append(flavors, map[string]interface{}{
				"vcpus":  val["vcpus"],
				"memory": val["ram"],
				"name":   val["spec_code"],
				"mode":   val["instance_mode"],
			})
		}
	}

	if len(flavors) == 0 {
		return fmtp.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	d.SetId("flavors")
	return d.Set("flavors", flavors)
}

func sendRdsFlavorV3ListRequest(client *golangsdk.ServiceClient, url string) (interface{}, error) {
	r := golangsdk.Result{}
	_, r.Err = client.Get(url, &r.Body, &golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"X-Language":   "en-us",
		}})
	if r.Err != nil {
		return nil, fmtp.Errorf("Error fetching flavors for rds v3, error: %s", r.Err)
	}

	v, err := navigateValue(r.Body, []string{"flavors"}, nil)
	if err != nil {
		return nil, err
	}
	return v, nil
}
