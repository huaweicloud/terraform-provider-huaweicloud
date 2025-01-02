package gaussdb

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GaussDB POST /v3/{project_id}/configurations/comparison
func ResourceOpenGaussParameterTemplateCompare() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOpenGaussParameterTemplateCompareCreate,
		ReadContext:   resourceOpenGaussParameterTemplateCompareRead,
		DeleteContext: resourceOpenGaussParameterTemplateCompareDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"target_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"differences": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     gaussDBTemplateCompareDifferencesSchema(),
			},
		},
	}
}

func gaussDBTemplateCompareDifferencesSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_value": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_value": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func resourceOpenGaussParameterTemplateCompareCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/configurations/comparison"
		product = "opengauss"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{KeepResponseBody: true}
	createOpt.JSONBody = utils.RemoveNil(buildCreateOpenGaussParameterTemplateCompareBodyParams(d))

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating GaussDB OpenGauss parameter template compare: %s", err)
	}

	sourceId := d.Get("source_id").(string)
	targetId := d.Get("target_id").(string)
	d.SetId(fmt.Sprintf("%s/%s", sourceId, targetId))

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(
		d.Set("differences", flattenOpenGaussParameterTemplateCompareResponseBody(createRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildCreateOpenGaussParameterTemplateCompareBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"source_id": d.Get("source_id"),
		"target_id": d.Get("target_id"),
	}
	return bodyParams
}

func flattenOpenGaussParameterTemplateCompareResponseBody(resp interface{}) []interface{} {
	differencesJson := utils.PathSearch("differences", resp, make([]interface{}, 0))
	differencesArray := differencesJson.([]interface{})
	if len(differencesArray) < 1 {
		return nil
	}
	rst := make([]interface{}, 0, len(differencesArray))
	for _, v := range differencesArray {
		rst = append(rst, map[string]interface{}{
			"name":         utils.PathSearch("name", v, nil),
			"source_value": utils.PathSearch("source_value", v, nil),
			"target_value": utils.PathSearch("target_value", v, nil),
		})
	}
	return rst
}

func resourceOpenGaussParameterTemplateCompareRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceOpenGaussParameterTemplateCompareDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting parameter template compare resource is not supported. The resource is only removed from the" +
		"state, the GaussDB OpenGauss instance remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
