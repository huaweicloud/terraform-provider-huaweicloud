package secmaster

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

// @API Secmaster GET /v1/{project_id}/workspaces/{workspace_id}/soc/playbooks/versions/approval
func DataSourcePlaybookApprovals() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePlaybookApprovalsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"resource_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"approve_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"data": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"result": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"content": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourcePlaybookApprovalsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/playbooks/versions/approval"
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{workspace_id}", d.Get("workspace_id").(string))
	listPath = fmt.Sprintf("%s?resource_id=%s&approve_type=%s", listPath, d.Get("resource_id").(string), d.Get("approve_type").(string))

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	listResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return diag.Errorf("error retrieving playbook version approval: %s", err)
	}

	listRespBody, err := utils.FlattenResponse(listResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataList := utils.PathSearch("data", listRespBody, make([]interface{}, 0)).([]interface{})

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("data", flattenPlaybookApprovals(dataList)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenPlaybookApprovals(playbookApprovals []interface{}) []interface{} {
	if len(playbookApprovals) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(playbookApprovals))
	for _, v := range playbookApprovals {
		rst = append(rst, map[string]interface{}{
			"id":          utils.PathSearch("id", v, nil),
			"result":      utils.PathSearch("result", v, nil),
			"content":     utils.PathSearch("content", v, nil),
			"type":        utils.PathSearch("type", v, nil),
			"resource_id": utils.PathSearch("resource_id", v, nil),
			"user_id":     utils.PathSearch("user_id", v, nil),
			"create_time": utils.PathSearch("create_time", v, nil),
			"update_time": utils.PathSearch("update_time", v, nil),
		})
	}

	return rst
}
