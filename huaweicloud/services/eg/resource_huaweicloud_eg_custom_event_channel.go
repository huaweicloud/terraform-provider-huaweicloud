package eg

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/eg/v1/channel/custom"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API EG POST /v1/{project_id}/channels
// @API EG GET /v1/{project_id}/channels/{channel_id}
// @API EG PUT /v1/{project_id}/channels/{channel_id}
// @API EG DELETE /v1/{project_id}/channels/{channel_id}
func ResourceCustomEventChannel() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCustomEventChannelCreate,
		ReadContext:   resourceCustomEventChannelRead,
		UpdateContext: resourceCustomEventChannelUpdate,
		DeleteContext: resourceCustomEventChannelDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceChannelImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the custom event channel is located.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the custom event channel.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the custom event channel.",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The ID of the enterprise project to which the custom event channel belongs.",
			},
			"cross_account_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "The list of domain IDs (other tenants) for the cross-account policy.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the custom event channel.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The latest update time of the custom event channel.",
			},
		},
	}
}

func buildCustomEventChannelCreateOpts(d *schema.ResourceData, region, domainId, epsId string) (*custom.CreateOpts, error) {
	var (
		channelName = d.Get("name").(string)
		result      = custom.CreateOpts{
			Name:                channelName,
			Description:         d.Get("description").(string),
			EnterpriseProjectId: epsId,
		}
	)

	if accountIds, ok := d.GetOk("cross_account_ids"); ok {
		if domainId == "" {
			return nil, fmt.Errorf("unable to find the domain ID, please check whether the 'domain_id' is " +
				"configured in your script or IAM query agency is allowed")
		}
		result.CrossAccount = utils.Bool(true)
		result.Policy = &custom.CrossAccountPolicy{
			Sid:    "allow_account_to_put_events",
			Effect: "Allow",
			Principal: custom.PrincipalInfo{
				IAM: utils.ExpandToStringListBySet(accountIds.(*schema.Set)),
			},
			Action:   "eg:channels:putEvents",
			Resource: fmt.Sprintf("urn:eg:%s:%s:channel:%s", region, domainId, channelName),
		}
	}

	return &result, nil
}

func resourceCustomEventChannelCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.EgV1Client(region)
	if err != nil {
		return diag.Errorf("error creating EG v1 client: %s", err)
	}

	opts, err := buildCustomEventChannelCreateOpts(d, region, cfg.DomainID, cfg.GetEnterpriseProjectID(d))
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := custom.Create(client, *opts)
	if err != nil {
		return diag.Errorf("error creating custom event channel: %s", err)
	}
	d.SetId(resp.ID)
	return resourceCustomEventChannelRead(ctx, d, meta)
}

func resourceCustomEventChannelRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		channelId = d.Id()
	)
	client, err := cfg.EgV1Client(region)
	if err != nil {
		return diag.Errorf("error creating EG v1 client: %s", err)
	}

	resp, err := custom.Get(client, channelId, cfg.GetEnterpriseProjectID(d))
	if err != nil {
		return common.CheckDeletedDiag(d, err, "custom event channel")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", resp.Name),
		d.Set("description", resp.Description),
		d.Set("enterprise_project_id", resp.EnterpriseProjectId),
		d.Set("cross_account_ids", resp.Policy.Principal.IAM),
		d.Set("created_at", resp.CreatedTime),
		d.Set("updated_at", resp.UpdatedTime),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving EG custom event channel (%s) fields: %s", channelId, err)
	}
	return nil
}

func buildCustomEventChannelUpdateOpts(d *schema.ResourceData, region, domainId, epsId string) (*custom.UpdateOpts, error) {
	result := custom.UpdateOpts{
		ChannelId:           d.Id(),
		Description:         utils.String(d.Get("description").(string)),
		EnterpriseProjectId: epsId,
	}

	if accountIds, ok := d.GetOk("cross_account_ids"); ok {
		if domainId == "" {
			return nil, fmt.Errorf("unable to find the domain ID, please check whether the 'domain_id' is " +
				"configured in your script or IAM query agency is allowed")
		}
		result.CrossAccount = utils.Bool(true)
		result.Policy = &custom.CrossAccountPolicy{
			Sid:    "allow_account_to_put_events",
			Effect: "Allow",
			Principal: custom.PrincipalInfo{
				IAM: utils.ExpandToStringListBySet(accountIds.(*schema.Set)),
			},
			Action:   "eg:channels:putEvents",
			Resource: fmt.Sprintf("urn:eg:%s:%s:channel:%s", region, domainId, d.Get("name").(string)),
		}
	}

	return &result, nil
}

func resourceCustomEventChannelUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.EgV1Client(region)
	if err != nil {
		return diag.Errorf("error creating EG v1 client: %s", err)
	}

	opts, err := buildCustomEventChannelUpdateOpts(d, region, cfg.DomainID, cfg.GetEnterpriseProjectID(d))
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = custom.Update(client, *opts)
	if err != nil {
		return diag.Errorf("error updating custom event channel (%s): %s", d.Id(), err)
	}
	return resourceCustomEventChannelRead(ctx, d, meta)
}

func resourceCustomEventChannelDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		channelId = d.Id()
	)
	client, err := cfg.EgV1Client(region)
	if err != nil {
		return diag.Errorf("error creating EG v1 client: %s", err)
	}

	err = custom.Delete(client, channelId, cfg.GetEnterpriseProjectID(d))
	if err != nil {
		return diag.Errorf("error deleting custom event channel (%s): %s", channelId, err)
	}
	return nil
}

func resourceChannelImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	var (
		err      error
		importId = d.Id()
		parts    = strings.Split(importId, "/")
	)
	if len(parts) < 1 || len(parts) > 2 {
		return nil, fmt.Errorf("invalid resource ID format for EG channel, want 'id' (without enterprise project "+
			"association) or '<id>/<enterprise_project_id>' (with enterprise project association), but got '%s'", importId)
	}
	d.SetId(parts[0])
	if len(parts) > 1 {
		err = d.Set("enterprise_project_id", parts[1])
	}
	return []*schema.ResourceData{d}, err
}
