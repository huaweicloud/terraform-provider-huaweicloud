package modelarts

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var (
	trainingExperimentNonUpdatableParams = []string{
		"metadata.*.workspace_id",
	}
	trainingExperimentNotFoundErrCodes = []string{
		"ModelArts.2842", // The resource does not exist.
	}
)

// @API ModelArts POST /v2/{project_id}/training-experiments
// @API ModelArts GET /v2/{project_id}/training-experiments/{experiment_id}
// @API ModelArts PUT /v2/{project_id}/training-experiments/{experiment_id}
// @API ModelArts DELETE /v2/{project_id}/training-experiments/{experiment_id}
func ResourceTrainingExperiment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTrainingExperimentCreate,
		ReadContext:   resourceTrainingExperimentRead,
		UpdateContext: resourceTrainingExperimentUpdate,
		DeleteContext: resourceTrainingExperimentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(trainingExperimentNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the training experiment is located.`,
			},

			// Required parameters.
			"metadata": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// Required parameters.
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The name of the training experiment.`,
						},

						// Optional parameters.
						"workspace_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `The ID of the workspace to which the training experiment belongs.`,
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The description of the training experiment.`,
						},
					},
				},
				Description: `The configuration of the training experiment.`,
			},

			// Attributes.
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the training experiment, in RFC3339 format.`,
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The latest update time of the training experiment, in RFC3339 format.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					`Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.`,
					utils.SchemaDescInput{Internal: true},
				),
			},
		},
	}
}

func buildTrainingExperimentCreateBodyParams(metadata []interface{}) map[string]interface{} {
	if len(metadata) < 1 {
		return nil
	}

	return map[string]interface{}{
		"metadata": utils.RemoveNil(map[string]interface{}{
			"name":         utils.PathSearch("name", metadata[0], nil),
			"description":  utils.ValueIgnoreEmpty(utils.PathSearch("description", metadata[0], nil)),
			"workspace_id": utils.ValueIgnoreEmpty(utils.PathSearch("workspace_id", metadata[0], nil)),
		}),
	}
}

func createTrainingExperiment(client *golangsdk.ServiceClient, params map[string]interface{}) (interface{}, error) {
	httpUrl := "v2/{project_id}/training-experiments"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(params),
	}

	requestResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(requestResp)
}

func resourceTrainingExperimentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("modelarts", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	resp, err := createTrainingExperiment(client, buildTrainingExperimentCreateBodyParams(d.Get("metadata").([]interface{})))
	if err != nil {
		return diag.Errorf("error creating training experiment: %s", err)
	}

	trainingExperimentId := utils.PathSearch("metadata.id", resp, "").(string)
	if trainingExperimentId == "" {
		return diag.Errorf("unable to find the ID of the training experiment from the API response")
	}

	d.SetId(trainingExperimentId)

	return resourceTrainingExperimentRead(ctx, d, meta)
}

func GetTrainingExperimentById(client *golangsdk.ServiceClient, trainingExperimentId string) (interface{}, error) {
	httpUrl := "v2/{project_id}/training-experiments/{experiment_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{experiment_id}", trainingExperimentId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(requestResp)
}

func resourceTrainingExperimentRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                  = meta.(*config.Config)
		region               = cfg.GetRegion(d)
		trainingExperimentId = d.Id()
	)

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	resp, err := GetTrainingExperimentById(client, trainingExperimentId)
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", trainingExperimentNotFoundErrCodes...),
			fmt.Sprintf("error retrieving training experiment (%s)", trainingExperimentId),
		)
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("metadata", flattenTrainingExperimentMetadata(resp)),
		// Attributes.
		d.Set("create_time", utils.FormatTimeStampRFC3339(int64(utils.PathSearch("metadata.create_time",
			resp, float64(0)).(float64))/1000, false)),
		d.Set("update_time", utils.FormatTimeStampRFC3339(int64(utils.PathSearch("metadata.update_time",
			resp, float64(0)).(float64))/1000, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenTrainingExperimentMetadata(resp interface{}) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"name":         utils.PathSearch("metadata.name", resp, nil),
			"description":  utils.PathSearch("metadata.description", resp, nil),
			"workspace_id": utils.PathSearch("metadata.workspace_id", resp, nil),
		},
	}
}

func buildTrainingExperimentUpdateBodyParams(metadata []interface{}) map[string]interface{} {
	if len(metadata) < 1 {
		return nil
	}

	return map[string]interface{}{
		"name":        utils.PathSearch("name", metadata[0], nil),
		"description": utils.PathSearch("description", metadata[0], nil),
	}
}

func updateTrainingExperiment(client *golangsdk.ServiceClient, trainingExperimentId string, params map[string]interface{}) error {
	httpUrl := "v2/{project_id}/training-experiments/{experiment_id}"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{experiment_id}", trainingExperimentId)

	opt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(params),
	}

	_, err := client.Request("PUT", updatePath, &opt)
	return err
}

func resourceTrainingExperimentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                  = meta.(*config.Config)
		trainingExperimentId = d.Id()
	)

	client, err := cfg.NewServiceClient("modelarts", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	err = updateTrainingExperiment(client, trainingExperimentId, buildTrainingExperimentUpdateBodyParams(d.Get("metadata").([]interface{})))
	if err != nil {
		return diag.Errorf("error updating training experiment (%s): %s", trainingExperimentId, err)
	}

	return resourceTrainingExperimentRead(ctx, d, meta)
}

func deleteTrainingExperiment(client *golangsdk.ServiceClient, trainingExperimentId string) error {
	httpUrl := "v2/{project_id}/training-experiments/{experiment_id}"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{experiment_id}", trainingExperimentId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err := client.Request("DELETE", deletePath, &opt)
	return err
}

func resourceTrainingExperimentDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                  = meta.(*config.Config)
		trainingExperimentId = d.Id()
	)

	client, err := cfg.NewServiceClient("modelarts", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	err = deleteTrainingExperiment(client, trainingExperimentId)
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", trainingExperimentNotFoundErrCodes...),
			fmt.Sprintf("error deleting training experiment (%s)", trainingExperimentId),
		)
	}

	return nil
}
