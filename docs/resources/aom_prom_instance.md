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
  prom_version          = "1.5"
  
  prom_limits {
    compactor_blocks_retention_period = "360h"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the prometheus instance is located.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `prom_name` - (Required, String) Specifies the name of the prometheus instance.  
  The value can contain `1` to `100` characters. Only Chinese and English letters, digits, underscores (_)
  and hyphens (-) are allowed, and it must start with letters, digits or Chinese characters.

* `prom_type` - (Required, String, ForceNew) Specifies the type of the prometheus instance.
  The value can be: **ECS**, **VPC**, **CCE**, **REMOTE_WRITE**, **KUBERNETES**,
  **CLOUD_SERVICE** or **ACROSS_ACCOUNT**.
  Changing this parameter will create a new resource.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project ID to which the prometheus
  instance belongs.  
  This field is only valid for enterprise users, if omitted, default enterprise project will be used.  
  Changing this parameter will create a new resource.

* `prom_version` - (Optional, String, ForceNew) Specifies the version of the prometheus instance.

* `prom_limits` - (Optional, List) Specifies the limit configurations of the prometheus instance.
  The [prom_limits](#aom_prom_instance_prom_limits) structure is documented below.

  -> The limits can only be updated once every `24` hours.

<a name="aom_prom_instance_prom_limits"></a>
The `prom_limits` block supports:

* `compactor_blocks_retention_period` - (Required, String) Specifies the retention period for the compactor blocks.  
  The valid values are as follows:
  + **360h**
  + **420h**
  + **1440h**
  + **2160h**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `prom_http_api_endpoint` - The HTTP URL for calling the prometheus instance.

* `remote_read_url` - The remote read address of the prometheus instance.

* `remote_write_url` - The remote write address of the Prometheus instance.

* `created_at` - The creation time of the prometheus instance, in RFC3339 format.

## Import

The prometheus instance can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_aom_prom_instance.test <id>
```
