package modelarts

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccTrainingImageStore_basic(t *testing.T) {
	var (
		name  = acceptance.RandomAccResourceNameWithDash()
		rName = "huaweicloud_modelarts_training_image_store.test"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckDomainName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"null": {
				Source:            "hashicorp/null",
				VersionConstraint: "3.2.1",
			},
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccTrainingImageStore_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(rName, "training_job_id",
						"huaweicloud_modelarts_training_job.test", "id"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrSet(rName, "task_id"),
					resource.TestCheckResourceAttr(rName, "tag", "v1.0"),
					resource.TestCheckResourceAttr(rName, "description", "Created by terraform script"),
					resource.TestCheckResourceAttr(rName, "namespace", acceptance.HW_DOMAIN_NAME),
				),
			},
		},
	})
}

func testAccTrainingImageStore_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "test" {
  bucket        = "%[1]s"
  acl           = "private"
  force_destroy = true
}

resource "huaweicloud_obs_bucket_object" "test" {
  bucket       = huaweicloud_obs_bucket.test.bucket
  key          = "%[1]s/bootfile.py"
  content      = <<EOF
#!/usr/bin/env python
import os
print(os.getcwd())
EOF
  content_type = "text/py"
}

data "huaweicloud_modelarts_training_job_engines" "test" {}

locals {
  engine = try([for v in data.huaweicloud_modelarts_training_job_engines.test.engines : v if v.run_user != ""][0], {})
}

data "huaweicloud_modelarts_training_job_flavors" "test" {}

locals {
  flavor_id = try(
    [
      for v in data.huaweicloud_modelarts_training_job_flavors.test.flavors : v.flavor_id if
      contains(startswith(trimspace(v.support_engines), "[") ?
  jsondecode(v.support_engines) : split(",", v.support_engines), local.engine.engine_name)][0], "")
}

resource "huaweicloud_modelarts_training_job" "test" {
  kind = "job"

  metadata {
    name = "%[1]s"
  }

  algorithm {
    code_dir  = "/${huaweicloud_obs_bucket.test.bucket}/%[1]s/"
    boot_file = "/${huaweicloud_obs_bucket.test.bucket}/${huaweicloud_obs_bucket_object.test.key}"

    engine {
      engine_id      = lookup(local.engine, "engine_id", "")
      engine_version = lookup(local.engine, "engine_version", "")
      engine_name    = lookup(local.engine, "engine_name", "")
    }
  }

  spec {
    resource {
      flavor_id  = local.flavor_id
      node_count = 1
    }
  }

  depends_on = [huaweicloud_obs_bucket_object.test]
}

resource "null_resource" "wait_for_running" {
  depends_on = [huaweicloud_modelarts_training_job.test]

  provisioner "local-exec" {
    command = <<-EOT
      while true; do
        terraform apply -refresh-only -target=huaweicloud_modelarts_training_job.test -auto-approve
        status=$(terraform state show -no-color huaweicloud_modelarts_training_job.test \
          | grep -E '^\s+status\s*=' | head -1 | awk -F'= ' '{print $2}' | tr -d '"')
        [ "$status" = "Running" ] && break
        [ "$status" = "Failed" ] || [ "$status" = "Abnormal" ] && exit 1
        sleep 30
      done
    EOT
  }
}

data "huaweicloud_modelarts_training_jobs" "test" {
  filters {
    key      = "id"
    operator = "in"
    value    = [huaweicloud_modelarts_training_job.test.id]
  }
}
`, name)
}

func testAccTrainingImageStore_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_modelarts_training_image_store" "test" {
  training_job_id = huaweicloud_modelarts_training_job.test.id
  task_id         = try(data.huaweicloud_modelarts_training_jobs.test.jobs[0].status[0].tasks[0], "")
  name            = "%[2]s"
  namespace       = "%[3]s"
  tag             = "v1.0"
  description     = "Created by terraform script"

  depends_on = [null_resource.wait_for_running]
}
`, testAccTrainingImageStore_base(name), name, acceptance.HW_DOMAIN_NAME)
}
