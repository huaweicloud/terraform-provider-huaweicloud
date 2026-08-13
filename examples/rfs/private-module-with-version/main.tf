resource "huaweicloud_rfs_private_module" "test" {
  module_name        = var.module_name
  module_description = var.module_description
}

resource "huaweicloud_rfs_private_module_version" "test" {
  module_name         = huaweicloud_rfs_private_module.test.module_name
  module_version      = var.module_version
  module_uri          = var.module_uri
  version_description = var.version_description
}
