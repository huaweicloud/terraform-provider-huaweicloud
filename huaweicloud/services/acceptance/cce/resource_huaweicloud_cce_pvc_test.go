package cce

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cce"

	"github.com/chnsz/golangsdk/openstack/cce/v1/persistentvolumeclaims"
)

func getPvcResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.CceV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating HuaweiCloud CCE v1 client: %s", err)
	}
	resp, err := cce.GetCcePvcInfoById(c, state.Primary.Attributes["cluster_id"],
		state.Primary.Attributes["namespace"], state.Primary.ID)
	if resp == nil && err == nil {
		return resp, fmt.Errorf("Unable to find the persistent volume claim (%s)", state.Primary.ID)
	}
	return resp, err
}

func TestAccCcePersistentVolumeClaimsV1_basic(t *testing.T) {
	var pvc persistentvolumeclaims.PersistentVolumeClaim
	resourceName := "huaweicloud_cce_pvc.test"
	randName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	rc := acceptance.InitResourceCheck(
		resourceName,
		&pvc,
		getPvcResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCcePersistentVolumeClaimsV1_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "cluster_id",
						"${huaweicloud_cce_cluster.test.id}"),
					resource.TestCheckResourceAttr(resourceName, "namespace", "default"),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "volume_type", cce.EvsVolume),
					resource.TestCheckResourceAttr(resourceName, "access_mode", "ReadWriteOnce"),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "volume_id",
						"${huaweicloud_evs_volume.test.id}"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCcePvcImportStateFunc(resourceName),
			},
		},
	})
}

func TestAccCcePersistentVolumeClaimsV1_obs(t *testing.T) {
	var pvc persistentvolumeclaims.PersistentVolumeClaim
	resourceName := "huaweicloud_cce_pvc.test"
	randName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	rc := acceptance.InitResourceCheck(
		resourceName,
		&pvc,
		getPvcResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCcePersistentVolumeClaimsV1_obs(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "cluster_id",
						"${huaweicloud_cce_cluster.test.id}"),
					resource.TestCheckResourceAttr(resourceName, "namespace", "default"),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "volume_type", cce.ObsVolume),
					resource.TestCheckResourceAttr(resourceName, "access_mode", "ReadWriteMany"),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "volume_id",
						"${huaweicloud_obs_bucket.test.id}"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCcePvcImportStateFunc(resourceName),
			},
		},
	})
}

func TestAccCcePersistentVolumeClaimsV1_sfs(t *testing.T) {
	var pvc persistentvolumeclaims.PersistentVolumeClaim
	resourceName := "huaweicloud_cce_pvc.test"
	randName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	rc := acceptance.InitResourceCheck(
		resourceName,
		&pvc,
		getPvcResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCcePersistentVolumeClaimsV1_sfs(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "cluster_id",
						"${huaweicloud_cce_cluster.test.id}"),
					resource.TestCheckResourceAttr(resourceName, "namespace", "default"),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "volume_type", cce.SfsVolume),
					resource.TestCheckResourceAttr(resourceName, "access_mode", "ReadWriteMany"),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "volume_id",
						"${huaweicloud_sfs_file_system.test.id}"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCcePvcImportStateFunc(resourceName),
			},
		},
	})
}

func testAccCcePvcImportStateFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("Resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.ID == "" || rs.Primary.Attributes["cluster_id"] == "" ||
			rs.Primary.Attributes["namespace"] == "" {
			return "", fmt.Errorf("resource not found: %s/%s/%s", rs.Primary.Attributes["cluster_id"],
				rs.Primary.Attributes["namespace"], rs.Primary.ID)
		}
		return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["cluster_id"], rs.Primary.Attributes["namespace"],
			rs.Primary.ID), nil
	}
}

func testAccCceCluster_config(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%s"
  cidr = "192.168.0.0/20"
}

resource "huaweicloud_vpc_subnet" "test" {
  name       = "%s"
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
  vpc_id     = huaweicloud_vpc.test.id
}

resource "huaweicloud_cce_cluster" "test" {
  name                   = "%s"
  flavor_id              = "cce.s1.small"
  vpc_id                 = huaweicloud_vpc.test.id
  subnet_id              = huaweicloud_vpc_subnet.test.id
  container_network_type = "overlay_l2"
}`, rName, rName, rName)
}

func testAccCcePersistentVolumeClaimsV1_basic(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

%s

resource "huaweicloud_evs_volume" "test" {
  name              = "%s"
  volume_type       = "SSD"
  size              = 40
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
}

resource "huaweicloud_cce_pvc" "test" {
  cluster_id  = huaweicloud_cce_cluster.test.id
  namespace   = "default"
  name        = "%s"
  volume_id   = huaweicloud_evs_volume.test.id
  volume_type = "bs"
}
`, testAccCceCluster_config(rName), rName, rName)
}

func testAccCcePersistentVolumeClaimsV1_obs(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_obs_bucket" "test" {
  bucket                = "%s"
  storage_class         = "STANDARD"
  acl                   = "private"
  enterprise_project_id = "%s"
}

resource "huaweicloud_cce_pvc" "test" {
  cluster_id  = huaweicloud_cce_cluster.test.id
  namespace   = "default"
  name        = "%s"
  volume_id   = huaweicloud_obs_bucket.test.id
  volume_type = "obs"
}
`, testAccCceCluster_config(rName), rName, acceptance.HW_ENTERPRISE_PROJECT_ID, rName)
}

func testAccCcePersistentVolumeClaimsV1_sfs(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_sfs_file_system" "test" {
  size      = 10
  name      = "%s"
  access_to = huaweicloud_vpc.test.id
}

resource "huaweicloud_cce_pvc" "test" {
  cluster_id  = huaweicloud_cce_cluster.test.id
  namespace   = "default"
  name        = "%s"
  volume_id   = huaweicloud_sfs_file_system.test.id
  volume_type = "nfs"
}
`, testAccCceCluster_config(rName), rName, rName)
}
