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

// @API ELB GET /v3/{project_id}/elb/l7policies/{l7policy_id}/rules
func DataSourceElbL7rules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceElbL7rulesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"l7policy_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"l7rule_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"compare_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"value": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"l7rules": {
				Type:     schema.TypeList,
				Elem:     l7rulesSchema(),
				Computed: true,
			},
		},
	}
}

func l7rulesSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"compare_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"value": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"conditions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
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

func dataSourceElbL7rulesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		listL7rulesHttpUrl = "v3/{project_id}/elb/l7policies/{l7policy_id}/rules"
		listL7rulesProduct = "elb"
	)
	listL7rulesClient, err := cfg.NewServiceClient(listL7rulesProduct, region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}
	listL7rulesPath := listL7rulesClient.Endpoint + listL7rulesHttpUrl
	listL7rulesPath = strings.ReplaceAll(listL7rulesPath, "{project_id}", listL7rulesClient.ProjectID)
	listL7rulesPath = strings.ReplaceAll(listL7rulesPath, "{l7policy_id}", d.Get("l7policy_id").(string))
	listL7rulesQueryParams := buildListL7rulesQueryParams(d, cfg.GetEnterpriseProjectID(d))
	listL7rulesPath += listL7rulesQueryParams
	listL7rulesResp, err := pagination.ListAllItems(
		listL7rulesClient,
		"marker",
		listL7rulesPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving ELB L7 rules")
	}

	listL7rulesRespJson, err := json.Marshal(listL7rulesResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var listL7rulesRespBody interface{}
	err = json.Unmarshal(listL7rulesRespJson, &listL7rulesRespBody)
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
		d.Set("l7rules", flattenListL7rulesBody(listL7rulesRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildListL7rulesQueryParams(d *schema.ResourceData, enterpriseProjectId string) string {
	res := ""
	if v, ok := d.GetOk("l7rule_id"); ok {
		res = fmt.Sprintf("%s&id=%v", res, v)
	}
	if v, ok := d.GetOk("compare_type"); ok {
		res = fmt.Sprintf("%s&compare_type=%v", res, v)
	}
	if v, ok := d.GetOk("value"); ok {
		res = fmt.Sprintf("%s&value=%v", res, v)
	}
	if v, ok := d.GetOk("type"); ok {
		res = fmt.Sprintf("%s&type=%v", res, v)
	}
	if enterpriseProjectId != "" {
		res = fmt.Sprintf("%s&enterprise_project_id=%v", res, enterpriseProjectId)
	}
	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func flattenListL7rulesBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("rules", resp, make([]interface{}, 0))
	if curJson == nil {
		return nil
	}
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":           utils.PathSearch("id", v, nil),
			"compare_type": utils.PathSearch("compare_type", v, nil),
			"type":         utils.PathSearch("type", v, nil),
			"value":        utils.PathSearch("value", v, nil),
			"conditions":   utils.PathSearch("conditions", v, nil),
			"created_at":   utils.PathSearch("created_at", v, nil),
			"updated_at":   utils.PathSearch("updated_at", v, nil),
		})
	}
	return rst
}
