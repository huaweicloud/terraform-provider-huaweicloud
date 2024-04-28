---
subcategory: "Organizations"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_organizations_policy"
description: ""
---

# huaweicloud_organizations_policy

Manages an Organizations policy resource within HuaweiCloud.

## Example Usage

### Create service control policy

```hcl
resource "huaweicloud_organizations_policy" "scp_policy"{
  name    = "test_policy_name"
  type    = "service_control_policy"
  content = jsonencode(
    {
      "Version":"5.0",
      "Statement":[
        {
          "Effect":"Deny",
          "Action":[]
        }
      ]
    }
  )
}
```

### Create tag policy

```hcl
resource "huaweicloud_organizations_policy" "tag_policy"{
  name    = "test_policy_name"
  type    = "tag_policy"
  content = jsonencode(
    {
      "tags":{
        "test_tag":{
          "test_key":{
            "@@assign":"test_tag"
          }
        }
      }
    }
  )
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Specifies the name to be assigned to the policy. It can contain 1 to 64 characters, only
  English and Chinese letters, digits, underscore (_), hyphens (-) and spaces are allowed and the first and last
  characters cannot be spaces.

* `content` - (Required, String) Specifies the policy text content to be added to the new policy. For details, see the
  following documents:
  <br/> For service control policy: [documentation](https://support.huaweicloud.com/intl/en-us/usermanual-organizations/org_03_0033.html).
  <br/> For tag policy: [documentation](https://support.huaweicloud.com/intl/en-us/usermanual-organizations/org_03_0068.html).

* `type` - (Required, String, ForceNew) Specifies the type of the policy to be created. Value options:
  + **service_control_policy**: service control policy.
  + **tag_policy**: tag policy.

  Changing this parameter will create a new resource.

* `description` - (Optional, String) Specifies the description to be assigned to the policy. It can contain 1 to 512
  characters.

* `tags` - (Optional, Map) Specifies the key/value to attach to the policy.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `urn` - Indicates the uniform resource name of the policy.

## Import

The organizations policy can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_organizations_policy.test <id>
```
