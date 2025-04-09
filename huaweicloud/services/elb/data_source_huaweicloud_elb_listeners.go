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

// @API ELB GET /v3/{project_id}/elb/listeners
func DataSourceElbListeners() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceElbListenersRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"loadbalancer_id": {
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
			"ca_certificate": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"request_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"default_pool_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"server_certificate": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_member_retry": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"true", "false",
				}, false),
			},
			"advanced_forwarding_enabled": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"true", "false",
				}, false),
			},
			"http2_enable": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"true", "false",
				}, false),
			},
			"idle_timeout": {
				Type:     schema.TypeInt,
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
			"response_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"protection_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"proxy_protocol_enable": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"true", "false",
				}, false),
			},
			"ssl_early_data_enable": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"true", "false",
				}, false),
			},
			"tls_ciphers_policy": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"protocol_port": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"listener_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"listeners": {
				Type:     schema.TypeList,
				Elem:     listenersSchema(),
				Computed: true,
			},
		},
	}
}

func listenersSchema() *schema.Resource {
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
			"protocol_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"default_pool_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"http2_enable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"forward_elb": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"forward_eip": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"forward_port": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"forward_request_port": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"forward_host": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"forward_proto": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"forward_tls_certificate": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"forward_tls_cipher": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"forward_tls_protocol": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"real_ip": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"ipgroup": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     listenerIpGroupsSchema(),
			},
			"port_ranges": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     listenerPortRangesSchema(),
			},
			"sni_certificate": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"server_certificate": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ca_certificate": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tls_ciphers_policy": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"idle_timeout": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"request_timeout": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"response_timeout": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"loadbalancer_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"advanced_forwarding_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"protection_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"protection_reason": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"proxy_protocol_enable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"quic_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     listenerQuicConfigSchema(),
			},
			"security_policy_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sni_match_algo": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ssl_early_data_enable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"max_connection": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cps": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"enable_member_retry": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gzip_enable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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

func listenerIpGroupsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"enable_ipgroup": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"ipgroup_id": {
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

func listenerPortRangesSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"start_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"end_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
	return &sc
}

func listenerQuicConfigSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"enable_quic_upgrade": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"quic_listener_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func dataSourceElbListenersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		listListenersHttpUrl = "v3/{project_id}/elb/listeners"
		listListenersProduct = "elb"
	)
	listListenersClient, err := cfg.NewServiceClient(listListenersProduct, region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}
	listListenersPath := listListenersClient.Endpoint + listListenersHttpUrl
	listListenersPath = strings.ReplaceAll(listListenersPath, "{project_id}", listListenersClient.ProjectID)
	listListenersQueryParams := buildListListenersQueryParams(d, cfg.GetEnterpriseProjectID(d))
	listListenersPath += listListenersQueryParams
	listListenersResp, err := pagination.ListAllItems(
		listListenersClient,
		"marker",
		listListenersPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving ELB listeners")
	}

	listListenersRespJson, err := json.Marshal(listListenersResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var listListenersRespBody interface{}
	err = json.Unmarshal(listListenersRespJson, &listListenersRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("listeners", flattenListListenersBody(listListenersRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildListListenersQueryParams(d *schema.ResourceData, enterpriseProjectId string) string {
	res := ""
	if v, ok := d.GetOk("loadbalancer_id"); ok {
		res = fmt.Sprintf("%s&loadbalancer_id=%v", res, v)
	}
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	if v, ok := d.GetOk("description"); ok {
		res = fmt.Sprintf("%s&description=%v", res, v)
	}
	if v, ok := d.GetOk("ca_certificate"); ok {
		res = fmt.Sprintf("%s&client_ca_tls_container_ref=%v", res, v)
	}
	if v, ok := d.GetOk("request_timeout"); ok {
		res = fmt.Sprintf("%s&client_timeout=%v", res, v)
	}
	if v, ok := d.GetOk("default_pool_id"); ok {
		res = fmt.Sprintf("%s&default_pool_id=%v", res, v)
	}
	if v, ok := d.GetOk("server_certificate"); ok {
		res = fmt.Sprintf("%s&default_tls_container_ref=%v", res, v)
	}
	if v, ok := d.GetOk("enable_member_retry"); ok {
		enableMemberRetry, _ := strconv.ParseBool(v.(string))
		res = fmt.Sprintf("%s&enable_member_retry=%v", res, enableMemberRetry)
	}
	if v, ok := d.GetOk("advanced_forwarding_enabled"); ok {
		enableMemberRetry, _ := strconv.ParseBool(v.(string))
		res = fmt.Sprintf("%s&enhance_l7policy_enable=%v", res, enableMemberRetry)
	}
	if v, ok := d.GetOk("http2_enable"); ok {
		http2Enable, _ := strconv.ParseBool(v.(string))
		res = fmt.Sprintf("%s&http2_enable=%v", res, http2Enable)
	}
	if v, ok := d.GetOk("idle_timeout"); ok {
		res = fmt.Sprintf("%s&keepalive_timeout=%v", res, v)
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
	if v, ok := d.GetOk("response_timeout"); ok {
		res = fmt.Sprintf("%s&member_timeout=%v", res, v)
	}
	if v, ok := d.GetOk("protection_status"); ok {
		res = fmt.Sprintf("%s&protection_status=%v", res, v)
	}
	if v, ok := d.GetOk("proxy_protocol_enable"); ok {
		proxyProtocolEnable, _ := strconv.ParseBool(v.(string))
		res = fmt.Sprintf("%s&proxy_protocol_enable=%v", res, proxyProtocolEnable)
	}
	if v, ok := d.GetOk("ssl_early_data_enable"); ok {
		sslEarlyDataEnable, _ := strconv.ParseBool(v.(string))
		res = fmt.Sprintf("%s&ssl_early_data_enable=%v", res, sslEarlyDataEnable)
	}
	if v, ok := d.GetOk("tls_ciphers_policy"); ok {
		res = fmt.Sprintf("%s&tls_ciphers_policy=%v", res, v)
	}
	if v, ok := d.GetOk("listener_id"); ok {
		res = fmt.Sprintf("%s&id=%v", res, v)
	}
	if v, ok := d.GetOk("protocol"); ok {
		res = fmt.Sprintf("%s&protocol=%v", res, v)
	}
	if v, ok := d.GetOk("protocol_port"); ok {
		res = fmt.Sprintf("%s&protocol_port=%v", res, v)
	}
	if enterpriseProjectId != "" {
		res = fmt.Sprintf("%s&enterprise_project_id=%v", res, enterpriseProjectId)
	}
	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func flattenListListenersBody(resp interface{}) []interface{} {
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
		insertHeaders := utils.PathSearch("insert_headers", v, nil)
		rst = append(rst, map[string]interface{}{
			"id":                          utils.PathSearch("id", v, nil),
			"name":                        utils.PathSearch("name", v, nil),
			"description":                 utils.PathSearch("description", v, nil),
			"protocol":                    utils.PathSearch("protocol", v, nil),
			"protocol_port":               utils.PathSearch("protocol_port", v, nil),
			"default_pool_id":             utils.PathSearch("default_pool_id", v, nil),
			"http2_enable":                utils.PathSearch("http2_enable", v, nil),
			"forward_elb":                 utils.PathSearch(`"X-Forwarded-ELB-ID"`, insertHeaders, nil),
			"forward_eip":                 utils.PathSearch(`"X-Forwarded-ELB-IP"`, insertHeaders, nil),
			"forward_port":                utils.PathSearch(`"X-Forwarded-Port"`, insertHeaders, nil),
			"forward_request_port":        utils.PathSearch(`"X-Forwarded-For-Port"`, insertHeaders, nil),
			"forward_host":                utils.PathSearch(`"X-Forwarded-Host"`, insertHeaders, nil),
			"forward_proto":               utils.PathSearch(`"X-Forwarded-Proto"`, insertHeaders, nil),
			"forward_tls_certificate":     utils.PathSearch(`"X-Forwarded-TLS-Certificate-ID"`, insertHeaders, nil),
			"forward_tls_cipher":          utils.PathSearch(`"X-Forwarded-TLS-Cipher"`, insertHeaders, nil),
			"forward_tls_protocol":        utils.PathSearch(`"X-Forwarded-TLS-Protocol"`, insertHeaders, nil),
			"real_ip":                     utils.PathSearch(`"X-Real-IP"`, insertHeaders, nil),
			"ipgroup":                     flattenListenersIpGroup(v),
			"port_ranges":                 flattenListenersPortRanges(v),
			"sni_certificate":             utils.PathSearch("sni_container_refs", v, nil),
			"server_certificate":          utils.PathSearch("default_tls_container_ref", v, nil),
			"ca_certificate":              utils.PathSearch("client_ca_tls_container_ref", v, nil),
			"tls_ciphers_policy":          utils.PathSearch("tls_ciphers_policy", v, nil),
			"idle_timeout":                utils.PathSearch("keepalive_timeout", v, nil),
			"request_timeout":             utils.PathSearch("client_timeout", v, nil),
			"response_timeout":            utils.PathSearch("member_timeout", v, nil),
			"loadbalancer_id":             utils.PathSearch("loadbalancers|[0].id", v, nil),
			"advanced_forwarding_enabled": utils.PathSearch("enhance_l7policy_enable", v, nil),
			"protection_status":           utils.PathSearch("protection_status", v, nil),
			"proxy_protocol_enable":       utils.PathSearch("proxy_protocol_enable", v, nil),
			"quic_config":                 flattenListenersQuicConfig(v),
			"security_policy_id":          utils.PathSearch("security_policy_id", v, nil),
			"sni_match_algo":              utils.PathSearch("sni_match_algo", v, nil),
			"ssl_early_data_enable":       utils.PathSearch("ssl_early_data_enable", v, nil),
			"max_connection":              utils.PathSearch("connection", v, nil),
			"cps":                         utils.PathSearch("cps", v, nil),
			"enable_member_retry":         utils.PathSearch("enable_member_retry", v, nil),
			"enterprise_project_id":       utils.PathSearch("enterprise_project_id", v, nil),
			"gzip_enable":                 utils.PathSearch("gzip_enable", v, nil),
			"tags":                        utils.FlattenTagsToMap(utils.PathSearch("tags", resp, make([]interface{}, 0))),
			"created_at":                  utils.PathSearch("created_at", v, nil),
			"updated_at":                  utils.PathSearch("updated_at", v, nil),
		})
	}
	return rst
}

func flattenListenersIpGroup(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("ipgroup", resp, nil)
	if curJson == nil {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"enable_ipgroup": utils.PathSearch("enable_ipgroup", curJson, nil),
			"ipgroup_id":     utils.PathSearch("ipgroup_id", curJson, nil),
			"type":           utils.PathSearch("type", curJson, nil),
		},
	}
	return rst
}

func flattenListenersPortRanges(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curJson := utils.PathSearch("port_ranges", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"start_port": utils.PathSearch("start_port", v, nil),
			"end_port":   utils.PathSearch("end_port", v, nil),
		})
	}
	return rst
}

func flattenListenersQuicConfig(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("quic_config", resp, nil)
	if curJson == nil {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"enable_quic_upgrade": utils.PathSearch("enable_quic_upgrade", curJson, nil),
			"quic_listener_id":    utils.PathSearch("quic_listener_id", curJson, nil),
		},
	}
	return rst
}
