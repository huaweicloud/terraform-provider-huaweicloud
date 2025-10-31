package iam

import (
	"context"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// DataSourceIdentityProjectQuota
// @API IAM GET /v3.0/OS-QUOTA/projects/{project_id}
func DataSourceIdentityProjectQuota() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceIdentityProjectQuotaRead,

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"resources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of resource, must be `project`",
						},
						"max": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum Quota",
						},
						"min": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Minimum Quota",
						},
						"quota": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Current quota",
						},
						"used": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Used quota",
						},
					},
				},
			},
		},
	}
}

func DataSourceIdentityProjectQuotaRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}
	showProjectQuotaPath := iamClient.Endpoint + "v3.0/OS-QUOTA/projects/{project_id}"
	showProjectQuotaPath = strings.ReplaceAll(showProjectQuotaPath, "{project_id}", d.Get("project_id").(string))
	options := golangsdk.RequestOpts{KeepResponseBody: true}
	response, err := iamClient.Request("GET", showProjectQuotaPath, &options)
	if err != nil {
		return diag.Errorf("ShowProjectQuota error : %s", err)
	}
	respBody, err := utils.FlattenResponse(response)
	if err != nil {
		return diag.FromErr(err)
	}
	id, _ := uuid.GenerateUUID()
	d.SetId(id)
	resources := utils.PathSearch("quotas.resources", respBody, make([]interface{}, 0))
	if err = d.Set("resources", resources); err != nil {
		return diag.Errorf("error setting regions fields: %s", err)
	}
	return nil
}
