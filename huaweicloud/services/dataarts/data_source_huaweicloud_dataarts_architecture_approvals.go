package dataarts

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DataArtsStudio GET /v2/{project_id}/design/approvals
func DataSourceArchitectureApprovals() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceArchitectureApprovalsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The region where the approvals are located.",
			},

			// Required parameters.
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the workspace to which the approvals belong.",
			},

			// Optional parameters.
			"biz_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The business ID of the approvals.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the approvals.",
			},
			"create_by": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The creator of the approvals.",
			},
			"approver": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The approver of the approvals.",
			},
			"approval_status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The approval status of the approvals.",
			},

			// Attributes.
			"approvals": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataArchitectureApprovalsElem(),
				Description: "The list of approvals that matched filter parameters.",
			},
		},
	}
}

func dataArchitectureApprovalsElem() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the approval, in UUID format.`,
			},
			"name_ch": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The chinese name of the approval.",
			},
			"name_en": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The english name of the approval.",
			},
			"biz_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The business ID of the approval.",
			},
			"biz_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The business type of the approval.",
			},
			"biz_info": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The business details of the approval.",
			},
			"biz_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The business status of the approval.",
			},
			"approval_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The approval status of the approval.",
			},
			"approval_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The approval type of the approval.",
			},
			"submit_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The submit time of the approval, in RFC3339 format.",
			},
			"create_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creator of the approval.",
			},
			"approval_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The approval time of the approval, in RFC3339 format.",
			},
			"approver": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The approver of the approval.",
			},
			"email": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The approver email of the approval.",
			},
			"msg": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The approval message of the approval.",
			},
		},
	}
	return &sc
}

func buildArchitectureApprovalsQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("biz_id"); ok {
		res = fmt.Sprintf("%s&biz_id=%v", res, v)
	}
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	if v, ok := d.GetOk("create_by"); ok {
		res = fmt.Sprintf("%s&create_by=%v", res, v)
	}
	if v, ok := d.GetOk("approver"); ok {
		res = fmt.Sprintf("%s&approver=%v", res, v)
	}
	if v, ok := d.GetOk("approval_status"); ok {
		res = fmt.Sprintf("%s&approval_status=%v", res, v)
	}

	return res
}

func listArchitectureApprovals(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/design/approvals?limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))
	listPath += buildArchitectureApprovalsQueryParams(d)

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildArchitectureMoreHeaders(d.Get("workspace_id").(string)),
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &listOpt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		approvals := utils.PathSearch("data.value.records", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, approvals...)

		if len(approvals) < limit {
			break
		}
		offset += len(approvals)
	}

	return result, nil
}

func flattenArchitectureApprovals(approvals []interface{}) []map[string]interface{} {
	if len(approvals) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(approvals))
	for _, approval := range approvals {
		result = append(result, map[string]interface{}{
			"id":              utils.PathSearch("id", approval, nil),
			"name_ch":         utils.PathSearch("name_ch", approval, nil),
			"name_en":         utils.PathSearch("name_en", approval, nil),
			"biz_id":          utils.PathSearch("biz_id", approval, nil),
			"biz_type":        utils.PathSearch("biz_type", approval, nil),
			"biz_info":        utils.PathSearch("biz_info", approval, nil),
			"biz_status":      utils.PathSearch("biz_status", approval, nil),
			"approval_status": utils.PathSearch("approval_status", approval, nil),
			"approval_type":   utils.PathSearch("approval_type", approval, nil),
			"submit_time":     utils.PathSearch("submit_time", approval, nil),
			"create_by":       utils.PathSearch("create_by", approval, nil),
			"approval_time":   utils.PathSearch("approval_time", approval, nil),
			"approver":        utils.PathSearch("approver", approval, nil),
			"email":           utils.PathSearch("email", approval, nil),
			"msg":             utils.PathSearch("msg", approval, nil),
		})
	}
	return result
}

func dataSourceArchitectureApprovalsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts client: %s", err)
	}

	approvals, err := listArchitectureApprovals(client, d)
	if err != nil {
		return diag.Errorf("error querying DataArts Architecture approvals: %s", err)
	}

	randUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("approvals", flattenArchitectureApprovals(approvals)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
