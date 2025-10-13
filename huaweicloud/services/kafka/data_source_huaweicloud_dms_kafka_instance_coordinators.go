package kafka

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

// @API Kafka GET /v2/{project_id}/instances/{instance_id}/management/coordinators
func DataSourceInstanceCoordinators() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceInstanceCoordinatorsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the coordinators are located.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the Kafka instance to which the coordinators belong.`,
			},
			"coordinators": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The ID of the broker of the corresponding to the coordinator.`,
						},
						"group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the consumer group.`,
						},
						"host": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The address of the broker of the corresponding to the coordinator.`,
						},
						"port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The port number of the corresponding to the coordinator.`,
						},
					},
				},
				Description: `The list of coordinators corresponding to all consumer groups.`,
			},
		},
	}
}

func dataSourceInstanceCoordinatorsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
	)
	client, err := cfg.NewServiceClient("dmsv2", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	httpUrl := "v2/{project_id}/instances/{instance_id}/management/coordinators"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", instanceId)

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf-8",
		},
	}

	listResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return diag.Errorf("error retrieving coordinator list of Kafka instance (%s): %s", instanceId, err)
	}

	listRespBody, err := utils.FlattenResponse(listResp)
	if err != nil {
		return diag.FromErr(err)
	}

	randomId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(randomId)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("coordinators", flattenCoordinators(utils.PathSearch("coordinators",
			listRespBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCoordinators(coordinators []interface{}) []interface{} {
	if len(coordinators) == 0 {
		return nil
	}

	results := make([]interface{}, 0, len(coordinators))
	for _, coordinator := range coordinators {
		results = append(results, map[string]interface{}{
			"group_id": utils.PathSearch("group_id", coordinator, nil),
			"id":       utils.PathSearch("id", coordinator, nil),
			"host":     utils.PathSearch("host", coordinator, nil),
			"port":     utils.PathSearch("port", coordinator, nil),
		})
	}
	return results
}
