package kafka

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

// @API Kafka GET /v2/{engine}/{project_id}/instances/{instance_id}/groups/{group}/members
func DataSourceConsumerGroupMembers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceConsumerGroupMembersRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the consumer group members are located.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the Kafka instance.`,
			},
			"group": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the consumer group.`,
			},
			"host": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The address of the consumer.`,
			},
			"member_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the consumer.`,
			},
			"members": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of consumer group members that match the filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the consumer.`,
						},
						"host": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The address of the consumer.`,
						},
						"client_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the client.`,
						},
					},
				},
			},
		},
	}
}

func buildConsumerGroupMembersQueryParams(d *schema.ResourceData) string {
	res := ""

	if host, ok := d.GetOk("host"); ok {
		res += fmt.Sprintf("&host=%s", host)
	}
	if memberId, ok := d.GetOk("member_id"); ok {
		res += fmt.Sprintf("&member_id=%s", memberId)
	}

	return res
}

func listConsumerGroupMembers(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		// The limit maximum value is 50, default is 10.
		httpUrl = "v2/kafka/{project_id}/instances/{instance_id}/groups/{group}/members?limit=50"
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", d.Get("instance_id").(string))
	listPath = strings.ReplaceAll(listPath, "{group}", d.Get("group").(string))
	listPath += buildConsumerGroupMembersQueryParams(d)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf-8",
		},
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		resp, err := client.Request("GET", listPathWithOffset, &getOpt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		members := utils.PathSearch("members", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, members...)
		// The `offset` cannot be greater than or equal to the total number of members. Otherwise, the response is as follows:
		// {"error_code": "DMS.00400062","error_msg": "Invalid {0} parameter in the request."}
		offset += len(members)
		if offset >= int(utils.PathSearch("total", respBody, float64(0)).(float64)) {
			break
		}
	}
	return result, nil
}

func dataSourceConsumerGroupMembersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("dmsv2", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	members, err := listConsumerGroupMembers(client, d)
	if err != nil {
		return diag.Errorf("error querying member list under consumer group (%s): %s", d.Get("group").(string), err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("members", flattenConsumerGroupMembers(members)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenConsumerGroupMembers(members []interface{}) []interface{} {
	if len(members) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(members))
	for _, v := range members {
		rst = append(rst, map[string]interface{}{
			"id":        utils.PathSearch("member_id", v, nil),
			"host":      utils.PathSearch("host", v, nil),
			"client_id": utils.PathSearch("client_id", v, nil),
		})
	}
	return rst
}
