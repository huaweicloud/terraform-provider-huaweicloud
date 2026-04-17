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

// @API RFS GET /v1/{project_id}/stacks/{stack_name}/resources
func DataSourceStackResources() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceStackResourcesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"stack_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"stack_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"stack_resources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"physical_resource_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"physical_resource_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"logical_resource_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"logical_resource_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"index_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status_message": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_attributes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeString,
										Computed: true,
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

func buildListStackResourcesQueryParams(d *schema.ResourceData, marker string) string {
	rst := ""

	if v, ok := d.GetOk("stack_id"); ok {
		rst += fmt.Sprintf("&stack_id=%s", v.(string))
	}

	if marker != "" {
		rst += fmt.Sprintf("&marker=%s", marker)
	}

	if rst != "" {
		rst = "?" + rst[1:]
	}

	return rst
}

func dataSourceStackResourcesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg               = meta.(*config.Config)
		region            = cfg.GetRegion(d)
		httpUrl           = "v1/{project_id}/stacks/{stack_name}/resources"
		allStackResources = make([]interface{}, 0)
		nextMarker        string
	)

	client, err := cfg.NewServiceClient("rfs", region)
	if err != nil {
		return diag.Errorf("error creating RFS client: %s", err)
	}

	reqUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	stackName := d.Get("stack_name").(string)
	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{stack_name}", stackName)
	requestOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Client-Request-Id": reqUUID},
		KeepResponseBody: true,
	}

	for {
		requestPathWithQueryParams := requestPath + buildListStackResourcesQueryParams(d, nextMarker)
		resp, err := client.Request("GET", requestPathWithQueryParams, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving RFS stack resources: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		stackResources := utils.PathSearch("stack_resources", respBody, make([]interface{}, 0)).([]interface{})
		if len(stackResources) == 0 {
			break
		}

		allStackResources = append(allStackResources, stackResources...)

		nextMarker = utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if nextMarker == "" {
			break
		}
	}

	d.SetId(reqUUID)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("stack_resources", flattenStackResources(allStackResources)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenStackResources(stackResources []interface{}) []interface{} {
	if len(stackResources) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(stackResources))
	for _, v := range stackResources {
		rst = append(rst, map[string]interface{}{
			"physical_resource_id":   utils.PathSearch("physical_resource_id", v, nil),
			"physical_resource_name": utils.PathSearch("physical_resource_name", v, nil),
			"logical_resource_name":  utils.PathSearch("logical_resource_name", v, nil),
			"logical_resource_type":  utils.PathSearch("logical_resource_type", v, nil),
			"index_key":              utils.PathSearch("index_key", v, nil),
			"resource_status":        utils.PathSearch("resource_status", v, nil),
			"status_message":         utils.PathSearch("status_message", v, nil),
			"resource_attributes":    flattenStackResourceAttributes(v),
		})
	}
	return rst
}

func flattenStackResourceAttributes(v interface{}) []interface{} {
	attributesResp := utils.PathSearch("resource_attributes", v, make([]interface{}, 0)).([]interface{})
	if len(attributesResp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(attributesResp))
	for _, v := range attributesResp {
		rst = append(rst, map[string]interface{}{
			"key":   utils.PathSearch("key", v, nil),
			"value": utils.PathSearch("value", v, nil),
		})
	}

	return rst
}
