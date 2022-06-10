package apig

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/responses"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccApigResponseV2_basic(t *testing.T) {
	var (
		// Only letters, digits and underscores (_) are allowed in the name.
		rName        = fmt.Sprintf("tf_acc_test_%s", acctest.RandString(5))
		resourceName = "huaweicloud_apig_response.test"
		resp         responses.Response
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t) // The creation of APIG instance needs the enterprise project ID.
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckApigResponseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccApigResponse_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApigResponseExists(resourceName, &resp),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
			{
				Config: testAccApigResponse_update(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApigResponseExists(resourceName, &resp),
					resource.TestCheckResourceAttr(resourceName, "name", rName+"_update"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccApigGroupSubResourceImportStateIdFunc(resourceName),
			},
		},
	})
}

func TestAccApigResponseV2_customRules(t *testing.T) {
	var (
		// Only letters, digits and underscores (_) are allowed in the name.
		rName        = fmt.Sprintf("tf_acc_test_%s", acctest.RandString(5))
		resourceName = "huaweicloud_apig_response.test"
		resp         responses.Response
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t) // The creation of APIG instance needs the enterprise project ID.
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckApigResponseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccApigResponse_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApigResponseExists(resourceName, &resp),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
			{
				// Add two custom rule.
				Config: testAccApigResponse_rules(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApigResponseExists(resourceName, &resp),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "2"),
				),
			},
			{
				// Remove one and update another.
				Config: testAccApigResponse_rulesUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApigResponseExists(resourceName, &resp),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "rule.#", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccApigGroupSubResourceImportStateIdFunc(resourceName),
			},
		},
	})
}

func testAccCheckApigResponseDestroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	client, err := config.ApigV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_apig_response" {
			continue
		}
		_, err := responses.Get(client, rs.Primary.Attributes["instance_id"], rs.Primary.Attributes["group_id"],
			rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("APIG v2 API custom response (%s) is still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckApigResponseExists(n string, resp *responses.Response) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource %s not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No APIG V2 API custom response Id")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		client, err := config.ApigV2Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
		}
		found, err := responses.Get(client, rs.Primary.Attributes["instance_id"], rs.Primary.Attributes["group_id"],
			rs.Primary.ID).Extract()
		if err != nil {
			return fmt.Errorf("APIG v2 API custom response not exist: %s", err)
		}
		*resp = *found
		return nil
	}
}

func isResponseImportIdValid(rs *terraform.ResourceState) bool {
	if rs.Primary.Attributes["instance_id"] == "" || rs.Primary.Attributes["group_id"] == "" ||
		rs.Primary.Attributes["name"] == "" {
		return false
	}
	return true
}

func testAccApigGroupSubResourceImportStateIdFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("Resource (%s) not found: %s", name, rs)
		}
		if isResponseImportIdValid(rs) {
			return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["instance_id"], rs.Primary.Attributes["group_id"],
				rs.Primary.Attributes["name"]), nil
		}
		return "", fmt.Errorf("resource not found: %s/%s/%s", rs.Primary.Attributes["instance_id"],
			rs.Primary.Attributes["group_id"], rs.Primary.Attributes["name"])
	}
}

func testAccApigResponse_base(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_apig_group" "test" {
  name        = "%s"
  instance_id = huaweicloud_apig_instance.test.id
  description = "Created by script"
}
`, testAccApigApplication_base(rName), rName)
}

func testAccApigResponse_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_apig_response" "test" {
  name        = "%s"
  instance_id = huaweicloud_apig_instance.test.id
  group_id    = huaweicloud_apig_group.test.id

  rule {
    error_type  = "AUTHORIZER_FAILURE"
    body        = "{\"code\":\"$context.authorizer.frontend.code\",\"message\":\"$context.authorizer.frontend.message\"}"
    status_code = 401
  }
}
`, testAccApigResponse_base(rName), rName)
}

func testAccApigResponse_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_apig_response" "test" {
  name        = "%s_update"
  instance_id = huaweicloud_apig_instance.test.id
  group_id    = huaweicloud_apig_group.test.id

  rule {
    error_type  = "AUTHORIZER_FAILURE"
    body        = "{\"code\":\"$context.authorizer.frontend.code\",\"message\":\"$context.authorizer.frontend.message\"}"
    status_code = 401
  }
}
`, testAccApigResponse_base(rName), rName)
}

func testAccApigResponse_rules(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_apig_response" "test" {
  name        = "%s"
  instance_id = huaweicloud_apig_instance.test.id
  group_id    = huaweicloud_apig_group.test.id

  rule {
    error_type  = "ACCESS_DENIED"
    body        = "{\"error_code\":\"$context.error.code\",\"error_msg\":\"$context.error.message\"}"
    status_code = 400
  }
  rule {
    error_type  = "AUTHORIZER_FAILURE"
    body        = "{\"code\":\"$context.authorizer.frontend.code\",\"message\":\"$context.authorizer.frontend.message\"}"
    status_code = 401
  }
}
`, testAccApigResponse_base(rName), rName)
}

func testAccApigResponse_rulesUpdate(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_apig_response" "test" {
  name        = "%s"
  instance_id = huaweicloud_apig_instance.test.id
  group_id    = huaweicloud_apig_group.test.id

  rule {
    error_type  = "AUTHORIZER_FAILURE"
    body        = "{\"code\":\"$context.authorizer.frontend.code\",\"message\":\"$context.authorizer.frontend.message\"}"
    status_code = 403
  }
}
`, testAccApigResponse_base(rName), rName)
}
