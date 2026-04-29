package dataarts

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

// @API DataArtsStudio GET /v2/{project_id}/design/standards/templates
func DataSourceArchitectureDataStandardTemplateCustoms() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceArchitectureDataStandardTemplateCustomsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the data standard template customs are located.`,
			},

			// Required parameters.
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the workspace to which the data standard template customs belong.`,
			},

			// Attributes.
			"customs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the data standard template custom, in UUID format.`,
						},
						"fd_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The field name of data standard template custom.`,
						},
						"fd_name_en": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The field english name of data standard template custom.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of data standard template custom.`,
						},
						"actived": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the data standard template custom field is visible.`,
						},
						"required": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the data standard template custom field is required.`,
						},
						"searchable": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the data standard template custom field is searchable.`,
						},
						"optional_values": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Valid range for the custom field of the data standard template.`,
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the data standard template custom, in RFC3339 format.`,
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The latest update time of the data standard template custom, in RFC3339 format.`,
						},
						"create_by": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creator of the data standard template custom.`,
						},
						"update_by": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The last editor of the data standard template custom.`,
						},
					},
				},
				Description: `The list of data standard template customs that matched filter parameters.`,
			},
		},
	}
}

func filterArchitectureDataStandardTemplateCustoms(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/design/standards/templates"
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildArchitectureMoreHeaders(d.Get("workspace_id").(string)),
	}

	requestResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}

	customs := utils.PathSearch("data.value.custom", respBody, make([]interface{}, 0)).([]interface{})
	result = append(result, customs...)

	return result, nil
}

func flattenArchitectureDataStandardTemplateCustoms(customs []interface{}) []map[string]interface{} {
	if customs == nil {
		return nil
	}
	result := make([]map[string]interface{}, 0, len(customs))
	for _, custom := range customs {
		result = append(result, map[string]interface{}{
			"id":              utils.PathSearch("id", custom, nil),
			"fd_name":         utils.PathSearch("fd_name", custom, nil),
			"fd_name_en":      utils.PathSearch("fd_name_en", custom, nil),
			"description":     utils.PathSearch("description", custom, nil),
			"actived":         utils.PathSearch("actived", custom, nil),
			"required":        utils.PathSearch("required", custom, nil),
			"searchable":      utils.PathSearch("searchable", custom, nil),
			"optional_values": utils.PathSearch("optional_values", custom, nil),
			"create_time":     utils.PathSearch("create_time", custom, nil),
			"update_time":     utils.PathSearch("update_time", custom, nil),
			"create_by":       utils.PathSearch("create_by", custom, nil),
			"update_by":       utils.PathSearch("update_by", custom, nil),
		})
	}
	return result
}

func resourceArchitectureDataStandardTemplateCustomsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts client: %s", err)
	}

	customs, err := filterArchitectureDataStandardTemplateCustoms(client, d)
	if err != nil {
		return diag.Errorf("error querying DataArts Architecture data standard template customs: %s", err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("customs", flattenArchitectureDataStandardTemplateCustoms(customs)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
