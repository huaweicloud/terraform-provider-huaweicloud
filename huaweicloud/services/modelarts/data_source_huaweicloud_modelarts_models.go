// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product ModelArts
// ---------------------------------------------------------------

package modelarts

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

// @API ModelArts GET /v1/{project_id}/models
func DataSourceModels() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceModelsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Model name. Fuzzy match is supported. Set **exact_match** to **true** to use exact match.`,
			},
			"exact_match": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Whether to use exact match.`,
			},
			"version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Model version.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Model status.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the model. Fuzzy match is supported.`,
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Workspace ID, which defaults to 0.`,
			},
			"model_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Model type, which is used for obtaining models of this type.`,
			},
			"not_model_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Model type, which is used for obtaining models of types except for this type.`,
			},
			"models": {
				Type:        schema.TypeList,
				Elem:        modelsSchema(),
				Computed:    true,
				Description: `The list of models.`,
			},
		},
	}
}

func modelsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Model ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Model name.`,
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Model version.`,
			},
			"model_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Model type.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Model description that consists of 1 to 100 characters.`,
			},
			"owner": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `User ID of the tenant to which the model belongs.`,
			},
			"source_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Model source type.`,
			},
			"model_source": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Model source.`,
			},
			"install_type": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: `Deployment types supported by the model.`,
			},
			"size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Model size, in bytes.`,
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Workspace ID.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Model status.`,
			},
			"market_flag": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the model is subscribed from AI Gallery.`,
			},
			"tunable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the model can be tuned.`,
			},
			"publishable_flag": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the model can be published to AI Gallery.`,
			},
			"subscription_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Model subscription ID.`,
			},
			"extra": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Extended parameter.`,
			},
			"specification": {
				Type:     schema.TypeList,
				Elem:     modelsSpecificationSchema(),
				Computed: true,
			},
		},
	}
	return &sc
}

func modelsSpecificationSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"min_cpu": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Minimal CPU specifications.`,
			},
			"min_gpu": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Minimal GPU specifications.`,
			},
			"min_memory": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Minimum memory.`,
			},
			"min_ascend": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Minimal Ascend specifications.`,
			},
		},
	}
	return &sc
}

func resourceModelsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// listModels: Query the list of ModelArts models
	var (
		listModelsHttpUrl = "v1/{project_id}/models"
		listModelsProduct = "modelarts"
	)
	listModelsClient, err := cfg.NewServiceClient(listModelsProduct, region)
	if err != nil {
		return diag.Errorf("error creating ModelArts Client: %s", err)
	}

	listModelsPath := listModelsClient.Endpoint + listModelsHttpUrl
	listModelsPath = strings.ReplaceAll(listModelsPath, "{project_id}", listModelsClient.ProjectID)

	listModelsqueryParams := buildListModelsQueryParams(d)
	listModelsPath += listModelsqueryParams

	listModelsResp, err := pagination.ListAllItems(
		listModelsClient,
		"offset",
		listModelsPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving models")
	}

	listModelsRespJson, err := json.Marshal(listModelsResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var listModelsRespBody interface{}
	err = json.Unmarshal(listModelsRespJson, &listModelsRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("models", flattenListModelsModels(listModelsRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListModelsModels(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("models", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":               utils.PathSearch("model_id", v, nil),
			"name":             utils.PathSearch("model_name", v, nil),
			"version":          utils.PathSearch("model_version", v, nil),
			"model_type":       utils.PathSearch("model_type", v, nil),
			"description":      utils.PathSearch("description", v, nil),
			"owner":            utils.PathSearch("owner", v, nil),
			"source_type":      utils.PathSearch("source_type", v, nil),
			"model_source":     utils.PathSearch("model_source", v, nil),
			"install_type":     utils.PathSearch("install_type", v, nil),
			"size":             utils.PathSearch("model_size", v, nil),
			"workspace_id":     utils.PathSearch("workspace_id", v, nil),
			"status":           utils.PathSearch("model_status", v, nil),
			"market_flag":      utils.PathSearch("market_flag", v, nil),
			"tunable":          utils.PathSearch("tunable", v, nil),
			"publishable_flag": utils.PathSearch("publishable_flag", v, nil),
			"subscription_id":  utils.PathSearch("subscription_id", v, nil),
			"extra":            utils.PathSearch("extra", v, nil),
			"specification":    flattenModelsSpecification(v),
		})
	}
	return rst
}

func flattenModelsSpecification(resp interface{}) []interface{} {
	var rst []interface{}
	curJson := utils.PathSearch("specification", resp, make(map[string]interface{})).(map[string]interface{})
	if len(curJson) < 1 {
		return rst
	}

	rst = []interface{}{
		map[string]interface{}{
			"min_cpu":    utils.PathSearch("min_cpu", curJson, nil),
			"min_gpu":    utils.PathSearch("min_gpu", curJson, nil),
			"min_memory": utils.PathSearch("min_memory", curJson, nil),
			"min_ascend": utils.PathSearch("min_ascend", curJson, nil),
		},
	}
	return rst
}

func buildListModelsQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&model_name=%v", res, v)
	}

	if v, ok := d.GetOk("exact_match"); ok {
		res = fmt.Sprintf("%s&exact_match=%v", res, v)
	}

	if v, ok := d.GetOk("version"); ok {
		res = fmt.Sprintf("%s&model_version=%v", res, v)
	}

	if v, ok := d.GetOk("status"); ok {
		res = fmt.Sprintf("%s&model_status=%v", res, v)
	}

	if v, ok := d.GetOk("description"); ok {
		res = fmt.Sprintf("%s&description=%v", res, v)
	}

	if v, ok := d.GetOk("workspace_id"); ok {
		res = fmt.Sprintf("%s&workspace_id=%v", res, v)
	}

	if v, ok := d.GetOk("model_type"); ok {
		res = fmt.Sprintf("%s&model_type=%v", res, v)
	}

	if v, ok := d.GetOk("not_model_type"); ok {
		res = fmt.Sprintf("%s&not_model_type=%v", res, v)
	}

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}
