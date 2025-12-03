package elb

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ELB GET /v3/{project_id}/elb/loadbalancers/{loadbalancer_id}/local-ports
func DataSourceElbLoadBalancerPorts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceElbLoadBalancerPortsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"loadbalancer_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"port_id": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ip_address": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ipv6_address": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"type": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"virsubnet_id": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ports": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     loadBalancerPortsSchema(),
			},
		},
	}
}

func loadBalancerPortsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"port_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv6_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"virsubnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func dataSourceElbLoadBalancerPortsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		httpUrl = "v3/{project_id}/elb/loadbalancers/{loadbalancer_id}/local-ports"
		product = "elb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{loadbalancer_id}", d.Get("loadbalancer_id").(string))
	getQueryParams := buildGetLoadBalancerPortsQueryParams(d)
	getPath += getQueryParams

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving ELB load balancer ports: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
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
		d.Set("ports", flattenLoadBalancerPortsBody(getRespBody)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetLoadBalancerPortsQueryParams(d *schema.ResourceData) string {
	res := ""
	if raw, ok := d.GetOk("port_id"); ok {
		for _, v := range raw.([]interface{}) {
			res = fmt.Sprintf("%s&port_id=%v", res, v)
		}
	}
	if raw, ok := d.GetOk("ip_address"); ok {
		for _, v := range raw.([]interface{}) {
			res = fmt.Sprintf("%s&ip_address=%v", res, v)
		}
	}
	if raw, ok := d.GetOk("ipv6_address"); ok {
		for _, v := range raw.([]interface{}) {
			res = fmt.Sprintf("%s&ipv6_address=%v", res, v)
		}
	}
	if raw, ok := d.GetOk("type"); ok {
		for _, v := range raw.([]interface{}) {
			res = fmt.Sprintf("%s&type=%v", res, v)
		}
	}
	if raw, ok := d.GetOk("virsubnet_id"); ok {
		for _, v := range raw.([]interface{}) {
			res = fmt.Sprintf("%s&virsubnet_id=%v", res, v)
		}
	}
	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func flattenLoadBalancerPortsBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("ports", resp, make([]interface{}, 0))
	if curJson == nil {
		return nil
	}
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"port_id":      utils.PathSearch("port_id", v, nil),
			"ip_address":   utils.PathSearch("ip_address", v, nil),
			"ipv6_address": utils.PathSearch("ipv6_address", v, nil),
			"type":         utils.PathSearch("type", v, nil),
			"virsubnet_id": utils.PathSearch("virsubnet_id", v, nil),
		})
	}
	return rst
}
