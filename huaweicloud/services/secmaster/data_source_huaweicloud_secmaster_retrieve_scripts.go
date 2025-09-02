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

// @API SecMaster GET /v2/{project_id}/workspaces/{workspace_id}/siem/retrieve-scripts
func DataSourceRetrieveScripts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRetrieveScriptsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the workspace ID.",
			},
			"table_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the table ID.",
			},
			"script_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the script name.",
			},
			"sort_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the attribute fields for sorting.",
			},
			"sort_dir": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the sorting order. Supported values are **ASC** and **DESC**.",
			},
			"records": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The retrieve scripts list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"script_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The script ID.",
						},
						"project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The project ID.",
						},
						"workspace_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The workspace ID.",
						},
						"script_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The script name.",
						},
						"table_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The table ID.",
						},
						"category": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The script category.",
						},
						"directory": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The script directory group name.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The script description.",
						},
						"script": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The script content.",
						},
						"create_by": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IAM user ID who created the script.",
						},
						"create_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The creation time (timestamp in milliseconds).",
						},
						"update_by": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IAM user ID who updated the script.",
						},
						"update_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The update time (timestamp in milliseconds).",
						},
					},
				},
			},
		},
	}
}

func buildRetrieveScriptsQueryParams(d *schema.ResourceData) string {
	queryParams := "?limit=200"

	if v, ok := d.GetOk("table_id"); ok {
		queryParams = fmt.Sprintf("%s&table_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("script_name"); ok {
		queryParams = fmt.Sprintf("%s&script_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("sort_key"); ok {
		queryParams = fmt.Sprintf("%s&sort_key=%v", queryParams, v)
	}
	if v, ok := d.GetOk("sort_dir"); ok {
		queryParams = fmt.Sprintf("%s&sort_dir=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceRetrieveScriptsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "secmaster"
		httpUrl = "v2/{project_id}/workspaces/{workspace_id}/siem/retrieve-scripts"
		offset  = 0
		result  = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", d.Get("workspace_id").(string))
	requestPath += buildRetrieveScriptsQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", requestPath, offset)
		requestResp, err := client.Request("GET", currentPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving SecMaster retrieve scripts: %s", err)
		}

		requestRespBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return diag.FromErr(err)
		}

		recordsResp := utils.PathSearch("records", requestRespBody, make([]interface{}, 0)).([]interface{})
		if len(recordsResp) == 0 {
			break
		}

		result = append(result, recordsResp...)
		offset += len(recordsResp)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("records", flattenRetrieveScriptsRecords(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenRetrieveScriptsRecords(recordsResp []interface{}) []interface{} {
	if len(recordsResp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(recordsResp))
	for _, v := range recordsResp {
		rst = append(rst, map[string]interface{}{
			"script_id":    utils.PathSearch("script_id", v, nil),
			"project_id":   utils.PathSearch("project_id", v, nil),
			"workspace_id": utils.PathSearch("workspace_id", v, nil),
			"script_name":  utils.PathSearch("script_name", v, nil),
			"table_id":     utils.PathSearch("table_id", v, nil),
			"category":     utils.PathSearch("category", v, nil),
			"directory":    utils.PathSearch("directory", v, nil),
			"description":  utils.PathSearch("description", v, nil),
			"script":       utils.PathSearch("script", v, nil),
			"create_by":    utils.PathSearch("create_by", v, nil),
			"create_time":  utils.PathSearch("create_time", v, nil),
			"update_by":    utils.PathSearch("update_by", v, nil),
			"update_time":  utils.PathSearch("update_time", v, nil),
		})
	}

	return rst
}
