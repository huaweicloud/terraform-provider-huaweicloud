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

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/dns/v2/ptrrecords"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DNS PATCH /v2/reverse/floatingips/{region}:{floatingip_id}
// @API DNS GET /v2/reverse/floatingips/{region}:{floatingip_id}
// @API DNS GET /v2/{project_id}/DNS-ptr_record/{resource_id}/tags
// @API DNS POST /v2/{project_id}/DNS-ptr_record/{resource_id}/tags/action
func ResourcePtrRecord() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePtrRecordCreate,
		ReadContext:   resourcePtrRecordRead,
		UpdateContext: resourcePtrRecordUpdate,
		DeleteContext: resourcePtrRecordDelete,
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
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(_, oldVal, newVal string, _ *schema.ResourceData) bool {
					return strings.TrimSuffix(oldVal, ".") == strings.TrimSuffix(newVal, ".")
				},
				Description: `The domain name of the PTR record.`,
			},
			"floatingip_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the FloatingIP/EIP.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the PTR record.`,
			},
			"ttl": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `The time to live (TTL) of the record set (in seconds).`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: `The enterprise project ID of the PTR record.`,
			},
			"tags": common.TagsSchema(`The key/value pairs to associate with the PTR record.`),
			"address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The address of the FloatingIP/EIP.`,
			},
		},
	}
}

func resourcePtrRecordCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	dnsClient, err := cfg.DnsV2Client(region)
	if err != nil {
		return diag.Errorf("error creating DNS client: %s", err)
	}

	createOpts := ptrrecords.CreateOpts{
		PtrName:             d.Get("name").(string),
		Description:         utils.String(d.Get("description").(string)),
		TTL:                 d.Get("ttl").(int),
		Tags:                getPtrRecordsTagList(d),
		EnterpriseProjectID: cfg.GetEnterpriseProjectID(d),
	}

	log.Printf("[DEBUG] Create options: %#v", createOpts)
	floatingIpId := d.Get("floatingip_id").(string)
	respBody, err := ptrrecords.Create(dnsClient, region, floatingIpId, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating DNS PTR record: %s", err)
	}

	ptrRecordId := respBody.ID
	if ptrRecordId == "" {
		return diag.Errorf("unable to find PTR record ID from API response")
	}

	err = waitForPtrRecordCreateOrUpdate(ctx, dnsClient, ptrRecordId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for PTR record (%s) to be created: %s", ptrRecordId, err)
	}
	d.SetId(ptrRecordId)

	log.Printf("[DEBUG] Created DNS PTR record %s: %#v", ptrRecordId, respBody)
	return resourcePtrRecordRead(ctx, d, meta)
}

func getPtrRecordsTagList(d *schema.ResourceData) []ptrrecords.Tag {
	tagMap := d.Get("tags").(map[string]interface{})
	var tagList []ptrrecords.Tag
	for k, v := range tagMap {
		tagList = append(tagList, ptrrecords.Tag{
			Key:   k,
			Value: v.(string),
		})
	}
	return tagList
}

func waitForPtrRecordCreateOrUpdate(ctx context.Context, dnsClient *golangsdk.ServiceClient, ptrRecordId string,
	timeout time.Duration) error {
	log.Printf("[DEBUG] Waiting for DNS PTR record (%s) to become ACTIVE", ptrRecordId)
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      waitForPtrRecord(dnsClient, ptrRecordId, false),
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 3 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourcePtrRecordRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		conf        = meta.(*config.Config)
		region      = conf.GetRegion(d)
		ptrRecordId = d.Id()
	)
	dnsClient, err := conf.DnsV2Client(region)
	if err != nil {
		return diag.Errorf("error creating DNS client: %s", err)
	}

	respBody, err := ptrrecords.Get(dnsClient, ptrRecordId).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DNS PTR record")
	}

	log.Printf("[DEBUG] Retrieved PTR record %s: %#v", ptrRecordId, respBody)

	// Obtain relevant info from parsing the ID
	floatingIpId, err := parsePtrRecordId(ptrRecordId)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", respBody.PtrName),
		d.Set("description", respBody.Description),
		d.Set("floatingip_id", floatingIpId),
		d.Set("ttl", respBody.TTL),
		d.Set("address", respBody.Address),
		d.Set("enterprise_project_id", respBody.EnterpriseProjectID),
	)
	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error setting resource: %s", mErr)
	}

	if err := utils.SetResourceTagsToState(d, dnsClient, "DNS-ptr_record", ptrRecordId); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourcePtrRecordUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		conf        = meta.(*config.Config)
		region      = conf.GetRegion(d)
		ptrRecordId = d.Id()
	)
	dnsClient, err := conf.DnsV2Client(region)
	if err != nil {
		return diag.Errorf("error creating DNS client: %s", err)
	}

	if d.HasChanges("name", "description", "ttl") {
		updateOpts := ptrrecords.CreateOpts{
			PtrName:     d.Get("name").(string),
			Description: utils.String(d.Get("description").(string)),
			TTL:         d.Get("ttl").(int),
		}

		log.Printf("[DEBUG] Update options: %#v", updateOpts)
		respBody, err := ptrrecords.Create(dnsClient, region, d.Get("floatingip_id").(string), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("error updating DNS PTR record: %s", err)
		}

		err = waitForPtrRecordCreateOrUpdate(ctx, dnsClient, ptrRecordId, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.Errorf("error waiting for PTR record (%s) to be updated: %s", ptrRecordId, err)
		}

		log.Printf("[DEBUG] Updated DNS PTR record %s: %#v", ptrRecordId, respBody)
	}

	// update tags
	tagErr := utils.UpdateResourceTags(dnsClient, d, "DNS-ptr_record", d.Id())
	if tagErr != nil {
		return diag.Errorf("error updating tags of DNS PTR record %s: %s", d.Id(), tagErr)
	}

	return resourcePtrRecordRead(ctx, d, meta)
}

func resourcePtrRecordDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		conf        = meta.(*config.Config)
		ptrRecordId = d.Id()
	)
	dnsClient, err := conf.DnsV2Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DNS client: %s", err)
	}

	err = ptrrecords.Delete(dnsClient, ptrRecordId).ExtractErr()
	if err != nil {
		return diag.Errorf("error deleting DNS PTR record: %s", err)
	}

	log.Printf("[DEBUG] Waiting for DNS PTR record (%s) to be deleted", ptrRecordId)
	stateConf := &resource.StateChangeConf{
		// Allows deletion of PTR record with status 'ERROR'.
		Pending:      []string{"PENDING"},
		Target:       []string{"DELETED"},
		Refresh:      waitForPtrRecord(dnsClient, ptrRecordId, true),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        5 * time.Second,
		PollInterval: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for PTR record (%s) to deleted: %s", ptrRecordId, err)
	}
	return nil
}

func waitForPtrRecord(dnsClient *golangsdk.ServiceClient, ptrRecordId string, isDelete bool) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ptrRecord, err := ptrrecords.Get(dnsClient, ptrRecordId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok && isDelete {
				return "Resource Not Found", "DELETED", nil
			}
			return nil, "ERROR", err
		}

		status := ptrRecord.Status
		log.Printf("[DEBUG] DNS PTR record (%s) current status: %s", ptrRecord.ID, status)
		if !isDelete {
			if status == "ACTIVE" {
				return ptrRecord, "COMPLETED", nil
			}

			if status == "ERROR" {
				return ptrRecord, "ERROR", fmt.Errorf("unexpect status (%s)", status)
			}
		}
		return ptrRecord, "PENDING", nil
	}
}

func parsePtrRecordId(ptrRecordId string) (string, error) {
	parts := strings.Split(ptrRecordId, ":")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid PTR record ID format (%s), want '<region>:<floatingip_id>'", ptrRecordId)
	}

	return parts[1], nil
}
