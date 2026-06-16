package modelarts

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ModelArts GET /v2/{project_id}/datasets
func DataSourceDatasets() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDatasetsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the datasets are located.`,
			},

			// Optional parameters.
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the dataset to be queried.`,
			},
			"type": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The type of the dataset to be queried.`,
			},

			// Attributes.
			"datasets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the dataset.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the dataset.`,
						},
						"type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The type of the dataset.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the dataset.`,
						},
						"output_path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The OBS storage path that used to store output files.`,
						},
						"data_source": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"data_type": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `The type of the data source.`,
									},
									"path": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The OBS storage path or MRS HDFS path.`,
									},
									"cluster_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The ID of the DWS/MRS cluster.`,
									},
									"database_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The name of the DWS/DLI database.`,
									},
									"table_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The name of the DWS/DLI table.`,
									},
									"user_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The name of the DWS database user.`,
									},
									"queue_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The name of the DLI queue.`,
									},
									"with_column_header": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: `Whether the data contains table header when the type of dataset is 400 (table type).`,
									},
								},
							},
							Description: `The data sources which be used to imported the source data (such as pictures/files/audio,
etc.) in this directory and subdirectories to the dataset.`,
						},
						"schemas": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The field type of the schema.`,
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The field name of the schema.`,
									},
								},
							},
							Description: `The schema configurations of the dataset.`,
						},
						"labels": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The name of the label.`,
									},
									"property_color": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The color of the label.`,
									},
									"property_shape": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The shape of the label.`,
									},
									"property_shortcut": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The shortcut of the label.`,
									},
								},
							},
							Description: `The labels of the dataset.`,
						},
						"data_format": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The data format of the dataset.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the dataset, in RFC3339 format.`,
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The status of the dataset.`,
						},
					},
				},
				Description: `The list of datasets that match the filter parameters.`,
			},
		},
	}
}

func buildDatasetsQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("name"); ok {
		// Currently, `with_labels` can only be used when the query returns a small amount of data;
		// large amounts of data will cause gateway timeouts.
		res = fmt.Sprintf("%s&search_content=%v&with_labels=true", res, v)
	}
	if v, ok := d.GetOk("type"); ok {
		res = fmt.Sprintf("%s&type=%v", res, v)
	}
	return res
}

func listDatasets(client *golangsdk.ServiceClient, queryParams string) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/datasets?limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))
	if queryParams != "" {
		listPath += queryParams
	}

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := listPath + fmt.Sprintf("&offset=%d", offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &listOpt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		datasets := utils.PathSearch("datasets", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, datasets...)
		if len(datasets) < limit {
			break
		}
	}
	return result, nil
}

func flattenDatasets(datasets []interface{}) []map[string]interface{} {
	if len(datasets) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(datasets))
	for _, dataset := range datasets {
		result = append(result, map[string]interface{}{
			"id":          utils.PathSearch("dataset_id", dataset, nil),
			"name":        utils.PathSearch("dataset_name", dataset, nil),
			"type":        utils.PathSearch("dataset_type", dataset, nil),
			"description": utils.PathSearch("description", dataset, nil),
			"output_path": utils.PathSearch("work_path", dataset, nil),
			"created_at":  utils.FormatTimeStampUTC(int64(utils.PathSearch("create_time", dataset, float64(0)).(float64)) / 1000),
			"status":      utils.PathSearch("status", dataset, nil),
			"data_format": utils.PathSearch("data_format", dataset, nil),
			"data_source": flattenDatasetDataSource(utils.PathSearch("data_sources", dataset, make([]interface{}, 0)).([]interface{})),
			"schemas":     flattenDatasetSchema(utils.PathSearch("schema", dataset, make([]interface{}, 0)).([]interface{})),
			"labels":      flattenDatasetLabels(utils.PathSearch("labels", dataset, make([]interface{}, 0)).([]interface{})),
		})
	}
	return result
}

func dataSourceDatasetsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts v2 client, err=%s", err)
	}

	datasets, err := listDatasets(client, buildDatasetsQueryParams(d))
	if err != nil {
		return diag.Errorf("error querying ModelArts datasets: %s ", err)
	}

	randUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("datasets", flattenDatasets(datasets)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
