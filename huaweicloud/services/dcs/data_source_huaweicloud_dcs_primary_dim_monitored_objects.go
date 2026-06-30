package dcs

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DCS GET /v2/{project_id}/dims/monitored-objects
func DataSourceDcsPrimaryDimMonitoredObjects() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDcsPrimaryDimMonitoredObjectsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"dim_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"router": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"children": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dim_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dim_route": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dcs_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceDcsPrimaryDimMonitoredObjectsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("dcs", region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	httpUrl := "v2/{project_id}/dims/monitored-objects"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath += buildGetPrimaryDimMonitoredObjectsQueryParams(d)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	listResp, err := client.Request("GET", listPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving DCS primary dim monitored objects: %s", err)
	}
	listRespBody, err := utils.FlattenResponse(listResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id.String())

	routerRaw := utils.PathSearch("router", listRespBody, nil)
	routerList := make([]string, 0)

	if routerRaw != nil {
		for _, v := range routerRaw.([]interface{}) {
			if str, ok := v.(string); ok {
				routerList = append(routerList, str)
			}
		}
	}

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("router", routerList),
		d.Set("children", flattenGetPrimaryDimMonitoredObjectsChildren(listRespBody)),
		d.Set("instances", flattenGetPrimaryDimMonitoredObjectsInstances(listRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetPrimaryDimMonitoredObjectsQueryParams(d *schema.ResourceData) string {
	res := fmt.Sprintf("?dim_name=%v", d.Get("dim_name"))
	return res
}

func flattenGetPrimaryDimMonitoredObjectsChildren(resp interface{}) []interface{} {
	curJson := utils.PathSearch("children", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"dim_name":  utils.PathSearch("dim_name", v, nil),
			"dim_route": utils.PathSearch("dim_route", v, nil),
		})
	}
	return res
}

func flattenGetPrimaryDimMonitoredObjectsInstances(resp interface{}) []interface{} {
	curJson := utils.PathSearch("instances", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"dcs_instance_id": utils.PathSearch("dcs_instance_id", v, nil),
			"name":            utils.PathSearch("name", v, nil),
			"status":          utils.PathSearch("status", v, nil),
		})
	}
	return res
}
