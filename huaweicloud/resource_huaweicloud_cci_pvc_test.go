package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/cci/v1/persistentvolumeclaims"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccCCIPersistentVolumeClaims_basic(t *testing.T) {
	var pvc persistentvolumeclaims.ListResp
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_cci_pvc.test"
	volumeType := "ssd"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCCI(t)
			testAccPreCheckCCINamespace(t)
			testAccPreCheckEpsID(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCCIPersistentVolumeClaimsDestroy(volumeType),
		Steps: []resource.TestStep{
			{
				Config: testAccCCIPersistentVolumeClaims_basic(rName, volumeType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCIPersistentVolumeClaimsExists(resourceName, volumeType, &pvc),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "namespace", HW_CCI_NAMESPACE),
					resource.TestCheckResourceAttr(resourceName, "volume_type", volumeType),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCCIPvcImportStateIdFunc(resourceName),
			},
		},
	})
}

func TestAccCCIPersistentVolumeClaims_obs(t *testing.T) {
	var pvc persistentvolumeclaims.ListResp
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_cci_pvc.test"
	volumeType := "obs"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCCI(t)
			testAccPreCheckCCINamespace(t)
			testAccPreCheckEpsID(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCCIPersistentVolumeClaimsDestroy(volumeType),
		Steps: []resource.TestStep{
			{
				Config: testAccCCIPersistentVolumeClaims_obs(rName, volumeType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCIPersistentVolumeClaimsExists(resourceName, volumeType, &pvc),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "namespace", HW_CCI_NAMESPACE),
					resource.TestCheckResourceAttr(resourceName, "volume_type", volumeType),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCCIPvcImportStateIdFunc(resourceName),
			},
		},
	})
}

func TestAccCCIPersistentVolumeClaims_nfs(t *testing.T) {
	var pvc persistentvolumeclaims.ListResp
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_cci_pvc.test"
	volumeType := "nfs-rw"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCCI(t)
			testAccPreCheckCCINamespace(t)
			testAccPreCheckEpsID(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCCIPersistentVolumeClaimsDestroy(volumeType),
		Steps: []resource.TestStep{
			{
				Config: testAccCCIPersistentVolumeClaims_nfs(rName, volumeType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCIPersistentVolumeClaimsExists(resourceName, volumeType, &pvc),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "namespace", HW_CCI_NAMESPACE),
					resource.TestCheckResourceAttr(resourceName, "volume_type", volumeType),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCCIPvcImportStateIdFunc(resourceName),
			},
		},
	})
}

func TestAccCCIPersistentVolumeClaims_efs(t *testing.T) {
	var pvc persistentvolumeclaims.ListResp
	suffix := acctest.RandString(5)
	rName := fmt.Sprintf("tf-acc-test-%s", suffix)
	resourceName := "huaweicloud_cci_pvc.test"
	volumeType := "efs-standard"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCCI(t)
			testAccPreCheckCCINamespace(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCCIPersistentVolumeClaimsDestroy(volumeType),
		Steps: []resource.TestStep{
			{
				Config: testAccCCIPersistentVolumeClaims_efs(rName, volumeType, suffix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCIPersistentVolumeClaimsExists(resourceName, volumeType, &pvc),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "namespace", HW_CCI_NAMESPACE),
					resource.TestCheckResourceAttr(resourceName, "volume_type", volumeType),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCCIPvcImportStateIdFunc(resourceName),
			},
		},
	})
}

func testAccCheckCCIPersistentVolumeClaimsDestroy(volumeType string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*config.Config)
		client, err := config.CciV1Client(HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating HuaweiCloud CCI client: %s", err)
		}
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "huaweicloud_cci_pvc" {
				continue
			}
			response, err := getCCIPvcInfoById(client, HW_CCI_NAMESPACE, volumeType, rs.Primary.ID)
			if err == nil && response != nil {
				return fmtp.Errorf("The PVC (%s) still exist", rs.Primary.ID)
			}
		}
		return nil
	}
}

func testAccCheckCCIPersistentVolumeClaimsExists(n, volumeType string,
	pvc *persistentvolumeclaims.ListResp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}
		config := testAccProvider.Meta().(*config.Config)
		client, err := config.CciV1Client(HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating HuaweiCloud CCI Client: %s", err)
		}
		response, err := getCCIPvcInfoById(client, HW_CCI_NAMESPACE, volumeType, rs.Primary.ID)
		if err != nil {
			return fmtp.Errorf("Unable to find the specifies PVC (%s) form server: %s", rs.Primary.ID, err)
		}
		if response != nil {
			*pvc = *response
			return nil
		}

		return fmtp.Errorf("PVC (%s) not found", rs.Primary.ID)
	}
}

func testAccCCIPvcImportStateIdFunc(pvcRes string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		pvc, ok := s.RootModule().Resources[pvcRes]
		if !ok {
			return "", fmtp.Errorf("Auto Scaling lifecycle hook not found: %s", pvc)
		}
		if pvc.Primary.Attributes["volume_type"] == "" || pvc.Primary.ID == "" {
			return "", fmtp.Errorf("Unable to find the resource by import infos: %s/%s/%s",
				HW_CCI_NAMESPACE, pvc.Primary.Attributes["volume_type"], pvc.Primary.ID)
		}
		return fmt.Sprintf("%s/%s/%s", HW_CCI_NAMESPACE, pvc.Primary.Attributes["volume_type"], pvc.Primary.ID), nil
	}
}

func testAccCCIPersistentVolumeClaims_basic(rName, volumeType string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_evs_volume" "test" {
  name                  = "%s"
  description           = "Created by acc test"
  availability_zone     = data.huaweicloud_availability_zones.test.names[0]
  volume_type           = "SAS"
  size                  = 12
  enterprise_project_id = "%s"
}

resource "huaweicloud_cci_pvc" "test" {
  name        = "%s"
  namespace   = "%s"
  volume_type = "%s"
  volume_id   = huaweicloud_evs_volume.test.id
}
`, rName, HW_ENTERPRISE_PROJECT_ID_TEST, rName, HW_CCI_NAMESPACE, volumeType)
}

func testAccCCIPersistentVolumeClaims_obs(rName, volumeType string) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "bucket" {
  bucket                = "%s"
  storage_class         = "STANDARD"
  acl                   = "private"
  enterprise_project_id = "%s"
}

resource "huaweicloud_cci_pvc" "test" {
  name        = "%s"
  namespace   = "%s"
  volume_type = "%s"
  volume_id   = huaweicloud_obs_bucket.bucket.id
}
`, rName, HW_ENTERPRISE_PROJECT_ID_TEST, rName, HW_CCI_NAMESPACE, volumeType)
}

func testAccCCIPersistentVolumeClaims_nfs(rName, volumeType string) string {
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
  availability_zone = data.huaweicloud_availability_zones.myaz.names[0]
  enterprise_project_id = "%[2]s"
}

resource "huaweicloud_cci_pvc" "test" {
  name              = "%[1]s"
  namespace         = "%[3]s"
  volume_type       = "%[4]s"
  volume_id         = huaweicloud_sfs_file_system.sfs_1.id
  device_mount_path = huaweicloud_sfs_file_system.sfs_1.export_location
}
`, rName, HW_ENTERPRISE_PROJECT_ID_TEST, HW_CCI_NAMESPACE, volumeType)
}

func testAccCCIPersistentVolumeClaims_efs(rName, volumeType, suffix string) string {
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
`, rName, HW_CCI_NAMESPACE, volumeType)
}
