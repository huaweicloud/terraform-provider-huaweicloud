package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccHSSContainerExportTask_basic(t *testing.T) {
	rName := "huaweicloud_hss_container_export_task.test"

	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testContainerExportTask_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(rName, "task_id"),
					resource.TestCheckResourceAttrSet(rName, "task_name"),
					resource.TestCheckResourceAttrSet(rName, "task_status"),
					resource.TestCheckResourceAttrSet(rName, "file_id"),
					resource.TestCheckResourceAttrSet(rName, "file_name"),
				),
			},
		},
	})
}

const testContainerExportTask_basic = `
resource "huaweicloud_hss_container_export_task" "test" {
  export_headers = [
    ["container_name", "Container Name"],
    ["cluster_name", "Cluster Name"],
    ["status", "Status"]
  ]

  cluster_container     = "true"
  cluster_type          = "cce"
  cluster_name          = "test-cluster-name"
  container_name        = "test-container-name"
  pod_name              = "test-pod-name"
  image_name            = "test-image-name"
  status                = "Running"
  risky                 = "true"
  cpu_limit             = "100m"
  memory_limit          = "300Mi"
  enterprise_project_id = "0"
  export_size           = 100

  create_time {
    start_time = 1690348800
    end_time   = 1690435200
  }
}
`
