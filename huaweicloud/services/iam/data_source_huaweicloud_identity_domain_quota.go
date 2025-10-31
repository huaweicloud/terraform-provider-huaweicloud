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

// DataSourceIdentityDomainQuota
// @API IAM GET /v3.0/OS-QUOTA/domains/{domain_id}
func DataSourceIdentityDomainQuota() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceIdentityDomainQuotaRead,

		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "Query the type of quota. If this parameter is not filled in, all types of quotas will " +
					"be returned. The value range is as follows: user, group, idp, agency, policy, " +
					"assigment_group_mp, assigment_agency_mp, assigment_group_ep, assigment_user_ep, mapping",
			},

			"resources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of resource",
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

func DataSourceIdentityDomainQuotaRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}
	showDomainQuotaPath := iamClient.Endpoint + "v3.0/OS-QUOTA/domains/{domain_id}"
	showDomainQuotaPath = strings.ReplaceAll(showDomainQuotaPath, "{domain_id}", cfg.DomainID)
	typeStr := d.Get("type").(string)
	if typeStr != "" {
		showDomainQuotaPath = showDomainQuotaPath + "?type=" + typeStr
	}
	options := golangsdk.RequestOpts{KeepResponseBody: true}
	response, err := iamClient.Request("GET", showDomainQuotaPath, &options)
	if err != nil {
		return diag.Errorf("error showDomainQuota: %s", err)
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
