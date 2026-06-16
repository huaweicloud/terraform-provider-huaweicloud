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

// @API DataArtsStudio GET /v1/{project_id}/service/authorizeapply/approver
func DataSourceDataServiceApprovers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDataServiceApproversRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the approvers are located.`,
			},

			// Required parameters.
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the workspace to which the approvers belong.`,
			},

			// Optional parameters.
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the approver to be queried.`,
			},

			// Attributes.
			"approvers": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataServiceApproverElem(),
				Description: `The list of approvers that match the filter parameters.`,
			},
		},
	}
}

func dataServiceApproverElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the approver.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the approver.`,
			},
			"user_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The user ID of the approver.`,
			},
			"user_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The user name of the approver.`,
			},
			"email": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: `The email of the approver.`,
			},
			"phone_number": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: `The phone number of the approver.`,
			},
			"department": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The department of the approver.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description of the approver.`,
			},
			"character": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The character of the approver.`,
			},
			"create_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creator of the approver.`,
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the approver, in RFC3339 format.`,
			},
			"app_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The application name.`,
			},
			"topic_urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The topic URN.`,
			},
			"template_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The template ID.`,
			},
			"project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The project ID.`,
			},
			"approver_type": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The type of the approver.`,
			},
		},
	}
}

func buildDataServiceApproversQueryParams(d *schema.ResourceData) string {
	res := ""
	if approverName, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&approver_name=%v", res, approverName)
	}
	return res
}

func listDataServiceApprovers(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/service/authorizeapply/approver?size={size}"
		page    = 1
		size    = 100
		result  = make([]interface{}, 0)
	)

	listPathWithSize := client.Endpoint + httpUrl
	listPathWithSize = strings.ReplaceAll(listPathWithSize, "{project_id}", client.ProjectID)
	listPathWithSize = strings.ReplaceAll(listPathWithSize, "{size}", strconv.Itoa(size))
	listPathWithSize += buildDataServiceApproversQueryParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"workspace":    d.Get("workspace_id").(string),
			"dlm-type":     "EXCLUSIVE",
		},
	}

	for {
		listPathWithPage := listPathWithSize + fmt.Sprintf("&page=%s", strconv.Itoa(page))
		requestResp, err := client.Request("GET", listPathWithPage, &opt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		approvers := utils.PathSearch("approvers", respBody, make([]interface{}, 0)).([]interface{})
		if len(approvers) < 1 {
			break
		}
		result = append(result, approvers...)
		page++
	}

	return result, nil
}

func flattenDataServiceApprovers(approvers []interface{}) []interface{} {
	if len(approvers) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(approvers))
	for _, approver := range approvers {
		result = append(result, map[string]interface{}{
			"id":            utils.PathSearch("id", approver, nil),
			"name":          utils.PathSearch("approver_name", approver, nil),
			"user_id":       utils.PathSearch("user_id", approver, nil),
			"user_name":     utils.PathSearch("user_name", approver, nil),
			"email":         utils.PathSearch("email", approver, nil),
			"phone_number":  utils.PathSearch("phone_number", approver, nil),
			"department":    utils.PathSearch("department", approver, nil),
			"description":   utils.PathSearch("description", approver, nil),
			"character":     utils.PathSearch("character", approver, nil),
			"create_by":     utils.PathSearch("create_by", approver, nil),
			"create_time":   utils.FormatTimeStampRFC3339(int64(utils.PathSearch("create_time", approver, float64(0)).(float64))/1000, false),
			"app_name":      utils.PathSearch("app_name", approver, nil),
			"topic_urn":     utils.PathSearch("topic_urn", approver, nil),
			"template_id":   utils.PathSearch("template_id", approver, nil),
			"project_id":    utils.PathSearch("project_id", approver, nil),
			"approver_type": utils.PathSearch("approver_type", approver, nil),
		})
	}

	return result
}

func dataSourceDataServiceApproversRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	approvers, err := listDataServiceApprovers(client, d)
	if err != nil {
		return diag.Errorf("error querying Data Service approvers: %s", err)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("approvers", flattenDataServiceApprovers(approvers)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
