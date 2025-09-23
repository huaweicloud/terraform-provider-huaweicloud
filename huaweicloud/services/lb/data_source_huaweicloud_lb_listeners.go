package lb

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ELB GET /v2/{project_id}/elb/listeners
func DataSourceListeners() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceListenersRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"protocol_port": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"client_ca_tls_container_ref": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"default_pool_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"default_tls_container_ref": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"http2_enable": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"true", "false",
				}, false),
			},
			"listener_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"loadbalancer_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tls_ciphers_policy": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"listeners": {
				Type:     schema.TypeList,
				Elem:     listenersListenersSchema(),
				Computed: true,
			},
		},
	}
}

func listenersListenersSchema() *schema.Resource {
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
			"client_ca_tls_container_ref": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"connection_limit": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"http2_enable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"insert_headers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     listenersListenerInsertHeadersSchema(),
			},
			"default_tls_container_ref": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sni_container_refs": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
			"tls_ciphers_policy": {
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
			"loadbalancers": {
				Type:     schema.TypeList,
				Elem:     listenersListenerLoadbalancersSchema(),
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
	return &sc
}

func listenersListenerInsertHeadersSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"x_forwarded_elb_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"x_forwarded_host": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func listenersListenerLoadbalancersSchema() *schema.Resource {
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

func resourceListenersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// listListeners: Query the list of Shared ELB listeners
	var (
		listListenersHttpUrl = "v2/{project_id}/elb/listeners"
		listListenersProduct = "elb"
	)
	listListenersClient, err := cfg.NewServiceClient(listListenersProduct, region)
	if err != nil {
		return diag.Errorf("error creating Listeners client: %s", err)
	}

	listListenersPath := listListenersClient.Endpoint + listListenersHttpUrl
	listListenersPath = strings.ReplaceAll(listListenersPath, "{project_id}", listListenersClient.ProjectID)

	listListenersPath += buildListListenersQueryParams(d)

	listListenersOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	listListenersResp, err := listListenersClient.Request("GET", listListenersPath, &listListenersOpt)
	if err != nil {
		return diag.Errorf("error retrieving listeners: %s", err)
	}

	listListenersRespBody, err := utils.FlattenResponse(listListenersResp)
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
		d.Set("listeners", flattenListListenersBodyListeners(listListenersRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListListenersBodyListeners(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("listeners", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":                          utils.PathSearch("id", v, nil),
			"name":                        utils.PathSearch("name", v, nil),
			"protocol":                    utils.PathSearch("protocol", v, nil),
			"protocol_port":               utils.PathSearch("protocol_port", v, nil),
			"default_pool_id":             utils.PathSearch("default_pool_id", v, nil),
			"client_ca_tls_container_ref": utils.PathSearch("client_ca_tls_container_ref", v, nil),
			"description":                 utils.PathSearch("description", v, nil),
			"connection_limit":            utils.PathSearch("connection_limit", v, nil),
			"http2_enable":                utils.PathSearch("http2_enable", v, nil),
			"insert_headers":              flattenListenerInsertHeaders(v),
			"default_tls_container_ref":   utils.PathSearch("default_tls_container_ref", v, nil),
			"sni_container_refs":          utils.PathSearch("sni_container_refs", v, nil),
			"protection_status":           utils.PathSearch("protection_status", v, nil),
			"protection_reason":           utils.PathSearch("protection_reason", v, nil),
			"tls_ciphers_policy":          utils.PathSearch("tls_ciphers_policy", v, nil),
			"created_at":                  utils.PathSearch("created_at", v, nil),
			"updated_at":                  utils.PathSearch("updated_at", v, nil),
			"loadbalancers":               flattenListenerLoadbalancers(v),
			"tags":                        utils.FlattenTagsToMap(utils.PathSearch("tags", resp, make([]interface{}, 0))),
		})
	}
	return rst
}

func flattenListenerInsertHeaders(resp interface{}) []map[string]interface{} {
	curJson := utils.PathSearch("insert_headers", resp, nil)
	if curJson == nil {
		return nil
	}

	xForwardedElbIp := utils.PathSearch(`"X-Forwarded-ELB-IP"`, curJson, false).(bool)
	xForwardedHost := utils.PathSearch(`"X-Forwarded-Host"`, curJson, false).(bool)
	rst := []map[string]interface{}{
		{
			"x_forwarded_elb_ip": strconv.FormatBool(xForwardedElbIp),
			"x_forwarded_host":   strconv.FormatBool(xForwardedHost),
		},
	}
	return rst
}

func flattenListenerLoadbalancers(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("loadbalancers", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id": utils.PathSearch("id", v, nil),
		})
	}
	return rst
}

func buildListListenersQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	if v, ok := d.GetOk("protocol"); ok {
		res = fmt.Sprintf("%s&protocol=%v", res, v)
	}
	if v, ok := d.GetOk("protocol_port"); ok {
		res = fmt.Sprintf("%s&protocol_port=%v", res, v)
	}
	if v, ok := d.GetOk("client_ca_tls_container_ref"); ok {
		res = fmt.Sprintf("%s&client_ca_tls_container_ref=%v", res, v)
	}
	if v, ok := d.GetOk("default_pool_id"); ok {
		res = fmt.Sprintf("%s&default_pool_id=%v", res, v)
	}
	if v, ok := d.GetOk("default_tls_container_ref"); ok {
		res = fmt.Sprintf("%s&default_tls_container_ref=%v", res, v)
	}
	if v, ok := d.GetOk("description"); ok {
		res = fmt.Sprintf("%s&description=%v", res, v)
	}
	if v, ok := d.GetOk("enterprise_project_id"); ok {
		res = fmt.Sprintf("%s&enterprise_project_id=%v", res, v)
	}
	if v, ok := d.GetOk("http2_enable"); ok {
		http2Enable, _ := strconv.ParseBool(v.(string))
		res = fmt.Sprintf("%s&http2_enable=%v", res, http2Enable)
	}
	if v, ok := d.GetOk("listener_id"); ok {
		res = fmt.Sprintf("%s&id=%v", res, v)
	}
	if v, ok := d.GetOk("loadbalancer_id"); ok {
		res = fmt.Sprintf("%s&loadbalancer_id=%v", res, v)
	}
	if v, ok := d.GetOk("tls_ciphers_policy"); ok {
		res = fmt.Sprintf("%s&tls_ciphers_policy=%v", res, v)
	}
	if res != "" {
		res = "?" + res[1:]
	}
	return res
}
