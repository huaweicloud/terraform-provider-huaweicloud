package dew

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

// @API DEW GET /v1/{project_id}/dew/cpcs/instances
func DataSourceInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceInstancesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// This parameter does not take effect.
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_normal": {
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
			"result": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_normal": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"specification": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"az": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expired_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildInstancesQueryParams(d *schema.ResourceData) string {
	rst := "?page_size=10"

	if v, ok := d.GetOk("cluster_id"); ok {
		rst += fmt.Sprintf("&cluster_id=%v", v)
	}

	if v, ok := d.GetOk("name"); ok {
		rst += fmt.Sprintf("&name=%v", v)
	}

	if v, ok := d.GetOk("is_normal"); ok {
		rst += fmt.Sprintf("&is_normal=%v", v)
	}

	if v, ok := d.GetOk("sort_key"); ok {
		rst += fmt.Sprintf("&sort_key=%v", v)
	}

	if v, ok := d.GetOk("sort_dir"); ok {
		rst += fmt.Sprintf("&sort_dir=%v", v)
	}

	return rst
}

// The pagination does not take effect, use `total_num` as quantitative comparison.
func dataSourceInstancesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		httpUrl      = "v1/{project_id}/dew/cpcs/instances"
		pageNum      = 1
		allInstances = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient("kms", region)
	if err != nil {
		return diag.Errorf("error creating KMS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += buildInstancesQueryParams(d)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		getPathWithPageNum := fmt.Sprintf("%s&page_num=%d", getPath, pageNum)
		resp, err := client.Request("GET", getPathWithPageNum, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving CPCS instances: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		results := utils.PathSearch("result", respBody, make([]interface{}, 0)).([]interface{})
		if len(results) == 0 {
			break
		}

		allInstances = append(allInstances, results...)

		totalNum := int(utils.PathSearch("total_num", respBody, float64(0)).(float64))
		if len(allInstances) >= totalNum {
			break
		}

		pageNum++
	}

	datasourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("error generating UUID: %s", err)
	}

	d.SetId(datasourceId)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("result", flattenInstancesResponseBody(allInstances)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenInstancesResponseBody(results []interface{}) []interface{} {
	if len(results) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(results))
	for _, association := range results {
		result = append(result, map[string]interface{}{
			"instance_id":   utils.PathSearch("instance_id", association, nil),
			"resource_id":   utils.PathSearch("resource_id", association, nil),
			"instance_name": utils.PathSearch("instance_name", association, nil),
			"service_type":  utils.PathSearch("service_type", association, nil),
			"cluster_id":    utils.PathSearch("cluster_id", association, nil),
			"is_normal":     utils.PathSearch("is_normal", association, nil),
			"status":        utils.PathSearch("status", association, nil),
			"image_name":    utils.PathSearch("image_name", association, nil),
			"specification": utils.PathSearch("specification", association, nil),
			"az":            utils.PathSearch("az", association, nil),
			"expired_time":  utils.PathSearch("expired_time", association, nil),
			"create_time":   utils.PathSearch("create_time", association, nil),
		})
	}

	return result
}
