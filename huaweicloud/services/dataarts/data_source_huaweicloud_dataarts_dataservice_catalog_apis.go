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

// @API DataArtsStudio GET /v1/{project_id}/service/servicecatalogs/{catalog_id}/apis
func DataSourceDataServiceCatalogApis() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDataServiceCatalogApisRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the catalog and APIs are located.`,
			},

			// Parameters in request header
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the workspace to which the catalog and APIs belong.`,
			},
			"dlm_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The type of DLM engine.`,
			},

			// Query argument
			"catalog_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the catalog to which the APIs belong.`,
			},

			// Attribute
			"apis": {
				Type:        schema.TypeList,
				Elem:        dataserviceCatelogApiSummarySchema(),
				Computed:    true,
				Description: `All API summaries that under the specified catalog.`,
			},
		},
	}
}

func dataserviceCatelogApiSummarySchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The API ID, in UUID format.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The API name.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description of the API.`,
			},
			"group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the group to which the shared API belongs.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the shared API.`,
			},
			"debug_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The debug status of the shared API.`,
			},
			"publish_messages": {
				Type:        schema.TypeList,
				Elem:        dataserviceCatelogApiPublishInfoSchema(),
				Computed:    true,
				Description: `All publish messages of the exclusive API.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The API type.`,
			},
			"manager": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The API reviewer.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the API, in RFC3339 format.`,
			},
			"authorization_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The authorization status of the API.`,
			},
		},
	}
	return &sc
}

func dataserviceCatelogApiPublishInfoSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The publish ID, in UUID format.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the instance used to publish the exclusive API.`,
			},
			"instance_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the instance used to publish the exclusive API.`,
			},
			"api_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The publish status of the exclusive API.`,
			},
			"api_debug": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The debug status of the exclusive API.`,
			},
		},
	}
	return &sc
}

func queryDataServiceApisUnderCatalog(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/service/servicecatalogs/{catalog_id}/apis?limit=100"
		offset  = 0
		result  = make([]interface{}, 0)
	)
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{catalog_id}", d.Get("catalog_id").(string))

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"workspace":    d.Get("workspace_id").(string),
			"dlm-type":     d.Get("dlm_type").(string),
		},
	}

	for {
		listPathWithOffset := listPath + fmt.Sprintf("&offset=%d", offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		apiRecords := utils.PathSearch("apis", respBody, make([]interface{}, 0)).([]interface{})
		if len(apiRecords) < 1 {
			break
		}
		result = append(result, apiRecords...)
		offset += len(apiRecords)
	}

	return result, nil
}

func flattenDataServiceCatalogApiPublishMessages(messages []interface{}) []interface{} {
	result := make([]interface{}, 0, len(messages))
	for _, message := range messages {
		result = append(result, map[string]interface{}{
			"id":            utils.PathSearch("id", message, nil),
			"instance_id":   utils.PathSearch("instance_id", message, nil),
			"instance_name": utils.PathSearch("instance_name", message, nil),
			"api_status":    utils.PathSearch("api_status", message, nil),
			"api_debug":     utils.PathSearch("api_debug", message, nil),
		})
	}
	return result
}

func flattenDataServiceCatalogApis(apiRecords []interface{}) []interface{} {
	result := make([]interface{}, 0, len(apiRecords))

	for _, apiRecord := range apiRecords {
		result = append(result, map[string]interface{}{
			"id":                   utils.PathSearch("id", apiRecord, nil),
			"name":                 utils.PathSearch("name", apiRecord, nil),
			"description":          utils.PathSearch("description", apiRecord, nil),
			"group_id":             utils.PathSearch("group_id", apiRecord, nil),
			"status":               utils.PathSearch("status", apiRecord, nil),
			"debug_status":         utils.PathSearch("debug_status", apiRecord, nil),
			"type":                 utils.PathSearch("type", apiRecord, nil),
			"manager":              utils.PathSearch("manager", apiRecord, nil),
			"authorization_status": utils.PathSearch("authorization_status", apiRecord, nil),
			"publish_messages": flattenDataServiceCatalogApiPublishMessages(utils.PathSearch("publish_messages",
				apiRecord, make([]interface{}, 0)).([]interface{})),
			"created_at": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("create_time",
				apiRecord, float64(0)).(float64))/1000, false),
		})
	}

	return result
}

func dataSourceDataServiceCatalogApisRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	apiRecords, err := queryDataServiceApisUnderCatalog(client, d)
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
		d.Set("apis", flattenDataServiceCatalogApis(apiRecords)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
