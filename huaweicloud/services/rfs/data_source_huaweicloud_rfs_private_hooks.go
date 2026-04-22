package rfs

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RFS GET /v1/private-hooks
func DataSourceRfsPrivateHooks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRfsPrivateHooksRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"sort_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sort_dir": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"hooks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hook_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hook_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hook_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"default_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"configuration": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"target_stacks": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"failure_mode": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
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

func buildPrivateHooksQueryParams(d *schema.ResourceData, marker string) string {
	rst := ""

	if v, ok := d.GetOk("sort_key"); ok {
		rst += fmt.Sprintf("&sort_key=%s", v.(string))
	}

	if v, ok := d.GetOk("sort_dir"); ok {
		rst += fmt.Sprintf("&sort_dir=%s", v.(string))
	}

	if marker != "" {
		rst += fmt.Sprintf("&marker=%s", marker)
	}

	if rst != "" {
		rst = "?" + rst[1:]
	}

	return rst
}

func dataSourceRfsPrivateHooksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		httpUrl    = "v1/private-hooks"
		allHooks   = make([]interface{}, 0)
		nextMarker string
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
	requestOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Client-Request-Id": uuid,
			"Content-Type":      "application/json",
		},
		KeepResponseBody: true,
	}

	for {
		requestPathWithQueryParams := requestPath + buildPrivateHooksQueryParams(d, nextMarker)
		resp, err := client.Request("GET", requestPathWithQueryParams, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving RFS private hooks: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		hooks := utils.PathSearch("hooks", respBody, make([]interface{}, 0)).([]interface{})
		if len(hooks) == 0 {
			break
		}

		allHooks = append(allHooks, hooks...)

		nextMarker = utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if nextMarker == "" {
			break
		}
	}

	d.SetId(uuid)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("hooks", flattenPrivateHooks(allHooks)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenPrivateHooks(hooks []interface{}) []interface{} {
	if len(hooks) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(hooks))
	for _, hook := range hooks {
		rst = append(rst, map[string]interface{}{
			"hook_id":          utils.PathSearch("hook_id", hook, nil),
			"hook_name":        utils.PathSearch("hook_name", hook, nil),
			"hook_description": utils.PathSearch("hook_description", hook, nil),
			"default_version":  utils.PathSearch("default_version", hook, nil),
			"configuration":    flattenRfsHookConfiguration(utils.PathSearch("configuration", hook, nil)),
			"create_time":      utils.PathSearch("create_time", hook, nil),
			"update_time":      utils.PathSearch("update_time", hook, nil),
		})
	}
	return rst
}

func flattenRfsHookConfiguration(configuration interface{}) []interface{} {
	if configuration == nil {
		return nil
	}

	result := []interface{}{
		map[string]interface{}{
			"target_stacks": utils.PathSearch("target_stacks", configuration, nil),
			"failure_mode":  utils.PathSearch("failure_mode", configuration, nil),
		},
	}

	return result
}
