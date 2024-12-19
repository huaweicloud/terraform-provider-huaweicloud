package workspace

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/workspace/v2/services"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getServiceFunc(conf *config.Config, _ *terraform.ResourceState) (interface{}, error) {
	client, err := conf.WorkspaceV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Workspace v2 client: %s", err)
	}
	resp, err := services.Get(client)
	if resp.Status == "CLOSED" {
		return nil, golangsdk.ErrDefault404{}
	}
	return resp, err
}

func TestAccService_basic(t *testing.T) {
	var (
		service      services.Service
		resourceName = "huaweicloud_workspace_service.test"
		rName        = acceptance.RandomAccResourceNameWithDash()
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&service,
		getServiceFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccService_basic_step1(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id", "huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "network_ids.0",
						"huaweicloud_vpc_subnet.master", "id"),
					resource.TestCheckResourceAttr(resourceName, "auth_type", "LITE_AS"),
					resource.TestCheckResourceAttr(resourceName, "access_mode", "INTERNET"),
					resource.TestCheckResourceAttrSet(resourceName, "management_subnet_cidr"),
					resource.TestCheckResourceAttrSet(resourceName, "infrastructure_security_group.0.name"),
					resource.TestCheckResourceAttrSet(resourceName, "infrastructure_security_group.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "desktop_security_group.0.name"),
					resource.TestCheckResourceAttrSet(resourceName, "desktop_security_group.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "internet_access_address"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttr(resourceName, "is_locked", "0"),
				),
			},
			{
				Config: testAccService_basic_step2(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "network_ids.0",
						"huaweicloud_vpc_subnet.master", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "network_ids.1",
						"huaweicloud_vpc_subnet.standby", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "internet_access_address"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_id", rName),
					resource.TestCheckResourceAttr(resourceName, "enterprise_id", rName),
					resource.TestCheckResourceAttr(resourceName, "otp_config_info.0.enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "otp_config_info.0.receive_mode", "VMFA"),
				),
			},
			{
				Config: testAccService_basic_step3(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "otp_config_info.0.rule_type", "ACCESS_MODE"),
					resource.TestCheckResourceAttr(resourceName, "otp_config_info.0.rule", "PRIVATE"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCustomImportStateIdFunc(resourceName, "NA"),
			},
		},
	})
}

func testAccCustomImportStateIdFunc(resourceName, customId string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		_, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("the resource (%s) is not found in the tfstate", resourceName)
		}
		return customId, nil
	}
}

func TestAccService_internetAccessPort(t *testing.T) {
	var (
		service      services.Service
		resourceName = "huaweicloud_workspace_service.test"
		rName        = acceptance.RandomAccResourceNameWithDash()
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&service,
		getServiceFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkspaceInternetAccessPort(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccService_basic_step1(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resourceName, "internet_access_port"),
				),
			},
			{
				Config: testAccService_internetAccessPort_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "internet_access_port", acceptance.HW_WORKSPACE_INTERNET_ACCESS_PORT),
				),
			},
		},
	})
}

func retrieveAdDomain(domainNames []string) string {
	if len(domainNames) < 1 {
		return ""
	}

	re := regexp.MustCompile(`^[^.]+\.(.*)$`)
	match := re.FindStringSubmatch(domainNames[0])
	if len(match) < 2 {
		log.Printf("[ERROR] No match found for the AD domain: %s", domainNames[0])
		return ""
	}

	return match[1]
}

// These method are parsed before precheck.
// To avoid subscript out of bounds, the length is checked before the value is taken.
func retrieveActiveDomainIpAddress(ipAddresses []string) string {
	if len(ipAddresses) > 0 {
		return ipAddresses[0]
	}
	return ""
}

func retrieveStandbyDomainIpAddress(ipAddresses []string) string {
	if len(ipAddresses) > 1 {
		return ipAddresses[1]
	}
	return ""
}

func retrieveActiveDomainName(domainNames []string) string {
	if len(domainNames) > 0 {
		return domainNames[0]
	}
	return ""
}

func retrieveStandbyDomainName(domainNames []string) string {
	if len(domainNames) > 1 {
		return domainNames[1]
	}
	return ""
}

func TestAccService_localAD(t *testing.T) {
	var (
		service      services.Service
		resourceName = "huaweicloud_workspace_service.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&service,
		getServiceFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkspaceAD(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccService_localAD_step1(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "auth_type", "LOCAL_AD"),
					resource.TestCheckResourceAttr(resourceName, "ad_domain.0.name",
						retrieveAdDomain(strings.Split(acceptance.HW_WORKSPACE_AD_DOMAIN_NAMES, ","))),
					resource.TestCheckResourceAttr(resourceName, "ad_domain.0.admin_account", acceptance.HW_WORKSPACE_AD_SERVER_ACCOUNT),
					resource.TestCheckResourceAttr(resourceName, "ad_domain.0.password", acceptance.HW_WORKSPACE_AD_SERVER_PWD),
					resource.TestCheckResourceAttr(resourceName, "ad_domain.0.active_domain_name",
						retrieveActiveDomainName(strings.Split(acceptance.HW_WORKSPACE_AD_DOMAIN_NAMES, ","))),
					resource.TestCheckResourceAttr(resourceName, "ad_domain.0.active_domain_ip",
						retrieveActiveDomainIpAddress(strings.Split(acceptance.HW_WORKSPACE_AD_DOMAIN_IPS, ","))),
					resource.TestCheckResourceAttr(resourceName, "ad_domain.0.active_dns_ip",
						retrieveActiveDomainIpAddress(strings.Split(acceptance.HW_WORKSPACE_AD_DOMAIN_IPS, ","))),
					resource.TestCheckResourceAttr(resourceName, "access_mode", "INTERNET"),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", acceptance.HW_WORKSPACE_AD_VPC_ID),
					resource.TestCheckResourceAttr(resourceName, "network_ids.0", acceptance.HW_WORKSPACE_AD_NETWORK_ID),
					resource.TestCheckResourceAttrSet(resourceName, "infrastructure_security_group.0.name"),
					resource.TestCheckResourceAttrSet(resourceName, "infrastructure_security_group.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "desktop_security_group.0.name"),
					resource.TestCheckResourceAttrSet(resourceName, "desktop_security_group.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "internet_access_address"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
				),
			},
			{
				Config: testAccService_localAD_step2(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "ad_domain.0.standby_domain_name",
						retrieveStandbyDomainName(strings.Split(acceptance.HW_WORKSPACE_AD_DOMAIN_NAMES, ","))),
					resource.TestCheckResourceAttr(resourceName, "ad_domain.0.standby_domain_ip",
						retrieveStandbyDomainIpAddress(strings.Split(acceptance.HW_WORKSPACE_AD_DOMAIN_IPS, ","))),
					resource.TestCheckResourceAttr(resourceName, "ad_domain.0.standby_dns_ip",
						retrieveStandbyDomainIpAddress(strings.Split(acceptance.HW_WORKSPACE_AD_DOMAIN_IPS, ","))),
					resource.TestCheckResourceAttrSet(resourceName, "internet_access_address"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_id", "custom-workspace-service"),
					resource.TestCheckResourceAttr(resourceName, "otp_config_info.0.enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "otp_config_info.0.receive_mode", "VMFA"),
				),
			},
			{
				Config: testAccService_localAD_step3(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "otp_config_info.0.rule_type", "ACCESS_MODE"),
					resource.TestCheckResourceAttr(resourceName, "otp_config_info.0.rule", "PRIVATE"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"ad_domain.0.password",
				},
			},
		},
	})
}

func TestAccService_internetAccessPort_localAD(t *testing.T) {
	var (
		service      services.Service
		resourceName = "huaweicloud_workspace_service.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&service,
		getServiceFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkspaceAD(t)
			acceptance.TestAccPreCheckWorkspaceInternetAccessPort(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccService_localAD_step1(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resourceName, "internet_access_port"),
				),
			},
			{
				Config: testAccService_localAD_internetAccessPort_update(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "internet_access_port", acceptance.HW_WORKSPACE_INTERNET_ACCESS_PORT),
				),
			},
		},
	})
}

func testAccService_base(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/20"
}

resource "huaweicloud_vpc_subnet" "master" {
  vpc_id = huaweicloud_vpc.test.id

  name       = "%[1]s-master"
  cidr       = cidrsubnet(huaweicloud_vpc.test.cidr, 4, 1)
  gateway_ip = cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 4, 1), 1)
}

resource "huaweicloud_vpc_subnet" "standby" {
  vpc_id = huaweicloud_vpc.test.id

  name       = "%[1]s-standby"
  cidr       = cidrsubnet(huaweicloud_vpc.test.cidr, 4, 2)
  gateway_ip = cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 4, 2), 1)
}
`, rName)
}

func testAccService_basic_step1(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_service" "test" {
  access_mode = "INTERNET"
  vpc_id      = huaweicloud_vpc.test.id
  network_ids = [
    huaweicloud_vpc_subnet.master.id,
  ]
}
`, testAccService_base(rName))
}

func testAccService_basic_step2(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_service" "test" {
  access_mode = "INTERNET"
  vpc_id      = huaweicloud_vpc.test.id
  network_ids = [
    huaweicloud_vpc_subnet.master.id,
    huaweicloud_vpc_subnet.standby.id,
  ]

  enterprise_id        = "%[2]s"

  otp_config_info {
    enable       = true
    receive_mode = "VMFA"
  }
}
`, testAccService_base(rName), rName)
}

func testAccService_basic_step3(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_service" "test" {
  access_mode = "INTERNET"
  vpc_id      = huaweicloud_vpc.test.id
  network_ids = [
    huaweicloud_vpc_subnet.master.id,
    huaweicloud_vpc_subnet.standby.id,
  ]

  enterprise_id        = "%[2]s"

  otp_config_info {
    enable       = true
    receive_mode = "VMFA"
    rule_type    = "ACCESS_MODE"
    rule         = "PRIVATE"
  }
}
`, testAccService_base(rName), rName)
}

func testAccService_localAD_step1() string {
	return fmt.Sprintf(`
resource "huaweicloud_workspace_service" "test" {
  ad_domain {
    name               = try(element(regexall("\\w+\\.(.*)", element(split(",", "%[1]s"), 0))[0], 0), "")
    active_domain_name = element(split(",", "%[1]s"), 0)
    active_domain_ip   = element(split(",", "%[2]s"), 0)
    active_dns_ip      = element(split(",", "%[2]s"), 0)
    admin_account      = "%[3]s"
    password           = "%[4]s"
  }

  auth_type   = "LOCAL_AD"
  access_mode = "INTERNET"
  vpc_id      = "%[5]s"
  network_ids = ["%[6]s"]
}
`, acceptance.HW_WORKSPACE_AD_DOMAIN_NAMES,
		acceptance.HW_WORKSPACE_AD_DOMAIN_IPS,
		acceptance.HW_WORKSPACE_AD_SERVER_ACCOUNT,
		acceptance.HW_WORKSPACE_AD_SERVER_PWD,
		acceptance.HW_WORKSPACE_AD_VPC_ID,
		acceptance.HW_WORKSPACE_AD_NETWORK_ID)
}

func testAccService_localAD_step2() string {
	return fmt.Sprintf(`
resource "huaweicloud_workspace_service" "test" {
  ad_domain {
    name                = try(element(regexall("\\w+\\.(.*)", element(split(",", "%[1]s"), 0))[0], 0), "")
    active_domain_name  = element(split(",", "%[1]s"), 0)
    active_domain_ip    = element(split(",", "%[2]s"), 0)
    active_dns_ip       = element(split(",", "%[2]s"), 0)
    standby_domain_name = element(split(",", "%[1]s"), 1)
    standby_domain_ip   = element(split(",", "%[2]s"), 1)
    standby_dns_ip      = element(split(",", "%[2]s"), 1)
    admin_account       = "%[3]s"
    password            = "%[4]s"
  }

  auth_type   = "LOCAL_AD"
  access_mode = "INTERNET"
  vpc_id      = "%[5]s"
  network_ids = ["%[6]s"]

  enterprise_id = "custom-workspace-service"

  otp_config_info {
    enable       = true
    receive_mode = "VMFA"
  }
}
`, acceptance.HW_WORKSPACE_AD_DOMAIN_NAMES,
		acceptance.HW_WORKSPACE_AD_DOMAIN_IPS,
		acceptance.HW_WORKSPACE_AD_SERVER_ACCOUNT,
		acceptance.HW_WORKSPACE_AD_SERVER_PWD,
		acceptance.HW_WORKSPACE_AD_VPC_ID,
		acceptance.HW_WORKSPACE_AD_NETWORK_ID)
}

func testAccService_localAD_step3() string {
	return fmt.Sprintf(`
resource "huaweicloud_workspace_service" "test" {
  ad_domain {
    name                = try(element(regexall("\\w+\\.(.*)", element(split(",", "%[1]s"), 0))[0], 0), "")
    active_domain_name  = element(split(",", "%[1]s"), 0)
    active_domain_ip    = element(split(",", "%[2]s"), 0)
    active_dns_ip       = element(split(",", "%[2]s"), 0)
    standby_domain_name = element(split(",", "%[1]s"), 1)
    standby_domain_ip   = element(split(",", "%[2]s"), 1)
    standby_dns_ip      = element(split(",", "%[2]s"), 1)
    admin_account       = "%[3]s"
    password            = "%[4]s"
  }

  auth_type   = "LOCAL_AD"
  access_mode = "INTERNET"
  vpc_id      = "%[5]s"
  network_ids = ["%[6]s"]

  enterprise_id = "custom-workspace-service"

  otp_config_info {
    enable       = true
    receive_mode = "VMFA"
    rule_type    = "ACCESS_MODE"
    rule         = "PRIVATE"
  }
}
`, acceptance.HW_WORKSPACE_AD_DOMAIN_NAMES,
		acceptance.HW_WORKSPACE_AD_DOMAIN_IPS,
		acceptance.HW_WORKSPACE_AD_SERVER_ACCOUNT,
		acceptance.HW_WORKSPACE_AD_SERVER_PWD,
		acceptance.HW_WORKSPACE_AD_VPC_ID,
		acceptance.HW_WORKSPACE_AD_NETWORK_ID)
}

func testAccService_internetAccessPort_update(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_service" "test" {
  access_mode = "INTERNET"
  vpc_id      = huaweicloud_vpc.test.id
  network_ids = [
    huaweicloud_vpc_subnet.master.id,
    huaweicloud_vpc_subnet.standby.id,
  ]
  
  internet_access_port = "%[2]s"
  enterprise_id        = "%[3]s"
}
`, testAccService_base(rName), acceptance.HW_WORKSPACE_INTERNET_ACCESS_PORT, rName)
}

func testAccService_localAD_internetAccessPort_update() string {
	return fmt.Sprintf(`
resource "huaweicloud_workspace_service" "test" {
  ad_domain {
    name               = try(element(regexall("\\w+\\.(.*)", element(split(",", "%[1]s"), 0))[0], 0), "")
    active_domain_name = element(split(",", "%[1]s"), 0)
    active_domain_ip   = element(split(",", "%[2]s"), 0)
    active_dns_ip      = element(split(",", "%[2]s"), 0)
    admin_account      = "%[3]s"
    password           = "%[4]s"
  }

  auth_type   = "LOCAL_AD"
  access_mode = "INTERNET"
  vpc_id      = "%[5]s"
  network_ids = ["%[6]s"]

  internet_access_port = "%[7]s"
  enterprise_id        = "custom-workspace-service"
}
`, acceptance.HW_WORKSPACE_AD_DOMAIN_NAMES,
		acceptance.HW_WORKSPACE_AD_DOMAIN_IPS,
		acceptance.HW_WORKSPACE_AD_SERVER_ACCOUNT,
		acceptance.HW_WORKSPACE_AD_SERVER_PWD,
		acceptance.HW_WORKSPACE_AD_VPC_ID,
		acceptance.HW_WORKSPACE_AD_NETWORK_ID,
		acceptance.HW_WORKSPACE_INTERNET_ACCESS_PORT)
}
