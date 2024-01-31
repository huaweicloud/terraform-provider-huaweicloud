// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product DataArts
// ---------------------------------------------------------------

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

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DataArtsStudio GET /v2/{project_id}/design/standards/templates
func DataSourceTemplateOptionalFields() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceTemplateOptionalFieldsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the workspace ID of DataArts Architecture.`,
			},
			"fd_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the name of the optional field.`,
			},
			"required": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether the field is required.`,
			},
			"searchable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether the field is search supported.`,
			},
			"optional_fields": {
				Type:        schema.TypeList,
				Elem:        templateOptionalFieldsOptionalFieldSchema(),
				Computed:    true,
				Description: `Indicates the list of DataArts Architecture data standard template optional fields.`,
			},
		},
	}
}

func templateOptionalFieldsOptionalFieldSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"fd_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the name of the field.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the description of the field.`,
			},
			"description_en": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the English description of the field.`,
			},
			"required": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether the field is required.`,
			},
			"searchable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether the field is search supported.`,
			},
		},
	}
	return &sc
}

func resourceTemplateOptionalFieldsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getTemplateOptionalFields: Query the List of DataArts Architecture data standard template optional fields
	var (
		getTemplateOptionalFieldsHttpUrl = "v2/{project_id}/design/standards/templates"
		getTemplateOptionalFieldsProduct = "dataarts"
	)
	getTemplateOptionalFieldsClient, err := cfg.NewServiceClient(getTemplateOptionalFieldsProduct, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	getTemplateOptionalFieldsPath := getTemplateOptionalFieldsClient.Endpoint + getTemplateOptionalFieldsHttpUrl
	getTemplateOptionalFieldsPath = strings.ReplaceAll(getTemplateOptionalFieldsPath, "{project_id}",
		getTemplateOptionalFieldsClient.ProjectID)

	getTemplateFieldsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": d.Get("workspace_id").(string)},
	}

	getTemplateOptionalFieldsResp, err := getTemplateOptionalFieldsClient.Request("GET",
		getTemplateOptionalFieldsPath, &getTemplateFieldsOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DataArts Architecture data standard template optional fields")
	}

	getTemplateOptionalFieldsRespBody, err := utils.FlattenResponse(getTemplateOptionalFieldsResp)
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
		d.Set("optional_fields", filterGetTemplateOptionalFieldsResponseBodyOptional(
			flattenGetTemplateOptionalFieldsResponseBodyOptional(getTemplateOptionalFieldsRespBody), d)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetTemplateOptionalFieldsResponseBodyOptional(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("data.value.preFields_optional", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"fd_name":        utils.PathSearch("fd_name", v, nil),
			"description":    utils.PathSearch("description", v, nil),
			"description_en": utils.PathSearch("descriptionEn", v, nil),
			"required":       utils.PathSearch("required", v, nil),
			"searchable":     utils.PathSearch("searchable", v, nil),
		})
	}
	return rst
}

func filterGetTemplateOptionalFieldsResponseBodyOptional(all []interface{}, d *schema.ResourceData) []interface{} {
	rst := make([]interface{}, 0, len(all))
	rawRequired, _ := d.GetOk("required")
	rawSearchable, _ := d.GetOk("searchable")
	for _, v := range all {
		if param, ok := d.GetOk("fd_name"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("fd_name", v, nil)) {
			continue
		}
		if fmt.Sprint(rawRequired) != fmt.Sprint(utils.PathSearch("required", v, nil)) {
			continue
		}
		if fmt.Sprint(rawSearchable) != fmt.Sprint(utils.PathSearch("searchable", v, nil)) {
			continue
		}

		rst = append(rst, v)
	}
	return rst
}
