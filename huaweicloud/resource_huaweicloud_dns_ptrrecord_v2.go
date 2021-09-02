package huaweicloud

import (
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/dns/v2/ptrrecords"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceDNSPtrRecordV2() *schema.Resource {
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
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"floatingip_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ttl": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 2147483647),
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"tags": tagsSchema(),
			"address": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDNSPtrRecordV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	region := GetRegion(d, config)
	dnsClient, err := config.DnsV2Client(region)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud DNS client: %s", err)
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
		PtrName:             d.Get("name").(string),
		Description:         d.Get("description").(string),
		TTL:                 d.Get("ttl").(int),
		Tags:                taglist,
		EnterpriseProjectID: GetEnterpriseProjectID(d, config),
	}

	logp.Printf("[DEBUG] Create Options: %#v", createOpts)
	fip_id := d.Get("floatingip_id").(string)
	n, err := ptrrecords.Create(dnsClient, region, fip_id, createOpts).Extract()
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud DNS PTR record: %s", err)
	}

	logp.Printf("[DEBUG] Waiting for DNS PTR record (%s) to become available", n.ID)
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
		return fmtp.Errorf(
			"Error waiting for PTR record (%s) to become ACTIVE for creation: %s",
			n.ID, err)
	}
	d.SetId(n.ID)

	logp.Printf("[DEBUG] Created HuaweiCloud DNS PTR record %s: %#v", n.ID, n)
	return resourceDNSPtrRecordV2Read(d, meta)
}

func resourceDNSPtrRecordV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	dnsClient, err := config.DnsV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud DNS client: %s", err)
	}

	n, err := ptrrecords.Get(dnsClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "ptr_record")
	}

	logp.Printf("[DEBUG] Retrieved PTR record %s: %#v", d.Id(), n)

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
	d.Set("enterprise_project_id", n.EnterpriseProjectID)

	// save tags
	if resourceTags, err := tags.Get(dnsClient, "DNS-ptr_record", d.Id()).Extract(); err == nil {
		tagmap := utils.TagsToMap(resourceTags.Tags)
		if err := d.Set("tags", tagmap); err != nil {
			return fmtp.Errorf("Error saving tags to state for DNS ptr record (%s): %s", d.Id(), err)
		}
	} else {
		logp.Printf("[WARN] Error fetching tags of DNS ptr record (%s): %s", d.Id(), err)
	}

	return nil
}

func resourceDNSPtrRecordV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	region := GetRegion(d, config)
	dnsClient, err := config.DnsV2Client(region)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud DNS client: %s", err)
	}

	if d.HasChanges("name", "description", "ttl") {
		updateOpts := ptrrecords.CreateOpts{
			PtrName:     d.Get("name").(string),
			Description: d.Get("description").(string),
			TTL:         d.Get("ttl").(int),
		}

		logp.Printf("[DEBUG] Update Options: %#v", updateOpts)
		fip_id := d.Get("floatingip_id").(string)
		n, err := ptrrecords.Create(dnsClient, region, fip_id, updateOpts).Extract()
		if err != nil {
			return fmtp.Errorf("Error updating HuaweiCloud DNS PTR record: %s", err)
		}

		logp.Printf("[DEBUG] Waiting for DNS PTR record (%s) to become available", n.ID)
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
			return fmtp.Errorf(
				"Error waiting for PTR record (%s) to become ACTIVE for update: %s",
				n.ID, err)
		}

		logp.Printf("[DEBUG] Updated HuaweiCloud DNS PTR record %s: %#v", n.ID, n)
	}

	// update tags
	tagErr := utils.UpdateResourceTags(dnsClient, d, "DNS-ptr_record", d.Id())
	if tagErr != nil {
		return fmtp.Errorf("Error updating tags of DNS PTR record %s: %s", d.Id(), tagErr)
	}

	return resourceDNSPtrRecordV2Read(d, meta)

}

func resourceDNSPtrRecordV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	dnsClient, err := config.DnsV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud DNS client: %s", err)
	}

	err = ptrrecords.Delete(dnsClient, d.Id()).ExtractErr()
	if err != nil {
		return fmtp.Errorf("Error deleting HuaweiCloud DNS PTR record: %s", err)
	}

	logp.Printf("[DEBUG] Waiting for DNS PTR record (%s) to be deleted", d.Id())
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
		return fmtp.Errorf(
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

		logp.Printf("[DEBUG] HuaweiCloud DNS PTR record (%s) current status: %s", ptrrecord.ID, ptrrecord.Status)
		return ptrrecord, ptrrecord.Status, nil
	}
}

func parseDNSV2PtrRecordID(id string) (string, error) {
	idParts := strings.Split(id, ":")
	if len(idParts) != 2 {
		return "", fmtp.Errorf("Unable to determine DNS PTR record ID from raw ID: %s", id)
	}

	fipID := idParts[1]
	return fipID, nil
}
