package drs

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceDrsJobV5_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceDrsJobV5_basic(name),
				ExpectError: regexp.MustCompile(`error waiting for DRS job`),
			},
		},
	})
}

// The parameters used in this test case are all dummy data and are only used to test failure scenarios.
func testAccResourceDrsJobV5_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_drs_job_v5" "test" {
  base_info {
    name                  = "%s"
    job_type              = "sync"
    engine_type           = "mysql-to-mysql"
    job_direction         = "up"
    task_type             = "FULL_INCR_TRANS"
    net_type              = "eip"
    charging_mode         = "on_demand"
    enterprise_project_id = "0"
    expired_days          = "14"
    is_open_fast_clean    = false

    tags {
      key   = "tag1"
      value = "value1"
    }
  }

  source_endpoint {
    db_type       = "mysql"
    endpoint_type = "cloud"
    endpoint_role = "so"

    endpoint {
      endpoint_name = "cloud_mysql"
      ip            = "192.168.0.141"
      db_port       = "3306"
      db_user       = "root"
      db_password   = "!@#eedd11"
      instance_id   = "3335c09cffc541d48b999b412531b27bin01"
      db_name       = "user"
    }

    cloud {
      region     = "cn-north-4"
      project_id = "0970dd7a1300f5673ff2c003c60ae111"
    }

    ssl {
      ssl_link = false
    }
  }

  target_endpoint {
    db_type       = "mysql"
    endpoint_type = "cloud"
    endpoint_role = "ta"

    endpoint {
      endpoint_name = "cloud_mysql"
      ip            = "192.168.0.105"
      db_port       = "3306"
      db_user       = "root"
      db_password   = "!@#eedd11"
      instance_id   = "dc939e01662142df975895a5a353b73fin01"
    }

    cloud {
      region     = "cn-north-4"
      project_id = "0970dd7a1300f5673ff2c003c60ae111"
      az_code    = "cn-north-4a"
    }

    vpc {
      vpc_id            = "0135b4e1-48f2-4981-8eba-eb802525ba31"
      subnet_id         = "8451d3ad-8014-477b-a735-4edf0cf407fa"
      security_group_id = "d3d3d1f2-c533-4132-9c2e-8457152413b5"
    }
  }

  node_info {
    spec {
      node_type = "medium"
    }

    vpc {
      vpc_id    = "0135b4e1-48f2-4981-8eba-eb802525ba31"
      subnet_id = "8451d3ad-8014-477b-a735-4edf0cf407fa"
    }
  }
}
`, name)
}
