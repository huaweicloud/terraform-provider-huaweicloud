package rgc

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var organizationalUnitRegisterNonUpdatableParams = []string{"organizational_unit_id"}

// @API RGC POST /v1/managed-organization/organizational-units/{organizational_unit_id}/register
// @API RGC POST /v1/managed-organization/organizational-units/{organizational_unit_id}/re-register
// @API RGC POST /v1/managed-organization/organizational-units/{organizational_unit_id}/de-register
// @API RGC GET /v1/managed-organization/managed-organizational-units/{managed_organizational_unit_id}
// @API RGC GET /v1/managed-organization/{operation_id}
func ResourceOrganizationalUnitRegister() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOrganizationalUnitRegisterCreate,
		UpdateContext: resourceOrganizationalUnitRegisterUpdate,
		ReadContext:   resourceOrganizationalUnitRegisterRead,
		DeleteContext: resourceOrganizationalUnitRegisterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(6 * time.Hour),
			Delete: schema.DefaultTimeout(6 * time.Hour),
		},

		CustomizeDiff: config.FlexibleForceNew(organizationalUnitRegisterNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"organizational_unit_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"organizational_unit_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"parent_organizational_unit_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"parent_organizational_unit_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"manage_account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"organizational_unit_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"organizational_unit_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOrganizationalUnitRegisterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	err := registerOrganizationalUnit(ctx, d, meta)
	if err != nil {
		return diag.Errorf("error register organizational unit: %s", err)
	}
	ouId := d.Get("organizational_unit_id").(string)
	d.SetId(ouId)
	return resourceOrganizationalUnitRegisterRead(ctx, d, meta)
}

func resourceOrganizationalUnitRegisterRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		getOrganizationUnitUrl     = "v1/managed-organization/managed-organizational-units/{managed_organizational_unit_id}"
		getOrganizationUnitProduct = "rgc"
	)

	getOrganizationUnitClient, err := cfg.NewServiceClient(getOrganizationUnitProduct, region)
	if err != nil {
		return diag.Errorf("error creating rgc client: %s", err)
	}

	getOrganizationUnitPath := getOrganizationUnitClient.Endpoint + getOrganizationUnitUrl
	getOrganizationUnitPath = strings.ReplaceAll(getOrganizationUnitPath, "{managed_organizational_unit_id}", d.Id())

	getOrganizationUnitOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getOrganizationUnitResp, err := getOrganizationUnitClient.Request("GET", getOrganizationUnitPath, &getOrganizationUnitOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving organizational unit")
	}

	getOrganizationUnitRespBody, err := utils.FlattenResponse(getOrganizationUnitResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("organizational_unit_id", utils.PathSearch("organizational_unit_id", getOrganizationUnitRespBody, nil)),
		d.Set("organizational_unit_name", utils.PathSearch("organizational_unit_name", getOrganizationUnitRespBody, nil)),
		d.Set("parent_organizational_unit_id", utils.PathSearch("parent_organizational_unit_id", getOrganizationUnitRespBody, nil)),
		d.Set("parent_organizational_unit_name", utils.PathSearch("parent_organizational_unit_name", getOrganizationUnitRespBody, nil)),
		d.Set("manage_account_id", utils.PathSearch("manage_account_id", getOrganizationUnitRespBody, nil)),
		d.Set("organizational_unit_type", utils.PathSearch("organizational_unit_type", getOrganizationUnitRespBody, nil)),
		d.Set("organizational_unit_status", utils.PathSearch("organizational_unit_status", getOrganizationUnitRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceOrganizationalUnitRegisterUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceOrganizationalUnitRegisterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	err := deRegisterOrganizationalUnit(ctx, d, meta)
	if err != nil {
		return diag.Errorf("error de-registing organizational unit: %s", err)
	}
	return nil
}

func registerOrganizationalUnit(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		registerOrganizationalUnitUrl     = "v1/managed-organization/organizational-units/{organizational_unit_id}/register"
		reRegisterOrganizationalUnitUrl   = "v1/managed-organization/organizational-units/{organizational_unit_id}/re-register"
		getRegisterOrganizationUnitUrl    = "v1/managed-organization/managed-organizational-units/{managed_organizational_unit_id}"
		registerOrganizationalUnitProduct = "rgc"
	)

	registerOrganizationalUnitClient, err := cfg.NewServiceClient(registerOrganizationalUnitProduct, region)
	if err != nil {
		return err
	}

	ouId := d.Get("organizational_unit_id").(string)
	getRegisterOrganizationUnitPath := registerOrganizationalUnitClient.Endpoint + getRegisterOrganizationUnitUrl
	getRegisterOrganizationUnitPath = strings.ReplaceAll(getRegisterOrganizationUnitPath, "{managed_organizational_unit_id}", ouId)

	getRegisterOrganizationUnitOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	var registerOrganizationalUnitPath string
	_, err = registerOrganizationalUnitClient.Request("GET", getRegisterOrganizationUnitPath,
		&getRegisterOrganizationUnitOpt)

	if err != nil {
		var errDefault404 golangsdk.ErrDefault404
		if !errors.As(err, &errDefault404) {
			return err
		}
		registerOrganizationalUnitPath = registerOrganizationalUnitClient.Endpoint + registerOrganizationalUnitUrl
	} else {
		registerOrganizationalUnitPath = registerOrganizationalUnitClient.Endpoint + reRegisterOrganizationalUnitUrl
	}

	registerOrganizationalUnitPath = strings.ReplaceAll(registerOrganizationalUnitPath, "{organizational_unit_id}", ouId)
	registerOrganizationalUnitOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	registerOrganizationalUnitResp, err := registerOrganizationalUnitClient.Request("POST", registerOrganizationalUnitPath,
		&registerOrganizationalUnitOpt)
	if err != nil {
		return err
	}
	registerOrganizationalUnitRespBody, err := utils.FlattenResponse(registerOrganizationalUnitResp)
	if err != nil {
		return err
	}

	operationId := utils.PathSearch("organizational_unit_operation_id", registerOrganizationalUnitRespBody, "").(string)

	stateConf := resource.StateChangeConf{
		Pending:      []string{"IN_PROGRESS"},
		Target:       []string{"SUCCEEDED"},
		Refresh:      organizationalUnitStateRefreshFunc(registerOrganizationalUnitClient, operationId),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return err
	}
	return nil
}

func deRegisterOrganizationalUnit(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		deRegisterOuHttpUrl = "v1/managed-organization/managed-organizational-units/{managed_organizational_unit_id}/de-register"
		deRegisterOuProduct = "rgc"
	)
	deRegisterOuClient, err := cfg.NewServiceClient(deRegisterOuProduct, region)
	if err != nil {
		return err
	}

	deRegisterOuPath := deRegisterOuClient.Endpoint + deRegisterOuHttpUrl
	deRegisterOuPath = strings.ReplaceAll(deRegisterOuPath, "{managed_organizational_unit_id}", d.Id())

	deRegisterOuOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	deRegisterOuResp, err := deRegisterOuClient.Request("POST", deRegisterOuPath, &deRegisterOuOpt)
	if err != nil {
		return err
	}

	deRegisterOuRespBody, err := utils.FlattenResponse(deRegisterOuResp)
	if err != nil {
		return err
	}

	operationId := utils.PathSearch("organizational_unit_operation_id", deRegisterOuRespBody, "").(string)

	stateConf := resource.StateChangeConf{
		Pending:      []string{"IN_PROGRESS"},
		Target:       []string{"SUCCEEDED"},
		Refresh:      organizationalUnitStateRefreshFunc(deRegisterOuClient, operationId),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return err
	}

	return nil
}

func organizationalUnitStateRefreshFunc(client *golangsdk.ServiceClient, operationId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getOuStatusHttpUrl := "v1/managed-organization/{operation_id}"
		getOuStatusPath := client.Endpoint + getOuStatusHttpUrl
		getOuStatusPath = strings.ReplaceAll(getOuStatusPath, "{operation_id}", operationId)

		getOuStatusOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		getOuStatusResp, err := client.Request("GET", getOuStatusPath, &getOuStatusOpt)
		if err != nil {
			return nil, "", err
		}

		getOuStatusRespBody, err := utils.FlattenResponse(getOuStatusResp)
		if err != nil {
			return nil, "", err
		}
		status := utils.PathSearch("status", getOuStatusRespBody, "").(string)
		if status == "FAILED" {
			message := utils.PathSearch("message", getOuStatusRespBody, nil)
			return nil, "FAILED", fmt.Errorf("status: %s; message: %s", status, message)
		}

		return getOuStatusRespBody, status, nil
	}
}
