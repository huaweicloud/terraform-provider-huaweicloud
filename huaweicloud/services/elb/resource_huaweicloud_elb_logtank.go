package elb

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/elb/v3/logtanks"
	"github.com/chnsz/golangsdk/openstack/lts/huawei/logstreams"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func ResourceLogTank() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLogTankCreate,
		ReadContext:   resourceLogTankRead,
		UpdateContext: resourceLogTankUpdate,
		DeleteContext: resourceLogTankDelete,
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
			"loadbalancer_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"log_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"log_topic_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceLogTankCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	diagnostics := checkGroupIdAndTopicId(config, d)
	if diagnostics != nil {
		return diagnostics
	}

	elbClient, err := config.ElbV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ELB v3 client: %s", err)
	}

	createOpts := logtanks.CreateOpts{
		LoadbalancerID: d.Get("loadbalancer_id").(string),
		LogGroupId:     d.Get("log_group_id").(string),
		LogTopicId:     d.Get("log_topic_id").(string),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	logTank, err := logtanks.Create(elbClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating logtank: %s", err)
	}

	d.SetId(logTank.ID)

	return resourceLogTankRead(ctx, d, meta)
}

func resourceLogTankRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	elbClient, err := config.ElbV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ELB v3 client: %s", err)
	}

	logTank, err := logtanks.Get(elbClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "logtanks")
	}

	log.Printf("[DEBUG] Retrieved logtank %s: %#v", d.Id(), logTank)

	mErr := multierror.Append(nil,
		d.Set("loadbalancer_id", logTank.LoadbalancerID),
		d.Set("log_group_id", logTank.LogGroupId),
		d.Set("log_topic_id", logTank.LogTopicId),
		d.Set("region", config.GetRegion(d)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting Dedicated ELB logtank fields: %s", err)
	}

	return nil
}

func resourceLogTankUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	diagnostics := checkGroupIdAndTopicId(config, d)
	if diagnostics != nil {
		return diagnostics
	}

	elbClient, err := config.ElbV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ELB v3 client: %s", err)
	}

	var updateOpts logtanks.UpdateOpts
	if d.HasChange("log_group_id") {
		updateOpts.LogGroupId = d.Get("log_group_id").(string)
	}
	if d.HasChange("log_topic_id") {
		updateOpts.LogTopicId = d.Get("log_topic_id").(string)
	}

	log.Printf("[DEBUG] Updating logtank %s with options: %#v", d.Id(), updateOpts)
	_, err = logtanks.Update(elbClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return diag.Errorf("unable to update logtank %s: %s", d.Id(), err)
	}

	return resourceLogTankRead(ctx, d, meta)
}

func resourceLogTankDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	elbClient, err := config.ElbV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ELB v3 client: %s", err)
	}

	log.Printf("[DEBUG] Attempting to delete logtank %s", d.Id())
	err = logtanks.Delete(elbClient, d.Id()).ExtractErr()
	if err != nil {
		return diag.Errorf("unable to delete logtank %s: %s", d.Id(), err)
	}
	return nil
}

func checkGroupIdAndTopicId(config *config.Config, d *schema.ResourceData) diag.Diagnostics {
	logStreamClient, err := config.LtsV2Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}
	logGroupId := d.Get("log_group_id").(string)
	logTopicId := d.Get("log_topic_id").(string)
	streams, err := logstreams.List(logStreamClient, logGroupId).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault400); ok {
			return diag.Errorf("the log group id %s is error: %s", logGroupId, err)
		}
		return diag.Errorf("error getting LTS log stream by log group id %s: %s", logGroupId, err)
	}
	containLogTopicId := false
	for _, stream := range streams.LogStreams {
		if stream.ID == logTopicId {
			containLogTopicId = true
			break
		}
	}
	if !containLogTopicId {
		return diag.Errorf("the log topic id %s not belong to the group id %s", logTopicId, logGroupId)
	}
	return nil
}
