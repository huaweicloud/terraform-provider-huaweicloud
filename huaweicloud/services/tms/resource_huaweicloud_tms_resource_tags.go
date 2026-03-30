package tms

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	tags "github.com/chnsz/golangsdk/openstack/tms/v1/resourcetags"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API TMS POST /v1.0/resource-tags/batch-create
// @API TMS POST /v1.0/resource-tags/batch-delete
// @API TMS GET /v2.0/resources/{resource_id}/tags
func ResourceResourceTags() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceResourceTagsCreate,
		ReadContext:   resourceResourceTagsRead,
		UpdateContext: resourceResourceTagsUpdate,
		DeleteContext: resourceResourceTagsDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The project ID of the resources.",
			},
			"resources": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The resource type.",
						},
						"resource_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The resource ID.",
						},
					},
				},
				Description: "The managed resource configuration.",
			},
			"tags": {
				Type:     schema.TypeMap,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The resource tags for batch management.",
			},
		},
	}
}

func buildResourcesInfo(resources []interface{}) []tags.Resource {
	if len(resources) < 1 {
		return nil
	}

	result := make([]tags.Resource, len(resources))
	for i, val := range resources {
		resource := val.(map[string]interface{})
		result[i] = tags.Resource{
			ResourceType: resource["resource_type"].(string),
			ResourceId:   resource["resource_id"].(string),
		}
	}
	return result
}

func expandResourceTags(tagsInput map[string]interface{}) []tags.ResourceTag {
	result := make([]tags.ResourceTag, 0, len(tagsInput))

	for key, value := range tagsInput {
		result = append(result, tags.ResourceTag{
			Key:   key,
			Value: utils.String(value.(string)),
		})
	}
	return result
}

func waitingResourceTagsConfigCompleted(ctx context.Context, cfg *config.Config, d *schema.ResourceData,
	timeout time.Duration) error {
	client, err := cfg.TmsV2Client(cfg.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating TMS v2 client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			// For some services, binding tags is an asynchronous operation, such as MRS.
			resResult, err := GetResourceTags(client, d)
			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					return "waiting config", "PENDING", nil
				}

				return nil, "ERROR", err
			}

			return resResult, "COMPLETED", nil
		},
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	return err
}

func resourceResourceTagsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.TmsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating TMS v1 client: %s", err)
	}

	opts := tags.BatchOpts{
		ProjectId: d.Get("project_id").(string),
		Resources: buildResourcesInfo(d.Get("resources").([]interface{})),
		Tags:      expandResourceTags(d.Get("tags").(map[string]interface{})),
	}
	failResp, err := tags.Create(client, opts)
	if err != nil {
		return diag.Errorf("error creating resource tags: %s", err)
	}
	if len(failResp) > 0 {
		return diag.Errorf("error creating resource tags: %#v", failResp)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate resource ID of the TMS tags management: %s", err)
	}
	d.SetId(randUUID)

	if err := waitingResourceTagsConfigCompleted(ctx, cfg, d, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.Errorf("error waiting for TMS resource tags config to complete: %s", err)
	}

	return resourceResourceTagsRead(ctx, d, meta)
}

func FlattenTagsToMap(tagsResp []tags.ResourceTag) map[string]interface{} {
	result := make(map[string]interface{})
	for _, val := range tagsResp {
		result[val.Key] = *val.Value
	}
	return result
}

func compareTwoTags(localTags, remoteTags map[string]interface{}) (same, diff map[string]interface{}) {
	same = make(map[string]interface{})
	diff = make(map[string]interface{})

	for localKey, localVal := range localTags {
		if remoteVal, ok := remoteTags[localKey]; ok {
			local, isTypeLocalOk := localVal.(string)
			if !isTypeLocalOk {
				log.Printf("[WARN] The type of tag key (%s) in the script is incorrect, want 'string', but got '%T'",
					localKey, localVal)
				continue
			}
			remote, isTypeRemoteOk := remoteVal.(string)
			if !isTypeRemoteOk {
				log.Printf("[WARN] The type of tag key (%s) in the remote response is incorrect, want 'string', but got '%T'",
					localKey, remoteVal)
				continue
			}
			if local == remote {
				same[localKey] = localVal
				continue
			}
		}
		diff[localKey] = localVal
	}
	return
}

func GetResourceTags(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		projectId = d.Get("project_id").(string)
		resources = d.Get("resources").([]interface{})
		resResult = make([]interface{}, 0, len(resources))
		tagsInput = d.Get("tags").(map[string]interface{})
	)

	// Check whether all tagged resources contain the expected tags correctly. If not, inconsistent tags information
	// will be printed in the log.
	for _, val := range resources {
		resource := val.(map[string]interface{})
		resourceId := resource["resource_id"].(string)
		opts := tags.QueryOpts{
			ResourceId:   resourceId,
			ResourceType: resource["resource_type"].(string),
			ProjectId:    projectId,
		}
		resp, err := tags.Get(client, opts)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				continue
			}
			return nil, fmt.Errorf("error query resource (%s) tags: %s", resourceId, err)
		}
		actualTags := FlattenTagsToMap(resp)
		same, diff := compareTwoTags(tagsInput, actualTags)
		if len(diff) > 0 {
			log.Printf("[ERROR] The tags of resource (%s) don't contain some tags that are expected to need to be set."+
				" It should contain tags (%#v), but some tags were not set successfully: %#v", resourceId, tagsInput,
				actualTags)
		}
		// If the tags are queried from the resource side (even only part of tags are queried), it means that the
		// creation or update action of the resource is (partially) successful.
		if len(same) > 0 {
			resResult = append(resResult, val)
		}
	}

	// All tags set failed.
	if len(resResult) < 1 {
		return nil, golangsdk.ErrDefault404{}
	}

	return resResult, nil
}

func resourceResourceTagsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.TmsV2Client(region)
	if err != nil {
		return diag.Errorf("error creating TMS v2 client: %s", err)
	}

	resResult, err := GetResourceTags(client, d)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error getting TMS resource tags")
	}

	mErr := multierror.Append(nil,
		d.Set("resources", resResult),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving resources and tags information: %s", err)
	}
	return nil
}

func resourceResourceTagsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.TmsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating TMS v1 client: %s", err)
	}

	var (
		projectId        = d.Get("project_id").(string)
		oldRes, newRes   = d.GetChange("resources")
		oldTags, newTags = d.GetChange("tags")
	)

	deleteOpts := tags.BatchOpts{
		ProjectId: projectId,
		Resources: buildResourcesInfo(oldRes.([]interface{})),
		Tags:      expandResourceTags(oldTags.(map[string]interface{})),
	}
	failResp, err := tags.Delete(client, deleteOpts)
	if err != nil {
		return diag.Errorf("error deleting resource tags: %s", err)
	}
	if len(failResp) > 0 {
		return diag.Errorf("some tags were not successfully removed: %#v", failResp)
	}

	opts := tags.BatchOpts{
		ProjectId: projectId,
		Resources: buildResourcesInfo(newRes.([]interface{})),
		Tags:      expandResourceTags(newTags.(map[string]interface{})),
	}
	failResp, err = tags.Create(client, opts)
	if err != nil {
		return diag.Errorf("error creating resource tags: %s", err)
	}
	if len(failResp) > 0 {
		return diag.Errorf("some tags were not set successfully: %#v", failResp)
	}

	if err := waitingResourceTagsConfigCompleted(ctx, cfg, d, d.Timeout(schema.TimeoutUpdate)); err != nil {
		return diag.Errorf("error waiting for TMS resource tags config to complete: %s", err)
	}

	return resourceResourceTagsRead(ctx, d, meta)
}

func waitingResourceTagsDeleteCompleted(ctx context.Context, cfg *config.Config, d *schema.ResourceData,
	timeout time.Duration) error {
	client, err := cfg.TmsV2Client(cfg.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating TMS v2 client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			// For some services, binding tags is an asynchronous operation, such as MRS.
			resResult, err := GetResourceTags(client, d)
			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					return "success deleted", "COMPLETED", nil
				}

				return nil, "ERROR", err
			}

			return resResult, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	return err
}

func resourceResourceTagsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.TmsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating TMS v1 client: %s", err)
	}

	deleteOpts := tags.BatchOpts{
		ProjectId: d.Get("project_id").(string),
		Resources: buildResourcesInfo(d.Get("resources").([]interface{})),
		Tags:      expandResourceTags(d.Get("tags").(map[string]interface{})),
	}
	failResp, err := tags.Delete(client, deleteOpts)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "TMS tags management")
	}
	if len(failResp) > 0 {
		return diag.Errorf("some tags were not successfully removed: %#v", failResp)
	}

	if err := waitingResourceTagsDeleteCompleted(ctx, cfg, d, d.Timeout(schema.TimeoutDelete)); err != nil {
		return diag.Errorf("error waiting for TMS resource tags delete to complete: %s", err)
	}

	return nil
}
