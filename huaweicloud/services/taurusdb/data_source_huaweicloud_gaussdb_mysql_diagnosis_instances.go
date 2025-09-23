package taurusdb

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GaussDBforMySQL GET /v3/{project_id}/instances/diagnosis-instance-infos
func DataSourceGaussDBMysqlDiagnosisInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGaussDBMysqlDiagnosisInstancesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"metric_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the metric name.`,
			},
			"instance_infos": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the information about the abnormal instances.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the instance ID.`,
						},
						"master_node_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the primary node ID.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceGaussDBMysqlDiagnosisInstancesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/diagnosis-instance-infos?offset=0&limit=100"
		product = "gaussdb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath += buildQueryParams(d)

	listResp, err := pagination.ListAllItems(
		client,
		"offset",
		listPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return diag.Errorf("error retrieving GaussDB MySQL diagnosis instances: %s", err)
	}

	listRespJson, err := json.Marshal(listResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var listRespBody interface{}
	err = json.Unmarshal(listRespJson, &listRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("instance_infos", flattenListDiagnosisInstancesBody(listRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildQueryParams(d *schema.ResourceData) string {
	return fmt.Sprintf("&metric_name=%s", d.Get("metric_name").(string))
}

func flattenListDiagnosisInstancesBody(resp interface{}) []interface{} {
	diagnosisInstancesJson := utils.PathSearch("instance_infos", resp, make([]interface{}, 0))
	diagnosisInstancesArray := diagnosisInstancesJson.([]interface{})
	if len(diagnosisInstancesArray) < 1 {
		return nil
	}
	rst := make([]interface{}, 0, len(diagnosisInstancesArray))

	for _, v := range diagnosisInstancesArray {
		rst = append(rst, map[string]interface{}{
			"instance_id":    utils.PathSearch("instance_id", v, nil),
			"master_node_id": utils.PathSearch("master_node_id", v, nil),
		})
	}
	return rst
}
