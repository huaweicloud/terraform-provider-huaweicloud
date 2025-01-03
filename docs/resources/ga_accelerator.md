---
subcategory: "Global Accelerator (GA)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ga_accelerator"
description: ""
---

# huaweicloud_ga_accelerator

Manages a GA accelerator resource within HuaweiCloud.

## Example Usage

### Accelerator With IPV4

```hcl
variable "name" {}
variable "description" {}

resource "huaweicloud_ga_accelerator" "test" {
  name        = var.name
  description = var.description

  ip_sets {
    area = "CM"
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
```

#### Accelerator With IPV4 And IPV6

```hcl
variable "name" {}
variable "description" {}

resource "huaweicloud_ga_accelerator" "test" {
  name        = var.name
  description = var.description

  ip_sets {
    ip_type = "IPV4"
    area    = "CM"
  }

  ip_sets {
    ip_type = "IPV6"
    area    = "CM"
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Specifies the global accelerator name.  
  The name can contain `1` to `64` characters, only letters, digits, and hyphens (-) are allowed.

* `ip_sets` - (Required, List, ForceNew) Specifies the IP addresses assigned to the global accelerator.
  The [AccelerateIp](#Accelerator_AccelerateIp) structure is documented below.

  Changing this parameter will create a new resource.

* `description` - (Optional, String) Specifies the description about the global accelerator.  
  The description contain a maximum of `255` characters, and the angle brackets (< and >) are not allowed.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project ID of the tenant.
  The value is **0** or a string that contains a maximum of 36 characters in UUID format with hyphens (-).
  **0** indicates the default enterprise project. Defaults to **0**.

  Changing this parameter will create a new resource.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the global accelerator.

<a name="Accelerator_AccelerateIp"></a>
The `AccelerateIp` block supports:

* `area` - (Required, String, ForceNew) Specifies the acceleration area. The value can be one of the following:
  + **OUTOFCM**: Outside the Chinese mainland
  + **CM**: Chinese mainland

  Changing this parameter will create a new resource.

* `ip_type` - (Optional, String, ForceNew) Specifies the IP address version. Defaults to **IPV4**.
  Changing this parameter will create a new resource.
  The valid values are as follows:
  + **IPV4**
  + **IPV6**

  -> If you want to set this parameter to **IPV6**, you must set **IPV4** at the same time.
    Please refer to the document sample.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - Indicates the provisioning status. The value can be one of the following:
  + **ACTIVE**: The resource is running.
  + **PENDING**: The status is to be determined.
  + **ERROR**: Failed to create the resource.
  + **DELETING**: The resource is being deleted.

* `flavor_id` - Indicates the specification ID.

* `ip_sets` - Indicates the IP addresses assigned to the global accelerator.
  The [AccelerateIp](#Accelerator_AccelerateIpResp) structure is documented below.

* `frozen_info` - Indicates the frozen details of cloud services or resources.
  The [FrozenInfo](#Accelerator_FrozenInfo) structure is documented below.

* `created_at` - Indicates when the global accelerator was created.

* `updated_at` - Indicates when the global accelerator was updated.

<a name="Accelerator_AccelerateIpResp"></a>
The `AccelerateIp` block supports:

* `ip_address` - Indicates the IP address.

<a name="Accelerator_FrozenInfo"></a>
The `FrozenInfo` block supports:

* `effect` - Indicates the status of the resource after being frozen. The value can be one of the following:
  + **1** (default): The resource is frozen and can be released.
  + **2**: The resource is frozen and cannot be released.
  + **3**: The resource is frozen and cannot be renewed.

* `scene` - Indicates the service scenario. The value can be one of the following:
  + **ARREAR**: The cloud service is in arrears, including expiration of yearly/monthly resources and fee deduction
    failure of pay-per-use  resources.
  + **POLICE**: The cloud service is frozen for public security.
  + **ILLEGAL**: The cloud service is frozen due to violation of laws and regulations.
  + **VERIFY**: The cloud service is frozen because the user fails to pass the real-name authentication.
  + **PARTNER**: A partner freezes their customer's resources.

* `status` - Indicates the status of a cloud service or resource. The value can be one of the following:
  + **0**: unfrozen/normal (The cloud service will recover after being unfrozen.)
  + **1**: frozen (Resources and data will be retained, but the cloud service cannot be used.)
  + **2**: deleted/terminated (Both resources and data will be cleared.)

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

The accelerator can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_ga_accelerator.test <id>
```
