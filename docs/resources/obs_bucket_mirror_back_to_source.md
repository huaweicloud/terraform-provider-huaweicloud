---
subcategory: "Object Storage Service (OBS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_obs_bucket_mirror_back_to_source"
description: |
  Manages an OBS bucket mirror back to source resource within HuaweiCloud.
---

# huaweicloud_obs_bucket_mirror_back_to_source

Manages an OBS bucket mirror back to source resource within HuaweiCloud.

## Example Usage

```hcl
variable "bucket_name" {}

# The rule ID must be unique in the bucket
resource "huaweicloud_obs_bucket_mirror_back_to_source" "test" {
  bucket = var.bucket_name
  rule   = jsonencode({
    "id" : "terraformtest",
    "condition" : {
      "httpErrorCodeReturnedEquals" : "404",
      "objectKeyPrefixEquals" : "test"
    },
    "redirect" : {
      "replaceKeyWith" : "test1$${key}test2",
      "publicSource" : {
        "sourceEndpoint" : {
          "master" : ["https://www.tftest1.com/xxx"],
          "slave" : ["https://www.tftest2.com/yyy"]
        }
      },
      "retryConditions" : ["4XX"],
      "returnBaseErrorConditions" : {
        "IFANY" : {
          "HTTP.STATUS_CODE" : ["4XX"]
        }
      },
      "agency" : try(data.huaweicloud_identity_agencies.test.agencies[0].name, ""),
      "ReturnBaseErrorBodyNewRule": true,
      "privateSource": {},
      "redirectHttpCode": "",
      "redirectOriginServer": "",
      "passQueryString" : true,
      "mirrorFollowRedirect" : true,
      "redirectWithoutReferer" : true,
      "mirrorHttpHeader" : {
        "passAll" : true
      }
    }
  })
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `bucket` - (Required, String, NonUpdatable) Specifies the name of the bucket to set mirror back to source.

* `rule` - (Required, String) Specifies the mirror back to source rules configuration, in JSON format.  
  For the rules syntax, please refer to the [documentation](https://support.huaweicloud.com/intl/en-us/api-obs/obs_04_0119.html).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which is the rule ID.

## Import

OBS bucket mirror back to source can be imported using `<bucket>/<id>` format, e.g.

```bash
$ terraform import huaweicloud_obs_bucket_mirror_back_to_source.test <bucket>/<rule_id>
```
