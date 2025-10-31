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

func getServiceProviderCertificateResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME

	var (
		getHttpUrl = "v1/identity-stores/{identity_store_id}/saml-certificates"
		getProduct = "identitystore"
	)
	client, err := cfg.NewServiceClient(getProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating IdentityCenter client: %s", err)
	}

	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{identity_store_id}", state.Primary.Attributes["identity_store_id"])

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	certificateId := state.Primary.ID

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Identity Center service provider certificate: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	certificate := utils.PathSearch(fmt.Sprintf("[?certificate_id =='%s']|[0]", certificateId), getRespBody, nil)
	if certificate == nil {
		return nil, fmt.Errorf("error get Identity Center service provider certificate")
	}
	return certificate, nil
}

func TestAccServiceProviderCertificate_basic(t *testing.T) {
	var obj interface{}

	rName := "huaweicloud_identitycenter_service_provider_certificate.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getServiceProviderCertificateResourceFunc,
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
				Config: testServiceProviderCertificate_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "identity_store_id",
						"data.huaweicloud_identitycenter_instance.test", "identity_store_id"),
					resource.TestCheckResourceAttrSet(rName, "x509certificate"),
					resource.TestCheckResourceAttrSet(rName, "algorithm"),
					resource.TestCheckResourceAttrSet(rName, "expiry_date"),
					resource.TestCheckResourceAttr(rName, "state", "INACTIVE"),
				),
			},
			{
				Config: testServiceProviderCertificate_update(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(rName, "identity_store_id",
						"data.huaweicloud_identitycenter_instance.test", "identity_store_id"),
					resource.TestCheckResourceAttrSet(rName, "x509certificate"),
					resource.TestCheckResourceAttrSet(rName, "algorithm"),
					resource.TestCheckResourceAttrSet(rName, "expiry_date"),
					resource.TestCheckResourceAttr(rName, "state", "ACTIVE"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateIdFunc: testServiceProviderCertificateImportState(rName),
				ImportStateVerify: true,
			},
		},
	})
}

func testServiceProviderCertificateImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}

		identityStoreID := rs.Primary.Attributes["identity_store_id"]
		if identityStoreID == "" {
			return "", fmt.Errorf("attribute (identity_store_id) of resource (%s) not found: %s", name, rs)
		}

		return identityStoreID + "/" + rs.Primary.ID, nil
	}
}

func testServiceProviderCertificate_basic() string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_identitycenter_service_provider_certificate" "test"{
  identity_store_id = data.huaweicloud_identitycenter_instance.test.identity_store_id
}
`, testAccDatasourceIdentityCenter_basic())
}

func testServiceProviderCertificate_update() string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_identitycenter_service_provider_certificate" "test"{
  identity_store_id = data.huaweicloud_identitycenter_instance.test.identity_store_id
  state             = "ACTIVE"
}
`, testAccDatasourceIdentityCenter_basic())
}
