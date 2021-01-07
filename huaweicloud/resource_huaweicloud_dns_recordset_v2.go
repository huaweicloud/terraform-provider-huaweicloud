package huaweicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/common/tags"
	"github.com/huaweicloud/golangsdk/openstack/dns/v2/recordsets"
	"github.com/huaweicloud/golangsdk/openstack/dns/v2/zones"
)

func ResourceDNSRecordSetV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceDNSRecordSetV2Create,
		Read:   resourceDNSRecordSetV2Read,
		Update: resourceDNSRecordSetV2Update,
		Delete: resourceDNSRecordSetV2Delete,
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
				ForceNew: true,
				Computed: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 255),
			},
			"records": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				MinItems: 1,
			},
			"ttl": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      300,
				ValidateFunc: validation.IntBetween(1, 2147483647),
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"A", "AAAA", "MX", "CNAME", "TXT", "NS", "SRV", "PTR", "CAA",
				}, false),
			},
			"value_specs": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceDNSRecordSetV2Create(d *schema.ResourceData, meta interface{}) error {
	zoneID := d.Get("zone_id").(string)
	dnsClient, zoneType, err := chooseDNSClientbyZoneID(d, zoneID, meta)
	if err != nil {
		return err
	}

	recordsraw := d.Get("records").([]interface{})
	records := make([]string, len(recordsraw))
	for i, recordraw := range recordsraw {
		records[i] = recordraw.(string)
	}

	createOpts := RecordSetCreateOpts{
		recordsets.CreateOpts{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
			Records:     records,
			TTL:         d.Get("ttl").(int),
			Type:        d.Get("type").(string),
		},
		MapValueSpecs(d),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	n, err := recordsets.Create(dnsClient, zoneID, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud DNS record set: %s", err)
	}

	id := fmt.Sprintf("%s/%s", zoneID, n.ID)
	d.SetId(id)

	log.Printf("[DEBUG] Waiting for DNS record set (%s) to become available", n.ID)
	stateConf := &resource.StateChangeConf{
		Target:     []string{"ACTIVE"},
		Pending:    []string{"PENDING"},
		Refresh:    waitForDNSRecordSet(dnsClient, zoneID, n.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()

	if err != nil {
		return fmt.Errorf(
			"Error waiting for record set (%s) to become ACTIVE for creation: %s",
			n.ID, err)
	}

	// set tags
	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		resourceType, err := getDNSRecordSetTagType(zoneType)
		if err != nil {
			return fmt.Errorf("Error getting resource type of DNS record set %s: %s", n.ID, err)
		}

		taglist := expandResourceTags(tagRaw)
		if tagErr := tags.Create(dnsClient, resourceType, n.ID, taglist).ExtractErr(); tagErr != nil {
			return fmt.Errorf("Error setting tags of DNS record set %s: %s", n.ID, tagErr)
		}
	}

	log.Printf("[DEBUG] Created HuaweiCloud DNS record set %s: %#v", n.ID, n)
	return resourceDNSRecordSetV2Read(d, meta)
}

func resourceDNSRecordSetV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	// Obtain relevant info from parsing the ID
	zoneID, recordsetID, err := parseDNSV2RecordSetID(d.Id())
	if err != nil {
		return err
	}

	dnsClient, zoneType, err := chooseDNSClientbyZoneID(d, zoneID, meta)
	if err != nil {
		return err
	}

	time.Sleep(2 * time.Second)
	n, err := recordsets.Get(dnsClient, zoneID, recordsetID).Extract()
	if err != nil {
		return CheckDeleted(d, err, "record_set")
	}

	log.Printf("[DEBUG] Retrieved  record set %s: %#v", recordsetID, n)

	d.Set("name", n.Name)
	d.Set("description", n.Description)
	d.Set("ttl", n.TTL)
	d.Set("type", n.Type)
	if err := d.Set("records", n.Records); err != nil {
		return fmt.Errorf("[DEBUG] Error saving records to state for HuaweiCloud DNS record set (%s): %s", d.Id(), err)
	}
	d.Set("region", GetRegion(d, config))
	d.Set("zone_id", zoneID)

	// save tags
	resourceType, err := getDNSRecordSetTagType(zoneType)
	if err != nil {
		return fmt.Errorf("Error getting resource type of DNS record set %s: %s", recordsetID, err)
	}
	resourceTags, err := tags.Get(dnsClient, resourceType, recordsetID).Extract()
	if err != nil {
		return fmt.Errorf("Error fetching HuaweiCloud DNS record set tags: %s", err)
	}

	tagmap := tagsToMap(resourceTags.Tags)
	if err := d.Set("tags", tagmap); err != nil {
		return fmt.Errorf("Error saving tags for HuaweiCloud DNS record set %s: %s", recordsetID, err)
	}

	return nil
}

func resourceDNSRecordSetV2Update(d *schema.ResourceData, meta interface{}) error {
	// Obtain relevant info from parsing the ID
	zoneID, recordsetID, err := parseDNSV2RecordSetID(d.Id())
	if err != nil {
		return err
	}

	dnsClient, zoneType, err := chooseDNSClientbyZoneID(d, zoneID, meta)
	if err != nil {
		return err
	}

	var updateOpts recordsets.UpdateOpts
	if d.HasChange("ttl") {
		updateOpts.TTL = d.Get("ttl").(int)
	}

	if d.HasChange("records") {
		recordsraw := d.Get("records").([]interface{})
		records := make([]string, len(recordsraw))
		for i, recordraw := range recordsraw {
			records[i] = recordraw.(string)
		}
		updateOpts.Records = records
	}

	if d.HasChange("description") {
		updateOpts.Description = d.Get("description").(string)
	}

	log.Printf("[DEBUG] Updating  record set %s with options: %#v", recordsetID, updateOpts)

	_, err = recordsets.Update(dnsClient, zoneID, recordsetID, updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error updating HuaweiCloud DNS  record set: %s", err)
	}

	log.Printf("[DEBUG] Waiting for DNS record set (%s) to update", recordsetID)
	stateConf := &resource.StateChangeConf{
		Target:     []string{"ACTIVE"},
		Pending:    []string{"PENDING"},
		Refresh:    waitForDNSRecordSet(dnsClient, zoneID, recordsetID),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for record set (%s) to become ACTIVE for updation: %s",
			recordsetID, err)
	}

	// update tags
	resourceType, err := getDNSRecordSetTagType(zoneType)
	if err != nil {
		return fmt.Errorf("Error getting resource type of DNS record set %s: %s", d.Id(), err)
	}

	tagErr := UpdateResourceTags(dnsClient, d, resourceType, recordsetID)
	if tagErr != nil {
		return fmt.Errorf("Error updating tags of DNS record set %s: %s", d.Id(), tagErr)
	}

	return resourceDNSRecordSetV2Read(d, meta)
}

func resourceDNSRecordSetV2Delete(d *schema.ResourceData, meta interface{}) error {
	// Obtain relevant info from parsing the ID
	zoneID, recordsetID, err := parseDNSV2RecordSetID(d.Id())
	if err != nil {
		return err
	}

	dnsClient, _, err := chooseDNSClientbyZoneID(d, zoneID, meta)
	if err != nil {
		return err
	}

	err = recordsets.Delete(dnsClient, zoneID, recordsetID).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error deleting HuaweiCloud DNS record set: %s", err)
	}

	log.Printf("[DEBUG] Waiting for DNS record set (%s) to be deleted", recordsetID)
	stateConf := &resource.StateChangeConf{
		Target:     []string{"DELETED"},
		Pending:    []string{"ACTIVE", "PENDING", "ERROR"},
		Refresh:    waitForDNSRecordSet(dnsClient, zoneID, recordsetID),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for record set (%s) to become DELETED for deletion: %s",
			recordsetID, err)
	}

	d.SetId("")
	return nil
}

func parseStatus(rawStatus string) string {
	splits := strings.Split(rawStatus, "_")
	// rawStatus maybe one of PENDING_CREATE, PENDING_UPDATE, PENDING_DELETE, ACTIVE, or ERROR
	return splits[0]
}

func waitForDNSRecordSet(dnsClient *golangsdk.ServiceClient, zoneID, recordsetId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		recordset, err := recordsets.Get(dnsClient, zoneID, recordsetId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return recordset, "DELETED", nil
			}

			return nil, "", err
		}

		log.Printf("[DEBUG] HuaweiCloud DNS record set (%s) current status: %s", recordset.ID, recordset.Status)
		return recordset, parseStatus(recordset.Status), nil
	}
}

func parseDNSV2RecordSetID(id string) (string, string, error) {
	idParts := strings.Split(id, "/")
	if len(idParts) != 2 {
		return "", "", fmt.Errorf("Unable to determine DNS record set ID from raw ID: %s", id)
	}

	zoneID := idParts[0]
	recordsetID := idParts[1]

	return zoneID, recordsetID, nil
}

func chooseDNSClientbyZoneID(d *schema.ResourceData, zoneID string, meta interface{}) (*golangsdk.ServiceClient, string, error) {
	config := meta.(*Config)
	region := GetRegion(d, config)

	var client *golangsdk.ServiceClient
	var zoneInfo *zones.Zone
	// Firstly, try to ues the DNS global endpoint
	client, err := config.DnsV2Client(region)
	if err != nil {
		return nil, "", fmt.Errorf("Error creating HuaweiCloud DNS client: %s", err)
	}

	// get zone with DNS global endpoint
	zoneInfo, err = zones.Get(client, zoneID).Extract()
	if err != nil {
		log.Printf("[WARN] fetching zone failed with DNS global endpoint: %s", err)

		// try to ues the DNS region endpoint
		client, clientErr := config.DnsWithRegionClient(region)
		if clientErr != nil {
			// it looks tricky as we return the fetching error rather than clientErr
			return nil, "", err
		}

		// get zone with DNS region endpoint
		zoneInfo, err = zones.Get(client, zoneID).Extract()
		if err != nil {
			return nil, "", err
		}
	}

	return client, zoneInfo.ZoneType, nil
}
