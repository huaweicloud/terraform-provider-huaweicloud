// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product ELB
// ---------------------------------------------------------------

package elb

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/chnsz/golangsdk/pagination"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/jmespath/go-jmespath"
)

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
			"healthmonitor_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Specifies the health monitor ID of the ELB pool.`,
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
				Type:     schema.TypeList,
				Elem:     poolsPoolPersistenceSchema(),
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
		},
	}
	return &sc
}

func resourcePoolsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)

	var mErr *multierror.Error

	// listPools: Query the List of ELB pools
	var (
		listPoolsHttpUrl = "v3/{project_id}/elb/pools"
		listPoolsProduct = "elb"
	)
	listPoolsClient, err := config.NewServiceClient(listPoolsProduct, region)
	if err != nil {
		return diag.Errorf("error creating Pools Client: %s", err)
	}

	listPoolsPath := listPoolsClient.Endpoint + listPoolsHttpUrl
	listPoolsPath = strings.Replace(listPoolsPath, "{project_id}", listPoolsClient.ProjectID, -1)

	listPoolsqueryParams := buildListPoolsQueryParams(d)
	listPoolsPath = listPoolsPath + listPoolsqueryParams

	listPoolsResp, err := pagination.ListAllItems(
		listPoolsClient,
		"marker",
		listPoolsPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Pools")
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

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

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
			"id":               utils.PathSearch("id", v, nil),
			"name":             utils.PathSearch("name", v, nil),
			"description":      utils.PathSearch("description", v, nil),
			"protocol":         utils.PathSearch("protocol", v, nil),
			"ip_version":       utils.PathSearch("ip_version", v, nil),
			"lb_method":        utils.PathSearch("lb_algorithm", v, nil),
			"healthmonitor_id": utils.PathSearch("healthmonitor_id", v, nil),
			"listeners":        flattenPoolListeners(v),
			"loadbalancers":    flattenPoolLoadbalancers(v),
			"members":          flattenPoolMembers(v),
			"persistence":      flattenPoolPersistence(v),
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

func flattenPoolLoadbalancers(resp interface{}) []interface{} {
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

func flattenPoolPersistence(resp interface{}) []interface{} {
	var rst []interface{}
	curJson, err := jmespath.Search("session_persistence", resp)
	if err != nil {
		log.Printf("[ERROR] error parsing persistence from response= %#v", resp)
		return rst
	}
	if curJson == nil {
		return nil
	}

	rst = []interface{}{
		map[string]interface{}{
			"type":        utils.PathSearch("type", curJson, nil),
			"cookie_name": utils.PathSearch("cookie_name", curJson, nil),
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
	if res != "" {
		res = "?" + res[1:]
	}
	return res
}
