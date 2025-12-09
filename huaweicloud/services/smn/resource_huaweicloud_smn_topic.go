package smn

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/smn/v2/topics"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API SMN GET /v2/{project_id}/notifications/topics/{id}/attributes
// @API SMN DELETE /v2/{project_id}/notifications/topics/{id}
// @API SMN GET /v2/{project_id}/notifications/topics/{id}
// @API SMN PUT /v2/{project_id}/notifications/topics/{id}
// @API SMN POST /v2/{project_id}/notifications/topics
// @API SMN POST /v2/{project_id}/{resource_type}/{id}/tags/action
// @API SMN GET /v2/{project_id}/{resource_type}/{id}/tags
// @API SMN PUT /v2/{project_id}/notifications/topics/{id}/attributes/{policyName}
func ResourceTopic() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTopicCreate,
		ReadContext:   resourceTopicRead,
		UpdateContext: resourceTopicUpdate,
		DeleteContext: resourceTopicDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"users_publish_allowed": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"services_publish_allowed": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"introduction": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": common.TagsSchema(),
			"access_policy": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "schema: Internal",
				ConflictsWith: []string{
					"users_publish_allowed", "services_publish_allowed",
				},
			},

			"topic_urn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"push_policy": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceTopicCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.SmnV2Client(region)
	if err != nil {
		return diag.Errorf("error creating SMN client: %s", err)
	}

	createOpts := topics.CreateOps{
		Name:                d.Get("name").(string),
		DisplayName:         d.Get("display_name").(string),
		EnterpriseProjectId: cfg.GetEnterpriseProjectID(d),
	}

	log.Printf("[DEBUG] create Options: %#v", createOpts)
	topic, err := topics.Create(client, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error getting SMN topic from result: %s", err)
	}

	log.Printf("[DEBUG] successfully created SMN topic: %s", topic.TopicUrn)
	d.SetId(topic.TopicUrn)

	// set policies
	err = updatePolicies(client, d, d.Id())
	if err != nil {
		diag.Errorf("error updating the policies of topic: %s", err)
	}

	// set tags
	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		taglist := utils.ExpandResourceTags(tagRaw)
		tagClient, err := cfg.SmnV2TagClient(region)
		if err != nil {
			return diag.Errorf("error creating SMN tag client: %s", err)
		}
		tagClient.MoreHeaders = map[string]string{
			"X-SMN-RESOURCEID-TYPE": "name",
		}
		if tagErr := tags.Create(tagClient, "smn_topic", d.Get("name").(string), taglist).ExtractErr(); tagErr != nil {
			return diag.Errorf("error setting tags of SMN topic %s: %s", d.Id(), tagErr)
		}
	}

	return resourceTopicRead(ctx, d, meta)
}

func resourceTopicRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.SmnV2Client(region)
	if err != nil {
		return diag.Errorf("error creating SMN client: %s", err)
	}

	topicUrn := d.Id()
	topicGet, err := topics.Get(client, topicUrn).ExtractGet()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving SMN topic")
	}

	log.Printf("[DEBUG] retrieved SMN topic %s: %#v", topicUrn, topicGet)
	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("topic_urn", topicGet.TopicUrn),
		d.Set("display_name", topicGet.DisplayName),
		d.Set("name", topicGet.Name),
		d.Set("push_policy", topicGet.PushPolicy),
		d.Set("update_time", topicGet.UpdateTime),
		d.Set("create_time", topicGet.CreateTime),
		d.Set("enterprise_project_id", topicGet.EnterpriseProjectId),
	)

	// fetch access policies
	policy, err := topics.GetPolicies(client, topicUrn, "access_policy").Extract()
	if err != nil {
		return diag.Errorf("error fetching the access_policy of the topic: %s", err)
	}

	var (
		usersPublishAllowed    string
		servicesPublishAllowed string
	)

	if policy.AccessPolicy != "" {
		var accessPolicy map[string]interface{}
		err = json.Unmarshal([]byte(policy.AccessPolicy), &accessPolicy)
		if err != nil {
			return diag.FromErr(err)
		}

		csp := utils.PathSearch("Statement[?Sid== '__user_pub_0'].Principal.CSP|[0]", accessPolicy, "")
		if csp == "" || csp == "*" {
			usersPublishAllowed = csp.(string)
		} else {
			usersPublishAllowed = strings.Join(utils.ExpandToStringList(csp.([]interface{})), ",")
		}

		services := utils.PathSearch("Statement[?Sid== '__service_pub_0'].Principal.Service|[0]", accessPolicy, []interface{}{})
		servicesPublishAllowed = strings.Join(utils.ExpandToStringList(services.([]interface{})), ",")
	}

	// users_publish_allowed and services_publish_allowed will not be set if access_policy is specified
	if _, ok := d.GetOk("access_policy"); !ok {
		mErr = multierror.Append(mErr,
			d.Set("users_publish_allowed", usersPublishAllowed),
			d.Set("services_publish_allowed", servicesPublishAllowed),
		)
	}

	// fetch introduction
	introduction, err := topics.GetPolicies(client, topicUrn, "introduction").Extract()
	if err != nil {
		return diag.Errorf("error fetching the introduction of the topic: %s", err)
	}

	mErr = multierror.Append(mErr, d.Set("introduction", introduction.Introduction))

	// fetch tags
	tagClient, err := cfg.SmnV2TagClient(region)
	if err != nil {
		return diag.Errorf("error creating SMN tag client: %s", err)
	}
	tagClient.MoreHeaders = map[string]string{
		"X-SMN-RESOURCEID-TYPE": "name",
	}
	if resourceTags, err := tags.Get(tagClient, "smn_topic", d.Get("name").(string)).Extract(); err == nil {
		tagmap := utils.TagsToMap(resourceTags.Tags)
		mErr = multierror.Append(mErr, d.Set("tags", tagmap))
	} else {
		log.Printf("[WARN] fetching tags of SMN topic failed: %s", err)
	}

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting SMN topic fields: %s", err)
	}

	return nil
}

func resourceTopicUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.SmnV2Client(region)
	if err != nil {
		return diag.Errorf("error creating SMN client: %s", err)
	}

	id := d.Id()
	if d.HasChange("display_name") {
		updateOpts := topics.UpdateOps{
			DisplayName: d.Get("display_name").(string),
		}
		if _, err = topics.Update(client, updateOpts, id).Extract(); err != nil {
			return diag.Errorf("error updating SMN topic %s: %s", id, err)
		}
	}

	if d.HasChanges("access_policy", "users_publish_allowed", "services_publish_allowed", "introduction") {
		err := updatePolicies(client, d, id)
		if err != nil {
			diag.Errorf("error updating the policies of topic: %s", err)
		}
	}

	// update tags
	if d.HasChange("tags") {
		tagClient, err := cfg.SmnV2TagClient(region)
		if err != nil {
			return diag.Errorf("error creating SMN tag client: %s", err)
		}
		tagClient.MoreHeaders = map[string]string{
			"X-SMN-RESOURCEID-TYPE": "name",
		}
		tagErr := utils.UpdateResourceTags(tagClient, d, "smn_topic", d.Get("name").(string))
		if tagErr != nil {
			return diag.Errorf("error updating tags of SMN topic %s: %s", id, tagErr)
		}
	}

	return resourceTopicRead(ctx, d, meta)
}

func updatePolicies(client *golangsdk.ServiceClient, d *schema.ResourceData, id string) error {
	if d.HasChanges("users_publish_allowed", "services_publish_allowed", "access_policy") {
		value, err := buildUpdateAccessPolicy(d)
		if err != nil {
			return err
		}

		opts := topics.UpdatePoliciesOpts{
			Value: value,
		}
		_, err = topics.UpdatePolicies(client, opts, id, "access_policy").Extract()
		if err != nil {
			return err
		}
	}

	if d.HasChange("introduction") {
		opts := topics.UpdatePoliciesOpts{
			Value: d.Get("introduction").(string),
		}
		_, err := topics.UpdatePolicies(client, opts, id, "introduction").Extract()
		if err != nil {
			return err
		}
	}

	return nil
}

func buildUpdateAccessPolicy(d *schema.ResourceData) (string, error) {
	if v, ok := d.GetOk("access_policy"); ok {
		return v.(string), nil
	}

	statement := []map[string]interface{}{}
	if v, ok := d.GetOk("users_publish_allowed"); ok {
		var csp interface{}
		if v == "*" {
			csp = "*"
		} else {
			csp = strings.Split(v.(string), ",")
		}

		statement = append(statement, map[string]interface{}{
			"Sid":    "__user_pub_0",
			"Effect": "Allow",
			"Principal": map[string]interface{}{
				"CSP": csp,
			},
			"Action": []interface{}{
				"SMN:Publish",
				"SMN:QueryTopicDetail",
			},
			"Resource": d.Id(),
		})
	}

	if v, ok := d.GetOk("services_publish_allowed"); ok {
		service := strings.Split(v.(string), ",")

		statement = append(statement, map[string]interface{}{
			"Sid":    "__service_pub_0",
			"Effect": "Allow",
			"Principal": map[string]interface{}{
				"Service": service,
			},
			"Action": []interface{}{
				"SMN:Publish",
				"SMN:QueryTopicDetail",
			},
			"Resource": d.Id(),
		})
	}

	if len(statement) == 0 {
		return "", nil
	}

	params := map[string]interface{}{
		"Version":   "2016-09-07",
		"Id":        "__default_policy_ID",
		"Statement": statement,
	}

	jsonData, err := json.Marshal(params)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

func resourceTopicDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.SmnV2Client(region)
	if err != nil {
		return diag.Errorf("error creating SMN client: %s", err)
	}

	err = topics.Delete(client, d.Id()).ExtractErr()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting SMN topic")
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForTopicDelete(client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error deleting SMN topic %s: %s", d.Id(), err)
	}

	return nil
}

func waitForTopicDelete(client *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r, err := topics.Get(client, id).ExtractGet()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] successfully deleted topic %s", id)
				return r, "DELETED", nil
			}
			return r, "ACTIVE", err
		}
		return r, "ACTIVE", nil
	}
}
