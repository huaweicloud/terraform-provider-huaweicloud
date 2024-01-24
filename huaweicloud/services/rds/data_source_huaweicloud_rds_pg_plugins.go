// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product RDS
// ---------------------------------------------------------------

package rds

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RDS GET /v3/{project_id}/instances/{instance_id}/extensions
func DataSourcePgPlugins() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourcePgPluginsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of a PostgreSQL instance.`,
			},
			"database_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the database name of a PostgreSQL instance.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the plugin name.`,
			},
			"version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the plugin version.`,
			},
			"created": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether the plugin has been created.`,
			},
			"plugins": {
				Type:        schema.TypeList,
				Elem:        pgPluginsPluginSchema(),
				Computed:    true,
				Description: `Indicates the plugin list.`,
			},
		},
	}
}

func pgPluginsPluginSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the plugin name.`,
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the plugin version.`,
			},
			"shared_preload_libraries": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the dependent preloaded library.`,
			},
			"created": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether the plugin has been created.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the plugin description.`,
			},
		},
	}
	return &sc
}

func resourcePgPluginsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// listPgPlugins: Query the List of RDS PostgreSQL plugins.
	var (
		listPgPluginsHttpUrl = "v3/{project_id}/instances/{instance_id}/extensions"
		listPgPluginsProduct = "rds"
	)
	listPgPluginsClient, err := cfg.NewServiceClient(listPgPluginsProduct, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	listPgPluginsPath := listPgPluginsClient.Endpoint + listPgPluginsHttpUrl
	listPgPluginsPath = strings.ReplaceAll(listPgPluginsPath, "{project_id}", listPgPluginsClient.ProjectID)
	listPgPluginsPath = strings.ReplaceAll(listPgPluginsPath, "{instance_id}", d.Get("instance_id").(string))

	listPgPluginsQueryParams := buildListPgPluginsQueryParams(d)
	listPgPluginsPath += listPgPluginsQueryParams

	listPgPluginsResp, err := pagination.ListAllItems(
		listPgPluginsClient,
		"offset",
		listPgPluginsPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving RDS PostgreSQL Plugins")
	}

	listPgPluginsRespJson, err := json.Marshal(listPgPluginsResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var listPgPluginsRespBody interface{}
	err = json.Unmarshal(listPgPluginsRespJson, &listPgPluginsRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("plugins", filterListPgPluginsBodyPlugin(flattenListPgPluginsBodyPlugin(listPgPluginsRespBody), d)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListPgPluginsBodyPlugin(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("extensions", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"name":                     utils.PathSearch("name", v, nil),
			"version":                  utils.PathSearch("version", v, nil),
			"shared_preload_libraries": utils.PathSearch("shared_preload_libraries", v, nil),
			"created":                  utils.PathSearch("created", v, nil),
			"description":              utils.PathSearch("description", v, nil),
		})
	}
	return rst
}

func filterListPgPluginsBodyPlugin(all []interface{}, d *schema.ResourceData) []interface{} {
	rst := make([]interface{}, 0, len(all))
	for _, v := range all {
		if param, ok := d.GetOk("name"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("name", v, nil)) {
			continue
		}
		if param, ok := d.GetOk("version"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("version", v, nil)) {
			continue
		}
		createdRaw := d.Get("created").(bool)
		created := utils.PathSearch("created", v, false).(bool)
		if createdRaw != created {
			continue
		}

		rst = append(rst, v)
	}
	return rst
}

func buildListPgPluginsQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("database_name"); ok {
		res = fmt.Sprintf("%s&database_name=%v", res, v)
	}

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}
