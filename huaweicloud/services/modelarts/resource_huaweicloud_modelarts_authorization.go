// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product ModelArts
// ---------------------------------------------------------------

package modelarts

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ModelArts POST /v2/{project_id}/agency
// @API ModelArts DELETE /v2/{project_id}/authorizations
// @API ModelArts GET /v2/{project_id}/authorizations
// @API ModelArts POST /v2/{project_id}/authorizations
func ResourceModelArtsAuthorization() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceModelArtsAuthorizationCreate,
		UpdateContext: resourceModelArtsAuthorizationCreate,
		ReadContext:   resourceModelArtsAuthorizationRead,
		DeleteContext: resourceModelArtsAuthorizationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"user_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `User ID.`,
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Authorization type. The valid value is **agency**.`,
			},
			"agency_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Agency name.`,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^(modelarts_agency)$|^(ma_agency_.+)$`),
					"The agency name can be modelarts_agency or prefixed with ma_agency_"),
			},
			"user_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `User Name.`,
			},
		},
	}
}

func resourceModelArtsAuthorizationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createAgency: create a ModelArts agency.
	var (
		createAgencyHttpUrl = "v2/{project_id}/agency"
		createAgencyProduct = "modelarts"
	)
	createAgencyClient, err := cfg.NewServiceClient(createAgencyProduct, region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	createAgencyPath := createAgencyClient.Endpoint + createAgencyHttpUrl
	createAgencyPath = strings.ReplaceAll(createAgencyPath, "{project_id}", createAgencyClient.ProjectID)

	createAgencyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	createAgencyOpt.JSONBody = utils.RemoveNil(buildCreateAgencyBodyParams(d))
	_, err = createAgencyClient.Request("POST", createAgencyPath, &createAgencyOpt)
	if err != nil && parseAgencyErrorIgnoreConflict(err) != nil {
		return diag.Errorf("error creating ModelArts agency: %s", err)
	}

	// createAuth: create a ModelArts authorization.
	var (
		createAuthHttpUrl = "v2/{project_id}/authorizations"
		createAuthProduct = "modelarts"
	)
	createAuthClient, err := cfg.NewServiceClient(createAuthProduct, region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	createAuthPath := createAuthClient.Endpoint + createAuthHttpUrl
	createAuthPath = strings.ReplaceAll(createAuthPath, "{project_id}", createAuthClient.ProjectID)

	createAuthOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	createAuthOpt.JSONBody = utils.RemoveNil(buildCreateAuthBodyParams(d))
	_, err = createAuthClient.Request("POST", createAuthPath, &createAuthOpt)
	if err != nil {
		return diag.Errorf("error creating ModelArts authorization: %s", err)
	}

	d.SetId(d.Get("user_id").(string))

	return resourceModelArtsAuthorizationRead(ctx, d, meta)
}

func buildCreateAgencyBodyParams(d *schema.ResourceData) map[string]interface{} {
	agencyNameSuffix := "" // when agency_name is modelarts_agency, set empty

	agencyName := d.Get("agency_name").(string)
	// when agency_name is prefixed with 'ma_agency_', take the suffix
	if strings.HasPrefix(agencyName, "ma_agency_") {
		agencyNameSuffix = strings.TrimPrefix(agencyName, "ma_agency_")
	}

	bodyParams := map[string]interface{}{
		"agency_name_suffix": agencyNameSuffix,
	}
	return bodyParams
}

func buildCreateAuthBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"user_id": utils.ValueIgnoreEmpty(d.Get("user_id")),
		"type":    utils.ValueIgnoreEmpty(d.Get("type")),
		"content": utils.ValueIgnoreEmpty(d.Get("agency_name")),
	}
	return bodyParams
}

func resourceModelArtsAuthorizationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getModelAuth: Query the ModelArts authorization.
	var (
		getModelAuthHttpUrl = "v2/{project_id}/authorizations"
		getModelAuthProduct = "modelarts"
	)
	getModelAuthClient, err := cfg.NewServiceClient(getModelAuthProduct, region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	getModelAuthPath := getModelAuthClient.Endpoint + getModelAuthHttpUrl
	getModelAuthPath = strings.ReplaceAll(getModelAuthPath, "{project_id}", getModelAuthClient.ProjectID)

	getModelAuthResp, err := pagination.ListAllItems(
		getModelAuthClient,
		"offset",
		getModelAuthPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving ModelArts authorization")
	}

	getModelAuthRespJson, err := json.Marshal(getModelAuthResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var getModelAuthRespBody interface{}
	err = json.Unmarshal(getModelAuthRespJson, &getModelAuthRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	getModelAuthRespBody = SearchAuthById(getModelAuthRespBody, d.Id())
	if getModelAuthRespBody == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "no data found")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("user_id", utils.PathSearch("user_id", getModelAuthRespBody, nil)),
		d.Set("user_name", utils.PathSearch("user_name", getModelAuthRespBody, nil)),
		d.Set("type", utils.PathSearch("type", getModelAuthRespBody, nil)),
		d.Set("agency_name", utils.PathSearch("content", getModelAuthRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

// There are two arrays in the raw response body, but in the return of pagination.ListAllItems,
// the key of array may be 'auth' or 'content', so need to search twice at most.
func SearchAuthById(getModelAuthRespBody interface{}, id string) interface{} {
	jsonPath1 := fmt.Sprintf("auth[?user_id =='%s']|[0]", id)
	rst := utils.PathSearch(jsonPath1, getModelAuthRespBody, nil)
	if rst != nil {
		return rst
	}

	jsonPath2 := fmt.Sprintf("content[?user_id =='%s']|[0]", id)
	return utils.PathSearch(jsonPath2, getModelAuthRespBody, nil)
}

func resourceModelArtsAuthorizationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteModelAuth: delete ModelArts authorization
	var (
		deleteModelAuthHttpUrl = "v2/{project_id}/authorizations"
		deleteModelAuthProduct = "modelarts"
	)
	deleteModelAuthClient, err := cfg.NewServiceClient(deleteModelAuthProduct, region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	deleteModelAuthPath := deleteModelAuthClient.Endpoint + deleteModelAuthHttpUrl
	deleteModelAuthPath = strings.ReplaceAll(deleteModelAuthPath, "{project_id}", deleteModelAuthClient.ProjectID)

	deleteModelAuthqueryParams := fmt.Sprintf("?user_id=%v", d.Id())
	deleteModelAuthPath += deleteModelAuthqueryParams

	deleteModelAuthOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	_, err = deleteModelAuthClient.Request("DELETE", deleteModelAuthPath, &deleteModelAuthOpt)
	if err != nil {
		return diag.Errorf("error deleting ModelArts authorization: %s", err)
	}

	return nil
}

func parseAgencyErrorIgnoreConflict(respErr error) error {
	var apiErr interface{}
	if errCode, ok := respErr.(golangsdk.ErrDefault400); ok {
		pErr := json.Unmarshal(errCode.Body, &apiErr)
		if pErr != nil {
			return pErr
		}
		errCode, err := jmespath.Search(`error_code`, apiErr)
		if err != nil {
			return fmt.Errorf("error parse errorCode from response body: %s", err.Error())
		}

		if errCode == `ModelArts.3818` { // ignore the conflict error.
			return nil
		}
	}
	return respErr
}
