resource "huaweicloud_ces_dashboard" "test" {
  name           = var.dashboard_name
  row_widget_num = var.dashboard_row_widget_num

  dynamic "extend_info" {
    for_each = var.dashboard_extend_info

    content {
      filter                  = extend_info.value.filter
      period                  = extend_info.value.period
      display_time            = extend_info.value.display_time
      refresh_time            = extend_info.value.refresh_time
      from                    = extend_info.value.from
      to                      = extend_info.value.to
      screen_color            = extend_info.value.screen_color
      enable_screen_auto_play = extend_info.value.enable_screen_auto_play
      time_interval           = extend_info.value.time_interval
      enable_legend           = extend_info.value.enable_legend
      full_screen_widget_num  = extend_info.value.full_screen_widget_num
    }
  }
}
