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

func getGatewayResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getGateway: Query the VPN gateway detail
	var (
		getGatewayHttpUrl = "v5/{project_id}/vpn-gateways/{id}"
		getGatewayProduct = "vpn"
	)
	getGatewayClient, err := conf.NewServiceClient(getGatewayProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating Gateway Client: %s", err)
	}

	getGatewayPath := getGatewayClient.Endpoint + getGatewayHttpUrl
	getGatewayPath = strings.ReplaceAll(getGatewayPath, "{project_id}", getGatewayClient.ProjectID)
	getGatewayPath = strings.ReplaceAll(getGatewayPath, "{id}", state.Primary.ID)

	getGatewayOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getGatewayResp, err := getGatewayClient.Request("GET", getGatewayPath, &getGatewayOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Gateway: %s", err)
	}
	return utils.FlattenResponse(getGatewayResp)
}

func TestAccGateway_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_vpn_gateway.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getGatewayResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testGateway_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "ha_mode", "active-active"),
					resource.TestCheckResourceAttrPair(rName, "connect_subnet", "huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "vpc_id", "huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
					resource.TestCheckResourceAttrPair(rName, "local_subnets.0", "huaweicloud_vpc_subnet.test", "cidr"),
					resource.TestCheckResourceAttrPair(rName, "eip1.0.id", "huaweicloud_vpc_eip.test1", "id"),
					resource.TestCheckResourceAttrPair(rName, "eip2.0.id", "huaweicloud_vpc_eip.test2", "id"),
					resource.TestCheckResourceAttrPair(rName, "availability_zones.0",
						"data.huaweicloud_vpn_gateway_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttrPair(rName, "availability_zones.1",
						"data.huaweicloud_vpn_gateway_availability_zones.test", "names.1"),
					resource.TestCheckResourceAttr(rName, "tags.key", "val"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
				),
			},
			{
				Config: testGateway_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"-update"),
					resource.TestCheckResourceAttrPair(rName, "local_subnets.0", "huaweicloud_vpc_subnet.test", "cidr"),
					resource.TestCheckResourceAttr(rName, "local_subnets.1", "192.168.2.0/24"),
					resource.TestCheckResourceAttr(rName, "flavor", "Professional2"),
					resource.TestCheckResourceAttrPair(rName, "availability_zones.0",
						"data.huaweicloud_vpn_gateway_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttrPair(rName, "availability_zones.1",
						"data.huaweicloud_vpn_gateway_availability_zones.test", "names.1"),
					resource.TestCheckResourceAttr(rName, "tags.key", "val"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar-update"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerifyIgnore: []string{"delete_eip_on_termination"},
			},
		},
	})
}

func TestAccGateway_UpdateWithEpsId(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_vpn_gateway.test"
	srcEPS := acceptance.HW_ENTERPRISE_PROJECT_ID_TEST
	destEPS := acceptance.HW_ENTERPRISE_MIGRATE_PROJECT_ID_TEST

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getGatewayResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMigrateEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testGateway_withEpsId(name, srcEPS),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "local_subnets.1", "192.168.2.0/24"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", srcEPS),
				),
			},
			{
				Config: testGateway_withEpsId(name, destEPS),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "local_subnets.1", "192.168.2.0/24"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", destEPS),
				),
			},
		},
	})
}

func TestAccGateway_activeStandbyHAMode(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_vpn_gateway.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getGatewayResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testGateway_activeStandbyHAMode(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "ha_mode", "active-standby"),
					resource.TestCheckResourceAttrPair(rName, "connect_subnet", "huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "vpc_id", "huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
					resource.TestCheckResourceAttrPair(rName, "local_subnets.0", "huaweicloud_vpc_subnet.test", "cidr"),
					resource.TestCheckResourceAttrPair(rName, "eip1.0.id", "huaweicloud_vpc_eip.test1", "id"),
					resource.TestCheckResourceAttrPair(rName, "eip2.0.id", "huaweicloud_vpc_eip.test2", "id"),
					resource.TestCheckResourceAttrPair(rName, "availability_zones.0",
						"data.huaweicloud_vpn_gateway_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttrPair(rName, "availability_zones.1",
						"data.huaweicloud_vpn_gateway_availability_zones.test", "names.1"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerifyIgnore: []string{"delete_eip_on_termination"},
			},
		},
	})
}

func TestAccGateway_certificate(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_vpn_gateway.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getGatewayResourceFunc,
	)

	cert := certificate{
		name:             "test_gateway_certificate",
		content:          acceptance.HW_GM_CERTIFICATE_CONTENT,
		privateKey:       acceptance.HW_GM_CERTIFICATE_PRIVATE_KEY,
		certificateChain: acceptance.HW_GM_CERTIFICATE_CHAIN,
		encCertificate:   acceptance.HW_GM_ENC_CERTIFICATE_CONTENT,
		encPrivateKey:    acceptance.HW_GM_ENC_CERTIFICATE_PRIVATE_KEY,
	}

	certUpdate := certificate{
		name:             "test_gateway_certificate_update",
		content:          acceptance.HW_NEW_GM_CERTIFICATE_CONTENT,
		privateKey:       acceptance.HW_NEW_GM_CERTIFICATE_PRIVATE_KEY,
		certificateChain: acceptance.HW_NEW_GM_CERTIFICATE_CHAIN,
		encCertificate:   acceptance.HW_NEW_GM_ENC_CERTIFICATE_CONTENT,
		encPrivateKey:    acceptance.HW_NEW_GM_ENC_CERTIFICATE_PRIVATE_KEY,
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGMCertificate(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testGateway_GMcertificate(name, cert),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "certificate.#", "1"),
					resource.TestCheckResourceAttr(rName, "certificate.0.name", cert.name),
					resource.TestCheckResourceAttrSet(rName, "certificate.0.content"),
					resource.TestCheckResourceAttrSet(rName, "certificate.0.private_key"),
					resource.TestCheckResourceAttrSet(rName, "certificate.0.certificate_chain"),
					resource.TestCheckResourceAttrSet(rName, "certificate.0.enc_certificate"),
					resource.TestCheckResourceAttrSet(rName, "certificate.0.enc_private_key"),
					resource.TestCheckResourceAttrSet(rName, "certificate.0.certificate_id"),
					resource.TestCheckResourceAttrSet(rName, "certificate.0.status"),
					resource.TestCheckResourceAttrSet(rName, "certificate.0.issuer"),
					resource.TestCheckResourceAttrSet(rName, "certificate.0.signature_algorithm"),
					resource.TestCheckResourceAttrSet(rName, "certificate.0.certificate_serial_number"),
					resource.TestCheckResourceAttrSet(rName, "certificate.0.certificate_subject"),
					resource.TestCheckResourceAttrSet(rName, "certificate.0.certificate_expire_time"),
					resource.TestCheckResourceAttrSet(rName, "certificate.0.certificate_chain_serial_number"),
					resource.TestCheckResourceAttrSet(rName, "certificate.0.certificate_chain_subject"),
					resource.TestCheckResourceAttrSet(rName, "certificate.0.certificate_chain_expire_time"),
					resource.TestCheckResourceAttrSet(rName, "certificate.0.enc_certificate_serial_number"),
					resource.TestCheckResourceAttrSet(rName, "certificate.0.enc_certificate_subject"),
					resource.TestCheckResourceAttrSet(rName, "certificate.0.enc_certificate_expire_time"),
					resource.TestCheckResourceAttrSet(rName, "certificate.0.created_at"),
					resource.TestCheckResourceAttrSet(rName, "certificate.0.updated_at"),
				),
			},
			{
				Config: testGateway_GMcertificate(name, certUpdate),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "certificate.#", "1"),
					resource.TestCheckResourceAttr(rName, "certificate.0.name", certUpdate.name),
					resource.TestCheckResourceAttrSet(rName, "certificate.0.content"),
					resource.TestCheckResourceAttrSet(rName, "certificate.0.private_key"),
					resource.TestCheckResourceAttrSet(rName, "certificate.0.certificate_chain"),
					resource.TestCheckResourceAttrSet(rName, "certificate.0.enc_certificate"),
					resource.TestCheckResourceAttrSet(rName, "certificate.0.enc_private_key"),
					resource.TestCheckResourceAttrSet(rName, "certificate.0.certificate_id"),
					resource.TestCheckResourceAttrSet(rName, "certificate.0.status"),
					resource.TestCheckResourceAttrSet(rName, "certificate.0.issuer"),
					resource.TestCheckResourceAttrSet(rName, "certificate.0.signature_algorithm"),
					resource.TestCheckResourceAttrSet(rName, "certificate.0.certificate_serial_number"),
					resource.TestCheckResourceAttrSet(rName, "certificate.0.certificate_subject"),
					resource.TestCheckResourceAttrSet(rName, "certificate.0.certificate_expire_time"),
					resource.TestCheckResourceAttrSet(rName, "certificate.0.certificate_chain_serial_number"),
					resource.TestCheckResourceAttrSet(rName, "certificate.0.certificate_chain_subject"),
					resource.TestCheckResourceAttrSet(rName, "certificate.0.certificate_chain_expire_time"),
					resource.TestCheckResourceAttrSet(rName, "certificate.0.enc_certificate_serial_number"),
					resource.TestCheckResourceAttrSet(rName, "certificate.0.enc_certificate_subject"),
					resource.TestCheckResourceAttrSet(rName, "certificate.0.enc_certificate_expire_time"),
					resource.TestCheckResourceAttrSet(rName, "certificate.0.created_at"),
					resource.TestCheckResourceAttrSet(rName, "certificate.0.updated_at"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerifyIgnore: []string{"certificate", "delete_eip_on_termination"},
			},
		},
	})
}

func TestAccGateway_deprecated(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_vpn_gateway.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getGatewayResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testGateway_deprecated(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "ha_mode", "active-standby"),
					resource.TestCheckResourceAttrPair(rName, "connect_subnet", "huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "vpc_id", "huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
					resource.TestCheckResourceAttrPair(rName, "local_subnets.0", "huaweicloud_vpc_subnet.test", "cidr"),
					resource.TestCheckResourceAttrPair(rName, "master_eip.0.id", "huaweicloud_vpc_eip.test1", "id"),
					resource.TestCheckResourceAttrPair(rName, "slave_eip.0.id", "huaweicloud_vpc_eip.test2", "id"),
					resource.TestCheckResourceAttrPair(rName, "availability_zones.0",
						"data.huaweicloud_vpn_gateway_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttrPair(rName, "availability_zones.1",
						"data.huaweicloud_vpn_gateway_availability_zones.test", "names.1"),
				),
			},
			{
				Config: testGateway_deprecated_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"-update"),
					resource.TestCheckResourceAttrPair(rName, "local_subnets.0", "huaweicloud_vpc_subnet.test", "cidr"),
					resource.TestCheckResourceAttr(rName, "local_subnets.1", "192.168.2.0/24"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerifyIgnore: []string{"delete_eip_on_termination"},
			},
		},
	})
}

func TestAccGateway_withER(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_vpn_gateway.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getGatewayResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testGateway_withER(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "network_type", "private"),
					resource.TestCheckResourceAttr(rName, "attachment_type", "er"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(rName, "access_private_ip_1", "172.16.0.99"),
					resource.TestCheckResourceAttr(rName, "access_private_ip_2", "172.16.0.100"),
					resource.TestCheckResourceAttrPair(rName, "er_id", "huaweicloud_er_instance.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "access_vpc_id", "huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "access_subnet_id", "huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "availability_zones.0",
						"data.huaweicloud_vpn_gateway_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttrPair(rName, "availability_zones.1",
						"data.huaweicloud_vpn_gateway_availability_zones.test", "names.1"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerifyIgnore: []string{"delete_eip_on_termination"},
			},
		},
	})
}

func TestAccGateway_deleteEipOnTermination(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_vpn_gateway.test"
	eip1RName := "huaweicloud_vpc_eip.test1"
	eip2RName := "huaweicloud_vpc_eip.test2"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getGatewayResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccGateway_deleteEipOnTermination(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "ha_mode", "active-active"),
					resource.TestCheckResourceAttr(rName, "delete_eip_on_termination", "false"),
					resource.TestCheckResourceAttrPair(rName, "connect_subnet", "huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "vpc_id", "huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
					resource.TestCheckResourceAttrPair(rName, "local_subnets.0", "huaweicloud_vpc_subnet.test", "cidr"),
					resource.TestCheckResourceAttrPair(rName, "eip1.0.id", "huaweicloud_vpc_eip.test1", "id"),
					resource.TestCheckResourceAttrPair(rName, "eip2.0.id", "huaweicloud_vpc_eip.test2", "id"),
					resource.TestCheckResourceAttrPair(rName, "availability_zones.0",
						"data.huaweicloud_vpn_gateway_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttrPair(rName, "availability_zones.1",
						"data.huaweicloud_vpn_gateway_availability_zones.test", "names.1"),
					resource.TestCheckResourceAttr(rName, "tags.key", "val"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
				),
			},
			{
				Config: testGateway_base(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(eip1RName, "id"),
					resource.TestCheckResourceAttrSet(eip2RName, "id"),
				),
			},
		},
	})
}

func testGateway_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_vpn_gateway_availability_zones" "test" {
  flavor          = "professional1"
  attachment_type = "vpc"
}

resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  name       = "%[1]s"
  vpc_id     = huaweicloud_vpc.test.id
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
}

resource "huaweicloud_vpc_eip" "test1" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "%[1]s-1"
    size        = 8
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "huaweicloud_vpc_eip" "test2" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "%[1]s-2"
    size        = 8
    share_type  = "PER"
    charge_mode = "traffic"
  }
}
`, name)
}

func testGateway_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_vpn_gateway" "test" {
  name               = "%s"
  vpc_id             = huaweicloud_vpc.test.id
  local_subnets      = [huaweicloud_vpc_subnet.test.cidr]
  connect_subnet     = huaweicloud_vpc_subnet.test.id
  availability_zones = [
    data.huaweicloud_vpn_gateway_availability_zones.test.names[0],
    data.huaweicloud_vpn_gateway_availability_zones.test.names[1]
  ]

  eip1 {
    id = huaweicloud_vpc_eip.test1.id
  }

  eip2 {
    id = huaweicloud_vpc_eip.test2.id
  }

  tags = {
    key = "val"
    foo = "bar"
  }
}
`, testGateway_base(name), name)
}

func testGateway_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_vpn_gateway" "test" {
  name               = "%s-update"
  vpc_id             = huaweicloud_vpc.test.id
  local_subnets      = [huaweicloud_vpc_subnet.test.cidr, "192.168.2.0/24"]
  connect_subnet     = huaweicloud_vpc_subnet.test.id
  flavor             = "Professional2"
  
  availability_zones = [
    data.huaweicloud_vpn_gateway_availability_zones.test.names[0],
    data.huaweicloud_vpn_gateway_availability_zones.test.names[1]
  ]

  eip1 {
    id = huaweicloud_vpc_eip.test1.id
  }

  eip2 {
    id = huaweicloud_vpc_eip.test2.id
  }

  tags = {
    key = "val"
    foo = "bar-update"
  }
}
`, testGateway_base(name), name)
}

func testGateway_withEpsId(name, epsId string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_vpn_gateway" "test" {
  name                  = "%s"
  vpc_id                = huaweicloud_vpc.test.id
  enterprise_project_id = "%s"
  local_subnets         = [huaweicloud_vpc_subnet.test.cidr, "192.168.2.0/24"]
  connect_subnet        = huaweicloud_vpc_subnet.test.id
  availability_zones    = [
    data.huaweicloud_vpn_gateway_availability_zones.test.names[0],
    data.huaweicloud_vpn_gateway_availability_zones.test.names[1]
  ]

  eip1 {
    id = huaweicloud_vpc_eip.test1.id
  }

  eip2 {
    id = huaweicloud_vpc_eip.test2.id
  }
}
`, testGateway_base(name), name, epsId)
}

func testGateway_activeStandbyHAMode(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_vpn_gateway" "test" {
  name               = "%s"
  vpc_id             = huaweicloud_vpc.test.id
  ha_mode            = "active-standby"
  local_subnets      = [huaweicloud_vpc_subnet.test.cidr]
  connect_subnet     = huaweicloud_vpc_subnet.test.id
  availability_zones = [
    data.huaweicloud_vpn_gateway_availability_zones.test.names[0],
    data.huaweicloud_vpn_gateway_availability_zones.test.names[1]
  ]

  eip1 {
    id = huaweicloud_vpc_eip.test1.id
  }

  eip2 {
    id = huaweicloud_vpc_eip.test2.id
  }
}
`, testGateway_base(name), name)
}

func testGateway_withER(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "172.16.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  name       = "%[1]s"
  vpc_id     = huaweicloud_vpc.test.id
  cidr       = "172.16.0.0/24"
  gateway_ip = "172.16.0.1"
}

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_er_instance" "test" {
  availability_zones = [
    data.huaweicloud_availability_zones.test.names[0],
    data.huaweicloud_availability_zones.test.names[3]
  ]

  name = "%[1]s"
  asn  = "65000"
}

data "huaweicloud_vpn_gateway_availability_zones" "test" {
  flavor          = "professional1"
  attachment_type = "er"
}

resource "huaweicloud_vpn_gateway" "test" {
  name               = "%[1]s"
  network_type       = "private"
  attachment_type    = "er"
  er_id              = huaweicloud_er_instance.test.id
  availability_zones = [
    data.huaweicloud_vpn_gateway_availability_zones.test.names[0],
    data.huaweicloud_vpn_gateway_availability_zones.test.names[1]
  ]

  access_vpc_id    = huaweicloud_vpc.test.id
  access_subnet_id = huaweicloud_vpc_subnet.test.id
  
  access_private_ip_1 = "172.16.0.99"
  access_private_ip_2 = "172.16.0.100"
}
`, name)
}

func testGateway_deprecated(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_vpn_gateway" "test" {
  name               = "%s"
  vpc_id             = huaweicloud_vpc.test.id
  local_subnets      = [huaweicloud_vpc_subnet.test.cidr]
  connect_subnet     = huaweicloud_vpc_subnet.test.id
  availability_zones = [
    data.huaweicloud_vpn_gateway_availability_zones.test.names[0],
    data.huaweicloud_vpn_gateway_availability_zones.test.names[1]
  ]

  master_eip {
    id = huaweicloud_vpc_eip.test1.id
  }

  slave_eip {
    id = huaweicloud_vpc_eip.test2.id
  }
}
`, testGateway_base(name), name)
}

func testGateway_deprecated_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_vpn_gateway" "test" {
  name               = "%s-update"
  vpc_id             = huaweicloud_vpc.test.id
  local_subnets      = [huaweicloud_vpc_subnet.test.cidr, "192.168.2.0/24"]
  connect_subnet     = huaweicloud_vpc_subnet.test.id
  availability_zones = [
    data.huaweicloud_vpn_gateway_availability_zones.test.names[0],
    data.huaweicloud_vpn_gateway_availability_zones.test.names[1]
  ]

  master_eip {
    id = huaweicloud_vpc_eip.test1.id
  }

  slave_eip {
    id = huaweicloud_vpc_eip.test2.id
  }
}
`, testGateway_base(name), name)
}

func testGateway_GMcertificate(name string, cert certificate) string {
	return fmt.Sprintf(`
data "huaweicloud_vpn_gateway_availability_zones" "test" {
  attachment_type = "vpc"
  flavor          = "gm"
}

resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  name       = "%[1]s"
  vpc_id     = huaweicloud_vpc.test.id
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
}

resource "huaweicloud_vpn_gateway" "test" {
  name               = "%[1]s"
  vpc_id             = huaweicloud_vpc.test.id
  flavor             = "GM"
  network_type       = "private"
  local_subnets      = [huaweicloud_vpc_subnet.test.cidr]
  connect_subnet     = huaweicloud_vpc_subnet.test.id
  availability_zones = [
    data.huaweicloud_vpn_gateway_availability_zones.test.names[0],
    data.huaweicloud_vpn_gateway_availability_zones.test.names[1]
  ]

  certificate {
    name              = "%[2]s"
    content           = "%[3]s"
    private_key       = "%[4]s"
    certificate_chain = "%[5]s"
    enc_certificate   = "%[6]s"
    enc_private_key   = "%[7]s"
  }
}
`, name, cert.name, cert.content, cert.privateKey, cert.certificateChain, cert.encCertificate, cert.encPrivateKey)
}

func testAccGateway_deleteEipOnTermination(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_vpn_gateway" "test" {
  name                      = "%s"
  vpc_id                    = huaweicloud_vpc.test.id
  local_subnets             = [huaweicloud_vpc_subnet.test.cidr]
  connect_subnet            = huaweicloud_vpc_subnet.test.id
  delete_eip_on_termination = false

  availability_zones = [
    data.huaweicloud_vpn_gateway_availability_zones.test.names[0],
    data.huaweicloud_vpn_gateway_availability_zones.test.names[1]
  ]

  eip1 {
    id = huaweicloud_vpc_eip.test1.id
  }

  eip2 {
    id = huaweicloud_vpc_eip.test2.id
  }

  tags = {
    key = "val"
    foo = "bar"
  }
}
`, testGateway_base(name), name)
}

type certificate struct {
	name             string
	content          string
	privateKey       string
	certificateChain string
	encCertificate   string
	encPrivateKey    string
}
