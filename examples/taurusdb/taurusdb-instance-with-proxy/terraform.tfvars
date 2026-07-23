vpc_name                     = "your_vpc"
subnet_name                  = "your_subnet"
security_group_name          = "your_security_group"
instance_name                = "your_taurusdb_instance"
instance_backup_time_window  = "09:00-10:00"
instance_backup_keep_days    = 7
proxy_name                   = "your_proxy"
proxy_node_num               = 2
proxy_port                   = 3339
proxy_access_control_ip_list = [
  {
    ip          = "3.3.3.3"
    description = "test description"
  }
]
proxy_parameters             = [
  {
    name      = "multiStatementType"
    value     = "Loose"
    elem_type = "system"
  }
]
