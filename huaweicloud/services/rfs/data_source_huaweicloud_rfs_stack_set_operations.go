package rfs

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RFS GET /v1/stack-sets/{stack_set_name}/operations
func DataSourceStackSetOperations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceStackSetOperationsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"stack_set_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"stack_set_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sort_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sort_dir": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"call_identity": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"stack_set_operations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"operation_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"stack_set_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"stack_set_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"action": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status_message": {
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

func buildListStackSetOperationsQueryParams(d *schema.ResourceData, marker string) string {
	rst := ""

	if v, ok := d.GetOk("stack_set_id"); ok {
		rst += fmt.Sprintf("&stack_set_id=%s", v.(string))
	}

	if v, ok := d.GetOk("filter"); ok {
		rst += fmt.Sprintf("&filter=%s", url.QueryEscape(v.(string)))
	}

	if v, ok := d.GetOk("sort_key"); ok {
		rst += fmt.Sprintf("&sort_key=%s", v.(string))
	}

	if v, ok := d.GetOk("sort_dir"); ok {
		rst += fmt.Sprintf("&sort_dir=%s", v.(string))
	}

	if v, ok := d.GetOk("call_identity"); ok {
		rst += fmt.Sprintf("&call_identity=%s", v.(string))
	}

	if marker != "" {
		rst += fmt.Sprintf("&marker=%s", marker)
	}

	if rst != "" {
		rst = "?" + rst[1:]
	}

	return rst
}

func dataSourceStackSetOperationsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                   = meta.(*config.Config)
		region                = cfg.GetRegion(d)
		httpUrl               = "v1/stack-sets/{stack_set_name}/operations"
		allStackSetOperations = make([]interface{}, 0)
		nextMarker            string
		stackSetName          = d.Get("stack_set_name").(string)
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
	requestPath = strings.ReplaceAll(requestPath, "{stack_set_name}", stackSetName)
	requestOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Client-Request-Id": reqUUID},
		KeepResponseBody: true,
	}

	for {
		requestPathWithQueryParams := requestPath + buildListStackSetOperationsQueryParams(d, nextMarker)
		resp, err := client.Request("GET", requestPathWithQueryParams, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving RFS stack set operations: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		stackSetOperations := utils.PathSearch("stack_set_operations", respBody, make([]interface{}, 0)).([]interface{})
		if len(stackSetOperations) == 0 {
			break
		}

		allStackSetOperations = append(allStackSetOperations, stackSetOperations...)

		nextMarker = utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if nextMarker == "" {
			break
		}
	}

	d.SetId(reqUUID)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("stack_set_operations", flattenStackSetOperations(allStackSetOperations)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenStackSetOperations(stackSetOperations []interface{}) []interface{} {
	if len(stackSetOperations) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(stackSetOperations))
	for _, v := range stackSetOperations {
		rst = append(rst, map[string]interface{}{
			"operation_id":   utils.PathSearch("operation_id", v, nil),
			"stack_set_id":   utils.PathSearch("stack_set_id", v, nil),
			"stack_set_name": utils.PathSearch("stack_set_name", v, nil),
			"action":         utils.PathSearch("action", v, nil),
			"status":         utils.PathSearch("status", v, nil),
			"status_message": utils.PathSearch("status_message", v, nil),
			"create_time":    utils.PathSearch("create_time", v, nil),
			"update_time":    utils.PathSearch("update_time", v, nil),
		})
	}

	return rst
}
