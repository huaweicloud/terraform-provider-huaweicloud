package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk/openstack/cci/v1/persistentvolumeclaims"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccCCIPersistentVolumeClaims_basic(t *testing.T) {
	var pvc persistentvolumeclaims.ListResp
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_cci_pvc.test"
	ns := "terraform-test"
	volumeType := "ssd"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCCI(t)
			testAccPreCheckEpsID(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCCIPersistentVolumeClaimsDestroy(ns, volumeType),
		Steps: []resource.TestStep{
			{
				Config: testAccCCIPersistentVolumeClaims_basic(ns, rName, volumeType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCIPersistentVolumeClaimsExists(resourceName, ns, volumeType, &pvc),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "namespace", ns),
					resource.TestCheckResourceAttr(resourceName, "volume_type", volumeType),
					resource.TestCheckResourceAttr(resourceName, "volume_size", "20"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCCIPvcImportStateIdFunc(resourceName, ns),
			},
		},
	})
}

func TestAccCCIPersistentVolumeClaims_obs(t *testing.T) {
	var pvc persistentvolumeclaims.ListResp
	rInt := acctest.RandInt()
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_cci_pvc.test"
	ns := "terraform-test"
	volumeType := "obs"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCCI(t)
			testAccPreCheckEpsID(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCCIPersistentVolumeClaimsDestroy(ns, volumeType),
		Steps: []resource.TestStep{
			{
				Config: testAccCCIPersistentVolumeClaims_obs(ns, rName, volumeType, rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCIPersistentVolumeClaimsExists(resourceName, ns, volumeType, &pvc),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "namespace", ns),
					resource.TestCheckResourceAttr(resourceName, "volume_type", volumeType),
					resource.TestCheckResourceAttr(resourceName, "volume_size", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCCIPvcImportStateIdFunc(resourceName, ns),
			},
		},
	})
}

func TestAccCCIPersistentVolumeClaims_nfs(t *testing.T) {
	var pvc persistentvolumeclaims.ListResp
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_cci_pvc.test"
	ns := "terraform-test"
	volumeType := "nfs-rw"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCCI(t)
			testAccPreCheckEpsID(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCCIPersistentVolumeClaimsDestroy(ns, volumeType),
		Steps: []resource.TestStep{
			{
				Config: testAccCCIPersistentVolumeClaims_nfs(ns, rName, volumeType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCIPersistentVolumeClaimsExists(resourceName, ns, volumeType, &pvc),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "namespace", ns),
					resource.TestCheckResourceAttr(resourceName, "volume_type", volumeType),
					resource.TestCheckResourceAttr(resourceName, "volume_size", "100"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCCIPvcImportStateIdFunc(resourceName, ns),
			},
		},
	})
}

func TestAccCCIPersistentVolumeClaims_efs(t *testing.T) {
	var pvc persistentvolumeclaims.ListResp
	suffix := acctest.RandString(5)
	rName := fmt.Sprintf("tf-acc-test-%s", suffix)
	resourceName := "huaweicloud_cci_pvc.test"
	ns := "terraform-test"
	volumeType := "efs-standard"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCCI(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCCIPersistentVolumeClaimsDestroy(ns, volumeType),
		Steps: []resource.TestStep{
			{
				Config: testAccCCIPersistentVolumeClaims_efs(ns, rName, suffix, volumeType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCIPersistentVolumeClaimsExists(resourceName, ns, volumeType, &pvc),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "namespace", ns),
					resource.TestCheckResourceAttr(resourceName, "volume_type", volumeType),
					resource.TestCheckResourceAttr(resourceName, "volume_size", "500"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCCIPvcImportStateIdFunc(resourceName, ns),
			},
		},
	})
}

func testAccCheckCCIPersistentVolumeClaimsDestroy(ns, volumeType string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*config.Config)
		client, err := config.CciV1Client(HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud CCI client: %s", err)
		}
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "huaweicloud_cci_pvc" {
				continue
			}
			_, err := getPvcInfoFromServer(client, ns, volumeType, rs.Primary.ID)
			if err != nil {
				return fmt.Errorf("Unable to find the specifies PVC (%s) form server: %s", rs.Primary.ID, err)
			}
		}
		return nil
	}
}

func testAccCheckCCIPersistentVolumeClaimsExists(n, ns, volumeType string,
	pvc *persistentvolumeclaims.ListResp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		config := testAccProvider.Meta().(*config.Config)
		client, err := config.CciV1Client(HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud CCI Client: %s", err)
		}
		response, err := getPvcInfoFromServer(client, ns, volumeType, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Unable to find the specifies PVC (%s) form server: %s", rs.Primary.ID, err)
		}
		if response != nil {
			*pvc = *response
			return nil
		}

		return fmt.Errorf("PVC (%s) not found", rs.Primary.ID)
	}
}

func testAccCCIPvcImportStateIdFunc(pvcRes, ns string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		pvc, ok := s.RootModule().Resources[pvcRes]
		if !ok {
			return "", fmt.Errorf("Auto Scaling lifecycle hook not found: %s", pvc)
		}
		if ns == "" || pvc.Primary.Attributes["volume_type"] == "" || pvc.Primary.ID == "" {
			return "", fmt.Errorf("Unable to find the resource by import ID: %s/%s/%s",
				ns, pvc.Primary.Attributes["volume_type"], pvc.Primary.ID)
		}
		return fmt.Sprintf("%s/%s/%s", ns, pvc.Primary.Attributes["volume_type"], pvc.Primary.ID), nil
	}
}

func testAccCCIPersistentVolumeClaims_basic(ns, rName, volumeType string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cci_pvc" "test" {
  name        = "%s"
  namespace   = "%s"
  volume_type = "%s"
  volume_id   = huaweicloud_evs_volume.test.id
  volume_size = huaweicloud_evs_volume.test.size
}
`, testAccEvsVolume_epsID(rName), rName, ns, volumeType)
}

func testAccCCIPersistentVolumeClaims_obs(ns, rName, volumeType string, suffix int) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cci_pvc" "test" {
  name        = "%s"
  namespace   = "%s"
  volume_type = "%s"
  volume_id   = huaweicloud_obs_bucket.bucket.id
  volume_size = 1
}
`, testAccObsBucket_epsId(suffix), rName, ns, volumeType)
}

func testAccCCIPersistentVolumeClaims_nfs(ns, rName, volumeType string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cci_pvc" "test" {
  name              = "%s"
  namespace         = "%s"
  volume_type       = "%s"
  volume_id         = huaweicloud_sfs_file_system.sfs_1.id
  volume_size       = 100
  device_mount_path = huaweicloud_sfs_file_system.sfs_1.export_location
}
`, testAccSFSFileSystemV2_epsId(rName), rName, ns, volumeType)
}

func testAccCCIPersistentVolumeClaims_efs(ns, rName, suffix, volumeType string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cci_pvc" "test" {
  name              = "%s"
  namespace         = "%s"
  volume_type       = "%s"
  volume_id         = huaweicloud_sfs_turbo.sfs-turbo1.id
  volume_size       = 500
  device_mount_path = huaweicloud_sfs_turbo.sfs-turbo1.export_location
}
`, testAccSFSTurbo_basic(suffix), rName, ns, volumeType)
}
