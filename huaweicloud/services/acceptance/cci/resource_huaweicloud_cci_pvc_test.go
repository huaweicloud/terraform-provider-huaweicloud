package cci

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/cci/v1/persistentvolumeclaims"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cci"
)

func getPvcResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.CciV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating HuaweiCloud CCI v1 client: %s", err)
	}
	return cci.GetPvcInfoById(c, acceptance.HW_CCI_NAMESPACE, state.Primary.Attributes["volume_type"], state.Primary.ID)
}

func TestAccPersistentVolumeClaims_basic(t *testing.T) {
	var pvc persistentvolumeclaims.PersistentVolumeClaim
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_cci_pvc.test"
	volumeType := "ssd"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&pvc,
		getPvcResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCCINamespace(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccPersistentVolumeClaims_basic(rName, volumeType),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "namespace", acceptance.HW_CCI_NAMESPACE),
					resource.TestCheckResourceAttr(resourceName, "volume_type", volumeType),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccPvcImportStateIdFunc(resourceName),
			},
		},
	})
}

func TestAccPersistentVolumeClaims_obs(t *testing.T) {
	var pvc persistentvolumeclaims.PersistentVolumeClaim
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_cci_pvc.test"
	volumeType := "obs"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&pvc,
		getPvcResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCCINamespace(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccPersistentVolumeClaims_obs(rName, volumeType),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "namespace", acceptance.HW_CCI_NAMESPACE),
					resource.TestCheckResourceAttr(resourceName, "volume_type", volumeType),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccPvcImportStateIdFunc(resourceName),
			},
		},
	})
}

func TestAccPersistentVolumeClaims_nfs(t *testing.T) {
	var pvc persistentvolumeclaims.PersistentVolumeClaim
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_cci_pvc.test"
	volumeType := "nfs-rw"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&pvc,
		getPvcResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCCINamespace(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccPersistentVolumeClaims_nfs(rName, volumeType),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "namespace", acceptance.HW_CCI_NAMESPACE),
					resource.TestCheckResourceAttr(resourceName, "volume_type", volumeType),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccPvcImportStateIdFunc(resourceName),
			},
		},
	})
}

func TestAccPersistentVolumeClaims_efs(t *testing.T) {
	var pvc persistentvolumeclaims.PersistentVolumeClaim
	suffix := acctest.RandString(5)
	rName := fmt.Sprintf("tf-acc-test-%s", suffix)
	resourceName := "huaweicloud_cci_pvc.test"
	volumeType := "efs-standard"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&pvc,
		getPvcResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCCINamespace(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccPersistentVolumeClaims_efs(rName, volumeType, suffix),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "namespace", acceptance.HW_CCI_NAMESPACE),
					resource.TestCheckResourceAttr(resourceName, "volume_type", volumeType),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccPvcImportStateIdFunc(resourceName),
			},
		},
	})
}

func testAccPvcImportStateIdFunc(pvcRes string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		pvc, ok := s.RootModule().Resources[pvcRes]
		if !ok {
			return "", fmt.Errorf("auto Scaling lifecycle hook not found: %s", pvc)
		}
		namespace := acceptance.HW_CCI_NAMESPACE
		volumeType := pvc.Primary.Attributes["volume_type"]
		pvcId := pvc.Primary.ID

		if volumeType == "" || pvcId == "" {
			return "", fmt.Errorf("unable to find the resource by import infos: %s/%s/%s",
				namespace, volumeType, pvcId)
		}
		return fmt.Sprintf("%s/%s/%s", namespace, volumeType, pvcId), nil
	}
}

func testAccPersistentVolumeClaims_basic(rName, volumeType string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_evs_volume" "test" {
  name                  = "%[1]s"
  description           = "Created by acc test"
  availability_zone     = data.huaweicloud_availability_zones.test.names[0]
  volume_type           = "SAS"
  size                  = 12
  enterprise_project_id = "%[2]s"
}

resource "huaweicloud_cci_pvc" "test" {
  name        = "%[1]s"
  namespace   = "%[3]s"
  volume_type = "%[4]s"
  volume_id   = huaweicloud_evs_volume.test.id
}
`, rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST, acceptance.HW_CCI_NAMESPACE, volumeType)
}

func testAccPersistentVolumeClaims_obs(rName, volumeType string) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "bucket" {
  bucket                = "%[1]s"
  storage_class         = "STANDARD"
  acl                   = "private"
  enterprise_project_id = "%[2]s"
}

resource "huaweicloud_cci_pvc" "test" {
  name        = "%[1]s"
  namespace   = "%[3]s"
  volume_type = "%[4]s"
  volume_id   = huaweicloud_obs_bucket.bucket.id
}
`, rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST, acceptance.HW_CCI_NAMESPACE, volumeType)
}

func testAccPersistentVolumeClaims_nfs(rName, volumeType string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

data "huaweicloud_availability_zones" "myaz" {}

resource "huaweicloud_sfs_file_system" "sfs_1" {
  share_proto  = "NFS"
  size         = 10
  name         = "%[1]s"
  description  = "sfs_c2c_test-file"
  access_to    = huaweicloud_vpc.test.id
  access_type  = "cert"
  access_level = "rw"

  availability_zone     = data.huaweicloud_availability_zones.myaz.names[0]
  enterprise_project_id = "%[2]s"
}

resource "huaweicloud_cci_pvc" "test" {
  name              = "%[1]s"
  namespace         = "%[3]s"
  volume_type       = "%[4]s"
  volume_id         = huaweicloud_sfs_file_system.sfs_1.id
  device_mount_path = huaweicloud_sfs_file_system.sfs_1.export_location
}
`, rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST, acceptance.HW_CCI_NAMESPACE, volumeType)
}

func testAccPersistentVolumeClaims_efs(rName, volumeType, _ string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

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

resource "huaweicloud_sfs_turbo" "test" {
  name              = "%[1]s"
  size              = 500
  share_proto       = "NFS"
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
}

resource "huaweicloud_cci_pvc" "test" {
  name              = "%[1]s"
  namespace         = "%[2]s"
  volume_type       = "%[3]s"
  volume_id         = huaweicloud_sfs_turbo.test.id
  device_mount_path = huaweicloud_sfs_turbo.test.export_location
}
`, rName, acceptance.HW_CCI_NAMESPACE, volumeType)
}
