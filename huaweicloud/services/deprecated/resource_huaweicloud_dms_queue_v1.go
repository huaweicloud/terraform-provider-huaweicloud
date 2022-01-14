package deprecated

import (
	"context"

	"github.com/chnsz/golangsdk/openstack/dms/v1/queues"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceDmsQueues() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDmsQueuesCreate,
		ReadContext:   resourceDmsQueuesRead,
		DeleteContext: resourceDmsQueuesDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		DeprecationMessage: "Deprecated, Distributed Message Service (Shared Edition) has withdrawn, " +
			"please use DMS for Kafka instead.",

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
			"queue_mode": {
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
			"redrive_policy": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"max_consume_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"retention_hours": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"reservation": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"max_msg_size_byte": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"produced_messages": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"group_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceDmsQueuesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	dmsV1Client, err := config.DmsV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud dms queue client: %s", err)
	}

	createOpts := &queues.CreateOps{
		Name:            d.Get("name").(string),
		QueueMode:       d.Get("queue_mode").(string),
		Description:     d.Get("description").(string),
		RedrivePolicy:   d.Get("redrive_policy").(string),
		MaxConsumeCount: d.Get("max_consume_count").(int),
		RetentionHours:  d.Get("retention_hours").(int),
	}

	logp.Printf("[DEBUG] Create Options: %#v", createOpts)
	v, err := queues.Create(dmsV1Client, createOpts).Extract()
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud queue: %s", err)
	}
	logp.Printf("[INFO] Queue ID: %s", v.ID)

	// Store the queue ID now
	d.SetId(v.ID)

	return resourceDmsQueuesRead(ctx, d, meta)
}

func resourceDmsQueuesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)

	dmsV1Client, err := config.DmsV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud dms queue client: %s", err)
	}
	v, err := queues.Get(dmsV1Client, d.Id(), true).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	logp.Printf("[DEBUG] Dms queue %s: %+v", d.Id(), v)

	d.SetId(v.ID)
	d.Set("name", v.Name)
	d.Set("created", v.Created)
	d.Set("description", v.Description)
	d.Set("queue_mode", v.QueueMode)
	d.Set("reservation", v.Reservation)
	d.Set("max_msg_size_byte", v.MaxMsgSizeByte)
	d.Set("produced_messages", v.ProducedMessages)
	d.Set("redrive_policy", v.RedrivePolicy)
	d.Set("max_consume_count", v.MaxConsumeCount)
	d.Set("group_count", v.GroupCount)

	return nil
}

func resourceDmsQueuesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	dmsV1Client, err := config.DmsV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud dms queue client: %s", err)
	}

	v, err := queues.Get(dmsV1Client, d.Id(), false).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "queue")
	}

	err = queues.Delete(dmsV1Client, d.Id()).ExtractErr()
	if err != nil {
		return fmtp.DiagErrorf("Error deleting HuaweiCloud queue: %s", err)
	}

	logp.Printf("[DEBUG] Dms queue %s: %+v deactivated.", d.Id(), v)
	d.SetId("")
	return nil
}
