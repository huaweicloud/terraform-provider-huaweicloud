# SecMaster configuration dictionary
resource "huaweicloud_secmaster_configuration_dictionary" "test" {
  dict_id      = var.dict_id
  dict_key     = var.dict_key
  dict_code    = var.dict_code
  dict_val     = var.dict_val
  language     = var.language
  version      = var.dict_version
  dict_pkey    = var.dict_pkey
  dict_pcode   = var.dict_pcode
  scope        = var.scope
  description  = var.description
  extend_field = var.extend_field
  is_built_in  = false

  lifecycle {
    ignore_changes = [
      is_built_in,
    ]
  }
}
