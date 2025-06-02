package smn

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var topicAttributesNonUpdatableParams = []string{
	"topic_urn",
	"name",
}

// @API SMN PUT /v2/{project_id}/notifications/topics/{topic_urn}/attributes/{name}
// @API SMN GET /v2/{project_id}/notifications/topics/{topic_urn}/attributes
func ResourceTopicAttributes() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTopicAttributesCreate,
		ReadContext:   resourceTopicAttributesRead,
		UpdateContext: resourceTopicAttributesUpdate,
		DeleteContext: resourceTopicAttributesDelete,

		CustomizeDiff: config.FlexibleForceNew(topicAttributesNonUpdatableParams),

		Importer: &schema.ResourceImporter{
			StateContext: resourceTopicAttributesImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the topic is located.`,
			},
			"topic_urn": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The topic URN.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The topic attribute name.`,
			},
			"value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The topic attribute value, in JSON format.`,
			},
		},
	}
}

func updateTopicAttribute(client *golangsdk.ServiceClient, topicUrn, name, value string) error {
	httpUrl := "v2/{project_id}/notifications/topics/{topic_urn}/attributes/{name}"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{topic_urn}", topicUrn)
	updatePath = strings.ReplaceAll(updatePath, "{name}", name)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: map[string]interface{}{
			"value": value,
		},
	}

	_, err := client.Request("PUT", updatePath, &opt)
	return err
}

func resourceTopicAttributesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.SmnV2Client(region)
	if err != nil {
		return diag.Errorf("error creating SMN client: %s", err)
	}

	topicUrn := d.Get("topic_urn").(string)
	name := d.Get("name").(string)
	value := d.Get("value").(string)

	if err := updateTopicAttribute(client, topicUrn, name, value); err != nil {
		return diag.Errorf("error setting attributes (names %s) for topic %s: %s", name, topicUrn, err)
	}

	d.SetId(fmt.Sprintf("%s/%s", topicUrn, name))

	return resourceTopicAttributesRead(ctx, d, meta)
}

func GetTopicAttributes(client *golangsdk.ServiceClient, topicUrn, name string) (interface{}, error) {
	httpUrl := "v2/{project_id}/notifications/topics/{topic_urn}/attributes?name={name}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{topic_urn}", topicUrn)
	getPath = strings.ReplaceAll(getPath, "{name}", name)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	requestResp, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch(fmt.Sprintf("attributes.%s", name), respBody, nil), nil
}

func resourceTopicAttributesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.SmnV2Client(region)
	if err != nil {
		return diag.Errorf("error creating SMN client: %s", err)
	}

	topicUrn := d.Get("topic_urn").(string)
	name := d.Get("name").(string)

	attributes, err := GetTopicAttributes(client, topicUrn, name)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving SMN topic attributes")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("value", attributes),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceTopicAttributesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.SmnV2Client(region)
	if err != nil {
		return diag.Errorf("error creating SMN client: %s", err)
	}

	topicUrn := d.Get("topic_urn").(string)
	name := d.Get("name").(string)
	value := d.Get("value").(string)

	if err := updateTopicAttribute(client, topicUrn, name, value); err != nil {
		return diag.Errorf("error updating attributes (names %s) for topic %s: %s", name, topicUrn, err)
	}

	return resourceTopicAttributesRead(ctx, d, meta)
}

func resourceTopicAttributesDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `Deleting this resource will not reset the topic attributes, but will only remove the resource
information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

func resourceTopicAttributesImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<topic_urn>/<name>', but got '%s'", d.Id())
	}

	topicUrn := parts[0]
	name := parts[1]

	mErr := multierror.Append(nil,
		d.Set("topic_urn", topicUrn),
		d.Set("name", name),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
