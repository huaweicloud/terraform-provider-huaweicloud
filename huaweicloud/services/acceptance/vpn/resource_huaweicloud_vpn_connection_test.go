package vpn

import (
	"fmt"
	"strings"
	"testing"

	"github.com/chnsz/golangsdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getConnectionResourceFunc(config *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getConnection: Query the VPN Connection detail
	var (
		getConnectionHttpUrl = "v5/{project_id}/vpn-connection/{id}"
		getConnectionProduct = "vpn"
	)
	getConnectionClient, err := config.NewServiceClient(getConnectionProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating Connection Client: %s", err)
	}

	getConnectionPath := getConnectionClient.Endpoint + getConnectionHttpUrl
	getConnectionPath = strings.ReplaceAll(getConnectionPath, "{project_id}", getConnectionClient.ProjectID)
	getConnectionPath = strings.ReplaceAll(getConnectionPath, "{id}", state.Primary.ID)

	getConnectionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getConnectionResp, err := getConnectionClient.Request("GET", getConnectionPath, &getConnectionOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Connection: %s", err)
	}
	return utils.FlattenResponse(getConnectionResp)
}

func TestAccConnection_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_vpn_connection.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getConnectionResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testConnection_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "vpn_type", "STATIC"),
					resource.TestCheckResourceAttr(rName, "ikepolicy.0.authentication_algorithm", "sha2-256"),
					resource.TestCheckResourceAttr(rName, "ikepolicy.0.encryption_algorithm", "aes-128"),
					resource.TestCheckResourceAttr(rName, "ikepolicy.0.lifetime_seconds", "86400"),
					resource.TestCheckResourceAttr(rName, "ipsecpolicy.0.authentication_algorithm", "sha2-256"),
					resource.TestCheckResourceAttr(rName, "ipsecpolicy.0.encryption_algorithm", "aes-128"),
					resource.TestCheckResourceAttr(rName, "ipsecpolicy.0.lifetime_seconds", "3600"),
					resource.TestCheckResourceAttrPair(rName, "gateway_id",
						"huaweicloud_vpn_gateway.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "gateway_ip",
						"huaweicloud_vpn_gateway.test", "master_eip.0.id"),
					resource.TestCheckResourceAttrPair(rName, "customer_gateway_id",
						"huaweicloud_vpn_customer_gateway.test", "id"),
				),
			},
			{
				Config: testConnection_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"-update"),
					resource.TestCheckResourceAttr(rName, "ikepolicy.0.authentication_algorithm", "sha2-512"),
					resource.TestCheckResourceAttr(rName, "ikepolicy.0.encryption_algorithm", "aes-256"),
					resource.TestCheckResourceAttr(rName, "ikepolicy.0.lifetime_seconds", "172800"),
					resource.TestCheckResourceAttr(rName, "ipsecpolicy.0.authentication_algorithm", "sha2-512"),
					resource.TestCheckResourceAttr(rName, "ipsecpolicy.0.encryption_algorithm", "aes-256"),
					resource.TestCheckResourceAttr(rName, "ipsecpolicy.0.lifetime_seconds", "7200"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"psk",
				},
			},
		},
	})
}

func testConnection_basic(name string) string {
	return fmt.Sprintf(`
%s
%s

resource "huaweicloud_vpn_connection" "test" {
  name                = "%s"
  gateway_id          = huaweicloud_vpn_gateway.test.id
  gateway_ip          = huaweicloud_vpn_gateway.test.master_eip[0].id
  customer_gateway_id = huaweicloud_vpn_customer_gateway.test.id
  peer_subnets        = ["192.168.55.0/24"]
  vpn_type            = "static"
  psk                 = "Test@123"
}
`, testGateway_basic(name), testCustomerGateway_basic(name), name)
}

func testConnection_update(name string) string {
	return fmt.Sprintf(`
%s
%s

resource "huaweicloud_vpn_connection" "test" {
  name                = "%s-update"
  gateway_id          = huaweicloud_vpn_gateway.test.id
  gateway_ip          = huaweicloud_vpn_gateway.test.master_eip[0].id
  customer_gateway_id = huaweicloud_vpn_customer_gateway.test.id
  peer_subnets        = ["192.168.55.0/24"]
  vpn_type            = "static"
  psk                 = "Test@123"

  ikepolicy {
    authentication_algorithm = "sha2-512"
    encryption_algorithm     = "aes-256"
    lifetime_seconds         = 172800
  }

  ipsecpolicy {
    authentication_algorithm = "sha2-512"
    encryption_algorithm     = "aes-256"
    lifetime_seconds         = 7200
  }
}
`, testGateway_basic(name), testCustomerGateway_basic(name), name)
}
