resource "huaweicloud_vpc" "test" {
  name = var.vpc_name
  cidr = var.vpc_cidr
}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id     = huaweicloud_vpc.test.id
  name       = var.subnet_name
  cidr       = var.subnet_cidr == "" ? cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0) : var.subnet_cidr
  gateway_ip = var.subnet_gateway_ip == "" ? cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0), 1) : var.subnet_gateway_ip
}

resource "huaweicloud_networking_secgroup" "test" {
  name = var.security_group_name
}

resource "huaweicloud_fgs_function" "test" {
  name        = var.function_name
  memory_size = var.function_memory_size
  runtime     = var.function_runtime
  timeout     = var.function_timeout
  handler     = var.function_handler
  code_type   = var.function_code_type
  app         = var.function_app
  func_code   = var.function_code
}

data "huaweicloud_availability_zones" "test" {
  count = length(var.availability_zones) == 0 ? 1 : 0
}

resource "huaweicloud_apig_instance" "test" {
  name                  = var.instance_name
  edition               = var.instance_edition
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  enterprise_project_id = var.enterprise_project_id
  availability_zones    = length(var.availability_zones) == 0 ? try(slice(data.huaweicloud_availability_zones.test[0].names, 0, var.availability_zones_count), null) : var.availability_zones
}

resource "huaweicloud_apig_custom_authorizer" "test" {
  instance_id      = huaweicloud_apig_instance.test.id
  name             = var.custom_authorizer_name
  function_urn     = huaweicloud_fgs_function.test.urn
  function_version = var.function_version
  type             = var.custom_authorizer_type
  network_type     = var.custom_authorizer_network_type
  cache_age        = var.custom_authorizer_cache_age
  is_body_send     = var.custom_authorizer_is_body_send
  user_data        = var.custom_authorizer_use_data

  dynamic "identity" {
    for_each = var.custom_authorizer_identity

    content {
      name       = identity.value.name
      location   = identity.value.location
      validation = identity.value.validation
    }
  }
}

resource "huaweicloud_apig_response" "test" {
  name        = var.response_name
  instance_id = huaweicloud_apig_instance.test.id
  group_id    = huaweicloud_apig_group.test.id

  dynamic "rule" {
    for_each = var.response_rules

    content {
      error_type  = rule.value["error_type"]
      body        = rule.value["body"]
      status_code = rule.value["status_code"]

      dynamic "headers" {
        for_each = rule.value["headers"]

        content {
          key   = headers.value["key"]
          value = headers.value["value"]
        }
      }
    }
  }
}

resource "huaweicloud_apig_group" "test" {
  name        = var.group_name
  instance_id = huaweicloud_apig_instance.test.id
}

resource "huaweicloud_apig_api" "test" {
  instance_id             = huaweicloud_apig_instance.test.id
  group_id                = huaweicloud_apig_group.test.id
  type                    = var.api_type
  name                    = var.api_name
  request_protocol        = var.api_request_protocol
  request_method          = var.api_request_method
  request_path            = var.api_request_path
  security_authentication = "AUTHORIZER"
  matching                = var.api_matching
  response_id             = huaweicloud_apig_response.test.id
  authorizer_id           = huaweicloud_apig_custom_authorizer.test.id

  dynamic "backend_params" {
    for_each = var.api_backend_params

    content {
      type              = backend_params.value["type"]
      name              = backend_params.value["name"]
      location          = backend_params.value["location"]
      value             = backend_params.value["value"]
      system_param_type = backend_params.value["system_param_type"]
    }
  }

  func_graph {
    function_urn     = huaweicloud_fgs_function.test.urn
    version          = var.function_version
    network_type     = var.api_func_graph_network_type
    request_protocol = var.api_func_graph_request_protocol
  }
}
