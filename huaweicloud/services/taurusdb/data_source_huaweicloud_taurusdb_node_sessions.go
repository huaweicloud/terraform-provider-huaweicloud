package taurusdb

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API TaurusDB GET /v3/{project_id}/instances/{instance_id}/nodes/{node_id}/processes
func DataSourceTaurusDBNodeSessions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTaurusDBNodeSessionsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the region in which to query the resource.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the TaurusDB instance.`,
			},
			"node_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the node in the TaurusDB instance.`,
			},
			"processes": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of user session threads in the node in the TaurusDB instance.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The ID of the user session thread.`,
						},
						"user": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The user who starts the session thread.`,
						},
						"host": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the database that is being accessed.`,
						},
						"db": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the database that is being accessed.",
						},
						"command": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The command that is being executed.",
						},
						"time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The time in seconds that the user session thread remains in the current state.",
						},
						"state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of the SQL statement that is being executed.",
						},
						"info": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The additional information, which is usually the statement that is being executed.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTaurusDBNodeSessionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		httpUrl    = "v3/{project_id}/instances/{instance_id}/nodes/{node_id}/processes"
		product    = "gaussdb"
		result     = make([]interface{}, 0)
		totalCount float64
		offset     = 0
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating TaurusDB Client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProviderClient.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{instance_id}", d.Get("instance_id").(string))
	requestPath = strings.ReplaceAll(requestPath, "{node_id}", d.Get("node_id").(string))

	listOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	for {
		currentPath := fmt.Sprintf("%s?limit=100&offset=%v", requestPath, offset)
		resp, err := client.Request("GET", currentPath, &listOpts)
		if err != nil {
			return diag.Errorf("error retrieving TaurusDB node sessions: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		processes := utils.PathSearch("processes", respBody, make([]interface{}, 0)).([]interface{})
		if len(processes) == 0 {
			break
		}
		result = append(result, processes...)

		totalCount = utils.PathSearch("total_count", respBody, float64(0)).(float64)
		if int(totalCount) == len(result) {
			break
		}

		offset += len(processes)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("processes", flattenGetNodeSessionsResponseBody(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetNodeSessionsResponseBody(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	res := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		res = append(res, map[string]interface{}{
			"id":      utils.PathSearch("id", v, nil),
			"user":    utils.PathSearch("user", v, nil),
			"host":    utils.PathSearch("host", v, nil),
			"db":      utils.PathSearch("db", v, nil),
			"command": utils.PathSearch("command", v, nil),
			"time":    utils.PathSearch("time", v, nil),
			"state":   utils.PathSearch("state", v, nil),
			"info":    utils.PathSearch("info", v, nil),
		})
	}
	return res
}
