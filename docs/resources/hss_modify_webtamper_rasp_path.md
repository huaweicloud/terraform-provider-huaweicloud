---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_modify_webtamper_rasp_path"
description: |-
  Manages a resource to update the dynamic web tamper protection Tomcat bin directory within HuaweiCloud.
---

# huaweicloud_hss_modify_webtamper_rasp_path

Manages a resource to update the dynamic web tamper protection Tomcat bin directory within HuaweiCloud.

-> This resource is a one-time action resource. Deleting this resource will not clear the corresponding request record,
  but will only remove the resource information from the tf state file.

## Example Usage

```hcl
variable "host_id" {}
variable "rasp_path" {}

resource "huaweicloud_hss_modify_webtamper_rasp_path" "test" {
  host_id   = var.host_id
  rasp_path = var.rasp_path
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `host_id` - (Required, String, NonUpdatable) Specifies the server ID.
  
  -> Only Linux servers are supported. For this parameter to be valid, the server must have WTP enabled, or the WTP
    policy is not deleted after WTP is disabled.

* `rasp_path` - (Required, String, NonUpdatable) Specifies the Tomcat bin directory for dynamic WTP.
  The value length is `1` to `256` characters. The value must start with a slash (/) and cannot end with a slash (/).
  Only letters, numbers, underscores (_), hyphens (-), and periods (.) are allowed.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If it is necessary to operate the hosts under all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

* `host_name` - (Optional, String, NonUpdatable) Specifies the server name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
