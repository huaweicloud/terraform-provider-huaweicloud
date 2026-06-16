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

// @API DataArtsStudio GET /v2/{project_id}/factory/scripts
func DataSourceFactoryScripts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFactoryScriptsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the scripts are located.`,
			},

			// Required parameters.
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the workspace to which the scripts belong.`,
			},

			// Optional parameters.
			"script_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of script to be queried.`,
			},

			// Attributes.
			"scripts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The script ID.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The script name.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The script type.`,
						},
						"directory": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The directory path where the script is located.`,
						},
						"create_user": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The user who created the script.`,
						},
						"connection_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The connection name associated with the script.`,
						},
						"database": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The database associated with the script.`,
						},
						"queue_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The DLI queue name associated with the script.`,
						},
						"configuration": {
							Type:        schema.TypeMap,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The user-defined configuration parameters associated with the script.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the script.`,
						},
						"modify_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The last modification time of the script, in RFC3339 format.`,
						},
						"owner": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The owner of the script.`,
						},
						"version": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The version number of the script.`,
						},
					},
				},
				Description: `The list of scripts that matched filter parameters.`,
			},
		},
	}
}

func buildFactoryScriptsQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("script_name"); ok {
		res = fmt.Sprintf("%s&script_name=%v", res, v)
	}

	return res
}

func listFactoryScripts(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/factory/scripts?limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))
	listPath += buildFactoryScriptsQueryParams(d)

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildFactoryMoreHeaders(d.Get("workspace_id").(string), false),
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

		scripts := utils.PathSearch("scripts", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, scripts...)

		if len(scripts) < limit {
			break
		}
		offset += len(scripts)
	}

	return result, nil
}

func flattenFactoryScripts(scripts []interface{}) []map[string]interface{} {
	if len(scripts) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(scripts))
	for _, script := range scripts {
		result = append(result, map[string]interface{}{
			"id":              utils.PathSearch("id", script, nil),
			"name":            utils.PathSearch("name", script, nil),
			"type":            utils.PathSearch("type", script, nil),
			"directory":       utils.PathSearch("directory", script, nil),
			"create_user":     utils.PathSearch("create_user", script, nil),
			"connection_name": utils.PathSearch("connection_name", script, nil),
			"database":        utils.PathSearch("database", script, nil),
			"queue_name":      utils.PathSearch("queue_name", script, nil),
			"configuration":   utils.PathSearch("configuration", script, nil),
			"description":     utils.PathSearch("description", script, nil),
			"modify_time":     utils.FormatTimeStampRFC3339(int64(utils.PathSearch("modify_time", script, float64(0)).(float64))/1000, false),
			"owner":           utils.PathSearch("owner", script, nil),
			"version":         int(utils.PathSearch("version", script, float64(0)).(float64)),
		})
	}
	return result
}

func dataSourceFactoryScriptsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts client: %s", err)
	}

	items, err := listFactoryScripts(client, d)
	if err != nil {
		return diag.Errorf("error querying DataArts factory scripts: %s", err)
	}

	randUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("scripts", flattenFactoryScripts(items)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
