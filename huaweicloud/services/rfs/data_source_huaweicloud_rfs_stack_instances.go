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

// @API RFS GET /v1/stack-sets/{stack_set_name}/stack-instances
func DataSourceStackInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceStackInstancesRead,

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
			"stack_instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"stack_set_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"stack_set_name": {
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
						"stack_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"stack_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"stack_domain_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"latest_stack_set_operation_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region": {
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

func buildListStackInstancesQueryParams(d *schema.ResourceData, marker string) string {
	rst := ""

	if v, ok := d.GetOk("stack_set_id"); ok {
		rst += fmt.Sprintf("&stack_set_id=%s", v.(string))
	}

	if v, ok := d.GetOk("filter"); ok {
		rst += fmt.Sprintf("&filter=%s", v.(string))
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

func dataSourceStackInstancesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg               = meta.(*config.Config)
		region            = cfg.GetRegion(d)
		httpUrl           = "v1/stack-sets/{stack_set_name}/stack-instances"
		allStackInstances = make([]interface{}, 0)
		nextMarker        string
	)

	client, err := cfg.NewServiceClient("rfs", region)
	if err != nil {
		return diag.Errorf("error creating RFS client: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{stack_set_name}", d.Get("stack_set_name").(string))
	requestOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Client-Request-Id": uuid,
		},
		KeepResponseBody: true,
	}

	for {
		requestPathWithQueryParams := requestPath + buildListStackInstancesQueryParams(d, nextMarker)
		resp, err := client.Request("GET", requestPathWithQueryParams, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving stack instances: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		stackInstances := utils.PathSearch("stack_instances", respBody, make([]interface{}, 0)).([]interface{})
		if len(stackInstances) == 0 {
			break
		}

		allStackInstances = append(allStackInstances, stackInstances...)

		nextMarker = utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if nextMarker == "" {
			break
		}
	}

	d.SetId(uuid)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("stack_instances", flattenStackInstances(allStackInstances)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenStackInstances(stackInstances []interface{}) []interface{} {
	if len(stackInstances) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(stackInstances))
	for _, v := range stackInstances {
		rst = append(rst, map[string]interface{}{
			"stack_set_id":                  utils.PathSearch("stack_set_id", v, nil),
			"stack_set_name":                utils.PathSearch("stack_set_name", v, nil),
			"status":                        utils.PathSearch("status", v, nil),
			"status_message":                utils.PathSearch("status_message", v, nil),
			"stack_id":                      utils.PathSearch("stack_id", v, nil),
			"stack_name":                    utils.PathSearch("stack_name", v, nil),
			"stack_domain_id":               utils.PathSearch("stack_domain_id", v, nil),
			"latest_stack_set_operation_id": utils.PathSearch("latest_stack_set_operation_id", v, nil),
			"region":                        utils.PathSearch("region", v, nil),
			"create_time":                   utils.PathSearch("create_time", v, nil),
			"update_time":                   utils.PathSearch("update_time", v, nil),
		})
	}
	return rst
}
