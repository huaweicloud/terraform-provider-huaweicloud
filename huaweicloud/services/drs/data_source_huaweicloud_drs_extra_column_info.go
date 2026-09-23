package drs

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DRS GET /v5/{project_id}/jobs/{job_id}/extra-column-info
func DataSourceDrsExtraColumnInfo() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDrsExtraColumnInfoRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the region in which to query the data source. If omitted, the provider-level region will be used.",
			},
			"job_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the ID of the DRS job to query extra column info.",
			},
			"is_only_show_sent": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Specifies whether to query only the processed objects that have been sent.",
			},
			"fetch_all": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Specifies whether to query all data. When set to true, offset and limit parameters are ignored.",
			},
			"column_process_objects": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of column process objects.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"object_alias_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The alias name of the object after mapping.",
						},
						"object_source_names": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "The source object names of the selected source database.",
						},
						"is_sent": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the extra column has been sent.",
						},
						"extra_column_infos": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The extra column information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"column_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of the column.",
									},
									"column_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type of the column.",
									},
									"column_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The value of the column.",
									},
									"data_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The data type of the column.",
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

func dataSourceDrsExtraColumnInfoRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		httpUrl = "v5/{project_id}/jobs/{job_id}/extra-column-info"
		product = "drs"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{job_id}", d.Get("job_id").(string))
	listPath += buildListExtraColumnInfoQueryParams(d)

	listResp, err := pagination.ListAllItems(
		client,
		"offset",
		listPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return diag.Errorf("error retrieving DRS extra column info: %s", err)
	}

	listRespJson, err := json.Marshal(listResp)
	if err != nil {
		return diag.FromErr(err)
	}

	var listRespBody interface{}
	err = json.Unmarshal(listRespJson, &listRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(
		d.Set("column_process_objects", flattenExtraColumnProcessObjects(listRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildListExtraColumnInfoQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("is_only_show_sent"); ok {
		res = fmt.Sprintf("%s&is_only_show_sent=%v", res, v.(bool))
	}
	if v, ok := d.GetOk("fetch_all"); ok {
		res = fmt.Sprintf("%s&fetch_all=%v", res, v.(bool))
	}

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func flattenExtraColumnProcessObjects(resp interface{}) []interface{} {
	curJson := utils.PathSearch("column_process_objects", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"object_alias_name":   utils.PathSearch("object_alias_name", v, nil),
			"object_source_names": utils.PathSearch("object_source_names", v, make([]interface{}, 0)),
			"is_sent":             utils.PathSearch("is_sent", v, nil),
			"extra_column_infos":  flattenExtraColumnInfos(v),
		})
	}
	return res
}

func flattenExtraColumnInfos(obj interface{}) []interface{} {
	curJson := utils.PathSearch("extra_column_infos", obj, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"column_name":  utils.PathSearch("column_name", v, nil),
			"column_type":  utils.PathSearch("column_type", v, nil),
			"column_value": utils.PathSearch("column_value", v, nil),
			"data_type":    utils.PathSearch("data_type", v, nil),
		})
	}
	return res
}
