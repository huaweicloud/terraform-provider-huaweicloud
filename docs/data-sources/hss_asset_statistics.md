---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_asset_statistics"
description: |-
  Use this data source to get the list of HSS asset statistics within HuaweiCloud.
---

# huaweicloud_hss_asset_statistics

Use this data source to get the list of HSS asset statistics within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_asset_statistics" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `host_id` - (Optional, String) Specifies the host ID.

* `category` - (Optional, String) Specifies the type. The default value is **host**.
  The valid values are as follows:
  + **host**: Host.
  + **container**: Container.

* `enterprise_project_id` - (Optional, String) Specifies the ID of the enterprise project to which the resource belongs.
  This parameter is valid only when the enterprise project function is enabled.
  The value **all_granted_eps** indicates all enterprise projects.
  If omitted, the default enterprise project will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `account_num` - The number of server accounts.

* `port_num` - The number of open ports.

* `process_num` - The number of processes.

* `app_num` - The number of applications.

* `auto_launch_num` - The number of auto launch startup processes.

* `web_framework_num` - The number of web frameworks.

* `web_site_num` - The number of web sites.

* `jar_package_num` - The number of JAR packages.

* `kernel_module_num` - The number of kernel modules.

* `web_service_num` - The number of web services.

* `web_app_num` - The number of web applications.

* `database_num` - The number of databases.
