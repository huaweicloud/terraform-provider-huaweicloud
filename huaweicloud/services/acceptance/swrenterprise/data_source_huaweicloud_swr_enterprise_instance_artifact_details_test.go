package swrenterprise

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSwrEnterpriseInstanceArtifactDetails_basic(t *testing.T) {
	dataSource := "data.huaweicloud_swr_enterprise_instance_artifact_details.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSwrSignatureWithImageEnabled(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSwrEnterpriseInstanceArtifactDetails_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "artifact_id"),
					resource.TestCheckResourceAttrSet(dataSource, "namespace_id"),
					resource.TestCheckResourceAttrSet(dataSource, "repository_id"),
					resource.TestCheckResourceAttrSet(dataSource, "media_type"),
					resource.TestCheckResourceAttrSet(dataSource, "size"),
					resource.TestCheckResourceAttrSet(dataSource, "digest"),
					resource.TestCheckResourceAttrSet(dataSource, "type"),
					resource.TestCheckResourceAttrSet(dataSource, "manifest_media_type"),
					resource.TestCheckResourceAttrSet(dataSource, "pull_time"),
					resource.TestCheckResourceAttrSet(dataSource, "push_time"),
					resource.TestCheckResourceAttrSet(dataSource, "tags.#"),
					resource.TestCheckResourceAttrSet(dataSource, "tags.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "tags.0.repository_id"),
					resource.TestCheckResourceAttrSet(dataSource, "tags.0.artifact_id"),
					resource.TestCheckResourceAttrSet(dataSource, "tags.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "tags.0.push_time"),
					resource.TestCheckResourceAttrSet(dataSource, "tags.0.pull_time"),
					resource.TestCheckResourceAttrSet(dataSource, "scan_overview.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "scan_overview.0.overview.#"),
					resource.TestCheckResourceAttrSet(dataSource, "scan_overview.0.overview.0.report_id"),
					resource.TestCheckResourceAttrSet(dataSource, "scan_overview.0.overview.0.scan_status"),
					resource.TestCheckResourceAttrSet(dataSource, "scan_overview.0.overview.0.severity"),
					resource.TestCheckResourceAttrSet(dataSource, "scan_overview.0.overview.0.duration"),
					resource.TestCheckResourceAttrSet(dataSource, "scan_overview.0.overview.0.start_time"),
					resource.TestCheckResourceAttrSet(dataSource, "scan_overview.0.overview.0.end_time"),
					resource.TestCheckResourceAttrSet(dataSource, "scan_overview.0.overview.0.complete_percent"),
					resource.TestCheckResourceAttrSet(dataSource, "scan_overview.0.overview.0.scanner.#"),
					resource.TestCheckResourceAttrSet(dataSource, "scan_overview.0.overview.0.scanner.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "scan_overview.0.overview.0.scanner.0.vendor"),
					resource.TestCheckResourceAttrSet(dataSource, "scan_overview.0.overview.0.scanner.0.version"),
					resource.TestCheckResourceAttrSet(dataSource, "scan_overview.0.overview.0.summary.#"),
					resource.TestCheckResourceAttrSet(dataSource, "scan_overview.0.overview.0.summary.0.total"),
					resource.TestCheckResourceAttrSet(dataSource, "scan_overview.0.overview.0.summary.0.fixable"),
					resource.TestCheckResourceAttrSet(dataSource, "scan_overview.0.overview.0.summary.0.summary.%"),
				),
			},
		},
	})
}

func testDataSourceSwrEnterpriseInstanceArtifactDetails_basic() string {
	return `
data "huaweicloud_swr_enterprise_instances" "test" {}

data "huaweicloud_swr_enterprise_repositories" "test" {
  instance_id = data.huaweicloud_swr_enterprise_instances.test.instances[0].id
}

data "huaweicloud_swr_enterprise_instance_artifacts" "test" {
  instance_id     = data.huaweicloud_swr_enterprise_instances.test.instances[0].id
  namespace_name  = data.huaweicloud_swr_enterprise_repositories.test.repositories[0].namespace_name
  repository_name = data.huaweicloud_swr_enterprise_repositories.test.repositories[0].name
}

data "huaweicloud_swr_enterprise_instance_artifact_details" "test" {
  instance_id        = data.huaweicloud_swr_enterprise_instances.test.instances[0].id
  namespace_name     = data.huaweicloud_swr_enterprise_repositories.test.repositories[0].namespace_name
  repository_name    = data.huaweicloud_swr_enterprise_repositories.test.repositories[0].name
  reference          = data.huaweicloud_swr_enterprise_instance_artifacts.test.artifacts[0].digest
  with_scan_overview = "true"
}
`
}
