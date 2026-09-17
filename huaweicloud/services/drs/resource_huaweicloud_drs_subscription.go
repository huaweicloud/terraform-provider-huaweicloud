package drs

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var drsSubscriptionNonUpdatableParams = []string{
	"name",
	"description",
	"instance_type",
	"enterprise_project_id",
	"tags",
	"source_endpoint_info",
	"is_grant_new_agency",
}

// @API DRS POST /v5/{project_id}/subscription
func ResourceDrsSubscription() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDrsSubscriptionCreate,
		ReadContext:   resourceDrsSubscriptionRead,
		UpdateContext: resourceDrsSubscriptionUpdate,
		DeleteContext: resourceDrsSubscriptionDelete,

		CustomizeDiff: config.FlexibleForceNew(drsSubscriptionNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Specifies the name of the subscription task.",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Specifies the enterprise project ID.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Specifies the description of the subscription task.",
			},
			"instance_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Specifies the instance type of the subscription task.",
			},
			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Specifies the tags of the subscription task.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specifies the tag key.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specifies the tag value.",
						},
					},
				},
			},
			"source_endpoint_info": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "Specifies the source database information of the subscription task.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specifies the database instance ID.",
						},
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies the database instance type.",
						},
					},
				},
			},
			"is_grant_new_agency": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Specifies whether to create a new agency.",
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func resourceDrsSubscriptionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		httpUrl = "v5/{project_id}/subscription"
		product = "drs"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateDrsSubscriptionBodyParams(d))

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating DRS subscription: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", createRespBody, "").(string)
	if jobId == "" {
		return diag.Errorf("error creating DRS subscription: job_id is not found in API response")
	}
	d.SetId(jobId)

	return resourceDrsSubscriptionRead(ctx, d, meta)
}

func buildCreateDrsSubscriptionBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":                  d.Get("name"),
		"description":           utils.ValueIgnoreEmpty(d.Get("description")),
		"instance_type":         utils.ValueIgnoreEmpty(d.Get("instance_type")),
		"enterprise_project_id": d.Get("enterprise_project_id"),
		"tags":                  utils.ValueIgnoreEmpty(buildDrsSubscriptionTagsBodyParams(d.Get("tags"))),
		"source_endpoint_info":  buildDrsSubscriptionSourceEndpointInfoBodyParams(d.Get("source_endpoint_info")),
		"is_grant_new_agency":   utils.ValueIgnoreEmpty(d.Get("is_grant_new_agency")),
	}
	return bodyParams
}

func buildDrsSubscriptionTagsBodyParams(tagsRaw interface{}) []map[string]interface{} {
	tags := tagsRaw.([]interface{})
	if len(tags) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, 0, len(tags))
	for _, v := range tags {
		tag := v.(map[string]interface{})
		result = append(result, map[string]interface{}{
			"key":   tag["key"],
			"value": tag["value"],
		})
	}
	return result
}

func buildDrsSubscriptionSourceEndpointInfoBodyParams(sourceEndpointRaw interface{}) map[string]interface{} {
	sourceEndpointList := sourceEndpointRaw.([]interface{})
	if len(sourceEndpointList) == 0 {
		return nil
	}
	sourceEndpoint := sourceEndpointList[0].(map[string]interface{})
	return map[string]interface{}{
		"id":   sourceEndpoint["id"],
		"type": utils.ValueIgnoreEmpty(sourceEndpoint["type"]),
	}
}

func resourceDrsSubscriptionRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDrsSubscriptionUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDrsSubscriptionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting DRS subscription resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
