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

// @API RDS GET /v3/{project_id}/instances/diagnosis/info
func DataSourceRdsDiagnosisInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRdsDiagnosisInstancesRead,
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
				Type:     schema.TypeString,
				Required: true,
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceRdsDiagnosisInstancesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v3/{project_id}/instances/diagnosis/info"
		product = "rds"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	var allInstances []interface{}
	offset := 0
	for {
		path := getPath + buildDiagnosisInstancesQueryParams(d, offset)
		getResp, err := client.Request("GET", path, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving RDS diagnosis result: %s", err)
		}
		body, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}

		instances := flattenDiagnosisInfoBody(body)
		allInstances = append(allInstances, instances...)

		if len(instances) < 10 {
			break
		}
		offset += 10
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("instances", allInstances),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

// ?engine=sqlserver&diagnosis=high_pressure&offset=1&limit=10
func buildDiagnosisInstancesQueryParams(d *schema.ResourceData, offset int) string {
	return fmt.Sprintf("?engine=%s&diagnosis=%s&offset=%d&limit=10",
		d.Get("engine").(string),
		d.Get("diagnosis").(string),
		offset,
	)
}

func flattenDiagnosisInfoBody(resp interface{}) (instances []interface{}) {
	instancesJson := utils.PathSearch("instances", resp, make([]interface{}, 0))
	instancessArr := instancesJson.([]interface{})
	instances = make([]interface{}, 0, len(instancessArr))
	for _, in := range instancessArr {
		instances = append(instances, map[string]interface{}{
			"id": utils.PathSearch("id", in, nil),
		})
	}
	return
}
