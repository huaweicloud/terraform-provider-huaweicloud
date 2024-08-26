package dew

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DEW PUT /v1/{project_id}/secrets/{secret_name}/stages/{stage_name}
// @API DEW GET /v1/{project_id}/secrets/{secret_name}/stages/{stage_name}
// @API DEW DELETE /v1/{project_id}/secrets/{secret_name}/stages/{stage_name}
func ResourceSecretVersionState() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSecretVersionStateCreate,
		ReadContext:   resourceSecretVersionStateRead,
		UpdateContext: resourceSecretVersionStateUpdate,
		DeleteContext: resourceSecretVersionStateDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceSecretVersionStateImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"secret_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"version_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSecretVersionStateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                       = meta.(*config.Config)
		region                    = cfg.GetRegion(d)
		createVersionStateHttpUrl = "v1/{project_id}/secrets/{secret_name}/stages/{stage_name}"
		product                   = "kms"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating KMS client: %s", err)
	}

	createVersionStatePath := client.Endpoint + createVersionStateHttpUrl
	createVersionStatePath = strings.ReplaceAll(createVersionStatePath, "{project_id}", client.ProjectID)
	createVersionStatePath = strings.ReplaceAll(createVersionStatePath, "{secret_name}", d.Get("secret_name").(string))
	createVersionStatePath = strings.ReplaceAll(createVersionStatePath, "{stage_name}", d.Get("name").(string))

	createVersionStateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildSecretVersionStateBodyParams(d),
	}

	creatVersionStateResp, err := client.Request("PUT", createVersionStatePath, &createVersionStateOpt)
	if err != nil {
		return diag.Errorf("error creating secret version state: %s", err)
	}

	createVersionStateRespBody, err := utils.FlattenResponse(creatVersionStateResp)
	if err != nil {
		return diag.FromErr(err)
	}

	stateName := utils.PathSearch("stage.name", createVersionStateRespBody, "").(string)
	if stateName == "" {
		return diag.Errorf("error creating secret version state: the name is not found in API response")
	}
	d.SetId(stateName)

	return resourceSecretVersionStateRead(ctx, d, meta)
}

func buildSecretVersionStateBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"version_id": d.Get("version_id"),
	}
	return bodyParams
}

func getVersionStateInfo(client *golangsdk.ServiceClient, d *schema.ResourceData) (*http.Response, error) {
	getVersionStatehttpUrl := "v1/{project_id}/secrets/{secret_name}/stages/{stage_name}"
	getVersionStatePath := client.Endpoint + getVersionStatehttpUrl
	getVersionStatePath = strings.ReplaceAll(getVersionStatePath, "{project_id}", client.ProjectID)
	getVersionStatePath = strings.ReplaceAll(getVersionStatePath, "{secret_name}", d.Get("secret_name").(string))
	getVersionStatePath = strings.ReplaceAll(getVersionStatePath, "{stage_name}", d.Id())
	getVersionStateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	return client.Request("GET", getVersionStatePath, &getVersionStateOpt)
}

func resourceSecretVersionStateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		mErr    *multierror.Error
		product = "kms"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating KMS client: %s", err)
	}

	getVersionStateResp, err := getVersionStateInfo(client, d)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving secret version state")
	}

	getVersionStateRespBody, err := utils.FlattenResponse(getVersionStateResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("name", utils.PathSearch("stage.name", getVersionStateRespBody, nil)),
		d.Set("secret_name", utils.PathSearch("stage.secret_name", getVersionStateRespBody, nil)),
		d.Set("version_id", utils.PathSearch("stage.version_id", getVersionStateRespBody, nil)),
		d.Set("updated_at", utils.FormatTimeStampRFC3339(
			int64(utils.PathSearch("stage.update_time", getVersionStateRespBody, float64(0)).(float64))/1000, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceSecretVersionStateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                       = meta.(*config.Config)
		region                    = cfg.GetRegion(d)
		updateVersionStateHttpUrl = "v1/{project_id}/secrets/{secret_name}/stages/{stage_name}"
		product                   = "kms"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating KMS client: %s", err)
	}

	updateVersionStatePath := client.Endpoint + updateVersionStateHttpUrl
	updateVersionStatePath = strings.ReplaceAll(updateVersionStatePath, "{project_id}", client.ProjectID)
	updateVersionStatePath = strings.ReplaceAll(updateVersionStatePath, "{secret_name}", d.Get("secret_name").(string))
	updateVersionStatePath = strings.ReplaceAll(updateVersionStatePath, "{stage_name}", d.Id())
	updateVersionStateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildSecretVersionStateBodyParams(d),
	}

	_, err = client.Request("PUT", updateVersionStatePath, &updateVersionStateOpt)
	if err != nil {
		return diag.Errorf("error updating secret version state: %s", err)
	}

	return resourceSecretVersionStateRead(ctx, d, meta)
}

func resourceSecretVersionStateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                       = meta.(*config.Config)
		region                    = cfg.GetRegion(d)
		deleteVersionStateHttpUrl = "v1/{project_id}/secrets/{secret_name}/stages/{stage_name}"
		product                   = "kms"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating KMS client: %s", err)
	}

	deleteVersionStatePath := client.Endpoint + deleteVersionStateHttpUrl
	deleteVersionStatePath = strings.ReplaceAll(deleteVersionStatePath, "{project_id}", client.ProjectID)
	deleteVersionStatePath = strings.ReplaceAll(deleteVersionStatePath, "{secret_name}", d.Get("secret_name").(string))
	deleteVersionStatePath = strings.ReplaceAll(deleteVersionStatePath, "{stage_name}", d.Id())
	deleteVersionStateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}

	// Before deleting, call the query details API, if query no result , then process `CheckDeleted` logic.
	_, err = getVersionStateInfo(client, d)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving secret version state")
	}

	_, err = client.Request("DELETE", deleteVersionStatePath, &deleteVersionStateOpt)
	if err != nil {
		return diag.Errorf("error deleting secret version state: %s", err)
	}

	// Successful deletion API call does not guarantee that the resource is successfully deleted.
	// Call the details API to confirm that the resource has been successfully deleted.
	_, err = getVersionStateInfo(client, d)
	if err == nil {
		return diag.Errorf("error deleting secret version state: the version state still exists")
	}

	return nil
}

func resourceSecretVersionStateImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format for import ID, want '<secret_name>/<id>', but got '%s'", d.Id())
	}

	d.SetId(parts[1])
	mErr := multierror.Append(nil,
		d.Set("secret_name", parts[0]),
	)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
