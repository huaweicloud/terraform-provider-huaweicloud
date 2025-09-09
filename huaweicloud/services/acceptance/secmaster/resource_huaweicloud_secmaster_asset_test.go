package secmaster

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/secmaster"
)

func getAssetResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("secmaster", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating SecMaster client: %s", err)
	}

	workspaceID := state.Primary.Attributes["workspace_id"]
	assetID := state.Primary.ID

	return secmaster.ReadAssetDetail(client, workspaceID, assetID)
}

func TestAccAsset_basic(t *testing.T) {
	var obj interface{}
	rName := "huaweicloud_secmaster_asset.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getAssetResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
			acceptance.TestAccPreCheckSecMasterRDSAssetID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAsset_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workspace_id", acceptance.HW_SECMASTER_WORKSPACE_ID),
					resource.TestCheckResourceAttr(rName, "asset_id", acceptance.HW_SECMASTER_RDS_ASSET_ID),
					resource.TestCheckResourceAttr(rName, "data_object.0.name", "tf_test_zoj7e"),
					resource.TestCheckResourceAttr(rName, "data_object.0.department.0.id", "test-department-id"),
					resource.TestCheckResourceAttr(rName, "data_object.0.environment.0.idc_id", "test-idc-id"),
					resource.TestCheckResourceAttr(rName, "data_object.0.environment.0.idc_name", "test-idc-name"),
					resource.TestCheckResourceAttr(rName, "data_object.0.governance_user.0.name", "llrds"),
					resource.TestCheckResourceAttr(rName, "data_object.0.governance_user.0.type", "test-governance-user-type"),
					resource.TestCheckResourceAttr(rName, "data_object.0.properties.0.hwc_rds.0.name", "tf_test_zoj7e"),
					resource.TestCheckResourceAttr(rName, "data_object.0.properties.0.hwc_rds.0.port", "3306"),
				),
			},
			{
				Config: testAsset_basic_update(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "data_object.0.name", "tf_test_zoj7e_update"),
					resource.TestCheckResourceAttr(rName, "data_object.0.properties.0.hwc_rds.0.alias", "test-alias-update"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccAssetImportStateIdFunc(rName),
			},
		},
	})
}

func testAsset_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_secmaster_asset" "test" {
  workspace_id = "%[1]s"
  asset_id     = "%[2]s"

  data_object {
    checksum           = "testchecksum"
    created            = "2025-08-25T00:30:54.820Z+0800"
    id                 = "%[2]s"
    level              = 1
    name               = "tf_test_zoj7e"
    provider           = "rds"
    provisioning_state = "temp-state"
    type               = "instances"

    department {
      id   = "test-department-id"
      name = "test-department-name"
    }

    environment {
      domain_id   = "ono-exist-domain-id"
      ep_id       = "0"
      ep_name     = "default"
      idc_id      = "test-idc-id"
      idc_name    = "test-idc-name"
      project_id  = "ono-exist-project-id"
      region_id   = "ono-exist-region-id"
      vendor_name = "test-vendor-name"
      vendor_type = "CloudService"
    }

    governance_user {
      name = "llrds"
      type = "test-governance-user-type"
    }

    properties {
      hwc_rds {
        alias                 = "test-alias"
        associated_with_ddm   = false
        backup_used_space     = 0
        cpu                   = "2"
        created               = "2025-08-24T15:33:15+0000"
        db_user_name          = "root"
        disk_encryption_id    = "test-disk-encryption-id"
        enable_ssl            = false
        enterprise_project_id = "ono-exist-enterprise-project-id"
        expiration_time       = "2025-08-24T15:33:15Z"
        flavor_ref            = "rds.mysql.x1.large.2"
        id                    = "ono-exist-id"
        maintenance_window    = "18:00-22:00"
        max_iops              = 0
        mem                   = "4"
        name                  = "tf_test_zoj7e"
        port                  = 3306
        private_dns_names = [
          "ono-exist-private-dns-name",
        ]
        private_ips = [
          "ono-exist-private-ips",
        ]
        project_id       = "ono-exist-project-id"
        protected_status = "CLOSE"
        public_ips = [
          "ono-exist-public-ips",
        ]
        read_only_by_user  = false
        region             = "ono-exist-region-id"
        security_group_id  = "ono-exist-security-group-id"
        status             = "ACTIVE"
        storage_used_space = 0
        subnet_id          = "ono-exist-subnet-id"
        switch_strategy    = "reliability"
        time_zone          = "UTC"
        type               = "Single"
        updated            = "2025-08-24T15:50:01+0000"
        vpc_id             = "ono-exist-vpc-id"

        backup_strategy {
          keep_days  = 7
          start_time = "02:00-03:00"
        }

        datastore {
          complete_version = "8.0.28.231003"
          type             = "MySQL"
          version          = "8.0"
        }

        ha {
          replication_mode = "async"
        }

        nodes {
          availability_zone = "ono-exist-availability-zone"
          id                = "ono-exist-id"
          name              = "ono-exist-name"
          role              = "master"
          status            = "ACTIVE"
        }

        related_instance {
          id   = "ono-exist-instance-id"
          type = "replica_of"
        }

        tags {
          key = "test-key"
          values = [
            "test-value1",
            "test-value2",
          ]
        }

        volume {
          size = 40
          type = "CLOUDSSD"
        }
      }
    }
  }
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, acceptance.HW_SECMASTER_RDS_ASSET_ID)
}

func testAsset_basic_update() string {
	return fmt.Sprintf(`
resource "huaweicloud_secmaster_asset" "test" {
  workspace_id = "%[1]s"
  asset_id     = "%[2]s"

  data_object {
    checksum           = "testchecksum"
    created            = "2025-08-25T00:30:54.820Z+0800"
    id                 = "%[2]s"
    level              = 1
    name               = "tf_test_zoj7e_update"
    provider           = "rds"
    provisioning_state = "temp-state"
    type               = "instances"

    department {
      id   = "test-department-id"
      name = "test-department-name"
    }

    environment {
      domain_id   = "ono-exist-domain-id"
      ep_id       = "0"
      ep_name     = "default"
      idc_id      = "test-idc-id"
      idc_name    = "test-idc-name"
      project_id  = "ono-exist-project-id"
      region_id   = "ono-exist-region-id"
      vendor_name = "test-vendor-name"
      vendor_type = "CloudService"
    }

    governance_user {
      name = "llrds"
      type = "test-governance-user-type"
    }

    properties {
      hwc_rds {
        alias                 = "test-alias-update"
        associated_with_ddm   = false
        backup_used_space     = 0
        cpu                   = "2"
        created               = "2025-08-24T15:33:15+0000"
        db_user_name          = "root"
        disk_encryption_id    = "test-disk-encryption-id"
        enable_ssl            = false
        enterprise_project_id = "ono-exist-enterprise-project-id"
        expiration_time       = "2025-08-24T15:33:15Z"
        flavor_ref            = "rds.mysql.x1.large.2"
        id                    = "ono-exist-id"
        maintenance_window    = "18:00-22:00"
        max_iops              = 0
        mem                   = "4"
        name                  = "tf_test_zoj7e"
        port                  = 3306
        private_dns_names = [
          "ono-exist-private-dns-name",
        ]
        private_ips = [
          "ono-exist-private-ips",
        ]
        project_id       = "ono-exist-project-id"
        protected_status = "CLOSE"
        public_ips = [
          "ono-exist-public-ips",
        ]
        read_only_by_user  = false
        region             = "ono-exist-region-id"
        security_group_id  = "ono-exist-security-group-id"
        status             = "ACTIVE"
        storage_used_space = 0
        subnet_id          = "ono-exist-subnet-id"
        switch_strategy    = "reliability"
        time_zone          = "UTC"
        type               = "Single"
        updated            = "2025-08-24T15:50:01+0000"
        vpc_id             = "ono-exist-vpc-id"

        backup_strategy {
          keep_days  = 7
          start_time = "02:00-03:00"
        }

        datastore {
          complete_version = "8.0.28.231003"
          type             = "MySQL"
          version          = "8.0"
        }

        ha {
          replication_mode = "async"
        }

        nodes {
          availability_zone = "ono-exist-availability-zone"
          id                = "ono-exist-id"
          name              = "ono-exist-name"
          role              = "master"
          status            = "ACTIVE"
        }

        related_instance {
          id   = "ono-exist-instance-id"
          type = "replica_of"
        }

        tags {
          key = "test-key"
          values = [
            "test-value1",
            "test-value2",
          ]
        }

        volume {
          size = 40
          type = "CLOUDSSD"
        }
      }
    }
  }
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, acceptance.HW_SECMASTER_RDS_ASSET_ID)
}

func testAccAssetImportStateIdFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", name)
		}
		if rs.Primary.ID == "" || rs.Primary.Attributes["workspace_id"] == "" {
			return "", errors.New("invalid format specified for import ID, must be <workspace_id>/<id>")
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["workspace_id"], rs.Primary.ID), nil
	}
}
