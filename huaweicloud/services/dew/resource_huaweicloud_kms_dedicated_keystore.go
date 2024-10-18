// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product DEW
// ---------------------------------------------------------------

package dew

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

// @API DEW POST /v1.0/{project_id}/keystores
// @API DEW GET /v1.0/{project_id}/keystores/{keystore_id}
// @API DEW DELETE /v1.0/{project_id}/keystores/{keystore_id}
func ResourceKmsDedicatedKeystore() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKmsDedicatedKeystoreCreate,
		ReadContext:   resourceKmsDedicatedKeystoreRead,
		DeleteContext: resourceKmsDedicatedKeystoreDelete,
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
			"alias": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the alias of a dedicated keystore.`,
			},
			"hsm_cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of a dedicated HSM cluster that has no dedicated keystore.`,
			},
			"hsm_ca_cert": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the CA certificate of the dedicated HSM cluster.`,
			},
		},
	}
}

func resourceKmsDedicatedKeystoreCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1.0/{project_id}/keystores"
		product = "kms"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating KMS client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
		JSONBody:         buildCreateDedicatedKeystoreBodyParams(d),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating KMS dedicated keystore: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	keystoreId := utils.PathSearch("keystore.keystore_id", createRespBody, "").(string)
	if keystoreId == "" {
		return diag.Errorf("unable to find the KMS dedicated keystore ID from the API response")
	}
	d.SetId(keystoreId)

	return resourceKmsDedicatedKeystoreRead(ctx, d, meta)
}

func buildCreateDedicatedKeystoreBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"keystore_alias": d.Get("alias"),
		"hsm_cluster_id": d.Get("hsm_cluster_id"),
		"hsm_ca_cert":    d.Get("hsm_ca_cert"),
	}
}

func resourceKmsDedicatedKeystoreRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		mErr    *multierror.Error
		httpUrl = "v1.0/{project_id}/keystores/{keystore_id}"
		product = "kms"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating KMS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{keystore_id}", d.Id())
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error.error_code", "KMS.8003"),
			"error retrieving KMS dedicated keystore")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("alias", utils.PathSearch("keystore.keystore_alias", getRespBody, nil)),
		d.Set("hsm_cluster_id", utils.PathSearch("keystore.hsm_cluster_id", getRespBody, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceKmsDedicatedKeystoreDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1.0/{project_id}/keystores/{keystore_id}"
		product = "kms"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating KMS client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{keystore_id}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error deleting KMS dedicated keystore: %s", err)
	}
	return nil
}
