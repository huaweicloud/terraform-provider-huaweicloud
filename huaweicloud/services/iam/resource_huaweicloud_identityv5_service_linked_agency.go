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

// @API IAM PUT /v5/service-linked-agencies
func ResourceV5ServiceLinkedAgency() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV5ServiceLinkedAgencyCreate,
		ReadContext:   resourceV5ServiceLinkedAgencyRead,
		DeleteContext: resourceV5ServiceLinkedAgencyDelete,

		Schema: map[string]*schema.Schema{
			"service_principal": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The service principal of the service-linked agency.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `The description of the service-linked agency.`,
			},
			"urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The uniform resource name of the service-linked agency.`,
			},
			"trust_policy": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The policy document of the service-linked agency, in JSON format.`,
			},
			"agency_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the service-linked agency.`,
			},
			"agency_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the service-linked agency.`,
			},
			"path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The resource path of the service-linked agency, in 'service-linked-agency/<service_principal>/' format.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the service-linked agency.`,
			},
			"max_session_duration": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The maximum session duration of the service-linked agency, in seconds.`,
			},
		},
	}
}

func resourceV5ServiceLinkedAgencyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		return diag.Errorf("error creating service-linked agency: %s", err)
	}

	responseBody, err := utils.FlattenResponse(response)
	if err != nil {
		return diag.FromErr(err)
	}

	agencyId := utils.PathSearch("agency.agency_id", responseBody, nil).(string)
	if agencyId == "" {
		return diag.Errorf("unable to find the agency ID from the API response")
	}

	d.SetId(agencyId)
	return resourceV5ServiceLinkedAgencyRead(ctx, d, meta)
}

func GetV5ServiceLinkedAgencyById(client *golangsdk.ServiceClient, agencyId string) (interface{}, error) {
	getAgencyHttpUrl := "v5/agencies/{agency_id}"
	getAgencyPath := client.Endpoint + getAgencyHttpUrl
	getAgencyPath = strings.ReplaceAll(getAgencyPath, "{agency_id}", agencyId)
	getAgencyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getAgencyResp, err := client.Request("GET", getAgencyPath, &getAgencyOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getAgencyResp)
}

func resourceV5ServiceLinkedAgencyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("iam", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	resp, err := GetV5ServiceLinkedAgencyById(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving service-linked agency")
	}

	mErr := multierror.Append(
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

func resourceV5ServiceLinkedAgencyDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for creating the service-linked agency. Deleting this resource will
not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
