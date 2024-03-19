package dli

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getDatasourceAuthResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getDatasourceAuth: Query the DLI datasource authentication.
	var (
		getDatasourceAuthHttpUrl = "v3/{project_id}/datasource/auth-infos"
		getDatasourceAuthProduct = "dli"
	)
	getDatasourceAuthClient, err := cfg.NewServiceClient(getDatasourceAuthProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DLI Client: %s", err)
	}

	getDatasourceAuthPath := getDatasourceAuthClient.Endpoint + getDatasourceAuthHttpUrl
	getDatasourceAuthPath = strings.ReplaceAll(getDatasourceAuthPath, "{project_id}", getDatasourceAuthClient.ProjectID)

	getDatasourceAuthPath += "?auth_info_name=" + state.Primary.ID

	getDatasourceAuthOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getDatasourceAuthResp, err := getDatasourceAuthClient.Request("GET", getDatasourceAuthPath, &getDatasourceAuthOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DatasourceAuth: %s", err)
	}
	getDatasourceAuthRespBody, err := utils.FlattenResponse(getDatasourceAuthResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DatasourceAuth: %s", err)
	}
	v := utils.PathSearch("auth_infos[0]", getDatasourceAuthRespBody, nil)
	if v == nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return v, nil
}

func TestAccDatasourceAuth_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_dli_datasource_auth.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDatasourceAuthResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDatasourceAuth_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "type", "passwd"),
					resource.TestCheckResourceAttr(rName, "user_name", "test"),
					resource.TestCheckResourceAttrSet(rName, "owner"),
				),
			},
			{
				Config: testDatasourceAuth_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "user_name", "test123"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password",
					"truststore_password",
					"keystore_password",
					"key_password",
				},
			},
		},
	})
}

func testDatasourceAuth_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dli_datasource_auth" "test" {
  name      = "%s"
  type      = "passwd"
  user_name = "test"
  password  = "Huawei12!"
}
`, name)
}

func testDatasourceAuth_basic_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dli_datasource_auth" "test" {
  name      = "%s"
  type      = "passwd"
  user_name = "test123"
  password  = "Huawei12!"
}
`, name)
}

func TestAccDatasourceAuth_css(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_dli_datasource_auth.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDatasourceAuthResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDliDsAuthCss(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDatasourceAuth_css(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "type", "CSS"),
					resource.TestCheckResourceAttr(rName, "user_name", "test"),
					resource.TestCheckResourceAttrSet(rName, "owner"),
				),
			},
			{
				Config: testDatasourceAuth_css_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "user_name", "test_update"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password",
					"truststore_password",
					"keystore_password",
					"key_password",
				},
			},
		},
	})
}

func testDatasourceAuth_css(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dli_datasource_auth" "test" {
  name                 = "%s"
  type                 = "CSS"
  user_name            = "test"
  password             = "Huawei12!"
  certificate_location = "%s"
}
`, name, acceptance.HW_DLI_DS_AUTH_CSS_OBS_PATH)
}

func testDatasourceAuth_css_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dli_datasource_auth" "test" {
  name                 = "%s"
  type                 = "CSS"
  user_name            = "test_update"
  password             = "Huawei12!"
  certificate_location = "%s"
}
`, name, acceptance.HW_DLI_DS_AUTH_CSS_OBS_PATH)
}

func TestAccDatasourceAuth_Kafka_SSL(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_dli_datasource_auth.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDatasourceAuthResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDliDsAuthKafka(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDatasourceAuth_Kafka_SSL(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "type", "Kafka_SSL"),
					resource.TestCheckResourceAttrSet(rName, "owner"),
				),
			},
			{
				Config: testDatasourceAuth_Kafka_SSL_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password",
					"truststore_password",
					"keystore_password",
					"key_password",
				},
			},
		},
	})
}

func testDatasourceAuth_Kafka_SSL(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dli_datasource_auth" "test" {
  name                = "%s"
  type                = "Kafka_SSL"
  truststore_location = "%s"
  truststore_password = "Huawei12!"
  keystore_location   = "%s"
  keystore_password   = "Huawei12!"
  key_password        = "Huawei12!"
}
`, name, acceptance.HW_DLI_DS_AUTH_KAFKA_TRUST_OBS_PATH, acceptance.HW_DLI_DS_AUTH_KAFKA_KEY_OBS_PATH)
}

func testDatasourceAuth_Kafka_SSL_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dli_datasource_auth" "test" {
  name                = "%s"
  type                = "Kafka_SSL"
  truststore_location = "%s"
  truststore_password = "Huawei123!"
  keystore_location   = "%s"
  keystore_password   = "Huawei123!"
  key_password        = "Huawei123!"
}
`, name, acceptance.HW_DLI_DS_AUTH_KAFKA_TRUST_OBS_PATH, acceptance.HW_DLI_DS_AUTH_KAFKA_KEY_OBS_PATH)
}

func TestAccDatasourceAuth_KRB(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_dli_datasource_auth.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDatasourceAuthResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDliDsAuthKrb(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDatasourceAuth_KRB(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "type", "KRB"),
					resource.TestCheckResourceAttrSet(rName, "owner"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password",
					"truststore_password",
					"keystore_password",
					"key_password",
				},
			},
		},
	})
}

func testDatasourceAuth_KRB(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dli_datasource_auth" "test" {
  name       = "%s"
  type       = "KRB"
  user_name  = "test"
  krb5_conf  = "%s"
  keytab     = "%s"
}
`, name, acceptance.HW_DLI_DS_AUTH_KRB_CONF_OBS_PATH, acceptance.HW_DLI_DS_AUTH_KRB_TAB_OBS_PATH)
}
