package hss

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

// @API HSS GET /v5/{project_id}/files/change-files
func DataSourceChangeFiles() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceChangeFilesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"host_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// The parameters `begin_time` and `end_time` are of type int in the API documentation, but here they are
			// changed to type string. Because it needs to meet the scenario that can be set to `0`.
			"begin_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"file_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"file_path": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"change_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"change_category": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"file_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"file_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"change_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"change_category": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"after_change": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"before_change": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"latest_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildChangeFilesQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := "?limit=200"

	if v, ok := d.GetOk("host_name"); ok {
		queryParams = fmt.Sprintf("%s&host_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("begin_time"); ok {
		queryParams = fmt.Sprintf("%s&begin_time=%v", queryParams, v)
	}
	if v, ok := d.GetOk("end_time"); ok {
		queryParams = fmt.Sprintf("%s&end_time=%v", queryParams, v)
	}
	if v, ok := d.GetOk("file_name"); ok {
		queryParams = fmt.Sprintf("%s&file_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("file_path"); ok {
		queryParams = fmt.Sprintf("%s&file_path=%v", queryParams, v)
	}
	if v, ok := d.GetOk("change_type"); ok {
		queryParams = fmt.Sprintf("%s&change_type=%v", queryParams, v)
	}
	if v, ok := d.GetOk("change_category"); ok {
		queryParams = fmt.Sprintf("%s&change_category=%v", queryParams, v)
	}
	if v, ok := d.GetOk("status"); ok {
		queryParams = fmt.Sprintf("%s&status=%v", queryParams, v)
	}
	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}

	return queryParams
}

func dataSourceChangeFilesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		epsId   = cfg.GetEnterpriseProjectID(d)
		product = "hss"
		httpUrl = "v5/{project_id}/files/change-files"
		offset  = 0
		result  = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildChangeFilesQueryParams(d, epsId)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", requestPath, offset)
		resp, err := client.Request("GET", currentPath, &requestOpts)
		if err != nil {
			return diag.Errorf("error retrieving HSS change files: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		dataResp := utils.PathSearch("data_list", respBody, make([]interface{}, 0)).([]interface{})
		if len(dataResp) == 0 {
			break
		}

		result = append(result, dataResp...)
		offset += len(dataResp)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("data_list", flattenChangeFiles(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenChangeFiles(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"id":              utils.PathSearch("id", v, nil),
			"file_name":       utils.PathSearch("file_name", v, nil),
			"file_path":       utils.PathSearch("file_path", v, nil),
			"status":          utils.PathSearch("status", v, nil),
			"host_name":       utils.PathSearch("host_name", v, nil),
			"change_type":     utils.PathSearch("change_type", v, nil),
			"change_category": utils.PathSearch("change_category", v, nil),
			"after_change":    utils.PathSearch("after_change", v, nil),
			"before_change":   utils.PathSearch("before_change", v, nil),
			"latest_time":     utils.PathSearch("latest_time", v, nil),
		})
	}

	return rst
}
