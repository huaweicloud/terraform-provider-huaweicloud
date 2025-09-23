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

// @API ELB GET /v3/{project_id}/elb/pools/{pool_id}/members
func DataSourceElbMembers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceElbMembersRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"pool_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"protocol_port": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"weight": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"member_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"operating_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"member_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"members": {
				Type:     schema.TypeList,
				Elem:     membersSchema(),
				Computed: true,
			},
		},
	}
}

func membersSchema() *schema.Resource {
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
			"weight": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"protocol_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"member_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"operating_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"reason": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     memberReasonSchema(),
			},
			"status": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     memberStatusSchema(),
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

func memberReasonSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"expected_response": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"healthcheck_response": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"reason_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func memberStatusSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"listener_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"operating_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"reason": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     memberReasonSchema(),
			},
		},
	}
}

func dataSourceElbMembersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	var (
		listMembersHttpUrl = "v3/{project_id}/elb/pools/{pool_id}/members"
		listMembersProduct = "elb"
	)
	listMembersClient, err := cfg.NewServiceClient(listMembersProduct, cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}
	listMembersPath := listMembersClient.Endpoint + listMembersHttpUrl
	listMembersPath = strings.ReplaceAll(listMembersPath, "{project_id}", listMembersClient.ProjectID)
	listMembersPath = strings.ReplaceAll(listMembersPath, "{pool_id}", d.Get("pool_id").(string))
	listMembersQueryParams := buildListMembersQueryParams(d, cfg.GetEnterpriseProjectID(d))
	listMembersPath += listMembersQueryParams

	listMembersResp, err := pagination.ListAllItems(
		listMembersClient,
		"marker",
		listMembersPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Members")
	}

	listMembersRespJson, err := json.Marshal(listMembersResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var listMembersRespBody interface{}
	err = json.Unmarshal(listMembersRespJson, &listMembersRespBody)
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
		d.Set("members", flattenListMembersBody(listMembersRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildListMembersQueryParams(d *schema.ResourceData, enterpriseProjectId string) string {
	res := ""
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	if v, ok := d.GetOk("address"); ok {
		res = fmt.Sprintf("%s&address=%v", res, v)
	}
	if v, ok := d.GetOk("protocol_port"); ok {
		res = fmt.Sprintf("%s&protocol_port=%v", res, v)
	}
	if v, ok := d.GetOk("weight"); ok {
		res = fmt.Sprintf("%s&weight=%v", res, v)
	}
	if v, ok := d.GetOk("subnet_id"); ok {
		res = fmt.Sprintf("%s&subnet_cidr_id=%v", res, v)
	}
	if enterpriseProjectId != "" {
		res = fmt.Sprintf("%s&enterprise_project_id=%v", res, enterpriseProjectId)
	}
	if v, ok := d.GetOk("member_id"); ok {
		res = fmt.Sprintf("%s&id=%v", res, v)
	}
	if v, ok := d.GetOk("member_type"); ok {
		res = fmt.Sprintf("%s&member_type=%v", res, v)
	}
	if v, ok := d.GetOk("instance_id"); ok {
		res = fmt.Sprintf("%s&instance_id=%v", res, v)
	}
	if v, ok := d.GetOk("ip_version"); ok {
		res = fmt.Sprintf("%s&ip_version=%v", res, v)
	}
	if v, ok := d.GetOk("operating_status"); ok {
		res = fmt.Sprintf("%s&operating_status=%v", res, v)
	}
	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func flattenListMembersBody(resp interface{}) []interface{} {
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
			"id":               utils.PathSearch("id", v, nil),
			"name":             utils.PathSearch("name", v, nil),
			"weight":           utils.PathSearch("weight", v, nil),
			"subnet_id":        utils.PathSearch("subnet_cidr_id", v, nil),
			"address":          utils.PathSearch("address", v, nil),
			"protocol_port":    utils.PathSearch("protocol_port", v, nil),
			"member_type":      utils.PathSearch("member_type", v, nil),
			"instance_id":      utils.PathSearch("instance_id", v, nil),
			"ip_version":       utils.PathSearch("ip_version", v, nil),
			"operating_status": utils.PathSearch("operating_status", v, nil),
			"reason":           flattenMemberReason(utils.PathSearch("reason", v, nil)),
			"status":           flattenMemberStatus(utils.PathSearch("status", v, nil)),
			"created_at":       utils.PathSearch("created_at", v, nil),
			"updated_at":       utils.PathSearch("updated_at", v, nil),
		})
	}
	return rst
}
