package dsc

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DSC POST /v1/{project_id}/metadata/tags
// @API DSC DELETE /v1/{project_id}/metadata/tags
func ResourceMetadataTags() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMetadataTagsCreate,
		ReadContext:   resourceMetadataTagsRead,
		UpdateContext: resourceMetadataTagsUpdate,
		DeleteContext: resourceMetadataTagsDelete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"names": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
			"msg": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildMetadataTagsBodyParams(names []interface{}) map[string]interface{} {
	return map[string]interface{}{
		"names": utils.ExpandToStringList(names),
	}
}

func createMetadataTags(client *golangsdk.ServiceClient, names []interface{}) (interface{}, error) {
	if len(names) == 0 {
		return nil, nil
	}

	requestPath := client.Endpoint + "v1/{project_id}/metadata/tags"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildMetadataTagsBodyParams(names)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func deleteMetadataTags(client *golangsdk.ServiceClient, names []interface{}) error {
	if len(names) == 0 {
		return nil
	}

	requestPath := client.Endpoint + "v1/{project_id}/metadata/tags"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildMetadataTagsBodyParams(names)),
	}

	_, err := client.Request("DELETE", requestPath, &requestOpt)
	return err
}

func resourceMetadataTagsCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		names  = d.Get("names").(*schema.Set).List()
	)

	client, err := cfg.NewServiceClient("dsc", region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	respBody, err := createMetadataTags(client, names)
	if err != nil {
		return diag.Errorf("error creating DSC metadata tags: %s", err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("names", names),
		d.Set("msg", utils.PathSearch("msg", respBody, nil)),
		d.Set("status", utils.PathSearch("status", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceMetadataTagsRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because there is no query API to verify the tags.
	return nil
}

func resourceMetadataTagsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("dsc", region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	if d.HasChange("names") {
		oldRaws, newRaws := d.GetChange("names")
		deleteNameList := oldRaws.(*schema.Set).Difference(newRaws.(*schema.Set)).List()
		addNameList := newRaws.(*schema.Set).Difference(oldRaws.(*schema.Set)).List()

		if err := deleteMetadataTags(client, deleteNameList); err != nil {
			return diag.Errorf("error deleting DSC metadata tags in update operation: %s", err)
		}

		if _, err := createMetadataTags(client, addNameList); err != nil {
			return diag.Errorf("error creating DSC metadata tags in update operation: %s", err)
		}
	}

	return resourceMetadataTagsRead(ctx, d, meta)
}

func resourceMetadataTagsDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		names  = d.Get("names").(*schema.Set).List()
	)

	client, err := cfg.NewServiceClient("dsc", region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	if err := deleteMetadataTags(client, names); err != nil {
		return diag.Errorf("error deleting DSC metadata tags: %s", err)
	}

	return nil
}
