resource "huaweicloud_dsc_mask_algorithm" "test" {
  algorithm_name = var.mask_algorithm_name
  algorithm      = "PRESNM"
  algorithm_type = "MASK_BY_OVERWRITE"
  category       = "BUILT_SELF"

  parameter = jsonencode({
    type   = "CHAR"
    first  = var.mask_algorithm_prefix_length
    second = var.mask_algorithm_suffix_length
    method = var.mask_algorithm_replacement
  })
}
