resource "huaweicloud_vpc" "test" {
  name = var.vpc_name
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_servicestagev3_environment" "test" {
  name        = var.environment_name
  description = var.environment_description
  vpc_id      = huaweicloud_vpc.test.id
}

resource "huaweicloud_servicestagev3_application" "test" {
  name        = var.application_name
  description = var.application_description
}

resource "huaweicloud_servicestagev3_application_configuration" "test" {
  environment_id = huaweicloud_servicestagev3_environment.test.id
  application_id = huaweicloud_servicestagev3_application.test.id

  configuration {
    env {
      name  = var.application_configuration_env_name
      value = var.application_configuration_env_value
    }
  }
}
