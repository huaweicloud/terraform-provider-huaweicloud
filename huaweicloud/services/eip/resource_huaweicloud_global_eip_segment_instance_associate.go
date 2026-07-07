package eip

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API EIP POST /v3/{domain_id}/global-eip-segments/{global_eip_segment_id}/associate-instance
// @API EIP POST /v3/{domain_id}/global-eip-segments/{global_eip_segment_id}/disassociate-instance
// @API EIP GET /v3/{domain_id}/global-eip-segments/{global_eip_segment_id}
// @API EIP GET /v3/{domain_id}/geip/jobs/{job_id}
// This resource has not been tested yet
func ResourceSegmentInstanceAssociate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSegmentInstanceAssociateCreate,
		ReadContext:   resourceSegmentInstanceAssociateRead,
		UpdateContext: resourceSegmentInstanceAssociateUpdate,
		DeleteContext: resourceSegmentInstanceAssociateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew([]string{
			"global_eip_segment_id",
			"global_eip_segment",
			"global_eip_segment.*.region",
			"global_eip_segment.*.instance_id",
			"global_eip_segment.*.instance_type",
			"global_eip_segment.*.project_id",
			"global_eip_segment.*.instance_site",
			"global_eip_segment.*.service_id",
			"global_eip_segment.*.service_type",
		}),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"global_eip_segment_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"global_eip_segment": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     globalEipSegmentInfoSchema(),
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"associate_instance": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     globalEipAssociateInstanceSchema(),
			},
		},
	}
}

func globalEipSegmentInfoSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_site": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"service_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"service_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
	return &sc
}

func globalEipAssociateInstanceSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"instance_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func buildSegmentInstanceBodyParams(d *schema.ResourceData) map[string]interface{} {
	rawArray, ok := d.Get("global_eip_segment").([]interface{})
	if !ok || len(rawArray) == 0 {
		return nil
	}

	item, ok := rawArray[0].(map[string]interface{})

	if !ok {
		return nil
	}

	bodyParams := map[string]interface{}{
		"global_eip_segment": map[string]interface{}{
			"region":        item["region"],
			"instance_type": item["instance_type"],
			"instance_id":   item["instance_id"],
			"project_id":    item["project_id"],
			"instance_site": utils.ValueIgnoreEmpty(item["instance_site"]),
			"service_id":    utils.ValueIgnoreEmpty(item["service_id"]),
			"service_type":  utils.ValueIgnoreEmpty(item["service_type"]),
		},
	}

	return bodyParams
}

func resourceSegmentInstanceAssociateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                = meta.(*config.Config)
		region             = cfg.GetRegion(d)
		httpUrl            = "v3/{domain_id}/global-eip-segments/{global_eip_segment_id}/associate-instance"
		globalEipSegmentId = d.Get("global_eip_segment_id").(string)
	)

	client, err := cfg.NewServiceClient("geip", region)
	if err != nil {
		return diag.Errorf("error creating GEIP client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{domain_id}", cfg.DomainID)
	requestPath = strings.ReplaceAll(requestPath, "{global_eip_segment_id}", globalEipSegmentId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildSegmentInstanceBodyParams(d)),
	}

	createAssociateResp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating global EIP segment instance association: %s", err)
	}

	segmentInstanceAssociateRespBody, err := utils.FlattenResponse(createAssociateResp)
	if err != nil {
		return diag.FromErr(err)
	}
	jobID := utils.PathSearch("job_id", segmentInstanceAssociateRespBody, "").(string)
	if jobID == "" {
		return diag.Errorf("unable to find job ID from the API response")
	}

	// wait for job status become SUCCESS
	err = waitForJobComplete(ctx, d.Timeout(schema.TimeoutCreate), jobID, cfg.DomainID, client)
	if err != nil {
		return diag.Errorf("error waiting for global EIP associating with instance: %s", err)
	}

	d.SetId(globalEipSegmentId)

	return resourceSegmentInstanceAssociateRead(ctx, d, meta)
}

func GetGlobalsegment(client *golangsdk.ServiceClient, domainID, globalEipId string) (interface{}, error) {
	getPath := client.Endpoint + "v3/{domain_id}/global-eip-segments/{global_eip_segment_id}"
	getPath = strings.ReplaceAll(getPath, "{domain_id}", domainID)
	getPath = strings.ReplaceAll(getPath, "{global_eip_segment_id}", globalEipId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}
	target := utils.PathSearch("global_eip_segment.associate_instance", respBody, nil)
	if target == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return utils.PathSearch("global_eip_segment", respBody, nil), nil
}

func resourceSegmentInstanceAssociateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("geip", region)
	if err != nil {
		return diag.Errorf("error creating GEIP client: %s", err)
	}

	segmentInfo, err := GetGlobalsegment(client, cfg.DomainID, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving global EIP segment")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("global_eip_segment_id", utils.PathSearch("id", segmentInfo, nil)),
		d.Set("associate_instance", flattenAssociateInstanceInfo(segmentInfo)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAssociateInstanceInfo(segmentInfo interface{}) []interface{} {
	if segmentInfo == nil {
		return nil
	}

	associateInstanceInfo := utils.PathSearch("associate_instance", segmentInfo, nil)
	if associateInstanceInfo == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"instance_id":   utils.PathSearch("instance_id", associateInstanceInfo, nil),
			"instance_type": utils.PathSearch("instance_type", associateInstanceInfo, nil),
		},
	}
}

func resourceSegmentInstanceAssociateUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceSegmentInstanceAssociateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("geip", region)
	if err != nil {
		return diag.Errorf("error creating GEIP client: %s", err)
	}

	deleteGEIPHttpUrl := "v3/{domain_id}/global-eip-segments/{global_eip_segment_id}/disassociate-instance"
	deleteGEIPPath := client.Endpoint + deleteGEIPHttpUrl
	deleteGEIPPath = strings.ReplaceAll(deleteGEIPPath, "{domain_id}", cfg.DomainID)
	deleteGEIPPath = strings.ReplaceAll(deleteGEIPPath, "{global_eip_segment_id}", d.Id())

	deleteGEIPOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	deleteAssociateResp, err := client.Request("POST", deleteGEIPPath, &deleteGEIPOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertUndefinedErrInto404Err(err, 409, "error_code", "GEIP.5002"),
			"error deleting segment instance associate")
	}

	segmentInstanceAssociateRespBody, err := utils.FlattenResponse(deleteAssociateResp)
	if err != nil {
		return diag.FromErr(err)
	}
	jobID := utils.PathSearch("job_id", segmentInstanceAssociateRespBody, "").(string)
	if jobID == "" {
		return diag.Errorf("unable to find job ID from the API response")
	}

	// wait for job status become SUCCESS
	err = waitForJobComplete(ctx, d.Timeout(schema.TimeoutDelete), jobID, cfg.DomainID, client)
	if err != nil {
		return diag.Errorf("error waiting for global EIP disassociating  with instance: %s", err)
	}

	return nil
}

func waitForJobComplete(ctx context.Context, timeout time.Duration, id string, domainID string,
	client *golangsdk.ServiceClient) error {
	stateConf := &retry.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"SUCCESS"},
		Refresh:      jobRefreshFunc(client, id, domainID),
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func jobRefreshFunc(client *golangsdk.ServiceClient, jobID, domainID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getJobHttpUrl := "v3/{domain_id}/geip/jobs/{job_id}"
		getJobPath := client.Endpoint + getJobHttpUrl
		getJobPath = strings.ReplaceAll(getJobPath, "{domain_id}", domainID)
		getJobPath = strings.ReplaceAll(getJobPath, "{job_id}", jobID)
		getJobOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}

		getJobResp, err := client.Request("GET", getJobPath, &getJobOpt)
		if err != nil {
			return nil, "ERROR", err
		}
		getJobRespBody, err := utils.FlattenResponse(getJobResp)
		if err != nil {
			return nil, "ERROR", err
		}
		status := utils.PathSearch("job.status", getJobRespBody, "").(string)
		if status == "" {
			return nil, "ERROR", errors.New("unable to find job status from the API response")
		}

		if status == "FINISH_SUCC" {
			return getJobRespBody, "SUCCESS", nil
		}
		return getJobRespBody, "PENDING", nil
	}
}
