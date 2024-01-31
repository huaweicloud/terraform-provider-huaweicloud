// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product DCS
// ---------------------------------------------------------------

package dcs

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

// @API DCS GET /v2/{project_id}/config-templates
func DataSourceTemplates() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceTemplatesRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the type of the template.`,
			},
			"template_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the ID of the template.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the name of the template.`,
			},
			"engine": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the cache engine.`,
			},
			"engine_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the cache engine version.`,
			},
			"cache_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the DCS instance type.`,
			},
			"product_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the product edition.`,
			},
			"storage_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the storage type.`,
			},
			"templates": {
				Type:        schema.TypeList,
				Elem:        templatesTemplateSchema(),
				Computed:    true,
				Description: `Indicates the list of DCS templates.`,
			},
		},
	}
}

func templatesTemplateSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of the template.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the name of the template.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the type of the template.`,
			},
			"engine": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the cache engine.`,
			},
			"engine_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the cache engine version.`,
			},
			"cache_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the DCS instance type.`,
			},
			"product_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the product edition.`,
			},
			"storage_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the storage type.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the description of the template.`,
			},
		},
	}
	return &sc
}

func resourceTemplatesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getDcsTemplates: Query the List of DCS templates.
	var (
		getDcsTemplatesHttpUrl = "v2/{project_id}/config-templates"
		getDcsTemplatesProduct = "dcs"
	)
	getDcsTemplatesClient, err := cfg.NewServiceClient(getDcsTemplatesProduct, region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	getDcsTemplatesPath := getDcsTemplatesClient.Endpoint + getDcsTemplatesHttpUrl
	getDcsTemplatesPath = strings.ReplaceAll(getDcsTemplatesPath, "{project_id}", getDcsTemplatesClient.ProjectID)

	getDcsTemplatesQueryParams := buildGetDcsTemplatesQueryParams(d)
	getDcsTemplatesPath += getDcsTemplatesQueryParams

	getDcsTemplatesResp, err := pagination.ListAllItems(
		getDcsTemplatesClient,
		"offset",
		getDcsTemplatesPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DCS templates")
	}

	getDcsTemplatesRespJson, err := json.Marshal(getDcsTemplatesResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var getDcsTemplatesRespBody interface{}
	err = json.Unmarshal(getDcsTemplatesRespJson, &getDcsTemplatesRespBody)
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
		d.Set("templates", filterGetDcsTemplatesResponseBodyTemplate(
			flattenGetDcsTemplatesResponseBodyTemplate(getDcsTemplatesRespBody), d)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetDcsTemplatesResponseBodyTemplate(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("templates", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"template_id":    utils.PathSearch("template_id", v, nil),
			"name":           utils.PathSearch("name", v, nil),
			"type":           utils.PathSearch("type", v, nil),
			"engine":         utils.PathSearch("engine", v, nil),
			"engine_version": utils.PathSearch("engine_version", v, nil),
			"cache_mode":     utils.PathSearch("cache_mode", v, nil),
			"product_type":   utils.PathSearch("product_type", v, nil),
			"storage_type":   utils.PathSearch("storage_type", v, nil),
			"description":    utils.PathSearch("description", v, nil),
		})
	}
	return rst
}

func filterGetDcsTemplatesResponseBodyTemplate(all []interface{}, d *schema.ResourceData) []interface{} {
	rst := make([]interface{}, 0, len(all))
	for _, v := range all {
		if param, ok := d.GetOk("template_id"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("template_id", v, nil)) {
			continue
		}
		if param, ok := d.GetOk("engine"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("engine", v, nil)) {
			continue
		}
		if param, ok := d.GetOk("engine_version"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("engine_version", v, nil)) {
			continue
		}
		if param, ok := d.GetOk("cache_mode"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("cache_mode", v, nil)) {
			continue
		}
		if param, ok := d.GetOk("product_type"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("product_type", v, nil)) {
			continue
		}
		if param, ok := d.GetOk("storage_type"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("storage_type", v, nil)) {
			continue
		}

		rst = append(rst, v)
	}
	return rst
}

func buildGetDcsTemplatesQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("type"); ok {
		res = fmt.Sprintf("%s&type=%v", res, v)
	}

	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}
