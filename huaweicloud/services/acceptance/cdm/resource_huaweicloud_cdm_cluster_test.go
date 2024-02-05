package cdm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/cdm/v1/clusters"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func getCdmClusterResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.CdmV11Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CDM v1 client, err=%s", err)
	}
	return clusters.Get(client, state.Primary.ID)
}

func TestAccResourceCdmCluster_basic(t *testing.T) {
	var obj clusters.ClusterCreateOpts
	resourceName := "huaweicloud_cdm_cluster.test"
	name := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getCdmClusterResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCdmCluster_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "status", "Normal"),
					resource.TestCheckResourceAttrSet(resourceName, "version"),
					resource.TestCheckResourceAttrSet(resourceName, "created"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCdmCluster_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_cdm_flavors" "test" {}

resource "huaweicloud_cdm_cluster" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  flavor_id         = data.huaweicloud_cdm_flavors.test.flavors[0].id
  name              = "%s"
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  vpc_id            = huaweicloud_vpc.test.id
}
`, common.TestBaseNetwork(name), name)
}

func TestAccResourceCdmCluster_all(t *testing.T) {
	var obj clusters.ClusterCreateOpts
	resourceName := "huaweicloud_cdm_cluster.test"
	name := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getCdmClusterResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCdmCluster_all(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "is_auto_off", "true"),
					resource.TestCheckResourceAttr(resourceName, "status", "Normal"),
					resource.TestCheckResourceAttrSet(resourceName, "created"),
				),
			},
			{
				Config: testAccCdmCluster_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "is_auto_off", "false"),
					resource.TestCheckResourceAttr(resourceName, "status", "Normal"),
					resource.TestCheckResourceAttrSet(resourceName, "created"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"email", "phone_num"},
			},
		},
	})
}

func testAccCdmCluster_all(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_cdm_flavors" "test" {}

resource "huaweicloud_cdm_cluster" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  flavor_id         = data.huaweicloud_cdm_flavors.test.flavors[0].id
  name              = "%s"
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  vpc_id            = huaweicloud_vpc.test.id
  is_auto_off       = true
  email             = ["test@test.com"]
  phone_num         = ["12345678910"]
}
`, common.TestBaseNetwork(name), name)
}

func testAccCdmCluster_update(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_cdm_flavors" "test" {}

resource "huaweicloud_cdm_cluster" "test" {
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  flavor_id          = data.huaweicloud_cdm_flavors.test.flavors[0].id
  name               = "%s"
  security_group_id  = huaweicloud_networking_secgroup.test.id
  subnet_id          = huaweicloud_vpc_subnet.test.id
  vpc_id             = huaweicloud_vpc.test.id
  email              = ["test@test.com"]
  phone_num          = ["12345678910"]
  schedule_boot_time = "00:00:00"
  schedule_off_time  = "10:00:00"
}
`, common.TestBaseNetwork(name), name)
}

func TestAccResourceCdmCluster_updateWithEpsId(t *testing.T) {
	var obj clusters.ClusterCreateOpts
	resourceName := "huaweicloud_cdm_cluster.test"
	name := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getCdmClusterResourceFunc,
	)
	srcEPS := acceptance.HW_ENTERPRISE_PROJECT_ID_TEST
	destEPS := acceptance.HW_ENTERPRISE_MIGRATE_PROJECT_ID_TEST

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMigrateEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCdmCluster_withEpsId(name, srcEPS),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", srcEPS),
				),
			},
			{
				Config: testAccCdmCluster_withEpsId(name, destEPS),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", destEPS),
				),
			},
		},
	})
}

func testAccCdmCluster_withEpsId(name, epsId string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_cdm_flavors" "test" {}

resource "huaweicloud_cdm_cluster" "test" {
  availability_zone      = data.huaweicloud_availability_zones.test.names[0]
  flavor_id              = data.huaweicloud_cdm_flavors.test.flavors[0].id
  name                   = "%s"
  security_group_id      = huaweicloud_networking_secgroup.test.id
  subnet_id              = huaweicloud_vpc_subnet.test.id
  vpc_id                 = huaweicloud_vpc.test.id
  enterprise_project_id  = "%s"
}
`, common.TestBaseNetwork(name), name, epsId)
}
