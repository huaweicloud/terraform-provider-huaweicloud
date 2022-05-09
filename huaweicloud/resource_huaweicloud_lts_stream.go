package huaweicloud

import (
	"github.com/chnsz/golangsdk/openstack/lts/huawei/logstreams"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceLTSStreamV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceStreamV2Create,
		Read:   resourceStreamV2Read,
		Delete: resourceStreamV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"stream_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"filter_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceStreamV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.LtsV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud LTS client: %s", err)
	}

	groupId := d.Get("group_id").(string)
	createOpts := &logstreams.CreateOpts{
		LogStreamName: d.Get("stream_name").(string),
	}

	logp.Printf("[DEBUG] Create Options: %#v", createOpts)

	streamCreate, err := logstreams.Create(client, groupId, createOpts).Extract()
	if err != nil {
		return fmtp.Errorf("Error creating log stream: %s", err)
	}

	d.SetId(streamCreate.ID)
	return resourceStreamV2Read(d, meta)
}

func resourceStreamV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.LtsV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud LTS client: %s", err)
	}

	groupId := d.Get("group_id").(string)
	streams, err := logstreams.List(client, groupId).Extract()
	if err != nil {
		return fmtp.Errorf("Error getting HuaweiCloud log stream %s: %s", d.Id(), err)
	}

	for _, stream := range streams.LogStreams {
		if stream.ID == d.Id() {
			logp.Printf("[DEBUG] Retrieved log stream %s: %#v", d.Id(), stream)
			d.SetId(stream.ID)
			d.Set("stream_name", stream.Name)
			d.Set("filter_count", stream.FilterCount)
			return nil
		}
	}

	return fmtp.Errorf("Error HuaweiCloud log group stream %s: No Found", d.Id())
}

func resourceStreamV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.LtsV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud LTS client: %s", err)
	}

	groupId := d.Get("group_id").(string)
	err = logstreams.Delete(client, groupId, d.Id()).ExtractErr()
	if err != nil {
		return CheckDeleted(d, err, "Error deleting log stream")
	}

	d.SetId("")
	return nil
}
