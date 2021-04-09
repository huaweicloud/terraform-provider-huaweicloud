package huaweicloud

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/swr/v2/namespaces"
)

func resourceSWROrganization() *schema.Resource {
	return &schema.Resource{
		Create: resourceSWROrganizationCreate,
		Read:   resourceSWROrganizationRead,
		Delete: resourceSWROrganizationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		//request and response parameters
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"creator": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"permission": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"login_server": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSWROrganizationCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	swrClient, err := config.SwrV2Client(GetRegion(d, config))

	if err != nil {
		return fmt.Errorf("Unable to create HuaweiCloud SWR client : %s", err)
	}

	name := d.Get("name").(string)
	createOpts := namespaces.CreateOpts{
		Namespace: name,
	}

	err = namespaces.Create(swrClient, createOpts).ExtractErr()

	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud SWR Organization: %s", err)
	}

	d.SetId(name)

	return resourceSWROrganizationRead(d, meta)
}

func resourceSWROrganizationRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	swrClient, err := config.SwrV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud SWR client: %s", err)
	}

	n, err := namespaces.Get(swrClient, d.Id()).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving HuaweiCloud SWR: %s", err)
	}

	permission := "Unknown"
	switch n.Auth {
	case 7:
		permission = "Manage"
	case 3:
		permission = "Write"
	case 1:
		permission = "Read"
	}

	d.Set("region", GetRegion(d, config))
	d.Set("name", n.Name)
	d.Set("creator", n.CreatorName)
	d.Set("permission", permission)

	login := fmt.Sprintf("swr.%s.%s", GetRegion(d, config), config.Cloud)
	d.Set("login_server", login)

	return nil
}

func resourceSWROrganizationDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	swrClient, err := config.SwrV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud SWR Client: %s", err)
	}

	err = namespaces.Delete(swrClient, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error deleting HuaweiCloud SWR Organization: %s", err)
	}

	d.SetId("")
	return nil
}
