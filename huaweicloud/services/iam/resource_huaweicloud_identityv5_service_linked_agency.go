package iam

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// ResourceIdentityv5ServiceLinkedAgency
// @API IAM PUT /v5/service-linked-agencies
func ResourceIdentityv5ServiceLinkedAgency() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServiceLinkedAgencyCreate,
		ReadContext:   resourceServiceLinkedAgencyRead,
		DeleteContext: resourceServiceLinkedAgencyDelete,

		Schema: map[string]*schema.Schema{
			"service_principal": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"urn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"trust_policy": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"agency_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"agency_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"max_session_duration": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceServiceLinkedAgencyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("iam", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}
	serviceLinkedAgencyPath := client.Endpoint + "v5/service-linked-agencies"
	options := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"service_principal": d.Get("service_principal"),
			"description":       d.Get("description"),
		},
	}
	response, err := client.Request("PUT", serviceLinkedAgencyPath, &options)
	if err != nil {
		return diag.Errorf("error createFederationToken: %s", err)
	}
	responseBody, err := utils.FlattenResponse(response)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(utils.PathSearch("agency.agency_id", responseBody, nil).(string))
	return resourceServiceLinkedAgencyRead(ctx, d, meta)
}

func resourceServiceLinkedAgencyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("iam", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	path := client.Endpoint + "v5/agencies/{agency_id}"
	path = strings.ReplaceAll(path, "{agency_id}", d.Id())
	reqOpt := &golangsdk.RequestOpts{KeepResponseBody: true}
	r, err := client.Request("GET", path, reqOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving agency")
	}
	resp, err := utils.FlattenResponse(r)
	if err != nil {
		return diag.FromErr(err)
	}
	mErr := multierror.Append(nil,
		d.Set("urn", utils.PathSearch("agency.urn", resp, nil)),
		d.Set("trust_policy", utils.PathSearch("agency.trust_policy", resp, nil)),
		d.Set("agency_id", utils.PathSearch("agency.agency_id", resp, nil)),
		d.Set("agency_name", utils.PathSearch("agency.agency_name", resp, nil)),
		d.Set("path", utils.PathSearch("agency.path", resp, nil)),
		d.Set("created_at", utils.PathSearch("agency.created_at", resp, nil)),
		d.Set("max_session_duration", utils.PathSearch("agency.max_session_duration", resp, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceServiceLinkedAgencyDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting service linked agency resource is not supported. The service-linked-agency is only removed " +
		"from the state, but it remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
