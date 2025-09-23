package elb

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ELB GET /v3/{project_id}/elb/master-slave-pools
func DataSourceActiveStandbyPools() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceActiveStandbyPoolsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"pool_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"connection_drain": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"true", "false",
				}, false),
			},
			"ip_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"lb_algorithm": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"loadbalancer_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"healthmonitor_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"listener_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"member_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"member_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"pools": {
				Type:     schema.TypeList,
				Elem:     poolsActiveStandbyPoolsSchema(),
				Computed: true,
			},
		},
	}
}

func poolsActiveStandbyPoolsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"any_port_enable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lb_algorithm": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"connection_drain_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"connection_drain_timeout": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"listeners": {
				Type:     schema.TypeList,
				Elem:     poolsActiveStandbyPoolListenersSchema(),
				Computed: true,
			},
			"loadbalancers": {
				Type:     schema.TypeList,
				Elem:     poolsActiveStandbyPoolLoadbalancersSchema(),
				Computed: true,
			},
			"members": {
				Type:     schema.TypeList,
				Elem:     poolsActiveStandbyPoolMembersSchema(),
				Computed: true,
			},
			"healthmonitor": {
				Type:     schema.TypeList,
				Elem:     poolsActiveStandbyPoolHealthmonitorSchema(),
				Computed: true,
			},
			"quic_cid_hash_strategy": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     poolsActiveStandbyPoolQuicCidHashStrategySchema(),
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
	return &sc
}

func poolsActiveStandbyPoolListenersSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func poolsActiveStandbyPoolLoadbalancersSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func poolsActiveStandbyPoolMembersSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"protocol_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"role": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"operating_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"member_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"reason": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     poolsActiveStandbyPoolMemberReasonSchema(),
			},
			"status": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     poolsActiveStandbyPoolMemberStatusSchema(),
			},
		},
	}
	return &sc
}

func poolsActiveStandbyPoolMemberReasonSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"reason_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expected_response": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"healthcheck_response": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func poolsActiveStandbyPoolMemberStatusSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"listener_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"operating_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func poolsActiveStandbyPoolHealthmonitorSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"delay": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"max_retries": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"timeout": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expected_codes": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"http_method": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"max_retries_down": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"monitor_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"url_path": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func poolsActiveStandbyPoolQuicCidHashStrategySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"len": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"offset": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceActiveStandbyPoolsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		listActiveStandbyPoolsUrl     = "v3/{project_id}/elb/master-slave-pools"
		listActiveStandbyPoolsProduct = "elb"
	)
	elbClient, err := cfg.NewServiceClient(listActiveStandbyPoolsProduct, region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}
	listActiveStandbyPoolsPath := elbClient.Endpoint + listActiveStandbyPoolsUrl
	listActiveStandbyPoolsPath = strings.ReplaceAll(listActiveStandbyPoolsPath, "{project_id}", elbClient.ProjectID)

	listActiveStandbyPoolsQueryParams := buildListActiveStandbyPoolsQueryParams(d)
	listActiveStandbyPoolsPath += listActiveStandbyPoolsQueryParams

	listActiveStandbyPoolsResp, err := pagination.ListAllItems(
		elbClient,
		"marker",
		listActiveStandbyPoolsPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Pools")
	}

	listActiveStandbyPoolsRespJson, err := json.Marshal(listActiveStandbyPoolsResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var listActiveStandbyPoolBody interface{}
	err = json.Unmarshal(listActiveStandbyPoolsRespJson, &listActiveStandbyPoolBody)
	if err != nil {
		return diag.FromErr(err)
	}
	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("pools", flattenListActiveStandbyPoolsBodyPools(listActiveStandbyPoolBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListActiveStandbyPoolsBodyPools(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("pools", resp, make([]interface{}, 0))
	if curJson == nil {
		return nil
	}
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":                       utils.PathSearch("id", v, nil),
			"name":                     utils.PathSearch("name", v, nil),
			"description":              utils.PathSearch("description", v, nil),
			"protocol":                 utils.PathSearch("protocol", v, nil),
			"type":                     utils.PathSearch("type", v, nil),
			"any_port_enable":          utils.PathSearch("any_port_enable", v, nil),
			"enterprise_project_id":    utils.PathSearch("enterprise_project_id", v, nil),
			"ip_version":               utils.PathSearch("ip_version", v, nil),
			"lb_algorithm":             utils.PathSearch("lb_algorithm", v, nil),
			"vpc_id":                   utils.PathSearch("vpc_id", v, nil),
			"connection_drain_enabled": utils.PathSearch("connection_drain.enable", v, nil),
			"connection_drain_timeout": utils.PathSearch("connection_drain.timeout", v, nil),
			"listeners":                flattenActiveStandbyPoolListeners(v),
			"loadbalancers":            flattenActiveStandbyPoolLoadbalancers(v),
			"members":                  flattenActiveStandbyPoolMember(v),
			"healthmonitor":            flattenActiveStandbyPoolMonitor(v),
			"quic_cid_hash_strategy":   flattenActiveStandbyPoolQuicCidHashStrategy(v),
			"created_at":               utils.PathSearch("created_at", v, nil),
			"updated_at":               utils.PathSearch("updated_at", v, nil),
		})
	}
	return rst
}

func flattenActiveStandbyPoolListeners(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("listeners", resp, make([]interface{}, 0))
	if curJson == nil {
		return nil
	}
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id": utils.PathSearch("id", v, nil),
		})
	}
	return rst
}

func flattenActiveStandbyPoolLoadbalancers(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("loadbalancers", resp, make([]interface{}, 0))
	if curJson == nil {
		return nil
	}
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id": utils.PathSearch("id", v, nil),
		})
	}
	return rst
}

func flattenActiveStandbyPoolMember(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("members", resp, make([]interface{}, 0))
	if curJson == nil {
		return nil
	}
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"address":          utils.PathSearch("address", v, nil),
			"protocol_port":    utils.PathSearch("protocol_port", v, nil),
			"role":             utils.PathSearch("role", v, nil),
			"id":               utils.PathSearch("id", v, nil),
			"name":             utils.PathSearch("name", v, nil),
			"subnet_id":        utils.PathSearch("subnet_cidr_id", v, nil),
			"member_type":      utils.PathSearch("member_type", v, nil),
			"instance_id":      utils.PathSearch("instance_id", v, nil),
			"operating_status": utils.PathSearch("operating_status", v, nil),
			"ip_version":       utils.PathSearch("ip_version", v, nil),
			"reason":           flattenActiveStandbyPoolMemberReason(v),
			"status":           flattenActiveStandbyPoolMemberStatus(v),
		})
	}
	return rst
}

func flattenActiveStandbyPoolMonitor(resp interface{}) []interface{} {
	var rst []interface{}
	curJson := utils.PathSearch("healthmonitor", resp, make(map[string]interface{})).(map[string]interface{})
	if len(curJson) < 1 {
		return nil
	}

	rst = []interface{}{
		map[string]interface{}{
			"name":             utils.PathSearch("name", curJson, nil),
			"monitor_port":     utils.PathSearch("monitor_port", curJson, nil),
			"domain_name":      utils.PathSearch("domain_name", curJson, nil),
			"type":             utils.PathSearch("type", curJson, nil),
			"url_path":         utils.PathSearch("url_path", curJson, nil),
			"max_retries_down": utils.PathSearch("max_retries_down", curJson, nil),
			"max_retries":      utils.PathSearch("max_retries", curJson, nil),
			"expected_codes":   utils.PathSearch("expected_codes", curJson, nil),
			"http_method":      utils.PathSearch("http_method", curJson, nil),
			"timeout":          utils.PathSearch("timeout", curJson, nil),
			"delay":            utils.PathSearch("delay", curJson, nil),
			"id":               utils.PathSearch("id", curJson, nil),
		},
	}
	return rst
}

func flattenActiveStandbyPoolQuicCidHashStrategy(resp interface{}) []interface{} {
	curJson := utils.PathSearch("quic_cid_hash_strategy", resp, nil)
	if curJson == nil {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"len":    utils.PathSearch("len", curJson, nil),
			"offset": utils.PathSearch("offset", curJson, nil),
		},
	}
	return rst
}

func buildListActiveStandbyPoolsQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("pool_id"); ok {
		res = fmt.Sprintf("%s&id=%v", res, v)
	}
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	if v, ok := d.GetOk("description"); ok {
		res = fmt.Sprintf("%s&description=%v", res, v)
	}
	if v, ok := d.GetOk("connection_drain"); ok {
		connectionDrain, _ := strconv.ParseBool(v.(string))
		res = fmt.Sprintf("%s&connection_drain=%v", res, connectionDrain)
	}
	if v, ok := d.GetOk("ip_version"); ok {
		res = fmt.Sprintf("%s&ip_version=%v", res, v)
	}
	if v, ok := d.GetOk("lb_algorithm"); ok {
		res = fmt.Sprintf("%s&lb_algorithm=%v", res, v)
	}
	if v, ok := d.GetOk("loadbalancer_id"); ok {
		res = fmt.Sprintf("%s&loadbalancer_id=%v", res, v)
	}
	if v, ok := d.GetOk("healthmonitor_id"); ok {
		res = fmt.Sprintf("%s&healthmonitor_id=%v", res, v)
	}
	if v, ok := d.GetOk("protocol"); ok {
		res = fmt.Sprintf("%s&protocol=%v", res, v)
	}
	if v, ok := d.GetOk("listener_id"); ok {
		res = fmt.Sprintf("%s&listener_id=%v", res, v)
	}
	if v, ok := d.GetOk("type"); ok {
		res = fmt.Sprintf("%s&type=%v", res, v)
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		res = fmt.Sprintf("%s&vpc_id=%v", res, v)
	}
	if v, ok := d.GetOk("member_address"); ok {
		res = fmt.Sprintf("%s&member_address=%v", res, v)
	}
	if v, ok := d.GetOk("member_instance_id"); ok {
		res = fmt.Sprintf("%s&member_instance_id=%v", res, v)
	}
	if res != "" {
		res = "?" + res[1:]
	}
	return res
}
