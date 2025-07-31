package secmaster

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var nonUpdatableParamsOrder = []string{"operate_type", "tags", "product_list",
	"product_list.*.id",
	"product_list.*.product_id",
	"product_list.*.cloud_service_type",
	"product_list.*.resource_type",
	"product_list.*.resource_spec_code",
	"product_list.*.usage_measure_id",
	"product_list.*.usage_value",
	"product_list.*.resource_size",
	"product_list.*.usage_factor",
	"product_list.*.resource_id",
}

// @API SecMaster POST /v1/{project_id}/subscriptions/orders
func ResourcePostPaidOrder() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePostPaidOrderCreate,
		UpdateContext: resourcePostPaidOrderUpdate,
		ReadContext:   resourcePostPaidOrderRead,
		DeleteContext: resourcePostPaidOrderDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: customdiff.All(
			config.FlexibleForceNew(nonUpdatableParamsOrder),
			config.MergeDefaultTags(),
		),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"operate_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": common.TagsSchema(),
			"product_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"product_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"cloud_service_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"resource_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"resource_spec_code": {
							Type:     schema.TypeString,
							Required: true,
						},
						"usage_measure_id": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"usage_value": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"resource_size": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"usage_factor": {
							Type:     schema.TypeString,
							Required: true,
						},
						"resource_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
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

func resourcePostPaidOrderCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		createPostPaidOrderHttpUrl = "v1/{project_id}/subscriptions/orders"
		createPostPaidOrderProduct = "secmaster"
	)
	createPostPaidOrderClient, err := cfg.NewServiceClient(createPostPaidOrderProduct, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	createPostPaidOrderPath := createPostPaidOrderClient.Endpoint + createPostPaidOrderHttpUrl
	createPostPaidOrderPath = strings.ReplaceAll(createPostPaidOrderPath, "{project_id}", createPostPaidOrderClient.ProjectID)
	createPostPaidOrderOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"X-Language": "en-us"},
	}

	createOpts := map[string]interface{}{
		"region_id":    region,
		"domain_id":    cfg.DomainID,
		"tags":         utils.ExpandResourceTagsMap(d.Get("tags").(map[string]interface{})),
		"limit":        utils.ValueIgnoreEmpty(d.Get("limit")),
		"product_list": buildProductListParams(d),
	}

	createPostPaidOrderOpt.JSONBody = utils.RemoveNil(createOpts)

	_, err = createPostPaidOrderClient.Request("POST", createPostPaidOrderPath, &createPostPaidOrderOpt)
	if err != nil {
		return diag.Errorf("error creating SecMaster post paid order: %s", err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = createPostPaidOrderWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the SecMaster post paid order (%s) creation to complete: %s", d.Id(), err)
	}

	return nil
}

func buildProductListParams(d *schema.ResourceData) []map[string]interface{} {
	v, ok := d.GetOk("product_list")
	if !ok {
		return nil
	}

	productList := v.([]interface{})

	res := make([]map[string]interface{}, len(productList))

	for i, v := range productList {
		res[i] = map[string]interface{}{
			"id":                 utils.PathSearch("id", v, nil),
			"product_id":         utils.PathSearch("product_id", v, nil),
			"cloud_service_type": utils.PathSearch("cloud_service_type", v, nil),
			"resource_type":      utils.PathSearch("resource_type", v, nil),
			"resource_spec_code": utils.PathSearch("resource_spec_code", v, nil),
			"usage_measure_id":   utils.PathSearch("usage_measure_id", v, nil),
			"usage_value":        utils.PathSearch("usage_value", v, nil),
			"resource_size":      utils.PathSearch("resource_size", v, nil),
			"usage_factor":       utils.PathSearch("usage_factor", v, nil),
			"resource_id":        utils.ValueIgnoreEmpty(utils.PathSearch("resource_id", v, nil)),
		}
	}

	return res
}

func resourcePostPaidOrderRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourcePostPaidOrderUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourcePostPaidOrderDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for SecMaster post paid order resource.
		Deleting this resource will not change the status of the currently SecMaster post paid order resource,
		but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

func createPostPaidOrderWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			cfg := meta.(*config.Config)
			region := cfg.GetRegion(d)

			var (
				getSubscriptionsHttpUrl = "v1/{project_id}/subscriptions/version"
				getSubscriptionsProduct = "secmaster"
			)
			getSubscriptionsClient, err := cfg.NewServiceClient(getSubscriptionsProduct, region)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error creating SecMaster client: %s", err)
			}

			getSubscriptionsPath := getSubscriptionsClient.Endpoint + getSubscriptionsHttpUrl
			getSubscriptionsPath = strings.ReplaceAll(getSubscriptionsPath, "{project_id}", getSubscriptionsClient.ProjectID)

			getSubscriptionsOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
				MoreHeaders: map[string]string{
					"Content-Type": "application/json",
					"X-Language":   "en-us",
				},
			}

			getSubscriptionsResp, err := getSubscriptionsClient.Request("GET", getSubscriptionsPath, &getSubscriptionsOpt)
			if err != nil {
				return nil, "ERROR", err
			}

			getSubscriptionsRespBody, err := utils.FlattenResponse(getSubscriptionsResp)
			if err != nil {
				return nil, "ERROR", err
			}
			csbVersion := utils.PathSearch(`csb_version`, getSubscriptionsRespBody, nil)
			if csbVersion == nil {
				return nil, "ERROR", fmt.Errorf("error parsing %s from response body", `csb_version`)
			}

			version := fmt.Sprintf("%v", csbVersion)

			if version != "NA" {
				return getSubscriptionsRespBody, "COMPLETED", nil
			}

			return getSubscriptionsRespBody, "PENDING", nil
		},
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}
