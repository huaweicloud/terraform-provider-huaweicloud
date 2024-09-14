---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_trust_agency"
description: |-
  Manages a trust agency resource within HuaweiCloud.
---
# huaweicloud_identity_trust_agency

Manages a trust agency resource within HuaweiCloud.

-> **NOTE:** You *must* have admin privileges to use this resource.

## Example Usage

### Delegate another HUAWEI CLOUD service to perform operations on your resources

```hcl
variable "agency_name" {}
variable "policy_name" {}

resource "huaweicloud_identity_trust_agency" "test" {
  name         = var.agency_name
  policy_names = [var.policy_name]
  description  = "test demo"
  trust_policy = jsonencode(
    {
      Statement = [
        {
          Action = [
            "sts:agencies:assume",
          ]
          Effect = "Allow"
          Principal = {
            Service = [
              "service.APIG",
            ]
          }
        },
      ]
      Version = "5.0"
    }
  )
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String, ForceNew) Specifies the name of trust agency. The name is a string of `1` to `64`
  characters. Only English letters, digits, underscores (_), plus (+), equals (=), commas (,), dots (.), ats (@) and
  hyphens (-) are allowed. Changing this will create a new trust agency.

* `trust_policy` - (Required, String, ForceNew) Specifies the trust policy of the trust agency. It's a JSON string.

* `policy_names` - (Required, List) Specifies a string list of one or more policy names that you would like to attach to
  the trust agency.

* `path` - (Optional, String, ForceNew) Specifies the resource path. It is made of several strings, each containing one
  or more English letters, digits, underscores (_), plus (+), equals (=), comma (,), dots (.), at (@) and hyphens (-),
  and must be ended with slash (/). Such as **foo/bar/**. It's a part of the uniform resource name. Default is empty.
  Changing this will create a new trust agency.

* `duration` - (Optional, Int) Specifies the validity period of a trust agency.
  Default value is `3,600`. The unit is seconds.

* `tags` - (Optional, Map) Specifies the tags of the trust agency.

* `description` - (Optional, String) Specifies the description of the trust agency.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The trust agency ID.

* `urn` - The uniform resource name of the trust agency. Format is `iam::$accountID:agency:$path$agencyName` where
  `$accountID` is IAM account ID, `$path` is `path`, `$agencyName` is `name`.

* `created_at` - The time when the trust agency was created.

## Import

Service agencies can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_identity_trust_agency.test <id>
```
