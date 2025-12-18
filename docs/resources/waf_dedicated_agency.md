---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_dedicated_agency"
description: |-
  Manages a WAF dedicated agency resource within HuaweiCloud.
---

# huaweicloud_waf_dedicated_agency

Manages a WAF dedicated agency resource within HuaweiCloud.

-> All WAF resources depend on WAF instances, and the WAF instances need to be purchased before they can be used.

## Example Usage

```hcl
variable "role_name_list" {
  type = list(string)
}

resource "huaweicloud_waf_dedicated_agency" "test" {
  role_name_list = var.role_name_list
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `role_name_list` - (Required, List) Specifies the list of dedicated engine agent policy names to create.
  Valid values are:
  + **evs_to_waf_operate_policy**
  + **vpc_to_waf_operate_policy**
  + **ecs_to_waf_operate_policy**
  + **elb_to_waf_operate_policy**

* `purged` - (Optional, Bool) Specifies whether to delete delegates synchronously.
  This field is valid only when destroying this resource. Defaults to **false**.
  When `purged` is set to **true**, the `premium_waf_svc_trust` delegate created in IAM will be deleted synchronously.
  When `purged` is set to **false**, the `premium_waf_svc_trust` delegate created in IAM will not be deleted synchronously.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (dedicated agency ID).

* `name` - The agency name.

* `version` - The version.

* `duration` - The agent existence time period.

* `domain_id` - The domain ID.

* `is_valid` - Whether the agency is legal.

* `role_list` - The list of dedicated engine agent policies.
  The [role_list](#agency_role_list) structure is documented below.

<a name="agency_role_list"></a>
The `role_list` block supports:

* `description` - The description.

* `catalog` - The catalog.

* `id` - The role ID.

* `display_name` - The display name.

* `is_granted` - Whether it is granted.

## Import

The resource can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_waf_dedicated_agency.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `role_name_list`, `purged`.
It is generally recommended running `terraform plan` after importing a resource.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to align
with the resource. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_waf_dedicated_agency" "test" {
  ...
  
  lifecycle {
    ignore_changes = [
      role_name_list,
      purged,
    ]
  }
}
```
