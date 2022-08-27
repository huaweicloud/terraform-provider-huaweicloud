// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product ELB
// ---------------------------------------------------------------

package lb

import (
	"context"
	"fmt"
	"strings"

	"github.com/chnsz/golangsdk"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

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
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The listener name.`,
			},
			"protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The listener protocol.`,
			},
			"protocol_port": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The front-end listening port of the listener.`,
			},
			"listeners": {
				Type:        schema.TypeList,
				Elem:        listenersListenersSchema(),
				Computed:    true,
				Description: `Listener list.`,
			},
		},
	}
}

func listenersListenersSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ELB listener ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The listener name.`,
			},
			"protocol": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The listener protocol.`,
			},
			"protocol_port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The front-end listening port of the listener.`,
			},
			"default_pool_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the default pool with which the ELB listener is associated.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description of the ELB listener.`,
			},
			"connection_limit": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The maximum number of connections allowed for the listener.`,
			},
			"http2_enable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the ELB listener uses HTTP/2.`,
			},
			"default_tls_container_ref": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the server certificate used by the listener.`,
			},
			"sni_container_refs": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: `List of the SNI certificate (server certificates with a domain name) IDs used by the listener.`,
			},
			"loadbalancers": {
				Type:        schema.TypeList,
				Elem:        listenersListenerLoadbalancersSchema(),
				Computed:    true,
				Description: `Loadbalancer list. For details, see Data structure of the loadbalancer field.`,
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

func resourceListenersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)

	var mErr *multierror.Error

	// listListeners: Query the list of Shared ELB listeners
	var (
		listListenersHttpUrl = "v2/{project_id}/elb/listeners"
		listListenersProduct = "elb"
	)
	listListenersClient, err := config.NewServiceClient(listListenersProduct, region)
	if err != nil {
		return diag.Errorf("error creating Listeners Client: %s", err)
	}

	listListenersPath := listListenersClient.Endpoint + listListenersHttpUrl
	listListenersPath = strings.Replace(listListenersPath, "{project_id}", listListenersClient.ProjectID, -1)

	listListenersqueryParams := buildListListenersQueryParams(d)
	listListenersPath = listListenersPath + listListenersqueryParams

	listListenersOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	listListenersResp, err := listListenersClient.Request("GET", listListenersPath, &listListenersOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Listeners")
	}

	listListenersRespBody, err := utils.FlattenResponse(listListenersResp)
	if err != nil {
		return diag.FromErr(err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

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
			"id":                        utils.PathSearch("id", v, nil),
			"name":                      utils.PathSearch("name", v, nil),
			"protocol":                  utils.PathSearch("protocol", v, nil),
			"protocol_port":             utils.PathSearch("protocol_port", v, nil),
			"default_pool_id":           utils.PathSearch("default_pool_id", v, nil),
			"description":               utils.PathSearch("description", v, nil),
			"connection_limit":          utils.PathSearch("connection_limit", v, nil),
			"http2_enable":              utils.PathSearch("http2_enable", v, nil),
			"default_tls_container_ref": utils.PathSearch("default_tls_container_ref", v, nil),
			"sni_container_refs":        utils.PathSearch("sni_container_refs", v, nil),
			"loadbalancers":             flattenListenerLoadbalancers(v),
		})
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
	if res != "" {
		res = "?" + res[1:]
	}
	return res
}
