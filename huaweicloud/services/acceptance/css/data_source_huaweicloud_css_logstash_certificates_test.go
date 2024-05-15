package css

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCssLogstashCertificates_basic(t *testing.T) {
	dataSource := "data.huaweicloud_css_logstash_certificates.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	tmpFile, err := os.CreateTemp("", "tf-css-cert-test.cer")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	// write some data to the tempfile
	err = os.WriteFile(tmpFile.Name(), []byte("initial only test"), 0600)
	if err != nil {
		t.Fatal(err)
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCssLogstashCertificates_basic(rName, tmpFile.Name()),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "certificates.0.file_name"),
					resource.TestCheckResourceAttrSet(dataSource, "certificates.0.status"),

					resource.TestCheckOutput("file_name_filter_is_useful", "true"),
					resource.TestCheckOutput("status_filter_is_useful", "true"),
					resource.TestCheckOutput("certs_type_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceCssLogstashCertificates_basic(name, source string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_css_logstash_certificates" "test" {
  depends_on = [huaweicloud_css_logstash_custom_certificate.test]

  cluster_id = huaweicloud_css_logstash_cluster.test.id
}

locals {
  file_name = data.huaweicloud_css_logstash_certificates.test.certificates[0].file_name
  status    = data.huaweicloud_css_logstash_certificates.test.certificates[0].status
}

data "huaweicloud_css_logstash_certificates" "filter_by_file_name" {
  cluster_id = huaweicloud_css_logstash_cluster.test.id
  file_name  = local.file_name
}

data "huaweicloud_css_logstash_certificates" "filter_by_status" {
  cluster_id = huaweicloud_css_logstash_cluster.test.id
  status     = local.status
}

data "huaweicloud_css_logstash_certificates" "filter_by_certs_type" {
  cluster_id = huaweicloud_css_logstash_cluster.test.id
  certs_type = "defaultCerts"
}

locals {
  list_by_file_name  = data.huaweicloud_css_logstash_certificates.filter_by_file_name.certificates
  list_by_status     = data.huaweicloud_css_logstash_certificates.filter_by_status.certificates
  list_by_certs_type = data.huaweicloud_css_logstash_certificates.filter_by_certs_type.certificates
}

output "file_name_filter_is_useful" {
  value = length(local.list_by_file_name) > 0 && alltrue(
    [for v in local.list_by_file_name[*].file_name : v == local.file_name]
  )
}

output "status_filter_is_useful" {
  value = length(local.list_by_status) > 0 && alltrue(
    [for v in local.list_by_status[*].status : v == local.status]
  )
}

output "certs_type_filter_is_useful" {
  value = length(local.list_by_certs_type) == 1
}  
`, testLogstashCertificate_basic(name, source))
}
