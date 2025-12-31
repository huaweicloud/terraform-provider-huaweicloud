# The ID of the parent organizational unit (Required)
# This is required for both account enrollment and OU creation
# Replace with your actual parent organizational unit ID
parent_organizational_unit_id = "ou-xxxxxxxxxxxxx"

# The ID of the account to be enrolled with blueprint configuration
# Replace with your actual managed account ID
blueprint_managed_account_id            = "account-xxxxxxxxxxxxx"
# Blueprint product configuration
# Replace with your actual blueprint product ID and version
blueprint_product_id                    = "blueprint-xxxxxxxxxxxxx"
blueprint_product_version               = "1.0.0"
# Blueprint variables in JSON string format
# Customize these variables according to your blueprint requirements
blueprint_variables                     = "{\"environment\":\"production\",\"region\":\"cn-north-4\"}"
# Whether the blueprint has multi-account resources
is_blueprint_has_multi_account_resource = false
