package rfs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/rfs"
)

func getStackesourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("rfs", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating RFS client: %s", err)
	}

	return rfs.QueryStackById(client, state.Primary.ID)
}

func TestAccStack_basic(t *testing.T) {
	var (
		obj interface{}

		rName = "huaweicloud_rfs_stack.test"
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
resource "huaweicloud_rfs_stack" "test" {
  name        = "%[1]s"
  description = "Create by acc test"
}
`, name)
}

func TestAccStack_withBody(t *testing.T) {
	var (
		obj interface{}

		rName = "huaweicloud_rfs_stack.test"
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
				Config: testAccStack_withBody(name, basicTemplateInJsonFormat(name), "null"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
				),
			},
			{
				Config: testAccStack_withBody(name, updateTemplateInJsonFormat(), basicVariablesInVarsFormat(name)),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
				),
			},
			{
				Config: testAccStack_withBody(name, basicTemplateInHclFormat(name), "null"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
				),
			},
			{
				Config: testAccStack_withBody(name, updateTemplateInHclFormat(), basicVariablesInVarsFormat(name)),
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

func updateTemplateInJsonFormat() string {
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
        "name": "$${var.resource_name}",
        "cidr": "192.168.0.0/16"
      }
    },
    "huaweicloud_vpc_subnet": {
      "test": {
        "vpc_id": "$${huaweicloud_vpc.test.id}",
        "name": "$${var.resource_name}",
        "cidr": "$${cidrsubnet(huaweicloud_vpc.test.cidr, 4, 1)}",
        "gateway_ip": "$${cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 4, 1), 1)}"
      }
    }
  },
  "variable": {
    "resource_name": {
      "type": "string"
    }
  }
}
EOT
`, acceptance.HW_REGION_NAME)
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

func updateTemplateInHclFormat() string {
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
  name = var.resource_name
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id     = huaweicloud_vpc.test.id
  name       = var.resource_name
  cidr       = cidrsubnet(huaweicloud_vpc.test.cidr, 4, 1)
  gateway_ip = cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 4, 1), 1)
}

variable "resource_name" {
  type = string
}
EOT
`, acceptance.HW_REGION_NAME)
}

func basicVariablesInVarsFormat(name string) string {
	return fmt.Sprintf(`<<EOT
resource_name = "%[1]s"
EOT
`, name)
}

func testAccStack_withBody(name, template, vars string) string {
	return fmt.Sprintf(`
resource "huaweicloud_rfs_stack" "test" {
  name        = "%[1]s"
  description = "Create by acc test"

  agency {
    name          = "rf_admin_trust" # System RF agency
    provider_name = "huaweicloud"
  }

  template_body = %[2]s
  vars_body     = %[3]s
}
`, name, template, vars)
}

func TestAccStack_withUri(t *testing.T) {
	var (
		obj interface{}

		rName = "huaweicloud_rfs_stack.test"
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
				Config: testAccStack_withUri_jsonBody(name, basicTemplateInJsonFormat(name), "null"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
				),
			},
			{
				Config: testAccStack_withUri_jsonBody(name, updateTemplateInJsonFormat(), basicVariablesInVarsFormat(name)),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
				),
			},
			{
				Config: testAccStack_withUri_hclBody(name, basicTemplateInHclFormat(name), "null"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
				),
			},
			{
				Config: testAccStack_withUri_hclBody(name, updateTemplateInHclFormat(), basicVariablesInVarsFormat(name)),
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

func testAccStack_withUri_base(name, templatem, vars string) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "test" {
  bucket = "%[1]s"
  acl    = "private"
}

variable "script_template" {
  default = %[2]s
}

variable "script_variables" {
  default = %[3]s
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

resource "huaweicloud_obs_bucket_object" "variables" {
  count = var.script_variables != null ? 1 : 0

  bucket       = huaweicloud_obs_bucket.test.bucket
  key          = "rf/resource_stack/uri_test/variable.tfvars"
  content_type = "application/octet-stream"
  content      = var.script_variables
  
  provisioner "local-exec" {
    command = "sleep 30"
  }
}
`, name, templatem, vars)
}

func testAccStack_withUri_jsonBody(name, template, vars string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_obs_bucket_object" "template" {
  bucket       = huaweicloud_obs_bucket.test.bucket
  key          = "rf/resource_stack/uri_test/template.tf.json"
  content_type = "application/json"
  content      = var.script_template
  
  provisioner "local-exec" {
    command = "sleep 30"
  }
}

resource "huaweicloud_rfs_stack" "test" {
  depends_on = [
    huaweicloud_obs_bucket_policy.test,
    huaweicloud_obs_bucket_object.variables,
  ]

  name = "%[2]s"

  agency {
    name          = "rf_admin_trust" # System RF agency
    provider_name = "huaweicloud"
  }

  template_uri = format("https://%%s/%%s",
    huaweicloud_obs_bucket.test.bucket_domain_name,
    huaweicloud_obs_bucket_object.template.id)
  vars_uri     = var.script_variables != null ? format("https://%%s/%%s",
    huaweicloud_obs_bucket.test.bucket_domain_name,
    huaweicloud_obs_bucket_object.variables[0].id) : null

  lifecycle {
    replace_triggered_by = [
      huaweicloud_obs_bucket_object.template.content,
    ]
  }
}
`, testAccStack_withUri_base(name, template, vars), name)
}

func testAccStack_withUri_hclBody(name, template, vars string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_obs_bucket_object" "template" {
  bucket       = huaweicloud_obs_bucket.test.bucket
  key          = "rf/resource_stack/uri_test/template.tf"
  content_type = "text/plain"
  content      = var.script_template
  
  provisioner "local-exec" {
    command = "sleep 30"
  }
}

resource "huaweicloud_rfs_stack" "test" {
  depends_on = [
    huaweicloud_obs_bucket_policy.test,
    huaweicloud_obs_bucket_object.variables,
  ]

  name = "%[2]s"

  agency {
    name          = "rf_admin_trust" # System RF agency
    provider_name = "huaweicloud"
  }

  template_uri = format("https://%%s/%%s",
    huaweicloud_obs_bucket.test.bucket_domain_name,
    huaweicloud_obs_bucket_object.template.id)
  vars_uri     = var.script_variables != null ? format("https://%%s/%%s",
    huaweicloud_obs_bucket.test.bucket_domain_name,
    huaweicloud_obs_bucket_object.variables[0].id) : null

  lifecycle {
    replace_triggered_by = [
      # Rebuild to ensure that the correct storage objects are accessed when the resource stack is deployed.
      huaweicloud_obs_bucket_object.template.content,
    ]
  }
}
`, testAccStack_withUri_base(name, template, vars), name)
}

func TestAccStack_archive(t *testing.T) {
	var (
		obj interface{}

		rName = "huaweicloud_rfs_stack.test"
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
resource "huaweicloud_rfs_stack" "test" {
  name = "%[1]s"

  agency {
    name          = "rf_admin_trust" # System RF agency
    provider_name = "huaweicloud"
  }

  template_uri = "%[2]s"
}
`, name, acceptance.HW_RF_TEMPLATE_ARCHIVE_NO_VARS_URI)
}

func testAccStack_archive_step2(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_rfs_stack" "test" {
  name = "%[1]s"

  agency {
    name          = "rf_admin_trust" # System RF agency
    provider_name = "huaweicloud"
  }

  template_uri = "%[2]s"
  vars_uri     = "%[3]s"
}
`, name, acceptance.HW_RF_TEMPLATE_ARCHIVE_URI, acceptance.HW_RF_VARIABLES_ARCHIVE_URI)
}
