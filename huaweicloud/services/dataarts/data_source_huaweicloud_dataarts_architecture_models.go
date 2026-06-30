package dataarts

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DataArtsStudio GET /v2/{project_id}/design/workspaces
func DataSourceArchitectureModels() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceArchitectureModelsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the architecture models are located.`,
			},

			// Required parameters.
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the workspace.`,
			},

			// Optional parameters.
			"workspace_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The model workspace type.`,
			},
			"dw_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The data connection type.`,
			},

			// Attributes.
			"models": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        architectureModelSchema(),
				Description: `The list of the architecture models that matched filter parameters.`,
			},
		},
	}
}

func architectureModelSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the model.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the workspace.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description of the model.`,
			},
			"is_physical": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether it is a physical table.`,
			},
			"frequent": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether it is frequently used.`,
			},
			"top": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether it is a top-level governance.`,
			},
			"level": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The data governance layer.`,
			},
			"dw_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The data connection type.`,
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the model, in RFC3339 format.`,
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The update time of the model, in RFC3339 format.`,
			},
			"create_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creator of the model.`,
			},
			"update_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The updater of the model.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The workspace type.`,
			},
			"biz_catalog_ids": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The associated business catalog IDs.`,
			},
			"databases": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The database name list.`,
			},
			"table_model_prefix": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The table model validation prefix.`,
			},
		},
	}
}

func buildArchitectureModelsQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("workspace_type"); ok {
		res = fmt.Sprintf("%s&workspace_type=%v", res, v)
	}
	if v, ok := d.GetOk("dw_type"); ok {
		res = fmt.Sprintf("%s&dw_type=%v", res, v)
	}

	return res
}

func listArchitectureModels(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/design/workspaces?limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))
	listPath += buildArchitectureModelsQueryParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildArchitectureMoreHeaders(d.Get("workspace_id").(string)),
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		models := utils.PathSearch("data.value.records", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, models...)
		if len(models) < limit {
			break
		}
		offset += len(models)
	}

	return result, nil
}

func flattenArchitectureModels(models []interface{}) []interface{} {
	if len(models) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(models))
	for _, model := range models {
		result = append(result, map[string]interface{}{
			"id":                 utils.PathSearch("id", model, nil),
			"name":               utils.PathSearch("name", model, nil),
			"description":        utils.PathSearch("description", model, nil),
			"is_physical":        utils.PathSearch("is_physical", model, false),
			"frequent":           utils.PathSearch("frequent", model, false),
			"top":                utils.PathSearch("top", model, false),
			"level":              utils.PathSearch("level", model, nil),
			"dw_type":            utils.PathSearch("dw_type", model, nil),
			"create_time":        utils.PathSearch("create_time", model, nil),
			"update_time":        utils.PathSearch("update_time", model, nil),
			"create_by":          utils.PathSearch("create_by", model, nil),
			"update_by":          utils.PathSearch("update_by", model, nil),
			"type":               utils.PathSearch("type", model, nil),
			"biz_catalog_ids":    utils.PathSearch("biz_catalog_ids", model, nil),
			"databases":          utils.PathSearch("databases", model, make([]interface{}, 0)),
			"table_model_prefix": utils.PathSearch("table_model_prefix", model, nil),
		})
	}

	return result
}

func dataSourceArchitectureModelsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	models, err := listArchitectureModels(client, d)
	if err != nil {
		return diag.Errorf("error querying architecture models: %s", err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("models", flattenArchitectureModels(models)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
