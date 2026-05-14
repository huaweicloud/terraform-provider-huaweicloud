package eip

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
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

// 1. When creating this resource, the function of binding to the global public network bandwidth was removed, and it
// can only support the creation of global elastic public IP segment.
// 2. If users need to bind to the global public network bandwidth, it is recommended to use the
// `huaweicloud_global_eip_segment_bandwidth_associate` resource

// @API EIP POST /v3/{domain_id}/global-eip-segments
// @API EIP PUT /v3/{domain_id}/global-eip-segments/{global_eip_segment_id}
// @API EIP GET /v3/{domain_id}/global-eip-segments/{global_eip_segment_id}
// @API EIP DELETE /v3/{domain_id}/global-eip-segments/{global_eip_segment_id}
// @API EIP GET /v3/{domain_id}/geip/jobs/{job_id}
// @API EIP POST /v3/global-eip-segment/{resource_id}/tags/create
// @API EIP POST /v3/global-eip-segment/{resource_id}/tags/delete

var globalEipSegmentNonUpdatableParams = []string{
	"geip_pool_name",
	"access_site",
	"mask",
	"enterprise_project_id",
}

func ResourceGlobalEipSegment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGlobalEipSegmentCreate,
		ReadContext:   resourceGlobalEipSegmentRead,
		UpdateContext: resourceGlobalEipSegmentUpdate,
		DeleteContext: resourceGlobalEipSegmentDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: customdiff.All(
			config.FlexibleForceNew(globalEipSegmentNonUpdatableParams),
		),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"geip_pool_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"access_site": {
				Type:     schema.TypeString,
				Required: true,
			},
			"mask": {
				Type:     schema.TypeInt,
				Required: true,
			},
			// The `name` can be left blank, so no `Computed` is added.
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// The `description` can be left blank, so no `Computed` is added.
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// The `tags` can be left blank, so no `Computed` is added.
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"isp": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cidr": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cidr_v6": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"freezen": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_pre_paid": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_charged": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func buildCreateGlobalEipSegmentTagsBodyParams(tagList []interface{}) []map[string]interface{} {
	rst := make([]map[string]interface{}, 0, len(tagList))
	for _, item := range tagList {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		tag := map[string]interface{}{
			"key":   m["key"],
			"value": m["value"],
		}

		rst = append(rst, tag)
	}

	return rst
}

func buildCreateGlobalEipSegmentBodyParams(d *schema.ResourceData, epsID string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"geip_pool_name":        d.Get("geip_pool_name"),
		"access_site":           d.Get("access_site"),
		"mask":                  d.Get("mask"),
		"name":                  utils.ValueIgnoreEmpty(d.Get("name")),
		"description":           utils.ValueIgnoreEmpty(d.Get("description")),
		"enterprise_project_id": utils.ValueIgnoreEmpty(epsID),
	}
	if tagList := d.Get("tags").([]interface{}); len(tagList) > 0 {
		bodyParams["tags"] = buildCreateGlobalEipSegmentTagsBodyParams(tagList)
	}

	return map[string]interface{}{
		"global_eip_segment": utils.RemoveNil(bodyParams),
	}
}

func resourceGlobalEipSegmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "geip"
		httpUrl = "v3/{domain_id}/global-eip-segments"
		epsID   = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GEIP client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{domain_id}", cfg.DomainID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildCreateGlobalEipSegmentBodyParams(d, epsID),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating global EIP segment: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("global_eip_segment.id", respBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating global EIP segment: ID is not found in API response")
	}

	d.SetId(id)

	jobId := utils.PathSearch("job_id", respBody, "").(string)
	if jobId == "" {
		return diag.Errorf("error creating global EIP segment: job_id is not found in API response")
	}

	if err = waitGlobalEipSegmentCreateJobSuccess(ctx, client, d.Timeout(schema.TimeoutCreate), cfg.DomainID, jobId); err != nil {
		return diag.Errorf("error waiting for global EIP segment create job (%s) to succeed: %s", jobId, err)
	}

	return resourceGlobalEipSegmentRead(ctx, d, meta)
}

func waitGlobalEipSegmentCreateJobSuccess(ctx context.Context, client *golangsdk.ServiceClient,
	timeout time.Duration, domainId, jobId string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			jobDetail, err := getGlobalEipSegmentJobDetail(client, domainId, jobId)
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
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func getGlobalEipSegmentJobDetail(client *golangsdk.ServiceClient, domainId, jobId string) (interface{}, error) {
	requestPath := client.Endpoint + "v3/{domain_id}/geip/jobs/{job_id}"
	requestPath = strings.ReplaceAll(requestPath, "{domain_id}", domainId)
	requestPath = strings.ReplaceAll(requestPath, "{job_id}", jobId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving global EIP segment job detail: %s", err)
	}

	return utils.FlattenResponse(resp)
}

func resourceGlobalEipSegmentRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "geip"
		httpUrl = "v3/{domain_id}/global-eip-segments/{global_eip_segment_id}"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GEIP client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{domain_id}", cfg.DomainID)
	getPath = strings.ReplaceAll(getPath, "{global_eip_segment_id}", d.Id())
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		// When the resource does not exist, the get API returns a `404` error code.
		return common.CheckDeletedDiag(d, err, "error retrieving global EIP segment")
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	seg := utils.PathSearch("global_eip_segment", respBody, nil)
	if seg == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("geip_pool_name", utils.PathSearch("geip_pool_name", seg, nil)),
		d.Set("access_site", utils.PathSearch("access_site", seg, nil)),
		d.Set("name", utils.PathSearch("name", seg, nil)),
		d.Set("description", utils.PathSearch("description", seg, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", seg, nil)),
		d.Set("tags", flattenGlobalEipSegmentTagsList(
			utils.PathSearch("tags", seg, make([]interface{}, 0)).([]interface{}))),
		d.Set("domain_id", utils.PathSearch("domain_id", seg, nil)),
		d.Set("isp", utils.PathSearch("isp", seg, nil)),
		d.Set("ip_version", utils.PathSearch("ip_version", seg, nil)),
		d.Set("cidr", utils.PathSearch("cidr", seg, nil)),
		d.Set("cidr_v6", utils.PathSearch("cidr_v6", seg, nil)),
		d.Set("freezen", utils.PathSearch("freezen", seg, false)),
		d.Set("status", utils.PathSearch("status", seg, nil)),
		d.Set("created_at", utils.PathSearch("created_at", seg, nil)),
		d.Set("updated_at", utils.PathSearch("updated_at", seg, nil)),
		d.Set("is_pre_paid", utils.PathSearch("is_pre_paid", seg, false)),
		d.Set("is_charged", utils.PathSearch("is_charged", seg, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGlobalEipSegmentTagsList(tagsResp []interface{}) []interface{} {
	if len(tagsResp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(tagsResp))
	for _, tag := range tagsResp {
		rst = append(rst, map[string]interface{}{
			"key":   utils.PathSearch("key", tag, ""),
			"value": utils.PathSearch("value", tag, ""),
		})
	}

	return rst
}

func updateGlobalEipSegmentTags(client *golangsdk.ServiceClient, d *schema.ResourceData, id string) error {
	oRaw, nRaw := d.GetChange("tags")
	oList := oRaw.([]interface{})
	nList := nRaw.([]interface{})

	manageTagsHttpUrl := "v3/global-eip-segment/{resource_id}/tags/{action}"
	manageTagsPath := client.Endpoint + manageTagsHttpUrl
	manageTagsPath = strings.ReplaceAll(manageTagsPath, "{resource_id}", id)
	manageTagsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200, 201, 204,
		},
	}

	if len(oList) > 0 {
		manageTagsOpt.JSONBody = map[string]interface{}{
			"tags": buildCreateGlobalEipSegmentTagsBodyParams(oList),
		}

		deleteTagsPath := strings.ReplaceAll(manageTagsPath, "{action}", "delete")
		_, err := client.Request("POST", deleteTagsPath, &manageTagsOpt)
		if err != nil {
			return fmt.Errorf("error deleting old tags: %s", err)
		}
	}

	if len(nList) > 0 {
		manageTagsOpt.JSONBody = map[string]interface{}{
			"tags": buildCreateGlobalEipSegmentTagsBodyParams(nList),
		}
		createTagsPath := strings.ReplaceAll(manageTagsPath, "{action}", "create")
		_, err := client.Request("POST", createTagsPath, &manageTagsOpt)
		if err != nil {
			return fmt.Errorf("error creating new tags: %s", err)
		}
	}

	return nil
}

func buildUpdateGlobalEipSegmentBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"global_eip_segment": map[string]interface{}{
			"name":        d.Get("name"),
			"description": d.Get("description"),
		},
	}
}

func resourceGlobalEipSegmentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "geip"
		httpUrl = "v3/{domain_id}/global-eip-segments/{global_eip_segment_id}"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GEIP client: %s", err)
	}

	if d.HasChanges("name", "description") {
		updatePath := client.Endpoint + httpUrl
		updatePath = strings.ReplaceAll(updatePath, "{domain_id}", cfg.DomainID)
		updatePath = strings.ReplaceAll(updatePath, "{global_eip_segment_id}", d.Id())
		updateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         buildUpdateGlobalEipSegmentBodyParams(d),
		}

		resp, err := client.Request("PUT", updatePath, &updateOpt)
		if err != nil {
			return diag.Errorf("error updating global EIP segment: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		jobID := utils.PathSearch("job_id", respBody, "").(string)
		if jobID != "" {
			if err = waitGlobalEipSegmentCreateJobSuccess(ctx, client, d.Timeout(schema.TimeoutUpdate), cfg.DomainID, jobID); err != nil {
				return diag.Errorf("error waiting for global EIP segment update job (%s) to succeed: %s", jobID, err)
			}
		}
	}

	if d.HasChange("tags") {
		if tagErr := updateGlobalEipSegmentTags(client, d, d.Id()); tagErr != nil {
			return diag.Errorf("error updating tags of global EIP segment (%s): %s", d.Id(), tagErr)
		}
	}

	return resourceGlobalEipSegmentRead(ctx, d, meta)
}

func resourceGlobalEipSegmentDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "geip"
		httpUrl = "v3/{domain_id}/global-eip-segments/{global_eip_segment_id}"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GEIP client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{domain_id}", cfg.DomainID)
	deletePath = strings.ReplaceAll(deletePath, "{global_eip_segment_id}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200, 201, 204,
		},
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error deleting global EIP segment: %s", err)
	}

	return nil
}
