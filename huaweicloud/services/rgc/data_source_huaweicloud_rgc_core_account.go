package rgc

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RGC GET /v1/managed-organization/managed-core-accounts
func DataSourceCoreAccounts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCoreAccountsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"account_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"core_resource_mappings": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceCoreAccountsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var getCoreAccountsProduct = "rgc"
	getCoreAccountsClient, err := cfg.NewServiceClient(getCoreAccountsProduct, region)
	if err != nil {
		return diag.Errorf("error creating RGC client: %s", err)
	}

	getCoreAccountsRespBody, err := getCoreAccounts(getCoreAccountsClient, d)

	if err != nil {
		return diag.Errorf("error retrieving RGC core account: %s", err)
	}

	coreResourceMappings, err := parseCoreResourceMapping(utils.PathSearch("core_resource_mappings",
		getCoreAccountsRespBody, make(map[string]interface{})).(map[string]interface{}))

	if err != nil {
		return diag.FromErr(err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr = multierror.Append(mErr,
		d.Set("account_id", utils.PathSearch("account_id", getCoreAccountsRespBody, nil)),
		d.Set("core_resource_mappings", coreResourceMappings),
		d.Set("region", region),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func parseCoreResourceMapping(mapping map[string]interface{}) (string, error) {
	if len(mapping) == 0 {
		return "", nil
	}

	jsonBytes, err := json.Marshal(mapping)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

func getCoreAccounts(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	var (
		getCoreAccountsHttpUrl = "v1/managed-organization/managed-core-accounts"
	)
	getCoreAccountsHttpPath := client.Endpoint + getCoreAccountsHttpUrl
	getCoreAccountsHttpPath += buildCoreAccountsQueryParams(d)

	getCoreAccountsHttpOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getCoreAccountsHttpResp, err := client.Request("GET", getCoreAccountsHttpPath, &getCoreAccountsHttpOpt)
	if err != nil {
		return nil, err
	}
	getCoreAccountsRespBody, err := utils.FlattenResponse(getCoreAccountsHttpResp)
	if err != nil {
		return nil, err
	}
	return getCoreAccountsRespBody, nil
}

func buildCoreAccountsQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("account_type"); ok {
		res = fmt.Sprintf("%s&account_type=%v", res, v)
	}

	if res != "" {
		res = "?" + res[1:]
	}

	return res
}
