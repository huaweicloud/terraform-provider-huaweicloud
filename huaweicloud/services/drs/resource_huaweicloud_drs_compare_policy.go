package drs

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var comparePolicyNonUpdatableParams = []string{
	"job_id",
}

// @API DRS PUT /v5/{project_id}/jobs/{job_id}/compare-policy
// @API DRS GET /v5/{project_id}/jobs/{job_id}/compare-policy
func ResourceDrsComparePolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDrsComparePolicyCreate,
		ReadContext:   resourceDrsComparePolicyRead,
		UpdateContext: resourceDrsComparePolicyUpdate,
		DeleteContext: resourceDrsComparePolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(comparePolicyNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"job_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"period": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"begin_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"compare_type": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"compare_policy": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"interval_hour": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"next_compare_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func convertToStringList(rawArray []interface{}) interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	return utils.ExpandToStringList(rawArray)
}

func buildComparePolicyBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"action":         "open",
		"period":         utils.ValueIgnoreEmpty(d.Get("period")),
		"begin_time":     utils.ValueIgnoreEmpty(d.Get("begin_time")),
		"end_time":       utils.ValueIgnoreEmpty(d.Get("end_time")),
		"compare_type":   convertToStringList(d.Get("compare_type").([]interface{})),
		"compare_policy": utils.ValueIgnoreEmpty(d.Get("compare_policy")),
		"interval_hour":  utils.ValueIgnoreEmpty(d.Get("interval_hour")),
	}
}

func updateComparePolicy(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	requestPath := client.Endpoint + "v5/{project_id}/jobs/{job_id}/compare-policy"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{job_id}", d.Get("job_id").(string))
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildComparePolicyBodyParams(d)),
	}

	_, err := client.Request("PUT", requestPath, &requestOpt)
	return err
}

func resourceDrsComparePolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		jobId  = d.Get("job_id").(string)
	)

	client, err := cfg.NewServiceClient("drs", region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	if err := updateComparePolicy(client, d); err != nil {
		return diag.Errorf("error configing DRS compare policy in create operation: %s", err)
	}

	d.SetId(jobId)

	return resourceDrsComparePolicyRead(ctx, d, meta)
}

func GetComparePolicy(client *golangsdk.ServiceClient, jobId string) (interface{}, error) {
	requestPath := client.Endpoint + "v5/{project_id}/jobs/{job_id}/compare-policy"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{job_id}", jobId)
	reqOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	resp, err := client.Request("GET", requestPath, &reqOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	status := utils.PathSearch("status", respBody, "").(string)
	if status == "CLOSED" {
		return nil, golangsdk.ErrDefault404{}
	}

	return respBody, nil
}

func resourceDrsComparePolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "drs"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	respBody, err := GetComparePolicy(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DRS compare policy")
	}

	mErr := multierror.Append(
		nil,
		d.Set("region", region),
		d.Set("job_id", d.Id()),
		d.Set("period", utils.PathSearch("period", respBody, nil)),
		d.Set("begin_time", utils.PathSearch("begin_time", respBody, nil)),
		d.Set("end_time", utils.PathSearch("end_time", respBody, nil)),
		d.Set("compare_type", utils.PathSearch("compare_type", respBody, nil)),
		d.Set("compare_policy", utils.PathSearch("compare_policy", respBody, nil)),
		d.Set("interval_hour", utils.PathSearch("interval_hour", respBody, nil)),
		d.Set("status", utils.PathSearch("status", respBody, nil)),
		d.Set("next_compare_time", utils.PathSearch("next_compare_time", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDrsComparePolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("drs", region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	if err := updateComparePolicy(client, d); err != nil {
		return diag.Errorf("error configing DRS compare policy in update operation: %s", err)
	}

	return resourceDrsComparePolicyRead(ctx, d, meta)
}

func resourceDrsComparePolicyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/{project_id}/jobs/{job_id}/compare-policy"
	)

	client, err := cfg.NewServiceClient("drs", region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{job_id}", d.Id())

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody: map[string]interface{}{
			"action": "closed",
		},
	}

	_, err = client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error closing DRS compare policy: %s", err)
	}

	return nil
}
