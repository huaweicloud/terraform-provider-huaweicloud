package rfs

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RFS GET /v1/stack-sets
func DataSourceRfsStackSets() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRfsStackSetsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
			"stack_sets": {
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
						"stack_set_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"permission_model": {
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
					},
				},
			},
		},
	}
}

func buildStackSetsQueryParams(d *schema.ResourceData, marker string) string {
	rst := ""

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

func dataSourceRfsStackSetsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		httpUrl    = "v1/stack-sets"
		allSets    = make([]interface{}, 0)
		nextMarker string
	)

	client, err := cfg.NewServiceClient("rfs", region)
	if err != nil {
		return diag.Errorf("error creating RFS client: %s", err)
	}

	uuid, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Client-Request-Id": uuid.String(),
			"Content-Type":      "application/json",
		},
		KeepResponseBody: true,
	}

	for {
		requestPathWithQueryParams := requestPath + buildStackSetsQueryParams(d, nextMarker)
		resp, err := client.Request("GET", requestPathWithQueryParams, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving RFS stack sets: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		sets := utils.PathSearch("stack_sets", respBody, make([]interface{}, 0)).([]interface{})
		if len(sets) == 0 {
			break
		}

		allSets = append(allSets, sets...)

		nextMarker = utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if nextMarker == "" {
			break
		}
	}

	d.SetId(uuid.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("stack_sets", flattenRfsStackSets(allSets)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenRfsStackSets(sets []interface{}) []interface{} {
	if len(sets) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(sets))
	for _, set := range sets {
		rst = append(rst, map[string]interface{}{
			"stack_set_id":          utils.PathSearch("stack_set_id", set, nil),
			"stack_set_name":        utils.PathSearch("stack_set_name", set, nil),
			"stack_set_description": utils.PathSearch("stack_set_description", set, nil),
			"permission_model":      utils.PathSearch("permission_model", set, nil),
			"status":                utils.PathSearch("status", set, nil),
			"create_time":           utils.PathSearch("create_time", set, nil),
			"update_time":           utils.PathSearch("update_time", set, nil),
		})
	}
	return rst
}
