package das

import (
	"context"
	"errors"
	"fmt"
	"strconv"
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

var emailTemplateNonUpdatableParams = []string{
	"datastore_type",
}

// @API DAS POST /v3/{project_id}/batch-inspection/email-template
// @API DAS GET /v3/{project_id}/batch-inspection/email-template
// @API DAS PUT /v3/{project_id}/batch-inspection/email-template
// @API DAS DELETE /v3/{project_id}/batch-inspection/email-template
func ResourceEmailTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEmailTemplateCreate,
		ReadContext:   resourceEmailTemplateRead,
		UpdateContext: resourceEmailTemplateUpdate,
		DeleteContext: resourceEmailTemplateDelete,

		CustomizeDiff: config.FlexibleForceNew(emailTemplateNonUpdatableParams),

		Importer: &schema.ResourceImporter{
			StateContext: resourceEmailTemplateImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the email template is located.`,
			},

			// Required parameters.
			"datastore_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The database type.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the email template.`,
			},
			"groups": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: `The list of instance group IDs.`,
			},
			"health_rank": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of health ranks.`,
			},
			"inspection_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The diagnosis time.`,
			},
			"send_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The send time.`,
			},
			"time_zone": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The time zone.`,
			},

			// Optional parameters.
			"email": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: `The email address.`,
			},
			"topic": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The topic ID.`,
			},
			"topic_urn": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The topic URN.`,
			},
			"obs_bucket_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The OBS bucket name.`,
			},

			// Attributes.
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The update time, in RFC3339 format.`,
			},
			"user_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the user who last modified the template.`,
			},
			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The status of the email template.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					`Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.`,
					utils.SchemaDescInput{Internal: true},
				),
			},
		},
	}
}

func buildEmailTemplateCreateBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"datastore_type":  d.Get("datastore_type"),
		"template_name":   d.Get("name"),
		"group_id":        d.Get("groups"),
		"health_rank":     d.Get("health_rank"),
		"inspection_time": d.Get("inspection_time"),
		"send_time":       d.Get("send_time"),
		"time_zone":       d.Get("time_zone"),
		"email":           utils.ValueIgnoreEmpty(d.Get("email")),
		"topic":           utils.ValueIgnoreEmpty(d.Get("topic")),
		"topic_urn":       utils.ValueIgnoreEmpty(d.Get("topic_urn")),
		"obs_bucket_name": utils.ValueIgnoreEmpty(d.Get("obs_bucket_name")),
	}
}

func listEmailTemplates(client *golangsdk.ServiceClient, datastoreType string) ([]interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/batch-inspection/email-template?limit={limit}"
		limit   = 200
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))
	listPath = fmt.Sprintf("%s&datastore_type=%v", listPath, datastoreType)

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &listOpt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		templates := utils.PathSearch("email_template_list", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, templates...)
		if len(templates) < limit {
			break
		}

		offset += len(templates)
	}

	return result, nil
}

func GetEmailTemplateById(client *golangsdk.ServiceClient, datastoreType, templateId string) (interface{}, error) {
	templates, err := listEmailTemplates(client, datastoreType)
	if err != nil {
		return nil, err
	}

	for _, template := range templates {
		currentTemplateId := fmt.Sprintf("%v", utils.PathSearch("template_id", template, ""))
		if currentTemplateId == templateId {
			return template, nil
		}
	}

	return nil, golangsdk.ErrDefault404{
		ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
			Method:    "GET",
			URL:       "/v3/{project_id}/batch-inspection/email-template",
			RequestId: "NONE",
			Body:      []byte(fmt.Sprintf("the email template (%s) has been removed", templateId)),
		},
	}
}

func resourceEmailTemplateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/batch-inspection/email-template"
	)

	client, err := cfg.NewServiceClient("das", region)
	if err != nil {
		return diag.Errorf("error creating DAS client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildEmailTemplateCreateBodyParams(d)),
	}

	requestResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating DAS email template: %s", err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	templateId := utils.PathSearch("template_id", respBody, float64(0)).(float64)
	if templateId == 0 {
		return diag.Errorf("unable to find the template ID from the API response")
	}
	d.SetId(strconv.Itoa(int(templateId)))

	return resourceEmailTemplateRead(ctx, d, meta)
}

func resourceEmailTemplateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		datastoreType = d.Get("datastore_type").(string)
		templateId    = d.Id()
	)

	client, err := cfg.NewServiceClient("das", region)
	if err != nil {
		return diag.Errorf("error creating DAS client: %s", err)
	}

	respBody, err := GetEmailTemplateById(client, datastoreType, templateId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error retrieving DAS email template (%s)", templateId))
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("datastore_type", utils.PathSearch("datastore_type", respBody, nil)),
		d.Set("name", utils.PathSearch("template_name", respBody, nil)),
		d.Set("groups", utils.PathSearch("group_id", respBody, make([]interface{}, 0))),
		d.Set("health_rank", utils.PathSearch("health_rank", respBody, make([]interface{}, 0))),
		d.Set("inspection_time", utils.PathSearch("inspection_time", respBody, nil)),
		d.Set("send_time", utils.PathSearch("send_time", respBody, nil)),
		d.Set("time_zone", utils.PathSearch("time_zone", respBody, nil)),
		d.Set("email", utils.PathSearch("email", respBody, nil)),
		d.Set("topic", utils.PathSearch("topic", respBody, nil)),
		d.Set("topic_urn", utils.PathSearch("topic_urn", respBody, nil)),
		d.Set("obs_bucket_name", utils.PathSearch("obs_bucket_name", respBody, nil)),
		d.Set("user_id", utils.PathSearch("user_id", respBody, nil)),
		d.Set("status", utils.PathSearch("status", respBody, nil)),
		d.Set("updated_at", utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(
			utils.PathSearch("create_time", respBody, "").(string))/1000, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildEmailTemplateUpdateBodyParams(d *schema.ResourceData) (map[string]interface{}, error) {
	templateId, err := strconv.Atoi(d.Id())
	if err != nil {
		return nil, fmt.Errorf("error parsing template ID: %s", err)
	}

	bodyParams := map[string]interface{}{
		"template_id":     templateId,
		"template_name":   utils.ValueIgnoreEmpty(d.Get("name")),
		"group_id":        utils.ValueIgnoreEmpty(d.Get("groups")),
		"email":           utils.ValueIgnoreEmpty(d.Get("email")),
		"topic":           utils.ValueIgnoreEmpty(d.Get("topic")),
		"topic_urn":       utils.ValueIgnoreEmpty(d.Get("topic_urn")),
		"health_rank":     utils.ValueIgnoreEmpty(d.Get("health_rank")),
		"obs_bucket_name": utils.ValueIgnoreEmpty(d.Get("obs_bucket_name")),
		"inspection_time": utils.ValueIgnoreEmpty(d.Get("inspection_time")),
		"send_time":       utils.ValueIgnoreEmpty(d.Get("send_time")),
		"time_zone":       utils.ValueIgnoreEmpty(d.Get("time_zone")),
	}

	return bodyParams, nil
}

func resourceEmailTemplateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/batch-inspection/email-template"
	)

	client, err := cfg.NewServiceClient("das", region)
	if err != nil {
		return diag.Errorf("error creating DAS client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updateBodyParams, err := buildEmailTemplateUpdateBodyParams(d)
	if err != nil {
		return diag.FromErr(err)
	}

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(updateBodyParams),
	}

	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating DAS email template (%s): %s", d.Id(), err)
	}

	return resourceEmailTemplateRead(ctx, d, meta)
}

func resourceEmailTemplateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/batch-inspection/email-template"

		templateId = d.Id()
	)

	client, err := cfg.NewServiceClient("das", region)
	if err != nil {
		return diag.Errorf("error creating DAS client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: map[string]interface{}{
			"template_id": templateId,
		},
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error deleting DAS email template (%s)", templateId))
	}

	return nil
}

func resourceEmailTemplateImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, errors.New("invalid format specified for import ID, must be <datastore_type>/<template_id>")
	}

	d.SetId(parts[1])
	return []*schema.ResourceData{d}, d.Set("datastore_type", parts[0])
}
