package eg

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk/openstack/eg/v1/source/custom"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API EG POST /v1/{project_id}/sources
// @API EG GET /v1/{project_id}/sources/{source_id}
// @API EG PUT /v1/{project_id}/sources/{source_id}
// @API EG DELETE /v1/{project_id}/sources/{source_id}
func ResourceCustomEventSource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCustomEventSourceCreate,
		ReadContext:   resourceCustomEventSourceRead,
		UpdateContext: resourceCustomEventSourceUpdate,
		DeleteContext: resourceCustomEventSourceDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the custom event source is located.",
			},
			"channel_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the custom event channel to which the custom event source belong.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the custom event source.",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The type of the custom event source.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the custom event source.",
			},
			"detail": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsJSON,
				Description:  "The configuration detail of the event source, in JSON format.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the custom event source.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the custom event source.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The update time of the custom event source.",
			},
		},
	}
}

func resourceCustomEventSourceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)

		opts = custom.CreateOpts{
			ChannelId:   d.Get("channel_id").(string),
			Type:        d.Get("type").(string),
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
			Detail:      unmarshalEventSubscriptionParamsters("event source detail", d.Get("detail").(string)),
		}
	)
	client, err := cfg.EgV1Client(region)
	if err != nil {
		return diag.Errorf("error creating EG v1 client: %s", err)
	}

	resp, err := custom.Create(client, opts)
	if err != nil {
		return diag.Errorf("error creating custom event source: %s", err)
	}
	d.SetId(resp.ID)
	return resourceCustomEventSourceRead(ctx, d, meta)
}

func parseCustomEventSourceDetail(detail interface{}) interface{} {
	jsonDetail, err := json.Marshal(detail)
	if err != nil {
		log.Printf("[ERROR] unable to convert the detail of the custom event source, not json format")
		return nil
	}
	return string(jsonDetail)
}

func resourceCustomEventSourceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		sourceId = d.Id()
	)
	client, err := cfg.EgV1Client(region)
	if err != nil {
		return diag.Errorf("error creating EG v1 client: %s", err)
	}

	resp, err := custom.Get(client, sourceId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Custom Event Source")
	}

	mErr := multierror.Append(
		nil,
		d.Set("region", region),
		d.Set("channel_id", resp.ChannelId),
		d.Set("type", resp.Type),
		d.Set("name", resp.Name),
		d.Set("description", resp.Description),
		d.Set("detail", parseCustomEventSourceDetail(resp.Detail)),
		d.Set("status", resp.Status),
		d.Set("created_at", resp.CreatedTime),
		d.Set("updated_at", resp.UpdatedTime),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving EG custom event source fields: %s", err)
	}
	return nil
}

func resourceCustomEventSourceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		opts   = custom.UpdateOpts{
			SourceId:    d.Id(),
			Description: utils.String(d.Get("description").(string)),
			Detail:      unmarshalEventSubscriptionParamsters("event source detail", d.Get("detail").(string)),
		}
	)
	client, err := cfg.EgV1Client(region)
	if err != nil {
		return diag.Errorf("error creating EG v1 client: %s", err)
	}

	_, err = custom.Update(client, opts)
	if err != nil {
		return diag.Errorf("error updating custom event source: %s", err)
	}
	return resourceCustomEventSourceRead(ctx, d, meta)
}

func resourceCustomEventSourceDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		sourceId = d.Id()
	)
	client, err := cfg.EgV1Client(region)
	if err != nil {
		return diag.Errorf("error creating EG v1 client: %s", err)
	}

	err = custom.Delete(client, sourceId)
	if err != nil {
		return diag.Errorf("error deleting custom event source: %s", err)
	}
	return nil
}
