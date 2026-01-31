package lakeformation

import (
	"context"
	"fmt"
	"strconv"
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

// @API LakeFormation GET /v1/{project_id}/instances
func DataSourceInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceInstancesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the instances are located.`,
			},

			// Optional parameters.
			"in_recycle_bin": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to query instances in the recycle bin.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the enterprise project to which the instances belong.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the instance to be queried.`,
			},

			// Attributes.
			"instances": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of instances that matched filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the instance.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the instance.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the instance.`,
						},
						"enterprise_project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the enterprise project to which the instance belongs.`,
						},
						"shared": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the instance is shared.`,
						},
						"default_instance": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the instance is the default instance.`,
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation timestamp of the instance.`,
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The update timestamp of the instance.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The status of the instance.`,
						},
						"in_recycle_bin": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the instance is in the recycle bin.`,
						},
						"tags": common.TagsComputedSchema(`The key/value pairs to associate with the instance.`),
					},
				},
			},
		},
	}
}

func buildInstancesQueryParams(d *schema.ResourceData) string {
	var (
		enterpriseProjectId = "all_granted_eps"
		// in_recycle_bin is required by API, defaults to false.
		res = fmt.Sprintf("&in_recycle_bin=%v", d.Get("in_recycle_bin").(bool))
	)

	// enterprise_project_id is required by API, default is "all_granted_eps"
	if epsId, ok := d.GetOk("enterprise_project_id"); ok {
		enterpriseProjectId = epsId.(string)
	}
	res = fmt.Sprintf("%s&enterprise_project_id=%v", res, enterpriseProjectId)

	if name, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, name.(string))
	}

	return res
}

func listInstances(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/instances?limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))
	listPath += buildInstancesQueryParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
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
		instances := utils.PathSearch("instances", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, instances...)
		if len(instances) < limit {
			break
		}
		offset += len(instances)
	}

	return result, nil
}

func flattenInstances(items []interface{}) []map[string]interface{} {
	if len(items) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, map[string]interface{}{
			"instance_id":           utils.PathSearch("instance_id", item, nil),
			"name":                  utils.PathSearch("name", item, nil),
			"description":           utils.PathSearch("description", item, nil),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", item, nil),
			"shared":                utils.PathSearch("shared", item, nil),
			"default_instance":      utils.PathSearch("default_instance", item, nil),
			"create_time": utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("create_time",
				item, "").(string))/1000, false),
			"update_time": utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("update_time",
				item, "").(string))/1000, false),
			"status":         utils.PathSearch("status", item, nil),
			"in_recycle_bin": utils.PathSearch("in_recycle_bin", item, nil),
			"tags":           utils.FlattenTagsToMap(utils.PathSearch("tags", item, make([]interface{}, 0)).([]interface{})),
		})
	}

	return result
}

func dataSourceInstancesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("lakeformation", region)
	if err != nil {
		return diag.Errorf("error creating LakeFormation client: %s", err)
	}

	resp, err := listInstances(client, d)
	if err != nil {
		return diag.Errorf("error querying instances: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("instances", flattenInstances(resp)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
