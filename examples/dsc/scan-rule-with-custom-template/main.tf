resource "huaweicloud_dsc_scan_security_level" "test" {
  security_level_name = var.security_level_name
  color_number        = var.security_level_color_number
  security_level_desc = var.security_level_description
}

resource "huaweicloud_dsc_scan_template" "test" {
  action             = "ADD"
  name               = var.scan_template_name
  description        = var.scan_template_description
  add_built_in_rules = false
}

resource "huaweicloud_dsc_scan_template_classification" "test" {
  template_id         = huaweicloud_dsc_scan_template.test.id
  classification_name = var.classification_name
}

resource "huaweicloud_dsc_scan_rule" "test" {
  rule_name      = var.scan_rule_name
  rule_type      = "REGEX"
  category       = "BUILT_SELF"
  logic_operator = "AND"
  match_rate     = var.scan_rule_match_rate
  min_match      = 1
  rule_desc      = var.scan_rule_description

  content {
    effective_mode = "NOT_IN"
    location       = "NAME"
    rule_content   = "bphone"
  }

  content {
    effective_mode = "IN"
    location       = "REMARK"
    rule_content   = "telephone number"
  }

  templates {
    template_id       = huaweicloud_dsc_scan_template.test.id
    classification_id = huaweicloud_dsc_scan_template_classification.test.id
    security_level_id = huaweicloud_dsc_scan_security_level.test.id
    is_used           = "true"
  }
}
