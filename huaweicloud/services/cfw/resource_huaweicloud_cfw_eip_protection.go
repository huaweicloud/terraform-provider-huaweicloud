package cfw

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
	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

type eipProtectionPage struct {
	pagination.OffsetPageBase
}

const (
	eipProtectionProduct string = "cfw"

	openEipProtection  int = 0
	closeEipProtection int = 1

	modifyProtectHttpUrl string = "v1/{project_id}/eip/protect"
	queryHttpUrl         string = "v1/{project_id}/eips/protect"
)

// @API CFW POST /v1/{project_id}/eip/protect
// @API CFW GET /v1/{project_id}/eips/protect
func ResourceEipProtection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEipProtectionCreate,
		ReadContext:   resourceEipProtectionRead,
		UpdateContext: resourceEipProtectionUpdate,
		DeleteContext: resourceEipProtectionDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Importer: &schema.ResourceImporter{
			StateContext: resourceEipProtectionImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"object_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The protected object ID.`,
			},
			"protected_eip": {
				Type:        schema.TypeSet,
				Elem:        eipProtectionProtectedEipSchema(),
				Required:    true,
				Description: `The protected EIP configurations.`,
				Set: func(v interface{}) int {
					m := v.(map[string]interface{})
					return hashcode.String(m["id"].(string))
				},
			},
		},
	}
}

func eipProtectionProtectedEipSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the protected EIP.`,
			},
			"public_ipv4": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The IPv4 address of the protected EIP.`,
			},
			"public_ipv6": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The IPv6 address of the protected EIP.`,
			},
		},
	}
}

func buildProtectedEipBodyAndIds(rawParams *schema.Set) ([]map[string]interface{}, []string) {
	if rawParams.Len() < 1 {
		return nil, nil
	}

	requestParams := make([]map[string]interface{}, rawParams.Len())
	protectedEipIds := make([]string, rawParams.Len())
	for i, v := range rawParams.List() {
		raw := v.(map[string]interface{})
		requestParams[i] = map[string]interface{}{
			"id":          raw["id"],
			"public_ip":   utils.ValueIgnoreEmpty(raw["public_ipv4"]),
			"public_ipv6": utils.ValueIgnoreEmpty(raw["public_ipv6"]),
		}
		protectedEipIds[i] = raw["id"].(string)
	}
	return requestParams, protectedEipIds
}

func buildModifyProtectedEipsParams(objectId string, protectedEips *schema.Set,
	operation int) (map[string]interface{}, []string) {
	protectEipInfo, protectedEipIds := buildProtectedEipBodyAndIds(protectedEips)
	requestParams := map[string]interface{}{
		"object_id": objectId,
		"ip_infos":  protectEipInfo,
		"status":    operation,
	}
	return requestParams, protectedEipIds
}

func resourceEipProtectionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)

		objectId      = d.Get("object_id").(string)
		protectedEips = d.Get("protected_eip").(*schema.Set)
	)
	client, err := cfg.NewServiceClient(eipProtectionProduct, region)
	if err != nil {
		return diag.Errorf("error creating CFW Client: %s", err)
	}

	err = protectEips(ctx, client, objectId, protectedEips, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(objectId)

	return resourceEipProtectionRead(ctx, d, meta)
}

func extractSyncedEips(p pagination.Page) ([]interface{}, error) {
	var s struct {
		Records []interface{} `json:"records"`
	}
	err := p.(eipProtectionPage).Result.ExtractIntoStructPtr(&s, "data")
	return s.Records, err
}

// IsEmpty checks whether current page is empty.
func (b eipProtectionPage) IsEmpty() (bool, error) {
	arr, err := extractSyncedEips(b)
	return len(arr) == 0, err
}

func buildSyncedEipsQueryParams(objectId string) string {
	// sync is a required query parameter, which is used to synchronize the public IP from the EIP side to the CFW side.
	return fmt.Sprintf("?offset=0&limit=100&sync=1&object_id=%s", objectId)
}

// QuerySyncedEips is the method used to query synced EIPs for CFW service.
func QuerySyncedEips(client *golangsdk.ServiceClient, url, objectId string) ([]interface{}, error) {
	path := client.Endpoint + url
	path = strings.ReplaceAll(path, "{project_id}", client.ProjectID)

	queryParams := buildSyncedEipsQueryParams(objectId)
	path += queryParams

	var records []interface{}
	err := pagination.NewPager(client, path, func(r pagination.PageResult) pagination.Page {
		p := eipProtectionPage{pagination.OffsetPageBase{PageResult: r}}
		return p
	}).EachPage(func(page pagination.Page) (bool, error) {
		resp, err := extractSyncedEips(page)
		if err != nil {
			return false, err
		}
		records = append(records, resp...)
		return true, nil
	})
	if err != nil {
		return nil, err
	}

	return records, nil
}

func flattenProtectedEips(eips []interface{}) []interface{} {
	rst := make([]interface{}, 0, len(eips))
	for _, eip := range eips {
		if status := utils.PathSearch("status", eip, float64(1)); int(status.(float64)) != openEipProtection {
			continue
		}
		rst = append(rst, map[string]interface{}{
			"id":          utils.PathSearch("id", eip, nil),
			"public_ipv4": utils.PathSearch("public_ip", eip, nil),
			"public_ipv6": utils.PathSearch("public_ipv6", eip, nil),
		})
	}
	return rst
}

// ProtectedEipExist method will return true if a protected public IP exists under the object.
func ProtectedEipExist(eips []interface{}) bool {
	for _, eip := range eips {
		status := utils.PathSearch("status", eip, float64(1))
		if int(status.(float64)) == openEipProtection {
			return true
		}
	}
	return false
}

func resourceEipProtectionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		objectId = d.Id()
	)
	client, err := cfg.NewServiceClient(eipProtectionProduct, region)
	if err != nil {
		return diag.Errorf("error creating CFW Client: %s", err)
	}

	resp, err := getEipProtection(client, queryHttpUrl, objectId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving EIP protection")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("protected_eip", flattenProtectedEips(resp)),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving EIP protection resource fields: %s", err)
	}
	return nil
}

func getEipProtection(client *golangsdk.ServiceClient, queryHttpUrl, objectId string) ([]interface{}, error) {
	resp, err := QuerySyncedEips(client, queryHttpUrl, objectId)
	if err != nil {
		return nil, common.ConvertExpected400ErrInto404Err(err, "error_code", "CFW.00200005")
	}
	if !ProtectedEipExist(resp) {
		return nil, golangsdk.ErrDefault404{}
	}

	return resp, nil
}

func resourceEipProtectionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		objectId = d.Id()
	)
	client, err := cfg.NewServiceClient(eipProtectionProduct, region)
	if err != nil {
		return diag.Errorf("error creating CFW Client: %s", err)
	}

	oldRaw, newRaw := d.GetChange("protected_eip")
	protectedEipSet := newRaw.(*schema.Set).Difference(oldRaw.(*schema.Set))
	unprotectedEipSet := oldRaw.(*schema.Set).Difference(newRaw.(*schema.Set))

	if unprotectedEipSet.Len() > 0 {
		err = unprotectEips(ctx, client, objectId, unprotectedEipSet, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if protectedEipSet.Len() > 0 {
		err = protectEips(ctx, client, objectId, protectedEipSet, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceEipProtectionRead(ctx, d, meta)
}

func resourceEipProtectionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		objectId     = d.Id()
		unprotectSet = d.Get("protected_eip").(*schema.Set)
	)
	client, err := cfg.NewServiceClient(eipProtectionProduct, region)
	if err != nil {
		return diag.Errorf("error creating CFW Client: %s", err)
	}

	_, err = getEipProtection(client, queryHttpUrl, objectId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving EIP protection")
	}

	err = unprotectEips(ctx, client, objectId, unprotectSet, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", "CFW.00200005"),
			"error deleting EIP protection",
		)
	}
	return nil
}

func protectEips(ctx context.Context, client *golangsdk.ServiceClient, objectId string, protectedEipSet *schema.Set,
	timeout time.Duration) error {
	requestParams, protectedEipIds := buildModifyProtectedEipsParams(objectId, protectedEipSet, openEipProtection)

	// Before modifying the EIP protection, synchronize all public IPs to the CFW service.
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      syncedEipsRefreshFunc(client, queryHttpUrl, objectId, protectedEipIds, closeEipProtection),
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("the EIPs to be protected have not been fully synchronized to CFW: %s", err)
	}

	opts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes:          []int{200},
	}

	path := client.Endpoint + modifyProtectHttpUrl
	path = strings.ReplaceAll(path, "{project_id}", client.ProjectID)
	opts.JSONBody = utils.RemoveNil(requestParams)
	_, err = client.Request("POST", path, &opts)
	if err != nil {
		return fmt.Errorf("error enabling EIP protection: %s", err)
	}

	// After modifying the EIP protection, check the status of all protected public IPs.
	stateConf = &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      syncedEipsRefreshFunc(client, queryHttpUrl, objectId, protectedEipIds, openEipProtection),
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("an error occurred while protecting EIPs: %s", err)
	}
	return nil
}

func unprotectEips(ctx context.Context, client *golangsdk.ServiceClient, objectId string, unprotectedEipSet *schema.Set,
	timeout time.Duration) error {
	requestParams, protectedEipIds := buildModifyProtectedEipsParams(objectId, unprotectedEipSet, closeEipProtection)
	opts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes:          []int{200},
	}

	path := client.Endpoint + modifyProtectHttpUrl
	path = strings.ReplaceAll(path, "{project_id}", client.ProjectID)
	opts.JSONBody = utils.RemoveNil(requestParams)
	_, err := client.Request("POST", path, &opts)
	if err != nil {
		return fmt.Errorf("error disabling EIP protection: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      syncedEipsRefreshFunc(client, queryHttpUrl, objectId, protectedEipIds, closeEipProtection),
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("an error occurred while unprotecting the EIP: %s", err)
	}
	return nil
}

func syncedEipsRefreshFunc(client *golangsdk.ServiceClient, url, objectId string, syncEips []string,
	targetStatus int) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := QuerySyncedEips(client, url, objectId)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] No synchronized EIP found")
				return resp, "PENDING", nil
			}
			return nil, "ERROR", err
		}

		syncedEipCount := 0
		for _, eipId := range syncEips {
			for _, val := range resp {
				// The valid values for status are 0 and 1.
				// If status is not found in the response, the default value is 2.
				if eipId == utils.PathSearch("id", val, "").(string) &&
					int(utils.PathSearch("status", val, float64(2)).(float64)) == targetStatus {
					syncedEipCount++
					break
				}
			}
		}
		if len(syncEips) > syncedEipCount {
			return resp, "PENDING", nil
		}
		return resp, "COMPLETED", nil
	}
}

func resourceEipProtectionImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	// object_id is both a resource parameter and an resource ID.
	err := d.Set("object_id", d.Id())
	return []*schema.ResourceData{d}, err
}
