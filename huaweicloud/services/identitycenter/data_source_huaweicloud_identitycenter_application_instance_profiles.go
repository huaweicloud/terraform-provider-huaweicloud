package identitycenter

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IdentityCenter GET /v1/instances/{instance_id}/application-instances/{application_instance_id}/profiles
func DataSourceIdentityCenterApplicationInstanceProfiles() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceIdentityCenterApplicationInstanceProfilesRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"application_instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"profile_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIdentityCenterApplicationInstanceProfilesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		listHttpUrl = "v1/instances/{instance_id}/application-instances/{application_instance_id}/profiles"
		product     = "identitycenter"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center Client: %s", err)
	}

	listPath := client.Endpoint + listHttpUrl
	listPath = strings.ReplaceAll(listPath, "{instance_id}", d.Get("instance_id").(string))
	listPath = strings.ReplaceAll(listPath, "{application_instance_id}", d.Get("application_instance_id").(string))

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	listResp, err := client.Request("GET", listPath, &listOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Identity Center application instance profile")
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

	mErr = multierror.Append(
		mErr,
		d.Set("name", utils.PathSearch("applicationProfiles|[0].name", listRespBody, nil)),
		d.Set("status", utils.PathSearch("applicationProfiles|[0].status", listRespBody, nil)),
		d.Set("profile_id", utils.PathSearch("applicationProfiles|[0].profile_id", listRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
