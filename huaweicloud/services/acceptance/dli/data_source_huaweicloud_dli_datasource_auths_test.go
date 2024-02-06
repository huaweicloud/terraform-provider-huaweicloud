package dli

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceAuths_basic(t *testing.T) {
	var (
		name     = acceptance.RandomAccResourceName()
		password = acceptance.RandomPassword()
		byName   = "data.huaweicloud_dli_datasource_auths.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDatasourceAuths_basic(name, password),
				Check: resource.ComposeTestCheckFunc(
					dcByName.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(byName, "auths.#"),
					resource.TestCheckResourceAttrSet(byName, "auths.0.created_at"),
					resource.TestCheckResourceAttrSet(byName, "auths.0.updated_at"),
					resource.TestCheckResourceAttrSet(byName, "auths.0.owner"),
					resource.TestCheckResourceAttr(byName, "auths.0.type", "passwd"),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testDatasourceAuths_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_dli_datasource_auths" "filter_by_name" {
  depends_on = [
    huaweicloud_dli_datasource_auth.test
  ] 

  name = "%[1]s"
}

output "is_name_filter_useful" {
  value = length(data.huaweicloud_dli_datasource_auths.filter_by_name.auths) == 1
}
`, name)
}

func testDatasourceAuths_basic(name string, password string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dli_datasource_auth" "test" {
  name     = "%[1]s"
  type     = "passwd"
  username = "%[1]s"
  password = "%[2]s"
}

%[3]s
`, name, password, testDatasourceAuths_base(name))
}

func TestAccDatasourceAuths_CSS(t *testing.T) {
	var (
		name     = acceptance.RandomAccResourceName()
		password = acceptance.RandomPassword()
		byName   = "data.huaweicloud_dli_datasource_auths.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDliDsAuthCss(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDatasourceAuths_CSS(name, password),
				Check: resource.ComposeTestCheckFunc(
					dcByName.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(byName, "auths.#"),
					resource.TestCheckResourceAttr(byName, "auths.0.type", "CSS"),
					resource.TestCheckResourceAttr(byName, "auths.0.certificate_location", acceptance.HW_DLI_DS_AUTH_CSS_OBS_PATH),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testDatasourceAuths_CSS(name string, password string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dli_datasource_auth" "test" {
  name                 = "%[1]s"
  type                 = "CSS"
  username             = "%[1]s"
  password             = "%[2]s"
  certificate_location = "%[3]s"
}

%[4]s
`, name, password, acceptance.HW_DLI_DS_AUTH_CSS_OBS_PATH, testDatasourceAuths_base(name))
}

func TestAccDatasourceAuths_Kafka_SSL(t *testing.T) {
	var (
		name     = acceptance.RandomAccResourceName()
		password = acceptance.RandomPassword()
		byName   = "data.huaweicloud_dli_datasource_auths.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDliDsAuthKafka(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDatasourceAuths_Kafka_SSL(name, password),
				Check: resource.ComposeTestCheckFunc(
					dcByName.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(byName, "auths.#"),
					resource.TestCheckResourceAttr(byName, "auths.0.type", "Kafka_SSL"),
					resource.TestCheckResourceAttr(byName, "auths.0.truststore_location", acceptance.HW_DLI_DS_AUTH_KAFKA_TRUST_OBS_PATH),
					resource.TestCheckResourceAttr(byName, "auths.0.keystore_location", acceptance.HW_DLI_DS_AUTH_KAFKA_KEY_OBS_PATH),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testDatasourceAuths_Kafka_SSL(name string, password string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dli_datasource_auth" "test" {
  name                = "%[1]s"
  type                = "Kafka_SSL"
  truststore_location = "%[3]s"
  truststore_password = "%[2]s"
  keystore_location   = "%[4]s"
  keystore_password   = "%[2]s"
  key_password        = "%[2]s"
}

%[5]s
`, name, password, acceptance.HW_DLI_DS_AUTH_KAFKA_TRUST_OBS_PATH, acceptance.HW_DLI_DS_AUTH_KAFKA_KEY_OBS_PATH, testDatasourceAuths_base(name))
}

func TestAccDatasourceAuths_KRB(t *testing.T) {
	var (
		name     = acceptance.RandomAccResourceName()
		byName   = "data.huaweicloud_dli_datasource_auths.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDliDsAuthKrb(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDatasourceAuths_KRB(name),
				Check: resource.ComposeTestCheckFunc(
					dcByName.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(byName, "auths.#"),
					resource.TestCheckResourceAttr(byName, "auths.0.type", "KRB"),
					resource.TestCheckResourceAttr(byName, "auths.0.krb5_conf", acceptance.HW_DLI_DS_AUTH_KRB_CONF_OBS_PATH),
					resource.TestCheckResourceAttr(byName, "auths.0.keytab", acceptance.HW_DLI_DS_AUTH_KRB_TAB_OBS_PATH),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testDatasourceAuths_KRB(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dli_datasource_auth" "test" {
  name      = "%[1]s"
  type      = "KRB"
  username  = "%[1]s"
  krb5_conf = "%[2]s"
  keytab    = "%[3]s"
}

%[4]s
`, name, acceptance.HW_DLI_DS_AUTH_KRB_CONF_OBS_PATH, acceptance.HW_DLI_DS_AUTH_KRB_TAB_OBS_PATH, testDatasourceAuths_base(name))
}
