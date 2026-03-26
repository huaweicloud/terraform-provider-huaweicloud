package rds

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var configurationCompareNonUpdatableParams = []string{
	"source_id", "target_id",
}

// @API RDS PUT /v3/{project_id}/configurations/difference
func ResourceRdsConfigurationCompare() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRdsConfigurationCompareCreate,
		ReadContext:   resourceRdsConfigurationCompareRead,
		UpdateContext: resourceRdsConfigurationCompareUpdate,
		DeleteContext: resourceRdsConfigurationCompareDelete,

		CustomizeDiff: config.FlexibleForceNew(configurationCompareNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"source_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"target_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"source_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"parameters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
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
				},
			},
		},
	}
}

func resourceRdsConfigurationCompareCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/configurations/difference"
		product = "rds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	createOpt.JSONBody = buildCreateConfigurationCompareBodyParams(d)

	createResp, err := client.Request("PUT", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating RDS configuration compare: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/%s", d.Get("source_id").(string), d.Get("target_id").(string)))

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("source_id", utils.PathSearch("source_id", createRespBody, nil)),
		d.Set("target_id", utils.PathSearch("target_id", createRespBody, nil)),
		d.Set("source_name", utils.PathSearch("source_name", createRespBody, nil)),
		d.Set("target_name", utils.PathSearch("target_name", createRespBody, nil)),
		d.Set("parameters", flattenCompareResponseBodyDifferences(createRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildCreateConfigurationCompareBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"source_id": d.Get("source_id"),
		"target_id": d.Get("target_id"),
	}
	return bodyParams
}

func flattenCompareResponseBodyDifferences(resp interface{}) []interface{} {
	parameters := utils.PathSearch("parameters", resp, make([]interface{}, 0)).([]interface{})
	if len(parameters) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(parameters))
	for _, parameter := range parameters {
		rst = append(rst, map[string]interface{}{
			"name":         utils.PathSearch("name", parameter, nil),
			"source_value": utils.PathSearch("source_value", parameter, nil),
			"target_value": utils.PathSearch("target_value", parameter, nil),
		})
	}
	return rst
}

func resourceRdsConfigurationCompareRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceRdsConfigurationCompareUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceRdsConfigurationCompareDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting RDS parameter template compare resource is not supported. The resource is only removed from " +
		"the state, the parameter template still remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
