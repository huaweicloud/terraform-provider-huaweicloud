---
subcategory: "Application Operations Management (AOM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aom_prom_instance"
description: ""
---

# huaweicloud_aom_prom_instance

Manages an AOM prometheus instance resource within HuaweiCloud.

## Example Usage

```hcl
variable "enterprise_project_id" {}

resource "huaweicloud_aom_prom_instance" "instance" {
  prom_name             = "test_demo"
  prom_type             = "ECS"
  enterprise_project_id = var.enterprise_project_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the AOM prometheus instance.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `prom_name` - (Required, String, ForceNew) Specifies the name of the AOM prometheus instance.
  The value can contain `1` to `100` characters. Only Chinese and English letters, digits, underscores (_)
  and hyphens (-) are allowed, and it must start with letters, digits or Chinese characters.
  Changing this parameter will create a new resource.

* `prom_type` - (Required, String, ForceNew) Specifies the type of the AOM prometheus instance.
  The value can be: **ECS**, **VPC**, **CCE**, **REMOTE_WRITE**, **KUBERNETES**,
  **CLOUD_SERVICE** or **ACROSS_ACCOUNT**.
  Changing this parameter will create a new resource.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project ID of the
  AOM prometheus instance. Defaults to **0**. Changing this parameter will create a new resource.

* `prom_version` - (Optional, String, ForceNew) Specifies the version of AOM prometheus instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `created_at` - The creation time of AOM prometheus instance.

* `prom_http_api_endpoint` - The url for calling the AOM Prometheus instance.

* `remote_read_url` - The remote read address of AOM Prometheus instance.

* `remote_write_url` - The remote write address of AOM Prometheus instance.

## Import

The AOM prometheus instance can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_aom_prom_instance.test <id>
```
