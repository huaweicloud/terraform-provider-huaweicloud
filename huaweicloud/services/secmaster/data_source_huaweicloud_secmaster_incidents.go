package secmaster

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

// @API SecMaster POST /v1/{project_id}/workspaces/{workspace_id}/soc/incidents/search
func DataSourceIncidents() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIncidentsRead,

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
				Description: `Specifies the ID of the workspace to which the incident belongs.`,
			},
			"from_date": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the search start time.`,
			},
			"to_date": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the search end time.`,
			},
			"condition": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Specifies the search condition expression.`,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"conditions": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: `Specifies the condition expression list.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `Specifies the expression name.`,
									},
									"data": {
										Type:        schema.TypeList,
										Optional:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: `Specifies the expression content.`,
									},
								},
							},
						},
						"logics": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `Specifies the expression logic.`,
						},
					},
				},
			},
			"incidents": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The incident list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The incident ID.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The incident name.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The incident description.`,
						},
						"type": {
							Type:        schema.TypeList,
							Elem:        IncidentsTypeSchema(),
							Computed:    true,
							Description: `The incident type configuration.`,
						},
						"level": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The incident level.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The incident status.`,
						},
						"data_source": {
							Type:        schema.TypeList,
							Elem:        IncidentsDataSourceSchema(),
							Computed:    true,
							Description: `The data source configuration.`,
						},
						"first_occurrence_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The first occurrence time of the incident.`,
						},
						"owner": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The user name of the owner.`,
						},
						"last_occurrence_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The last occurrence time of the incident.`,
						},
						"planned_closure_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The planned closure time of the incident.`,
						},
						"verification_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The verification status.`,
						},
						"stage": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The stage of the incident.`,
						},
						"debugging_data": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Whether it's a debugging data.`,
						},
						"labels": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The labels.`,
						},
						"close_reason": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The close reason.`,
						},
						"close_comment": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The close comment.`,
						},
						"creator": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name creator name.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The update time.`,
						},
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The version of the data source of an incident.`,
						},
						"domain_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the account (domain_id) to whom the data is delivered and hosted.`,
						},
						"project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of project where the account to whom the data is delivered and hosted belongs to.`,
						},
						"region_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the region where the account to whom the data is delivered and hosted belongs to.`,
						},
						"workspace_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the current workspace.`,
						},
						"arrive_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The data receiving time.`,
						},
						"count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The times of the incident occurrences.`,
						},
						"ipdrr_phase": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The handling phase No.`,
						},
					},
				},
			},
		},
	}
}

func IncidentsTypeSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"category": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The category.`,
			},
			"incident_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The incident type.`,
			},
		},
	}
	return &sc
}

func IncidentsDataSourceSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"product_feature": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The product feature.`,
			},
			"product_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The product name.`,
			},
			"source_type": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The source type.`,
			},
		},
	}
	return &sc
}

func dataSourceIncidentsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getIncident: Query the SecMaster incidents
	var (
		listIncidentHttpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/incidents/search"
		listIncidentProduct = "secmaster"
	)
	client, err := cfg.NewServiceClient(listIncidentProduct, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	listIncidentPath := client.Endpoint + listIncidentHttpUrl
	listIncidentPath = strings.ReplaceAll(listIncidentPath, "{project_id}", client.ProjectID)
	listIncidentPath = strings.ReplaceAll(listIncidentPath, "{workspace_id}", fmt.Sprintf("%v", d.Get("workspace_id")))

	listIncidentOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	bodyParams, err := buildIncidentsBodyParams(d)
	if err != nil {
		return diag.FromErr(err)
	}

	incidents := make([]interface{}, 0)
	offset := 0
	for {
		bodyParams["offset"] = offset
		listIncidentOpt.JSONBody = bodyParams
		listIncidentResp, err := client.Request("POST", listIncidentPath, &listIncidentOpt)
		if err != nil {
			return diag.FromErr(err)
		}
		listIncidentRespBody, err := utils.FlattenResponse(listIncidentResp)
		if err != nil {
			return diag.FromErr(err)
		}
		data := utils.PathSearch("data", listIncidentRespBody, make([]interface{}, 0)).([]interface{})
		incidents = append(incidents, data...)

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

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("incidents", flattenIncidentsResponseBody(incidents)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildIncidentsBodyParams(d *schema.ResourceData) (map[string]interface{}, error) {
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
	if v, ok := d.GetOk("condition.0"); ok {
		bodyParams["condition"] = v
	}

	return bodyParams, nil
}

func flattenIncidentsResponseBody(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	incidents := make([]interface{}, len(resp))
	for i, v := range resp {
		dataObject := utils.PathSearch("data_object", v, nil)
		incidents[i] = map[string]interface{}{
			"id":                    utils.PathSearch("id", dataObject, nil),
			"name":                  utils.PathSearch("title", dataObject, nil),
			"description":           utils.PathSearch("description", dataObject, nil),
			"type":                  flattenGetIncidentResponseBodyType(dataObject),
			"level":                 utils.PathSearch("severity", dataObject, nil),
			"status":                utils.PathSearch("handle_status", dataObject, nil),
			"owner":                 utils.PathSearch("owner", dataObject, nil),
			"data_source":           flattenGetIncidentResponseBodyDataSource(dataObject),
			"first_occurrence_time": utils.PathSearch("first_observed_time", dataObject, nil),
			"last_occurrence_time":  utils.PathSearch("last_observed_time", dataObject, nil),
			"verification_status":   utils.PathSearch("verification_state", dataObject, nil),
			"stage":                 utils.PathSearch("ipdrr_phase", dataObject, nil),
			"debugging_data":        fmt.Sprintf("%v", utils.PathSearch("simulation", dataObject, nil)),
			"labels":                utils.PathSearch("labels", dataObject, nil),
			"close_reason":          utils.PathSearch("close_reason", dataObject, nil),
			"close_comment":         utils.PathSearch("close_comment", dataObject, nil),
			"creator":               utils.PathSearch("creator", dataObject, nil),
			"created_at":            utils.PathSearch("create_time", dataObject, nil),
			"updated_at":            utils.PathSearch("update_time", dataObject, nil),
			"planned_closure_time":  utils.PathSearch("sla", dataObject, nil),
			"version":               utils.PathSearch("version", dataObject, nil),
			"domain_id":             utils.PathSearch("domain_id", dataObject, nil),
			"project_id":            utils.PathSearch("project_id", dataObject, nil),
			"region_id":             utils.PathSearch("region_id", dataObject, nil),
			"workspace_id":          utils.PathSearch("workspace_id", dataObject, nil),
			"arrive_time":           utils.PathSearch("arrive_time", dataObject, nil),
			"count":                 utils.PathSearch("count", dataObject, 0),
			"ipdrr_phase":           utils.PathSearch("ipdrr_phase", dataObject, nil),
		}
	}
	return incidents
}
