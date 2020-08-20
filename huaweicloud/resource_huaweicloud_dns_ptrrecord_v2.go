package huaweicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/dns/v2/ptrrecords"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceDNSPtrRecordV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceDNSPtrRecordV2Create,
		Read:   resourceDNSPtrRecordV2Read,
		Update: resourceDNSPtrRecordV2Update,
		Delete: resourceDNSPtrRecordV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"floatingip_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"ttl": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     false,
				ValidateFunc: resourceValidateTTL,
			},
			"tags": {
				Type:         schema.TypeMap,
				Optional:     true,
				ForceNew:     false,
				ValidateFunc: validateECSTagValue,
				Elem:         &schema.Schema{Type: schema.TypeString},
			},
			"address": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDNSPtrRecordV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	region := GetRegion(d, config)
	dnsClient, err := config.dnsV2Client(region)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud DNS client: %s", err)
	}

	tagmap := d.Get("tags").(map[string]interface{})
	taglist := []ptrrecords.Tag{}
	for k, v := range tagmap {
		tag := ptrrecords.Tag{
			Key:   k,
			Value: v.(string),
		}
		taglist = append(taglist, tag)
	}

	createOpts := ptrrecords.CreateOpts{
		PtrName:     d.Get("name").(string),
		Description: d.Get("description").(string),
		TTL:         d.Get("ttl").(int),
		Tags:        taglist,
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	fip_id := d.Get("floatingip_id").(string)
	n, err := ptrrecords.Create(dnsClient, region, fip_id, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud DNS PTR record: %s", err)
	}

	log.Printf("[DEBUG] Waiting for DNS PTR record (%s) to become available", n.ID)
	stateConf := &resource.StateChangeConf{
		Target:     []string{"ACTIVE"},
		Pending:    []string{"PENDING_CREATE"},
		Refresh:    waitForDNSPtrRecord(dnsClient, n.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()

	if err != nil {
		return fmt.Errorf(
			"Error waiting for PTR record (%s) to become ACTIVE for creation: %s",
			n.ID, err)
	}
	d.SetId(n.ID)

	log.Printf("[DEBUG] Created HuaweiCloud DNS PTR record %s: %#v", n.ID, n)
	return resourceDNSPtrRecordV2Read(d, meta)
}

func resourceDNSPtrRecordV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	dnsClient, err := config.dnsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud DNS client: %s", err)
	}

	n, err := ptrrecords.Get(dnsClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "ptr_record")
	}

	log.Printf("[DEBUG] Retrieved PTR record %s: %#v", d.Id(), n)

	// Obtain relevant info from parsing the ID
	fipID, err := parseDNSV2PtrRecordID(d.Id())
	if err != nil {
		return err
	}

	d.Set("name", n.PtrName)
	d.Set("description", n.Description)
	d.Set("floatingip_id", fipID)
	d.Set("ttl", n.TTL)
	d.Set("address", n.Address)

	return nil
}

func resourceDNSPtrRecordV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	region := GetRegion(d, config)
	dnsClient, err := config.dnsV2Client(region)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud DNS client: %s", err)
	}

	tagmap := d.Get("tags").(map[string]interface{})
	taglist := []ptrrecords.Tag{}
	for k, v := range tagmap {
		tag := ptrrecords.Tag{
			Key:   k,
			Value: v.(string),
		}
		taglist = append(taglist, tag)
	}

	createOpts := ptrrecords.CreateOpts{
		PtrName:     d.Get("name").(string),
		Description: d.Get("description").(string),
		TTL:         d.Get("ttl").(int),
		Tags:        taglist,
	}

	log.Printf("[DEBUG] Update Options: %#v", createOpts)
	fip_id := d.Get("floatingip_id").(string)
	n, err := ptrrecords.Create(dnsClient, region, fip_id, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error updating HuaweiCloud DNS PTR record: %s", err)
	}

	log.Printf("[DEBUG] Waiting for DNS PTR record (%s) to become available", n.ID)
	stateConf := &resource.StateChangeConf{
		Target:     []string{"ACTIVE"},
		Pending:    []string{"PENDING_CREATE"},
		Refresh:    waitForDNSPtrRecord(dnsClient, n.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()

	if err != nil {
		return fmt.Errorf(
			"Error waiting for PTR record (%s) to become ACTIVE for update: %s",
			n.ID, err)
	}

	log.Printf("[DEBUG] Updated HuaweiCloud DNS PTR record %s: %#v", n.ID, n)
	return resourceDNSPtrRecordV2Read(d, meta)

}

func resourceDNSPtrRecordV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	dnsClient, err := config.dnsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud DNS client: %s", err)
	}

	err = ptrrecords.Delete(dnsClient, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error deleting HuaweiCloud DNS PTR record: %s", err)
	}

	log.Printf("[DEBUG] Waiting for DNS PTR record (%s) to be deleted", d.Id())
	stateConf := &resource.StateChangeConf{
		Target:     []string{"DELETED"},
		Pending:    []string{"ACTIVE", "PENDING_DELETE", "ERROR"},
		Refresh:    waitForDNSPtrRecord(dnsClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for PTR record (%s) to become DELETED for deletion: %s",
			d.Id(), err)
	}

	d.SetId("")
	return nil
}

func waitForDNSPtrRecord(dnsClient *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ptrrecord, err := ptrrecords.Get(dnsClient, id).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return ptrrecord, "DELETED", nil
			}

			return nil, "", err
		}

		log.Printf("[DEBUG] HuaweiCloud DNS PTR record (%s) current status: %s", ptrrecord.ID, ptrrecord.Status)
		return ptrrecord, ptrrecord.Status, nil
	}
}

func parseDNSV2PtrRecordID(id string) (string, error) {
	idParts := strings.Split(id, ":")
	if len(idParts) != 2 {
		return "", fmt.Errorf("Unable to determine DNS PTR record ID from raw ID: %s", id)
	}

	fipID := idParts[1]
	return fipID, nil
}
