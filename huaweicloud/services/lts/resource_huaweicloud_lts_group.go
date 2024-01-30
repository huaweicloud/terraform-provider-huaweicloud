package lts

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/lts/huawei/loggroups"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API LTS POST /v2/{project_id}/groups/{id}
// @API LTS DELETE /v2/{project_id}/groups/{id}
// @API LTS POST /v2/{project_id}/groups
// @API LTS GET /v2/{project_id}/groups
func ResourceLTSGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGroupCreate,
		ReadContext:   resourceGroupRead,
		UpdateContext: resourceGroupUpdate,
		DeleteContext: resourceGroupDelete,
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
			"group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ttl_in_days": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"tags": common.TagsSchema(),

			// Attributes
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.LtsV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	createOpts := &loggroups.CreateOpts{
		LogGroupName: d.Get("group_name").(string),
		TTL:          d.Get("ttl_in_days").(int),
		Tags:         utils.ExpandResourceTags(d.Get("tags").(map[string]interface{})),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	groupCreate, err := loggroups.Create(client, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating log group: %s", err)
	}

	d.SetId(groupCreate.ID)
	return resourceGroupRead(ctx, d, meta)
}

func resourceGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.LtsV2Client(region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	groups, err := loggroups.List(client).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error getting log group")
	}

	groupID := d.Id()
	for _, group := range groups.LogGroups {
		if group.ID == groupID {
			log.Printf("[DEBUG] Retrieved log group %s: %#v", groupID, group)
			mErr := multierror.Append(nil,
				d.Set("region", region),
				d.Set("group_name", group.Name),
				d.Set("ttl_in_days", group.TTLinDays),
				d.Set("tags", group.Tags),
				d.Set("created_at", utils.FormatTimeStampRFC3339(group.CreationTime/1000, false)),
			)
			return diag.FromErr(mErr.ErrorOrNil())
		}
	}

	// can not find the log group by ID
	return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
}

func resourceGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.LtsV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	updateOpts := &loggroups.UpdateOpts{
		TTL: d.Get("ttl_in_days").(int),
	}

	if d.HasChanges("tags") {
		// NOTE: the key in tags can not be removed due to the API restrictions.
		tagRaw := d.Get("tags").(map[string]interface{})
		taglist := utils.ExpandResourceTags(tagRaw)
		updateOpts.Tags = taglist
	}

	log.Printf("[DEBUG] Update Options: %#v", updateOpts)
	_, err = loggroups.Update(client, updateOpts, d.Id()).Extract()
	if err != nil {
		return diag.Errorf("error updating log group: %s", err)
	}

	return resourceGroupRead(ctx, d, meta)
}

func resourceGroupDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.LtsV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	err = loggroups.Delete(client, d.Id()).ExtractErr()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting log group")
	}

	return nil
}
