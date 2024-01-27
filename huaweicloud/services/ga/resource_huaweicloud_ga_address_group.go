package ga

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GA POST /v1/ip-groups
// @API GA GET /v1/ip-groups/{ip_group_id}
// @API GA PUT /v1/ip-groups/{ip_group_id}
// @API GA DELETE /v1/ip-groups/{ip_group_id}
func ResourceIpAddressGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIpAddressGroupCreate,
		ReadContext:   resourceIpAddressGroupRead,
		UpdateContext: resourceIpAddressGroupUpdate,
		DeleteContext: resourceIpAddressGroupDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip_addresses": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem:     ipAddressGroupSchema(),
				MaxItems: 20,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"listeners": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     associatedListenersSchema(),
			},
		},
	}
}

func ipAddressGroupSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"cidr": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func associatedListenersSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func buildCreateIpAddressGroupBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"ip_group": map[string]interface{}{
			"name":        d.Get("name"),
			"description": utils.ValueIngoreEmpty(d.Get("description")),
			"ip_list":     buildIpAddressOptionBodyParams(d.Get("ip_addresses").(*schema.Set)),
		},
	}
	return bodyParams
}

func buildIpAddressOptionBodyParams(rawParams *schema.Set) []map[string]interface{} {
	if rawParams.Len() == 0 {
		return nil
	}
	rst := make([]map[string]interface{}, rawParams.Len())
	for _, val := range rawParams.List() {
		raw := val.(map[string]interface{})
		params := map[string]interface{}{
			"cidr":        raw["cidr"],
			"description": utils.ValueIngoreEmpty(raw["description"]),
		}
		rst = append(rst, params)
	}

	return rst
}

func resourceIpAddressGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		createIpAddressGroupHttpUrl = "v1/ip-groups"
		createEndpointProduct       = "ga"
	)
	createIpAddressGroupClient, err := cfg.NewServiceClient(createEndpointProduct, region)
	if err != nil {
		return diag.Errorf("error creating IP address group client: %s", err)
	}

	createIpAddressGroupPath := createIpAddressGroupClient.Endpoint + createIpAddressGroupHttpUrl
	createIpAddressGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			201,
		},
	}

	createIpAddressGroupOpt.JSONBody = utils.RemoveNil(buildCreateIpAddressGroupBodyParams(d))
	createIpAddressGroupResp, err := createIpAddressGroupClient.Request("POST", createIpAddressGroupPath, &createIpAddressGroupOpt)
	if err != nil {
		return diag.Errorf("error creating IP address group: %s", err)
	}

	createIpAddressGroupRespBody, err := utils.FlattenResponse(createIpAddressGroupResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := jmespath.Search("ip_group.id", createIpAddressGroupRespBody)
	if err != nil {
		return diag.Errorf("error creating IP address group: ID is not found in API response")
	}

	d.SetId(id.(string))

	err = createIpAddressGroupWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the creation of IP address group (%s) to complete: %s", d.Id(), err)
	}

	return resourceIpAddressGroupRead(ctx, d, meta)
}

func getIpAddressGroupInfo(d *schema.ResourceData, meta interface{}) (*http.Response, error) {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		getIpAddressGroupHttpUrl = "v1/ip-groups/{ip_group_id}"
		getIpAddressGroupProduct = "ga"
	)

	getIpAddressGroupClient, err := cfg.NewServiceClient(getIpAddressGroupProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating IP address group client: %s", err)
	}

	getIpAddressGroupPath := getIpAddressGroupClient.Endpoint + getIpAddressGroupHttpUrl
	getIpAddressGroupPath = strings.ReplaceAll(getIpAddressGroupPath, "{ip_group_id}", d.Id())
	getIpAddressGroupOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	resp, err := getIpAddressGroupClient.Request("GET", getIpAddressGroupPath, &getIpAddressGroupOpts)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func flattenGetIpListResponseBody(rawParams interface{}) []interface{} {
	rawArray, _ := rawParams.([]interface{})
	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(rawArray))
	for _, v := range rawArray {
		params := map[string]interface{}{
			"cidr":        utils.PathSearch("cidr", v, nil),
			"description": utils.PathSearch("description", v, nil),
			"created_at":  utils.PathSearch("created_at", v, nil),
		}
		rst = append(rst, params)
	}

	return rst
}

func flattenGetListenersResponseBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curJson := utils.PathSearch("ip_group.associated_listeners", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":   utils.PathSearch("listener_id", v, nil),
			"type": utils.PathSearch("type", v, nil),
		})
	}

	return rst
}

func resourceIpAddressGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	getIpAddressGroupResp, err := getIpAddressGroupInfo(d, meta)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "GA IP address group")
	}

	getIpAddressGroupRespBody, err := utils.FlattenResponse(getIpAddressGroupResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(
		nil,
		d.Set("name", utils.PathSearch("ip_group.name", getIpAddressGroupRespBody, nil)),
		d.Set("description", utils.PathSearch("ip_group.description", getIpAddressGroupRespBody, nil)),
		d.Set("ip_addresses", flattenGetIpListResponseBody(utils.PathSearch("ip_group.ip_list", getIpAddressGroupRespBody, make([]interface{}, 0)))),
		d.Set("status", utils.PathSearch("ip_group.status", getIpAddressGroupRespBody, nil)),
		d.Set("created_at", utils.PathSearch("ip_group.created_at", getIpAddressGroupRespBody, nil)),
		d.Set("updated_at", utils.PathSearch("ip_group.updated_at", getIpAddressGroupRespBody, nil)),
		d.Set("listeners", flattenGetListenersResponseBody(getIpAddressGroupRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildUpdateIpAddressGroupBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"ip_group": map[string]interface{}{
			"name":        utils.ValueIngoreEmpty(d.Get("name")),
			"description": d.Get("description"),
		},
	}
	return bodyParams
}

func resourceIpAddressGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	if d.HasChanges("name", "description") {
		var (
			updateIpAddressGroupHttpUrl = "v1/ip-groups/{ip_group_id}"
			updateIpAddressGroupProduct = "ga"
		)

		updateIpAddressGroupClient, err := cfg.NewServiceClient(updateIpAddressGroupProduct, region)
		if err != nil {
			return diag.Errorf("error creating IP address group client: %s", err)
		}

		updateIpAddressGroupPath := updateIpAddressGroupClient.Endpoint + updateIpAddressGroupHttpUrl
		updateIpAddressGroupPath = strings.ReplaceAll(updateIpAddressGroupPath, "{ip_group_id}", d.Id())
		updateIpAddressGroupOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}

		updateIpAddressGroupOpt.JSONBody = utils.RemoveNil(buildUpdateIpAddressGroupBodyParams(d))
		_, err = updateIpAddressGroupClient.Request("PUT", updateIpAddressGroupPath, &updateIpAddressGroupOpt)
		if err != nil {
			return diag.Errorf("error updating IP address group: %s", err)
		}
	}

	return resourceIpAddressGroupRead(ctx, d, meta)
}

func resourceIpAddressGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		deleteIpAddressGroupHttpUrl = "v1/ip-groups/{ip_group_id}"
		deleteIpAddressGroupProduct = "ga"
	)

	deleteIpAddressGroupClient, err := cfg.NewServiceClient(deleteIpAddressGroupProduct, region)
	if err != nil {
		return diag.Errorf("error creating IP address group client: %s", err)
	}

	deleteIpAddressGroupPath := deleteIpAddressGroupClient.Endpoint + deleteIpAddressGroupHttpUrl
	deleteIpAddressGroupPath = strings.ReplaceAll(deleteIpAddressGroupPath, "{ip_group_id}", d.Id())
	deleteIpAddressGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}

	_, err = deleteIpAddressGroupClient.Request("DELETE", deleteIpAddressGroupPath, &deleteIpAddressGroupOpt)
	if err != nil {
		return diag.Errorf("error deleting IP address group: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"DELETED"},
		Refresh:      waitIpAddressGroupStatusRefreshFunc(d, meta, true),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func createIpAddressGroupWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      waitIpAddressGroupStatusRefreshFunc(d, meta, false),
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func waitIpAddressGroupStatusRefreshFunc(d *schema.ResourceData, meta interface{}, isDelete bool) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := getIpAddressGroupInfo(d, meta)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok && isDelete {
				return resp, "DELETED", nil
			}

			return nil, "ERROR", err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, "ERROR", err
		}

		statusRaw, err := jmespath.Search("ip_group.status", respBody)
		if err != nil {
			return nil, "ERROR", fmt.Errorf("error parsing %s from response body", statusRaw)
		}

		status := fmt.Sprintf("%v", statusRaw)

		if utils.StrSliceContains([]string{"ERROR"}, status) {
			return respBody, "ERROR", fmt.Errorf("unexpected address group status: %s", status)
		}

		if utils.StrSliceContains([]string{"ACTIVE"}, status) {
			return respBody, "COMPLETED", nil
		}

		return respBody, "PENDING", nil
	}
}
