package elb

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

// @API ELB POST /v3/{project_id}/elb/certificates/settings/private-key-echo
// @API ELB GET /v3/{project_id}/elb/certificates/settings/private-key-echo
func ResourceCertificatePrivateKeyEcho() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCertificatePrivateKeyEchoCreate,
		ReadContext:   resourceCertificatePrivateKeyEchoRead,
		DeleteContext: resourceCertificatePrivateKeyEchoDelete,
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
		},
	}
}

func resourceCertificatePrivateKeyEchoCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	var (
		product = "elb"
	)
	client, err := conf.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}
	err = updateCertificatePrivateKeyEcho(client, true, "creating")
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(client.ProjectID)

	return resourceCertificatePrivateKeyEchoRead(ctx, d, meta)
}

func resourceCertificatePrivateKeyEchoRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v3/{project_id}/elb/certificates/settings/private-key-echo"
		product = "elb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving ELB certificate private key echo")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	privateKeyEcho := utils.PathSearch("private_key_echo", getRespBody, false).(bool)
	if !privateKeyEcho {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceCertificatePrivateKeyEchoDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		product = "elb"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}
	err = updateCertificatePrivateKeyEcho(client, false, "deleting")
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func updateCertificatePrivateKeyEcho(client *golangsdk.ServiceClient, privateKeyEcho bool, operation string) error {
	var (
		httpUrl = "v3/{project_id}/elb/certificates/settings/private-key-echo"
	)
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateOpt.JSONBody = utils.RemoveNil(buildUpdateCertificatePrivateKeyEchoBodyParams(privateKeyEcho))
	_, err := client.Request("POST", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("error %s ELB certificate private key echo: %s", operation, err)
	}

	return nil
}

func buildUpdateCertificatePrivateKeyEchoBodyParams(privateKeyEcho bool) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"private_key_echo": privateKeyEcho,
	}
	return bodyParams
}
