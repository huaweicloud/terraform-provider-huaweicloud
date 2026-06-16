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

// @API DataArtsStudio GET /v1/{project_id}/service/servicecatalogs/{catalog_id}/catalogs
func DataSourceDataServiceCatalogs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDataServiceCatalogsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The region where the catalogs are located.",
			},

			// Required parameters.
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the workspace to which the catalogs belong.",
			},

			// Optional parameters.
			"catalog_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "0", // Root path.
				Description: "The ID of the catalog.",
			},

			// Attributes.
			"catalogs": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataDataServiceCatalogsElem(),
				Description: "The list of catalogs that matched filter parameters.",
			},
		},
	}
}

func dataDataServiceCatalogsElem() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the catalog, in UUID format.`,
			},
			"parent_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the parent catalog for the current catalog.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the catalog.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the catalog.",
			},
			"api_catalog_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the API catalog.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the catalog, in RFC3339 format.",
			},
			"create_user": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creator of the catalog.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The last update time of the catalog, in RFC3339 format.",
			},
			"update_user": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The last updater of the catalog.",
			},
		},
	}
	return &sc
}

func buildDataServiceMoreHeaders(workspaceId string) map[string]string {
	result := map[string]string{
		"Content-Type": "application/json",
		"Dlm-Type":     "EXCLUSIVE",
	}

	if workspaceId != "" {
		result["workspace"] = workspaceId
	}

	return result
}

func listDataServiceCatalogs(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/service/servicecatalogs/{catalog_id}/catalogs?limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{catalog_id}", d.Get("catalog_id").(string))
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildDataServiceMoreHeaders(d.Get("workspace_id").(string)),
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

		catalogs := utils.PathSearch("catalogs", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, catalogs...)

		if len(catalogs) < limit {
			break
		}
		offset += len(catalogs)
	}

	return result, nil
}

func flattenDataServiceCatalogs(catalogs []interface{}) []map[string]interface{} {
	if len(catalogs) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(catalogs))
	for _, catalog := range catalogs {
		result = append(result, map[string]interface{}{
			"id":               utils.PathSearch("catalog_id", catalog, nil),
			"parent_id":        utils.PathSearch("pid", catalog, nil),
			"name":             utils.PathSearch("name", catalog, nil),
			"description":      utils.PathSearch("description", catalog, nil),
			"api_catalog_type": utils.PathSearch("api_catalog_type", catalog, nil),
			"created_at":       utils.FormatTimeStampRFC3339(int64(utils.PathSearch("create_time", catalog, float64(0)).(float64))/1000, false),
			"create_user":      utils.PathSearch("create_user", catalog, nil),
			"updated_at":       utils.FormatTimeStampRFC3339(int64(utils.PathSearch("update_time", catalog, float64(0)).(float64))/1000, false),
			"update_user":      utils.PathSearch("update_user", catalog, nil),
		})
	}
	return result
}

func dataSourceDataServiceCatalogsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts client: %s", err)
	}

	catalogs, err := listDataServiceCatalogs(client, d)
	if err != nil {
		return diag.Errorf("error querying DataArts DataService catalogs: %s", err)
	}

	randUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("catalogs", flattenDataServiceCatalogs(catalogs)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
