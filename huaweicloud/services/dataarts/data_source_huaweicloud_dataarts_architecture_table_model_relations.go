package dataarts

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

// @API DataArtsStudio GET /v2/{project_id}/design/{model_id}/table-model/relation
func DataSourceArchitectureTableModelRelations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceArchitectureTableModelRelationsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the table model relations are located.`,
			},

			// Required parameters.
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The workspace ID to which the table model belongs.`,
			},
			"model_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The model ID to which the table model relations belong.`,
			},

			// Optional parameters.
			"biz_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The business type.`,
			},

			// Attributes.
			"relations": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        architectureTableModelRelationSchema(),
				Description: `The list of the table model relations matched filter parameters.`,
			},
		},
	}
}

func architectureTableModelRelationSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the relation.`,
			},
			"relation_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the relation.`,
			},
			"source_table_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the source table.`,
			},
			"source_table_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the source table.`,
			},
			"target_table_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the target table.`,
			},
			"target_table_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the target table.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the relation, in RFC3339 format.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The update time of the relation, in RFC3339 format.`,
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creator of the relation.`,
			},
			"updated_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The updater of the relation.`,
			},
		},
	}
}

func buildArchitectureTableModelRelationsQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("biz_type"); ok {
		res = fmt.Sprintf("%s&biz_type=%v", res, v)
	}

	if res != "" {
		res = "?" + res[1:]
	}

	return res
}

func listArchitectureTableModelRelations(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/design/{model_id}/table-model/relation"
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{model_id}", d.Get("model_id").(string))
	listPath += buildArchitectureTableModelRelationsQueryParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildArchitectureMoreHeaders(d.Get("workspace_id").(string)),
	}

	requestResp, err := client.Request("GET", listPath, &opt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}

	records := utils.PathSearch("data.value.records", respBody, make([]interface{}, 0)).([]interface{})
	for _, record := range records {
		relations := utils.PathSearch("relations", record, make([]interface{}, 0)).([]interface{})
		result = append(result, relations...)
	}

	return result, nil
}

func flattenArchitectureTableModelRelations(relations []interface{}) []interface{} {
	if len(relations) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(relations))
	for _, relation := range relations {
		result = append(result, map[string]interface{}{
			"id":                utils.PathSearch("id", relation, nil),
			"relation_type":     utils.PathSearch("relation_type", relation, nil),
			"source_table_id":   utils.PathSearch("source_table_id", relation, nil),
			"source_table_name": utils.PathSearch("source_table_name", relation, nil),
			"target_table_id":   utils.PathSearch("target_table_id", relation, nil),
			"target_table_name": utils.PathSearch("target_table_name", relation, nil),
			"created_at":        utils.PathSearch("create_time", relation, nil),
			"updated_at":        utils.PathSearch("update_time", relation, nil),
			"created_by":        utils.PathSearch("create_by", relation, nil),
			"updated_by":        utils.PathSearch("update_by", relation, nil),
		})
	}

	return result
}

func dataSourceArchitectureTableModelRelationsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	relations, err := listArchitectureTableModelRelations(client, d)
	if err != nil {
		return diag.FromErr(err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("relations", flattenArchitectureTableModelRelations(relations)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
