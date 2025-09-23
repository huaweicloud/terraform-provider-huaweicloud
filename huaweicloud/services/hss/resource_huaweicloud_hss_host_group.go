package hss

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

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

type ProtectStatus string

const (
	ProtectStatusClosed ProtectStatus = "closed"
	ProtectStatusOpened ProtectStatus = "opened"
)

// @API HSS DELETE /v5/{project_id}/host-management/groups
// @API HSS GET /v5/{project_id}/host-management/groups
// @API HSS POST /v5/{project_id}/host-management/groups
// @API HSS PUT /v5/{project_id}/host-management/groups
// @API HSS GET /v5/{project_id}/host-management/hosts
func ResourceHostGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceHostGroupCreate,
		ReadContext:   resourceHostGroupRead,
		UpdateContext: resourceHostGroupUpdate,
		DeleteContext: resourceHostGroupDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: &schema.ResourceImporter{
			StateContext: resourceHostGroupImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the host group is located.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the host group.",
			},
			"host_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "schema: Required; The list of host IDs.",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The ID of the enterprise project to which the host group belongs.",
			},
			// Attributes
			"host_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The host number.",
			},
			"risk_host_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The current status of the virtual interface.",
			},
			"unprotect_host_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The creation time of the virtual interface.",
			},
			"unprotect_host_ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The ID list of the unprotect hosts.",
			},
		},
	}
}

func checkAllHostsAvailable(ctx context.Context, client *golangsdk.ServiceClient, epsId string, hostIDs []string,
	timeout time.Duration) ([]string, error) {
	unprotectedIDs := make([]string, 0)
	for _, hostId := range hostIDs {
		log.Printf("[DEBUG] Waiting for the host (%s) status to become available.", hostId)
		stateConf := &resource.StateChangeConf{
			Pending:      []string{"PENDING"},
			Target:       []string{"COMPLETED"},
			Refresh:      hostStatusRefreshFunc(client, epsId, hostId),
			Timeout:      timeout,
			Delay:        30 * time.Second,
			PollInterval: 30 * time.Second,
		}
		unprotectedHostId, err := stateConf.WaitForStateContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("error waiting for the host (%s) status to become completed: %s", hostId, err)
		}

		if unprotectedHostId != nil && unprotectedHostId.(string) != "" {
			unprotectedIDs = append(unprotectedIDs, unprotectedHostId.(string))
		}
	}

	return unprotectedIDs, nil
}

func getHostFunc(client *golangsdk.ServiceClient, epsId, hostId string) (interface{}, error) {
	getPath := client.Endpoint + "v5/{project_id}/host-management/hosts"
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += fmt.Sprintf("?enterprise_project_id=%v&refresh=%v&host_id=%v", epsId, true, hostId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving HSS host: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	hostResp := utils.PathSearch("data_list[0]", getRespBody, nil)

	return hostResp, nil
}

func hostStatusRefreshFunc(client *golangsdk.ServiceClient, epsId, hostId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var unprotectedHostId string
		if epsId == "" {
			epsId = QueryAllEpsValue
		}

		hostResp, err := getHostFunc(client, epsId, hostId)
		if err != nil {
			return "", "ERROR", err
		}

		protectStatus := utils.PathSearch("protect_status", hostResp, "").(string)
		if hostResp == nil || protectStatus == "" {
			return "", "PENDING", nil
		}

		if protectStatus == string(ProtectStatusClosed) {
			unprotectedHostId = hostId
		}

		return unprotectedHostId, "COMPLETED", nil
	}
}

func buildCreateOrUpdateHostGroupQueryParams(epsId string) string {
	if epsId != "" {
		return fmt.Sprintf("?enterprise_project_id=%v", epsId)
	}

	return ""
}

func buildCreateHostGroupBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"group_name":   d.Get("name"),
		"host_id_list": utils.ExpandToStringListBySet(d.Get("host_ids").(*schema.Set)),
	}
	return bodyParams
}

func resourceHostGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		product   = "hss"
		epsId     = cfg.GetEnterpriseProjectID(d)
		groupName = d.Get("name").(string)
		hostIds   = utils.ExpandToStringListBySet(d.Get("host_ids").(*schema.Set))
		mErr      *multierror.Error
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	// Before creating, check if all hosts can be accessed and obtain a list of all host IDs that have not enabled
	// host protection.
	unprotectedIDs, err := checkAllHostsAvailable(ctx, client, epsId, hostIds, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	if len(unprotectedIDs) > 0 {
		mErr = multierror.Append(nil, d.Set("unprotect_host_ids", unprotectedIDs))
	}
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving `unprotect_host_ids` field in creation operation: %s", err)
	}

	createPath := client.Endpoint + "v5/{project_id}/host-management/groups"
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath += buildCreateOrUpdateHostGroupQueryParams(epsId)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"region": region},
		JSONBody:         buildCreateHostGroupBodyParams(d),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating HSS host group: %s", err)
	}

	hostGroups, err := queryHostGroupsByName(client, region, epsId, groupName)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(hostGroups) < 1 {
		return diag.Errorf("error creating HSS host group: after successful creation, host group is not found " +
			"in query API response")
	}

	groupId := utils.PathSearch("group_id", hostGroups[0], "").(string)
	if groupId == "" {
		return diag.Errorf("error creating HSS host group: ID is not found in API response")
	}

	d.SetId(groupId)

	return resourceHostGroupRead(ctx, d, meta)
}

func buildQueryHostGroupsByNameQueryParams(epsId, groupName string) string {
	queryParams := "?limit=20"
	if epsId != "" {
		queryParams += fmt.Sprintf("&enterprise_project_id=%v", epsId)
	}
	if groupName != "" {
		queryParams += fmt.Sprintf("&group_name=%v", groupName)
	}

	return queryParams
}

func queryHostGroupsByName(client *golangsdk.ServiceClient, region, epsId, groupName string) ([]interface{}, error) {
	getPath := client.Endpoint + "v5/{project_id}/host-management/groups"
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += buildQueryHostGroupsByNameQueryParams(epsId, groupName)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"region": region},
	}

	var (
		offset = 0
		result = make([]interface{}, 0)
	)

	// The `name` parameter is a fuzzy match, so pagination must be used to retrieve all data related to that name.
	for {
		currentPath := fmt.Sprintf("%s&offset=%v", getPath, offset)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return nil, fmt.Errorf("error retrieving HSS host groups, %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return nil, err
		}

		hostGroupsResp := utils.PathSearch("data_list", getRespBody, make([]interface{}, 0)).([]interface{})
		if len(hostGroupsResp) == 0 {
			break
		}

		result = append(result, hostGroupsResp...)
		offset += len(hostGroupsResp)
	}

	return result, nil
}

func filterHostGroupById(allHostGroups []interface{}, groupId string) interface{} {
	for _, hostGroup := range allHostGroups {
		if utils.PathSearch("group_id", hostGroup, "").(string) == groupId {
			return hostGroup
		}
	}

	return nil
}

func QueryHostGroupById(client *golangsdk.ServiceClient, region, epsId, groupId string) (interface{}, error) {
	allHostGroups, err := queryHostGroupsByName(client, region, epsId, "")
	if err != nil {
		return nil, err
	}

	hostGroup := filterHostGroupById(allHostGroups, groupId)
	if hostGroup == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return hostGroup, nil
}

func resourceHostGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		epsId   = cfg.GetEnterpriseProjectID(d)
		groupId = d.Id()
		product = "hss"
		mErr    *multierror.Error
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	hostGroup, err := QueryHostGroupById(client, region, epsId, groupId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "HSS host group")
	}

	mErr = multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("group_name", hostGroup, nil)),
		d.Set("host_ids", utils.ExpandToStringList(utils.PathSearch("host_id_list", hostGroup, make([]interface{}, 0)).([]interface{}))),
		d.Set("host_num", utils.PathSearch("host_num", hostGroup, nil)),
		d.Set("risk_host_num", utils.PathSearch("risk_host_num", hostGroup, nil)),
		d.Set("unprotect_host_num", utils.PathSearch("unprotect_host_num", hostGroup, nil)),
	)

	if len(d.Get("unprotect_host_ids").([]interface{})) == 0 {
		// The reason for writing an empty array to `unprotect_host_ids` is to avoid unexpected changes
		mErr = multierror.Append(mErr, d.Set("unprotect_host_ids", make([]string, 0)))
	}

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving host group fields: %s", err)
	}

	return nil
}

func buildUpdateHostGroupBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"group_name":   d.Get("name"),
		"group_id":     d.Id(),
		"host_id_list": utils.ExpandToStringListBySet(d.Get("host_ids").(*schema.Set)),
	}

	return bodyParams
}

func resourceHostGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		epsId   = cfg.GetEnterpriseProjectID(d)
		product = "hss"
		hostIds = utils.ExpandToStringListBySet(d.Get("host_ids").(*schema.Set))
		mErr    *multierror.Error
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	unprotectedIDs, err := checkAllHostsAvailable(ctx, client, epsId, hostIds, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return diag.FromErr(err)
	}

	if len(unprotectedIDs) > 0 {
		mErr = multierror.Append(nil, d.Set("unprotect_host_ids", unprotectedIDs))
	}
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving `unprotect_host_ids` field in update operation: %s", err)
	}

	updatePath := client.Endpoint + "v5/{project_id}/host-management/groups"
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath += buildCreateOrUpdateHostGroupQueryParams(epsId)
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"region": region},
		JSONBody:         buildUpdateHostGroupBodyParams(d),
	}

	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating HSS host group: %s", err)
	}

	return resourceHostGroupRead(ctx, d, meta)
}

func buildDeleteHostGroupQueryParams(epsId, groupId string) string {
	queryParams := fmt.Sprintf("?group_id=%v", groupId)
	if epsId != "" {
		queryParams += fmt.Sprintf("&enterprise_project_id=%v", epsId)
	}

	return queryParams
}

func resourceHostGroupDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		epsId   = cfg.GetEnterpriseProjectID(d)
		groupId = d.Id()
		product = "hss"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	deletePath := client.Endpoint + "v5/{project_id}/host-management/groups"
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath += buildDeleteHostGroupQueryParams(epsId, groupId)
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"region": region},
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	// The API error when deleting non-existent resource is as follows:
	// {
	// "error_code": "HSS.1019",
	// "error_description": "查询服务器组信息失败",
	// "error_msg": "查询服务器组信息失败"
	// }
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", "HSS.1019"),
			"error deleting HSS host group")
	}

	return nil
}

func resourceHostGroupImportState(_ context.Context, d *schema.ResourceData, _ interface{}) (
	[]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format of import ID, must be <enterprise_project_id>/<id>")
	}

	d.SetId(parts[1])

	return []*schema.ResourceData{d}, d.Set("enterprise_project_id", parts[0])
}
