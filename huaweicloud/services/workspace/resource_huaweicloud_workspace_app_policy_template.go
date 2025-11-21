package workspace

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

// @API Workspace POST /v1/{project_id}/policy-templates
// @API Workspace GET /v1/{project_id}/policy-templates
// @API Workspace PUT /v1/{project_id}/policy-templates/{policy_template_id}
// @API Workspace DELETE /v1/{project_id}/policy-templates/{policy_template_id}
func ResourceAppPolicyTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppPolicyTemplateCreate,
		ReadContext:   resourceAppPolicyTemplateRead,
		UpdateContext: resourceAppPolicyTemplateUpdate,
		DeleteContext: resourceAppPolicyTemplateDelete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the policy template is located.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the policy group.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the policy template.`,
			},
			"policies": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsJSON,
				Description:  `The policies configuration in JSON format.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the policy template, in RFC3339 format.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The update time of the policy template, in RFC3339 format.`,
			},
		},
	}
}

func buildAppPolicyTemplateBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"policy_group": map[string]interface{}{
			// Required parameters
			"policy_group_name": d.Get("name"),
			"policies":          utils.StringToJson(d.Get("policies").(string)),
			// Optional parameters
			"description": d.Get("description"),
		},
	}
}

func resourceAppPolicyTemplateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	httpUrl := "v1/{project_id}/policy-templates"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildAppPolicyTemplateBodyParams(d),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating policy template: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceId := utils.PathSearch("id", respBody, "").(string)
	if resourceId == "" {
		return diag.Errorf("unable to find the policy template ID from the API response")
	}
	d.SetId(resourceId)

	return resourceAppPolicyTemplateRead(ctx, d, meta)
}

func listAppPolicyTemplates(client *golangsdk.ServiceClient) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/policy-templates"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = fmt.Sprintf("%s?limit=%d", listPath, limit)

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)

		opt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders: map[string]string{
				"Content-Type": "application/json",
			},
		}

		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		policyTemplates := utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, policyTemplates...)

		if len(policyTemplates) < limit {
			break
		}
		offset += len(policyTemplates)
	}

	return result, nil
}

func GetAppPolicyTemplateById(client *golangsdk.ServiceClient, templateId string) (interface{}, error) {
	policyTemplates, err := listAppPolicyTemplates(client)
	if err != nil {
		return nil, err
	}

	template := utils.PathSearch(fmt.Sprintf("[?id=='%s']|[0]", templateId), policyTemplates, nil)
	if template == nil {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v1/{project_id}/policy-templates",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("the policy template (%s) has been deleted", templateId)),
			},
		}
	}

	return template, nil
}

func resourceAppPolicyTemplateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	respBody, err := GetAppPolicyTemplateById(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error retrieving policy template (%s)", d.Id()))
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", respBody, nil)),
		d.Set("description", utils.PathSearch("description", respBody, nil)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("create_time",
			respBody, "").(string))/1000, false)),
		d.Set("updated_at", utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("update_time",
			respBody, "").(string))/1000, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceAppPolicyTemplateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	httpUrl := "v1/{project_id}/policy-templates/{policy_template_id}"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{policy_template_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildAppPolicyTemplateBodyParams(d),
	}

	_, err = client.Request("PATCH", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating policy template (%s): %s", d.Id(), err)
	}

	return resourceAppPolicyTemplateRead(ctx, d, meta)
}

func resourceAppPolicyTemplateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	httpUrl := "v1/{project_id}/policy-templates/{policy_template_id}"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{policy_template_id}", d.Id())

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error deleting policy template (%s): %s", d.Id(), err)
	}

	return nil
}
