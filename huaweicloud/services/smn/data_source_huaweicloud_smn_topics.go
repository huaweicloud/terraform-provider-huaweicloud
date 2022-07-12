package smn

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/smn/v2/topics"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func DataSourceTopics() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTopicsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"topic_urn": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"topics": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"topic_urn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"display_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enterprise_project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"push_policy": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataSourceTopicsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	client, err := config.SmnV2Client(region)
	if err != nil {
		return diag.Errorf("error creating SMN client: %s", err)
	}
	// create tagClient to fetch tags
	tagClient, err := config.SmnV2TagClient(region)
	if err != nil {
		return diag.Errorf("error creating SMN tag client: %s", err)
	}
	tagClient.MoreHeaders = map[string]string{
		"X-SMN-RESOURCEID-TYPE": "name",
	}

	allTopics, err := topics.List(client).Extract()
	if err != nil {
		return diag.Errorf("unable to list topics: %s ", err)
	}

	log.Printf("[DEBUG] Retrieved SMN topics: %#v", allTopics)

	filter := map[string]interface{}{
		"TopicUrn":            d.Get("topic_urn"),
		"Name":                d.Get("name"),
		"EnterpriseProjectId": d.Get("enterprise_project_id"),
		"DisplayName":         d.Get("display_name"),
	}

	filterTopics, err := utils.FilterSliceWithField(allTopics, filter)
	if err != nil {
		return diag.Errorf("filter topics failed: %s", err)
	}

	ids := make([]string, len(filterTopics))
	stateTopics := make([]map[string]interface{}, len(filterTopics))

	for i, item := range filterTopics {
		topic := item.(topics.TopicGet)
		ids[i] = topic.TopicUrn
		stateTopics[i] = flattenSourceTopic(tagClient, topic)
	}

	if len(ids) == 1 {
		d.SetId(ids[0])
	} else {
		d.SetId(hashcode.Strings(ids))
	}

	if err := d.Set("topics", stateTopics); err != nil {
		diag.Errorf("error setting SMN topics: %s", err)
	}

	return nil
}

func flattenSourceTopic(tagClient *golangsdk.ServiceClient, topic topics.TopicGet) map[string]interface{} {
	stateTopic := map[string]interface{}{
		"topic_urn":             topic.TopicUrn,
		"id":                    topic.TopicUrn,
		"display_name":          topic.DisplayName,
		"name":                  topic.Name,
		"push_policy":           topic.PushPolicy,
		"enterprise_project_id": topic.EnterpriseProjectId,
	}

	if resourceTags, err := tags.Get(tagClient, "smn_topic", topic.Name).Extract(); err == nil {
		tagmap := utils.TagsToMap(resourceTags.Tags)
		stateTopic["tags"] = tagmap
	} else {
		log.Printf("[WARN] fetching tags of SMN topic failed: %s", err)
	}

	return stateTopic
}
