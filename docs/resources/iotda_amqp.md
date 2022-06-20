---
subcategory: "IoT Device Access (IoTDA)"
---

# huaweicloud_iotda_amqp

Manages an IoTDA AMQP queue within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_iotda_amqp" "queue" {
  name = "queue_name"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the IoTDA AMQP queue resource.
If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the AMQP queue name, which contains 8 to 128 characters.
Only letters, digits, hyphens (-), underscores (_), dots (.) and colons (:) are allowed.
Changing this parameter will create a new resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

## Import

AMQP queues can be imported using the `id`, e.g.

```
$ terraform import huaweicloud_iotda_amqp.test 10022532f4f94f26b01daa1e424853e1
```
