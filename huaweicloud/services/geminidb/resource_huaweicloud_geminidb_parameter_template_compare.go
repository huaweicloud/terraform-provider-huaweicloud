package geminidb

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

var geminiDBParameterTemplateCompareNonUpdatableParams = []string{
	"target_configuration_id", "target_configuration_id",
}

// @API GeminiDB POST /v3/{project_id}/configurations/comparison
func ResourceGeminiDBParameterTemplateCompare() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGeminiDBParameterTemplateCompareCreate,
		UpdateContext: resourceGeminiDBParameterTemplateCompareUpdate,
		ReadContext:   resourceGeminiDBParameterTemplateCompareRead,
		DeleteContext: resourceGeminiDBParameterTemplateCompareDelete,

		CustomizeDiff: config.FlexibleForceNew(geminiDBParameterTemplateCompareNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"source_configuration_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"target_configuration_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"differences": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     geminiDBParameterTemplateCompareDifferenceSchema(),
			},
		},
	}
}

func geminiDBParameterTemplateCompareDifferenceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"parameter_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_value": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_value": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceGeminiDBParameterTemplateCompareCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("geminidb", region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	httpUrl := "v3/{project_id}/configurations/comparison"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createReq := map[string]interface{}{
		"source_configuration_id": d.Get("source_configuration_id"),
		"target_configuration_id": d.Get("target_configuration_id"),
	}

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         createReq,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error comparing GeminiDB parameter template compare: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.Errorf("error flattening GeminiDB parameter template compare response: %s", err)
	}

	differences := flattenGeminiDBTemplateCompareDifferences(createRespBody)

	resourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(resourceId.String())

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("differences", differences),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGeminiDBTemplateCompareDifferences(resp interface{}) []map[string]interface{} {
	differencesRaw := utils.PathSearch("differences", resp, nil)
	if differencesRaw == nil {
		return nil
	}

	differencesSlice, ok := differencesRaw.([]interface{})
	if !ok {
		return nil
	}

	differences := make([]map[string]interface{}, 0, len(differencesSlice))
	for _, diffRaw := range differencesSlice {
		diffMap := map[string]interface{}{
			"parameter_name": utils.PathSearch("parameter_name", diffRaw, nil),
			"source_value":   utils.PathSearch("source_value", diffRaw, nil),
			"target_value":   utils.PathSearch("target_value", diffRaw, nil),
		}
		differences = append(differences, diffMap)
	}

	return differences
}

func resourceGeminiDBParameterTemplateCompareUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceGeminiDBParameterTemplateCompareRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceGeminiDBParameterTemplateCompareDelete(_ context.Context, _ *schema.ResourceData,
	_ interface{}) diag.Diagnostics {
	errorMsg := "Deleting a parameter template compare resource is not supported. The comparison operation is a one-time action " +
		"and does not create any persistent resources."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Resource deletion is not supported",
			Detail:   errorMsg,
		},
	}
}
