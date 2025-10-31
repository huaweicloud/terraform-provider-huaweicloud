package identitycenter

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getIdentityCenterTenantResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME

	var (
		listHttpUrl = "v1/identity-stores/{identity_store_id}/provision-tenant"
		listProduct = "identitystore"
	)
	client, err := cfg.NewServiceClient(listProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating Identity Center Client: %s", err)
	}

	listPath := client.Endpoint + listHttpUrl
	listPath = strings.ReplaceAll(listPath, "{identity_store_id}", state.Primary.Attributes["identity_store_id"])

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	tenantId := state.Primary.ID

	listResp, err := client.Request("GET",
		listPath, &listOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Identity Center tenant: %s", err)
	}

	listRespBody, err := utils.FlattenResponse(listResp)
	if err != nil {
		return nil, err
	}

	tenant := utils.PathSearch(fmt.Sprintf("provisioning_tenants|[?tenant_id =='%s']|[0]", tenantId), listRespBody, nil)
	if tenant == nil {
		return nil, fmt.Errorf("error get Identity Center tenant")
	}
	return tenant, nil
}

func TestAccIdentityCenterTenant_basic(t *testing.T) {
	var obj interface{}

	rName := "huaweicloud_identitycenter_tenant.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getIdentityCenterTenantResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testIdentityCenterTenant_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "identity_store_id",
						"data.huaweicloud_identitycenter_instance.test", "identity_store_id"),
					resource.TestCheckResourceAttrSet(rName, "creation_time"),
					resource.TestCheckResourceAttrSet(rName, "scim_endpoint"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testIdentityCenterTenantImportState(rName),
			},
		},
	})
}

func testIdentityCenterTenant_basic() string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_identitycenter_tenant" "test" {
  depends_on        = [huaweicloud_identitycenter_identity_provider.test]
  identity_store_id = data.huaweicloud_identitycenter_instance.test.identity_store_id
}
`, testIdentityCenterIdentityProvider_create_without_metadata_basic())
}

func testIdentityCenterTenantImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}

		identityStoreId := rs.Primary.Attributes["identity_store_id"]
		if identityStoreId == "" {
			return "", fmt.Errorf("attribute (identity_store_id) of Resource (%s) not found: %s", name, rs)
		}
		return identityStoreId + "/" + rs.Primary.ID, nil
	}
}
