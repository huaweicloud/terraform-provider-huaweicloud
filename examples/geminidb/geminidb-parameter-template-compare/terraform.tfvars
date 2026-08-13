parameter_template_configuration = {
  "source_parameter_template" = {
    template_name        = "your_source_parameter_template"
    template_description = ""
    parameter_values     = { request_timeout_in_ms = "30000" }
  }
  "target_parameter_template" = {
    template_name        = "your_target_parameter_template"
    template_description = ""
    parameter_values     = {request_timeout_in_ms = "20000"}
  }
}
