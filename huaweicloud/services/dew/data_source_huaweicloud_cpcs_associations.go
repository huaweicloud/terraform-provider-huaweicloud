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

// @API DEW GET /v1/{project_id}/dew/cpcs/associations
func DataSourceAssociations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAssociationsRead,

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
			"app_id": {
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
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"app_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"app_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subnet_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_server_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpcep_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_time": {
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

func buildAssociationsQueryParams(d *schema.ResourceData) string {
	rst := "?page_size=10"

	if v, ok := d.GetOk("cluster_id"); ok {
		rst += fmt.Sprintf("&cluster_id=%v", v)
	}

	if v, ok := d.GetOk("app_id"); ok {
		rst += fmt.Sprintf("&app_id=%v", v)
	}

	if v, ok := d.GetOk("sort_key"); ok {
		rst += fmt.Sprintf("&sort_key=%v", v)
	}

	if v, ok := d.GetOk("sort_dir"); ok {
		rst += fmt.Sprintf("&sort_dir=%v", v)
	}

	return rst
}

// The first page of page_num is `1`.
func dataSourceAssociationsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg             = meta.(*config.Config)
		region          = cfg.GetRegion(d)
		httpUrl         = "v1/{project_id}/dew/cpcs/associations"
		pageNum         = 1
		allAssociations = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient("kms", region)
	if err != nil {
		return diag.Errorf("error creating KMS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += buildAssociationsQueryParams(d)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		getPathWithPageNum := fmt.Sprintf("%s&page_num=%d", getPath, pageNum)
		resp, err := client.Request("GET", getPathWithPageNum, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving CPCS associations: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		results := utils.PathSearch("result", respBody, make([]interface{}, 0)).([]interface{})
		if len(results) == 0 {
			break
		}

		allAssociations = append(allAssociations, results...)

		pageNum++
	}

	datasourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("error generating UUID: %s", err)
	}

	d.SetId(datasourceId)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("result", flattenAssociationsResponseBody(allAssociations)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAssociationsResponseBody(results []interface{}) []interface{} {
	if len(results) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(results))
	for _, association := range results {
		result = append(result, map[string]interface{}{
			"id":                  utils.PathSearch("id", association, nil),
			"cluster_id":          utils.PathSearch("cluster_id", association, nil),
			"cluster_name":        utils.PathSearch("cluster_name", association, nil),
			"app_id":              utils.PathSearch("app_id", association, nil),
			"app_name":            utils.PathSearch("app_name", association, nil),
			"vpc_name":            utils.PathSearch("vpc_name", association, nil),
			"subnet_name":         utils.PathSearch("subnet_name", association, nil),
			"cluster_server_type": utils.PathSearch("cluster_server_type", association, nil),
			"vpcep_address":       utils.PathSearch("address", association, nil),
			"update_time":         utils.PathSearch("update_time", association, nil),
			"create_time":         utils.PathSearch("create_time", association, nil),
		})
	}

	return result
}
