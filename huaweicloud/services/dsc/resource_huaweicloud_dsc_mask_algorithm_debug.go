package dsc

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var maskAlgorithmDebugNonUpdatableParams = []string{
	"algorithm",
	"algorithm_type",
	"data",
	"parameter",
}

// @API DSC POST /v1/{project_id}/sdg/server/mask/algorithms/test
func ResourceMaskAlgorithmDebug() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMaskAlgorithmDebugCreate,
		ReadContext:   resourceMaskAlgorithmDebugRead,
		UpdateContext: resourceMaskAlgorithmDebugUpdate,
		DeleteContext: resourceMaskAlgorithmDebugDelete,

		CustomizeDiff: config.FlexibleForceNew(maskAlgorithmDebugNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the mask algorithm to be tested is located.`,
			},

			// Required parameters.
			"algorithm": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The mask algorithm.`,
			},
			"algorithm_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The type of the mask algorithm.`,
			},
			"data": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: `The data content to be processed by the mask algorithm.`,
			},

			// Optional parameters.
			"parameter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The configuration parameters of the mask algorithm, in JSON format.`,
			},

			// Attributes.
			"processed_data": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The data content after the mask algorithm processing.`,
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

func buildMaskAlgorithmDebugBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		// Required parameters.
		"algorithm":      d.Get("algorithm"),
		"algorithm_type": d.Get("algorithm_type"),
		"data":           d.Get("data"),
		// Optional parameters.
		"parameter": utils.ValueIgnoreEmpty(d.Get("parameter")),
	}
}

func debugMaskAlgorithm(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	httpUrl := "v1/{project_id}/sdg/server/mask/algorithms/test"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildMaskAlgorithmDebugBodyParams(d)),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func resourceMaskAlgorithmDebugCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("dsc", region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	respBody, err := debugMaskAlgorithm(client, d)
	if err != nil {
		return diag.Errorf("error testing mask algorithm: %s", err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("processed_data", utils.PathSearch("processed_data", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceMaskAlgorithmDebugRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceMaskAlgorithmDebugUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceMaskAlgorithmDebugDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to test a DSC mask algorithm. Deleting this resource
will not clear the corresponding test record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
