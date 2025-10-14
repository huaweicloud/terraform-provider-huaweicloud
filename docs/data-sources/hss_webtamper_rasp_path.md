---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_webtamper_rasp_path"
description: |-
  Use this data source to query the Tomcat bin directory configured for dynamic web tamper protection.
---

# huaweicloud_hss_webtamper_rasp_path

Use this data source to query the Tomcat bin directory configured for dynamic web tamper protection.

## Example Usage

```hcl
variable "host_id" {}

data "huaweicloud_hss_webtamper_rasp_path" "test" {
  host_id = var.host_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `host_id` - (Required, String) Specifies the ID of the server.

  -> Only Linux servers are supported. The server must have WTP (web tamper protection) enabled or the WTP policy
    is not deleted after WTP is disabled.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `rasp_path` - The Tomcat bin directory for dynamic web tamper protection.
