package modelarts

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ModelArts GET /v2/{project_id}/training-job-engines
func DataSourceTrainingJobEngines() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTrainingJobEnginesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the training job engines are located.`,
			},
			"engines": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"engine_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the engine.`,
						},
						"engine_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the engine.`,
						},
						"engine_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The version of the engine.`,
						},
						"v1_compatible": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the engine is v1 compatible.`,
						},
						"run_user": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The default startup user UID of the engine.`,
						},
						"image_info": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cpu_image_url": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The CPU image URL of the engine.`,
									},
									"gpu_image_url": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The GPU or Ascend image URL of the engine.`,
									},
									"image_version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The image version of the engine.`,
									},
								},
							},
							Description: `The image information of the engine.`,
						},
					},
				},
				Description: `The list of training job engines.`,
			},
		},
	}
}

func listTrainingJobEngines(client *golangsdk.ServiceClient) ([]interface{}, error) {
	httpUrl := "v2/{project_id}/training-job-engines"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", listPath, &opt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{}), nil
}

func dataSourceTrainingJobEnginesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	engines, err := listTrainingJobEngines(client)
	if err != nil {
		return diag.Errorf("error querying training job engines: %s", err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("engines", flattenTrainingJobEngines(engines)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenTrainingJobEngines(engines []interface{}) []map[string]interface{} {
	if len(engines) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(engines))
	for _, engine := range engines {
		result = append(result, map[string]interface{}{
			"engine_id":      utils.PathSearch("engine_id", engine, nil),
			"engine_name":    utils.PathSearch("engine_name", engine, nil),
			"engine_version": utils.PathSearch("engine_version", engine, nil),
			"v1_compatible":  utils.PathSearch("v1_compatible", engine, nil),
			"run_user":       utils.PathSearch("run_user", engine, nil),
			"image_info": flattenTrainingJobEngineImageInfo(utils.PathSearch("image_info",
				engine, make(map[string]interface{})).(map[string]interface{})),
		})
	}

	return result
}

func flattenTrainingJobEngineImageInfo(imageInfo map[string]interface{}) []map[string]interface{} {
	if len(imageInfo) < 1 {
		return nil
	}

	return []map[string]interface{}{
		{
			"cpu_image_url": utils.PathSearch("cpu_image_url", imageInfo, nil),
			"gpu_image_url": utils.PathSearch("gpu_image_url", imageInfo, nil),
			"image_version": utils.PathSearch("image_version", imageInfo, nil),
		},
	}
}
