package rfs

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

// @API RFS GET /v1/{project_id}/stacks
func DataSourceStacks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceStacksRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"stacks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"stack_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"stack_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
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
						"status_message": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildListStacksQueryParams(marker string) string {
	if marker != "" {
		return fmt.Sprintf("?marker=%s", marker)
	}

	return ""
}

func dataSourceStacksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		httpUrl    = "v1/{project_id}/stacks"
		allStacks  = make([]interface{}, 0)
		nextMarker string
	)

	client, err := cfg.NewServiceClient("rfs", region)
	if err != nil {
		return diag.Errorf("error creating RFS client: %s", err)
	}

	reqUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate RFS request UUID: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Client-Request-Id": reqUUID},
		KeepResponseBody: true,
	}

	for {
		requestPathWithQueryParams := requestPath + buildListStacksQueryParams(nextMarker)
		resp, err := client.Request("GET", requestPathWithQueryParams, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving RFS stacks: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		stacks := utils.PathSearch("stacks", respBody, make([]interface{}, 0)).([]interface{})
		if len(stacks) == 0 {
			break
		}

		allStacks = append(allStacks, stacks...)

		nextMarker = utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if nextMarker == "" {
			break
		}
	}

	d.SetId(reqUUID)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("stacks", flattenStacks(allStacks)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenStacks(stacks []interface{}) []interface{} {
	if len(stacks) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(stacks))
	for _, v := range stacks {
		rst = append(rst, map[string]interface{}{
			"stack_name":     utils.PathSearch("stack_name", v, nil),
			"description":    utils.PathSearch("description", v, nil),
			"stack_id":       utils.PathSearch("stack_id", v, nil),
			"status":         utils.PathSearch("status", v, nil),
			"create_time":    utils.PathSearch("create_time", v, nil),
			"update_time":    utils.PathSearch("update_time", v, nil),
			"status_message": utils.PathSearch("status_message", v, nil),
		})
	}

	return rst
}
