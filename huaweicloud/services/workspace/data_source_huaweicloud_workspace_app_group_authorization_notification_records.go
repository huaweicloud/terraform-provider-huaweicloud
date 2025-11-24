package workspace

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

// @API Workspace GET /v1/{project_id}/mails
func DataSourceAppGroupAuthorizationNotificationRecords() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppGroupAuthorizationNotificationRecordsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the authorization notification records are located.`,
			},
			"app_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the application group.`,
			},
			"account": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the authorized user (group).`,
			},
			"mail_send_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The type of authorization operation.`,
			},
			"mail_send_result": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The result of the notification sending.`,
			},
			"records": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the record.`,
						},
						"account": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the user (group).`,
						},
						"account_auth_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the account.`,
						},
						"account_auth_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the authorized object.`,
						},
						"app_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the application group.`,
						},
						"app_group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the application group.`,
						},
						"mail_send_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the authorization.`,
						},
						"mail_send_result": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The result of the notification sending.`,
						},
						"error_msg": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The error message when the notification failed to be sent.`,
						},
						"send_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The time when the authorization notification was sent.`,
						},
					},
				},
				Description: `The authorization notification record list that match filter parameters.`,
			},
		},
	}
}

func buildAppGroupAuthorizationNotificationRecordsQueryParams(d *schema.ResourceData) string {
	res := fmt.Sprintf("&app_group_id=%s", d.Get("app_group_id").(string))

	if v, ok := d.GetOk("account"); ok {
		res = fmt.Sprintf("%s&account=%s", res, v)
	}

	if v, ok := d.GetOk("mail_send_type"); ok {
		res = fmt.Sprintf("%s&mail_send_type=%s", res, v)
	}

	if v, ok := d.GetOk("mail_send_result"); ok {
		res = fmt.Sprintf("%s&mail_send_result=%s", res, v)
	}

	return res
}

func listAppGroupAuthorizationNotificationRecords(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/mails"
		// The `limit` default value is 10, maximum value is 100.
		limit  = 100
		offset = 0
		result = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = fmt.Sprintf("%s?limit=%d", listPath, limit)
	listPath += buildAppGroupAuthorizationNotificationRecordsQueryParams(d)

	requestOpts := &golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, requestOpts)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		records := utils.PathSearch("records", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, records...)
		if len(records) < limit {
			break
		}

		offset += len(records)
	}

	return result, nil
}

func flattenAuthorizationNotificationRecords(records []interface{}) []map[string]interface{} {
	if len(records) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(records))
	for _, item := range records {
		result = append(result, map[string]interface{}{
			"id":                utils.PathSearch("id", item, nil),
			"account":           utils.PathSearch("account", item, nil),
			"account_auth_type": utils.PathSearch("account_auth_type", item, nil),
			"account_auth_name": utils.PathSearch("account_auth_name", item, nil),
			"app_group_id":      utils.PathSearch("app_group_id", item, nil),
			"app_group_name":    utils.PathSearch("app_group_name", item, nil),
			"mail_send_type":    utils.PathSearch("mail_send_type", item, nil),
			"mail_send_result":  utils.PathSearch("mail_send_result", item, nil),
			"error_msg":         utils.PathSearch("error_msg", item, nil),
			"send_at":           utils.PathSearch("send_at", item, nil),
		})
	}

	return result
}

func dataSourceAppGroupAuthorizationNotificationRecordsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	records, err := listAppGroupAuthorizationNotificationRecords(client, d)
	if err != nil {
		return diag.Errorf("error querying application group authorization notification records: %s", err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("records", flattenAuthorizationNotificationRecords(records)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
