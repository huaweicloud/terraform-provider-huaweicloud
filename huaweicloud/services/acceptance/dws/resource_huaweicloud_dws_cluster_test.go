package dws

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dws"
)

func getClusterResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getDwsCluster: Query the DWS cluster.
	getDwsClusterClient, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return nil, fmt.Errorf("error creating DWS client: %s", err)
	}

	getDwsClusterRespBody, err := dws.GetClusterInfoByClusterId(getDwsClusterClient, state.Primary.ID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DWS cluster: %s", err)
	}

	return getDwsClusterRespBody, nil
}

func TestAccResourceCluster_basicV1(t *testing.T) {
	var obj interface{}

	resourceName := "huaweicloud_dws_cluster.test"
	name := acceptance.RandomAccResourceName()
	password := "TF" + acceptance.RandomPassword()
	updatePassword := "TF" + acceptance.RandomPassword()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getClusterResourceFunc,
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
				Config: testAccDwsCluster_basic_step1(name, 3, dws.PublicBindTypeAuto, password),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "number_of_node", "3"),
					resource.TestCheckResourceAttr(resourceName, "logical_cluster_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "public_ip.0.public_bind_type", dws.PublicBindTypeAuto),
					resource.TestCheckResourceAttrSet(resourceName, "public_ip.0.eip_id"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "val"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttrPair(resourceName, "security_group_id", "huaweicloud_networking_secgroup.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "description", "Created v1 cluster by terraform script"),
				),
			},
			{
				Config: testAccDwsCluster_basic_step2(name, updatePassword),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "number_of_node", "6"),
					resource.TestCheckResourceAttr(resourceName, "logical_cluster_enable", "false"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_MIGRATE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "public_ip.0.public_bind_type", dws.PublicBindTypeBindExisting),
					resource.TestCheckResourceAttrPair(resourceName, "public_ip.0.eip_id", "huaweicloud_vpc_eip.test.0", "id"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "val"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar1"),
					resource.TestCheckResourceAttrPair(resourceName, "security_group_id", "huaweicloud_networking_secgroup.test2", "id"),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated cluster info"),
				),
			},
			{
				Config: testAccDwsCluster_basic_step3(name, updatePassword),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "number_of_node", "3"),
					resource.TestCheckResourceAttr(resourceName, "public_ip.0.public_bind_type", dws.PublicBindTypeNotUse),
					resource.TestCheckResourceAttr(resourceName, "public_ip.0.eip_id", ""),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
			// Assert whether the restart is successful.
			{
				Config: testAccDwsCluster_basic_step4(name, updatePassword),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"user_pwd", "number_of_cn", "volume", "endpoints", "logical_cluster_enable", "force_backup"},
			},
		},
	})
}

func testAccDwsCluster_eip(name string, count int) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc_eip" "test" {
  count = %d

  publicip {
    type = "5_bgp"
  }

  bandwidth {
    name        = "%s_${count.index}"
    size        = 1
    share_type  = "PER"
    charge_mode = "bandwidth"
  }
}
`, count, name)
}

func testAccDwsCluster_basic_step1(rName string, numberOfNode int, publicIpBindType, password string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_dws_cluster" "test" {
  name                   = "%[2]s"
  node_type              = "dwsk2.xlarge"
  number_of_node         = %[3]d
  vpc_id                 = huaweicloud_vpc.test.id
  network_id             = huaweicloud_vpc_subnet.test.id
  security_group_id      = huaweicloud_networking_secgroup.test.id
  availability_zone      = data.huaweicloud_availability_zones.test.names[0]
  user_name              = "test_cluster_admin"
  user_pwd               = "%[4]s"
  logical_cluster_enable = true
  enterprise_project_id  = "%[5]s"
  description            = "Created v1 cluster by terraform script"

  public_ip {
    # Automatically purchase EIP when creating cluster.
    public_bind_type = "%[6]s"
  }

  tags = {
    key = "val"
    foo = "bar"
  }
}
`, common.TestBaseNetwork(rName), rName, numberOfNode, password, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST, publicIpBindType)
}

func testAccDwsCluster_basic_step2(rName, password string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "huaweicloud_networking_secgroup" "test2" {
  name                 = "%[3]s_test2"
  delete_default_rules = true
}

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_dws_cluster" "test" {
  name                   = "%[3]s"
  node_type              = "dwsk2.xlarge"
  number_of_node         = 6
  vpc_id                 = huaweicloud_vpc.test.id
  network_id             = huaweicloud_vpc_subnet.test.id
  security_group_id      = huaweicloud_networking_secgroup.test2.id
  availability_zone      = data.huaweicloud_availability_zones.test.names[0]
  user_name              = "test_cluster_admin"
  user_pwd               = "%[4]s"
  logical_cluster_enable = false
  enterprise_project_id  = "%[5]s"
  description            = "Updated cluster info"

  public_ip {
    # Modify the associated EIP.
    public_bind_type = "bind_existing"
    eip_id           = huaweicloud_vpc_eip.test[0].id
  }

  tags = {
    key = "val"
    foo = "bar1"
  }
}
`, common.TestBaseNetwork(rName), testAccDwsCluster_eip(rName, 1), rName, password,
		acceptance.HW_ENTERPRISE_MIGRATE_PROJECT_ID_TEST)
}

func testAccDwsCluster_basic_step3(rName, password string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "huaweicloud_networking_secgroup" "test2" {
  name                 = "%[3]s_test2"
  delete_default_rules = true
}

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_dws_cluster" "test" {
  name                   = "%[3]s"
  node_type              = "dwsk2.xlarge"
  number_of_node         = 3
  vpc_id                 = huaweicloud_vpc.test.id
  network_id             = huaweicloud_vpc_subnet.test.id
  security_group_id      = huaweicloud_networking_secgroup.test2.id
  availability_zone      = data.huaweicloud_availability_zones.test.names[0]
  user_name              = "test_cluster_admin"
  user_pwd               = "%[4]s"
  logical_cluster_enable = false
  enterprise_project_id  = "%[5]s"

  public_ip {
    # Unbind the associated EIP.
    public_bind_type = "not_use"
    eip_id           = ""
  }
}
`, common.TestBaseNetwork(rName), testAccDwsCluster_eip(rName, 1), rName, password,
		acceptance.HW_ENTERPRISE_MIGRATE_PROJECT_ID_TEST)
}

func testAccDwsCluster_basic_step4(name, password string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dws_cluster_restart" "test" {
  cluster_id = huaweicloud_dws_cluster.test.id
}
`, testAccDwsCluster_basic_step3(name, password))
}

func TestAccResourceCluster_basicV2(t *testing.T) {
	var obj interface{}

	resourceName := "huaweicloud_dws_cluster.testv2"
	name := acceptance.RandomAccResourceName()
	// The password must be between 12 and 32 characters.
	password := "TF" + acceptance.RandomPassword()
	updatePassword := "TF" + acceptance.RandomPassword()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getClusterResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMigrateEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			// Assert basic configuration.
			{
				Config: testAccDwsCluster_basicV2_step1(name, password),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "number_of_node", "3"),
					resource.TestCheckResourceAttr(resourceName, "public_ip.0.public_bind_type", dws.PublicBindTypeBindExisting),
					resource.TestCheckResourceAttrPair(resourceName, "public_ip.0.eip_id", "huaweicloud_vpc_eip.test.0", "id"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "val"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.capacity", "100"),
					resource.TestCheckResourceAttrSet(resourceName, "version"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "elb.0.name", name+"_elb0"),
					resource.TestCheckResourceAttrPair(resourceName, "security_group_id", "huaweicloud_networking_secgroup.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "description", "Created v2 cluster by terraform script"),
				),
			},
			// Assert all modifiable parameters.
			{
				Config: testAccDwsCluster_basicV2_step2(name, updatePassword),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "number_of_node", "6"),
					resource.TestCheckResourceAttrPair(resourceName, "public_ip.0.eip_id", "huaweicloud_vpc_eip.test.1", "id"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "val"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "cat"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.capacity", "150"),
					resource.TestCheckResourceAttrSet(resourceName, "version"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_MIGRATE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "elb.0.name", name+"_elb1"),
					resource.TestCheckResourceAttrPair(resourceName, "security_group_id", "huaweicloud_networking_secgroup.test2", "id"),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated cluster info"),
				),
			},
			// Assert that ELB and EIP are unbound and delete all tags.
			{
				Config: testAccDwsCluster_basicV2_step3(name, updatePassword),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "number_of_node", "3"),
					resource.TestCheckResourceAttr(resourceName, "public_ip.0.public_bind_type", dws.PublicBindTypeNotUse),
					resource.TestCheckResourceAttr(resourceName, "public_ip.0.eip_id", ""),
					resource.TestCheckResourceAttr(resourceName, "elb.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
				),
			},
			// Assert whether the restart is successful.
			{
				Config: testAccDwsCluster_basicV2_step4(name, updatePassword),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"user_pwd", "number_of_cn", "volume", "endpoints", "lts_enable",
					"logical_cluster_enable", "elb_id", "force_backup"},
			},
		},
	})
}

func testAccDwsCluster_basicV2_base(rName string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "huaweicloud_networking_secgroup" "test2" {
  name                 = "%[3]s_test2"
  delete_default_rules = true
}

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_dws_flavors" "test" {
  vcpus          = 4
  memory         = 32
  datastore_type = "dws"
}

resource "huaweicloud_elb_loadbalancer" "test" {
  count          = 2
  name           = "%[3]s_elb${count.index}"
  vpc_id         = huaweicloud_vpc.test.id
  ipv4_subnet_id = huaweicloud_vpc_subnet.test.ipv4_subnet_id

  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0]
  ]

  backend_subnets = [
    huaweicloud_vpc_subnet.test.id
  ]

  protection_status = "nonProtection"
}
`, common.TestBaseNetwork(rName), testAccDwsCluster_eip(rName, 2), rName)
}

func testAccDwsCluster_basicV2_step1(rName, password string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dws_cluster" "testv2" {
  name                   = "%[2]s"
  node_type              = "dwsk2.xlarge"
  number_of_node         = 3
  vpc_id                 = huaweicloud_vpc.test.id
  network_id             = huaweicloud_vpc_subnet.test.id
  security_group_id      = huaweicloud_networking_secgroup.test.id
  availability_zone      = data.huaweicloud_availability_zones.test.names[0]
  user_name              = "test_cluster_admin"
  user_pwd               = "%[3]s"
  version                = data.huaweicloud_dws_flavors.test.flavors[0].datastore_version
  number_of_cn           = 3
  lts_enable             = true
  logical_cluster_enable = false
  enterprise_project_id  = "%[4]s"
  elb_id                 = huaweicloud_elb_loadbalancer.test[0].id
  description            = "Created v2 cluster by terraform script"

  public_ip {
    # Binging a EIP for cluster.
    public_bind_type = "bind_existing"
    eip_id           = huaweicloud_vpc_eip.test[0].id
  }

  volume {
    type     = "SSD"
    capacity = 100
  }

  tags = {
    key = "val"
    foo = "bar"
  }
}
`, testAccDwsCluster_basicV2_base(rName), rName, password, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccDwsCluster_basicV2_step2(rName, password string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dws_cluster" "testv2" {
  name                   = "%[2]s"
  node_type              = "dwsk2.xlarge"
  number_of_node         = 6
  vpc_id                 = huaweicloud_vpc.test.id
  network_id             = huaweicloud_vpc_subnet.test.id
  security_group_id      = huaweicloud_networking_secgroup.test2.id
  availability_zone      = data.huaweicloud_availability_zones.test.names[0]
  user_name              = "test_cluster_admin"
  user_pwd               = "%[3]s"
  version                = data.huaweicloud_dws_flavors.test.flavors[0].datastore_version
  number_of_cn           = 3
  lts_enable             = false
  logical_cluster_enable = false
  enterprise_project_id  = "%[4]s"
  elb_id                 = huaweicloud_elb_loadbalancer.test[1].id
  description            = "Updated cluster info"

  public_ip {
    # Modify the associated EIP.
    eip_id = huaweicloud_vpc_eip.test[1].id
  }

  volume {
    type     = "SSD"
    capacity = 150
  }

  tags = {
    key = "val"
    foo = "cat"
  }
}
`, testAccDwsCluster_basicV2_base(rName), rName, password, acceptance.HW_ENTERPRISE_MIGRATE_PROJECT_ID_TEST)
}

func testAccDwsCluster_basicV2_step3(rName, password string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dws_cluster" "testv2" {
  name                   = "%[2]s"
  node_type              = "dwsk2.xlarge"
  number_of_node         = 3
  force_backup           = false
  vpc_id                 = huaweicloud_vpc.test.id
  network_id             = huaweicloud_vpc_subnet.test.id
  security_group_id      = huaweicloud_networking_secgroup.test2.id
  availability_zone      = data.huaweicloud_availability_zones.test.names[0]
  user_name              = "test_cluster_admin"
  user_pwd               = "%[3]s"
  version                = data.huaweicloud_dws_flavors.test.flavors[0].datastore_version
  number_of_cn           = 3
  lts_enable             = false
  logical_cluster_enable = true
  enterprise_project_id  = "%[4]s"

  public_ip {
    # Unbind the associated EIP.
    eip_id = ""
  }

  volume {
    type     = "SSD"
    capacity = 150
  }
}
`, testAccDwsCluster_basicV2_base(rName), rName, password, acceptance.HW_ENTERPRISE_MIGRATE_PROJECT_ID_TEST)
}

func testAccDwsCluster_basicV2_step4(name, password string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dws_cluster_restart" "test" {
  cluster_id = huaweicloud_dws_cluster.testv2.id
}
`, testAccDwsCluster_basicV2_step3(name, password))
}

// Test the scenarios with multiple AZs and volumes is local disk.
func TestAccResourceCluster_basicV2_mutilAZs(t *testing.T) {
	var obj interface{}

	resourceName := "huaweicloud_dws_cluster.test"
	name := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getClusterResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMutilAZ(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDwsCluster_basicV2_mutilAZs(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "number_of_node", "3"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "val"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "availability_zone", acceptance.HW_DWS_MUTIL_AZS),
					resource.TestCheckResourceAttrSet(resourceName, "version"),
					resource.TestCheckResourceAttr(resourceName, "public_ip.0.public_bind_type", dws.PublicBindTypeNotUse),
					resource.TestCheckResourceAttr(resourceName, "public_ip.0.eip_id", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// After the resource is created successfully, the "updated" attribute value is not refreshed immediately,
				// resulting in inconsistencies between the actual value and the expected value when importing the resource.
				ImportStateVerifyIgnore: []string{"user_pwd", "number_of_cn", "volume", "endpoints", "updated"},
			},
		},
	})
}

func testAccDwsCluster_basicV2_mutilAZs(rName string) string {
	baseNetwork := common.TestBaseNetwork(rName)
	password := "TF" + acceptance.RandomPassword()
	return fmt.Sprintf(`
%s

data "huaweicloud_dws_flavors" "test" {
  vcpus          = 4
  memory         = 16
  datastore_type = "dws"
}

resource "huaweicloud_dws_cluster" "test" {
  name = "%s"
  // The specification of the local disk.
  node_type         = "dws2.olap.4xlarge.i3"
  number_of_node    = 3
  vpc_id            = huaweicloud_vpc.test.id
  network_id        = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  availability_zone = "%s"
  user_name         = "test_cluster_admin"
  user_pwd          = "%s"
  version           = data.huaweicloud_dws_flavors.test.flavors[0].datastore_version
  number_of_cn      = 3

  tags = {
    key = "val"
    foo = "bar"
  }
}
`, baseNetwork, rName, acceptance.HW_DWS_MUTIL_AZS, password)
}
