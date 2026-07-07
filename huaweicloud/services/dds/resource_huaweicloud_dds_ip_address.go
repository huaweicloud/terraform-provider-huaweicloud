package dds

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var ipAddressNonUpdatableParams = []string{"instance_id", "type", "password"}

// @API DDS POST /v3/{project_id}/instances/{instance_id}/create-ip
// @API DDS DELETE /v3/{project_id}/instances/{instance_id}/ip
// @API DDS GET /v3/{project_id}/instances
// @API DDS GET /v3/{project_id}/jobs
func ResourceIpAddress() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIpAddressCreate,
		ReadContext:   resourceIpAddressRead,
		UpdateContext: resourceIpAddressUpdate,
		DeleteContext: resourceIpAddressDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(ipAddressNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"target_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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

func buildCreateIpAddressBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"type":     d.Get("type"),
		"password": d.Get("password"),
	}

	return bodyParams
}

func resourceIpAddressCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
		nodeType   = d.Get("type").(string)
		targetIds  = d.Get("target_ids").(*schema.Set).List()
		httpUrl    = "v3/{project_id}/instances/{instance_id}/create-ip"
	)

	client, err := cfg.NewServiceClient("dds", region)
	if err != nil {
		return diag.Errorf("error creating DDS client: %s", err)
	}

	if len(targetIds) > 0 {
		return diag.Errorf("error creating %s IP address: parameter `target_ids` does not support in creating", nodeType)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)
	createOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		KeepResponseBody: true,
		JSONBody:         buildCreateIpAddressBodyParams(d),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating %s IP address: %s", nodeType, err)
	}

	d.SetId(instanceId)

	err = waitForAddIpAddressCompleted(ctx, client, d.Id(), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceIpAddressRead(ctx, d, meta)
}

func waitForAddIpAddressCompleted(ctx context.Context, client *golangsdk.ServiceClient, instanceId string, timeout time.Duration) error {
	stateConf := &retry.StateChangeConf{
		Pending:      []string{"updating"},
		Target:       []string{"normal"},
		Refresh:      ddsInstanceStateRefreshFunc(client, instanceId),
		Timeout:      timeout,
		Delay:        15 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for adding IP address to complete: %s ", err)
	}

	return nil
}

func GetInstanceInfo(client *golangsdk.ServiceClient, instanceId string) (interface{}, error) {
	getHttpUrl := "v3/{project_id}/instances?id={instance_id}"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)
	getOpts := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch("instances|[0]", respBody, nil), nil
}

func resourceIpAddressRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		nodeType = d.Get("type").(string)
	)

	client, err := cfg.NewServiceClient("dds", region)
	if err != nil {
		return diag.Errorf("error creating DDS client: %s", err)
	}

	instacneInfo, err := GetInstanceInfo(client, d.Id())
	if err != nil {
		// When the instance does not exist, the response HTTP status code of the query API is 400
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", instanceNotFoundCodes...),
			fmt.Sprintf("error retrieving %s IP address", nodeType))
	}

	privateIpInfo := utils.PathSearch(fmt.Sprintf("groups[?type=='%s'].nodes[].private_ip", nodeType),
		instacneInfo, make([]interface{}, 0)).([]interface{})
	if len(privateIpInfo) == 0 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, fmt.Sprintf("error retrieving %s IP address", nodeType))
	}

	return diag.FromErr(d.Set("region", region))
}

func buildUpdateIpAddressBodyParams(d *schema.ResourceData, targetId string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"type":      d.Get("type"),
		"target_id": targetId,
		"password":  d.Get("password"),
	}

	return bodyParams
}

func resourceIpAddressUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		nodeType = d.Get("type").(string)
	)

	client, err := cfg.NewServiceClient("dds", region)
	if err != nil {
		return diag.Errorf("error creating DDS client: %s", err)
	}

	if d.HasChanges("target_ids") {
		oldTargetIds, newTargetIds := d.GetChange("target_ids")
		removeTargetIds := oldTargetIds.(*schema.Set).Difference(newTargetIds.(*schema.Set))
		if len(removeTargetIds.List()) > 0 {
			return diag.Errorf("error updating %s IP address: remove the IP address operation does not support", nodeType)
		}

		targetIds := newTargetIds.(*schema.Set).Difference(oldTargetIds.(*schema.Set))
		addTargetIds := utils.ExpandToStringList(targetIds.List())
		for _, targetId := range addTargetIds {
			err = updateIpAddress(ctx, client, d, d.Id(), targetId, nodeType)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return resourceIpAddressRead(ctx, d, meta)
}

func updateIpAddress(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	instanceId, targetId, nodeType string) error {
	httpUrl := "v3/{project_id}/instances/{instance_id}/create-ip"
	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{instance_id}", instanceId)

	requestOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		KeepResponseBody: true,
		JSONBody:         buildUpdateIpAddressBodyParams(d, targetId),
	}

	_, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return fmt.Errorf("error updating %s IP address: %s", nodeType, err)
	}

	err = waitForAddIpAddressCompleted(ctx, client, d.Id(), d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return err
	}

	return nil
}

func buildDeleteIpAddressBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"type": d.Get("type"),
	}

	return bodyParams
}

func resourceIpAddressDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		nodeType = d.Get("type").(string)
		httpUrl  = "v3/{project_id}/instances/{instance_id}/ip"
	)

	client, err := cfg.NewServiceClient("dds", region)
	if err != nil {
		return diag.Errorf("error creating DDS client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		KeepResponseBody: true,
		JSONBody:         buildDeleteIpAddressBodyParams(d),
	}

	resp, err := client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error deleting %s IP addresses: %s", nodeType, err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("jobId", respBody, "").(string)
	if jobId == "" {
		return diag.Errorf("error deleting %s IP addresses: unable to find job ID from the API response", nodeType)
	}

	stateConf := &retry.StateChangeConf{
		Pending:      []string{"Running"},
		Target:       []string{"Completed"},
		Refresh:      JobStateRefreshFunc(client, jobId),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        60 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for the job (%s) completed: %s ", jobId, err)
	}

	return nil
}
