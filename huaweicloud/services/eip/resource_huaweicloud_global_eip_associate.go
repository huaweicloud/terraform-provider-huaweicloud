package eip

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API EIP POST /v3/{domain_id}/global-eips/{global_eip_id}/associate-instance
// @API EIP POST /v3/{domain_id}/global-eips/{global_eip_id}/disassociate-instance
// @API EIP GET /v3/{domain_id}/geip/jobs/{job_id}
// @API EIP GET /v3/{domain_id}/global-eips/{id}
// @API CC GET /v3/{domain_id}/gcb/gcbandwidths/{id}
// @API CC PUT /v3/{domain_id}/gcb/gcbandwidths/{id}
func ResourceGlobalEIPAssociate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGlobalEIPAssociateCreate,
		ReadContext:   resourceGlobalEIPAssociateRead,
		UpdateContext: resourceGlobalEIPAssociateUpdate,
		DeleteContext: resourceGlobalEIPAssociateDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"global_eip_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"associate_instance": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"project_id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"instance_type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"service_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"service_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
					},
				},
			},
			"gc_bandwidth": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"charge_mode": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"enterprise_project_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"tags": common.TagsForceNewSchema(),
					},
				},
			},
			"is_reserve_gcb": {
				Type:     schema.TypeBool,
				Required: true,
			},
		},
	}
}

func resourceGlobalEIPAssociateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("geip", region)
	if err != nil {
		return diag.Errorf("error creating GEIP client: %s", err)
	}

	createGEIPAssociateHttpUrl := "v3/{domain_id}/global-eips/{global_eip_id}/associate-instance"
	createGEIPAssociatePath := client.Endpoint + createGEIPAssociateHttpUrl
	createGEIPAssociatePath = strings.ReplaceAll(createGEIPAssociatePath, "{domain_id}", cfg.DomainID)
	createGEIPAssociatePath = strings.ReplaceAll(createGEIPAssociatePath, "{global_eip_id}", d.Get("global_eip_id").(string))
	createGEIPAssociateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"global_eip": utils.RemoveNil(buildCreateGEIPAssociateBodyParams(d, cfg.GetEnterpriseProjectID(d))),
		},
	}

	createGEIPAssociateResp, err := client.Request("POST", createGEIPAssociatePath, &createGEIPAssociateOpt)
	if err != nil {
		return diag.FromErr(err)
	}
	createGEIPAssociateRespBody, err := utils.FlattenResponse(createGEIPAssociateResp)
	if err != nil {
		return diag.FromErr(err)
	}
	jobID := utils.PathSearch("job_id", createGEIPAssociateRespBody, "").(string)
	if jobID == "" {
		return diag.Errorf("unable to find job ID from the API response")
	}

	// wait for job status become SUCCESS
	err = waitForJobStatusComplete(ctx, d.Timeout(schema.TimeoutCreate), jobID, cfg.DomainID, client)
	if err != nil {
		return diag.Errorf("error waiting for global EIP associating with instance: %s", err)
	}

	d.SetId(d.Get("global_eip_id").(string))

	// wait for GEIP update complete
	geip, err := waitForGEIPCompleteWithRespBody(ctx, d.Timeout(schema.TimeoutCreate),
		d.Get("global_eip_id").(string), cfg.DomainID, client)
	if err != nil {
		return diag.Errorf("error waiting for global EIP associating with instance: %s", err)
	}
	gcbID := utils.PathSearch("global_eip.global_connection_bandwidth_info.gcb_id", geip, "").(string)
	if gcbID == "" {
		if v, ok := d.GetOk("gc_bandwidth"); ok && len(v.([]interface{})) > 0 {
			return diag.Errorf("unable to find global connection bandwidth ID from the API response")
		}
	}

	// if bandwidth charge_mode is not "bwd", call Update GCB
	ccClient, err := cfg.NewServiceClient("cc", region)
	if err != nil {
		return diag.Errorf("error creating CC client: %s", err)
	}
	if v, ok := d.GetOk("gc_bandwidth"); ok && len(v.([]interface{})) > 0 {
		gcb := d.Get("gc_bandwidth").([]interface{})[0].(map[string]interface{})
		if gcb["id"].(string) == "" && gcb["charge_mode"].(string) == "95" {
			err = updateGCB(ccClient, gcbID, cfg.DomainID)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return resourceGlobalEIPAssociateRead(ctx, d, meta)
}

func buildCreateGEIPAssociateBodyParams(d *schema.ResourceData, epsID string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"associate_instance_info": buildGEIPAssociateRequestBodyAssociateInstance(d.Get("associate_instance")),
		"gc_bandwidth_info":       buildGEIPAssociateRequestBodyGCB(d.Get("gc_bandwidth"), epsID),
	}
	return bodyParams
}

func buildGEIPAssociateRequestBodyAssociateInstance(rawParams interface{}) map[string]interface{} {
	rawArray, _ := rawParams.([]interface{})
	raw := rawArray[0].(map[string]interface{})
	params := map[string]interface{}{
		"region":        raw["region"],
		"project_id":    raw["project_id"],
		"instance_id":   raw["instance_id"],
		"instance_type": raw["instance_type"],
		"service_id":    utils.ValueIgnoreEmpty(raw["service_id"]),
		"service_type":  utils.ValueIgnoreEmpty(raw["service_type"]),
	}
	return params
}

func buildGEIPAssociateRequestBodyGCB(rawParams interface{}, epsID string) map[string]interface{} {
	rawArray, _ := rawParams.([]interface{})
	if len(rawArray) < 1 {
		return nil
	}
	raw := rawArray[0].(map[string]interface{})
	if raw["id"].(string) == "" {
		return map[string]interface{}{
			"type":                  "Region",
			"name":                  raw["name"],
			"charge_mode":           "bwd",
			"size":                  raw["size"],
			"description":           utils.ValueIgnoreEmpty(raw["description"]),
			"tags":                  utils.ValueIgnoreEmpty(utils.ExpandResourceTags(raw["tags"].(map[string]interface{}))),
			"enterprise_project_id": utils.ValueIgnoreEmpty(epsID),
		}
	}
	return map[string]interface{}{
		"id": raw["id"],
	}
}

func waitForJobStatusComplete(ctx context.Context, timeout time.Duration, id string, domainID string,
	client *golangsdk.ServiceClient) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"SUCCESS"},
		Refresh:      jobStatusRefreshFunc(client, id, domainID),
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func jobStatusRefreshFunc(client *golangsdk.ServiceClient, jobID, domainID string) resource.StateRefreshFunc {
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
			return nil, "ERROR", fmt.Errorf("unable to find job status from the API response")
		}

		// status are FINISH_ROLLBACK_SUCC, FINISH_SUCC and WAIT_TO_SCHEDULE
		if status == "FINISH_ROLLBACK_SUCC" {
			return nil, "FAILURE", fmt.Errorf("job fail: %s", utils.PathSearch("job.error_message", getJobRespBody, nil))
		} else if status == "FINISH_SUCC" {
			return getJobRespBody, "SUCCESS", nil
		}
		return getJobRespBody, "PENDING", nil
	}
}

func waitForGEIPCompleteWithRespBody(ctx context.Context, timeout time.Duration, id string, domainID string,
	client *golangsdk.ServiceClient) (interface{}, error) {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"SUCCESS"},
		Refresh:      geipStatusRefreshFunc(id, domainID, client),
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	return stateConf.WaitForStateContext(ctx)
}

func updateGCB(client *golangsdk.ServiceClient, id, domainID string) error {
	updateGCBHttpUrl := "v3/{domain_id}/gcb/gcbandwidths/{id}"
	updateGCBPath := client.Endpoint + updateGCBHttpUrl
	updateGCBPath = strings.ReplaceAll(updateGCBPath, "{domain_id}", domainID)
	updateGCBPath = strings.ReplaceAll(updateGCBPath, "{id}", id)
	updateGCBOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"globalconnection_bandwidth": map[string]interface{}{
				"charge_mode": "95",
			},
		},
	}

	_, err := client.Request("PUT", updateGCBPath, &updateGCBOpt)
	if err != nil {
		return fmt.Errorf("error updating global connection bandwidth: %s", err)
	}

	return nil
}

func resourceGlobalEIPAssociateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("geip", region)
	if err != nil {
		return diag.Errorf("error creating GEIP client: %s", err)
	}

	getGEIPHttpUrl := "v3/{domain_id}/global-eips/{id}"
	getGEIPPath := client.Endpoint + getGEIPHttpUrl
	getGEIPPath = strings.ReplaceAll(getGEIPPath, "{domain_id}", cfg.DomainID)
	getGEIPPath = strings.ReplaceAll(getGEIPPath, "{id}", d.Id())
	getGEIPOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getGEIPResp, err := client.Request("GET", getGEIPPath, &getGEIPOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving global EIP")
	}
	getGEIPRespBody, err := utils.FlattenResponse(getGEIPResp)
	if err != nil {
		return diag.FromErr(err)
	}
	geip := utils.PathSearch("global_eip", getGEIPRespBody, nil)
	if geip == nil {
		return diag.Errorf("unable to find global EIP from the API response")
	}

	// Call GET GCB API to get more info, because charge_mode is not in return.
	ccClient, err := cfg.NewServiceClient("cc", region)
	if err != nil {
		return diag.Errorf("error creating CC client: %s", err)
	}
	gcbID := utils.PathSearch("global_connection_bandwidth_info.gcb_id", geip, "").(string)
	var gcbInfo interface{}
	if gcbID != "" {
		gcbInfo, err = getGCBInfo(ccClient, cfg.DomainID, gcbID)
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		log.Printf("[WARN] the global EIP(%s) is not associating with global connection bandwidth", d.Id())
	}

	mErr := multierror.Append(nil,
		d.Set("global_eip_id", d.Id()),
		d.Set("associate_instance", flattenAssociateInstance(utils.PathSearch("associate_instance_info", geip, nil))),
		d.Set("gc_bandwidth", gcbInfo),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getGCBInfo(client *golangsdk.ServiceClient, domainID, gcbID string) (interface{}, error) {
	getGCBHttpUrl := "v3/{domain_id}/gcb/gcbandwidths/{id}"
	getGCBPath := client.Endpoint + getGCBHttpUrl
	getGCBPath = strings.ReplaceAll(getGCBPath, "{domain_id}", domainID)
	getGCBPath = strings.ReplaceAll(getGCBPath, "{id}", gcbID)
	getGCBOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getGCBResp, err := client.Request("GET", getGCBPath, &getGCBOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving GCB: %s", err)
	}
	getGCBRespBody, err := utils.FlattenResponse(getGCBResp)
	if err != nil {
		return nil, fmt.Errorf("error flattening GCB: %s", err)
	}
	gcb := utils.PathSearch("globalconnection_bandwidth", getGCBRespBody, nil)
	if gcb == nil {
		return nil, fmt.Errorf("unable to find global connection bandwidth from the API response")
	}

	result := make([]interface{}, 0)
	result = append(result, map[string]interface{}{
		"id":                    gcbID,
		"name":                  utils.PathSearch("name", gcb, nil),
		"size":                  int(utils.PathSearch("size", gcb, float64(0)).(float64)),
		"charge_mode":           utils.PathSearch("charge_mode", gcb, nil),
		"description":           utils.PathSearch("description", gcb, nil),
		"tags":                  utils.FlattenTagsToMap(utils.PathSearch("tags", gcb, nil)),
		"enterprise_project_id": utils.PathSearch("enterprise_project_id", gcb, nil),
	})

	return result, nil
}

func flattenAssociateInstance(rawAttrs interface{}) []interface{} {
	curAttrs := rawAttrs.(map[string]interface{})
	rst := make([]interface{}, 0)
	rst = append(rst, map[string]interface{}{
		"region":        curAttrs["region"],
		"project_id":    curAttrs["project_id"],
		"instance_type": curAttrs["instance_type"],
		"instance_id":   curAttrs["instance_id"],
		"service_id":    curAttrs["service_id"],
		"service_type":  curAttrs["service_type"],
	})
	return rst
}

func resourceGlobalEIPAssociateUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceGlobalEIPAssociateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("geip", region)
	if err != nil {
		return diag.Errorf("error creating GEIP client: %s", err)
	}

	deleteGEIPHttpUrl := "v3/{domain_id}/global-eips/{global_eip_id}/disassociate-instance"
	deleteGEIPPath := client.Endpoint + deleteGEIPHttpUrl
	deleteGEIPPath = strings.ReplaceAll(deleteGEIPPath, "{domain_id}", cfg.DomainID)
	deleteGEIPPath = strings.ReplaceAll(deleteGEIPPath, "{global_eip_id}", d.Id())

	deleteGEIPOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"is_reserve_gcb": strconv.FormatBool(d.Get("is_reserve_gcb").(bool)),
		},
	}

	deleteGEIPAssociateResp, err := client.Request("POST", deleteGEIPPath, &deleteGEIPOpt)
	if err != nil {
		return diag.FromErr(err)
	}
	deleteGEIPAssociateRespBody, err := utils.FlattenResponse(deleteGEIPAssociateResp)
	if err != nil {
		return diag.FromErr(err)
	}
	jobID := utils.PathSearch("job_id", deleteGEIPAssociateRespBody, "").(string)
	if jobID == "" {
		return diag.Errorf("unable to find job ID from the API response")
	}

	err = waitForJobStatusComplete(ctx, d.Timeout(schema.TimeoutDelete), jobID, cfg.DomainID, client)
	if err != nil {
		return diag.Errorf("error waiting for global EIP disassociating with instance: %s", err)
	}

	// wait for GEIP update
	err = waitForGEIPComplete(ctx, d.Timeout(schema.TimeoutDelete), d.Id(), cfg.DomainID, client)
	if err != nil {
		return diag.Errorf("error waiting for global EIP disassociating with instance: %s", err)
	}

	return nil
}
