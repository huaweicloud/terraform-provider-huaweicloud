package secmaster

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API SecMaster POST /v2/{project_id}/workspaces/{workspace_id}/sa/baseline/compliance-packages/search
func DataSourceCompliancePackages() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCompliancePackagesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"state": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"builtin_compliance_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"customized_compliance_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"disabled_compliance_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"enabled_compliance_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"compliance_packages": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     compliancePackageSchema(),
			},
		},
	}
}

func compliancePackageSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"spec_catalog_vo_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     baselineCatalogSchema(),
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"classify": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"areas": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"check_items_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"has_auto_check_items": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func baselineCatalogSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"serial_number": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"level_number": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"root": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"parent": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_leaf": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"check_items": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     checkitemCatalogSchema(),
			},
		},
	}
}

func checkitemCatalogSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildCompliancePackagesBodyParams(d *schema.ResourceData, limit, offset int) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"limit":  limit,
		"offset": offset,
	}
	if v, ok := d.GetOk("name"); ok {
		bodyParams["name"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		bodyParams["description"] = v
	}
	if v, ok := d.GetOk("type"); ok {
		bodyParams["type"] = v
	}
	if v, ok := d.GetOk("state"); ok {
		bodyParams["state"] = v
	}

	return bodyParams
}

func dataSourceCompliancePackagesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "secmaster"
		httpUrl = "v2/{project_id}/workspaces/{workspace_id}/sa/baseline/compliance-packages/search"
		result  = make([]interface{}, 0)
		limit   = 1000
		offset  = 0
		mErr    *multierror.Error
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", d.Get("workspace_id").(string))

	reqOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"content-type": "application/json;charset=UTF-8",
			"x-language":   "en-us",
		},
	}

	for {
		currentBodyParams := buildCompliancePackagesBodyParams(d, limit, offset)
		reqOpt.JSONBody = utils.RemoveNil(currentBodyParams)

		resp, err := client.Request("POST", requestPath, &reqOpt)
		if err != nil {
			return diag.Errorf("error retrieving SecMaster compliance packages: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		packagesRaw := utils.PathSearch("compliance_packages", respBody, make([]interface{}, 0)).([]interface{})
		if len(packagesRaw) == 0 {
			break
		}

		result = append(result, packagesRaw...)

		if len(packagesRaw) < limit {
			break
		}

		offset += len(packagesRaw)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("builtin_compliance_num", utils.PathSearch("builtin_compliance_num", result, nil)),
		d.Set("customized_compliance_num", utils.PathSearch("customized_compliance_num", result, nil)),
		d.Set("disabled_compliance_num", utils.PathSearch("disabled_compliance_num", result, nil)),
		d.Set("enabled_compliance_num", utils.PathSearch("enabled_compliance_num", result, nil)),
		d.Set("compliance_packages", flattenCompliancePackages(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCompliancePackages(result []interface{}) []interface{} {
	if len(result) == 0 {
		return nil
	}

	packages := make([]interface{}, 0, len(result))
	for _, pkg := range result {
		packages = append(packages, map[string]interface{}{
			"uuid":                 utils.PathSearch("uuid", pkg, nil),
			"name":                 utils.PathSearch("name", pkg, nil),
			"version":              utils.PathSearch("version", pkg, nil),
			"owner":                utils.PathSearch("owner", pkg, nil),
			"spec_catalog_vo_list": flattenBaselineCatalogs(utils.PathSearch("spec_catalog_vo_list", pkg, make([]interface{}, 0)).([]interface{})),
			"description":          utils.PathSearch("description", pkg, nil),
			"classify":             utils.PathSearch("classify", pkg, nil),
			"areas":                utils.PathSearch("areas", pkg, nil),
			"region":               utils.PathSearch("region", pkg, nil),
			"state":                utils.PathSearch("state", pkg, nil),
			"type":                 utils.PathSearch("type", pkg, nil),
			"check_items_num":      utils.PathSearch("check_items_num", pkg, nil),
			"has_auto_check_items": utils.PathSearch("has_auto_check_items", pkg, nil),
		})
	}

	return packages
}

func flattenBaselineCatalogs(catalogs []interface{}) []interface{} {
	if len(catalogs) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(catalogs))
	for _, catalog := range catalogs {
		result = append(result, map[string]interface{}{
			"uuid":          utils.PathSearch("uuid", catalog, nil),
			"serial_number": utils.PathSearch("serial_number", catalog, nil),
			"level_number":  utils.PathSearch("level_number", catalog, nil),
			"root":          utils.PathSearch("root", catalog, nil),
			"parent":        utils.PathSearch("parent", catalog, nil),
			"is_leaf":       utils.PathSearch("is_leaf", catalog, nil),
			"check_items":   flattenCheckitemCatalogs(utils.PathSearch("check_items", catalog, make([]interface{}, 0)).([]interface{})),
		})
	}

	return result
}

func flattenCheckitemCatalogs(checkItems []interface{}) []interface{} {
	if len(checkItems) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(checkItems))
	for _, item := range checkItems {
		result = append(result, map[string]interface{}{
			"uuid": utils.PathSearch("uuid", item, nil),
			"name": utils.PathSearch("name", item, nil),
		})
	}

	return result
}
