// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product LTS
// ---------------------------------------------------------------

package lts

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API LTS GET /v3/{project_id}/lts/struct/customtemplate
func DataSourceCustomTemplates() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceCustomTemplatesRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"template_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the ID of the template to be queried.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the custom template name to be queried.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the custom template type to be queried.`,
			},
			"templates": {
				Type:        schema.TypeList,
				Elem:        customTemplateSchema(),
				Computed:    true,
				Description: `The list of structuring custom templates`,
			},
		},
	}
}

func customTemplateSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The structuring custom template ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The structuring custom template name.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The structuring custom template type.`,
			},
			"demo_log": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The sample log event.`,
			},
		},
	}
	return &sc
}

func resourceCustomTemplatesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		mErr    *multierror.Error
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/lts/struct/customtemplate"
		product = "lts"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath += buildListCustomTemplatesQueryParams(d)
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}

	listResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return diag.Errorf("error retrieving LTS structuring custom templates: %s", err)
	}

	listRespBody, err := utils.FlattenResponse(listResp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(generateUUID)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("templates", flattenAndFilterCustomTemplates(listRespBody, d)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAndFilterCustomTemplates(resp interface{}, d *schema.ResourceData) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("results", resp, make([]interface{}, 0))
	curArray, isArray := curJson.([]interface{})
	if !isArray {
		log.Printf("[WARN] error flatten LTS structuring custom templates: " +
			"the results value is not array in API response")
		return nil
	}

	rawName, nameExist := d.GetOk("name")
	rawType, typeExist := d.GetOk("type")
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		templateName := utils.PathSearch("template_name", v, "")
		templateType := utils.PathSearch("template_type", v, "")
		if nameExist && rawName.(string) != templateName.(string) {
			continue
		}
		if typeExist && rawType.(string) != templateType.(string) {
			continue
		}

		rst = append(rst, map[string]interface{}{
			"name":     templateName,
			"type":     templateType,
			"id":       utils.PathSearch("id", v, nil),
			"demo_log": utils.PathSearch("demo_log", v, nil),
		})
	}
	return rst
}

func buildListCustomTemplatesQueryParams(d *schema.ResourceData) string {
	if v, ok := d.GetOk("template_id"); ok {
		return fmt.Sprintf("?id=%s", v.(string))
	}
	return ""
}
