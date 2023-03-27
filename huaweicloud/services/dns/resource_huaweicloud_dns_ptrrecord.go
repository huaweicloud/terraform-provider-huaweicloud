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
	"github.com/chnsz/golangsdk/openstack/dns/v2/ptrrecords"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func ResourceDNSPtrRecord() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDNSPtrRecordCreate,
		ReadContext:   resourceDNSPtrRecordRead,
		UpdateContext: resourceDNSPtrRecordUpdate,
		DeleteContext: resourceDNSPtrRecordDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
			"tags": common.TagsSchema(),
			"address": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDNSPtrRecordCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	dnsClient, err := conf.DnsV2Client(region)
	if err != nil {
		return diag.Errorf("error creating DNS client: %s", err)
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
		EnterpriseProjectID: common.GetEnterpriseProjectID(d, conf),
	}

	log.Printf("[DEBUG] Create options: %#v", createOpts)
	fipId := d.Get("floatingip_id").(string)
	n, err := ptrrecords.Create(dnsClient, region, fipId, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating DNS PTR record: %s", err)
	}

	log.Printf("[DEBUG] Waiting for DNS PTR record (%s) to become ACTIVE", n.ID)
	stateConf := &resource.StateChangeConf{
		Target:     []string{"ACTIVE"},
		Pending:    []string{"PENDING_CREATE"},
		Refresh:    waitForDNSPtrRecord(dnsClient, n.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)

	if err != nil {
		return diag.Errorf(
			"error waiting for PTR record (%s) to become ACTIVE for creation: %s",
			n.ID, err)
	}
	d.SetId(n.ID)

	log.Printf("[DEBUG] Created DNS PTR record %s: %#v", n.ID, n)
	return resourceDNSPtrRecordRead(ctx, d, meta)
}

func resourceDNSPtrRecordRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	dnsClient, err := conf.DnsV2Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DNS client: %s", err)
	}

	n, err := ptrrecords.Get(dnsClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "ptr_record")
	}

	log.Printf("[DEBUG] Retrieved PTR record %s: %#v", d.Id(), n)

	// Obtain relevant info from parsing the ID
	fipID, err := parseDNSPtrRecordID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(nil,
		d.Set("name", n.PtrName),
		d.Set("description", n.Description),
		d.Set("floatingip_id", fipID),
		d.Set("ttl", n.TTL),
		d.Set("address", n.Address),
		d.Set("enterprise_project_id", n.EnterpriseProjectID),
	)
	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error setting resource: %s", mErr)
	}

	// save tags
	if resourceTags, err := tags.Get(dnsClient, "DNS-ptr_record", d.Id()).Extract(); err == nil {
		tagmap := utils.TagsToMap(resourceTags.Tags)
		if err := d.Set("tags", tagmap); err != nil {
			return diag.Errorf("error saving tags to state for DNS PTR record (%s): %s", d.Id(), err)
		}
	} else {
		log.Printf("[WARN] Error fetching tags of DNS PTR record (%s): %s", d.Id(), err)
	}

	return nil
}

func resourceDNSPtrRecordUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	dnsClient, err := conf.DnsV2Client(region)
	if err != nil {
		return diag.Errorf("error creating DNS client: %s", err)
	}

	if d.HasChanges("name", "description", "ttl") {
		updateOpts := ptrrecords.CreateOpts{
			PtrName:     d.Get("name").(string),
			Description: d.Get("description").(string),
			TTL:         d.Get("ttl").(int),
		}

		log.Printf("[DEBUG] Update options: %#v", updateOpts)
		fipId := d.Get("floatingip_id").(string)
		n, err := ptrrecords.Create(dnsClient, region, fipId, updateOpts).Extract()
		if err != nil {
			return diag.Errorf("error updating DNS PTR record: %s", err)
		}

		log.Printf("[DEBUG] Waiting for DNS PTR record (%s) to become ACTIVE", n.ID)
		stateConf := &resource.StateChangeConf{
			Target:     []string{"ACTIVE"},
			Pending:    []string{"PENDING_CREATE"},
			Refresh:    waitForDNSPtrRecord(dnsClient, n.ID),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      5 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf(
				"error waiting for PTR record (%s) to become ACTIVE for update: %s",
				n.ID, err)
		}

		log.Printf("[DEBUG] Updated DNS PTR record %s: %#v", n.ID, n)
	}

	// update tags
	tagErr := utils.UpdateResourceTags(dnsClient, d, "DNS-ptr_record", d.Id())
	if tagErr != nil {
		return diag.Errorf("error updating tags of DNS PTR record %s: %s", d.Id(), tagErr)
	}

	return resourceDNSPtrRecordRead(ctx, d, meta)
}

func resourceDNSPtrRecordDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	dnsClient, err := conf.DnsV2Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DNS client: %s", err)
	}

	err = ptrrecords.Delete(dnsClient, d.Id()).ExtractErr()
	if err != nil {
		return diag.Errorf("error deleting DNS PTR record: %s", err)
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

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(
			"error waiting for PTR record (%s) to become DELETED for deletion: %s",
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

		log.Printf("[DEBUG] DNS PTR record (%s) current status: %s", ptrrecord.ID, ptrrecord.Status)
		return ptrrecord, ptrrecord.Status, nil
	}
}

func parseDNSPtrRecordID(id string) (string, error) {
	idParts := strings.Split(id, ":")
	if len(idParts) != 2 {
		return "", fmt.Errorf("unable to determine DNS PTR record ID from raw ID: %s", id)
	}

	fipID := idParts[1]
	return fipID, nil
}
