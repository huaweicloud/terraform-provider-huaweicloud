package modelarts

import (
	"context"

	"github.com/chnsz/golangsdk/openstack/modelarts/v2/dataset"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func DataSourceDatasets() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDatasetsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"type": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{0, 1, 3, 100, 101, 102, 200, 201, 202, 400, 600, 900}),
			},

			"datasets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"type": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"output_path": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"data_source": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"data_type": {
										Type:     schema.TypeInt,
										Computed: true,
									},

									"path": {
										Type:     schema.TypeString,
										Computed: true,
									},

									"with_column_header": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},

						"schemas": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},

									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},

						"labels": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},

									"property_color": {
										Type:     schema.TypeString,
										Computed: true,
									},

									"property_shape": {
										Type:     schema.TypeString,
										Computed: true,
									},

									"property_shortcut": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},

						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"status": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"data_format": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceDatasetsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	client, err := config.ModelArtsV2Client(region)
	if err != nil {
		return fmtp.DiagErrorf("Error creating ModelArts v2 client, err=%s", err)
	}

	opts := dataset.ListOpts{
		WithLabels:    true,
		SearchContent: d.Get("name").(string),
	}

	if v, ok := d.GetOk("type"); ok {
		dataType := v.(int)
		opts.DatasetType = &dataType
	}

	page, err := dataset.List(client, opts)
	if err != nil {
		return fmtp.DiagErrorf("Error querying ModelArts datasets: %s ", err)
	}

	p, err := page.AllPages()
	if err != nil {
		return fmtp.DiagErrorf("Error querying ModelArts datasets: %s", err)
	}

	datasets, err := dataset.ExtractDatasets(p)
	if err != nil {
		return fmtp.DiagErrorf("Error querying ModelArts datasets: %s", err)
	}

	var rst []map[string]interface{}
	var ids []string
	for _, v := range datasets {
		item := map[string]interface{}{
			"id":          v.DatasetId,
			"name":        v.DatasetName,
			"type":        v.DatasetType,
			"description": v.Description,
			"output_path": v.WorkPath,
			"created_at":  utils.FormatTimeStampUTC(int64(v.CreateTime) / 1000),
			"status":      v.Status,
			"data_format": v.DataFormat,
			"data_source": parseDataSource(v.DataSources),
			"schemas":     parseSchemas(v.Schema),
			"labels":      parseLabels(v.Labels),
		}
		rst = append(rst, item)
		ids = append(ids, v.DatasetId)
	}

	err = d.Set("datasets", rst)
	if err != nil {
		return fmtp.DiagErrorf("set datasets err:%s", err)
	}

	d.SetId(hashcode.Strings(ids))
	return nil
}

func parseLabels(in []dataset.Label) []interface{} {
	var result []interface{}
	for _, v := range in {
		item := map[string]interface{}{
			"name":              v.Name,
			"property_color":    v.Property.Color,
			"property_shape":    v.Property.DefaultShape,
			"property_shortcut": v.Property.Shortcut,
		}

		result = append(result, item)
	}
	return result
}

func parseDataSource(in []dataset.DataSource) []interface{} {
	result := make([]interface{}, 1)
	if len(in) >= 1 {
		result[0] = map[string]interface{}{
			"data_type":          in[0].DataType,
			"path":               in[0].DataPath,
			"with_column_header": in[0].WithColumnHeader,
		}
	}
	return result
}

func parseSchemas(in []dataset.Field) []interface{} {
	result := make([]interface{}, len(in))
	for i, v := range in {
		result[i] = map[string]interface{}{
			"name": v.Name,
			"type": v.Type,
		}
	}
	return result
}
