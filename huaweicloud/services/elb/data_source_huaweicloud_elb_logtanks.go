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

// @API ELB GET /v3/{project_id}/elb/logtanks
func DataSourceElbLogtanks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceElbLogtanksRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"logtank_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"log_topic_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"log_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"loadbalancer_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"logtanks": {
				Type:     schema.TypeList,
				Elem:     logtanksSchema(),
				Computed: true,
			},
		},
	}
}

func logtanksSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"loadbalancer_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"log_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"log_topic_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func dataSourceElbLogtanksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	var (
		listLogtanksHttpUrl = "v3/{project_id}/elb/logtanks"
		listLogtanksProduct = "elb"
	)
	listLogtanksClient, err := cfg.NewServiceClient(listLogtanksProduct, cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}
	listLogtanksPath := listLogtanksClient.Endpoint + listLogtanksHttpUrl
	listLogtanksPath = strings.ReplaceAll(listLogtanksPath, "{project_id}", listLogtanksClient.ProjectID)
	listLogtanksQueryParams := buildListLogtanksQueryParams(d, cfg.GetEnterpriseProjectID(d))
	listLogtanksPath += listLogtanksQueryParams

	listLogtanksResp, err := pagination.ListAllItems(
		listLogtanksClient,
		"marker",
		listLogtanksPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving logtanks")
	}

	listLogtanksRespJson, err := json.Marshal(listLogtanksResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var listLogtanksRespBody interface{}
	err = json.Unmarshal(listLogtanksRespJson, &listLogtanksRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(
		d.Set("region", cfg.GetRegion(d)),
		d.Set("logtanks", flattenListLogtanksBody(listLogtanksRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildListLogtanksQueryParams(d *schema.ResourceData, enterpriseProjectId string) string {
	res := ""
	if v, ok := d.GetOk("logtank_id"); ok {
		res = fmt.Sprintf("%s&id=%v", res, v)
	}
	if v, ok := d.GetOk("log_topic_id)"); ok {
		res = fmt.Sprintf("%s&log_topic_id=%v", res, v)
	}
	if v, ok := d.GetOk("log_group_id"); ok {
		res = fmt.Sprintf("%s&log_group_id=%v", res, v)
	}
	if v, ok := d.GetOk("loadbalancer_id"); ok {
		res = fmt.Sprintf("%s&loadbalancer_id=%v", res, v)
	}
	if enterpriseProjectId != "" {
		res = fmt.Sprintf("%s&enterprise_project_id=%v", res, enterpriseProjectId)
	}
	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func flattenListLogtanksBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("logtanks", resp, make([]interface{}, 0))
	if curJson == nil {
		return nil
	}
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":              utils.PathSearch("id", v, nil),
			"log_topic_id":    utils.PathSearch("log_topic_id", v, nil),
			"log_group_id":    utils.PathSearch("log_group_id", v, nil),
			"loadbalancer_id": utils.PathSearch("loadbalancer_id", v, nil),
		})
	}
	return rst
}
