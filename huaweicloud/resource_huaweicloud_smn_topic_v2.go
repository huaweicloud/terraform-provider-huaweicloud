package huaweicloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/smn/v2/topics"
)

func resourceTopic() *schema.Resource {
	return &schema.Resource{
		Create: resourceTopicCreate,
		Read:   resourceTopicRead,
		Delete: resourceTopicDelete,
		Update: resourceTopicUpdate,

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

func resourceTopicCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.SmnV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud smn client: %s", err)
	}

	createOpts := topics.CreateOps{
		Name:        d.Get("name").(string),
		DisplayName: d.Get("display_name").(string),
	}
	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	topic, err := topics.Create(client, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error getting topic from result: %s", err)
	}
	log.Printf("[DEBUG] Create : topic.TopicUrn %s", topic.TopicUrn)
	if topic.TopicUrn != "" {
		d.SetId(topic.TopicUrn)
		return resourceTopicRead(d, meta)
	}

	return fmt.Errorf("Unexpected conversion error in resourceTopicCreate.")
}

func resourceTopicRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.SmnV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud smn client: %s", err)
	}

	topicUrn := d.Id()
	topicGet, err := topics.Get(client, topicUrn).ExtractGet()
	if err != nil {
		return CheckDeleted(d, err, "topic")
	}

	log.Printf("[DEBUG] Retrieved topic %s: %#v", topicUrn, topicGet)

	d.Set("topic_urn", topicGet.TopicUrn)
	d.Set("display_name", topicGet.DisplayName)
	d.Set("name", topicGet.Name)
	d.Set("push_policy", topicGet.PushPolicy)
	d.Set("update_time", topicGet.UpdateTime)
	d.Set("create_time", topicGet.CreateTime)

	return nil
}

func resourceTopicDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.SmnV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud smn client: %s", err)
	}

	log.Printf("[DEBUG] Deleting topic %s", d.Id())

	id := d.Id()
	result := topics.Delete(client, id)
	if result.Err != nil {
		return result.Err
	}

	for {
		_, err = topics.Get(client, id).ExtractGet()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				break
			}
		}
	}

	log.Printf("[DEBUG] Successfully deleted topic %s", id)
	return nil
}

func resourceTopicUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.SmnV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud smn client: %s", err)
	}

	log.Printf("[DEBUG] Updating topic %s", d.Id())
	id := d.Id()

	var updateOpts topics.UpdateOps
	if d.HasChange("display_name") {
		updateOpts.DisplayName = d.Get("display_name").(string)
	}

	_, err = topics.Update(client, updateOpts, id).Extract()
	if err != nil {
		return fmt.Errorf("Error updating topic from result: %s", err)
	}

	log.Printf("[DEBUG] Update : topic.TopicUrn: %s", id)
	return resourceTopicRead(d, meta)
}
