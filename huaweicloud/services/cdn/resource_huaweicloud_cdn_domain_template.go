package cdn

import (
	"context"
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

// @API CDN POST /v1.0/cdn/configuration/templates
// @API CDN GET /v1.0/cdn/configuration/templates
// @API CDN PUT /v1.0/cdn/configuration/templates/{tml_id}
// @API CDN DELETE /v1.0/cdn/configuration/templates/{tml_id}
func ResourceDomainTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDomainTemplateCreate,
		ReadContext:   resourceDomainTemplateRead,
		UpdateContext: resourceDomainTemplateUpdate,
		DeleteContext: resourceDomainTemplateDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceDomainTemplateImportState,
		},

		Schema: map[string]*schema.Schema{
			// Required parameters.
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the domain template.`,
			},
			"configs": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsJSON,
				Description:  `The configuration of the domain template, in JSON format.`,
			},

			// Optional parameters.
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the domain template.`,
			},

			// Attributes.
			"type": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The type of the domain template.`,
			},
			"account_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The account ID.`,
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the domain template.`,
			},
			"modify_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The modification time of the domain template.`,
			},
		},
	}
}

func buildDomainTemplateBodyParams(d *schema.ResourceData) map[string]interface{} {
	return utils.RemoveNil(map[string]interface{}{
		"tml_name": d.Get("name").(string),
		"configs":  utils.StringToJson(d.Get("configs").(string)),
		"remark":   utils.ValueIgnoreEmpty(d.Get("description")),
	})
}

func createDomainTemplate(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	httpUrl := "v1.0/cdn/configuration/templates"
	createPath := client.Endpoint + httpUrl

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildDomainTemplateBodyParams(d),
		OkCodes:          []int{204},
	}

	_, err := client.Request("POST", createPath, &createOpt)
	return err
}

func listDomainTemplates(client *golangsdk.ServiceClient) ([]interface{}, error) {
	var (
		httpUrl = "v1.0/cdn/configuration/templates?limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := listPath + fmt.Sprintf("&offset=%d", offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		templates := utils.PathSearch("templates", respBody, make([]interface{}, 0)).([]interface{})
		if len(templates) < 1 {
			break
		}
		result = append(result, templates...)
		if len(templates) < limit {
			break
		}
		offset += len(templates)
	}

	return result, nil
}

func GetDomainTemplateById(client *golangsdk.ServiceClient, tmlId string) (interface{}, error) {
	templates, err := listDomainTemplates(client)
	if err != nil {
		return nil, err
	}

	template := utils.PathSearch(fmt.Sprintf("[?tml_id =='%s']|[0]", tmlId), templates, nil)
	if template == nil {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v1.0/cdn/configuration/templates",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("the template with ID '%s' has been removed", tmlId)),
			},
		}
	}
	return template, nil
}

func GetDomainTemplateByName(client *golangsdk.ServiceClient, tmlName string) (interface{}, error) {
	templates, err := listDomainTemplates(client)
	if err != nil {
		return nil, err
	}

	template := utils.PathSearch(fmt.Sprintf("[?tml_name =='%s']|[0]", tmlName), templates, nil)
	if template == nil {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v1.0/cdn/configuration/templates",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("the template with name '%s' has been removed", tmlName)),
			},
		}
	}
	return template, nil
}

func resourceDomainTemplateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		tmlName = d.Get("name").(string)
	)

	client, err := cfg.NewServiceClient("cdn", "")
	if err != nil {
		return diag.Errorf("error creating CDN client: %s", err)
	}

	if err := createDomainTemplate(client, d); err != nil {
		return diag.Errorf("error creating CDN domain template: %s", err)
	}

	template, err := GetDomainTemplateByName(client, tmlName)
	if err != nil {
		return diag.Errorf("unable to find the created template with name '%s': %s", tmlName, err)
	}

	tmlId := utils.PathSearch("tml_id", template, "").(string)
	if tmlId == "" {
		return diag.Errorf("unable to find the template ID from the API response")
	}

	d.SetId(tmlId)

	return resourceDomainTemplateRead(ctx, d, meta)
}

func resourceDomainTemplateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg   = meta.(*config.Config)
		tmlId = d.Id()
	)

	client, err := cfg.NewServiceClient("cdn", "")
	if err != nil {
		return diag.Errorf("error creating CDN client: %s", err)
	}

	template, err := GetDomainTemplateById(client, tmlId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CDN domain template")
	}

	mErr := multierror.Append(nil,
		d.Set("name", utils.PathSearch("tml_name", template, nil)),
		d.Set("description", utils.PathSearch("remark", template, nil)),
		d.Set("type", utils.PathSearch("type", template, nil)),
		d.Set("account_id", utils.PathSearch("account_id", template, nil)),
		d.Set("configs", utils.JsonToString(utils.PathSearch("configs", template, nil))),
		d.Set("create_time", utils.FormatTimeStampRFC3339(int64(utils.PathSearch("create_time", template,
			float64(0)).(float64))/1000, false)),
		d.Set("modify_time", utils.FormatTimeStampRFC3339(int64(utils.PathSearch("modify_time", template,
			float64(0)).(float64))/1000, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func updateDomainTemplate(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	httpUrl := "v1.0/cdn/configuration/templates/{tml_id}"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{tml_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildDomainTemplateBodyParams(d),
		OkCodes:          []int{204},
	}

	_, err := client.Request("PUT", updatePath, &updateOpt)
	return err
}

func resourceDomainTemplateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg = meta.(*config.Config)
	)

	client, err := cfg.NewServiceClient("cdn", "")
	if err != nil {
		return diag.Errorf("error creating CDN client: %s", err)
	}

	err = updateDomainTemplate(client, d)
	if err != nil {
		return diag.Errorf("error updating CDN domain template: %s", err)
	}

	return resourceDomainTemplateRead(ctx, d, meta)
}

func deleteDomainTemplate(client *golangsdk.ServiceClient, tmlId string) error {
	httpUrl := "v1.0/cdn/configuration/templates/{tml_id}"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{tml_id}", tmlId)

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes:          []int{204},
	}

	_, err := client.Request("DELETE", deletePath, &deleteOpt)
	return err
}

func resourceDomainTemplateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg   = meta.(*config.Config)
		tmlId = d.Id()
	)

	client, err := cfg.NewServiceClient("cdn", "")
	if err != nil {
		return diag.Errorf("error creating CDN client: %s", err)
	}

	err = deleteDomainTemplate(client, tmlId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error deleting CDN domain template (%s)", tmlId))
	}

	return nil
}

func resourceDomainTemplateImportState(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	importId := d.Id()

	// If the ID is a UUID, set the ID directly
	if utils.IsUUID(importId) {
		d.SetId(importId)
		return []*schema.ResourceData{d}, nil
	}

	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("cdn", "")
	if err != nil {
		return nil, fmt.Errorf("error creating CDN client: %s", err)
	}

	// If the ID is not a UUID, get the template by name and set the ID
	template, err := GetDomainTemplateByName(client, importId)
	if err != nil {
		return nil, err
	}
	d.SetId(utils.PathSearch("tml_id", template, "").(string))

	return []*schema.ResourceData{d}, nil
}
