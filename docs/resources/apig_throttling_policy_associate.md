---
subcategory: "API Gateway (Dedicated APIG)"
---

# huaweicloud_apig_throttling_policy_associate

Use this resource to bind the APIs to the throttling policy within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "policy_id" {}
variable "api_publish_id1" {}
variable "api_publish_id2" {}

resource "huaweicloud_apig_throttling_policy_associate" "test" {
  instance_id = var.instance_id
  policy_id   = var.policy_id

  publish_ids = [
    var.api_publish_id1,
    var.api_publish_id2,
  ]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the API instance and throttling policy are located.
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the APIG dedicated instance to which the APIs and the
  throttling policy belongs. Changing this will create a new resource.

* `policy_id` - (Required, String, ForceNew) Specifies the ID of the API group to which the API response belongs to.
  Changing this will create a new resource.

* `publish_ids` - (Required, List) Specifies the publish ID corresponding to the API bound by the throttling policy.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Resource ID. The format is `<instance_id>/<policy_id>`.

## Import

Associate resources can be imported using their `policy_id` and the APIG dedicated instance ID to which the policy
belongs, separated by a slash, e.g.

```
$ terraform import huaweicloud_apig_throttling_policy_associate.test &ltinstance id&gt/&ltpolicy_id&gt
```
