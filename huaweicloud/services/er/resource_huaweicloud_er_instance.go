// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product ER
// ---------------------------------------------------------------

package er

import (
	"context"
	"fmt"
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

// @API ER POST /v3/{project_id}/enterprise-router/instances
// @API ER PUT /v3/{project_id}/enterprise-router/instances/{er_id}
// @API ER POST /v3/{project_id}/{resource_type}/{resource_id}/tags/action
// @API ER GET /v3/{project_id}/enterprise-router/instances/{er_id}
// @API ER POST /v3/{project_id}/enterprise-router/instances/{er_id}/change-availability-zone-ids
// @API ER DELETE /v3/{project_id}/enterprise-router/instances/{er_id}
func ResourceInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceInstanceCreate,
		UpdateContext: resourceInstanceUpdate,
		ReadContext:   resourceInstanceRead,
		DeleteContext: resourceInstanceDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		CustomizeDiff: config.MergeDefaultTags(),

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
				Description: `The router name.`,
			},
			"availability_zones": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				Description: `The availability zone list where the Enterprise router is located.`,
			},
			"asn": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: `The BGP AS number of the Enterprise router.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the Enterprise router.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The enterprise project ID to which the Enterprise router belongs.`,
			},
			"enable_default_propagation": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to enable the propagation of the default route table.`,
			},
			"enable_default_association": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to enable the association of the default route table.`,
			},
			"auto_accept_shared_attachments": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to automatically accept the creation of shared attachment.`,
			},
			"default_propagation_route_table_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The ID of the default propagation route table.`,
			},
			"default_association_route_table_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The ID of the default association route table.`,
			},
			"tags": common.TagsSchema(),
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Current status of the router.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The latest update time.`,
			},
		},
	}
}

func resourceInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/enterprise-router/instances"
	)
	client, err := cfg.NewServiceClient("er", region)
	if err != nil {
		return diag.Errorf("error creating ER Client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createInstanceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateInstanceBodyParams(d, cfg)),
	}

	createInstanceResp, err := client.Request("POST", createPath, &createInstanceOpt)
	if err != nil {
		return diag.Errorf("error creating instance: %s", err)
	}

	createInstanceRespBody, err := utils.FlattenResponse(createInstanceResp)
	if err != nil {
		return diag.FromErr(err)
	}

	instanceId := utils.PathSearch("instance.id", createInstanceRespBody, "").(string)
	if instanceId == "" {
		return diag.Errorf("unable to find the ER instance ID from the API response")
	}
	d.SetId(instanceId)

	err = instanceWaitingForStateCompleted(ctx, client, d, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the create operation of the instance (%s) to complete: %s", d.Id(), err)
	}

	if _, ok := d.GetOk("tags"); ok {
		err = utils.UpdateResourceTags(client, d, "instance", d.Id())
		if err != nil {
			return diag.Errorf("error creating instance tags: %s", err)
		}
	}

	return resourceInstanceRead(ctx, d, meta)
}

func buildCreateInstanceBodyParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"instance": buildCreateInstanceInstanceChildBody(d, cfg),
	}
	return bodyParams
}

func buildCreateInstanceInstanceChildBody(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	params := map[string]interface{}{
		"name":                           utils.ValueIgnoreEmpty(d.Get("name")),
		"availability_zone_ids":          utils.ValueIgnoreEmpty(d.Get("availability_zones")),
		"asn":                            utils.ValueIgnoreEmpty(d.Get("asn")),
		"description":                    utils.ValueIgnoreEmpty(d.Get("description")),
		"enterprise_project_id":          utils.ValueIgnoreEmpty(cfg.GetEnterpriseProjectID(d)),
		"enable_default_propagation":     utils.ValueIgnoreEmpty(d.Get("enable_default_propagation")),
		"enable_default_association":     utils.ValueIgnoreEmpty(d.Get("enable_default_association")),
		"auto_accept_shared_attachments": utils.ValueIgnoreEmpty(d.Get("auto_accept_shared_attachments")),
	}
	return params
}

func instanceWaitingForStateCompleted(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			httpUrl := "v3/{project_id}/enterprise-router/instances/{er_id}"
			refreshPath := client.Endpoint + httpUrl
			refreshPath = strings.ReplaceAll(refreshPath, "{project_id}", client.ProjectID)
			refreshPath = strings.ReplaceAll(refreshPath, "{er_id}", d.Id())

			createInstanceWaitingOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
				OkCodes: []int{
					200,
				},
			}
			requestResp, err := client.Request("GET", refreshPath, &createInstanceWaitingOpt)
			if err != nil {
				return nil, "ERROR", err
			}

			createInstanceWaitingRespBody, err := utils.FlattenResponse(requestResp)
			if err != nil {
				return nil, "ERROR", err
			}
			status := utils.PathSearch(`instance.state`, createInstanceWaitingRespBody, "").(string)

			if utils.StrSliceContains([]string{"fail"}, status) {
				return createInstanceWaitingRespBody, "", fmt.Errorf("unexpected status '%s'", status)
			}

			if utils.StrSliceContains([]string{"available"}, status) {
				return createInstanceWaitingRespBody, "COMPLETED", nil
			}

			return createInstanceWaitingRespBody, "PENDING", nil
		},
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceInstanceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/enterprise-router/instances/{er_id}"
	)
	client, err := cfg.NewServiceClient("er", region)
	if err != nil {
		return diag.Errorf("error creating ER Client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{er_id}", d.Id())

	getInstanceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getInstanceResp, err := client.Request("GET", getPath, &getInstanceOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving instance")
	}
	getInstanceRespBody, err := utils.FlattenResponse(getInstanceResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("instance.name", getInstanceRespBody, nil)),
		d.Set("description", utils.PathSearch("instance.description", getInstanceRespBody, nil)),
		d.Set("status", utils.PathSearch("instance.state", getInstanceRespBody, nil)),
		// The time results are not the time in RF3339 format without milliseconds.
		d.Set("created_at", utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("instance.created_at",
			getInstanceRespBody, "").(string))/1000, false)),
		d.Set("updated_at", utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("instance.updated_at",
			getInstanceRespBody, "").(string))/1000, false)),
		d.Set("enterprise_project_id", utils.PathSearch("instance.enterprise_project_id", getInstanceRespBody, nil)),
		d.Set("asn", utils.PathSearch("instance.asn", getInstanceRespBody, nil)),
		d.Set("enable_default_propagation", utils.PathSearch("instance.enable_default_propagation", getInstanceRespBody, nil)),
		d.Set("enable_default_association", utils.PathSearch("instance.enable_default_association", getInstanceRespBody, nil)),
		d.Set("default_propagation_route_table_id", utils.PathSearch("instance.default_propagation_route_table_id", getInstanceRespBody, nil)),
		d.Set("default_association_route_table_id", utils.PathSearch("instance.default_association_route_table_id", getInstanceRespBody, nil)),
		d.Set("availability_zones", utils.PathSearch("instance.availability_zone_ids", getInstanceRespBody, nil)),
		d.Set("auto_accept_shared_attachments", utils.PathSearch("instance.auto_accept_shared_attachments", getInstanceRespBody, nil)),
		d.Set("tags", utils.FlattenTagsToMap(utils.PathSearch("instance.tags", getInstanceRespBody, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Id()
	)
	client, err := cfg.NewServiceClient("er", region)
	if err != nil {
		return diag.Errorf("error creating ER Client: %s", err)
	}

	updateInstancehasChanges := []string{
		"name",
		"description",
		"enable_default_propagation",
		"enable_default_association",
		"default_propagation_route_table_id",
		"default_association_route_table_id",
		"auto_accept_shared_attachments",
	}

	if d.HasChanges(updateInstancehasChanges...) {
		// Update the basic configuration of instance
		httpUrl := "v3/{project_id}/enterprise-router/instances/{er_id}"
		updatePath := client.Endpoint + httpUrl
		updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
		updatePath = strings.ReplaceAll(updatePath, "{er_id}", instanceId)

		updateInstanceOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}
		updateInstanceOpt.JSONBody = utils.RemoveNil(buildUpdateInstanceBodyParams(d))
		_, err = client.Request("PUT", updatePath, &updateInstanceOpt)
		if err != nil {
			return diag.Errorf("error updating Instance: %s", err)
		}
		err = instanceWaitingForStateCompleted(ctx, client, d, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.Errorf("error waiting for the update operation of the instance (%s) to complete: %s", instanceId, err)
		}
	}

	updateInstanceAvailabilityZoneshasChanges := []string{
		"availability_zones",
	}

	if d.HasChanges(updateInstanceAvailabilityZoneshasChanges...) {
		// updateInstanceAvailabilityZones: Update the availability zone list where the Enterprise router instance is located
		httpUrl := "v3/{project_id}/enterprise-router/instances/{er_id}/change-availability-zone-ids"
		updatePath := client.Endpoint + httpUrl
		updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
		updatePath = strings.ReplaceAll(updatePath, "{er_id}", instanceId)

		updateInstanceAvailabilityZonesOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				202,
			},
		}
		updateInstanceAvailabilityZonesOpt.JSONBody = utils.RemoveNil(buildUpdateInstanceAvailabilityZonesBodyParams(d))
		_, err = client.Request("POST", updatePath, &updateInstanceAvailabilityZonesOpt)
		if err != nil {
			return diag.Errorf("error updating instance: %s", err)
		}
		err = instanceWaitingForStateCompleted(ctx, client, d, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.Errorf("error waiting for the update operation of the instance (%s) availability zones to complete: %s", instanceId, err)
		}
	}

	if d.HasChange("tags") {
		err = utils.UpdateResourceTags(client, d, "instance", instanceId)
		if err != nil {
			return diag.Errorf("error updating instance tags: %s", err)
		}
	}

	if d.HasChange("enterprise_project_id") {
		migrateOpts := config.MigrateResourceOpts{
			ResourceId:   instanceId,
			ResourceType: "instance",
			RegionId:     region,
			ProjectId:    client.ProjectID,
		}
		if err := cfg.MigrateEnterpriseProject(ctx, d, migrateOpts); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceInstanceRead(ctx, d, meta)
}

func buildUpdateInstanceBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"instance": buildUpdateInstanceInstanceChildBody(d),
	}
	return bodyParams
}

func buildUpdateInstanceInstanceChildBody(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"name":                           utils.ValueIgnoreEmpty(d.Get("name")),
		"description":                    d.Get("description"),
		"auto_accept_shared_attachments": utils.ValueIgnoreEmpty(d.Get("auto_accept_shared_attachments")),
		"enable_default_propagation":     d.Get("auto_accept_shared_attachments"),
		"enable_default_association":     d.Get("auto_accept_shared_attachments"),
	}
	if d.Get("enable_default_propagation").(bool) {
		params["default_propagation_route_table_id"] = utils.ValueIgnoreEmpty(d.Get("default_propagation_route_table_id"))
	}
	if d.Get("enable_default_association").(bool) {
		params["default_association_route_table_id"] = utils.ValueIgnoreEmpty(d.Get("default_association_route_table_id"))
	}
	return params
}

func buildUpdateInstanceAvailabilityZonesBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"availability_zone_ids": utils.ValueIgnoreEmpty(d.Get("availability_zones")),
	}
	return bodyParams
}

func resourceInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/enterprise-router/instances/{er_id}"
	)
	client, err := cfg.NewServiceClient("er", region)
	if err != nil {
		return diag.Errorf("error creating ER Client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{er_id}", d.Id())

	deleteInstanceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			202,
		},
	}
	_, err = client.Request("DELETE", deletePath, &deleteInstanceOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting instance")
	}

	err = deleteInstanceWaitingForStateCompleted(ctx, client, d, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("error waiting for the delete operation of instance (%s) to complete: %s", d.Id(), err)
	}
	return nil
}

func deleteInstanceWaitingForStateCompleted(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			httpUrl := "v3/{project_id}/enterprise-router/instances/{er_id}"
			refreshPath := client.Endpoint + httpUrl
			refreshPath = strings.ReplaceAll(refreshPath, "{project_id}", client.ProjectID)
			refreshPath = strings.ReplaceAll(refreshPath, "{er_id}", d.Id())

			deleteInstanceWaitingOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
				OkCodes: []int{
					200,
				},
			}
			deleteInstanceWaitingResp, err := client.Request("GET", refreshPath, &deleteInstanceWaitingOpt)
			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					// When the error code is 404, the value of respBody is nil, and a non-null value is returned to avoid continuing the loop check.
					return "Resource Not Found", "COMPLETED", nil
				}
				return nil, "ERROR", err
			}

			deleteInstanceWaitingRespBody, err := utils.FlattenResponse(deleteInstanceWaitingResp)
			if err != nil {
				return nil, "ERROR", err
			}
			status := utils.PathSearch(`instance.state`, deleteInstanceWaitingRespBody, "").(string)

			if utils.StrSliceContains([]string{"fail"}, status) {
				return deleteInstanceWaitingRespBody, "", fmt.Errorf("unexpected status '%s'", status)
			}

			return deleteInstanceWaitingRespBody, "PENDING", nil
		},
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}
