package smn

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/smn/v2/logtank"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// TopicURNNotExistsCode is smn error code means `topic information is not found`
const TopicURNNotExistsCode = "SMN.00010008"

// @API SMN DELETE /v2/{project_id}/notifications/topics/{topicUrn}/logtanks/{logTankID}
// @API SMN PUT /v2/{project_id}/notifications/topics/{topicUrn}/logtanks/{logTankID}
// @API SMN GET /v2/{project_id}/notifications/topics/{topicUrn}/logtanks
// @API SMN POST /v2/{project_id}/notifications/topics/{topicUrn}/logtanks
func ResourceSmnLogtank() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSmnLogtankCreate,
		UpdateContext: resourceSmnLogtankUpdate,
		ReadContext:   ResourceSmnLogtankRead,
		DeleteContext: resourceSmnLogtankDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSmnLogtankImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"topic_urn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"log_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"log_stream_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"logtank_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSmnLogtankCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.SmnV2Client(region)
	if err != nil {
		return diag.Errorf("error creating SMN client: %s", err)
	}

	createOps := logtank.Opts{
		LogGroupID:  d.Get("log_group_id").(string),
		LogStreamID: d.Get("log_stream_id").(string),
	}
	topicUrn := d.Get("topic_urn").(string)

	result, err := logtank.Create(client, topicUrn, createOps).Extract()
	if err != nil {
		return diag.Errorf("error creating SMN logtank: %s", err)
	}

	d.SetId(topicUrn)
	mErr := multierror.Append(nil, d.Set("logtank_id", result.ID))
	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error creating SMN logtank when set logtank_id: %s", mErr)
	}
	return ResourceSmnLogtankRead(ctx, d, meta)
}

func resourceSmnLogtankUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.SmnV2Client(region)
	if err != nil {
		return diag.Errorf("error creating SMN client: %s", err)
	}

	topicUrn := d.Id()
	id := d.Get("logtank_id").(string)
	if d.HasChanges("log_group_id", "log_stream_id") {
		updateOpts := logtank.Opts{
			LogGroupID:  d.Get("log_group_id").(string),
			LogStreamID: d.Get("log_stream_id").(string),
		}
		_, err = logtank.Update(client, topicUrn, id, updateOpts).Extract()
		if err != nil {
			return diag.Errorf("error updating logtank: %s", err)
		}
	}

	return ResourceSmnLogtankRead(ctx, d, meta)
}

func ResourceSmnLogtankRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.SmnV2Client(region)
	if err != nil {
		return diag.Errorf("error creating SMN client: %s", err)
	}

	topicUrn := d.Id()
	logtanks, err := logtank.List(client, topicUrn).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", TopicURNNotExistsCode),
			"error retrieving SMN logtank")
	}
	logtankID := d.Get("logtank_id").(string)
	logtankGet := GetLogtankById(logtanks, logtankID)
	if logtankGet == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving SMN logtank")
	}

	mErr := multierror.Append(
		d.Set("logtank_id", logtankID),
		d.Set("region", region),
		d.Set("topic_urn", topicUrn),
		d.Set("log_group_id", logtankGet.LogGroupID),
		d.Set("log_stream_id", logtankGet.LogStreamID),
		d.Set("updated_at", logtankGet.UpdateTime),
		d.Set("created_at", logtankGet.CreateTime),
	)

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting logtank: %s", err)
	}
	return nil
}

func GetLogtankById(logtanks []logtank.LogtankGet, id string) *logtank.LogtankGet {
	if len(logtanks) == 0 {
		return nil
	}
	if id == "" {
		return &logtanks[0]
	}
	for _, logtankItem := range logtanks {
		if logtankItem.ID == id {
			return &logtankItem
		}
	}
	return nil
}

func resourceSmnLogtankDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.SmnV2Client(region)
	if err != nil {
		return diag.Errorf("error creating SMN client: %s", err)
	}

	topicUrn := d.Id()
	logtankID := d.Get("logtank_id").(string)

	if err = logtank.Delete(client, topicUrn, logtankID).ExtractErr(); err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting SMN logtank")
	}

	return nil
}

func resourceSmnLogtankImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	var mErr *multierror.Error
	parts := strings.Split(d.Id(), "/")
	if len(parts) < 1 || len(parts) > 2 {
		return nil, fmt.Errorf("the imported ID specifies an invalid format, must be <topic_urn> or <topic_urn>/<id>")
	}
	d.SetId(parts[0])
	if len(parts) == 2 {
		mErr = multierror.Append(d.Set("logtank_id", parts[1]))
	}
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
