package eip

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API EIP POST /v3/{domain_id}/global-eip-segments/batch-attach-internet-bandwidths
// @API EIP POST /v3/{domain_id}/global-eip-segments/batch-detach-internet-bandwidths
// @API EIP GET /v3/{domain_id}/global-eip-segments/{global_eip_segment_id}
// @API EIP GET /v3/{domain_id}/geip/jobs/{job_id}
func ResourceSegmentBandwidthAssociate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSegmentBandwidthAssociateCreate,
		ReadContext:   resourceSegmentBandwidthAssociateRead,
		UpdateContext: resourceSegmentBandwidthAssociateUpdate,
		DeleteContext: resourceSegmentBandwidthAssociateDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew([]string{
			"global_eip_segment_id",
			"internet_bandwidth_id",
		}),

		Schema: map[string]*schema.Schema{
			"global_eip_segment_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"internet_bandwidth_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"internet_bandwidth": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildAttachSegmentBandwidthAssociateBodyParams(segmentId, bandwidthId string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"global_eip_segment_id": segmentId,
		"internet_bandwidth_id": bandwidthId,
	}

	return map[string]interface{}{
		"global_eip_segments": []map[string]interface{}{bodyParams},
	}
}

func getEipJobDetail(client *golangsdk.ServiceClient, domainId, jobId string) (interface{}, error) {
	requestPath := client.Endpoint + "v3/{domain_id}/geip/jobs/{job_id}"
	requestPath = strings.ReplaceAll(requestPath, "{domain_id}", domainId)
	requestPath = strings.ReplaceAll(requestPath, "{job_id}", jobId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving EIP job detail: %s", err)
	}

	return utils.FlattenResponse(resp)
}

func waitSegmentBandwidthAssociateJobSuccess(ctx context.Context, client *golangsdk.ServiceClient,
	timeout time.Duration, domainId, jobId string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			jobDetail, err := getEipJobDetail(client, domainId, jobId)
			if err != nil {
				return nil, "ERROR", err
			}

			status := utils.PathSearch("job.status", jobDetail, "").(string)
			if status == "" {
				return jobDetail, "ERROR", errors.New("status is not found in job detail response")
			}

			if status == "FINISH_SUCC" {
				return jobDetail, "COMPLETED", nil
			}

			// Due to the unclear API description, and for program security reasons, all other states are categorized as
			// types requiring PENDING.
			return jobDetail, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceSegmentBandwidthAssociateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		domainId    = cfg.DomainID
		httpUrl     = "v3/{domain_id}/global-eip-segments/batch-attach-internet-bandwidths"
		segmentId   = d.Get("global_eip_segment_id").(string)
		bandwidthId = d.Get("internet_bandwidth_id").(string)
	)

	client, err := cfg.NewServiceClient("geip", region)
	if err != nil {
		return diag.Errorf("error creating GEIP client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{domain_id}", domainId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildAttachSegmentBandwidthAssociateBodyParams(segmentId, bandwidthId),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error attaching EIP bandwidth (%s) to segment (%s): %s", bandwidthId, segmentId, err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", respBody, "").(string)
	if jobId == "" {
		return diag.Errorf("error attaching EIP bandwidth (%s) to segment (%s): Job ID is not found in API response",
			bandwidthId, segmentId)
	}

	d.SetId(segmentId)

	err = waitSegmentBandwidthAssociateJobSuccess(ctx, client, d.Timeout(schema.TimeoutCreate), domainId, jobId)
	if err != nil {
		return diag.Errorf("error waiting for EIP segment bandwidth association job (%s) to succeed: %s", jobId, err)
	}

	return resourceSegmentBandwidthAssociateRead(ctx, d, meta)
}

func GetEipSegmentDetail(client *golangsdk.ServiceClient, domainId, segmentId string) (interface{}, error) {
	requestPath := client.Endpoint + "v3/{domain_id}/global-eip-segments/{global_eip_segment_id}"
	requestPath = strings.ReplaceAll(requestPath, "{domain_id}", domainId)
	requestPath = strings.ReplaceAll(requestPath, "{global_eip_segment_id}", segmentId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	internetBandwidthRespBody := utils.PathSearch("global_eip_segment.internet_bandwidth", respBody, nil)
	if internetBandwidthRespBody == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return respBody, nil
}

func resourceSegmentBandwidthAssociateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		domainId  = cfg.DomainID
		segmentId = d.Id()
	)

	client, err := cfg.NewServiceClient("geip", region)
	if err != nil {
		return diag.Errorf("error creating GEIP client: %s", err)
	}

	respBody, err := GetEipSegmentDetail(client, domainId, segmentId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving EIP segment detail")
	}

	mErr := multierror.Append(
		d.Set("global_eip_segment_id", utils.PathSearch("global_eip_segment.id", respBody, nil)),
		d.Set("internet_bandwidth_id", utils.PathSearch("global_eip_segment.internet_bandwidth.id", respBody, nil)),
		d.Set("internet_bandwidth", flattenInternetBandwidth(respBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenInternetBandwidth(respBody interface{}) []interface{} {
	internetBandwidth := utils.PathSearch("global_eip_segment.internet_bandwidth", respBody, nil)
	if internetBandwidth == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"id":   utils.PathSearch("id", internetBandwidth, nil),
			"size": utils.PathSearch("size", internetBandwidth, nil),
		},
	}
}

func resourceSegmentBandwidthAssociateUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func buildDetachSegmentBandwidthAssociateBodyParams(segmentId string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"global_eip_segment_id": segmentId,
	}

	return map[string]interface{}{
		"global_eip_segments": []map[string]interface{}{bodyParams},
	}
}

func resourceSegmentBandwidthAssociateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		domainId  = cfg.DomainID
		httpUrl   = "v3/{domain_id}/global-eip-segments/batch-detach-internet-bandwidths"
		segmentId = d.Get("global_eip_segment_id").(string)
	)

	client, err := cfg.NewServiceClient("geip", region)
	if err != nil {
		return diag.Errorf("error creating GEIP client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{domain_id}", domainId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildDetachSegmentBandwidthAssociateBodyParams(segmentId),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error detaching EIP segment (%s) bandwidth: %s", segmentId, err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", respBody, "").(string)
	if jobId == "" {
		return diag.Errorf("error detaching EIP segment (%s) bandwidth: Job ID is not found in API response", segmentId)
	}

	if err := waitSegmentBandwidthAssociateJobSuccess(ctx, client, d.Timeout(schema.TimeoutDelete), domainId, jobId); err != nil {
		return diag.Errorf("error waiting for EIP segment bandwidth disassociation job (%s) to succeed: %s", jobId, err)
	}

	return nil
}
