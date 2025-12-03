package hss

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API HSS POST /v5/{project_id}/setting/login-common-location
// @API HSS GET /v5/{project_id}/setting/login-common-location
func ResourceLoginCommonLocation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLoginCommonLocationCreate,
		ReadContext:   resourceLoginCommonLocationRead,
		UpdateContext: resourceLoginCommonLocationUpdate,
		DeleteContext: resourceLoginCommonLocationDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceLoginCommonLocationImportState,
		},

		CustomizeDiff: config.FlexibleForceNew([]string{
			"enterprise_project_id",
			"area_code",
		}),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"area_code": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"host_id_list": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
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

func buildLoginCommonLocationModifyQueryParams(epsId string) string {
	if epsId != "" {
		return fmt.Sprintf("?enterprise_project_id=%s", epsId)
	}

	return ""
}

func updateLoginCommonLocation(client *golangsdk.ServiceClient, d *schema.ResourceData, epsId string) error {
	requestPath := client.Endpoint + "v5/{project_id}/setting/login-common-location"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildLoginCommonLocationModifyQueryParams(epsId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"area_code":    d.Get("area_code"),
			"host_id_list": d.Get("host_id_list"),
		},
	}

	_, err := client.Request("POST", requestPath, &requestOpt)
	return err
}

func resourceLoginCommonLocationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "hss"
		epsId   = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	if err := updateLoginCommonLocation(client, d, epsId); err != nil {
		return diag.Errorf("error updating HSS login common location in create operation: %s", err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	return resourceLoginCommonLocationRead(ctx, d, meta)
}

func buildLoginCommonLocationQueryParams(epsId string, areaCode int) string {
	rst := fmt.Sprintf("?area_code=%d", areaCode)

	if epsId != "" {
		rst += fmt.Sprintf("&enterprise_project_id=%s", epsId)
	}

	return rst
}

func QueryLoginCommonLocation(client *golangsdk.ServiceClient, areaCode int, epsId string) (interface{}, error) {
	requestPath := client.Endpoint + "v5/{project_id}/setting/login-common-location"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildLoginCommonLocationQueryParams(epsId, areaCode)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	data := utils.PathSearch("data_list|[0]", respBody, nil)
	if data == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return data, nil
}

func resourceLoginCommonLocationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		product  = "hss"
		epsId    = cfg.GetEnterpriseProjectID(d)
		areaCode = d.Get("area_code").(int)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	respBody, err := QueryLoginCommonLocation(client, areaCode, epsId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving HSS login common location")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("host_id_list", utils.ExpandToStringList(utils.PathSearch("host_id_list", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceLoginCommonLocationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "hss"
		epsId   = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	if err := updateLoginCommonLocation(client, d, epsId); err != nil {
		return diag.Errorf("error updating HSS login common location in update operation: %s", err)
	}

	return resourceLoginCommonLocationRead(ctx, d, meta)
}

func resourceLoginCommonLocationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "hss"
		epsId   = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + "v5/{project_id}/setting/login-common-location"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildLoginCommonLocationModifyQueryParams(epsId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"area_code":    d.Get("area_code"),
			"host_id_list": make([]string, 0),
		},
	}

	_, err = client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error deleting HSS login common location: %s", err)
	}

	return nil
}

func resourceLoginCommonLocationImportState(_ context.Context, d *schema.ResourceData, _ interface{}) (
	[]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, errors.New("invalid format of import ID, must be <enterprise_project_id>/<area_code>")
	}

	areaCode, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, errors.New("invalid format of import ID, <area_code> must be number value")
	}

	mErr := multierror.Append(
		d.Set("enterprise_project_id", parts[0]),
		d.Set("area_code", areaCode),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
