package apig

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/certificates"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func getCertificateFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.ApigV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG v2 client: %s", err)
	}
	return certificates.Get(client, state.Primary.ID)
}

func TestAccCertificate_basic(t *testing.T) {
	var (
		cert certificates.Certificate

		rName      = "huaweicloud_apig_certificate.test"
		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&cert,
		getCertificateFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCertificateWithoutRootCA(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCertificate_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "type", "global"),
					resource.TestCheckResourceAttr(rName, "instance_id", "common"),
					resource.TestMatchResourceAttr(rName, "effected_at", regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`)),
					resource.TestMatchResourceAttr(rName, "expires_at", regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`)),
					resource.TestMatchResourceAttr(rName, "sans.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
				),
			},
			{
				Config: testAccCertificate_basic_step2(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "type", "global"),
					resource.TestCheckResourceAttr(rName, "instance_id", "common"),
					resource.TestMatchResourceAttr(rName, "effected_at", regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`)),
					resource.TestMatchResourceAttr(rName, "expires_at", regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`)),
					resource.TestMatchResourceAttr(rName, "sans.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"content", "private_key",
				},
			},
		},
	})
}

func testAccCertificate_basic_step1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_apig_certificate" "test" {
  name        = "%[1]s"
  content     = "%[2]s"
  private_key = "%[3]s"
}
`, name, acceptance.HW_CERTIFICATE_CONTENT, acceptance.HW_CERTIFICATE_PRIVATE_KEY)
}

func testAccCertificate_basic_step2(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_apig_certificate" "test" {
  name        = "%[1]s"
  content     = "%[2]s"
  private_key = "%[3]s"
}
`, name, acceptance.HW_NEW_CERTIFICATE_CONTENT, acceptance.HW_NEW_CERTIFICATE_PRIVATE_KEY)
}

func TestAccCertificate_instance(t *testing.T) {
	var (
		cert certificates.Certificate

		rName      = "huaweicloud_apig_certificate.test"
		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()
		baseConfig = testAccCertificate_instanceBase(name)
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&cert,
		getCertificateFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCertificateWithoutRootCA(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCertificate_instance_step1(baseConfig, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrPair(rName, "instance_id", "huaweicloud_apig_instance.test", "id"),
					resource.TestMatchResourceAttr(rName, "effected_at", regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`)),
					resource.TestMatchResourceAttr(rName, "expires_at", regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`)),
					resource.TestMatchResourceAttr(rName, "sans.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"content", "private_key",
				},
			},
			{
				Config: testAccCertificate_instance_step2(baseConfig, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttrPair(rName, "instance_id", "huaweicloud_apig_instance.test", "id"),
					resource.TestMatchResourceAttr(rName, "effected_at", regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`)),
					resource.TestMatchResourceAttr(rName, "expires_at", regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`)),
					resource.TestMatchResourceAttr(rName, "sans.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
				),
			},
		},
	})
}

func testAccCertificate_instanceBase(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_apig_instance" "test" {
  name                  = "%[2]s"
  edition               = "BASIC"
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  enterprise_project_id = "0"

  availability_zones = try(slice(data.huaweicloud_availability_zones.test.names, 0, 1), null)
}
`, common.TestBaseNetwork(name), name)
}

func testAccCertificate_instance_step1(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_certificate" "test" {
  type        = "instance"
  instance_id = huaweicloud_apig_instance.test.id
  name        = "%[2]s"
  content     = "%[3]s"
  private_key = "%[4]s"
}
`, baseConfig, name, acceptance.HW_CERTIFICATE_CONTENT,
		acceptance.HW_CERTIFICATE_PRIVATE_KEY)
}

func testAccCertificate_instance_step2(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_certificate" "test" {
  name        = "%[2]s"
  type        = "instance"
  instance_id = huaweicloud_apig_instance.test.id
  content     = "%[3]s"
  private_key = "%[4]s"
}
`, baseConfig, name, acceptance.HW_NEW_CERTIFICATE_CONTENT,
		acceptance.HW_NEW_CERTIFICATE_PRIVATE_KEY)
}

func TestAccCertificate_instanceWithRootCA(t *testing.T) {
	var (
		cert certificates.Certificate

		rName      = "huaweicloud_apig_certificate.test"
		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()
		baseConfig = testAccCertificate_instanceBase(name)
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&cert,
		getCertificateFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCertificateFull(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCertificate_instanceWithRootCA_step1(baseConfig, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "type", "instance"),
					resource.TestCheckResourceAttrPair(rName, "instance_id", "huaweicloud_apig_instance.test", "id"),
					resource.TestMatchResourceAttr(rName, "effected_at", regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`)),
					resource.TestMatchResourceAttr(rName, "expires_at", regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`)),
					resource.TestMatchResourceAttr(rName, "sans.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
				),
			},
			{
				Config: testAccCertificate_instanceWithRootCA_step2(baseConfig, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "type", "instance"),
					resource.TestCheckResourceAttrPair(rName, "instance_id", "huaweicloud_apig_instance.test", "id"),
					resource.TestMatchResourceAttr(rName, "effected_at", regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`)),
					resource.TestMatchResourceAttr(rName, "expires_at", regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`)),
					resource.TestMatchResourceAttr(rName, "sans.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"content", "private_key", "trusted_root_ca",
				},
			},
		},
	})
}

func testAccCertificate_instanceWithRootCA_step1(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_certificate" "test" {
  type            = "instance"
  instance_id     = huaweicloud_apig_instance.test.id
  name            = "%[2]s"
  content         = "%[3]s"
  private_key     = "%[4]s"
  trusted_root_ca = "%[5]s"
}
`, baseConfig, name, acceptance.HW_CERTIFICATE_CONTENT,
		acceptance.HW_CERTIFICATE_PRIVATE_KEY, acceptance.HW_CERTIFICATE_ROOT_CA)
}

func testAccCertificate_instanceWithRootCA_step2(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_certificate" "test" {
  type            = "instance"
  instance_id     = huaweicloud_apig_instance.test.id
  name            = "%[2]s"
  content         = "%[3]s"
  private_key     = "%[4]s"
  trusted_root_ca = "%[5]s"
}
`, baseConfig, name, acceptance.HW_NEW_CERTIFICATE_CONTENT,
		acceptance.HW_NEW_CERTIFICATE_PRIVATE_KEY, acceptance.HW_NEW_CERTIFICATE_ROOT_CA)
}
