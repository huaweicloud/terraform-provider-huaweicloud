package er

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/er/v3/attachments"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ER GET /v3/{project_id}/enterprise-router/{er_id}/attachments
func DataSourceAttachments() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAttachmentsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The region where the ER attachments are located.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ER instance ID to which the attachment belongs.`,
			},
			"attachment_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The specified attachment ID used to query.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The resource type to be filtered.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name used to filter the attachments.`,
			},
			"resource_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The associated resource ID.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The status used to filter the attachments.`,
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The key/value pairs used to filter the attachments.`,
			},
			// Attributes
			"attachments": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The attachment ID.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The attachment name.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the attachment.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The current status of the attachment.`,
						},
						"associated": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether this attachment has been associated.`,
						},
						"resource_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The associated resource ID.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the attachment.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The latest update time of the attachment.`,
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The key/value pairs to associate with the attachment.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The attachment type.`,
						},
						"route_table_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The associated route table ID.`,
						},
					},
				},
				Description: `All attachments that match the filter parameters.`,
			},
		},
	}
}

// Filter attachments by 'name' and 'attachment_id'.
func filterAttachments(d *schema.ResourceData, all []attachments.Attachment) ([]attachments.Attachment, error) {
	filter := map[string]interface{}{}
	if name, ok := d.GetOk("name"); ok {
		filter["Name"] = name
	}
	if attachmentId, ok := d.GetOk("attachment_id"); ok {
		filter["ID"] = attachmentId
	}

	if len(filter) < 1 {
		return all, nil
	}
	filterResult, err := utils.FilterSliceWithField(all, filter)
	if err != nil {
		return nil, fmt.Errorf("error filtering attachment list: %s", err)
	}
	result := make([]attachments.Attachment, 0, len(filterResult))
	for _, val := range filterResult {
		result = append(result, val.(attachments.Attachment))
	}
	return result, nil
}

func filterAttachmentsByTags(d *schema.ResourceData, all []attachments.Attachment) ([]attachments.Attachment, error) {
	tagFilter, ok := d.GetOk("tags")
	if !ok {
		return all, nil
	}
	result := make([]attachments.Attachment, 0, len(all))
	for _, val := range all {
		tagmap := utils.TagsToMap(val.Tags)

		// Filter attachment list by tags, if the filter is nil, skip and return all fileterResult elements.
		if utils.HasMapContains(tagmap, tagFilter.(map[string]interface{})) {
			result = append(result, val)
		}
	}
	return result, nil
}

func flattenAttachments(all []attachments.Attachment) []map[string]interface{} {
	if len(all) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(all))
	for i, attachment := range all {
		result[i] = map[string]interface{}{
			"id":          attachment.ID,
			"name":        attachment.Name,
			"description": attachment.Description,
			"status":      attachment.Status,
			"associated":  attachment.Associated,
			"resource_id": attachment.ResourceId,
			// The time results are not the time in RF3339 format without milliseconds.
			"created_at":     utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(attachment.CreatedAt)/1000, false),
			"updated_at":     utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(attachment.UpdatedAt)/1000, false),
			"tags":           utils.TagsToMap(attachment.Tags),
			"type":           attachment.ResourceType,
			"route_table_id": attachment.RouteTableId,
		}
	}
	return result
}

func buildAttachmentListOpts(d *schema.ResourceData) attachments.ListOpts {
	return attachments.ListOpts{
		Statuses:      buildSliceIgnoreEmptyElement(d.Get("status").(string)),
		ResourceTypes: buildSliceIgnoreEmptyElement(d.Get("type").(string)),
		ResourceIds:   buildSliceIgnoreEmptyElement(d.Get("resource_id").(string)),
		SortKey:       []string{"name"},
	}
}

func dataSourceAttachmentsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.ErV3Client(region)
	if err != nil {
		return diag.Errorf("error creating ER v3 client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	resp, err := attachments.List(client, instanceId, buildAttachmentListOpts(d))
	if err != nil {
		return diag.Errorf("error retrieving attachments: %s", err)
	}
	if resp, err = filterAttachments(d, resp); err != nil {
		return diag.FromErr(err)
	}
	if resp, err = filterAttachmentsByTags(d, resp); err != nil {
		return diag.FromErr(err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("attachments", flattenAttachments(resp)),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error saving attachments data source fields: %s", mErr)
	}
	return nil
}
