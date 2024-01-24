package elb

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

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
		rst = append(rst, map[string]interface{}{
			"id":                          utils.PathSearch("id", v, nil),
			"name":                        utils.PathSearch("name", v, nil),
			"description":                 utils.PathSearch("description", v, nil),
			"protocol":                    utils.PathSearch("protocol", v, nil),
			"protocol_port":               utils.PathSearch("protocol_port", v, nil),
			"default_pool_id":             utils.PathSearch("default_pool_id", v, nil),
			"http2_enable":                utils.PathSearch("http2_enable", v, nil),
			"forward_eip":                 utils.PathSearch("insert_headers.X-Forwarded-ELB-IP", v, nil),
			"forward_port":                utils.PathSearch("insert_headers.X-Forwarded-Port", v, nil),
			"forward_request_port":        utils.PathSearch("insert_headers.X-Forwarded-For-Port", v, nil),
			"forward_host":                utils.PathSearch("insert_headers.X-Forwarded-Host", v, nil),
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
			"protection_reason":           utils.PathSearch("protection_reason", v, nil),
		})
	}
	return rst
}
