package taurusdb

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

var starrocksParametersCompareNoneUpdatableParams = []string{
	"source_configuration_id",
}

// @API TaurusDB POST /v3/{project_id}/configurations/starrocks/comparison
func ResourceTaurusDBHtapStarrocksParametersCompare() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTaurusDBHtapStarrocksParametersCompareCreate,
		ReadContext:   resourceTaurusDBHtapStarrocksParametersCompareRead,
		UpdateContext: resourceTaurusDBHtapStarrocksParametersCompareUpdate,
		DeleteContext: resourceTaurusDBHtapStarrocksParametersCompareDelete,

		CustomizeDiff: config.FlexibleForceNew(starrocksParametersCompareNoneUpdatableParams),

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
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"differences": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     starrocksParametersCompareDifferencesSchema(),
			},
		},
	}
}

func starrocksParametersCompareDifferencesSchema() *schema.Resource {
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

func resourceTaurusDBHtapStarrocksParametersCompareCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("gaussdb", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	httpUrl := "v3/{project_id}/configurations/starrocks/comparison"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	createOpt.JSONBody = buildCompareStarrocksParametersBodyParams(d)

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error comparing HTAP StarRocks parameters: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id.String())

	differences := utils.PathSearch("differences", respBody, make([]interface{}, 0)).([]interface{})
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("differences", flattenStarrocksParametersCompareDifferences(differences)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildCompareStarrocksParametersBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"source_configuration_id": d.Get("source_configuration_id"),
	}
	return bodyParams
}

func flattenStarrocksParametersCompareDifferences(resp []interface{}) []interface{} {
	res := make([]interface{}, 0)
	for _, v := range resp {
		res = append(res, map[string]interface{}{
			"parameter_name": utils.PathSearch("parameter_name", v, nil),
			"source_value":   utils.PathSearch("source_value", v, nil),
			"target_value":   utils.PathSearch("target_value", v, nil),
		})
	}
	return res
}

func resourceTaurusDBHtapStarrocksParametersCompareRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceTaurusDBHtapStarrocksParametersCompareUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceTaurusDBHtapStarrocksParametersCompareDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting parameters compare resource is not supported. The resource is only removed from the state," +
		" the comparison result remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
