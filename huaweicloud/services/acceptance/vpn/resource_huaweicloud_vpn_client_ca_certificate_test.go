package vpn

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/vpn"
)

func getVPNClientCACertificate(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME

	product := "vpn"
	client, err := conf.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating VPN client: %s", err)
	}

	return vpn.GetClientCACertificate(client, state.Primary.Attributes["vpn_server_id"], state.Primary.ID)
}

func testClientCACertificateImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.Attributes["vpn_server_id"] == "" {
			return "", fmt.Errorf("attribute (vpn_server_id) of Resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.ID == "" {
			return "", fmt.Errorf("attribute (id) of Resource (%s) not found: %s", name, rs)
		}

		return rs.Primary.Attributes["vpn_server_id"] + "/" + rs.Primary.ID, nil
	}
}

func TestAccClientCACertificate_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_vpn_client_ca_certificate.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getVPNClientCACertificate,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckVPNP2cServer(t)
			acceptance.TestAccPreCheckVPNP2cClientCACertificate(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccClientCACertificate_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "content", acceptance.HW_VPN_P2C_CLIENT_CA_CERTIFICATE),
					resource.TestCheckResourceAttrSet(rName, "issuer"),
					resource.TestCheckResourceAttrSet(rName, "subject"),
					resource.TestCheckResourceAttrSet(rName, "serial_number"),
					resource.TestCheckResourceAttrSet(rName, "expiration_time"),
					resource.TestCheckResourceAttrSet(rName, "signature_algorithm"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				Config: testAccClientCACertificate_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"-update"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testClientCACertificateImportState(rName),
				ImportStateVerifyIgnore: []string{
					"content",
				},
			},
		},
	})
}

func testAccClientCACertificate_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpn_client_ca_certificate" "test" {
  vpn_server_id = "%[1]s"
  name          = "%[2]s"
  content       = "%[3]s"
}
`, acceptance.HW_VPN_P2C_SERVER, name, acceptance.HW_VPN_P2C_CLIENT_CA_CERTIFICATE)
}

func testAccClientCACertificate_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpn_client_ca_certificate" "test" {
  vpn_server_id = "%[1]s"
  name          = "%[2]s-update"
  content       = "%[3]s"
}
`, acceptance.HW_VPN_P2C_SERVER, name, acceptance.HW_VPN_P2C_CLIENT_CA_CERTIFICATE)
}
