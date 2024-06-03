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
// @API GA POST /v1/ip-groups/{ip_group_id}/add-ips
// @API GA POST /v1/ip-groups/{ip_group_id}/remove-ips
// @API GA POST /v1/ip-groups/{ip_group_id}/associate-listener
// @API GA POST /v1/ip-groups/{ip_group_id}/disassociate-listener
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
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
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
				Elem:     ipAddressGroupSchema(),
			},
			"listeners": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     associatedListenersSchema(),
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
		},
	}
}

func ipAddressGroupSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"cidr": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
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
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
	return &sc
}

func buildCreateIpAddressGroupBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"ip_group": map[string]interface{}{
			"name":        d.Get("name"),
			"description": utils.ValueIgnoreEmpty(d.Get("description")),
		},
	}
	return bodyParams
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

	err = waitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the creation of IP address group (%s) to complete: %s", d.Id(), err)
	}

	//  call addIps to support more than 20 ip addresses
	if val, ok := d.GetOk("ip_addresses"); ok {
		err = addIps(ctx, d, meta, createIpAddressGroupClient, val.(*schema.Set).List())
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if val, ok := d.GetOk("listeners"); ok {
		err = associateListener(ctx, d, meta, createIpAddressGroupClient, val.(*schema.Set).List())
		if err != nil {
			return diag.FromErr(err)
		}
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
	rawArray, _ := resp.([]interface{})
	rst := make([]interface{}, 0, len(rawArray))
	for _, v := range rawArray {
		params := map[string]interface{}{
			"id":   utils.PathSearch("listener_id", v, nil),
			"type": utils.PathSearch("type", v, nil),
		}
		rst = append(rst, params)
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
		d.Set("ip_addresses", flattenGetIpListResponseBody(utils.PathSearch("ip_group.ip_list",
			getIpAddressGroupRespBody, make([]interface{}, 0)))),
		d.Set("status", utils.PathSearch("ip_group.status", getIpAddressGroupRespBody, nil)),
		d.Set("created_at", utils.PathSearch("ip_group.created_at", getIpAddressGroupRespBody, nil)),
		d.Set("updated_at", utils.PathSearch("ip_group.updated_at", getIpAddressGroupRespBody, nil)),
		d.Set("listeners", flattenGetListenersResponseBody(utils.PathSearch("ip_group.associated_listeners",
			getIpAddressGroupRespBody, make([]interface{}, 0)))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildUpdateIpAddressGroupBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"ip_group": map[string]interface{}{
			"name":        utils.ValueIgnoreEmpty(d.Get("name")),
			"description": d.Get("description"),
		},
	}
	return bodyParams
}

func resourceIpAddressGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		updateIpAddressGroupHttpUrl = "v1/ip-groups/{ip_group_id}"
		updateIpAddressGroupProduct = "ga"
	)

	updateIpAddressGroupClient, err := cfg.NewServiceClient(updateIpAddressGroupProduct, region)
	if err != nil {
		return diag.Errorf("error creating IP address group client: %s", err)
	}

	if d.HasChanges("name", "description") {
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

	if d.HasChange("ip_addresses") {
		oldIpsRaw, newIpsRaw := d.GetChange("ip_addresses")
		addIpsRaw := newIpsRaw.(*schema.Set).Difference(oldIpsRaw.(*schema.Set))
		removeIpsRaw := oldIpsRaw.(*schema.Set).Difference(newIpsRaw.(*schema.Set))

		if removeIpsRaw.Len() > 0 {
			err = removeIps(ctx, d, meta, updateIpAddressGroupClient, removeIpsRaw.List())
			if err != nil {
				return diag.FromErr(err)
			}
		}

		if addIpsRaw.Len() > 0 {
			err = addIps(ctx, d, meta, updateIpAddressGroupClient, addIpsRaw.List())
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("listeners") {
		oldListenerRaw, newListenerRaw := d.GetChange("listeners")
		associateListenerRaw := newListenerRaw.(*schema.Set).Difference(oldListenerRaw.(*schema.Set))
		disassociateListenerRaw := oldListenerRaw.(*schema.Set).Difference(newListenerRaw.(*schema.Set))

		if disassociateListenerRaw.Len() > 0 {
			err = disassociateListener(ctx, d, meta, updateIpAddressGroupClient, disassociateListenerRaw.List())
			if err != nil {
				return diag.FromErr(err)
			}
		}

		if associateListenerRaw.Len() > 0 {
			err = associateListener(ctx, d, meta, updateIpAddressGroupClient, associateListenerRaw.List())
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return resourceIpAddressGroupRead(ctx, d, meta)
}

func addIps(ctx context.Context, d *schema.ResourceData, meta interface{}, client *golangsdk.ServiceClient,
	rawParams interface{}) error {
	addIpsHttpUrl := "v1/ip-groups/{ip_group_id}/add-ips"
	addIpsPath := client.Endpoint + addIpsHttpUrl
	addIpsPath = strings.ReplaceAll(addIpsPath, "{ip_group_id}", d.Id())

	addIpsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	if rawArray, ok := rawParams.([]interface{}); ok {
		batchSize := 20
		rst := make([]interface{}, len(rawArray))

		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			rst[i] = map[string]interface{}{
				"cidr":        raw["cidr"],
				"description": utils.ValueIgnoreEmpty(raw["description"]),
			}
		}

		for i := 0; i < len(rst); i += batchSize {
			endIndex := i + batchSize
			if endIndex > len(rst) {
				endIndex = len(rst)
			}

			batch := rst[i:endIndex]
			addIpsOpt.JSONBody = utils.RemoveNil(map[string]interface{}{
				"ip_list": batch,
			})

			_, err := client.Request("POST", addIpsPath, &addIpsOpt)
			if err != nil {
				return fmt.Errorf("error adding IP addresses to the IP address group: %s", err)
			}

			err = waitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutUpdate))
			if err != nil {
				return fmt.Errorf("error waiting for the completion of adding IP addresses to the IP address group: %s", err)
			}
		}
	}

	return nil
}

func removeIps(ctx context.Context, d *schema.ResourceData, meta interface{}, client *golangsdk.ServiceClient,
	rawParams interface{}) error {
	removeIpsHttpUrl := "v1/ip-groups/{ip_group_id}/remove-ips"
	removeIpsPath := client.Endpoint + removeIpsHttpUrl
	removeIpsPath = strings.ReplaceAll(removeIpsPath, "{ip_group_id}", d.Id())

	removeIpsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	if rawArray, ok := rawParams.([]interface{}); ok {
		batchSize := 20
		rst := make([]string, len(rawArray))
		for i, v := range rawArray {
			rst[i] = utils.PathSearch("cidr", v, nil).(string)
		}

		for i := 0; i < len(rst); i += batchSize {
			endIndex := i + batchSize
			if endIndex > len(rst) {
				endIndex = len(rst)
			}

			batch := rst[i:endIndex]
			removeIpsOpt.JSONBody = utils.RemoveNil(map[string]interface{}{
				"ip_list": batch,
			})

			_, err := client.Request("POST", removeIpsPath, &removeIpsOpt)
			if err != nil {
				return fmt.Errorf("error removing IP addresses from the IP address group: %s", err)
			}

			err = waitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutUpdate))
			if err != nil {
				return fmt.Errorf("error waiting for the completion of removing IP addresses from the IP address group: %s", err)
			}
		}
	}

	return nil
}

func associateListener(ctx context.Context, d *schema.ResourceData, meta interface{}, client *golangsdk.ServiceClient,
	rawParams interface{}) error {
	associateListenerHttpUrl := "v1/ip-groups/{ip_group_id}/associate-listener"
	associateListenerPath := client.Endpoint + associateListenerHttpUrl
	associateListenerPath = strings.ReplaceAll(associateListenerPath, "{ip_group_id}", d.Id())

	associateListenerOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	if rawArray, ok := rawParams.([]interface{}); ok {
		for _, v := range rawArray {
			raw := v.(map[string]interface{})
			listenerId := raw["id"]
			associateListenerOpt.JSONBody = map[string]interface{}{
				"listener_id": listenerId,
				"type":        raw["type"],
			}
			_, err := client.Request("POST", associateListenerPath, &associateListenerOpt)
			if err != nil {
				return fmt.Errorf("error associate listener (%s) for IP address group: %s", listenerId, err)
			}

			err = waitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutUpdate))
			if err != nil {
				return fmt.Errorf("error waiting for the completion of association listener (%s) to the IP address group: %s", listenerId, err)
			}
		}
	}

	return nil
}

func disassociateListener(ctx context.Context, d *schema.ResourceData, meta interface{}, client *golangsdk.ServiceClient,
	rawParams interface{}) error {
	disassociateListenerHttpUrl := "v1/ip-groups/{ip_group_id}/disassociate-listener"
	disassociateListenerPath := client.Endpoint + disassociateListenerHttpUrl
	disassociateListenerPath = strings.ReplaceAll(disassociateListenerPath, "{ip_group_id}", d.Id())

	disassociateListenerOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	if rawArray, ok := rawParams.([]interface{}); ok {
		for _, v := range rawArray {
			raw := v.(map[string]interface{})
			listenerId := raw["id"]
			disassociateListenerOpt.JSONBody = map[string]interface{}{
				"listener_id": listenerId,
			}
			_, err := client.Request("POST", disassociateListenerPath, &disassociateListenerOpt)
			if err != nil {
				return fmt.Errorf("error disassociated listener (%s) from IP address group: %s", listenerId, err)
			}

			err = waitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutUpdate))
			if err != nil {
				return fmt.Errorf("error waiting for the completion of disassociation listener (%s) from the IP address group: %s", listenerId, err)
			}
		}
	}

	return nil
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

	if val, ok := d.GetOk("listeners"); ok {
		err = disassociateListener(ctx, d, meta, deleteIpAddressGroupClient, val.(*schema.Set).List())
		if err != nil {
			return diag.FromErr(err)
		}
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

func waitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
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
				// When the error code is 404, the value of respBody is nil, and a non-null value is returned to avoid continuing the loop check.
				return "Resource Not Found", "DELETED", nil
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
