package modelarts

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ModelArts GET /v2/{project_id}/training-experiments
func DataSourceTrainingExperiments() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTrainingExperimentsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the training experiments are located.`,
			},

			// Optional parameters.
			"workspace_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the workspace to which the training experiments belong.`,
			},
			"sort_by": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The field used for sorting the training experiments.`,
			},
			"order": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The sort order of the training experiments.`,
			},

			// Attributes.
			"training_experiments": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of training experiments that matched the filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"metadata": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        trainingExperimentsMetadataSchema(),
							Description: `The metadata of the training experiment.`,
						},
						"statistic": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"job_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `The number of training jobs under the training experiment.`,
									},
								},
							},
							Description: `The statistics of the training experiment.`,
						},
					},
				},
			},
		},
	}
}

func trainingExperimentsMetadataSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the training experiment.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the training experiment.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description of the training experiment.`,
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the workspace to which the training experiment belongs.`,
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the training experiment, in RFC3339 format.`,
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The update time of the training experiment, in RFC3339 format.`,
			},
		},
	}
}

func buildTrainingExperimentsQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("workspace_id"); ok {
		res = fmt.Sprintf("%s&workspace_id=%v", res, v)
	}
	if v, ok := d.GetOk("sort_by"); ok {
		res = fmt.Sprintf("%s&sort_by=%v", res, v)
	}
	if v, ok := d.GetOk("order"); ok {
		res = fmt.Sprintf("%s&order=%v", res, v)
	}

	return res
}

func listTrainingExperiments(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/training-experiments"
		// Maximum is 50.
		limit    = 50
		pageSize = 0
		result   = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath += fmt.Sprintf("?limit=%d%s", limit, buildTrainingExperimentsQueryParams(d))

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithPageSize := listPath + fmt.Sprintf("&offset=%d", pageSize)
		requestResp, err := client.Request("GET", listPathWithPageSize, &opt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		items := utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{})
		if len(items) < 1 {
			break
		}

		result = append(result, items...)
		if len(items) < limit {
			break
		}

		pageSize++
	}

	return result, nil
}

func dataSourceTrainingExperimentsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	experiments, err := listTrainingExperiments(client, d)
	if err != nil {
		return diag.Errorf("error querying training experiments: %s", err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("training_experiments", flattenTrainingExperiments(experiments)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenTrainingExperiments(experiments []interface{}) []map[string]interface{} {
	if len(experiments) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(experiments))
	for _, experiment := range experiments {
		result = append(result, map[string]interface{}{
			"metadata":  flattenTrainingExperimentsMetadata(utils.PathSearch("metadata", experiment, nil)),
			"statistic": flattenTrainingExperimentsStatistic(utils.PathSearch("statistic", experiment, nil)),
		})
	}

	return result
}

func flattenTrainingExperimentsMetadata(metadata interface{}) []map[string]interface{} {
	if metadata == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"id":           utils.PathSearch("id", metadata, nil),
			"name":         utils.PathSearch("name", metadata, nil),
			"description":  utils.PathSearch("description", metadata, nil),
			"workspace_id": utils.PathSearch("workspace_id", metadata, nil),
			"create_time": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("create_time",
				metadata, float64(0)).(float64))/1000, false),
			"update_time": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("update_time",
				metadata, float64(0)).(float64))/1000, false),
		},
	}
}

func flattenTrainingExperimentsStatistic(statistic interface{}) []map[string]interface{} {
	if statistic == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"job_count": utils.PathSearch("job_count", statistic, nil),
		},
	}
}
