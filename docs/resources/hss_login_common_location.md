---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_login_common_location"
description: |-
  Manages a login common location resource within HuaweiCloud HSS.
---

# huaweicloud_hss_login_common_location

Manages a login common location resource within HuaweiCloud HSS.

## Example Usage

```hcl
resource "huaweicloud_hss_login_common_location" "test" {
  area_code             = 86
  host_id_list          = ["host_id_1", "host_id_2"]
  enterprise_project_id = "all_granted_eps"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `area_code` - (Required, Int, NonUpdatable) Specifies the area code for the login location.

* `host_id_list` - (Required, List) Specifies the list of host IDs to apply the login location settings.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.
  This parameter only needs to be configured after the Enterprise Project feature is enabled.
  For enterprise users, if omitted, default enterprise project will be used.
  Value **0** means default enterprise project.
  Value **all_granted_eps** means all enterprise projects to which the user has been granted access.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

Login common location can be imported using the `enterprise_project_id` and `area_code` separated by a slash, e.g.

### Import resource under the default enterprise project

```bash
$ terraform import huaweicloud_hss_login_common_location.test 0/<id>
```

### Import resource from non default enterprise project

```bash
$ terraform import huaweicloud_hss_login_common_location.test <enterprise_project_id>/<area_code>
```
