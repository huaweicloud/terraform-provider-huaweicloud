package rds

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RDS GET /v3/{project_id}/instances/diagnosis
func DataSourceRdsDiagnosis() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRdsDiagnosisRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"engine": {
				Type:     schema.TypeString,
				Required: true,
			},
			"diagnosis": {
				Type:     schema.TypeList,
				Elem:     diagnosisItemSchema(),
				Computed: true,
			},
		},
	}
}

func diagnosisItemSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceRdsDiagnosisRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v3/{project_id}/instances/diagnosis"
		product = "rds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += buildDiagnosisQueryParams(d)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving diagnosis counts: %s", err)
	}

	body, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		d.Set("diagnosis", flattenDiagnosis(body)),
		d.Set("region", region),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildDiagnosisQueryParams(d *schema.ResourceData) string {
	return fmt.Sprintf("?engine=%v", d.Get("engine"))
}

func flattenDiagnosis(resp interface{}) []interface{} {
	diagnosis := utils.PathSearch("diagnosis", resp, make([]interface{}, 0))
	diagnosisArr, ok := diagnosis.([]interface{})
	if !ok || len(diagnosisArr) == 0 {
		return nil
	}
	result := make([]interface{}, 0, len(diagnosisArr))
	for _, d := range diagnosisArr {
		item := map[string]interface{}{
			"name":  utils.PathSearch("name", d, nil),
			"count": utils.PathSearch("count", d, nil),
		}
		result = append(result, item)
	}
	return result
}
