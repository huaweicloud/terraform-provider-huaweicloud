// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product CMDB
// ---------------------------------------------------------------

package cmdb

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const (
	AppNameExistsCode string = "AOM.30004012"
	AppNotExistsCode  string = "AOM.30004003"
)

func ResourceAomApplication() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApplicationCreate,
		ReadContext:   resourceApplicationRead,
		UpdateContext: resourceApplicationUpdate,
		DeleteContext: resourceApplicationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"register_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "schema: Computed",
			},

			// attributes
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creator": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modified_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modifier": {
				Type:     schema.TypeString,
				Computed: true,
			},

			// deprecated
			"aom_id": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "this parameter is deprecated",
			},
			"app_id": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "this parameter is deprecated",
			},
		},
	}
}

func buildCreateApplicationBodyParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":          d.Get("name"),
		"description":   utils.ValueIngoreEmpty(d.Get("description")),
		"display_name":  utils.ValueIngoreEmpty(d.Get("display_name")),
		"register_type": utils.ValueIngoreEmpty(d.Get("register_type")),
		"eps_id":        utils.ValueIngoreEmpty(common.GetEnterpriseProjectID(d, cfg)),
	}
	return bodyParams
}

func resourceApplicationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	product := "cmdb"

	client, err := cfg.NewServiceClient(product, cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	createApplicationHttpUrl := "v1/applications"
	createApplicationPath := client.Endpoint + createApplicationHttpUrl

	createApplicationOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	var resID string
	createApplicationOpt.JSONBody = utils.RemoveNil(buildCreateApplicationBodyParams(d, cfg))
	createApplicationResp, err := client.Request("POST", createApplicationPath, &createApplicationOpt)
	if err != nil {
		// if the application already exists, we reuse it rather than throwing the error
		// this looks weird, but for compatibility with previous behavior
		if !hasErrorCode(err, AppNameExistsCode) {
			return diag.Errorf("error creating AOM application: %s", err)
		}

		var listErr error
		resID, listErr = getApplicationByName(d, client)
		if listErr != nil || resID == "" {
			log.Printf("[WARN] failed to retrieve AOM application: %s", listErr)
			return diag.Errorf("error creating AOM application: %s", err)
		}
	} else {
		createApplicationRespBody, err := utils.FlattenResponse(createApplicationResp)
		if err != nil {
			return diag.FromErr(err)
		}

		id, err := jmespath.Search("id", createApplicationRespBody)
		if err != nil {
			return diag.Errorf("error creating AOM application: ID is not found in API response")
		}
		resID = id.(string)
	}

	d.SetId(resID)
	return resourceApplicationRead(ctx, d, meta)
}

func resourceApplicationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	product := "cmdb"

	client, err := cfg.NewServiceClient(product, cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	getApplicationHttpUrl := "v1/applications/{id}"
	getApplicationPath := client.Endpoint + getApplicationHttpUrl
	getApplicationPath = strings.ReplaceAll(getApplicationPath, "{id}", d.Id())

	getApplicationOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getApplicationResp, err := client.Request("GET", getApplicationPath, &getApplicationOpt)
	if err != nil {
		if hasErrorCode(err, AppNotExistsCode) {
			err = golangsdk.ErrDefault404{}
		}
		return common.CheckDeletedDiag(d, err, "error retrieving Application")
	}

	getApplicationRespBody, err := utils.FlattenResponse(getApplicationResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(nil,
		d.Set("name", utils.PathSearch("name", getApplicationRespBody, nil)),
		d.Set("description", utils.PathSearch("description", getApplicationRespBody, nil)),
		d.Set("display_name", utils.PathSearch("display_name", getApplicationRespBody, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("eps_id", getApplicationRespBody, nil)),
		d.Set("register_type", utils.PathSearch("register_type", getApplicationRespBody, nil)),
		d.Set("create_time", utils.PathSearch("create_time", getApplicationRespBody, nil)),
		d.Set("creator", utils.PathSearch("creator", getApplicationRespBody, nil)),
		d.Set("modified_time", utils.PathSearch("modified_time", getApplicationRespBody, nil)),
		d.Set("modifier", utils.PathSearch("modifier", getApplicationRespBody, nil)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting Application fields: %s", err)
	}

	return nil
}

func buildUpdateApplicationBodyParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":          d.Get("name"),
		"description":   utils.ValueIngoreEmpty(d.Get("description")),
		"display_name":  utils.ValueIngoreEmpty(d.Get("display_name")),
		"register_type": utils.ValueIngoreEmpty(d.Get("register_type")),
		"eps_id":        utils.ValueIngoreEmpty(common.GetEnterpriseProjectID(d, cfg)),
	}
	return bodyParams
}

func resourceApplicationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	product := "cmdb"

	client, err := cfg.NewServiceClient(product, cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	updateApplicationHttpUrl := "v1/applications/{id}"
	updateApplicationPath := client.Endpoint + updateApplicationHttpUrl
	updateApplicationPath = strings.ReplaceAll(updateApplicationPath, "{id}", d.Id())

	updateApplicationOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	updateApplicationOpt.JSONBody = utils.RemoveNil(buildUpdateApplicationBodyParams(d, cfg))
	_, err = client.Request("PUT", updateApplicationPath, &updateApplicationOpt)
	if err != nil {
		return diag.Errorf("error updating Application: %s", err)
	}

	return resourceApplicationRead(ctx, d, meta)
}

func resourceApplicationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	product := "cmdb"

	client, err := cfg.NewServiceClient(product, cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	deleteApplicationHttpUrl := "v1/applications/{id}"
	deleteApplicationPath := client.Endpoint + deleteApplicationHttpUrl
	deleteApplicationPath = strings.ReplaceAll(deleteApplicationPath, "{id}", d.Id())

	deleteApplicationOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	_, err = client.Request("DELETE", deleteApplicationPath, &deleteApplicationOpt)
	if err != nil {
		return diag.Errorf("error deleting Application: %s", err)
	}

	return nil
}

func getApplicationByName(d *schema.ResourceData, client *golangsdk.ServiceClient) (string, error) {
	getApplicationByNameHttpUrl := "v1/applications"
	getApplicationByNamePath := client.Endpoint + getApplicationByNameHttpUrl

	getApplicationByNamequeryParams := buildGetApplicationByNameQueryParams(d)
	getApplicationByNamePath += getApplicationByNamequeryParams

	getApplicationByNameOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	name := d.Get("name").(string)
	respRaw, err := client.Request("GET", getApplicationByNamePath, &getApplicationByNameOpt)
	if err != nil {
		return "", fmt.Errorf("error retrieving AOM application %s: %s", name, err)
	}

	respBody, err := utils.FlattenResponse(respRaw)
	if err != nil {
		return "", fmt.Errorf("error parsing AOM application %s: %s", name, err)
	}

	id, err := jmespath.Search("app_id", respBody)
	if err != nil {
		return "", err
	}
	return id.(string), nil
}

func buildGetApplicationByNameQueryParams(d *schema.ResourceData) string {
	res := fmt.Sprintf("?name=%v", d.Get("name"))

	if v, ok := d.GetOk("display_name"); ok {
		res = fmt.Sprintf("%s&display_name=%v", res, v)
	}

	return res
}

func hasErrorCode(err error, expectCode string) bool {
	if errCode, ok := err.(golangsdk.ErrDefault400); ok {
		var response interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &response); jsonErr == nil {
			errorCode, parseErr := jmespath.Search("error_code", response)
			if parseErr != nil {
				log.Printf("[WARN] failed to parse error_code from response body: %s", parseErr)
			}

			if errorCode == expectCode {
				return true
			}
		}
	}

	return false
}
