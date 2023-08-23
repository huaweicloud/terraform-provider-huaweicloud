package lts

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/lts/huawei/logstreams"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func ResourceLTSStream() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceStreamCreate,
		ReadContext:   resourceStreamRead,
		DeleteContext: resourceStreamDelete,
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

func resourceStreamCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.LtsV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	groupId := d.Get("group_id").(string)
	createOpts := &logstreams.CreateOpts{
		LogStreamName: d.Get("stream_name").(string),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	streamCreate, err := logstreams.Create(client, groupId, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating log stream: %s", err)
	}

	d.SetId(streamCreate.ID)
	return resourceStreamRead(ctx, d, meta)
}

func resourceStreamRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.LtsV2Client(region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	streamID := d.Id()
	groupID := d.Get("group_id").(string)
	streams, err := logstreams.List(client, groupID).Extract()
	if err != nil {
		notFoundDiags := diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Resource not found",
				Detail:   fmt.Sprintf("the log group %s is gone and will be removed in Terraform state.", groupID),
			},
		}

		if _, ok := err.(golangsdk.ErrDefault404); ok {
			// 404 indicates the log group is not exist
			d.SetId("")
			return notFoundDiags
		}

		if apiError, ok := err.(golangsdk.ErrDefault400); ok {
			// "LTS.0201" indicates the log group is not exist
			if resp, pErr := common.ParseErrorMsg(apiError.Body); pErr == nil && resp.ErrorCode == "LTS.0201" {
				d.SetId("")
				return notFoundDiags
			}
		}

		return diag.Errorf("error getting log stream %s: %s", streamID, err)
	}

	for _, stream := range streams.LogStreams {
		if stream.ID == streamID {
			log.Printf("[DEBUG] Retrieved log stream %s: %#v", streamID, stream)
			mErr := multierror.Append(nil,
				d.Set("region", region),
				d.Set("stream_name", stream.Name),
				d.Set("filter_count", stream.FilterCount),
			)
			return diag.FromErr(mErr.ErrorOrNil())
		}
	}

	// can not find the log stream by ID
	return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
}

func resourceStreamDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.LtsV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	groupId := d.Get("group_id").(string)
	err = logstreams.Delete(client, groupId, d.Id()).ExtractErr()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting log stream")
	}

	return nil
}
