package elb

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func DataSourceElbListeners() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceElbListeners,

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
			"protocol": {
				Type: schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{
					"TCP", "UDP", "HTTP", "HTTPS",
				}, false),
				Optional: true,
			},
			"protocol_port": {
				Type:     schema.TypeString,
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
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"protocol": {
				Type:     schema.TypeList,
				Computed: true,
			},
			"protocol_port": {
				Type:     schema.TypeString,
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
				Type:     schema.TypeString,
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
			}, "protection_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"protection_reason": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func dataSourceElbListeners(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	var (
		listLoadListenersHttpUrl = "v3/{project_id}/elb/listeners"
		listLoadListenersProduct = "elb"
	)
	listLoadListenersClient, err := cfg.NewServiceClient(listLoadListenersProduct, cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}
	listLoadListenersPath := listLoadListenersClient.Endpoint + listLoadListenersHttpUrl
	listLoadListenersPath = strings.ReplaceAll(listLoadListenersPath, "{project_id}", listLoadListenersClient.ProjectID)

	listLoadListenersQueryParams := buildListListenersQueryParams(d)
	listLoadListenersPath += listLoadListenersQueryParams

	listLoadListenersResp, err := pagination.ListAllItems(
		listLoadListenersClient,
		"marker",
		listLoadListenersPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Listeners")
	}

	listListenersRespJson, err := json.Marshal(listLoadListenersResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var listListenersRespBody interface{}
	err = json.Unmarshal(listListenersRespJson, &listListenersRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(
		d.Set("region", cfg.GetRegion(d)),
		d.Set("loadbalancers", flattenListListenersBody(listListenersRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildListListenersQueryParams(d *schema.ResourceData) string {
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
	if v, ok := d.GetOk("listener_id"); ok {
		res = fmt.Sprintf("%s&listener_id=%v", res, v)
	}
	if v, ok := d.GetOk("protocol"); ok {
		res = fmt.Sprintf("%s&protocol=%v", res, v)
	}
	if v, ok := d.GetOk("protocol_port"); ok {
		res = fmt.Sprintf("%s&protocol_port=%v", res, v)
	}
	if v, ok := d.GetOk("enterprise_project_id"); ok {
		res = fmt.Sprintf("%s&enterprise_project_id=%v", res, v)
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
		rst = append(rst, map[string]interface{}{
			"listener_id":                 utils.PathSearch("id", v, nil),
			"name":                        utils.PathSearch("name", v, nil),
			"description":                 utils.PathSearch("description", v, nil),
			"protocol":                    utils.PathSearch("protocol", v, nil),
			"protocol_port":               utils.PathSearch("protocol_port", v, nil),
			"default_pool_id":             utils.PathSearch("default_pool_id", v, nil),
			"http2_enable":                utils.PathSearch("http2_enable", v, nil),
			"forward_eip":                 utils.PathSearch("X-Forwarded-ELB-IP", v, nil),
			"forward_port":                utils.PathSearch("X-Forwarded-Port", v, nil),
			"forward_request_port":        utils.PathSearch("X-Forwarded-For-Port", v, nil),
			"forward_host":                utils.PathSearch("X-Forwarded-Host", v, nil),
			"sni_certificate":             utils.PathSearch("sni_container_refs", v, nil),
			"server_certificate":          utils.PathSearch("default_tls_container_ref", v, nil),
			"ca_certificate":              utils.PathSearch("client_ca_tls_container_ref", v, nil),
			"tls_ciphers_policy":          utils.PathSearch("tls_ciphers_policy", v, nil),
			"idle_timeout":                utils.PathSearch("keepalive_timeout", v, nil),
			"request_timeout":             utils.PathSearch("client_timeout", v, nil),
			"response_timeout":            utils.PathSearch("member_timeout", v, nil),
			"loadbalancer_id":             utils.PathSearch("loadbalancers", v, nil),
			"advanced_forwarding_enabled": utils.PathSearch("enhance_l7policy_enable", v, nil),
			"protection_status":           utils.PathSearch("protection_status", v, nil),
			"protection_reason":           utils.PathSearch("protection_reason", v, nil),
		})
	}
	return rst
}
