package iotda

import (
	"context"
	"log"
	"regexp"

	"github.com/chnsz/golangsdk"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iotda/v5/model"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func ResourceAmqp() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAmqpCreate,
		ReadContext:   resourceAmqpRead,
		DeleteContext: resourceAmqpDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

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
				ValidateFunc: validation.All(
					validation.StringLenBetween(8, 128),
					validation.StringMatch(regexp.MustCompile(`^[A-Za-z-_0-9.:]*$`),
						"Only letters, digits, hyphens (-), underscores (_), dots (.) and colons (:) are allowed"),
				),
			},
		},
	}
}

func resourceAmqpCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	region := c.GetRegion(d)
	client, err := c.HcIoTdaV5Client(region)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	createOpts := model.AddQueueRequest{
		Body: &model.QueueInfo{
			QueueName: d.Get("name").(string),
		},
	}
	log.Printf("[DEBUG] Create IoTDA AMQP queue params: %#v", createOpts)

	resp, err := client.AddQueue(&createOpts)
	if err != nil {
		return diag.Errorf("error creating IoTDA AMQP queue: %s", err)
	}

	if resp.QueueId == nil {
		return diag.Errorf("error creating IoTDA AMQP queue: id is not found in API response")
	}

	d.SetId(*resp.QueueId)
	return resourceAmqpRead(ctx, d, meta)
}

func resourceAmqpRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	region := c.GetRegion(d)
	client, err := c.HcIoTdaV5Client(region)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	detail, err := client.ShowQueue(&model.ShowQueueRequest{QueueId: d.Id()})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving IoTDA AMQP queue")
	}

	// When the queue does not exist, it still returns a empty struct
	if detail == nil || detail.QueueId == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving IoTDA AMQP queue")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", detail.QueueName),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceAmqpDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	region := c.GetRegion(d)
	client, err := c.HcIoTdaV5Client(region)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	deleteOpts := &model.DeleteQueueRequest{
		QueueId: d.Id(),
	}
	_, err = client.DeleteQueue(deleteOpts)
	if err != nil {
		return diag.Errorf("error deleting IoTDA AMQP queue: %s", err)
	}

	return nil
}
