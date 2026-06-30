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

// @API DataArtsStudio GET /v3/{project_id}/design/subjects
func DataSourceArchitectureSubjects() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceArchitectureSubjectsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The region where the subjects are located.",
			},

			// Required parameters.
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the workspace to which the subjects belong.",
			},

			// Optional parameters.
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name or code of the subject (fuzzy query).",
			},
			"create_by": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The creator of the subject.",
			},
			"owner": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The owner of the subject.",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The business status of the subject.",
			},
			"begin_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The start time for filtering, in RFC3339 format.",
			},
			"end_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The end time for filtering, in RFC3339 format.",
			},
			"parent_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The parent directory ID of the subject.",
			},
			"level": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The level of the subject.",
			},
			"with_relation": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to include relation information.",
			},

			// Attributes.
			"subjects": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataArchitectureSubject(),
				Description: "The list of subjects that matched filter parameters.",
			},
		},
	}
}

func dataArchitectureSubject() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the subject.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Chinese name of the subject.",
			},
			"name_en": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The English name of the subject.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the subject.",
			},
			"qualified_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The qualified name of the subject.",
			},
			"guid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The GUID of the subject, automatically generated.",
			},
			"code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The code of the subject.",
			},
			"alias": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The alias of the subject.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The business status of the subject.",
			},
			"new_biz": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The business version management information, in JSON format.",
			},
			"data_owner": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The data owner of the subject.",
			},
			"data_owner_list": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The data owner list of the subject.",
			},
			"data_department": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The data department of the subject.",
			},
			"path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The path of the subject.",
			},
			"level": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The level of the subject.",
			},
			"ordinal": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The ordinal of the subject.",
			},
			"owner": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The owner of the subject.",
			},
			"parent_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The parent directory ID of the subject.",
			},
			"swap_order_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The swap order ID of the subject.",
			},
			"qualified_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The qualified ID of the subject.",
			},
			"from_public": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the subject is from public layer.",
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creator of the subject.",
			},
			"updated_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The last editor of the subject.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the subject, in RFC3339 format.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The latest update time of the subject, in RFC3339 format.",
			},
			"children_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of child processes.",
			},
			"children": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The child directories, in JSON format.",
			},
			"self_defined_field": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The self-defined field in JSON format.",
			},
			"relations": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataArchitectureSubjectRelation(),
				Description: "The list of relations.",
			},
		},
	}
}

func dataArchitectureSubjectRelation() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the relation.",
			},
			"source_table_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The source table ID of the relation.",
			},
			"target_table_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The target table ID of the relation.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the relation.",
			},
			"source_table_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The source table name of the relation.",
			},
			"target_table_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The target table name of the relation.",
			},
			"role": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The role of the relation.",
			},
			"source_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The source type of the relation.",
			},
			"target_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The target type of the relation.",
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creator of the relation.",
			},
			"updated_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The last editor of the relation.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the relation, in RFC3339 format.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The latest update time of the relation, in RFC3339 format.",
			},
			"mappings": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataArchitectureSubjectMapping(),
				Description: "The list of mappings.",
			},
		},
	}
}

func dataArchitectureSubjectMapping() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the mapping.",
			},
			"relation_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The relation ID of the mapping.",
			},
			"source_field_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The source field ID of the mapping.",
			},
			"target_field_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The target field ID of the mapping.",
			},
			"source_field_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The source field name of the mapping.",
			},
			"target_field_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The target field name of the mapping.",
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creator of the mapping.",
			},
			"updated_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The last editor of the mapping.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the mapping, in RFC3339 format.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The latest update time of the mapping, in RFC3339 format.",
			},
		},
	}
}

func buildArchitectureSubjectsQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	if v, ok := d.GetOk("create_by"); ok {
		res = fmt.Sprintf("%s&create_by=%v", res, v)
	}
	if v, ok := d.GetOk("owner"); ok {
		res = fmt.Sprintf("%s&owner=%v", res, v)
	}
	if v, ok := d.GetOk("status"); ok {
		res = fmt.Sprintf("%s&status=%v", res, v)
	}
	if v, ok := d.GetOk("begin_time"); ok {
		res = fmt.Sprintf("%s&begin_time=%v", res, v)
	}
	if v, ok := d.GetOk("end_time"); ok {
		res = fmt.Sprintf("%s&end_time=%v", res, v)
	}
	if v, ok := d.GetOk("parent_id"); ok {
		res = fmt.Sprintf("%s&parent_id=%v", res, v)
	}
	if v, ok := d.GetOk("level"); ok {
		res = fmt.Sprintf("%s&level=%v", res, v)
	}
	if v, ok := d.GetOk("with_relation"); ok {
		res = fmt.Sprintf("%s&with_relation=%v", res, v)
	}

	return res
}

func listArchitectureSubjects(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/design/subjects?limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))
	listPath += buildArchitectureSubjectsQueryParams(d)

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildArchitectureMoreHeaders(d.Get("workspace_id").(string)),
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &listOpt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		subjects := utils.PathSearch("data.value.records", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, subjects...)

		if len(subjects) < limit {
			break
		}
		offset += len(subjects)
	}

	return result, nil
}

func flattenArchitectureSubjectMappings(mappings []interface{}) []map[string]interface{} {
	if len(mappings) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(mappings))
	for _, mapping := range mappings {
		result = append(result, map[string]interface{}{
			"id":                utils.PathSearch("id", mapping, nil),
			"relation_id":       utils.PathSearch("relation_id", mapping, nil),
			"source_field_id":   utils.PathSearch("source_field_id", mapping, nil),
			"target_field_id":   utils.PathSearch("target_field_id", mapping, nil),
			"source_field_name": utils.PathSearch("source_field_name", mapping, nil),
			"target_field_name": utils.PathSearch("target_field_name", mapping, nil),
			"created_by":        utils.PathSearch("create_by", mapping, nil),
			"updated_by":        utils.PathSearch("update_by", mapping, nil),
			"created_at":        utils.PathSearch("create_time", mapping, nil),
			"updated_at":        utils.PathSearch("update_time", mapping, nil),
		})
	}
	return result
}

func flattenArchitectureSubjectRelations(relations []interface{}) []map[string]interface{} {
	if len(relations) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(relations))
	for _, relation := range relations {
		result = append(result, map[string]interface{}{
			"id":                utils.PathSearch("id", relation, nil),
			"source_table_id":   utils.PathSearch("source_table_id", relation, nil),
			"target_table_id":   utils.PathSearch("target_table_id", relation, nil),
			"name":              utils.PathSearch("name", relation, nil),
			"source_table_name": utils.PathSearch("source_table_name", relation, nil),
			"target_table_name": utils.PathSearch("target_table_name", relation, nil),
			"role":              utils.PathSearch("role", relation, nil),
			"source_type":       utils.PathSearch("source_type", relation, nil),
			"target_type":       utils.PathSearch("target_type", relation, nil),
			"created_by":        utils.PathSearch("create_by", relation, nil),
			"updated_by":        utils.PathSearch("update_by", relation, nil),
			"created_at":        utils.PathSearch("create_time", relation, nil),
			"updated_at":        utils.PathSearch("update_time", relation, nil),
			"mappings": flattenArchitectureSubjectMappings(
				utils.PathSearch("mappings", relation, make([]interface{}, 0)).([]interface{})),
		})
	}
	return result
}

func flattenArchitectureSubjects(subjects []interface{}) []map[string]interface{} {
	if len(subjects) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(subjects))
	for _, subject := range subjects {
		result = append(result, map[string]interface{}{
			"id":                 utils.PathSearch("id", subject, nil),
			"name":               utils.PathSearch("name_ch", subject, nil),
			"name_en":            utils.PathSearch("name_en", subject, nil),
			"description":        utils.PathSearch("description", subject, nil),
			"qualified_name":     utils.PathSearch("qualified_name", subject, nil),
			"guid":               utils.PathSearch("guid", subject, nil),
			"code":               utils.PathSearch("code", subject, nil),
			"alias":              utils.PathSearch("alias", subject, nil),
			"status":             utils.PathSearch("status", subject, nil),
			"new_biz":            utils.JsonToString(utils.PathSearch("new_biz", subject, nil)),
			"data_owner":         utils.PathSearch("data_owner", subject, nil),
			"data_owner_list":    utils.PathSearch("data_owner_list", subject, nil),
			"data_department":    utils.PathSearch("data_department", subject, nil),
			"path":               utils.PathSearch("path", subject, nil),
			"level":              utils.PathSearch("level", subject, nil),
			"ordinal":            utils.PathSearch("ordinal", subject, nil),
			"owner":              utils.PathSearch("owner", subject, nil),
			"parent_id":          utils.PathSearch("parent_id", subject, nil),
			"swap_order_id":      utils.PathSearch("swap_order_id", subject, nil),
			"qualified_id":       utils.PathSearch("qualified_id", subject, nil),
			"from_public":        utils.PathSearch("from_public", subject, nil),
			"created_by":         utils.PathSearch("create_by", subject, nil),
			"updated_by":         utils.PathSearch("update_by", subject, nil),
			"created_at":         utils.PathSearch("create_time", subject, nil),
			"updated_at":         utils.PathSearch("update_time", subject, nil),
			"children_num":       utils.PathSearch("children_num", subject, nil),
			"children":           utils.JsonToString(utils.PathSearch("children", subject, nil)),
			"self_defined_field": utils.JsonToString(utils.PathSearch("self_defined_field", subject, nil)),
			"relations": flattenArchitectureSubjectRelations(
				utils.PathSearch("relations", subject, make([]interface{}, 0)).([]interface{})),
		})
	}
	return result
}

func dataSourceArchitectureSubjectsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts client: %s", err)
	}

	subjects, err := listArchitectureSubjects(client, d)
	if err != nil {
		return diag.Errorf("error querying architecture subjects: %s", err)
	}

	randUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("subjects", flattenArchitectureSubjects(subjects)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
