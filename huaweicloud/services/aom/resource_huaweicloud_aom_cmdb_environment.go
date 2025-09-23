// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product AOM
// ---------------------------------------------------------------

package aom

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

// @API AOM POST /v1/environments
// @API AOM DELETE /v1/environments/{environment_id}
// @API AOM GET /v1/environments/{environment_id}
// @API AOM PUT /v1/environments/{environment_id}
func ResourceCmdbEnvironment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCmdbEnvironmentCreate,
		ReadContext:   resourceCmdbEnvironmentRead,
		UpdateContext: resourceCmdbEnvironmentUpdate,
		DeleteContext: resourceCmdbEnvironmentDelete,
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
			"component_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"os_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// attributes
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"register_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildCreateEnvironmentBodyParams(d *schema.ResourceData, region string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"region":       region,
		"component_id": d.Get("component_id"),
		"env_name":     d.Get("name"),
		"env_type":     d.Get("type"),
		"os_type":      d.Get("os_type"),
		"description":  utils.ValueIgnoreEmpty(d.Get("description")),
	}
	return bodyParams
}

func resourceCmdbEnvironmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	product := "aom"
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	createEnvironmentHttpUrl := "v1/environments"
	createEnvironmentPath := client.Endpoint + createEnvironmentHttpUrl

	createEnvironmentOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	createEnvironmentOpt.JSONBody = utils.RemoveNil(buildCreateEnvironmentBodyParams(d, region))
	createEnvironmentResp, err := client.Request("POST", createEnvironmentPath, &createEnvironmentOpt)
	if err != nil {
		return diag.Errorf("error creating CMDB environment: %s", err)
	}
	createEnvironmentRespBody, err := utils.FlattenResponse(createEnvironmentResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("id", createEnvironmentRespBody, "").(string)
	if id == "" {
		return diag.Errorf("unable to find environment ID from the API response")
	}

	d.SetId(id)
	return resourceCmdbEnvironmentRead(ctx, d, meta)
}

func resourceCmdbEnvironmentRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	product := "aom"
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient(product, cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	getEnvironmentHttpUrl := "v1/environments/{environment_id}"
	getEnvironmentPath := client.Endpoint + getEnvironmentHttpUrl
	getEnvironmentPath = strings.ReplaceAll(getEnvironmentPath, "{environment_id}", d.Id())

	getEnvironmentOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getEnvironmentResp, err := client.Request("GET", getEnvironmentPath, &getEnvironmentOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", EnvNotExistsCode),
			"error retrieving CMDB environment")
	}

	getEnvironmentRespBody, err := utils.FlattenResponse(getEnvironmentResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("component_id", utils.PathSearch("component_id", getEnvironmentRespBody, nil)),
		d.Set("name", utils.PathSearch("env_name", getEnvironmentRespBody, nil)),
		d.Set("type", utils.PathSearch("env_type", getEnvironmentRespBody, nil)),
		d.Set("os_type", utils.PathSearch("os_type", getEnvironmentRespBody, nil)),
		d.Set("register_type", utils.PathSearch("register_type", getEnvironmentRespBody, nil)),
		d.Set("description", utils.PathSearch("description", getEnvironmentRespBody, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("eps_id", getEnvironmentRespBody, nil)),
		d.Set("created_at", utils.PathSearch("create_time", getEnvironmentRespBody, nil)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting CMDB environment fields: %s", err)
	}

	return nil
}

func buildUpdateEnvironmentBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"component_id": d.Get("component_id"),
		"os_type":      d.Get("os_type"),
		"env_name":     d.Get("name"),
		"env_type":     d.Get("type"),
		"description":  d.Get("description"),
	}
	return bodyParams
}

func resourceCmdbEnvironmentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	product := "aom"

	client, err := cfg.NewServiceClient(product, cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	updateEnvironmentHttpUrl := "v1/environments/{environment_id}"
	updateEnvironmentPath := client.Endpoint + updateEnvironmentHttpUrl
	updateEnvironmentPath = strings.ReplaceAll(updateEnvironmentPath, "{environment_id}", d.Id())

	updateEnvironmentOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	updateEnvironmentOpt.JSONBody = utils.RemoveNil(buildUpdateEnvironmentBodyParams(d))
	_, err = client.Request("PUT", updateEnvironmentPath, &updateEnvironmentOpt)
	if err != nil {
		return diag.Errorf("error updating CMDB environment: %s", err)
	}

	return resourceCmdbEnvironmentRead(ctx, d, meta)
}

func resourceCmdbEnvironmentDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	product := "aom"

	client, err := cfg.NewServiceClient(product, cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	deleteEnvironmentHttpUrl := "v1/environments/{environment_id}"
	deleteEnvironmentPath := client.Endpoint + deleteEnvironmentHttpUrl
	deleteEnvironmentPath = strings.ReplaceAll(deleteEnvironmentPath, "{environment_id}", d.Id())

	deleteEnvironmentOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	_, err = client.Request("DELETE", deleteEnvironmentPath, &deleteEnvironmentOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", EnvNotExistsCode),
			"error deleting CMDB environment")
	}

	return nil
}
