package cfw

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var batchDeleteDomainSetsNonUpdatableParams = []string{"object_id", "set_ids", "fw_instance_id", "enterprise_project_id"}

// @API CFW POST /v1/{project_id}/domain-sets/batch-delete
func ResourceBatchDeleteDomainSets() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBatchDeleteDomainSetsCreate,
		ReadContext:   resourceBatchDeleteDomainSetsRead,
		UpdateContext: resourceBatchDeleteDomainSetsUpdate,
		DeleteContext: resourceBatchDeleteDomainSetsDelete,

		CustomizeDiff: config.FlexibleForceNew(batchDeleteDomainSetsNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"object_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"set_ids": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
			"fw_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"data": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"response_data": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func buildBatchDeleteDomainSetsQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := ""

	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}
	if v, ok := d.GetOk("fw_instance_id"); ok {
		queryParams = fmt.Sprintf("%s&fw_instance_id=%v", queryParams, v)
	}
	if queryParams != "" {
		queryParams = "?" + queryParams[1:]
	}

	return queryParams
}

func buildBatchDeleteDomainSetsBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"object_id": d.Get("object_id"),
		"set_ids":   utils.ExpandToStringList(d.Get("set_ids").([]interface{})),
	}

	return bodyParams
}

func resourceBatchDeleteDomainSetsCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		objectId = d.Get("object_id").(string)
		epsId    = cfg.GetEnterpriseProjectID(d)
		httpUrl  = "v1/{project_id}/domain-sets/batch-delete"
	)

	client, err := cfg.NewServiceClient("cfw", region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildBatchDeleteDomainSetsQueryParams(d, epsId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildBatchDeleteDomainSetsBodyParams(d),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error batch deleting CFW domain sets: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(objectId)

	return diag.FromErr(d.Set("data", flattenBatchDeleteDomainSetsDataResp(respBody)))
}

func flattenBatchDeleteDomainSetsDataResp(respBody interface{}) []map[string]interface{} {
	responseDataResp := utils.PathSearch("data.responseDatas", respBody, make([]interface{}, 0)).([]interface{})
	if len(responseDataResp) == 0 {
		return nil
	}

	responseData := make([]map[string]interface{}, 0)
	for _, v := range responseDataResp {
		responseData = append(responseData, map[string]interface{}{
			"name": utils.PathSearch("name", v, nil),
			"id":   utils.PathSearch("id", v, nil),
		})
	}

	return []map[string]interface{}{
		{
			"response_data": responseData,
		},
	}
}

func resourceBatchDeleteDomainSetsRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in 'Read()' method because resource is a one-time action resource.
	return nil
}

func resourceBatchDeleteDomainSetsUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in 'Update()' method because resource is a one-time action resource.
	return nil
}

func resourceBatchDeleteDomainSetsDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to batch delete domain sets. Deleting this resource
    will not clear the corresponding request record, but will only remove resource information from
    the tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
