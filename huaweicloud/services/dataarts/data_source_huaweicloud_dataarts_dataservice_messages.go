package dataarts

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

// @API DataArtsStudio GET /v1/{project_id}/service/messages
func DataSourceDataServiceMessages() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDataServiceMessagesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the approval messages are located.`,
			},

			// Parameter in request header
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The workspace ID of the exclusive API to which the approval message belongs.`,
			},

			// Query argument
			"api_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the API to be approved.`,
			},

			// Attribute
			"messages": {
				Type:        schema.TypeList,
				Elem:        dataserviceMessageSchema(),
				Computed:    true,
				Description: `All approval messages that match the filter parameters.`,
			},
		},
	}
}

func dataserviceMessageSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the approval message, in UUID format.`,
			},
			"api_apply_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The apply status for API.`,
			},
			"api_apply_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The apply type.`,
			},
			"api_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the exclusive API to which the approval message belongs.`,
			},
			"api_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the exclusive API to which the approval message belongs.`,
			},
			"api_using_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The expiration time used by the API, in RFC3339 format.`,
			},
			"app_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The application ID of the API that has been bound (or is to be bound).`,
			},
			"app_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The application name of the API that has been bound (or is to be bound).`,
			},
			"apply_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The apply time, in RFC3339 format.`,
			},
			"approval_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The approval time of the approval message, in RFC3339 format.`,
			},
			"approver_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The approver name.`,
			},
			"comment": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The approval comment.`,
			},
			"user_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of applicant.`,
			},
		},
	}
	return &sc
}

func buildDataServiceMessagesQueryParams(d *schema.ResourceData) string {
	res := ""
	if appName, ok := d.GetOk("api_name"); ok {
		res = fmt.Sprintf("%s&api_name=%v", res, appName)
	}
	return res
}

func queryDataServiceMessages(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/service/messages?limit=100"
		offset  = 0
		result  = make([]interface{}, 0)
	)
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath += buildDataServiceMessagesQueryParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"workspace":    d.Get("workspace_id").(string),
			"dlm-type":     "EXCLUSIVE",
		},
	}

	for {
		listPathWithOffset := listPath + fmt.Sprintf("&offset=%d", offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		messages := utils.PathSearch("records", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, messages...)
		offset += len(messages)
		// If the offset is greater than or equal to the total number, the first page of content will be queried instead
		// of returning an empty page (record list).
		if total := utils.PathSearch("total", respBody, float64(0)).(float64); offset >= int(total) {
			break
		}
	}

	return result, nil
}

func flattenDataServiceMessages(messages []interface{}) []interface{} {
	result := make([]interface{}, 0, len(messages))

	for _, message := range messages {
		result = append(result, map[string]interface{}{
			"id":               utils.PathSearch("id", message, nil),
			"api_apply_status": utils.PathSearch("api_apply_status", message, nil),
			"api_apply_type":   utils.PathSearch("api_apply_type", message, nil),
			"api_id":           utils.PathSearch("api_id", message, nil),
			"api_name":         utils.PathSearch("api_name", message, nil),
			"api_using_time":   utils.FormatTimeStampRFC3339(int64(utils.PathSearch("api_using_time", message, float64(0)).(float64))/1000, false),
			"app_id":           utils.PathSearch("app_id", message, nil),
			"app_name":         utils.PathSearch("app_name", message, nil),
			"apply_time":       utils.FormatTimeStampRFC3339(int64(utils.PathSearch("apply_time", message, float64(0)).(float64))/1000, false),
			"approval_time":    utils.FormatTimeStampRFC3339(int64(utils.PathSearch("approval_time", message, float64(0)).(float64))/1000, false),
			"approver_name":    utils.PathSearch("approver_name", message, nil),
			"comment":          utils.PathSearch("comment", message, nil),
			"user_name":        utils.PathSearch("user_name", message, nil),
		})
	}

	return result
}

func dataSourceDataServiceMessagesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	messages, err := queryDataServiceMessages(client, d)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("messages", flattenDataServiceMessages(messages)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
