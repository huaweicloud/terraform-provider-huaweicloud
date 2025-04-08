// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product ELB
// ---------------------------------------------------------------

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

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ELB GET /v3/{project_id}/elb/pools
func DataSourcePools() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourcePoolsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"pool_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the ID of the ELB pool.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the name of the ELB pool.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the description of the ELB pool.`,
			},
			"any_port_enable": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"true", "false",
				}, false),
			},
			"connection_drain": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"true", "false",
				}, false),
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"member_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"member_device_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"member_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"member_deletion_protection_enable": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"true", "false",
				}, false),
			},
			"pool_health": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"public_border_group": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"quic_cid_len": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"quic_cid_offset": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"loadbalancer_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the loadbalancer ID of the ELB pool.`,
			},
			"healthmonitor_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the health monitor ID of the ELB pool.`,
			},
			"protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the protocol of the ELB pool.`,
			},
			"lb_method": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the method of the ELB pool.`,
			},
			"listener_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the listener ID of the ELB pool.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the type of the backend server group.`,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the ID of the VPC where the backend server group works.`,
			},
			"protection_status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the protection status for update.`,
			},
			"pools": {
				Type:        schema.TypeList,
				Elem:        poolsPoolsSchema(),
				Computed:    true,
				Description: `Pool list. For details, see Data structure of the Pool field.`,
			},
		},
	}
}

func poolsPoolsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The pool ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The pool name.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description of pool.`,
			},
			"protocol": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The protocol of pool.`,
			},
			"ip_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The IP version of pool.`,
			},
			"lb_method": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The method of the ELB pool.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the backend server group.`,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the VPC where the backend server group works.`,
			},
			"any_port_enable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"member_deletion_protection_enable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"public_border_group": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"protection_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The protection status for update.`,
			},
			"protection_reason": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The reason for update protection.`,
			},
			"slow_start_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether to enable slow start.`,
			},
			"slow_start_duration": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The slow start duration, in seconds.`,
			},
			"connection_drain_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether to enable delayed logout.`,
			},
			"connection_drain_timeout": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The timeout of the delayed logout in seconds.`,
			},
			"minimum_healthy_member_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The timeout of the delayed logout in seconds.`,
			},
			"healthmonitor_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The health monitor ID of the ELB pool.`,
			},
			"listeners": {
				Type:        schema.TypeList,
				Elem:        poolsPoolListenersSchema(),
				Computed:    true,
				Description: `Listener list. For details, see Data structure of the listener field.`,
			},
			"loadbalancers": {
				Type:        schema.TypeList,
				Elem:        poolsPoolLoadbalancersSchema(),
				Computed:    true,
				Description: `Loadbalancer list. For details, see Data structure of the loadbalancer field.`,
			},
			"members": {
				Type:        schema.TypeList,
				Elem:        poolsPoolMembersSchema(),
				Computed:    true,
				Description: `Loadbalancer list. For details, see Data structure of the members field.`,
			},
			"persistence": {
				Type:        schema.TypeList,
				Elem:        poolsPoolPersistenceSchema(),
				Computed:    true,
				Description: `Whether connections in the same session will be processed by the same pool member or not.`,
			},
			"quic_cid_hash_strategy": {
				Type:     schema.TypeList,
				Elem:     poolsPoolQuicCidHashStrategySchema(),
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
	return &sc
}

func poolsPoolListenersSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The listener ID.`,
			},
		},
	}
	return &sc
}

func poolsPoolLoadbalancersSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The loadbalancer ID.`,
			},
		},
	}
	return &sc
}

func poolsPoolMembersSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The member ID.`,
			},
		},
	}
	return &sc
}

func poolsPoolPersistenceSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of persistence mode.`,
			},
			"cookie_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the cookie if persistence mode is set appropriately.`,
			},
			"timeout": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The stickiness duration, in minutes.`,
			},
		},
	}
	return &sc
}

func poolsPoolQuicCidHashStrategySchema() *schema.Resource {
	sc := schema.Resource{
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
	return &sc
}

func resourcePoolsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// listPools: Query the List of ELB pools
	var (
		listPoolsHttpUrl = "v3/{project_id}/elb/pools"
		listPoolsProduct = "elb"
	)
	listPoolsClient, err := cfg.NewServiceClient(listPoolsProduct, region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	listPoolsPath := listPoolsClient.Endpoint + listPoolsHttpUrl
	listPoolsPath = strings.ReplaceAll(listPoolsPath, "{project_id}", listPoolsClient.ProjectID)

	listPoolsQueryParams := buildListPoolsQueryParams(d)
	listPoolsPath += listPoolsQueryParams

	listPoolsResp, err := pagination.ListAllItems(
		listPoolsClient,
		"marker",
		listPoolsPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return diag.Errorf("error retrieving ELB pools: %s", err)
	}

	listPoolsRespJson, err := json.Marshal(listPoolsResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var listPoolsRespBody interface{}
	err = json.Unmarshal(listPoolsRespJson, &listPoolsRespBody)
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
		d.Set("pools", flattenListPoolsBodyPools(listPoolsRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListPoolsBodyPools(resp interface{}) []interface{} {
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
			"id":                                utils.PathSearch("id", v, nil),
			"name":                              utils.PathSearch("name", v, nil),
			"description":                       utils.PathSearch("description", v, nil),
			"protocol":                          utils.PathSearch("protocol", v, nil),
			"ip_version":                        utils.PathSearch("ip_version", v, nil),
			"lb_method":                         utils.PathSearch("lb_algorithm", v, nil),
			"healthmonitor_id":                  utils.PathSearch("healthmonitor_id", v, nil),
			"type":                              utils.PathSearch("type", v, nil),
			"vpc_id":                            utils.PathSearch("vpc_id", v, nil),
			"any_port_enable":                   utils.PathSearch("any_port_enable", v, nil),
			"enterprise_project_id":             utils.PathSearch("enterprise_project_id", v, nil),
			"member_deletion_protection_enable": utils.PathSearch("member_deletion_protection_enable", v, nil),
			"public_border_group":               utils.PathSearch("public_border_group", v, nil),
			"protection_status":                 utils.PathSearch("protection_status", v, nil),
			"protection_reason":                 utils.PathSearch("protection_reason", v, nil),
			"slow_start_enabled":                utils.PathSearch("slow_start.enable", v, nil),
			"slow_start_duration":               utils.PathSearch("slow_start.duration", v, nil),
			"connection_drain_enabled":          utils.PathSearch("connection_drain.enable", v, nil),
			"connection_drain_timeout":          utils.PathSearch("connection_drain.timeout", v, nil),
			"minimum_healthy_member_count":      utils.PathSearch("pool_health.minimum_healthy_member_count", v, nil),
			"listeners":                         flattenPoolListeners(v),
			"loadbalancers":                     flattenPoolLoadBalancers(v),
			"members":                           flattenPoolMembers(v),
			"persistence":                       flattenPersistence(v),
			"quic_cid_hash_strategy":            flattenPoolQuicCidHashStrategy(v),
			"created_at":                        utils.PathSearch("created_at", v, nil),
			"updated_at":                        utils.PathSearch("updated_at", v, nil),
		})
	}
	return rst
}

func flattenPoolListeners(resp interface{}) []interface{} {
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

func flattenPoolLoadBalancers(resp interface{}) []interface{} {
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

func flattenPoolMembers(resp interface{}) []interface{} {
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
			"id": utils.PathSearch("id", v, nil),
		})
	}
	return rst
}

func flattenPersistence(resp interface{}) []map[string]interface{} {
	persistence := utils.PathSearch("session_persistence", resp, nil)
	if persistence == nil {
		return nil
	}

	rst := []map[string]interface{}{
		{
			"cookie_name": utils.PathSearch("cookie_name", persistence, nil),
			"type":        utils.PathSearch("type", persistence, nil),
			"timeout":     utils.PathSearch("persistence_timeout", persistence, nil),
		},
	}
	return rst
}

func flattenPoolQuicCidHashStrategy(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
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

func buildListPoolsQueryParams(d *schema.ResourceData) string {
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
	if v, ok := d.GetOk("any_port_enable"); ok {
		anyPortEnable, _ := strconv.ParseBool(v.(string))
		res = fmt.Sprintf("%s&any_port_enable=%v", res, anyPortEnable)
	}
	if v, ok := d.GetOk("connection_drain"); ok {
		connectionDrain, _ := strconv.ParseBool(v.(string))
		res = fmt.Sprintf("%s&connection_drain=%v", res, connectionDrain)
	}
	if v, ok := d.GetOk("enterprise_project_id"); ok {
		res = fmt.Sprintf("%s&enterprise_project_id=%v", res, v)
	}
	if v, ok := d.GetOk("ip_version"); ok {
		res = fmt.Sprintf("%s&ip_version=%v", res, v)
	}
	if v, ok := d.GetOk("member_address"); ok {
		res = fmt.Sprintf("%s&member_address=%v", res, v)
	}
	if v, ok := d.GetOk("member_device_id"); ok {
		res = fmt.Sprintf("%s&member_device_id=%v", res, v)
	}
	if v, ok := d.GetOk("member_instance_id"); ok {
		res = fmt.Sprintf("%s&member_instance_id=%v", res, v)
	}
	if v, ok := d.GetOk("member_deletion_protection_enable"); ok {
		memberDeletionProtectionEnable, _ := strconv.ParseBool(v.(string))
		res = fmt.Sprintf("%s&member_deletion_protection_enable=%v", res, memberDeletionProtectionEnable)
	}
	if v, ok := d.GetOk("pool_health"); ok {
		res = fmt.Sprintf("%s&pool_health=%v", res, v)
	}
	if v, ok := d.GetOk("public_border_group"); ok {
		res = fmt.Sprintf("%s&public_border_group=%v", res, v)
	}
	if v, ok := d.GetOk("quic_cid_len"); ok {
		res = fmt.Sprintf("%s&quic_cid_len=%v", res, v)
	}
	if v, ok := d.GetOk("quic_cid_offset"); ok {
		res = fmt.Sprintf("%s&quic_cid_offset=%v", res, v)
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
	if v, ok := d.GetOk("lb_method"); ok {
		res = fmt.Sprintf("%s&lb_algorithm=%v", res, v)
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
	if v, ok := d.GetOk("protection_status"); ok {
		res = fmt.Sprintf("%s&protection_status=%v", res, v)
	}
	if res != "" {
		res = "?" + res[1:]
	}
	return res
}
