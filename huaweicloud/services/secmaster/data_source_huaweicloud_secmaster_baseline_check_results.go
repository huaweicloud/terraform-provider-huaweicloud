package secmaster

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API SecMaster POST /v2/{project_id}/workspaces/{workspace_id}/sa/baseline/search
func DataSourceSecmasterBaselineCheckResults() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSecmasterBaselineCheckResultsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the workspace ID.`,
			},
			"from_date": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the start time of the baseline check.`,
			},
			"to_date": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the end time of the baseline check.`,
			},
			"condition": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the condition expression.`,
			},
			"baseline_check_results": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of baseline check result.`,
			},
		},
	}
}

func dataSourceSecmasterBaselineCheckResultsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// getBaselineCheckResults: Query the SecMaster baseline check results.
	var (
		listBaselineCheckResultsHttpUrl = "v2/{project_id}/workspaces/{workspace_id}/sa/baseline/search"
		listBaselineCheckResultsProduct = "secmaster"
	)
	client, err := cfg.NewServiceClient(listBaselineCheckResultsProduct, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	listBaselineCheckResultsPath := client.Endpoint + listBaselineCheckResultsHttpUrl
	listBaselineCheckResultsPath = strings.ReplaceAll(listBaselineCheckResultsPath, "{project_id}", client.ProjectID)
	listBaselineCheckResultsPath = strings.ReplaceAll(listBaselineCheckResultsPath, "{workspace_id}", d.Get("workspace_id").(string))

	listBaselineCheckResultsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	bodyParams, err := buildBaselineCheckResultsBodyParams(d)
	if err != nil {
		return diag.FromErr(err)
	}

	baselineCheckResults := make([]interface{}, 0)
	offset := 0
	for {
		bodyParams["offset"] = offset
		listBaselineCheckResultsOpt.JSONBody = bodyParams
		listBaselineCheckResultsResp, err := client.Request("POST", listBaselineCheckResultsPath, &listBaselineCheckResultsOpt)
		if err != nil {
			return diag.FromErr(err)
		}
		listBaselineCheckResultsRespBody, err := utils.FlattenResponse(listBaselineCheckResultsResp)
		if err != nil {
			return diag.FromErr(err)
		}
		data := utils.PathSearch("data", listBaselineCheckResultsRespBody, make([]interface{}, 0)).([]interface{})
		baselineCheckResults = append(baselineCheckResults, data...)

		if len(data) < 1000 {
			break
		}
		offset += 1000
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("baseline_check_results", baselineCheckResults),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildBaselineCheckResultsBodyParams(d *schema.ResourceData) (map[string]interface{}, error) {
	bodyParams := map[string]interface{}{
		"limit": 1000,
	}

	if v, ok := d.GetOk("from_date"); ok {
		fromDateWithZ, err := formatInputTime(v.(string))
		if err != nil {
			return nil, err
		}

		bodyParams["from_date"] = fromDateWithZ
	}
	if v, ok := d.GetOk("to_date"); ok {
		toDateWithZ, err := formatInputTime(v.(string))
		if err != nil {
			return nil, err
		}

		bodyParams["to_date"] = toDateWithZ
	}
	if v, ok := d.GetOk("condition"); ok {
		bodyParams["condition"] = v
	}

	return bodyParams, nil
}
