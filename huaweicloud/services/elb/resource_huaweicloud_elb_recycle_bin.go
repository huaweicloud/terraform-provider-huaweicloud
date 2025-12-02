package elb

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ELB PUT /v3/{project_id}/elb/recycle-bin
// @API ELB PUT /v3/{project_id}/elb/recycle-bin/policy
// @API ELB GET /v3/{project_id}/elb/recycle-bin
func ResourceElbRecycleBin() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceElbRecycleBinCreate,
		ReadContext:   resourceElbRecycleBinRead,
		UpdateContext: resourceElbRecycleBinUpdate,
		DeleteContext: resourceElbRecycleBinDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"retention_hour": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: func(v interface{}, k string) ([]string, []error) {
					r := v.(string)
					if _, err := strconv.Atoi(r); err != nil {
						return nil, []error{fmt.Errorf("invalid value for %s (%s)", k, v)}
					}
					return nil, nil
				},
			},
			"recycle_threshold_hour": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: func(v interface{}, k string) ([]string, []error) {
					r := v.(string)
					if _, err := strconv.Atoi(r); err != nil {
						return nil, []error{fmt.Errorf("invalid value for %s (%s)", k, v)}
					}
					return nil, nil
				},
			},
		},
	}
}

func resourceElbRecycleBinCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		product = "elb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	err = updateRecycleBin(client, true)
	if err != nil {
		return diag.Errorf("error opening ELB recycle bin: %s", err)
	}

	d.SetId(client.ProjectID)

	_, retentionHourOk := d.GetOk("retention_hour")
	_, recycleThresholdHourOk := d.GetOk("recycle_threshold_hour")
	if retentionHourOk || recycleThresholdHourOk {
		err = updateRecycleBinPolicy(client, d)
		if err != nil {
			return diag.Errorf("error updating ELB recycle bin policy: %s", err)
		}
	}

	return resourceElbRecycleBinRead(ctx, d, meta)
}

func resourceElbRecycleBinRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v3/{project_id}/elb/recycle-bin"
		product = "elb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving ELB recycle bin policy")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	enable := utils.PathSearch("recycle_bin.enable", getRespBody, false).(bool)
	if !enable {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving ELB recycle bin")
	}

	retentionHour := utils.PathSearch("recycle_bin.policy.retention_hour", getRespBody, float64(0)).(float64)
	recycleThresholdHour := utils.PathSearch("recycle_bin.policy.recycle_threshold_hour", getRespBody, float64(0)).(float64)
	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("retention_hour", strconv.Itoa(int(retentionHour))),
		d.Set("recycle_threshold_hour", strconv.Itoa(int(recycleThresholdHour))),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceElbRecycleBinUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		product = "elb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	if d.HasChanges("retention_hour", "recycle_threshold_hour") {
		err = updateRecycleBinPolicy(client, d)
		if err != nil {
			return diag.Errorf("error updating ELB recycle bin policy: %s", err)
		}
	}
	return resourceElbRecycleBinRead(ctx, d, meta)
}

func resourceElbRecycleBinDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		product = "elb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	err = updateRecycleBin(client, false)
	if err != nil {
		return diag.Errorf("error closing ELB recycle bin: %s", err)
	}

	return nil
}

func updateRecycleBin(client *golangsdk.ServiceClient, enable bool) error {
	httpUrl := "v3/{project_id}/elb/recycle-bin"

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateOpt.JSONBody = buildUpdateRecycleBodyParams(enable)

	_, err := client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return err
	}
	return nil
}

func buildUpdateRecycleBodyParams(enable bool) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"enable": enable,
	}
	return map[string]interface{}{
		"recycle_bin": bodyParams,
	}
}

func updateRecycleBinPolicy(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	httpUrl := "v3/{project_id}/elb/recycle-bin/policy"

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateOpt.JSONBody = buildRecyclePolicyBodyParams(d)

	_, err := client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return err
	}
	return nil
}

func buildRecyclePolicyBodyParams(d *schema.ResourceData) map[string]interface{} {
	retentionHour, retentionHourOk := d.GetOk("retention_hour")
	recycleThresholdHour, recycleThresholdHourOk := d.GetOk("recycle_threshold_hour")
	bodyParams := make(map[string]interface{})
	if retentionHourOk {
		r, _ := strconv.Atoi(retentionHour.(string))
		bodyParams["retention_hour"] = r
	}
	if recycleThresholdHourOk {
		r, _ := strconv.Atoi(recycleThresholdHour.(string))
		bodyParams["recycle_threshold_hour"] = r
	}
	return map[string]interface{}{
		"recycle_bin": map[string]interface{}{
			"policy": bodyParams,
		},
	}
}
