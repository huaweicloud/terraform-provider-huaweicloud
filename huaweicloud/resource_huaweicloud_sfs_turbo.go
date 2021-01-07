package huaweicloud

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/sfs_turbo/v1/shares"
)

func resourceSFSTurbo() *schema.Resource {
	return &schema.Resource{
		Create: resourceSFSTurboCreate,
		Read:   resourceSFSTurboRead,
		Update: resourceSFSTurboUpdate,
		Delete: resourceSFSTurboDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(4, 64),
			},
			"size": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(500),
			},
			"share_proto": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "NFS",
				ValidateFunc: validation.StringInSlice([]string{"NFS"}, false),
			},
			"share_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "STANDARD",
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"crypt_key_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"enhanced": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"export_location": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"available_capacity": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSFSTurboCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	sfsClient, err := config.SfsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud SFS Turbo client: %s", err)
	}

	createOpts := shares.CreateOpts{
		Name:             d.Get("name").(string),
		Size:             d.Get("size").(int),
		ShareProto:       d.Get("share_proto").(string),
		ShareType:        d.Get("share_type").(string),
		VpcID:            d.Get("vpc_id").(string),
		SubnetID:         d.Get("subnet_id").(string),
		SecurityGroupID:  d.Get("security_group_id").(string),
		AvailabilityZone: d.Get("availability_zone").(string),
	}

	metaOpts := shares.Metadata{}
	if v, ok := d.GetOk("crypt_key_id"); ok {
		metaOpts.CryptKeyID = v.(string)
	}
	if _, ok := d.GetOk("enhanced"); ok {
		metaOpts.ExpandType = "bandwidth"
	}
	createOpts.Metadata = metaOpts

	log.Printf("[DEBUG] create sfs turbo with option: %+v", createOpts)
	create, err := shares.Create(sfsClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud SFS Turbo: %s", err)
	}
	d.SetId(create.ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"100"},
		Target:     []string{"200"},
		Refresh:    waitForSFSTurboStatus(sfsClient, create.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      20 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, StateErr := stateConf.WaitForState()
	if StateErr != nil {
		return fmt.Errorf("Error waiting for SFS Turbo (%s) to become ready: %s ", d.Id(), StateErr)
	}

	return resourceSFSTurboRead(d, meta)
}

func resourceSFSTurboRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	sfsClient, err := config.SfsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud SFS Turbo client: %s", err)
	}

	n, err := shares.Get(sfsClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "Error retrieving SFS Turbo")
	}

	d.Set("name", n.Name)
	d.Set("share_proto", n.ShareProto)
	d.Set("share_type", n.ShareType)
	d.Set("vpc_id", n.VpcID)
	d.Set("subnet_id", n.SubnetID)
	d.Set("security_group_id", n.SecurityGroupID)
	d.Set("version", n.Version)
	d.Set("region", GetRegion(d, config))
	d.Set("availability_zone", n.AvailabilityZone)
	d.Set("available_capacity", n.AvailCapacity)
	d.Set("export_location", n.ExportLocation)
	d.Set("crypt_key_id", n.CryptKeyID)

	// n.Size is a string of float64, should convert it to int
	if fsize, err := strconv.ParseFloat(n.Size, 64); err == nil {
		if err = d.Set("size", int(fsize)); err != nil {
			return fmt.Errorf("Error reading size of SFS Turbo: %s", err)
		}
	}

	if n.ExpandType == "bandwidth" {
		d.Set("enhanced", true)
	} else {
		d.Set("enhanced", false)
	}

	var status string
	if n.SubStatus != "" {
		status = n.SubStatus
	} else {
		status = n.Status
	}
	d.Set("status", status)
	return nil
}

func resourceSFSTurboUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	sfsClient, err := config.SfsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error updating HuaweiCloud SFS Turbo client: %s", err)
	}

	if d.HasChange("size") {
		old, newsize := d.GetChange("size")
		if old.(int) > newsize.(int) {
			return fmt.Errorf("Shrinking HuaweiCloud SFS Turbo size is not supported")
		}

		expandOpts := shares.ExpandOpts{
			Extend: shares.ExtendOpts{NewSize: newsize.(int)},
		}
		expand := shares.Expand(sfsClient, d.Id(), expandOpts)
		if expand.Err != nil {
			return fmt.Errorf("Error Expanding HuaweiCloud Share File size: %s", expand.Err)
		}

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"121"},
			Target:     []string{"221", "200"},
			Refresh:    waitForSFSTurboSubStatus(sfsClient, d.Id()),
			Timeout:    d.Timeout(schema.TimeoutDelete),
			Delay:      10 * time.Second,
			MinTimeout: 5 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf("Error deleting HuaweiCloud SFS Turbo: %s", err)
		}
	}

	return resourceSFSTurboRead(d, meta)
}

func resourceSFSTurboDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	sfsClient, err := config.SfsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud SFS Turbo client: %s", err)
	}

	err = shares.Delete(sfsClient, d.Id()).ExtractErr()
	if err != nil {
		return CheckDeleted(d, err, "Error deleting SFS Turbo")
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"100", "200"},
		Target:     []string{"deleted"},
		Refresh:    waitForSFSTurboStatus(sfsClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error deleting HuaweiCloud SFS Turbo: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForSFSTurboStatus(sfsClient *golangsdk.ServiceClient, shareId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r, err := shares.Get(sfsClient, shareId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully deleted HuaweiCloud shared File %s", shareId)
				return r, "deleted", nil
			}
			return r, "error", err
		}

		return r, r.Status, nil
	}
}

func waitForSFSTurboSubStatus(sfsClient *golangsdk.ServiceClient, shareId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r, err := shares.Get(sfsClient, shareId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully deleted HuaweiCloud shared File %s", shareId)
				return r, "deleted", nil
			}
			return r, "error", err
		}

		var status string
		if r.SubStatus != "" {
			status = r.SubStatus
		} else {
			status = r.Status
		}
		return r, status, nil
	}
}
