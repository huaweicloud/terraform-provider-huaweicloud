package dew

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceCsmsSecrets_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	rName := "data.huaweicloud_csms_secrets.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceCsmsSecrets_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "secrets.0.id"),
					resource.TestCheckResourceAttrSet(rName, "secrets.0.name"),
					resource.TestCheckResourceAttrSet(rName, "secrets.0.status"),
					resource.TestCheckResourceAttrSet(rName, "secrets.0.kms_key_id"),
					resource.TestCheckResourceAttrSet(rName, "secrets.0.created_at"),
					resource.TestCheckResourceAttrSet(rName, "secrets.0.updated_at"),
					resource.TestCheckResourceAttrSet(rName, "secrets.0.secret_type"),
					resource.TestCheckResourceAttrSet(rName, "secrets.0.auto_rotation"),
					resource.TestCheckResourceAttrSet(rName, "secrets.0.event_subscriptions.#"),
					resource.TestCheckResourceAttrSet(rName, "secrets.0.enterprise_project_id"),

					resource.TestCheckOutput("name_filter_is_useful", "true"),

					resource.TestCheckOutput("secretId_filter_is_useful", "true"),

					resource.TestCheckOutput("status_filter_is_useful", "true"),

					resource.TestCheckOutput("enterpriseProjectId_filter_is_useful", "true"),

					resource.TestCheckOutput("eventName_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceCsmsSecrets_basic(name string) string {
	expireTime := time.Now().Add(48*time.Hour).Unix() * 1000
	return fmt.Sprintf(`
%s

data "huaweicloud_csms_secrets" "test" {
  depends_on = [huaweicloud_csms_secret.test]
}

data "huaweicloud_csms_secrets" "name_filter" {
  name = huaweicloud_csms_secret.test.name
}

locals {
  name = huaweicloud_csms_secret.test.name
}

output "name_filter_is_useful" {
  value = length(data.huaweicloud_csms_secrets.name_filter.secrets) > 0 && alltrue(
    [for v in data.huaweicloud_csms_secrets.name_filter.secrets[*].name : v == local.name]
  )
}

data "huaweicloud_csms_secrets" "secretId_filter" {
  secret_id = huaweicloud_csms_secret.test.secret_id
}

locals {
  secret_id = huaweicloud_csms_secret.test.secret_id
}

output "secretId_filter_is_useful" {
  value = length(data.huaweicloud_csms_secrets.secretId_filter.secrets) > 0 && alltrue(
    [for v in data.huaweicloud_csms_secrets.secretId_filter.secrets[*].id : v == 
  local.secret_id]
 )  
}

data "huaweicloud_csms_secrets" "status_filter" {
  status = huaweicloud_csms_secret.test.status
}

locals {
  status = huaweicloud_csms_secret.test.status
}

output "status_filter_is_useful" {
  value = length(data.huaweicloud_csms_secrets.status_filter.secrets) > 0 && alltrue(
    [for v in data.huaweicloud_csms_secrets.status_filter.secrets[*].status : v == local.status]
  )
}

data "huaweicloud_csms_secrets" "enterpriseProjectId_filter" {
  enterprise_project_id = data.huaweicloud_csms_secrets.test.secrets.0.enterprise_project_id
}

locals {
  enterprise_project_id = data.huaweicloud_csms_secrets.test.secrets.0.enterprise_project_id
}

output "enterpriseProjectId_filter_is_useful" {
  value = length(data.huaweicloud_csms_secrets.enterpriseProjectId_filter.secrets) > 0 && alltrue(
    [for v in data.huaweicloud_csms_secrets.enterpriseProjectId_filter.secrets[*].enterprise_project_id : v == local.enterprise_project_id]
  )
}

data "huaweicloud_csms_secrets" "eventName_filter" {
  event_name = data.huaweicloud_csms_secrets.test.secrets.0.event_subscriptions.0
}

locals {
  event_name = data.huaweicloud_csms_secrets.test.secrets.0.event_subscriptions.0
}

output "eventName_filter_is_useful" {
  value = length(data.huaweicloud_csms_secrets.eventName_filter.secrets) > 0 && alltrue(
    [for v in data.huaweicloud_csms_secrets.eventName_filter.secrets[*].event_subscriptions.0 : v == local.event_name]
  )
}
`, testAccDewCsmsSecret_basic(name, expireTime))
}
