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

var v5ResourceTagNonUpdatableParams = []string{"resource_type", "resource_id"}

// @API IAM POST /v5/{resource_type}/{resource_id}/tags/create
// @API IAM GET /v5/{resource_type}/{resource_id}/tags
// @API IAM DELETE /v5/{resource_type}/{resource_id}/tags/delete
func ResourceV5ResourceTag() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV5ResourceTagCreate,
		ReadContext:   resourceV5ResourceTagRead,
		UpdateContext: resourceV5ResourceTagUpdate,
		DeleteContext: resourceV5ResourceTagDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceV5ResourceTagImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(v5ResourceTagNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"resource_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The resource type to be associated with the tags.`,
			},
			"resource_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The resource ID to be associated with the tags.`,
			},
			"tags": common.TagsSchema(`The key/value pairs to associated with the resource.`),
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func resourceV5ResourceTagCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("iam", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	resourceType := d.Get("resource_type").(string)
	resourceId := d.Get("resource_id").(string)
	tags := d.Get("tags").(map[string]interface{})
	if err = createTags(client, tags, resourceType, resourceId); err != nil {
		return diag.Errorf("error adding tags to %s (%s): %s", resourceType, resourceId, err)
	}

	d.SetId(resourceType + "/" + resourceId)
	return resourceV5ResourceTagRead(ctx, d, meta)
}

func GetV5ResourceTagsById(client *golangsdk.ServiceClient, resourceType string, resourceId string) ([]interface{}, error) {
	tags, err := getV5ResourceTags(client, resourceType, resourceId)
	if len(tags) == 0 {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v5/{resource_type}/{resource_id}/tags",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("the resource (%s/%s) does not have any tags", resourceType, resourceId)),
			},
		}
	}

	return tags, err
}

func resourceV5ResourceTagRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("iam", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	tags, err := GetV5ResourceTagsById(client, d.Get("resource_type").(string), d.Get("resource_id").(string))
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving IAM resource tag")
	}

	if err = d.Set("tags", flattenTagsToMap(tags)); err != nil {
		return diag.Errorf("error saving resource tag (%s) fields: %s", d.Id(), err)
	}
	return nil
}

func resourceV5ResourceTagUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("iam", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	resourceType := d.Get("resource_type").(string)
	resourceId := d.Get("resource_id").(string)
	if d.HasChange("tags") {
		oRaw, nRaw := d.GetChange("tags")
		oMap := oRaw.(map[string]interface{})
		nMap := nRaw.(map[string]interface{})
		// remove old tags
		if len(oMap) > 0 {
			if err = removeV5TagsFromResource(client, oMap, resourceType, resourceId); err != nil {
				return diag.Errorf("error deleting tags of %s (%s): %s", resourceType, resourceId, err)
			}
		}
		// set new tags
		if len(nMap) > 0 {
			if err = addV5TagsToResource(client, nMap, resourceType, resourceId); err != nil {
				return diag.Errorf("error adding tags to %s (%s): %s", resourceType, resourceId, err)
			}
		}
	}

	return resourceV5ResourceTagRead(ctx, d, meta)
}

func resourceV5ResourceTagDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("iam", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	if err = removeV5TagsFromResource(
		client,
		d.Get("tags").(map[string]interface{}),
		d.Get("resource_type").(string),
		d.Get("resource_id").(string),
	); err != nil {
		return diag.Errorf("error deleting tags of resource (%s): %s", d.Id(), err)
	}

	return nil
}

func resourceV5ResourceTagImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<resource_type>/<resource_id>', but got '%s'", importedId)
	}

	mErr := multierror.Append(
		d.Set("resource_type", parts[0]),
		d.Set("resource_id", parts[1]),
	)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
