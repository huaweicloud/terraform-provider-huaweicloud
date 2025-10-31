package identitycenter

import (
	"context"
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

// @API IdentityCenter POST /v1/service/start
// @API IdentityCenter POST /v1/instances/{instance_id}/alias
// @API IdentityCenter GET  /v1/instances
// @API IdentityCenter GET  /v1/identity-center-service/status
// @API IdentityCenter POST /v1/service/delete
func ResourceIdentityCenterInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityCenterInstanceCreate,
		UpdateContext: resourceIdentityCenterInstanceUpdate,
		ReadContext:   resourceIdentityCenterInstanceRead,
		DeleteContext: resourceIdentityCenterInstanceDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"alias": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"identity_store_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_urn": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIdentityCenterInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		startIdentityCenterHttpUrl            = "v1/service/start"
		getIdentityCenterServiceStatusHttpUrl = "v1/identity-center-service/status"
		createAliasHttpUrl                    = "v1/instances/{instance_id}/alias"
		listInstancesHttpUrl                  = "v1/instances"
		identityCenterProduct                 = "identitycenter"
	)
	client, err := cfg.NewServiceClient(identityCenterProduct, region)
	if err != nil {
		return diag.Errorf("error creating IdentityCenter client: %s", err)
	}

	startIdentityCenterPath := client.Endpoint + startIdentityCenterHttpUrl

	startIdentityCenterOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = client.Request("POST", startIdentityCenterPath, &startIdentityCenterOpt)
	if err != nil {
		return diag.Errorf("error creating IdentityCenter instance: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{"disable"},
		Target:  []string{"enable"},
		Refresh: identityCenterServiceStatusRefreshFunc(getIdentityCenterServiceStatusHttpUrl,
			"serviceStatus", client),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        1 * time.Second,
		PollInterval: 5 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for creating IAM Identity Center instance to complete: %s", err)
	}

	listInstancesPath := client.Endpoint + listInstancesHttpUrl

	listInstancesOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	listInstancesResp, err := client.Request("GET", listInstancesPath, &listInstancesOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Identity Center Instance")
	}

	listInstancesRespBody, err := utils.FlattenResponse(listInstancesResp)
	if err != nil {
		return diag.FromErr(err)
	}

	instanceId := utils.PathSearch("instances|[0].instance_id", listInstancesRespBody, "").(string)
	if instanceId == "" {
		return diag.Errorf("unable to find the Identity Center instance ID from the API response")
	}
	d.SetId(instanceId)

	if _, ok := d.GetOk("alias"); ok {
		createAliasPath := client.Endpoint + createAliasHttpUrl
		createAliasPath = strings.ReplaceAll(createAliasPath, "{instance_id}", instanceId)
		createAliasOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         utils.RemoveNil(buildCreateAliasBodyParams(d)),
		}
		_, err = client.Request("POST", createAliasPath, &createAliasOpt)
		if err != nil {
			return diag.Errorf("error creating IdentityCenter instance alias: %s", err)
		}
	}

	return resourceIdentityCenterInstanceRead(ctx, d, meta)
}

func buildCreateAliasBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"alias": utils.ValueIgnoreEmpty(d.Get("alias")),
	}
	return bodyParams
}

func identityCenterServiceStatusRefreshFunc(getRequestStatusHttpUrl, searchExpression string,
	client *golangsdk.ServiceClient) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getPath := client.Endpoint + getRequestStatusHttpUrl

		getOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		getResp, err := client.Request("GET", getPath, &getOpt)
		if err != nil {
			return nil, "", err
		}
		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return nil, "", err
		}

		status := utils.PathSearch(searchExpression, getRespBody, "").(string)
		if status == "enable" {
			return getRespBody, status, nil
		}
		return getRespBody, "disable", nil
	}
}

func resourceIdentityCenterInstanceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		listHttpUrl = "v1/instances"
		listProduct = "identitycenter"
	)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient(listProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center Client: %s", err)
	}

	listPath := client.Endpoint + listHttpUrl

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	listResp, err := client.Request("GET", listPath, &listOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Identity Center Instance")
	}

	listRespBody, err := utils.FlattenResponse(listResp)
	if err != nil {
		return diag.FromErr(err)
	}

	instance := utils.PathSearch(fmt.Sprintf("instances[?instance_id =='%s']|[0]", d.Id()), listRespBody, nil)

	if instance == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "no instance found.")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("identity_store_id", utils.PathSearch("identity_store_id", instance, nil)),
		d.Set("instance_urn", utils.PathSearch("instance_urn", instance, nil)),
		d.Set("alias", utils.PathSearch("alias", instance, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceIdentityCenterInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v1/instances/{instance_id}/alias"
		product = "identitycenter"
	)

	if d.HasChange("alias") {
		client, err := cfg.NewServiceClient(product, region)
		if err != nil {
			return diag.Errorf("error creating Identity Center Client: %s", err)
		}

		updateAliasPath := client.Endpoint + httpUrl
		updateAliasPath = strings.ReplaceAll(updateAliasPath, "{instance_id}", d.Id())
		updateAliasOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         utils.RemoveNil(buildCreateAliasBodyParams(d)),
		}

		_, err = client.Request("POST", updateAliasPath, &updateAliasOpt)

		if err != nil {
			return diag.Errorf("error creating IdentityCenter instance alias: %s", err)
		}
	}

	return resourceIdentityCenterInstanceRead(ctx, d, meta)
}

func resourceIdentityCenterInstanceDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		deleteHttpUrl = "v1/service/delete"
		deleteProduct = "identitycenter"
	)
	client, err := cfg.NewServiceClient(deleteProduct, region)
	if err != nil {
		return diag.Errorf("error creating IdentityCenter client: %s", err)
	}

	deletePath := client.Endpoint + deleteHttpUrl

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = client.Request("POST", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting IdentityCenter instance")
	}

	return nil
}
