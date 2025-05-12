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

func getServerFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("vpn", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating VPN client: %s", err)
	}

	return vpn.GetServer(client, state.Primary.Attributes["p2c_vgw_id"], state.Primary.ID)
}

func testServerImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.Attributes["p2c_vgw_id"] == "" {
			return "", fmt.Errorf("P2C VGW ID of Resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.ID == "" {
			return "", fmt.Errorf("ID of Resource (%s) not found: %s", name, rs)
		}

		return rs.Primary.Attributes["p2c_vgw_id"] + "/" + rs.Primary.ID, nil
	}
}

func TestAccServer_basic(t *testing.T) {
	var obj interface{}

	rName := "huaweicloud_vpn_server.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getServerFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckVPNP2cGatewayId(t)
			acceptance.TestAccPreCheckVPNP2cServerCertificateID(t)
			acceptance.TestAccPreCheckVPNP2cClientCACertificate(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccServer_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "p2c_vgw_id", acceptance.HW_VPN_P2C_GATEWAY_ID),
					resource.TestCheckResourceAttr(rName, "local_subnets.#", "1"),
					resource.TestCheckResourceAttr(rName, "client_cidr", "192.168.1.0/24"),
					resource.TestCheckResourceAttr(rName, "client_auth_type", "CERT"),
					resource.TestCheckResourceAttr(rName, "server_certificate.0.id", acceptance.HW_VPN_P2C_SERVER_CERTIFICATE_ID),
					resource.TestCheckResourceAttrSet(rName, "server_certificate.0.name"),
					resource.TestCheckResourceAttrSet(rName, "server_certificate.0.serial_number"),
					resource.TestCheckResourceAttrSet(rName, "server_certificate.0.expiration_time"),
					resource.TestCheckResourceAttrSet(rName, "server_certificate.0.signature_algorithm"),
					resource.TestCheckResourceAttrSet(rName, "server_certificate.0.issuer"),
					resource.TestCheckResourceAttrSet(rName, "server_certificate.0.subject"),
					resource.TestCheckResourceAttr(rName, "ssl_options.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(rName, "ssl_options.0.port", "443"),
					resource.TestCheckResourceAttr(rName, "ssl_options.0.encryption_algorithm", "AES-128-GCM"),
					resource.TestCheckResourceAttr(rName, "ssl_options.0.is_compressed", "false"),
					resource.TestCheckResourceAttrSet(rName, "ssl_options.0.authentication_algorithm"),
					resource.TestCheckResourceAttr(rName, "client_ca_certificates_uploaded.0.name", "test-cert"),
					resource.TestCheckResourceAttrSet(rName, "client_ca_certificates_uploaded.0.id"),
					resource.TestCheckResourceAttrSet(rName, "client_ca_certificates_uploaded.0.issuer"),
					resource.TestCheckResourceAttrSet(rName, "client_ca_certificates_uploaded.0.signature_algorithm"),
					resource.TestCheckResourceAttrSet(rName, "client_ca_certificates_uploaded.0.subject"),
					resource.TestCheckResourceAttrSet(rName, "client_ca_certificates_uploaded.0.serial_number"),
					resource.TestCheckResourceAttrSet(rName, "client_ca_certificates_uploaded.0.expiration_time"),
					resource.TestCheckResourceAttrSet(rName, "client_ca_certificates_uploaded.0.created_at"),
					resource.TestCheckResourceAttrSet(rName, "client_ca_certificates_uploaded.0.updated_at"),
					resource.TestCheckResourceAttrSet(rName, "tunnel_protocol"),
					resource.TestCheckResourceAttrSet(rName, "client_config"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				Config: testAccServer_update(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "p2c_vgw_id", acceptance.HW_VPN_P2C_GATEWAY_ID),
					resource.TestCheckResourceAttr(rName, "local_subnets.#", "2"),
					resource.TestCheckResourceAttr(rName, "client_cidr", "192.168.2.0/24"),
					resource.TestCheckResourceAttr(rName, "client_auth_type", "CERT"),
					resource.TestCheckResourceAttr(rName, "server_certificate.0.id", acceptance.HW_VPN_P2C_SERVER_CERTIFICATE_ID),
					resource.TestCheckResourceAttr(rName, "ssl_options.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(rName, "ssl_options.0.port", "1194"),
					resource.TestCheckResourceAttr(rName, "ssl_options.0.encryption_algorithm", "AES-256-GCM"),
					resource.TestCheckResourceAttr(rName, "ssl_options.0.is_compressed", "false"),
					resource.TestCheckResourceAttrSet(rName, "ssl_options.0.authentication_algorithm"),
					resource.TestCheckResourceAttr(rName, "client_ca_certificates_uploaded.0.name", "test-cert-update"),
					resource.TestCheckResourceAttrSet(rName, "tunnel_protocol"),
					resource.TestCheckResourceAttrSet(rName, "client_config"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				Config: testAccServer_updateAuthTypeToLocalPassword(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "p2c_vgw_id", acceptance.HW_VPN_P2C_GATEWAY_ID),
					resource.TestCheckResourceAttr(rName, "local_subnets.#", "1"),
					resource.TestCheckResourceAttr(rName, "client_cidr", "192.168.1.0/24"),
					resource.TestCheckResourceAttr(rName, "client_auth_type", "LOCAL_PASSWORD"),
					resource.TestCheckResourceAttr(rName, "server_certificate.0.id", acceptance.HW_VPN_P2C_SERVER_CERTIFICATE_ID),
					resource.TestCheckResourceAttr(rName, "ssl_options.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(rName, "ssl_options.0.port", "443"),
					resource.TestCheckResourceAttr(rName, "ssl_options.0.encryption_algorithm", "AES-128-GCM"),
					resource.TestCheckResourceAttr(rName, "ssl_options.0.is_compressed", "false"),
					resource.TestCheckResourceAttrSet(rName, "ssl_options.0.authentication_algorithm"),
					resource.TestCheckResourceAttr(rName, "client_ca_certificates_uploaded.#", "0"),
					resource.TestCheckResourceAttrSet(rName, "tunnel_protocol"),
					resource.TestCheckResourceAttrSet(rName, "client_config"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				Config: testAccServer_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "p2c_vgw_id", acceptance.HW_VPN_P2C_GATEWAY_ID),
					resource.TestCheckResourceAttr(rName, "local_subnets.#", "1"),
					resource.TestCheckResourceAttr(rName, "client_cidr", "192.168.1.0/24"),
					resource.TestCheckResourceAttr(rName, "client_auth_type", "CERT"),
					resource.TestCheckResourceAttr(rName, "server_certificate.0.id", acceptance.HW_VPN_P2C_SERVER_CERTIFICATE_ID),
					resource.TestCheckResourceAttr(rName, "ssl_options.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(rName, "ssl_options.0.port", "443"),
					resource.TestCheckResourceAttr(rName, "ssl_options.0.encryption_algorithm", "AES-128-GCM"),
					resource.TestCheckResourceAttr(rName, "ssl_options.0.is_compressed", "false"),
					resource.TestCheckResourceAttrSet(rName, "ssl_options.0.authentication_algorithm"),
					resource.TestCheckResourceAttr(rName, "client_ca_certificates_uploaded.0.name", "test-cert"),
					resource.TestCheckResourceAttrSet(rName, "client_ca_certificates_uploaded.0.id"),
					resource.TestCheckResourceAttrSet(rName, "client_ca_certificates_uploaded.0.issuer"),
					resource.TestCheckResourceAttrSet(rName, "client_ca_certificates_uploaded.0.signature_algorithm"),
					resource.TestCheckResourceAttrSet(rName, "client_ca_certificates_uploaded.0.subject"),
					resource.TestCheckResourceAttrSet(rName, "client_ca_certificates_uploaded.0.serial_number"),
					resource.TestCheckResourceAttrSet(rName, "client_ca_certificates_uploaded.0.expiration_time"),
					resource.TestCheckResourceAttrSet(rName, "client_ca_certificates_uploaded.0.created_at"),
					resource.TestCheckResourceAttrSet(rName, "client_ca_certificates_uploaded.0.updated_at"),
					resource.TestCheckResourceAttrSet(rName, "tunnel_protocol"),
					resource.TestCheckResourceAttrSet(rName, "client_config"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testServerImportState(rName),
				ImportStateVerifyIgnore: []string{"client_ca_certificates"},
			},
		},
	})
}

func testAccServer_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_vpn_server" "test" {
  p2c_vgw_id       = "%[1]s"
  local_subnets    = ["192.168.0.0/24"]
  client_cidr      = "192.168.1.0/24"
  client_auth_type = "CERT"

  server_certificate {
    id = "%[2]s"
  }

  ssl_options {
    protocol             = "TCP"
    port                 = 443
    encryption_algorithm = "AES-128-GCM"
    is_compressed        = false
  }
  
  client_ca_certificates {
    name    = "test-cert"
    content = "%[3]s"
  }
}
`, acceptance.HW_VPN_P2C_GATEWAY_ID, acceptance.HW_VPN_P2C_SERVER_CERTIFICATE_ID, acceptance.HW_VPN_P2C_CLIENT_CA_CERTIFICATE)
}

func testAccServer_update() string {
	return fmt.Sprintf(`
resource "huaweicloud_vpn_server" "test" {
  p2c_vgw_id       = "%[1]s"
  local_subnets    = ["192.168.0.0/24", "192.168.1.0/24"]
  client_cidr      = "192.168.2.0/24"
  client_auth_type = "CERT"

  server_certificate {
    id = "%[2]s"
  }

  ssl_options {
    protocol             = "TCP"
    port                 = 1194
    encryption_algorithm = "AES-256-GCM"
    is_compressed        = false
  }
  
  client_ca_certificates {
    name    = "test-cert-update"
    content = "%[3]s"
  }
}
`, acceptance.HW_VPN_P2C_GATEWAY_ID, acceptance.HW_VPN_P2C_SERVER_CERTIFICATE_ID, acceptance.HW_VPN_P2C_CLIENT_CA_CERTIFICATE)
}

func testAccServer_updateAuthTypeToLocalPassword() string {
	return fmt.Sprintf(`
resource "huaweicloud_vpn_server" "test" {
  p2c_vgw_id       = "%[1]s"
  local_subnets    = ["192.168.0.0/24"]
  client_cidr      = "192.168.1.0/24"
  client_auth_type = "LOCAL_PASSWORD"

  server_certificate {
    id = "%[2]s"
  }

  ssl_options {
    protocol             = "TCP"
    port                 = 443
    encryption_algorithm = "AES-128-GCM"
    is_compressed        = false
  }
}
`, acceptance.HW_VPN_P2C_GATEWAY_ID, acceptance.HW_VPN_P2C_SERVER_CERTIFICATE_ID, acceptance.HW_VPN_P2C_CLIENT_CA_CERTIFICATE)
}
