---
subcategory: "Application Operations Management (AOM 2.0)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aom_cmdb_environment"
description: ""
---

# huaweicloud_aom_cmdb_environment

Manages an AOM environment resource within HuaweiCloud.

## Example Usage

```hcl
variable "com_id" {}

resource "huaweicloud_aom_cmdb_environment" "test" {
  name         = "env_demo"
  component_id = var.com_id
  type         = "DEV"
  os_type      = "LINUX"
  description  = "environment description"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region to which the environment is to be associated.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `component_id` - (Required, String, ForceNew) Specifies the component ID. Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the name of the environment. The value can contain 2 to 64 characters.
  Only letters, digits, underscores (_), hyphens (-), and periods (.) are allowed.

* `type` - (Required, String) Specifies the type of the environment. The value can be **DEV**, **TEST**, **PRE** and **ONLINE**.

* `os_type` - (Required, String, ForceNew) Specifies the OS type. The value can be **WINDOWS** and **LINUX**.
  Changing this parameter will create a new resource.

* `description` - (Optional, String) Specifies the description about the environment.
  The description can contain a maximum of 255 characters.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
* `enterprise_project_id` - The enterprise project ID.
* `register_type` - The register type of the environment.
* `created_at` - The creation time.

## Import

The AOM environment can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_aom_cmdb_environment.test 4ee753f776114565863d260f1cc62695
```
