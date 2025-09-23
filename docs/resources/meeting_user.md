---
subcategory: "Meeting"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_meeting_user"
description: ""
---

# huaweicloud_meeting_user

Manages a meeting user resource within HuaweiCloud.

## Example Usage

### Create a user using third-party application authorization

```hcl
variable "app_id" {}
variable "app_key" {}
variable "user_name" {}
variable "user_password" {}
variable "english_name" {}
variable "signature" {}
variable "title" {}

resource "huaweicloud_meeting_user" "test" {
  app_id  = var.app_id
  app_key = var.app_key

  name         = var.user_name
  password     = var.user_password
  country      = "chinaPR"
  description  = "Created by script"
  email        = "123456789@example.com"`
  english_name = var.english_name
  phone        = "+8612345678987"
  signature    = var.signature
  title        = var.title
  sort_level   = 5
}
```

### Create a user using third-party application authorization and specifies the account parameters

```hcl
variable "app_id" {}
variable "app_key" {}
variable "user_account" {}
variable "third_account" {}
variable "user_name" {}
variable "user_password" {}
variable "english_name" {}
variable "signature" {}
variable "title" {}

resource "huaweicloud_meeting_user" "test" {
  app_id  = var.app_id
  app_key = var.app_key

  account       = var.user_account
  third_account = var.third_account
  name          = var.user_name
  password      = var.user_password
  country       = "chinaPR"
  description   = "Created by script"
  email         = "123456789@example.com"
  english_name  = var.english_name
  phone         = "+8612345678987"
  signature     = var.signature
  title         = var.title
  sort_level    = 5
}
```

### Create a user using account authorization

```hcl
variable "account_name" {}
variable "account_password" {}
variable "user_name" {}
variable "user_password" {}
variable "english_name" {}
variable "signature" {}
variable "title" {}

resource "huaweicloud_meeting_user" "test" {
  account_name     = var.account_name
  account_password = var.account_password

  name          = var.user_name
  password      = var.user_password
  country       = "chinaPR"
  email         = "123456789@example.com"
  english_name  = var.english_name
  phone         = "+8612345678987"
  signature     = var.signature
  title         = var.title
  sort_level    = 5
}
```

## Argument Reference

The following arguments are supported:

* `account_name` - (Optional, String, ForceNew) Specifies the (HUAWEI Cloud meeting) user account name to which the
  administrator belongs. Changing this parameter will create a new resource.

* `account_password` - (Optional, String, ForceNew) Specifies the user password.
  Required if `account_name` is set. Changing this parameter will create a new resource.

* `app_id` - (Optional, String, ForceNew) Specifies the ID of the Third-party application.
  Changing this parameter will create a new resource.

  -> You can apply for an application and obtain the App ID and App Key in the console.

* `app_key` - (Optional, String, ForceNew) Specifies the Key information of the Third-party APP.
  Required if `app_id` is set. Changing this parameter will create a new resource.

* `corp_id` - (Optional, String, ForceNew) Specifies the corporation ID.
  Required if the application is used in multiple enterprises. Only available if `app_id` is set.
  Changing this parameter will create a new resource.

* `user_id` - (Optional, String, ForceNew) Specifies the user ID of the administrator.
  Only available if `app_id` is set. If omitted, the user ID of default administrator will be used.
  Changing this parameter will create a new resource.

-> Exactly one of account authorization and application authorization you must select.

* `name` - (Required, String) Specifies the user name. The value can contain `1` to `64` characters.

* `password` - (Required, String) Specifies the user password.
  The following conditions must be met:
  + `8` to `32` characters.
  + It cannot be consistent with the positive and reverse order of the `account` parameter.
  + Contains at least two character types: lowercase letters, uppercase letters, numbers, special characters
    (**`~!@#$%^&*()-_=+|[{}];:",'<.>/?**).

* `account` - (Optional, String, ForceNew) Specifies the user account. The value can contain `1` to `64` characters.
  If omitted, the service will automatically generate a value.
  Changing this parameter will create a new resource.

* `third_account` - (Optional, String) Specifies the third-party account name.

* `department_code` - (Optional, String) Specifies the department code. Defaults to **1** (Root department).

* `description` - (Optional, String) Specifies the description. The value can contain `0` to `128` characters.

* `email` - (Optional, String) Specifies the email address.

* `english_name` - (Optional, String) Specifies the english name. The value can contain `0` to `64` characters.

* `country` - (Optional, String) Specifies the country to which the phone number belongs to.

* `phone` - (Optional, String) Specifies the phone number.
  The phone number must start with a [country (region) code](#phone_number_mapping).

* `hide_phone` - (Optional, Bool) Specifies whether to hide the phone number.

* `is_send_notify` - (Optional, Bool) Specifies whether to send email and SMS notifications for account opening.
  Defaults to **true**.

* `signature` - (Optional, String) Specifies the signature. The value can contain `0` to `512` characters.

* `status` - (Optional, Int) Specifies the status. The valid values are as follows:
  + **0**: Normal.
  + **1**: Disable.

  Defaults to `0`.

* `title` - (Optional, String) Specifies the title name. The value can contain `0` to `32` characters.

* `sort_level` - (Optional, Int) Specifies the address book sorting level.
  The lower the serial number, the higher the priority.
  The valid value is range from `1` to `10,000`. Defaults to `10,000`.

* `is_admin` - (Optional, Bool) Specifies whether to send email and SMS notifications for account opening.
  Defaults to **true**.

  -> Assign administrator role need to use the default administrator for authorization.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (user account).

* `sip_number` - The SIP number.

* `type` - The user type.
  + **2**: Enterprise member account.

* `department_name` - The department name.

* `department_name_path` - The department full name.

## Import

Users can be imported using their `id` and authorization parameters, separated by slashes, e.g.

Import a user and authenticated by account.

```bash
$ terraform import huaweicloud_meeting_user.test <id>/<account_name>/<account_password>
```

Import a user and authenticated by `APP ID`/`APP Key`.

```bash
$ terraform import huaweicloud_meeting_user.test <id>/<app_id>/<app_key>/<corp_id>/<user_id>
```

The slashes cannot be missing even corporation ID and user ID are empty.

Note that some parameters do not support import due to missing API responses or privacy, such as `password`,
`is_send_notify` and `is_admin`. You can ignore this change as below.

```hcl
resource "huaweicloud_meeting_user" "test" {
  ...

  lifecycle {
    ignore_changes = [
      password, is_send_notify,
    ]
  }
}
```

## Appendix

<a name="phone_number_mapping"></a>
The countries (or regions) and phone numbers mapping relationship supports:

| Country Or Region | Country Code |
| ---- | ---- |
| chinaPR | +86 |
| chinaHKG | +852|
| chinaOMA | +853 |
| chinaTPE | +886 |
| BVl | +1284  |
| Bolivia | +591 |
| CZ | +420 |
| GB | +245 |
| SVGrenadines | +1784 |
| TAT | +1868 |
| UK | +44 |
| afghanistan | +93 |
| albania | +355 |
| algeria | +213 |
| andorra | +376 |
| angola | +244 |
| argentina | +54 |
| armenia | +374 |
| australia | +61 |
| austria | +43 |
| azerbaijan | +994 |
| bahamas | +1242 |
| bahrain | +973 |
| bangladesh | +880 |
| belarus | +375 |
| belgium | +32 |
| belize | +501 |
| benin | +229 |
| bosniaAndHerzegovina | +387 |
| botswana | +267 |
| brazil | +55 |
| brunei | +673 |
| bulgaria | +359 |
| burkinaFaso | +226 |
| burundi | +257 |
| cambodia | +855 |
| cameroon | +237 |
| canada | +1 |
| capeVerde | +238 |
| caymanIslands | +1345 |
| centralAfrican | +236 |
| chad | +235 |
| chile | +56 |
| columbia | +57 |
| comoros | +269 |
| congoB | +242 |
| congoJ | +243 |
| costarica | +506 |
| croatia | +385 |
| curacao | +599 |
| cyprus | +357 |
| denmark | +45 |
| djibouti | +253 |
| dominica | +1809 |
| ecuador | +593 |
| egypt | +20 |
| equatorialGuinea | +240 |
| estonia | +372 |
| finland | +358 |
| france | +33 |
| gabon | +241 |
| gambia | +220 |
| georgia | +995 |
| germany | +49 |
| ghana | +233 |
| greece | +30 |
| grenada | +1473 |
| guatemala | +502 |
| guinea | +224 |
| guyana | +592 |
| honduras | +504 |
| hungary | +36 |
| india | +91 |
| indonesia | +62 |
| iraq | +964 |
| ireland | +353 |
| israel | +972 |
| italy | +39 |
| ivoryCoast | +225 |
| jamaica | +1876 |
| japan | +81 |
| jordan | +962 |
| kazakhstan | +7 |
| kenya | +254 |
| kuwait | +965 |
| kyrgyzstan | +996 |
| laos | +856 |
| latvia | +371 |
| lebanon | +961 |
| lesotho | +266 |
| liberia | +231 |
| libya | +218 |
| lithuania | +370 |
| luxembourg | +352 |
| macedonia | +389 |
| madagascar | +261 |
| malawi | +265 |
| malaysia | +60 |
| maldives | +960 |
| mali | +223 |
| malta | +356 |
| mauritania | +222 |
| mauritius | +230 |
| mexico | +52 |
| moldova | +373 |
| mongolia | +976 |
| montenegro | +382  |
| morocco | +212 |
| mozambique | +258 |
| myanmar | +95 |
| namibia | +264 |
| nepal | +977 |
| netherlands | +31 |
| newZealand | +64 |
| nicaragua | +505 |
| niger | +227 |
| nigeria | +234 |
| norway | +47 |
| oman | +968 |
| pakistan | +92 |
| palestine | +970 |
| panama | +507 |
| papuaNewGuinea | +675 |
| peru | +51 |
| philippines | +63 |
| poland | +48 |
| portugal | +351 |
| puertoRico | +1787 |
| qatar | +974 |
| romania | +40 |
| russia | +7 |
| rwanda | +250 |
| saintMartin | +590 |
| salvatore | +503 |
| saudiArabia | +966 |
| senegal | +221 |
| serbia | +381 |
| seychelles | +248 |
| sierraLeone | +232 |
| singapore | +65 |
| slovakia | +421 |
| slovenia | +386 |
| somalia | +252 |
| southAfrica | +27 |
| southKorea | +82 |
| spain | +34 |
| sriLanka | +94 |
| suriname | +597 |
| swaziland | +268 |
| sweden | +46 |
| switzerland | +41 |
| tajikistan | +992 |
| tanzania | +255 |
| thailand | +66 |
| togo | +228 |
| tunisia | +216 |
| turkey | +90 |
| turkmenistan | +993 |
| uae | +971 |
| uganda | +256 |
| ukraine | +380 |
| uruguay | +598 |
| usa | +1 |
| uzbekistan | +998 |
| venezuela | +58 |
| vietNam | +84 |
| yemen | +967 |
| zambia | +260 |
| zimbabwe | +263 |
