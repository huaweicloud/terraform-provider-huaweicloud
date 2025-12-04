package cdn

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CDN POST /v1.0/cdn/configuration/templates/{tml_id}/apply
func ResourceDomainTemplateApply() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDomainTemplateApplyCreate,
		ReadContext:   resourceDomainTemplateApplyRead,
		UpdateContext: resourceDomainTemplateApplyUpdate,
		DeleteContext: resourceDomainTemplateApplyDelete,

		Schema: map[string]*schema.Schema{
			// Required parameters.
			"template_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the domain template.`,
			},
			"resources": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The list of domain names to apply the template. Multiple domain names are separated by commas.`,
			},
		},
	}
}

func applyDomainTemplate(client *golangsdk.ServiceClient, templateId, resources string) error {
	httpUrl := "v1.0/cdn/configuration/templates/{tml_id}/apply"
	applyPath := client.Endpoint + httpUrl
	applyPath = strings.ReplaceAll(applyPath, "{tml_id}", templateId)

	applyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: map[string]interface{}{
			"resources": resources,
		},
	}

	requestResp, err := client.Request("POST", applyPath, &applyOpt)
	if err != nil {
		return err
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return err
	}
	if utils.PathSearch("status", respBody, nil) == "success" {
		return nil
	}

	jobDetails := utils.PathSearch("detail", respBody, make([]interface{}, 0)).([]interface{})
	if len(jobDetails) < 1 {
		log.Printf("[ERROR] Unable to find the job details (with field `detail`) in the API response: %+v", respBody)
		return nil
	}

	for _, jobDetail := range jobDetails {
		if utils.PathSearch("status", jobDetail, nil) == "success" {
			continue
		}
		log.Printf("[ERROR] The job (about domain `%v`) is applied failed, the error message is: %v",
			utils.PathSearch("domain_name", jobDetail, ""), utils.PathSearch("error_msg", jobDetail, ""))
	}
	return nil
}

func resourceDomainTemplateApplyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("cdn", "")
	if err != nil {
		return diag.Errorf("error creating CDN client: %s", err)
	}

	err = applyDomainTemplate(client, d.Get("template_id").(string), d.Get("resources").(string))
	if err != nil {
		return diag.Errorf("error applying CDN domain template: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	return resourceDomainTemplateApplyRead(ctx, d, meta)
}

func resourceDomainTemplateApplyRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDomainTemplateApplyUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDomainTemplateApplyDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to apply CDN domain template. Deleting this resource
	will not clear the corresponding request record, but will only remove the resource information from the tf state
    file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
