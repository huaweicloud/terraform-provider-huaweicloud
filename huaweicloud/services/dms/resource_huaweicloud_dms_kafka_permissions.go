package dms

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/kafka/v2/model"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func ResourceDmsKafkaPermissions() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDmsKafkaPermissionsCreateOrUpdate,
		UpdateContext: resourceDmsKafkaPermissionsCreateOrUpdate,
		DeleteContext: resourceDmsKafkaPermissionsDelete,
		ReadContext:   resourceDmsKafkaPermissionsRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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

func buildPoliciesOpts(rawPolicies []interface{}) ([]model.AccessPolicyEntity, error) {
	if len(rawPolicies) < 1 {
		return nil, nil
	}

	policies := make([]model.AccessPolicyEntity, len(rawPolicies))
	for i, v := range rawPolicies {
		policy := v.(map[string]interface{})
		var accessPolicy model.AccessPolicyEntityAccessPolicy
		if err := accessPolicy.UnmarshalJSON([]byte(policy["access_policy"].(string))); err != nil {
			return nil, fmt.Errorf("error parsing the argument access_policy: %s", err)
		}
		policies[i] = model.AccessPolicyEntity{
			UserName:     utils.String(policy["user_name"].(string)),
			AccessPolicy: &accessPolicy,
		}
	}

	return policies, nil
}

func resourceDmsKafkaPermissionsCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	client, err := c.HcDmsV2Client(c.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	topicName := d.Get("topic_name").(string)
	instanceId := d.Get("instance_id").(string)

	policies, err := buildPoliciesOpts(d.Get("policies").([]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}

	createOrUpdateOpts := &model.UpdateTopicAccessPolicyRequest{
		InstanceId: instanceId,
		Body: &model.UpdateTopicAccessPolicyReq{
			Topics: []model.AccessPolicyTopicEntity{
				{
					Name:     topicName,
					Policies: policies,
				},
			},
		},
	}

	_, err = client.UpdateTopicAccessPolicy(createOrUpdateOpts)
	if err != nil {
		return diag.Errorf("error creating DMS kafka permissions: %s", err)
	}

	id := instanceId + "/" + topicName
	d.SetId(id)
	return resourceDmsKafkaPermissionsRead(ctx, d, meta)
}

func resourceDmsKafkaPermissionsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	client, err := c.HcDmsV2Client(c.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	// Split instance_id and topic_name from resource id
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return diag.Errorf("invalid id format, must be <instance_id>/<topic_name>")
	}
	instanceId := parts[0]
	topicName := parts[1]

	request := &model.ShowTopicAccessPolicyRequest{
		InstanceId: instanceId,
		TopicName:  topicName,
	}

	response, err := client.ShowTopicAccessPolicy(request)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DMS kafka permission")
	}

	if response.Policies != nil && len(*response.Policies) != 0 {
		policies := *response.Policies
		d.Set("instance_id", instanceId)
		d.Set("topic_name", topicName)
		d.Set("policies", flattenPolicies(policies))
		return nil
	}

	// DB permission deleted
	log.Printf("[WARN] failed to fetch DMS kafka permission %s: deleted", d.Id())
	d.SetId("")

	return nil
}

func resourceDmsKafkaPermissionsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	client, err := c.HcDmsV2Client(c.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	topicName := d.Get("topic_name").(string)
	instanceId := d.Get("instance_id").(string)

	deleteOpts := &model.UpdateTopicAccessPolicyRequest{
		InstanceId: instanceId,
		Body: &model.UpdateTopicAccessPolicyReq{
			Topics: []model.AccessPolicyTopicEntity{
				{
					Name:     topicName,
					Policies: []model.AccessPolicyEntity{},
				},
			},
		},
	}

	_, err = client.UpdateTopicAccessPolicy(deleteOpts)
	if err != nil {
		return diag.Errorf("error deleting DMS kafka permissions: %s", err)
	}

	return nil
}

func flattenPolicies(policies []model.PolicyEntity) []map[string]interface{} {
	policiesToSet := make([]map[string]interface{}, len(policies))
	for i, v := range policies {
		policiesToSet[i] = map[string]interface{}{
			"user_name":     v.UserName,
			"access_policy": v.AccessPolicy.Value(),
		}
	}

	return policiesToSet
}
