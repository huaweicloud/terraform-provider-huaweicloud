---
subcategory: "Tag Management Service (TMS)"
---

# huaweicloud_tms_api

Provides an API tms_predefineTags API resource.

## Example Usage

```hcl
//create
resource "huaweicloud_tms_tags" "test" {
	tag {
		key = "xxn"
		value = "11111"
	}
}

//update
resource "huaweicloud_tms_tags" "test" {
	tag {
		key = "nxx"
		value = "11111"
	}
}
```

## Argument Reference

The following arguments are supported:

* `tag` - (Required, List) A list of tag, each containing a key and value, MaxItems is 1.

The `tag` object supports the following:

* `key` - (Required, String) The key. The value contains a maximum of 36 characters. Character set: a-z, a-z, and 0-9, '-', '_', UNICODE character (\ u4E00 - \ u9FFF).
* `value` - (Required, String) Value. Each value contains a maximum of 43 characters and can be an empty string. Character set: AZ, a to z, 0-9, '. ', '-', '_', UNICODE character (\ u4E00 - \ u9FFF).

## Attributes Reference

No exported properties.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 3 minute.
* `read`   - Default is 3 minute.
* `update` - Default is 3 minute.
* `delete` - Default is 3 minute.
