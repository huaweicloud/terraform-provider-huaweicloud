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

// @API ELB GET /v3/{project_id}/elb/healthmonitors
func DataSourceElbMonitors() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceElbMonitorsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"monitor_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"pool_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"domain_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"interval": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"max_retries": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"max_retries_down": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"timeout": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status_code": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"http_method": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"url_path": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"monitors": {
				Type:     schema.TypeList,
				Elem:     monitorsSchema(),
				Computed: true,
			},
		},
	}
}

func monitorsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"interval": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"domain_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"http_method": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"max_retries": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"max_retries_down": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pool_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"timeout": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"url_path": {
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
		},
	}
	return &sc
}

func dataSourceElbMonitorsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		listMonitorsHttpUrl = "v3/{project_id}/elb/healthmonitors"
		listMonitorsProduct = "elb"
	)
	listMonitorsClient, err := cfg.NewServiceClient(listMonitorsProduct, region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}
	listMonitorsPath := listMonitorsClient.Endpoint + listMonitorsHttpUrl
	listMonitorsPath = strings.ReplaceAll(listMonitorsPath, "{project_id}", listMonitorsClient.ProjectID)
	listMonitorsQueryParams := buildListMonitorsQueryParams(d, cfg.GetEnterpriseProjectID(d))
	listMonitorsPath += listMonitorsQueryParams
	listMonitorsResp, err := pagination.ListAllItems(
		listMonitorsClient,
		"marker",
		listMonitorsPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving ELB health monitors")
	}

	listMonitorsRespJson, err := json.Marshal(listMonitorsResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var listMonitorsRespBody interface{}
	err = json.Unmarshal(listMonitorsRespJson, &listMonitorsRespBody)
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
		d.Set("monitors", flattenListMonitorsBody(listMonitorsRespBody, d)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildListMonitorsQueryParams(d *schema.ResourceData, enterpriseProjectId string) string {
	res := ""
	if v, ok := d.GetOk("monitor_id"); ok {
		res = fmt.Sprintf("%s&id=%v", res, v)
	}
	if v, ok := d.GetOk("monitor_port"); ok {
		res = fmt.Sprintf("%s&port=%v", res, v)
	}
	if v, ok := d.GetOk("domain_name"); ok {
		res = fmt.Sprintf("%s&domain_name=%v", res, v)
	}
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	if v, ok := d.GetOk("interval"); ok {
		res = fmt.Sprintf("%s&delay=%v", res, v)
	}
	if v, ok := d.GetOk("max_retries"); ok {
		res = fmt.Sprintf("%s&max_retries=%v", res, v)
	}
	if v, ok := d.GetOk("max_retries_down"); ok {
		res = fmt.Sprintf("%s&max_retries_down=%v", res, v)
	}
	if v, ok := d.GetOk("timeout"); ok {
		res = fmt.Sprintf("%s&timeout=%v", res, v)
	}
	if v, ok := d.GetOk("protocol"); ok {
		res = fmt.Sprintf("%s&type=%v", res, v)
	}
	if v, ok := d.GetOk("status_code"); ok {
		res = fmt.Sprintf("%s&expected_codes=%v", res, v)
	}
	if v, ok := d.GetOk("http_method"); ok {
		res = fmt.Sprintf("%s&http_method=%v", res, v)
	}
	if v, ok := d.GetOk("url_path"); ok {
		res = fmt.Sprintf("%s&url_path=%v", res, v)
	}
	if enterpriseProjectId != "" {
		res = fmt.Sprintf("%s&enterprise_project_id=%v", res, enterpriseProjectId)
	}
	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func flattenListMonitorsBody(resp interface{}, d *schema.ResourceData) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("healthmonitors", resp, make([]interface{}, 0))
	if curJson == nil {
		return nil
	}
	rawPoolId, rawPoolIdOk := d.GetOk("pool_id")
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		poolId := utils.PathSearch("pools|[0].id", v, "")
		if rawPoolIdOk && rawPoolId.(string) != poolId.(string) {
			continue
		}
		rst = append(rst, map[string]interface{}{
			"id":               utils.PathSearch("id", v, nil),
			"interval":         utils.PathSearch("delay", v, nil),
			"domain_name":      utils.PathSearch("domain_name", v, nil),
			"status_code":      utils.PathSearch("expected_codes", v, nil),
			"http_method":      utils.PathSearch("http_method", v, nil),
			"max_retries":      utils.PathSearch("max_retries", v, nil),
			"max_retries_down": utils.PathSearch("max_retries_down", v, nil),
			"port":             utils.PathSearch("monitor_port", v, nil),
			"name":             utils.PathSearch("name", v, nil),
			"pool_id":          poolId,
			"timeout":          utils.PathSearch("timeout", v, nil),
			"protocol":         utils.PathSearch("type", v, nil),
			"url_path":         utils.PathSearch("url_path", v, nil),
			"created_at":       utils.PathSearch("created_at", v, nil),
			"updated_at":       utils.PathSearch("updated_at", v, nil),
		})
	}
	return rst
}
