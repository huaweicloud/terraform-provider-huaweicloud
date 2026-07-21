package dsc

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
	maskAlgorithmNonUpdatableParams = []string{
		"category",
	}
	maskAlgorithmNotFoundErrCodes = []string{
		"dsc.10000009", // The DSC instance does not exist.
	}
)

// @API DSC POST /v1/{project_id}/sdg/server/mask/algorithms
// @API DSC GET /v2/{project_id}/mask/algorithms
// @API DSC PUT /v1/{project_id}/sdg/server/mask/algorithms/{algorithm_id}
// @API DSC DELETE /v1/{project_id}/sdg/server/mask/algorithms/{algorithm_id}
func ResourceMaskAlgorithm() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMaskAlgorithmCreate,
		ReadContext:   resourceMaskAlgorithmRead,
		UpdateContext: resourceMaskAlgorithmUpdate,
		DeleteContext: resourceMaskAlgorithmDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(maskAlgorithmNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the mask algorithm is located.`,
			},

			// Required parameters.
			"algorithm_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the mask algorithm.`,
			},
			"algorithm": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The encryption mask algorithm type.`,
			},
			"algorithm_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The type of the mask algorithm.`,
			},
			"category": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The category of the mask algorithm.`,
			},
			"parameter": {
				Type:     schema.TypeString,
				Required: true,
				// The order of fields in the JSON strings returned by the creation and query is inconsistent, we need to ignore the differences.
				DiffSuppressFunc: func(_, o, n string, _ *schema.ResourceData) bool {
					equal, _ := utils.CompareJsonTemplateAreEquivalent(o, n)
					return equal
				},
				Description: `The configuration parameters of the mask algorithm, in JSON format.`,
			},

			// Optional parameters.
			"data": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: `The data content processed by the mask algorithm.`,
			},
			"processed_data": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The data content after masking.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					`Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
		},
	}
}

func buildMaskAlgorithmBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		// Required parameters.
		"algorithm_name": d.Get("algorithm_name"),
		"algorithm":      d.Get("algorithm"),
		"algorithm_type": d.Get("algorithm_type"),
		"category":       d.Get("category"),
		"parameter":      d.Get("parameter").(string),
		// Optional parameters.
		"data":           utils.ValueIgnoreEmpty(d.Get("data")),
		"processed_data": utils.ValueIgnoreEmpty(d.Get("processed_data")),
	}
}

func createMaskAlgorithm(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	httpUrl := "v1/{project_id}/sdg/server/mask/algorithms"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildMaskAlgorithmBodyParams(d)),
		MoreHeaders: map[string]string{
			"content-type": "application/json;charset=UTF-8",
		},
	}

	_, err := client.Request("POST", createPath, &createOpt)
	return err
}

func resourceMaskAlgorithmCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg  = meta.(*config.Config)
		name = d.Get("algorithm_name").(string)
	)
	client, err := cfg.NewServiceClient("dsc", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	err = createMaskAlgorithm(client, d)
	if err != nil {
		return diag.Errorf("error creating mask algorithm: %s", err)
	}

	// The create API does not return the algorithm ID, query the list to get it.
	algorithms, err := listMaskAlgorithms(client)
	if err != nil {
		return diag.Errorf("error getting mask algorithm (%s): %s", name, err)
	}

	algorithmId := utils.PathSearch(fmt.Sprintf("[?algorithm_name=='%s']|[0].algorithm_id", name), algorithms, "").(string)
	if algorithmId == "" {
		return diag.Errorf("unable to find the mask algorithm ID from API response")
	}

	d.SetId(algorithmId)

	return resourceMaskAlgorithmRead(ctx, d, meta)
}

func listMaskAlgorithms(client *golangsdk.ServiceClient) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/mask/algorithms"
		limit   = 200
		// The offset is the page number.
		offset        = 1
		allAlgorithms = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = fmt.Sprintf("%s?limit=%v", listPath, limit)

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"content-type": "application/json;charset=UTF-8",
		},
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		resp, err := client.Request("GET", listPathWithOffset, &listOpt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		algorithms := utils.PathSearch("algorithms", respBody, make([]interface{}, 0)).([]interface{})
		allAlgorithms = append(allAlgorithms, algorithms...)
		if len(algorithms) < limit {
			break
		}

		offset++
	}

	return allAlgorithms, nil
}

func GetMaskAlgorithmById(client *golangsdk.ServiceClient, algorithmId string) (interface{}, error) {
	algorithms, err := listMaskAlgorithms(client)
	if err != nil {
		return nil, err
	}

	algorithm := utils.PathSearch(fmt.Sprintf("[?algorithm_id=='%s']|[0]", algorithmId), algorithms, nil)
	if algorithm == nil {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v2/{project_id}/mask/algorithms",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("the mask algorithm (%s) does not exist", algorithmId)),
			},
		}
	}

	return algorithm, nil
}

func resourceMaskAlgorithmRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		algorithmId = d.Id()
	)

	client, err := cfg.NewServiceClient("dsc", region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	algorithm, err := GetMaskAlgorithmById(client, algorithmId)
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected401ErrInto404Err(err, "error_code", maskAlgorithmNotFoundErrCodes...),
			fmt.Sprintf("error retrieving mask algorithm (%s)", algorithmId),
		)
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		// Required parameters.
		d.Set("algorithm_name", utils.PathSearch("algorithm_name", algorithm, nil)),
		d.Set("algorithm", utils.PathSearch("algorithm", algorithm, nil)),
		d.Set("category", utils.PathSearch("category", algorithm, nil)),
		d.Set("parameter", utils.PathSearch("parameter", algorithm, nil)),
		// Optional parameters.
		d.Set("data", utils.PathSearch("data", algorithm, nil)),
		d.Set("processed_data", utils.PathSearch("processed_data", algorithm, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func updateMaskAlgorithm(client *golangsdk.ServiceClient, algorithmId string, d *schema.ResourceData) error {
	httpUrl := "v1/{project_id}/sdg/server/mask/algorithms/{algorithm_id}"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{algorithm_id}", algorithmId)

	bodyParams := buildMaskAlgorithmBodyParams(d)
	bodyParams["algorithm_id"] = algorithmId

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(bodyParams),
		MoreHeaders: map[string]string{
			"content-type": "application/json;charset=UTF-8",
		},
	}

	_, err := client.Request("PUT", updatePath, &updateOpt)
	return err
}

func resourceMaskAlgorithmUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("dsc", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	if d.HasChangesExcept("enable_force_new") {
		err = updateMaskAlgorithm(client, d.Id(), d)
		if err != nil {
			return diag.Errorf("error updating mask algorithm (%v): %s", d.Get("algorithm_name"), err)
		}
	}

	return resourceMaskAlgorithmRead(ctx, d, meta)
}

func deleteMaskAlgorithm(client *golangsdk.ServiceClient, algorithmId string) error {
	httpUrl := "v1/{project_id}/sdg/server/mask/algorithms/{algorithm_id}"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{algorithm_id}", algorithmId)
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"content-type": "application/json;charset=UTF-8",
		},
	}

	_, err := client.Request("DELETE", deletePath, &deleteOpt)
	return err
}

func resourceMaskAlgorithmDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("dsc", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	err = deleteMaskAlgorithm(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected401ErrInto404Err(err, "error_code", maskAlgorithmNotFoundErrCodes...),
			fmt.Sprintf("error deleting mask algorithm (%v)", d.Get("algorithm_name")),
		)
	}

	return nil
}
