package swrenterprise

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSwrEnterpriseInstanceArtifactVulnerabilities_basic(t *testing.T) {
	dataSource := "data.huaweicloud_swr_enterprise_instance_artifact_vulnerabilities.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSwrSignatureWithImageEnabled(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSwrEnterpriseInstanceArtifactVulnerabilities_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "reports.#"),
					resource.TestCheckResourceAttrSet(dataSource, "reports.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "reports.0.content.#"),
					resource.TestCheckResourceAttrSet(dataSource, "reports.0.content.0.generated_at"),
					resource.TestCheckResourceAttrSet(dataSource, "reports.0.content.0.severity"),
					resource.TestCheckResourceAttrSet(dataSource, "reports.0.content.0.scanner.#"),
					resource.TestCheckResourceAttrSet(dataSource, "reports.0.content.0.scanner.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "reports.0.content.0.scanner.0.vendor"),
					resource.TestCheckResourceAttrSet(dataSource, "reports.0.content.0.scanner.0.version"),
					resource.TestCheckResourceAttrSet(dataSource, "reports.0.content.0.vulnerabilities.#"),
					resource.TestCheckResourceAttrSet(dataSource, "reports.0.content.0.vulnerabilities.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "reports.0.content.0.vulnerabilities.0.package"),
					resource.TestCheckResourceAttrSet(dataSource, "reports.0.content.0.vulnerabilities.0.version"),
					resource.TestCheckResourceAttrSet(dataSource, "reports.0.content.0.vulnerabilities.0.severity"),
					resource.TestCheckResourceAttrSet(dataSource, "reports.0.content.0.vulnerabilities.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "reports.0.content.0.vulnerabilities.0.links.#"),
					resource.TestCheckResourceAttrSet(dataSource, "reports.0.content.0.vulnerabilities.0.artifact_digests.#"),
					resource.TestCheckResourceAttrSet(dataSource, "reports.0.content.0.vulnerabilities.0.preferred_cvss.#"),
					resource.TestCheckResourceAttrSet(dataSource, "reports.0.content.0.vulnerabilities.0.preferred_cvss.0.score_v3"),
				),
			},
		},
	})
}

func testDataSourceSwrEnterpriseInstanceArtifactVulnerabilities_basic() string {
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

data "huaweicloud_swr_enterprise_instance_artifact_vulnerabilities" "test" {
  instance_id     = data.huaweicloud_swr_enterprise_instances.test.instances[0].id
  namespace_name  = data.huaweicloud_swr_enterprise_repositories.test.repositories[0].namespace_name
  repository_name = data.huaweicloud_swr_enterprise_repositories.test.repositories[0].name
  reference       = data.huaweicloud_swr_enterprise_instance_artifacts.test.artifacts[0].digest
}
`
}
