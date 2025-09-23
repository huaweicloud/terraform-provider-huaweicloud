// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product CC
// ---------------------------------------------------------------

package cc

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CC POST /v3/{domain_id}/ccaas/authorisations
// @API CC GET /v3/{domain_id}/ccaas/authorisations
// @API CC DELETE /v3/{domain_id}/ccaas/authorisations/{id}
// @API CC PUT /v3/{domain_id}/ccaas/authorisations/{id}
func ResourceAuthorization() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAuthorizationCreate,
		UpdateContext: resourceAuthorizationUpdate,
		ReadContext:   resourceAuthorizationRead,
		DeleteContext: resourceAuthorizationDelete,
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
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The name of the cross-account authorization.`,
			},
			"instance_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The instance type.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The instance ID.`,
			},
			"cloud_connection_domain_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The peer account ID that you want to authorize.`,
			},
			"cloud_connection_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Peer cloud connection ID.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The description of the cross-account authorization.`,
			},
		},
	}
}

func resourceAuthorizationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		createAuthorizationHttpUrl = "v3/{domain_id}/ccaas/authorisations"
		createAuthorizationProduct = "cc"
	)
	createAuthorizationClient, err := cfg.NewServiceClient(createAuthorizationProduct, region)
	if err != nil {
		return diag.Errorf("error creating CC client: %s", err)
	}

	createAuthorizationPath := createAuthorizationClient.Endpoint + createAuthorizationHttpUrl
	createAuthorizationPath = strings.ReplaceAll(createAuthorizationPath, "{domain_id}", cfg.DomainID)

	createAuthorizationOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			201,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	createAuthorizationOpt.JSONBody = utils.RemoveNil(buildCreateAuthorizationBodyParams(d, cfg))
	createAuthorizationResp, err := createAuthorizationClient.Request("POST", createAuthorizationPath, &createAuthorizationOpt)
	if err != nil {
		return diag.Errorf("error creating authorization: %s", err)
	}

	createAuthorizationRespBody, err := utils.FlattenResponse(createAuthorizationResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("authorisation.id", createAuthorizationRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating authorization: ID is not found in API response")
	}
	d.SetId(id)

	return resourceAuthorizationRead(ctx, d, meta)
}

func buildCreateAuthorizationBodyParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	region := cfg.GetRegion(d)
	bodyParams := map[string]interface{}{
		"authorisation": map[string]interface{}{
			"name":                       utils.ValueIgnoreEmpty(d.Get("name")),
			"description":                utils.ValueIgnoreEmpty(d.Get("description")),
			"instance_type":              d.Get("instance_type"),
			"instance_id":                d.Get("instance_id"),
			"project_id":                 cfg.GetProjectID(region),
			"region_id":                  region,
			"cloud_connection_domain_id": d.Get("cloud_connection_domain_id"),
			"cloud_connection_id":        d.Get("cloud_connection_id"),
		},
	}
	return bodyParams
}

func resourceAuthorizationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getAuthorization: Query the cross-account authorization
	var (
		getAuthorizationHttpUrl = "v3/{domain_id}/ccaas/authorisations?id={id}"
		getAuthorizationProduct = "cc"
	)
	getAuthorizationClient, err := cfg.NewServiceClient(getAuthorizationProduct, region)
	if err != nil {
		return diag.Errorf("error creating CC client: %s", err)
	}

	getAuthorizationPath := getAuthorizationClient.Endpoint + getAuthorizationHttpUrl
	getAuthorizationPath = strings.ReplaceAll(getAuthorizationPath, "{domain_id}", cfg.DomainID)
	getAuthorizationPath = strings.ReplaceAll(getAuthorizationPath, "{id}", d.Id())

	getAuthorizationOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	getAuthorizationResp, err := getAuthorizationClient.Request("GET", getAuthorizationPath, &getAuthorizationOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving authorization")
	}

	getAuthorizationRespBody, err := utils.FlattenResponse(getAuthorizationResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonPath := fmt.Sprintf("authorisations[?id =='%s']|[0]", d.Id())
	getAuthorizationRespBody = utils.PathSearch(jsonPath, getAuthorizationRespBody, nil)
	if getAuthorizationRespBody == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "no data found")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", getAuthorizationRespBody, nil)),
		d.Set("description", utils.PathSearch("description", getAuthorizationRespBody, nil)),
		d.Set("instance_type", utils.PathSearch("instance_type", getAuthorizationRespBody, nil)),
		d.Set("instance_id", utils.PathSearch("instance_id", getAuthorizationRespBody, nil)),
		d.Set("cloud_connection_domain_id", utils.PathSearch("cloud_connection_domain_id", getAuthorizationRespBody, nil)),
		d.Set("cloud_connection_id", utils.PathSearch("cloud_connection_id", getAuthorizationRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceAuthorizationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateAuthorizationChanges := []string{
		"name",
		"description",
	}

	if d.HasChanges(updateAuthorizationChanges...) {
		var (
			updateAuthorizationHttpUrl = "v3/{domain_id}/ccaas/authorisations/{id}"
			updateAuthorizationProduct = "cc"
		)
		updateAuthorizationClient, err := cfg.NewServiceClient(updateAuthorizationProduct, region)
		if err != nil {
			return diag.Errorf("error creating CC client: %s", err)
		}

		updateAuthorizationPath := updateAuthorizationClient.Endpoint + updateAuthorizationHttpUrl
		updateAuthorizationPath = strings.ReplaceAll(updateAuthorizationPath, "{domain_id}", cfg.DomainID)
		updateAuthorizationPath = strings.ReplaceAll(updateAuthorizationPath, "{id}", d.Id())

		updateAuthorizationOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
			MoreHeaders: map[string]string{"Content-Type": "application/json"},
		}

		updateAuthorizationOpt.JSONBody = utils.RemoveNil(buildUpdateAuthorizationBodyParams(d))
		_, err = updateAuthorizationClient.Request("PUT", updateAuthorizationPath, &updateAuthorizationOpt)
		if err != nil {
			return diag.Errorf("error updating authorization: %s", err)
		}
	}
	return resourceAuthorizationRead(ctx, d, meta)
}

func buildUpdateAuthorizationBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"authorisation": map[string]interface{}{
			"name":        utils.ValueIgnoreEmpty(d.Get("name")),
			"description": utils.ValueIgnoreEmpty(d.Get("description")),
		},
	}
	return bodyParams
}

func resourceAuthorizationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		deleteAuthorizationHttpUrl = "v3/{domain_id}/ccaas/authorisations/{id}"
		deleteAuthorizationProduct = "cc"
	)
	deleteAuthorizationClient, err := cfg.NewServiceClient(deleteAuthorizationProduct, region)
	if err != nil {
		return diag.Errorf("error creating CC client: %s", err)
	}

	deleteAuthorizationPath := deleteAuthorizationClient.Endpoint + deleteAuthorizationHttpUrl
	deleteAuthorizationPath = strings.ReplaceAll(deleteAuthorizationPath, "{domain_id}", cfg.DomainID)
	deleteAuthorizationPath = strings.ReplaceAll(deleteAuthorizationPath, "{id}", d.Id())

	deleteAuthorizationOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	_, err = deleteAuthorizationClient.Request("DELETE", deleteAuthorizationPath, &deleteAuthorizationOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting authorization")
	}

	return nil
}
