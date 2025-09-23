package dataarts

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

// @API DataArtsStudio GET /v2/{project_id}/design/models/statistic
func DataSourceArchitectureModelStatistic() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceArchitectureModelStatisticRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region in which to query model statistic information.`,
			},

			// Parameters in request header
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the workspace in which to query model statistic information.`,
			},

			// Attributes
			"frequent": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataStatisticObjectElem(),
				Description: `The list of the frequently-used objects.`,
			},
			"tops": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataStatisticObjectElem(),
				Description: `The list of the first-layer models.`,
			},
			"logics": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataStatisticObjectElem(),
				Description: `The list of the logical models.`,
			},
			"physicals": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataStatisticObjectElem(),
				Description: `The list of the physical models.`,
			},
			"dwr": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataStatisticObjectElem(),
				Description: `The DWR data reporting layer.`,
			},
			"dm": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataStatisticObjectElem(),
				Description: `The DM data integration layer.`,
			},
		},
	}
}

func dataStatisticObjectElem() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			// Public attributes.
			"biz_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The service type.`,
			},
			"level": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The model level.`,
			},
			"db": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The total database number.`,
			},
			"tb": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The total data table number.`,
			},
			"tb_published": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The published data table number.`,
			},
			"fd": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The total field number.`,
			},
			"fd_published": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The published field number.`,
			},
			"st": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: `The standard coverage.`,
			},
			"st_published": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: `The published standard coverage.`,
			},
			"model": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataStatisticWorkspaceObjectElem(),
				Description: `The model detail.`,
			},
		},
	}
	return &sc
}

func dataStatisticWorkspaceObjectElem() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the model (workspace).`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the model (workspace).`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description of the model (workspace).`,
			},
			"is_physical": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether a table is a physical table.`,
			},
			"frequent": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the model (workspace) is frequently used.`,
			},
			"top": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the model (workspace) is hierarchical governance.`,
			},
			"level": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The level of the data governance layering.`,
			},
			"dw_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the data connection.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the model (workspace), in RFC3339 format.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The latest update time of the model (workspace), in RFC3339 format.`,
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The person who creates the model (workspace).`,
			},
			"updated_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The person who updates the model (workspace).`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the model (workspace).`,
			},
			"biz_catalog_ids": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID list of associated service catalogs.`,
			},
			"databases": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of database names.`,
			},
			"table_model_prefix": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The verification prefix of the model (workspace).`,
			},
		},
	}
	return &sc
}

func getStatisticInfo(client *golangsdk.ServiceClient, workspaceId string) (interface{}, error) {
	httpUrl := "v2/{project_id}/design/models/statistic"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"workspace":    workspaceId,
		},
	}

	requestResp, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return nil, fmt.Errorf("error querying architecture statistic information: %s", err)
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}
	return utils.PathSearch("data.value", respBody, nil), nil
}

func ignoreStrResultNaN(num interface{}) interface{} {
	if strNum, ok := num.(string); ok && strNum == "NaN" {
		return nil
	}
	return num
}

func flattenStatisticObject(frequent []interface{}) []interface{} {
	result := make([]interface{}, 0, len(frequent))
	for _, val := range frequent {
		result = append(result, map[string]interface{}{
			"biz_type":     utils.PathSearch("biz_type", val, nil),
			"level":        utils.PathSearch("level", val, nil),
			"db":           utils.PathSearch("db", val, nil),
			"tb":           utils.PathSearch("tb", val, nil),
			"tb_published": utils.PathSearch("tb_published", val, nil),
			"fd":           utils.PathSearch("fd", val, nil),
			"fd_published": utils.PathSearch("fd_published", val, nil),
			"st":           ignoreStrResultNaN(utils.PathSearch("st", val, nil)),
			"st_published": ignoreStrResultNaN(utils.PathSearch("st_published", val, nil)),
			"model": flattenStatisticWorkspaceObjectResult(utils.PathSearch("model", val,
				make(map[string]interface{})).(map[string]interface{})),
		})
	}
	return result
}

func flattenStatisticWorkspaceObjectResult(workspaceObj map[string]interface{}) []map[string]interface{} {
	if len(workspaceObj) < 1 {
		return nil
	}
	return []map[string]interface{}{
		{
			"id":                 utils.PathSearch("id", workspaceObj, nil),
			"name":               utils.PathSearch("name", workspaceObj, nil),
			"description":        utils.PathSearch("description", workspaceObj, nil),
			"is_physical":        utils.PathSearch("is_physical", workspaceObj, nil),
			"frequent":           utils.PathSearch("frequent", workspaceObj, nil),
			"top":                utils.PathSearch("top", workspaceObj, nil),
			"level":              utils.PathSearch("level", workspaceObj, nil),
			"dw_type":            utils.PathSearch("dw_type", workspaceObj, nil),
			"created_at":         utils.PathSearch("create_time", workspaceObj, nil),
			"updated_at":         utils.PathSearch("update_time", workspaceObj, nil),
			"created_by":         utils.PathSearch("create_by", workspaceObj, nil),
			"updated_by":         utils.PathSearch("update_by", workspaceObj, nil),
			"type":               utils.PathSearch("type", workspaceObj, nil),
			"biz_catalog_ids":    utils.PathSearch("biz_catalog_ids", workspaceObj, nil),
			"databases":          utils.PathSearch("databases", workspaceObj, nil),
			"table_model_prefix": utils.PathSearch("table_model_prefix", workspaceObj, nil),
		},
	}
}

func flattenStatisticFrequent(frequent []interface{}) []interface{} {
	return flattenStatisticObject(frequent)
}

func flattenStatisticTops(tops []interface{}) []interface{} {
	return flattenStatisticObject(tops)
}

func flattenStatisticLogics(logics []interface{}) []interface{} {
	return flattenStatisticObject(logics)
}

func flattenStatisticPhysicals(physicals []interface{}) []interface{} {
	return flattenStatisticObject(physicals)
}

func flattenStatisticDwr(dwrDetail interface{}) []interface{} {
	if dwrDetail == nil {
		return nil
	}
	return flattenStatisticObject([]interface{}{dwrDetail})
}

func flattenStatisticDm(dmDetail interface{}) []interface{} {
	if dmDetail == nil {
		return nil
	}
	return flattenStatisticObject([]interface{}{dmDetail})
}

func dataSourceArchitectureModelStatisticRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	queryResp, err := getStatisticInfo(client, d.Get("workspace_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("frequent", flattenStatisticFrequent(utils.PathSearch("frequent", queryResp, make([]interface{}, 0)).([]interface{}))),
		d.Set("tops", flattenStatisticTops(utils.PathSearch("top", queryResp, make([]interface{}, 0)).([]interface{}))),
		d.Set("logics", flattenStatisticLogics(utils.PathSearch("logic", queryResp, make([]interface{}, 0)).([]interface{}))),
		d.Set("physicals", flattenStatisticPhysicals(utils.PathSearch("physical", queryResp, make([]interface{}, 0)).([]interface{}))),
		d.Set("dwr", flattenStatisticDwr(utils.PathSearch("dwr", queryResp, nil))),
		d.Set("dm", flattenStatisticDm(utils.PathSearch("dm", queryResp, nil))),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
