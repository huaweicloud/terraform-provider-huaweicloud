package er

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/er/v3/vpcattachments"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ER POST /v3/{project_id}/enterprise-router/{er_id}/vpc-attachments
// @API ER PUT /v3/{project_id}/enterprise-router/{er_id}/vpc-attachments/{vpc_attachment_id}
// @API ER DELETE /v3/{project_id}/enterprise-router/{er_id}/vpc-attachments/{vpc_attachment_id}
// @API ER GET /v3/{project_id}/enterprise-router/{er_id}/vpc-attachments/{vpc_attachment_id}
func ResourceVpcAttachment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVpcAttachmentCreate,
		UpdateContext: resourceVpcAttachmentUpdate,
		ReadContext:   resourceVpcAttachmentRead,
		DeleteContext: resourceVpcAttachmentDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceVpcAttachmentImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},

		CustomizeDiff: config.MergeDefaultTags(),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the ER instance and the VPC attachment are located.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the ER instance to which the VPC attachment belongs.`,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the VPC to which the VPC attachment belongs.`,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the VPC subnet to which the VPC attachment belongs.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the VPC attachment.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the VPC attachment.`,
			},
			"auto_create_vpc_routes": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Whether to automatically configure routes for the VPC which pointing to the ER instance.`,
			},
			"tags": common.TagsSchema(),
			// Attributes
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The current status of the VPC attachment.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The latest update time.`,
			},
		},
	}
}

func resourceVpcAttachmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ErV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ER v3 client: %s", err)
	}

	opts := vpcattachments.CreateOpts{
		VpcId:               d.Get("vpc_id").(string),
		SubnetId:            d.Get("subnet_id").(string),
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		AutoCreateVpcRoutes: utils.Bool(d.Get("auto_create_vpc_routes").(bool)),
		Tags:                utils.ExpandResourceTags(d.Get("tags").(map[string]interface{})),
	}
	instanceId := d.Get("instance_id").(string)
	resp, err := vpcattachments.Create(client, instanceId, opts)
	if err != nil {
		return diag.Errorf("error creating VPC attachment: %s", err)
	}
	d.SetId(resp.ID)

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      vpcAttachmentStatusRefreshFunc(client, instanceId, d.Id(), []string{"available", "initiating_request"}),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceVpcAttachmentRead(ctx, d, meta)
}

func resourceVpcAttachmentRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		instanceId   = d.Get("instance_id").(string)
		attachmentId = d.Id()
	)

	client, err := cfg.ErV3Client(region)
	if err != nil {
		return diag.Errorf("error creating ER v3 client: %s", err)
	}

	resp, err := vpcattachments.Get(client, instanceId, attachmentId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "ER VPC attachment")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("vpc_id", resp.VpcId),
		d.Set("subnet_id", resp.SubnetId),
		d.Set("name", resp.Name),
		d.Set("description", resp.Description),
		d.Set("auto_create_vpc_routes", resp.AutoCreateVpcRoutes),
		d.Set("tags", utils.TagsToMap(resp.Tags)),
		d.Set("status", resp.Status),
		// The time results are not the time in RF3339 format without milliseconds.
		d.Set("created_at", utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(resp.CreatedAt)/1000, false)),
		d.Set("updated_at", utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(resp.UpdatedAt)/1000, false)),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error saving VPC attachment (%s) fields: %s", d.Id(), mErr)
	}
	return nil
}

func updateVpcAttachmentBasicInfo(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		instanceId   = d.Get("instance_id").(string)
		attachmentId = d.Id()
	)

	opts := vpcattachments.UpdateOpts{
		Name:        d.Get("name").(string),
		Description: utils.String(d.Get("description").(string)),
	}

	_, err := vpcattachments.Update(client, instanceId, d.Id(), opts)
	if err != nil {
		return fmt.Errorf("error getting VPC attachment (%s) details: %s", d.Id(), err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      vpcAttachmentStatusRefreshFunc(client, instanceId, attachmentId, []string{"available", "initiating_request"}),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	return err
}

func vpcAttachmentStatusRefreshFunc(client *golangsdk.ServiceClient, instanceId, attachmentId string, targets []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := vpcattachments.Get(client, instanceId, attachmentId)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok && len(targets) < 1 {
				return resp, "COMPLETED", nil
			}

			return nil, "", err
		}
		log.Printf("[DEBUG] The details of the VPC attachment (%s) is: %#v", attachmentId, resp)

		if utils.StrSliceContains([]string{"failed"}, resp.Status) {
			return resp, "", fmt.Errorf("unexpected status '%s'", resp.Status)
		}
		if utils.StrSliceContains(targets, resp.Status) {
			return resp, "COMPLETED", nil
		}

		return resp, "PENDING", nil
	}
}

func resourceVpcAttachmentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ErV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ER v3 client: %s", err)
	}

	if d.HasChanges("name", "description") {
		if err = updateVpcAttachmentBasicInfo(ctx, client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("tags") {
		err = utils.UpdateResourceTags(client, d, "vpc-attachment", d.Id())
		if err != nil {
			return diag.Errorf("error updating VPC attachment tags: %s", err)
		}
	}

	return resourceVpcAttachmentRead(ctx, d, meta)
}

func resourceVpcAttachmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ErV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ER v3 client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	attachmentId := d.Id()

	err = vpcattachments.Delete(client, instanceId, attachmentId)
	if err != nil {
		return diag.Errorf("error deleting VPC attachment (%s) form the ER instance: %s", attachmentId, err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      vpcAttachmentStatusRefreshFunc(client, instanceId, attachmentId, nil),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceVpcAttachmentImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format for import ID, want '<instance_id>/<attachment_id>', but '%s'", d.Id())
	}

	d.SetId(parts[1])
	return []*schema.ResourceData{d}, d.Set("instance_id", parts[0])
}
