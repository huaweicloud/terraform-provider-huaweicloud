package coc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/coc"
)

func getScriptBatchExecuteResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("coc", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating COC client: %s", err)
	}

	return coc.GetExecutionTicketDetail(client, state.Primary.ID)
}

func TestAccScriptBatchExecute_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		obj   interface{}
		rName = "huaweicloud_coc_script_batch_execute.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getScriptBatchExecuteResourceFunc)

		rNameWithNotSync = "huaweicloud_coc_script_batch_execute.with_not_sync"
		rcWithNotSync    = acceptance.InitResourceCheck(rNameWithNotSync, &obj, getScriptBatchExecuteResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckProjectID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: "0.12.1",
			},
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccScriptBatchExecute_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "script_id", "huaweicloud_coc_script.test", "id"),
					resource.TestCheckResourceAttr(rName, "execute_batches.#", "2"),
					resource.TestCheckResourceAttr(rName, "execute_batches.0.batch_index", "1"),
					resource.TestCheckResourceAttr(rName, "execute_batches.0.instance_ids.#", "1"),
					resource.TestCheckResourceAttr(rName, "execute_batches.1.batch_index", "2"),
					resource.TestCheckResourceAttr(rName, "execute_batches.1.instance_ids.#", "2"),
					resource.TestCheckResourceAttr(rName, "timeout", "600"),
					resource.TestCheckResourceAttr(rName, "execute_user", "root"),
					resource.TestCheckResourceAttr(rName, "parameters.#", "2"),
					resource.TestCheckResourceAttr(rName, "script_name", name),
					resource.TestCheckResourceAttr(rName, "is_sync", "true"),
					resource.TestCheckResourceAttr(rName, "status", "FINISHED"),
					resource.TestMatchResourceAttr(rName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(rName, "finished_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					rcWithNotSync.CheckResourceExists(),
					resource.TestCheckResourceAttr(rNameWithNotSync, "is_sync", "false"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"execute_batches", "parameters", "is_sync",
				},
			},
			{
				ResourceName:      rNameWithNotSync,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"execute_batches", "parameters", "is_sync",
				},
			},
		},
	})
}

func testAccScriptBatchExecute_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = try(data.huaweicloud_availability_zones.test.names[0], "")
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

data "huaweicloud_images_images" "test" {
  flavor_id  = try(data.huaweicloud_compute_flavors.test.ids[0], "")
  os         = "Ubuntu"
  visibility = "public"
}

# The default security group rules cannot be deleted, otherwise the UniAgent installation will fail.
resource "huaweicloud_networking_secgroup" "test" {
  name = "%[2]s"
}

# Create two ECS instances and install UniAgent
resource "huaweicloud_compute_instance" "test" {
  count = 3

  name                  = "%[2]s${count.index}"
  availability_zone     = try(data.huaweicloud_availability_zones.test.names[0], "")
  flavor_id             = try(data.huaweicloud_compute_flavors.test.flavors[0].id, "NOT_FOUND")
  image_id              = try(data.huaweicloud_images_images.test.images[0].id, "NOT_FOUND")
  security_group_ids    = [huaweicloud_networking_secgroup.test.id]
  enterprise_project_id = "%[3]s"

  user_data = <<EOF
#! /bin/bash
set +o history;
curl -k -X GET -m 20 --retry 1 --retry-delay 10 -o /tmp/install_uniagent https://aom-uniagent-%[4]s.obs.%[4]s.myhuaweicloud.com/install_uniagent.sh; \
bash /tmp/install_uniagent -p %[5]s -v 1.2.0 -e %[4]s;set -o history;
EOF

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }
}

# Wait for the ECS instances UniAgent to be install completed.
resource "time_sleep" "test" {
  create_duration = "1m"

  depends_on = [huaweicloud_compute_instance.test]
}

resource "huaweicloud_coc_script" "test" {
  name                  = "%[2]s"
  risk_level            = "LOW"
  version               = "1.0.1"
  type                  = "SHELL"
  enterprise_project_id = "%[3]s"
  description           = "Created by Terraform script"

  content = <<EOF
#! /bin/bash
echo "hello $${name}@$${organization}!"
EOF

  parameters {
    name        = "name"
    value       = "world"
    description = "the first parameter"
  }
  parameters {
    name        = "organization"
    value       = "Huawei"
    description = "the second parameter"
    sensitive   = true
  }
}`, common.TestVpc(name, acceptance.HW_ENTERPRISE_PROJECT_ID), name,
		acceptance.HW_ENTERPRISE_PROJECT_ID,
		acceptance.HW_REGION_NAME,
		acceptance.HW_PROJECT_ID,
	)
}

func testAccScriptBatchExecute_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_coc_script_batch_execute" "test" {
  script_id    = huaweicloud_coc_script.test.id
  timeout      = 600
  execute_user = "root"

  execute_batches {
    batch_index  = 1
    instance_ids = slice(huaweicloud_compute_instance.test[*].id, 0, 1)
  }
  execute_batches {
    batch_index  = 2
    instance_ids = slice(huaweicloud_compute_instance.test[*].id, 1, 3)
  }

  parameters {
    name  = "name"
    value = "TF acceptance test"
  }
  parameters {
    name  = "organization"
    value = "HuaweiCloud"
  }

  depends_on = [time_sleep.test]
}

resource "huaweicloud_coc_script_batch_execute" "with_not_sync" {
  script_id    = huaweicloud_coc_script.test.id
  timeout      = 600
  execute_user = "root"
  is_sync      = false

  execute_batches {
    batch_index  = 1
    instance_ids = huaweicloud_compute_instance.test[*].id
  }

  parameters {
    name  = "name"
    value = "TF"
  }
  parameters {
    name  = "organization"
    value = "HW TF"
  }

  depends_on = [time_sleep.test]
}
`, testAccScriptBatchExecute_base(name))
}
