package lts

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API LTS GET /v1/{project_id}/lts/{member_account_id}/all-streams
func DataSourceMemberGroupStreams() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMemberGroupStreamsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region in which to query the member group streams.`,
			},
			"member_account_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the member account.`,
			},
			"groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of the log groups.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"log_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the log group.`,
						},
						"log_group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the log group.`,
						},
						"log_streams": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The list of log streams.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"log_stream_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The ID of the log stream.`,
									},
									"log_stream_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The name of the log stream.`,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceMemberGroupStreamsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("lts", region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	httpUrl := "v1/{project_id}/lts/{member_account_id}/all-streams"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{member_account_id}", d.Get("member_account_id").(string))

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", listPath, &opt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving member group streams")
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("groups", flattenGroups(utils.PathSearch("results", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGroups(groups []interface{}) []map[string]interface{} {
	if len(groups) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(groups))
	for _, group := range groups {
		result = append(result, map[string]interface{}{
			"log_group_id":   utils.PathSearch("log_group_id", group, nil),
			"log_group_name": utils.PathSearch("log_group_name", group, nil),
			"log_streams": flattenLogStreams(
				utils.PathSearch("log_streams", group, make([]interface{}, 0)).([]interface{})),
		})
	}

	return result
}

func flattenLogStreams(streams []interface{}) []map[string]interface{} {
	if len(streams) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(streams))
	for _, stream := range streams {
		result = append(result, map[string]interface{}{
			"log_stream_id":   utils.PathSearch("log_stream_id", stream, nil),
			"log_stream_name": utils.PathSearch("log_stream_name", stream, nil),
		})
	}

	return result
}
