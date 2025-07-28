package kafka

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/dms/v2/kafka/instances"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Kafka GET /v1/{project_id}/instances/{instance_id}/topics/{topic_name}/accesspolicy
// @API Kafka POST /v1/{project_id}/instances/{instance_id}/topics/accesspolicy
// @API Kafka GET /v2/{project_id}/instances/{instance_id}
func ResourceDmsKafkaPermissions() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDmsKafkaPermissionsCreateOrUpdate,
		UpdateContext: resourceDmsKafkaPermissionsCreateOrUpdate,
		DeleteContext: resourceDmsKafkaPermissionsDelete,
		ReadContext:   resourceDmsKafkaPermissionsRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"topic_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"policies": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"access_policy": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"pub", "sub", "all",
							}, false),
						},
					},
				},
			},
		},
	}
}

func buildPoliciesOpts(rawPolicies []interface{}) []interface{} {
	if len(rawPolicies) < 1 {
		return nil
	}

	policies := make([]interface{}, 0, len(rawPolicies))
	for _, v := range rawPolicies {
		policy := v.(map[string]interface{})
		policies = append(policies, map[string]interface{}{
			"user_name":     policy["user_name"],
			"access_policy": policy["access_policy"],
		})
	}

	return policies
}

func buildCreateOrUpdateDmsKafkaPermissionsBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":     d.Get("topic_name"),
		"policies": buildPoliciesOpts(d.Get("policies").([]interface{})),
	}
	return map[string]interface{}{
		"topics": []interface{}{bodyParams},
	}
}

func resourceDmsKafkaPermissionsCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.DmsV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	topicName := d.Get("topic_name").(string)
	instanceId := d.Get("instance_id").(string)

	if err := updateKafkaPermissions(client, instanceId, buildCreateOrUpdateDmsKafkaPermissionsBodyParams(d)); err != nil {
		return diag.Errorf("error setting DMS kafka permissions: %s", err)
	}

	id := instanceId + "/" + topicName
	d.SetId(id)

	if err = waitForKafkaTopicAccessPolicyComplete(ctx, client, d, instanceId, schema.TimeoutCreate); err != nil {
		return diag.FromErr(err)
	}

	return resourceDmsKafkaPermissionsRead(ctx, d, meta)
}

func updateKafkaPermissions(client *golangsdk.ServiceClient, instanceId string, opt map[string]interface{}) error {
	createHttpUrl := "v1/{project_id}/instances/{instance_id}/topics/accesspolicy"
	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200, 204,
		},
		JSONBody: utils.RemoveNil(opt),
	}

	_, err := client.Request("POST", createPath, &createOpt)
	return err
}

func resourceDmsKafkaPermissionsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.DmsV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	// Split instance_id and topic_name from resource id
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return diag.Errorf("invalid ID format, must be <instance_id>/<topic_name>")
	}
	instanceId := parts[0]
	topicName := parts[1]

	policies, err := GetDmsKafkaPermissions(client, instanceId, topicName)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DMS kafka permission")
	}
	mErr := multierror.Append(nil,
		d.Set("region", cfg.GetRegion(d)),
		d.Set("instance_id", instanceId),
		d.Set("topic_name", topicName),
		d.Set("policies", flattenPolicies(policies)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func GetDmsKafkaPermissions(client *golangsdk.ServiceClient, instId, name string) ([]interface{}, error) {
	getHttpUrl := "v1/{project_id}/instances/{instance_id}/topics/{topic_name}/accesspolicy"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instId)
	getPath = strings.ReplaceAll(getPath, "{topic_name}", name)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	policies := utils.PathSearch("policies", getRespBody, make([]interface{}, 0)).([]interface{})
	if len(policies) == 0 {
		return nil, golangsdk.ErrDefault404{}
	}

	return policies, nil
}

func resourceDmsKafkaPermissionsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.DmsV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)

	if err := updateKafkaPermissions(client, instanceId, buildDeleteDmsKafkaPermissionsBodyParams(d)); err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(
			err, "failed_topics|[0].error_msg", "Topic policy is empty."), "error deleting DMS kafka permissions")
	}

	if err := waitForKafkaTopicAccessPolicyComplete(ctx, client, d, instanceId, schema.TimeoutDelete); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func buildDeleteDmsKafkaPermissionsBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":     d.Get("topic_name"),
		"policies": []interface{}{},
	}
	return map[string]interface{}{
		"topics": []interface{}{bodyParams},
	}
}

func flattenPolicies(policies []interface{}) []map[string]interface{} {
	policiesToSet := make([]map[string]interface{}, len(policies))
	for i, v := range policies {
		policiesToSet[i] = map[string]interface{}{
			"user_name":     utils.PathSearch("user_name", v, nil),
			"access_policy": utils.PathSearch("access_policy", v, nil),
		}
	}

	return policiesToSet
}

func waitForKafkaTopicAccessPolicyComplete(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	instanceID string, timeout string) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"CREATED"},
		Refresh:      kafkaInstancePolicyRefreshFunc(client, instanceID),
		Timeout:      d.Timeout(timeout),
		Delay:        1 * time.Second,
		PollInterval: 2 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for DMS Kafka instance (%s) task to be completed: %s", d.Id(), err)
	}
	return nil
}

func kafkaInstancePolicyRefreshFunc(client *golangsdk.ServiceClient, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := instances.Get(client, instanceID).Extract()
		if err != nil {
			return nil, "ERROR", err
		}
		if resp.Task.Name == "updateTopicPolicies" {
			return resp, "PENDING", nil
		}
		return resp, "CREATED", nil
	}
}
