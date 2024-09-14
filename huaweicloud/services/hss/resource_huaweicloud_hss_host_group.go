package hss

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"

	hssv5 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/hss/v5"
	hssv5model "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/hss/v5/model"

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

func checkAllHostsAvailable(ctx context.Context, client *hssv5.HssClient, epsId string, hostIDs []string,
	timeout time.Duration) ([]string, error) {
	unprotected := make([]string, 0)
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
			unprotected = append(unprotected, unprotectedHostId.(string))
		}
	}
	return unprotected, nil
}

func hostStatusRefreshFunc(client *hssv5.HssClient, epsId, hostId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var unprotectedHostId string
		if epsId == "" {
			epsId = "all_granted_eps"
		}

		request := hssv5model.ListHostStatusRequest{
			EnterpriseProjectId: utils.String(epsId),
			Refresh:             utils.Bool(true),
			HostId:              utils.String(hostId),
		}
		resp, err := client.ListHostStatus(&request)
		if err != nil {
			return unprotectedHostId, "ERROR", err
		}
		if resp == nil || len(*resp.DataList) < 1 {
			return unprotectedHostId, "PENDING", nil
		}

		hostList := *resp.DataList
		if *hostList[0].ProtectStatus == string(ProtectStatusClosed) {
			unprotectedHostId = *hostList[0].HostId
		}
		return unprotectedHostId, "COMPLETED", nil
	}
}

func resourceHostGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.HcHssV5Client(region)
	if err != nil {
		return diag.Errorf("error creating HSS v5 client: %s", err)
	}

	var (
		groupName = d.Get("name").(string)
		hostIds   = utils.ExpandToStringListBySet(d.Get("host_ids").(*schema.Set))
		epsId     = cfg.GetEnterpriseProjectID(d)

		request = hssv5model.AddHostsGroupRequest{
			Region:              region,
			EnterpriseProjectId: utils.StringIgnoreEmpty(epsId),
			Body: &hssv5model.AddHostsGroupRequestInfo{
				GroupName:  groupName,
				HostIdList: hostIds,
			},
		}
	)

	unprotected, err := checkAllHostsAvailable(ctx, client, epsId, hostIds, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] All hosts are availabile.")
	if len(unprotected) > 1 {
		log.Printf("[WARN] These hosts are not protected: %#v", unprotected)
		d.Set("unprotect_host_ids", unprotected)
	}
	_, err = client.AddHostsGroup(&request)
	if err != nil {
		return diag.Errorf("error creating host group: %s", err)
	}

	allHostGroups, err := queryHostGroups(client, region, epsId, groupName)
	if err != nil {
		return diag.FromErr(err)
	}
	if len(allHostGroups) < 1 {
		return common.CheckDeletedDiag(d, err, "host group")
	}
	d.SetId(*allHostGroups[0].GroupId)

	return resourceHostGroupRead(ctx, d, meta)
}

func queryHostGroups(client *hssv5.HssClient, region, epsId, name string) ([]hssv5model.HostGroupItem, error) {
	var (
		offset        int32
		limit         int32
		allHostGroups []hssv5model.HostGroupItem = make([]hssv5model.HostGroupItem, 0)
	)
	for {
		response, err := client.ListHostGroups(&hssv5model.ListHostGroupsRequest{
			Region:              region,
			EnterpriseProjectId: utils.StringIgnoreEmpty(epsId),
			GroupName:           utils.StringIgnoreEmpty(name),
			Offset:              utils.Int32IgnoreEmpty(offset),
			Limit:               utils.Int32IgnoreEmpty(limit),
		})
		if err != nil {
			return nil, fmt.Errorf("error fetching host group: %s", err)
		}

		if response != nil && response.DataList != nil {
			allHostGroups = append(allHostGroups, *response.DataList...)
		}

		if response == nil || offset >= *response.TotalNum || len(*response.DataList) == 0 {
			break
		}

		offset += *response.TotalNum
	}

	return allHostGroups, nil
}

func QueryHostGroupById(client *hssv5.HssClient, region, epsId, groupId string) (*hssv5model.HostGroupItem, error) {
	allHostGroups, err := queryHostGroups(client, region, epsId, "")
	if err != nil {
		return nil, err
	}
	filter := map[string]interface{}{
		"GroupId": groupId,
	}
	result, err := utils.FilterSliceWithField(allHostGroups, filter)
	if err != nil {
		return nil, fmt.Errorf("erroring filting security groups list: %s", err)
	}

	if len(result) < 1 {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Body: []byte(fmt.Sprintf("the host group (%s) does not exist", groupId)),
			},
		}
	}
	if item, ok := result[0].(hssv5model.HostGroupItem); ok {
		return &item, nil
	}
	return nil, fmt.Errorf("invalid host group list, want 'hssv5model.HostGroupItem', but '%T'", result[0])
}

func resourceHostGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		groupId = d.Id()
		epsId   = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.HcHssV5Client(region)
	if err != nil {
		return diag.Errorf("error creating HSS v5 client: %s", err)
	}

	resp, err := QueryHostGroupById(client, region, epsId, groupId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "host group")
	}
	log.Printf("[DEBUG] The response of host group is: %#v", resp)

	mErr := multierror.Append(nil,
		d.Set("region", cfg.GetRegion(d)),
		d.Set("name", resp.GroupName),
		d.Set("host_ids", resp.HostIdList),
		d.Set("host_num", resp.HostNum),
		d.Set("risk_host_num", resp.RiskHostNum),
		d.Set("unprotect_host_num", resp.UnprotectHostNum),
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

func resourceHostGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.HcHssV5Client(region)
	if err != nil {
		return diag.Errorf("error creating HSS v5 client: %s", err)
	}

	var (
		groupId   = d.Id()
		groupName = d.Get("name").(string)
		hostIds   = utils.ExpandToStringListBySet(d.Get("host_ids").(*schema.Set))
		epsId     = cfg.GetEnterpriseProjectID(d)

		request = hssv5model.ChangeHostsGroupRequest{
			Region:              region,
			EnterpriseProjectId: utils.StringIgnoreEmpty(epsId),
			Body: &hssv5model.ChangeHostsGroupRequestInfo{
				GroupId:    groupId,
				GroupName:  utils.StringIgnoreEmpty(groupName),
				HostIdList: &hostIds,
			},
		}
	)

	unprotected, err := checkAllHostsAvailable(ctx, client, epsId, hostIds, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] All hosts are availabile.")
	if len(unprotected) > 1 {
		log.Printf("[WARN] These hosts are not protected: %#v", unprotected)
		d.Set("unprotect_host_ids", unprotected)
	}
	_, err = client.ChangeHostsGroup(&request)
	if err != nil {
		return diag.Errorf("error updating host group: %s", err)
	}

	return resourceHostGroupRead(ctx, d, meta)
}

func resourceHostGroupDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.HcHssV5Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating HSS v5 client: %s", err)
	}

	var (
		groupId = d.Id()

		request = hssv5model.DeleteHostsGroupRequest{
			Region:              cfg.GetRegion(d),
			EnterpriseProjectId: utils.StringIgnoreEmpty(cfg.GetEnterpriseProjectID(d)),
			GroupId:             groupId,
		}
	)

	_, err = client.DeleteHostsGroup(&request)
	if err != nil {
		return common.CheckDeletedDiag(d, parseDeleteHostGroupResponseError(err), "error deleting host group")
	}

	return nil
}

// When the host group does not exist, the response code for deleting the API is `400`,
// and the response body is as follows:
// {"status_code":400,"request_id":"f17e56c2e92584cfd4614ab467cd6a1b","error_code":"",
// "error_message":"{\"error_code\":\"00100090\",\"error_description\":\"Failed to load server groups.\"}",
// "encoded_authorization_message":""}
func parseDeleteHostGroupResponseError(err error) error {
	var errObj map[string]interface{}
	if jsonErr := json.Unmarshal([]byte(err.Error()), &errObj); jsonErr != nil {
		log.Printf("[WARN] failed to unmarshal error object: %s", jsonErr)
		return err
	}

	statusCode, parseStatusCodeErr := jmespath.Search("status_code", errObj)
	if parseStatusCodeErr != nil || statusCode == nil {
		log.Printf("[WARN] failed to parse status_code from response body: %s", parseStatusCodeErr)
		return err
	}

	if statusCodeFloat, ok := statusCode.(float64); ok && int(statusCodeFloat) == 400 {
		errorMessage, parseErrorMessageErr := jmespath.Search("error_message", errObj)
		if parseErrorMessageErr != nil || errorMessage == nil {
			log.Printf("[WARN] failed to parse error_message: %s", parseErrorMessageErr)
			return err
		}

		var errMsgObj map[string]interface{}
		if errMsgJson := json.Unmarshal([]byte(errorMessage.(string)), &errMsgObj); errMsgJson != nil {
			log.Printf("[WARN] failed to unmarshal error_message: %s", errMsgJson)
			return err
		}

		errorCode, errorCodeErr := jmespath.Search("error_code", errMsgObj)
		if errorCodeErr != nil || errorCode == nil {
			log.Printf("[WARN] failed to extract error_code: %s", errorCodeErr)
			return err
		}

		if errorCode == "00100090" {
			return golangsdk.ErrDefault404{}
		}
	}

	return err
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
