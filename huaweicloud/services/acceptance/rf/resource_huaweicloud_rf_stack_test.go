package rf

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/rf/v1/stacks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/rf"
)

func getStackesourceFunc(config *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := config.AosV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating AOS v3 client: %s", err)
	}

	return rf.QueryStackById(client, state.Primary.ID)
}

func TestAccStack_basic(t *testing.T) { // the template file is json format.
	var (
		obj stacks.Stack

		rName = "huaweicloud_rf_stack.test"
		name  = acceptance.RandomAccResourceNameWithDash()
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getStackesourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccStack_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "Create by acc test"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				Config: testAccStack_withBody_step1(name, basicTemplateInJsonFormat(name)),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
				),
			},
			{
				Config: testAccStack_withBody_step2(name, updateTemplateInJsonFormat(name), variableContent),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"agency",
					"template_body",
					"vars_body",
				},
			},
		},
	})
}

func testAccStack_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_rf_stack" "test" {
  name        = "%[1]s"
  description = "Create by acc test"
}
`, name)
}

func testAccStack_withBody_step1(name, template string) string {
	return fmt.Sprintf(`
resource "huaweicloud_rf_stack" "test" {
  name        = "%[1]s"
  description = "Create by acc test"

  agency {
    name          = "rf_admin_trust" // System RF agency
    provider_name = "huaweicloud"
  }

  template_body = %[2]s
}
`, name, template)
}

func testAccStack_withBody_step2(name, template, vars string) string {
	return fmt.Sprintf(`
resource "huaweicloud_rf_stack" "test" {
  name        = "%[1]s"
  description = "Create by acc test"

  agency {
    name          = "rf_admin_trust" // System RF agency
    provider_name = "huaweicloud"
  }

  template_body = %[2]s
  vars_body     = %[3]s
}
`, name, updateTemplateInJsonFormat(name), variableContent)
}

func TestAccStack_withBody_HCL(t *testing.T) {
	var (
		obj stacks.Stack

		rName = "huaweicloud_rf_stack.test"
		name  = acceptance.RandomAccResourceNameWithDash()
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getStackesourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccStack_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "Create by acc test"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				Config: testAccStack_withBody_step1(name, basicTemplateInHclFormat(name)),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
				),
			},
			{
				Config: testAccStack_withBody_step2(name, updateTemplateInHclFormat(name), variableContent),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"agency",
					"template_body",
					"vars_body",
				},
			},
		},
	})
}

func TestAccStack_withUri_JSON(t *testing.T) {
	var (
		obj stacks.Stack

		rName = "huaweicloud_rf_stack.test"
		name  = acceptance.RandomAccResourceNameWithDash()
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getStackesourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccStack_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				Config: testAccStack_withUri_step1(name, basicTemplateInJsonFormat(name), variableContent),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
				),
			},
			{
				Config: testAccStack_withUri_step2(name, basicTemplateInJsonFormat(name), variableContent),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"agency",
					"template_uri",
					"vars_uri",
				},
			},
		},
	})
}

func testAccStack_withUri_base(name, template, vars string) string {
	return fmt.Sprintf(`

resource "huaweicloud_obs_bucket" "test" {
  bucket = "%[1]s"
  acl    = "private"
}

resource "huaweicloud_obs_bucket_policy" "test" {
  bucket = huaweicloud_obs_bucket.test.bucket
  policy = <<EOT
{
  "Statement": [
    {
      "Sid": "RF-Access",
      "Effect": "Allow",
      "Principal": {
        "ID": ["*"]
      },
      "Action": [
        "GetObject"
      ],
      "Resource": [
        "${huaweicloud_obs_bucket.test.bucket}/rf/resource_stack/uri_test/*"
      ]
    }
  ]
}
EOT
}

resource "huaweicloud_obs_bucket_object" "template" {
  bucket       = huaweicloud_obs_bucket.test.bucket
  key          = "rf/resource_stack/uri_test/template.tf.json"
  content_type = "application/json"
  content      = %[2]s
}

resource "huaweicloud_obs_bucket_object" "variable" {
  bucket       = huaweicloud_obs_bucket.test.bucket
  key          = "rf/resource_stack/uri_test/resource.tfvars"
  content_type = "application/octet-stream"
  content      = %[3]s
}
`, name, template, vars)
}

func testAccStack_withUri_step1(name, template, vars string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rf_stack" "test" {
  name = "%[2]s"

  agency {
    name          = "rf_admin_trust" // System RF agency
    provider_name = "huaweicloud"
  }

  template_uri = "https://${huaweicloud_obs_bucket.test.bucket_domain_name}/${huaweicloud_obs_bucket_object.template.id}"
}
`, testAccStack_withUri_base(name, template, vars), name)
}

func testAccStack_withUri_step2(name, template, vars string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rf_stack" "test" {
  name = "%[2]s"

  agency {
    name          = "rf_admin_trust" // System RF agency
    provider_name = "huaweicloud"
  }

  template_uri = "https://${huaweicloud_obs_bucket.test.bucket_domain_name}/${huaweicloud_obs_bucket_object.template.id}"
  vars_uri     = "https://${huaweicloud_obs_bucket.test.bucket_domain_name}/${huaweicloud_obs_bucket_object.variable.id}"
}
`, testAccStack_withUri_base(name, template, vars), name)
}

func TestAccStack_withUri_HCL(t *testing.T) {
	var (
		obj stacks.Stack

		rName = "huaweicloud_rf_stack.test"
		name  = acceptance.RandomAccResourceNameWithDash()
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getStackesourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccStack_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				Config: testAccStack_withUri_step1(name, basicTemplateInHclFormat(name), variableContent),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
				),
			},
			{
				Config: testAccStack_withUri_step2(name, basicTemplateInHclFormat(name), variableContent),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"agency",
					"template_uri",
					"vars_uri",
				},
			},
		},
	})
}

func TestAccStack_archive(t *testing.T) {
	var (
		obj stacks.Stack

		rName = "huaweicloud_rf_stack.test"
		name  = acceptance.RandomAccResourceNameWithDash()
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getStackesourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRfArchives(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccStack_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				Config: testAccStack_archive_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
				),
			},
			{
				Config: testAccStack_archive_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"agency",
					"template_uri",
					"vars_uri",
				},
			},
		},
	})
}

func testAccStack_archive_step1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_rf_stack" "test" {
  name = "%[1]s"

  agency {
    name          = "rf_admin_trust" // System RF agency
    provider_name = "huaweicloud"
  }

  template_uri = "%[2]s"
}
`, name, acceptance.HW_RF_TEMPLATE_ARCHIVE_NO_VARS_URI)
}

func testAccStack_archive_step2(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_rf_stack" "test" {
  name = "%[1]s"

  agency {
    name          = "rf_admin_trust" // System RF agency
    provider_name = "huaweicloud"
  }

  template_uri = "%[2]s"
  vars_uri     = "%[3]s"
}
`, name, acceptance.HW_RF_TEMPLATE_ARCHIVE_URI, acceptance.HW_RF_VARIABLES_ARCHIVE_URI)
}

func basicTemplateInJsonFormat(name string) string {
	return fmt.Sprintf(`<<EOT
{
  "terraform": {
    "required_providers": [
      {
        "huaweicloud": {
          "source": "huawei.com/provider/huaweicloud",
          "version": ">= 1.41.0"
        }
      }
    ]
  },
  "provider": {
    "huaweicloud": {
      "region": "%[1]s"
    }
  },
  "resource": {
    "huaweicloud_vpc": {
      "test": {
        "name": "%[2]s",
        "cidr": "192.168.0.0/16"
      }
    }
  }
}
EOT
`, acceptance.HW_REGION_NAME, name)
}

func updateTemplateInJsonFormat(name string) string {
	return fmt.Sprintf(`<<EOT
{
  "terraform": {
    "required_providers": [
      {
        "huaweicloud": {
          "source": "huawei.com/provider/huaweicloud",
          "version": ">= 1.41.0"
        }
      }
    ]
  },
  "provider": {
    "huaweicloud": {
      "region": "%[1]s"
    }
  },
  "resource": {
    "huaweicloud_vpc": {
      "test": {
        "name": "%[2]s",
        "cidr": "192.168.0.0/16"
      }
    },
    "huaweicloud_vpc_subnet": {
      "test": {
        "vpc_id": "$${huaweicloud_vpc.test.id}",
        "name": "$${var.subnet_name}",
        "cidr": "$${cidrsubnet(huaweicloud_vpc.test.cidr, 4, 1)}",
        "gateway_ip": "$${cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 4, 1), 1)}"
      }
    }
  },
  "variable": {
    "subnet_name": {
      "type": "string"
    }
  }
}
EOT
`, acceptance.HW_REGION_NAME, name)
}

func basicTemplateInHclFormat(name string) string {
	// lintignore:AT004
	return fmt.Sprintf(`<<EOT
terraform {
  required_providers {
    huaweicloud = {
      source  = "huawei.com/provider/huaweicloud"
      version = ">= 1.41.0"
    }
  }
}

provider "huaweicloud" {
  region = "%[1]s"
}

resource "huaweicloud_vpc" "test" {
  name = "%[2]s"
  cidr = "192.168.0.0/16"
}
EOT
`, acceptance.HW_REGION_NAME, name)
}

func updateTemplateInHclFormat(name string) string {
	// lintignore:AT004
	return fmt.Sprintf(`<<EOT
terraform {
  required_providers {
    huaweicloud = {
      source  = "huawei.com/provider/huaweicloud"
      version = ">= 1.41.0"
    }
  }
}

provider "huaweicloud" {
  region = "%[1]s"
}

resource "huaweicloud_vpc" "test" {
  name = "%[2]s",
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id      = "$${huaweicloud_vpc.test.id}"
  name        = "$${var.subnet_name}"
  cidr        = "$${cidrsubnet(huaweicloud_vpc.test.cidr, 4, 1)}"
  gateway_ip" = "$${cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 4, 1), 1)}"
}

variable "subnet_name" {
  type = "string"
}
EOT
`, acceptance.HW_REGION_NAME, name)
}

const variableContent = `<<EOT
subnet_name = "tf-test-demo"
EOT
`
