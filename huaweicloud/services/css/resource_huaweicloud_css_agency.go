package css

import (
	"context"
	"errors"
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

var agencyNameType = map[string]string{
	"vpc": "css_upgrade_agency",
	"obs": "css_obs_agency",
	"elb": "css_elb_agency",
	"smn": "css_smn_agency",
}

var agencyNonUpdatableParams = []string{
	"domain_id",
	"domain_name",
	"type",
}

// @API CSS POST /v1.0/{project_id}/agency/create
// @API IAM GET /v3.0/OS-AGENCY/agencies
// @API IAM GET /v3.0/OS-AGENCY/agencies/{agency_id}
// @API IAM DELETE /v3.0/OS-AGENCY/agencies/{agency_id}
func ResourceAgency() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAgencyCreate,
		ReadContext:   resourceAgencyRead,
		UpdateContext: resourceAgencyUpdate,
		DeleteContext: resourceAgencyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(agencyNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"domain_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"domain_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func buildAgencyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"domain_id":   d.Get("domain_id"),
		"domain_name": d.Get("domain_name"),
		"type":        d.Get("type"),
	}

	return bodyParams
}

func resourceAgencyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		domainId   = d.Get("domain_id").(string)
		agencyType = d.Get("type").(string)
		agencyName = ""
		httpUrl    = "v1.0/{project_id}/agency/create"
	)

	cssClient, err := cfg.NewServiceClient("css", region)
	if err != nil {
		return diag.Errorf("error creating CSS client: %s", err)
	}

	iamClient, err := cfg.NewServiceClient("iam", region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	createPath := cssClient.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", cssClient.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildAgencyBodyParams(d),
	}

	_, err = cssClient.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating agency: %s", err)
	}

	if v, ok := agencyNameType[agencyType]; ok {
		agencyName = v
	}

	agency, err := getAgencyId(iamClient, domainId, agencyName)
	if err != nil {
		return diag.FromErr(err)
	}

	agencyId := utils.PathSearch("id", agency, "").(string)
	if agencyId == "" {
		return diag.Errorf("error creating agency: unable to find agency ID")
	}

	d.SetId(agencyId)

	return resourceAgencyRead(ctx, d, meta)
}

func getAgencyId(client *golangsdk.ServiceClient, domainId, agencyName string) (interface{}, error) {
	httpUrl := "v3.0/OS-AGENCY/agencies?domain_id={domain_id}&name={name}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{domain_id}", domainId)
	getPath = strings.ReplaceAll(getPath, "{name}", agencyName)
	listOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json:charset=utf8",
		},
	}

	resp, err := client.Request("GET", getPath, &listOpts)
	if err != nil {
		return nil, fmt.Errorf("error retrieving the agency list: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	agency := utils.PathSearch("agencies[0]", respBody, nil)
	if agency == nil {
		return nil, errors.New("error creating agency: unable to find agency in list API")
	}

	return agency, nil
}

func GetAgency(client *golangsdk.ServiceClient, agencyId string) (interface{}, error) {
	httpUrl := "v3.0/OS-AGENCY/agencies/{agency_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{agency_id}", agencyId)
	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json:charset=utf8",
		},
	}

	resp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func resourceAgencyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("iam", region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	agency, err := GetAgency(client, d.Id())
	if err != nil {
		// When the agency does not exist, the response HTTP status code of the query API is `404`
		return common.CheckDeletedDiag(d, err, "error retrieving agency")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("domain_id", utils.PathSearch("agency.domain_id", agency, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceAgencyUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceAgencyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3.0/OS-AGENCY/agencies/{agency_id}"
	)

	client, err := cfg.NewServiceClient("iam", region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{agency_id}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json:charset=utf8",
		},
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		// If the agency does not exist, the response HTTP status code of the deletion API is `404`.
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error deleting agency, the error message: %s", err))
	}

	return nil
}
