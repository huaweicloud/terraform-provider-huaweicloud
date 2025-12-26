package secmaster

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

// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/siem/dataspaces
func DataSourceDataspaces() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDataspacesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dataspace_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dataspace_name": {
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
			"records": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dataspace_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dataspace_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dataspace_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_by": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_by": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildDataspacesQueryParams(d *schema.ResourceData) string {
	queryParams := "?limit=100"

	if v, ok := d.GetOk("dataspace_id"); ok {
		queryParams = fmt.Sprintf("%s&dataspace_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("dataspace_name"); ok {
		queryParams = fmt.Sprintf("%s&dataspace_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("sort_key"); ok {
		queryParams = fmt.Sprintf("%s&sort_key=%v", queryParams, v)
	}
	if v, ok := d.GetOk("sort_dir"); ok {
		queryParams = fmt.Sprintf("%s&sort_dir=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceDataspacesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/siem/dataspaces"
		result      = make([]interface{}, 0)
		offset      = 0
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{workspace_id}", workspaceId)
	getPath += buildDataspacesQueryParams(d)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", getPath, offset)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving dataspaces: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}

		dataspaces := utils.PathSearch("records", getRespBody, make([]interface{}, 0)).([]interface{})
		if len(dataspaces) == 0 {
			break
		}

		result = append(result, dataspaces...)
		offset += len(dataspaces)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("records", flattenDataspaces(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDataspaces(dataResp []interface{}) []interface{} {
	if len(dataResp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(dataResp))
	for _, v := range dataResp {
		rst = append(rst, map[string]interface{}{
			"dataspace_id":   utils.PathSearch("dataspace_id", v, nil),
			"dataspace_name": utils.PathSearch("dataspace_name", v, nil),
			"dataspace_type": utils.PathSearch("dataspace_type", v, nil),
			"description":    utils.PathSearch("description", v, nil),
			"domain_id":      utils.PathSearch("domain_id", v, nil),
			"project_id":     utils.PathSearch("project_id", v, nil),
			"region_id":      utils.PathSearch("region_id", v, nil),
			"create_by":      utils.PathSearch("create_by", v, nil),
			"update_by":      utils.PathSearch("update_by", v, nil),
			"create_time":    utils.PathSearch("create_time", v, nil),
			"update_time":    utils.PathSearch("update_time", v, nil),
		})
	}

	return rst
}
