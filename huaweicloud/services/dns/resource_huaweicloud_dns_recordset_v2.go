package dns

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/dns/v2/recordsets"
	"github.com/chnsz/golangsdk/openstack/dns/v2/zones"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DNS PUT /v2.1/recordsets/{recordset_id}/statuses/set
// @API DNS POST /v2.1/zones/{zone_id}/recordsets
// @API DNS DELETE /v2.1/zones/{zone_id}/recordsets/{recordset_id}
// @API DNS GET /v2.1/zones/{zone_id}/recordsets/{recordset_id}
// @API DNS PUT /v2.1/zones/{zone_id}/recordsets/{recordset_id}
func ResourceDNSRecordSetV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDNSRecordSetV2Create,
		ReadContext:   resourceDNSRecordSetV2Read,
		UpdateContext: resourceDNSRecordSetV2Update,
		DeleteContext: resourceDNSRecordSetV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: config.MergeDefaultTags(),

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
				DiffSuppressFunc: func(_, oldVal, newVal string, _ *schema.ResourceData) bool {
					return strings.TrimSuffix(oldVal, ".") == strings.TrimSuffix(newVal, ".")
				},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"records": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				MinItems: 1,
			},
			"ttl": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  300,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"A", "AAAA", "MX", "CNAME", "TXT", "NS", "SRV", "PTR", "CAA",
				}, false),
			},
			"tags": common.TagsSchema(),
		},
	}
}

func resourceDNSRecordSetV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zoneID := d.Get("zone_id").(string)
	dnsClient, zoneType, err := chooseDNSClientbyZoneID(d, zoneID, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	recordsraw := d.Get("records").([]interface{})
	records := make([]string, len(recordsraw))
	for i, recordraw := range recordsraw {
		records[i] = recordraw.(string)
	}

	createOpts := recordsets.CreateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Records:     records,
		TTL:         d.Get("ttl").(int),
		Type:        d.Get("type").(string),
	}

	log.Printf("[DEBUG] Create options: %#v", createOpts)
	n, err := recordsets.Create(dnsClient, zoneID, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating DNS record set: %s", err)
	}

	id := fmt.Sprintf("%s/%s", zoneID, n.ID)
	d.SetId(id)

	log.Printf("[DEBUG] Waiting for DNS record set (%s) to become ACTIVE", n.ID)
	stateConf := &resource.StateChangeConf{
		Target:     []string{"ACTIVE"},
		Pending:    []string{"PENDING"},
		Refresh:    waitForDNSRecordSet(dnsClient, zoneID, n.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)

	if err != nil {
		return diag.Errorf(
			"error waiting for record set (%s) to become ACTIVE for creation: %s",
			n.ID, err)
	}

	// set tags
	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		resourceType, err := utils.GetDNSRecordSetTagType(zoneType)
		if err != nil {
			return diag.Errorf("error getting resource type of DNS record set %s: %s", n.ID, err)
		}

		taglist := utils.ExpandResourceTags(tagRaw)
		if tagErr := tags.Create(dnsClient, resourceType, n.ID, taglist).ExtractErr(); tagErr != nil {
			return diag.Errorf("error setting tags of DNS record set %s: %s", n.ID, tagErr)
		}
	}

	log.Printf("[DEBUG] Created DNS record set %s: %#v", n.ID, n)
	return resourceDNSRecordSetV2Read(ctx, d, meta)
}

func resourceDNSRecordSetV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	// Obtain relevant info from parsing the ID
	zoneID, recordsetID, err := parseDNSV2RecordSetID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	dnsClient, zoneType, err := chooseDNSClientbyZoneID(d, zoneID, meta)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error creating DNS client")
	}

	n, err := recordsets.Get(dnsClient, zoneID, recordsetID).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "record_set")
	}

	log.Printf("[DEBUG] Retrieved record set %s: %#v", recordsetID, n)

	mErr := multierror.Append(nil,
		d.Set("name", n.Name),
		d.Set("description", n.Description),
		d.Set("ttl", n.TTL),
		d.Set("type", n.Type),
		d.Set("records", n.Records),
		d.Set("region", conf.GetRegion(d)),
		d.Set("zone_id", zoneID),
	)
	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error setting resource: %s", mErr)
	}

	// save tags
	resourceType, err := utils.GetDNSRecordSetTagType(zoneType)
	if err != nil {
		return diag.Errorf("error getting resource type of DNS record set %s: %s", recordsetID, err)
	}
	if resourceTags, err := tags.Get(dnsClient, resourceType, recordsetID).Extract(); err == nil {
		tagmap := utils.TagsToMap(resourceTags.Tags)
		if err := d.Set("tags", tagmap); err != nil {
			return diag.Errorf("error saving tags to state for DNS record set (%s): %s", recordsetID, err)
		}
	} else {
		log.Printf("[WARN] Error fetching tags of DNS record set (%s): %s", recordsetID, err)
	}

	return nil
}

func resourceDNSRecordSetV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Obtain relevant info from parsing the ID
	zoneID, recordsetID, err := parseDNSV2RecordSetID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	dnsClient, zoneType, err := chooseDNSClientbyZoneID(d, zoneID, meta)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error creating DNS client")
	}

	if d.HasChanges("description", "ttl", "records") {
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

		log.Printf("[DEBUG] Updating record set %s with options: %#v", recordsetID, updateOpts)
		_, err = recordsets.Update(dnsClient, zoneID, recordsetID, updateOpts).Extract()
		if err != nil {
			return diag.Errorf("error updating DNS record set: %s", err)
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

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf(
				"error waiting for record set (%s) to become ACTIVE for updating: %s",
				recordsetID, err)
		}
	}

	// update tags
	resourceType, err := utils.GetDNSRecordSetTagType(zoneType)
	if err != nil {
		return diag.Errorf("error getting resource type of DNS record set %s: %s", d.Id(), err)
	}

	tagErr := utils.UpdateResourceTags(dnsClient, d, resourceType, recordsetID)
	if tagErr != nil {
		return diag.Errorf("error updating tags of DNS record set %s: %s", d.Id(), tagErr)
	}

	return resourceDNSRecordSetV2Read(ctx, d, meta)
}

func resourceDNSRecordSetV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Obtain relevant info from parsing the ID
	zoneID, recordsetID, err := parseDNSV2RecordSetID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	dnsClient, _, err := chooseDNSClientbyZoneID(d, zoneID, meta)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error creating DNS client")
	}

	err = recordsets.Delete(dnsClient, zoneID, recordsetID).ExtractErr()
	if err != nil {
		return diag.Errorf("error deleting DNS record set: %s", err)
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

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(
			"error waiting for record set (%s) to become DELETED for deletion: %s",
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

		log.Printf("[DEBUG] DNS record set (%s) current status: %s", recordset.ID, recordset.Status)
		return recordset, parseStatus(recordset.Status), nil
	}
}

func parseDNSV2RecordSetID(id string) (zoneID string, recordsetID string, err error) {
	idParts := strings.Split(id, "/")
	if len(idParts) != 2 {
		err = fmt.Errorf("unable to determine DNS record set ID from raw ID: %s", id)
		return
	}

	zoneID = idParts[0]
	recordsetID = idParts[1]
	return
}

// Use this function to build the client by DNS zone ID
// For a public zone, the endpoint of client should be https://dns.myhuaweicloud.com
// For a private zone, the endpoint of client should be https://dns.{region}.myhuaweicloud.com
// In most regions, the both endpoints can work well, but it's very useful for regions like `la-north-2`
func chooseDNSClientbyZoneID(d *schema.ResourceData, zoneID string, meta interface{}) (*golangsdk.ServiceClient, string, error) {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	var client *golangsdk.ServiceClient
	var zoneInfo *zones.Zone
	// Firstly, try to ues the DNS global endpoint
	client, err := conf.DnsV2Client(region)
	if err != nil {
		return nil, "", golangsdk.ErrDefault400{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Body: []byte(fmt.Sprintf("error creating DNS client: %s", err)),
			},
		}
	}

	// get zone with DNS global endpoint
	zoneInfo, err = zones.Get(client, zoneID).Extract()
	if err != nil {
		log.Printf("[WARN] fetching zone failed with DNS global endpoint: %s", err)

		// try to ues the DNS region endpoint
		var clientErr error
		client, clientErr = conf.DnsWithRegionClient(region)
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
