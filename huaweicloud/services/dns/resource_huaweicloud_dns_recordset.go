package dns

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DNS POST /v2.1/zones/{zone_id}/recordsets
// @API DNS GET /v2.1/zones/{zone_id}/recordsets/{recordset_id}
// @API DNS PUT /v2.1/zones/{zone_id}/recordsets/{recordset_id}
// @API DNS PUT /v2.1/recordsets/{recordset_id}/statuses/set
// @API DNS DELETE /v2.1/zones/{zone_id}/recordsets/{recordset_id}
// @API DNS POST /v2/zones/{zone_id}/recordsets
// @API DNS GET /v2/zones/{zone_id}/recordsets/{recordset_id}
// @API DNS PUT /v2/zones/{zone_id}/recordsets/{recordset_id}
// @API DNS DELETE /v2/zones/{zone_id}/recordsets/{recordset_id}

func ResourceDNSRecordset() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDNSRecordsetCreate,
		UpdateContext: resourceDNSRecordsetUpdate,
		ReadContext:   resourceDNSRecordsetRead,
		DeleteContext: resourceDNSRecordsetDelete,
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
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the zone to which the record set belongs.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the record set.`,
				DiffSuppressFunc: func(_, oldVal, newVal string, _ *schema.ResourceData) bool {
					return strings.TrimSuffix(oldVal, ".") == strings.TrimSuffix(newVal, ".")
				},
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"A", "AAAA", "MX", "CNAME", "TXT", "NS", "SRV", "CAA",
				}, false),
				Description: `The type of the record set.`,
			},
			"records": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				MinItems:    1,
				Required:    true,
				Description: `The list of the records of the record set.`,
			},
			"ttl": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     300,
				Description: `The time to live (TTL) of the record set (in seconds).`,
			},
			"line_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The resolution line ID.`,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "ENABLE",
				ValidateFunc: validation.StringInSlice([]string{"ENABLE", "DISABLE"}, false),
				Description:  `The status of the record set.`,
			},
			"tags": common.TagsSchema(),
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The description of the record set.`,
			},
			"weight": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `The weight of the record set.`,
			},
			"zone_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the zone to which the record set belongs.`,
			},
			"zone_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the zone to which the record set belongs.`,
			},
		},
	}
}

type WaitForConfig struct {
	ZoneID      string
	RecordsetID string
	ZoneType    string
	Timeout     time.Duration
}

func resourceDNSRecordsetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zoneID := d.Get("zone_id").(string)
	dnsClient, zoneType, err := chooseDNSClientbyZoneID(d, zoneID, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	if zoneType == "private" {
		if _, ok := d.GetOk("line_id"); ok {
			return diag.Errorf("private zone do not support line_id.")
		}
		if _, ok := d.GetOk("weight"); ok {
			return diag.Errorf("private zone do not support weight.")
		}
	}

	// createDNSRecordset: create DNS recordset.
	if err := createDNSRecordset(dnsClient, d, zoneType); err != nil {
		return diag.FromErr(err)
	}

	zoneID, recordsetID, err := parseDNSRecordsetID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	waitForConfig := &WaitForConfig{
		ZoneID:      zoneID,
		RecordsetID: recordsetID,
		ZoneType:    zoneType,
		Timeout:     d.Timeout(schema.TimeoutCreate),
	}
	if err := waitForDNSRecordsetCreateOrUpdate(ctx, dnsClient, waitForConfig); err != nil {
		return diag.FromErr(err)
	}

	return resourceDNSRecordsetRead(ctx, d, meta)
}

func createDNSRecordset(recordsetClient *golangsdk.ServiceClient, d *schema.ResourceData, zoneType string) error {
	version := getApiVersionByZoneType(zoneType)
	createDNSRecordsetHttpUrl := fmt.Sprintf("%s/zones/{zone_id}/recordsets", version)

	zoneID := d.Get("zone_id").(string)
	createDNSRecordsetPath := recordsetClient.Endpoint + createDNSRecordsetHttpUrl
	createDNSRecordsetPath = strings.ReplaceAll(createDNSRecordsetPath, "{zone_id}", zoneID)

	createDNSRecordsetOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			202,
		},
	}
	createDNSRecordsetOpt.JSONBody = utils.RemoveNil(buildCreateDNSRecordsetBodyParams(d))
	createDNSRecordsetResp, err := recordsetClient.Request("POST", createDNSRecordsetPath,
		&createDNSRecordsetOpt)
	if err != nil {
		return fmt.Errorf("error creating DNS recordset: %s", err)
	}

	createDNSRecordsetRespBody, err := utils.FlattenResponse(createDNSRecordsetResp)
	if err != nil {
		return err
	}

	recordSetID := utils.PathSearch("id", createDNSRecordsetRespBody, "").(string)
	if recordSetID == "" {
		return fmt.Errorf("unable to find the DNS recordset ID from the API response")
	}
	d.SetId(fmt.Sprintf("%s/%s", zoneID, recordSetID))
	return nil
}

func waitForDNSRecordsetCreateOrUpdate(ctx context.Context, recordsetClient *golangsdk.ServiceClient,
	waitForConfig *WaitForConfig) error {
	stateConf := &resource.StateChangeConf{
		Target:       []string{"ACTIVE", "DISABLE"},
		Pending:      []string{"PENDING"},
		Refresh:      dnsRecordsetStatusRefreshFunc(recordsetClient, waitForConfig),
		Timeout:      waitForConfig.Timeout,
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for DNS recordset (%s) to be ACTIVE or DISABLE : %s",
			waitForConfig.RecordsetID, err)
	}
	return nil
}

func buildCreateDNSRecordsetBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        utils.ValueIgnoreEmpty(d.Get("name")),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
		"type":        utils.ValueIgnoreEmpty(d.Get("type")),
		"status":      utils.ValueIgnoreEmpty(d.Get("status")),
		"ttl":         utils.ValueIgnoreEmpty(d.Get("ttl")),
		"records":     utils.ValueIgnoreEmpty(d.Get("records").(*schema.Set).List()),
		"line":        utils.ValueIgnoreEmpty(d.Get("line_id")),
		"tags":        utils.ExpandResourceTagsMap(d.Get("tags").(map[string]interface{})),
		"weight":      utils.ValueIgnoreEmpty(d.Get("weight")),
	}
	return bodyParams
}

func resourceDNSRecordsetRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	zoneID, recordsetID, err := parseDNSRecordsetID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	dnsClient, zoneType, err := chooseDNSClientbyZoneID(d, zoneID, meta)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error creating DNS client")
	}

	version := getApiVersionByZoneType(zoneType)
	getDNSRecordsetHttpUrl := fmt.Sprintf("%s/zones/{zone_id}/recordsets/{recordset_id}", version)

	getDNSRecordsetPath := dnsClient.Endpoint + getDNSRecordsetHttpUrl
	getDNSRecordsetPath = strings.ReplaceAll(getDNSRecordsetPath, "{zone_id}", zoneID)
	getDNSRecordsetPath = strings.ReplaceAll(getDNSRecordsetPath, "{recordset_id}", recordsetID)

	getDNSRecordsetOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getDNSRecordsetResp, err := dnsClient.Request("GET", getDNSRecordsetPath, &getDNSRecordsetOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DNS recordset")
	}

	getDNSRecordsetRespBody, err := utils.FlattenResponse(getDNSRecordsetResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", getDNSRecordsetRespBody, nil)),
		d.Set("description", utils.PathSearch("description", getDNSRecordsetRespBody, nil)),
		d.Set("zone_id", utils.PathSearch("zone_id", getDNSRecordsetRespBody, nil)),
		d.Set("zone_name", utils.PathSearch("zone_name", getDNSRecordsetRespBody, nil)),
		d.Set("type", utils.PathSearch("type", getDNSRecordsetRespBody, nil)),
		d.Set("ttl", utils.PathSearch("ttl", getDNSRecordsetRespBody, nil)),
		d.Set("records", utils.PathSearch("records", getDNSRecordsetRespBody, nil)),
		d.Set("status", getDNSRecordsetStatus(getDNSRecordsetRespBody)),
		d.Set("line_id", utils.PathSearch("line", getDNSRecordsetRespBody, nil)),
		d.Set("weight", utils.PathSearch("weight", getDNSRecordsetRespBody, nil)),
		d.Set("zone_type", zoneType),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	// set tags
	if err := setDNSRecordsetTags(d, dnsClient, recordsetID, zoneType); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func setDNSRecordsetTags(d *schema.ResourceData, client *golangsdk.ServiceClient, id, zoneType string) error {
	resourceType, err := utils.GetDNSRecordSetTagType(zoneType)
	if err != nil {
		return err
	}
	return utils.SetResourceTagsToState(d, client, resourceType, id)
}

func getDNSRecordsetStatus(getDNSRecordsetRespBody interface{}) string {
	status := utils.PathSearch("status", getDNSRecordsetRespBody, "").(string)
	if status == "ACTIVE" {
		return "ENABLE"
	}
	return status
}

func resourceDNSRecordsetUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zoneID, recordsetID, err := parseDNSRecordsetID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	dnsClient, zoneType, err := chooseDNSClientbyZoneID(d, zoneID, meta)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error creating DNS client")
	}

	if zoneType == "private" {
		if _, ok := d.GetOk("weight"); ok {
			return diag.Errorf("private zone do not support weight.")
		}
	}

	updateDNSRecordsetChanges := []string{
		"name",
		"description",
		"type",
		"ttl",
		"records",
		"weight",
	}
	if d.HasChanges(updateDNSRecordsetChanges...) {
		// updateDNSRecordset: Update DNS recordset
		if err := updateDNSRecordset(dnsClient, d, zoneID, recordsetID, zoneType); err != nil {
			return diag.FromErr(err)
		}

		waitForConfig := &WaitForConfig{
			ZoneID:      zoneID,
			RecordsetID: recordsetID,
			ZoneType:    zoneType,
			Timeout:     d.Timeout(schema.TimeoutUpdate),
		}
		if err := waitForDNSRecordsetCreateOrUpdate(ctx, dnsClient, waitForConfig); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("status") {
		// updateDNSRecordsetStatus: Update DNS recordset status
		if err := updateDNSRecordsetStatus(dnsClient, d, recordsetID); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("tags") {
		resourceType, err := utils.GetDNSRecordSetTagType(zoneType)
		if err != nil {
			return diag.FromErr(err)
		}

		err = utils.UpdateResourceTags(dnsClient, d, resourceType, recordsetID)
		if err != nil {
			return diag.Errorf("error updating DNS recordset tags: %s", err)
		}
	}
	return resourceDNSRecordsetRead(ctx, d, meta)
}

func updateDNSRecordset(recordsetClient *golangsdk.ServiceClient, d *schema.ResourceData, zoneID,
	recordsetID, zoneType string) error {
	version := getApiVersionByZoneType(zoneType)
	updateDNSRecordsetHttpUrl := fmt.Sprintf("%s/zones/{zone_id}/recordsets/{recordset_id}", version)

	updateDNSRecordsetPath := recordsetClient.Endpoint + updateDNSRecordsetHttpUrl
	updateDNSRecordsetPath = strings.ReplaceAll(updateDNSRecordsetPath, "{zone_id}", zoneID)
	updateDNSRecordsetPath = strings.ReplaceAll(updateDNSRecordsetPath, "{recordset_id}", recordsetID)

	updateDNSRecordsetOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			202,
		},
	}
	updateDNSRecordsetOpt.JSONBody = utils.RemoveNil(buildUpdateDNSRecordsetBodyParams(d))
	_, err := recordsetClient.Request("PUT", updateDNSRecordsetPath, &updateDNSRecordsetOpt)
	if err != nil {
		return fmt.Errorf("error updating DNS recordset: %s", err)
	}
	return nil
}

func updateDNSRecordsetStatus(recordsetClient *golangsdk.ServiceClient, d *schema.ResourceData,
	recordsetID string) error {
	var (
		updateDNSRecordsetStatusHttpUrl = "v2.1/recordsets/{recordset_id}/statuses/set"
	)

	updateDNSRecordsetStatusPath := recordsetClient.Endpoint + updateDNSRecordsetStatusHttpUrl
	updateDNSRecordsetStatusPath = strings.ReplaceAll(updateDNSRecordsetStatusPath, "{recordset_id}", recordsetID)

	updateDNSRecordsetStatusOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			202,
		},
	}
	updateDNSRecordsetStatusOpt.JSONBody = utils.RemoveNil(buildUpdateDNSRecordsetStatusBodyParams(d))
	_, err := recordsetClient.Request("PUT", updateDNSRecordsetStatusPath, &updateDNSRecordsetStatusOpt)
	if err != nil {
		return fmt.Errorf("error updating DNS recordset status: %s", err)
	}
	return nil
}

func buildUpdateDNSRecordsetBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        utils.ValueIgnoreEmpty(d.Get("name")),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
		"type":        utils.ValueIgnoreEmpty(d.Get("type")),
		"ttl":         utils.ValueIgnoreEmpty(d.Get("ttl")),
		"records":     utils.ValueIgnoreEmpty(d.Get("records").(*schema.Set).List()),
		"weight":      utils.ValueIgnoreEmpty(d.Get("weight")),
	}
	return bodyParams
}

func buildUpdateDNSRecordsetStatusBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"status": utils.ValueIgnoreEmpty(d.Get("status")),
	}
	return bodyParams
}

func resourceDNSRecordsetDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zoneID, recordsetID, err := parseDNSRecordsetID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	dnsClient, zoneType, err := chooseDNSClientbyZoneID(d, zoneID, meta)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error creating DNS client")
	}

	version := getApiVersionByZoneType(zoneType)
	deleteDNSRecordsetHttpUrl := fmt.Sprintf("%s/zones/{zone_id}/recordsets/{recordset_id}", version)

	deleteDNSRecordsetPath := dnsClient.Endpoint + deleteDNSRecordsetHttpUrl
	deleteDNSRecordsetPath = strings.ReplaceAll(deleteDNSRecordsetPath, "{zone_id}", zoneID)
	deleteDNSRecordsetPath = strings.ReplaceAll(deleteDNSRecordsetPath, "{recordset_id}", recordsetID)

	deleteDNSRecordsetOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			202,
		},
	}
	_, err = dnsClient.Request("DELETE", deleteDNSRecordsetPath, &deleteDNSRecordsetOpt)
	if err != nil {
		return diag.Errorf("error deleting DNS recordset: %s", err)
	}

	waitForConfig := &WaitForConfig{
		ZoneID:      zoneID,
		RecordsetID: recordsetID,
		ZoneType:    zoneType,
		Timeout:     d.Timeout(schema.TimeoutDelete),
	}
	if err := waitForDNSRecordsetDeleted(ctx, dnsClient, waitForConfig); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func waitForDNSRecordsetDeleted(ctx context.Context, recordsetClient *golangsdk.ServiceClient,
	waitForConfig *WaitForConfig) error {
	stateConf := &resource.StateChangeConf{
		Target:       []string{"DELETED"},
		Pending:      []string{"ACTIVE", "PENDING", "ERROR"},
		Refresh:      dnsRecordsetStatusRefreshFunc(recordsetClient, waitForConfig),
		Timeout:      waitForConfig.Timeout,
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for DNS recordset (%s) to be DELETED: %s",
			waitForConfig.RecordsetID, err)
	}
	return nil
}

func dnsRecordsetStatusRefreshFunc(client *golangsdk.ServiceClient, waitForConfig *WaitForConfig) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		version := getApiVersionByZoneType(waitForConfig.ZoneType)
		getDNSRecordsetHttpUrl := fmt.Sprintf("%s/zones/{zone_id}/recordsets/{recordset_id}", version)

		getDNSRecordsetPath := client.Endpoint + getDNSRecordsetHttpUrl
		getDNSRecordsetPath = strings.ReplaceAll(getDNSRecordsetPath, "{zone_id}", waitForConfig.ZoneID)
		getDNSRecordsetPath = strings.ReplaceAll(getDNSRecordsetPath, "{recordset_id}", waitForConfig.RecordsetID)

		getDNSRecordsetOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}
		getDNSRecordsetResp, err := client.Request("GET", getDNSRecordsetPath, &getDNSRecordsetOpt)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				// When the error code is 404, the value of respBody is nil, and a non-null value is returned to avoid continuing the loop check.
				return "Resource Not Found", "DELETED", nil
			}
			return nil, "", err
		}

		getDNSRecordsetRespBody, err := utils.FlattenResponse(getDNSRecordsetResp)
		if err != nil {
			return nil, "", err
		}

		status := utils.PathSearch("status", getDNSRecordsetRespBody, "")
		return getDNSRecordsetRespBody, parseStatus(status.(string)), nil
	}
}

func parseDNSRecordsetID(id string) (zoneID, recordsetID string, err error) {
	idArrays := strings.SplitN(id, "/", 2)
	if len(idArrays) != 2 {
		err = fmt.Errorf("invalid format specified for ID. Format must be <zone_id>/<recordset_id>")
		return
	}
	zoneID = idArrays[0]
	recordsetID = idArrays[1]
	return
}

func getApiVersionByZoneType(zoneType string) string {
	if zoneType == "private" {
		return "v2"
	}
	// v2.1 can support Multi-line Record Set which applies only to public zones
	return "v2.1"
}
