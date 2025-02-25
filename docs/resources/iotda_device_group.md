---
subcategory: "IoT Device Access (IoTDA)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_iotda_device_group"
description: |-
  Manages an IoTDA device group resource within HuaweiCloud.
---

# huaweicloud_iotda_device_group

Manages an IoTDA device group resource within HuaweiCloud.

-> When accessing an IoTDA **standard** or **enterprise** edition instance, you need to specify the IoTDA service
  endpoint in `provider` block.
  You can login to the IoTDA console, choose the instance **Overview** and click **Access Details**
  to view the HTTPS application access address. An example of the access address might be
  **9bc34xxxxx.st1.iotda-app.ap-southeast-1.myhuaweicloud.com**, then you need to configure the
  `provider` block as follows:

  ```hcl
  provider "huaweicloud" {
    endpoints = {
      iotda = "https://9bc34xxxxx.st1.iotda-app.ap-southeast-1.myhuaweicloud.com"
    }
  }
  ```

## Example Usage

```hcl
variable "name" {}
variable "space_id" {}
variable "device_id" {}

resource "huaweicloud_iotda_device_group" "test" {
  name       = var.name
  space_id   = var.space_id
  device_ids = [var.device_id]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the IoTDA device group resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the name of device group. The name contains a maximum of `64` characters.
  Only letters, digits, hyphens (-) and underscores (_) are allowed.

* `space_id` - (Required, String, ForceNew) Specifies the resource space ID to which the device group belongs.
  Changing this parameter will create a new resource.

* `description` - (Optional, String) Specifies the description of device group.
  The description contains a maximum of `64` characters. Only letters, Chinese characters, digits, hyphens (-),
  underscores (_) and the following special characters are allowed: `?'#().,&%@!`.

* `parent_group_id` - (Optional, String, ForceNew) Specifies the parent group id.
  Changing this parameter will create a new resource.

* `type` - (Optional, String, ForceNew) Specifies the device group type.
  The valid values are as follows:
  + **STATIC**: Static device group.
  + **DYNAMIC**: Dynamic device group.

  Defaults to **STATIC**.

-> When the `type` parameter is set to **DYNAMIC**:
  <br/>1. The `dynamic_group_rule` parameter is required.
  <br/>2. The `parent_group_id` parameter cannot be set, because dynamic group do not support parent-child nesting.
  <br/>3. The `device_ids` does not need to be set and cannot be updated, only as an attribute. Because devices in
  dynamic group can only be managed automatically based on dynamic group rule.
  <br/>4. Only the standard and enterprise versions of IoTDA instances support dynamic device groups.

* `dynamic_group_rule` - (Optional, String, ForceNew) Specifies the dynamic device group rule, just fill in the content
  of the **where** clause, the remaining clauses do not need to be filled in.
  e.g. **device_name = 'xxx' or device_name = 'xxx'**.
  More grammar rules, please see [API docs](https://support.huaweicloud.com/intl/en-us/api-iothub/SearchDevices.html).

* `device_ids` - (Optional, List) Specifies the list of device IDs bound to the group.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

## Import

The device group can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_iotda_device_group.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `space_id`. It is generally
recommended running `terraform plan` after importing the resource. You can then decide if changes should be applied to
the resource, or the resource definition should be updated to align with the group. Also, you can ignore changes as
below.

```hcl
resource "huaweicloud_iotda_device_group" "test" {
    ...

  lifecycle {
    ignore_changes = [
      space_id,
    ]
  }
}
```
