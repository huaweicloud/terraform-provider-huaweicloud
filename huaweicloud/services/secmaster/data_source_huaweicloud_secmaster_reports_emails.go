package secmaster

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Secmaster POST /v1/{project_id}/workspaces/{workspace_id}/sa/reports/emails/search
func DataSourceReportsEmails() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceReportsEmailsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"email_address": {
				Type:     schema.TypeString,
				Required: true,
			},
			"emails": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"report_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"email_status": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildReportsEmailsBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"email_address": d.Get("email_address"),
	}

	return bodyParams
}

func dataSourceReportsEmailsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/sa/reports/emails/search"
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{workspace_id}", workspaceId)

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildReportsEmailsBodyParams(d),
	}

	listResp, err := client.Request("POST", listPath, &listOpt)
	if err != nil {
		return diag.Errorf("error retrieving the recipient email status: %s", err)
	}

	listRespBody, err := utils.FlattenResponse(listResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("emails", flattenReportsEmails(listRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenReportsEmails(emailResp interface{}) []interface{} {
	if emailInfos, ok := emailResp.([]interface{}); ok {
		rst := make([]interface{}, 0, len(emailInfos))
		for _, v := range emailInfos {
			rst = append(rst, map[string]interface{}{
				"report_address": utils.PathSearch("address", v, nil),
				"email_status":   utils.PathSearch("status", v, nil),
			})
		}
		return rst
	}

	return make([]interface{}, 0)
}
