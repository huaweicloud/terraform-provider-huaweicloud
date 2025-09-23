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

// @API ELB GET /v3/{project_id}/elb/availability-zones
func DataSourceAvailabilityZones() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAvailabilityZonesRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"public_border_group": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"loadbalancer_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"availability_zones": {
				Type:     schema.TypeList,
				Elem:     availabilityZonesListSchema(),
				Computed: true,
			},
		},
	}
}

func availabilityZonesListSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"list": {
				Type:     schema.TypeList,
				Elem:     availabilityZonesSchema(),
				Computed: true,
			},
		},
	}
	return &sc
}

func availabilityZonesSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"protocol": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"public_border_group": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"category": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
	return &sc
}

func dataSourceAvailabilityZonesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v3/{project_id}/elb/availability-zones"
		product = "elb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)

	getPath += buildGetAvailabilityZonesQueryParams(d)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving ELB availability zones: %s", err)
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

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("availability_zones", flattenGetAvailabilityZonesBodyPools(getRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetAvailabilityZonesBodyPools(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("availability_zones", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		azList := v.([]interface{})
		list := make([]interface{}, 0, len(azList))
		for _, az := range azList {
			list = append(list, map[string]interface{}{
				"code":                utils.PathSearch("code", az, nil),
				"state":               utils.PathSearch("state", az, nil),
				"protocol":            utils.PathSearch("protocol", az, nil),
				"public_border_group": utils.PathSearch("public_border_group", az, nil),
				"category":            utils.PathSearch("category", az, nil),
			})
		}
		rst = append(rst, map[string]interface{}{
			"list": list,
		})
	}
	return rst
}

func buildGetAvailabilityZonesQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("public_border_group"); ok {
		res = fmt.Sprintf("%s&public_border_group=%v", res, v)
	}
	if v, ok := d.GetOk("loadbalancer_id"); ok {
		res = fmt.Sprintf("%s&loadbalancer_id=%v", res, v)
	}
	if res != "" {
		res = "?" + res[1:]
	}
	return res
}
