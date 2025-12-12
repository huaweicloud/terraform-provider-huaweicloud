package iam

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var resourceTagNonUpdatableParams = []string{"resource_type", "resource_id"}

// ResourceIdentityV5ResourceTag
// @API IAM POST /v5/{resource_type}/{resource_id}/tags/create
// @API IAM GET /v5/{resource_type}/{resource_id}/tags
// @API IAM DELETE /v5/{resource_type}/{resource_id}/tags/delete
func ResourceIdentityV5ResourceTag() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityV5ResourceTagCreate,
		ReadContext:   resourceIdentityV5ResourceTagRead,
		UpdateContext: resourceIdentityV5ResourceTagUpdate,
		DeleteContext: resourceIdentityV5ResourceTagDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceIdentityV5ResourceTagImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(resourceTagNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"resource_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"resource_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tags": common.TagsSchema(),
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func resourceIdentityV5ResourceTagCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMNoVersionClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM Client: %s", err)
	}

	resourceType := d.Get("resource_type").(string)
	resourceId := d.Get("resource_id").(string)
	tags := d.Get("tags").(map[string]interface{})
	if err = createTags(iamClient, tags, resourceType, resourceId); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resourceType + "/" + resourceId)
	return resourceIdentityV5ResourceTagRead(ctx, d, meta)
}

func resourceIdentityV5ResourceTagRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMNoVersionClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM Client: %s", err)
	}

	resourceType := d.Get("resource_type").(string)
	resourceId := d.Get("resource_id").(string)
	getResourceTagHttpUrl := "v5/{resource_type}/{resource_id}/tags"
	getResourceTagPath := iamClient.Endpoint + getResourceTagHttpUrl
	getResourceTagPath = strings.ReplaceAll(getResourceTagPath, "{resource_type}", resourceType)
	getResourceTagPath = strings.ReplaceAll(getResourceTagPath, "{resource_id}", resourceId)
	getResourceTagOpt := golangsdk.RequestOpts{KeepResponseBody: true}
	getResourceTagResp, err := iamClient.Request("GET", getResourceTagPath, &getResourceTagOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving IAM resource tag")
	}

	getResourceTagRespBody, err := utils.FlattenResponse(getResourceTagResp)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("tags", flattenTagsToMap(utils.PathSearch("tags", getResourceTagRespBody, nil))); err != nil {
		return diag.Errorf("error saving resource tag (%s) fields: %s", d.Id(), err)
	}
	return nil
}

func resourceIdentityV5ResourceTagUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMNoVersionClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM Client: %s", err)
	}

	resourceType := d.Get("resource_type").(string)
	resourceId := d.Get("resource_id").(string)
	if d.HasChange("tags") {
		oRaw, nRaw := d.GetChange("tags")
		oMap := oRaw.(map[string]interface{})
		nMap := nRaw.(map[string]interface{})
		// remove old tags
		if len(oMap) > 0 {
			if err = deleteTags(iamClient, oMap, resourceType, resourceId); err != nil {
				return diag.FromErr(err)
			}
		}
		// set new tags
		if len(nMap) > 0 {
			if err = createTags(iamClient, nMap, resourceType, resourceId); err != nil {
				return diag.FromErr(err)
			}
		}
	}
	return resourceIdentityV5ResourceTagRead(ctx, d, meta)
}

func resourceIdentityV5ResourceTagDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMNoVersionClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM Client: %s", err)
	}

	resourceType := d.Get("resource_type").(string)
	resourceId := d.Get("resource_id").(string)
	tags := d.Get("tags").(map[string]interface{})
	if err = deleteTags(iamClient, tags, resourceType, resourceId); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceIdentityV5ResourceTagImportState(
	_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid id: %s, id must be {resource_type}/{resource_id}", d.Id())
	}
	mErr := multierror.Append(nil,
		d.Set("resource_type", parts[0]),
		d.Set("resource_id", parts[1]),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return nil, fmt.Errorf("failed to set value to state when import resource tag, %s", err)
	}
	return []*schema.ResourceData{d}, nil
}
