package vpn

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

func getCustomerGatewayResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getCustomerGateway: Query the VPN customer gateway detail
	var (
		getCustomerGatewayHttpUrl = "v5/{project_id}/customer-gateways/{id}"
		getCustomerGatewayProduct = "vpn"
	)
	getCustomerGatewayClient, err := conf.NewServiceClient(getCustomerGatewayProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating CustomerGateway Client: %s", err)
	}

	getCustomerGatewayPath := getCustomerGatewayClient.Endpoint + getCustomerGatewayHttpUrl
	getCustomerGatewayPath = strings.ReplaceAll(getCustomerGatewayPath, "{project_id}", getCustomerGatewayClient.ProjectID)
	getCustomerGatewayPath = strings.ReplaceAll(getCustomerGatewayPath, "{id}", state.Primary.ID)

	getCustomerGatewayOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getCustomerGatewayResp, err := getCustomerGatewayClient.Request("GET", getCustomerGatewayPath, &getCustomerGatewayOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving CustomerGateway: %s", err)
	}
	return utils.FlattenResponse(getCustomerGatewayResp)
}

func TestAccCustomerGateway_basic_withDeprecatedFields(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	nameUpdate := name + "-update"
	rName := "huaweicloud_vpn_customer_gateway.test"
	ipAddress := "172.16.1.2"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getCustomerGatewayResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCustomerGateway_basic_withDeprecatedFields(name, ipAddress),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "ip", ipAddress),
					resource.TestCheckResourceAttr(rName, "tags.key", "val"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
				),
			},
			{
				Config: testCustomerGateway_update_withDeprecatedFields(nameUpdate, ipAddress),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", nameUpdate),
					resource.TestCheckResourceAttr(rName, "ip", ipAddress),
					resource.TestCheckResourceAttr(rName, "tags.key", "val"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar-update"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"ip", "route_mode",
				},
			},
		},
	})
}

func TestAccCustomerGateway_certificate_withDeprecatedFields(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_vpn_customer_gateway.test"
	ipAddress := "172.16.2.3"
	certificateContent := acceptance.HW_CERTIFICATE_CONTENT
	certificateContentUpdate := acceptance.HW_CERTIFICATE_CONTENT_UPDATE

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getCustomerGatewayResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckUpdateCertificateContent(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCustomerGateway_certificate_withDeprecatedFields(name, ipAddress, certificateContent),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "ip", ipAddress),
					resource.TestCheckResourceAttrSet(rName, "serial_number"),
					resource.TestCheckResourceAttrSet(rName, "signature_algorithm"),
					resource.TestCheckResourceAttrSet(rName, "issuer"),
					resource.TestCheckResourceAttrSet(rName, "subject"),
					resource.TestCheckResourceAttrSet(rName, "expire_time"),
					resource.TestCheckResourceAttrSet(rName, "is_updatable"),
				),
			},
			{
				Config: testCustomerGateway_certificate_withDeprecatedFields(name, ipAddress, certificateContentUpdate),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "ip", ipAddress),
					resource.TestCheckResourceAttrSet(rName, "serial_number"),
					resource.TestCheckResourceAttrSet(rName, "signature_algorithm"),
					resource.TestCheckResourceAttrSet(rName, "issuer"),
					resource.TestCheckResourceAttrSet(rName, "subject"),
					resource.TestCheckResourceAttrSet(rName, "expire_time"),
					resource.TestCheckResourceAttrSet(rName, "is_updatable"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"certificate_content", "ip", "route_mode",
				},
			},
		},
	})
}

func TestAccCustomerGateway_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	nameUpdate := name + "-update"
	rName := "huaweicloud_vpn_customer_gateway.test"
	ipAddress := "172.16.1.4"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getCustomerGatewayResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCustomerGateway_basic(name, ipAddress),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "id_value", ipAddress),
					resource.TestCheckResourceAttr(rName, "tags.key", "val"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
				),
			},
			{
				Config: testCustomerGateway_update(nameUpdate, ipAddress),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", nameUpdate),
					resource.TestCheckResourceAttr(rName, "id_value", ipAddress),
					resource.TestCheckResourceAttr(rName, "tags.key", "val"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar-update"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"route_mode",
				},
			},
		},
	})
}

func TestAccCustomerGateway_certificate(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_vpn_customer_gateway.test"
	ipAddress := "172.16.2.5"
	certificateContent := acceptance.HW_CERTIFICATE_CONTENT
	certificateContentUpdate := acceptance.HW_CERTIFICATE_CONTENT_UPDATE

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getCustomerGatewayResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckUpdateCertificateContent(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCustomerGateway_certificate(name, ipAddress, certificateContent),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "id_value", ipAddress),
					resource.TestCheckResourceAttrSet(rName, "serial_number"),
					resource.TestCheckResourceAttrSet(rName, "signature_algorithm"),
					resource.TestCheckResourceAttrSet(rName, "issuer"),
					resource.TestCheckResourceAttrSet(rName, "subject"),
					resource.TestCheckResourceAttrSet(rName, "expire_time"),
					resource.TestCheckResourceAttrSet(rName, "is_updatable"),
				),
			},
			{
				Config: testCustomerGateway_certificate(name, ipAddress, certificateContentUpdate),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "id_value", ipAddress),
					resource.TestCheckResourceAttrSet(rName, "serial_number"),
					resource.TestCheckResourceAttrSet(rName, "signature_algorithm"),
					resource.TestCheckResourceAttrSet(rName, "issuer"),
					resource.TestCheckResourceAttrSet(rName, "subject"),
					resource.TestCheckResourceAttrSet(rName, "expire_time"),
					resource.TestCheckResourceAttrSet(rName, "is_updatable"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"certificate_content", "ip", "route_mode",
				},
			},
		},
	})
}

func testCustomerGateway_basic_withDeprecatedFields(name, ipAddress string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpn_customer_gateway" "test" {
  name = "%s"
  ip   = "%s"

  tags = {
    key = "val"
    foo = "bar"
  }
}`, name, ipAddress)
}

func testCustomerGateway_update_withDeprecatedFields(name, ipAddress string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpn_customer_gateway" "test" {
  name = "%s"
  ip   = "%s"

  tags = {
    key = "val"
    foo = "bar-update"
  }
}`, name, ipAddress)
}

func testCustomerGateway_certificate_withDeprecatedFields(name, ipAddress, certificateContent string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpn_customer_gateway" "test" {
  name                = "%s"
  ip                  = "%s"
  certificate_content = "%s"
}`, name, ipAddress, certificateContent)
}

func testCustomerGateway_basic(name, ipAddress string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpn_customer_gateway" "test" {
  name     = "%s"
  id_value = "%s"

  tags = {
    key = "val"
    foo = "bar"
  }
}`, name, ipAddress)
}

func testCustomerGateway_update(name, ipAddress string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpn_customer_gateway" "test" {
  name     = "%s"
  id_value = "%s"

  tags = {
    key = "val"
    foo = "bar-update"
  }
}`, name, ipAddress)
}

func testCustomerGateway_certificate(name, ipAddress, certificateContent string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpn_customer_gateway" "test" {
  name                = "%s"
  id_value            = "%s"
  certificate_content = "%s"
}`, name, ipAddress, certificateContent)
}
