---
subcategory: "Log Tank Service (LTS)"
---

# huaweicloud\_lts\_group

Manages a log group resource within HuaweiCloud.

## Example Usage

### create a log group

```hcl
resource "huaweicloud_lts_group" "log_group1" {
  group_name  = "log_group1"
  ttl_in_days = 1
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the log group resource. If omitted, the provider-level region will work as default. Changing this creates a new log group resource.

* `group_name` - (Required)
  Specifies the log group name.
  Changing this parameter will create a new resource.

* `ttl_in_days` - (Required)
  Specifies the log expiration time(days), value range: 1-30.

## Attributes Reference

The following attributes are exported:

* `id` - The log group ID.

* `group_name` - See Argument Reference above.

* `ttl_in_days` -
  Specifies the log expiration time(days).

## Import

Log group can be imported using the `id`, e.g.

```
$ terraform import huaweicloud_lts_group.group_1 7117d38e-4c8f-4624-a505-bd96b97d024c
```
