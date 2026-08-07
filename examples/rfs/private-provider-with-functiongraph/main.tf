resource "huaweicloud_fgs_function" "test" {
  name        = var.function_name
  app         = var.function_app
  handler     = var.function_handler
  memory_size = var.function_memory_size
  timeout     = var.function_timeout
  code_type   = "inline"
  runtime     = var.function_runtime
  func_code   = base64encode(var.function_code)
}

resource "huaweicloud_rfs_private_provider" "test" {
  provider_name        = var.private_provider_name
  function_graph_urn   = huaweicloud_fgs_function.test.urn
  provider_description = var.private_provider_description
  provider_version     = var.private_provider_version
  version_description  = var.private_provider_version_description
}

resource "huaweicloud_rfs_private_provider_version" "test" {
  provider_name       = huaweicloud_rfs_private_provider.test.provider_name
  provider_version    = var.provider_version_number
  function_graph_urn  = huaweicloud_fgs_function.test.urn
  version_description = var.provider_version_description
}
