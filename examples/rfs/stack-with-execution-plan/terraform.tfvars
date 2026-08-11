stack_name          = "tf-test-stack"
execution_plan_name = "tf-test-plan"

execution_plan_template_body = <<EOT
terraform {
  required_providers {
    huaweicloud = {
      source  = "huawei.com/provider/huaweicloud"
      version = ">= 1.41.0"
    }
  }
}

resource "huaweicloud_vpc" "test" {
  name = "tf-test-vpc"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id     = huaweicloud_vpc.test.id
  name       = "tf-test-subnet"
  cidr       = "192.168.1.0/24"
  gateway_ip = "192.168.1.1"
}
EOT
