variable "iden_group_default_name" {
    type    = string
    default = "iden_group_default"
}

variable "iden_group_default_description" {
    type    = string
    default = "This is a default identity group."
}

variable "iden_group" {
    description = "List of identity group."
    type = list(object({
        name        = string
        description = string
    }))
    default = [
        {name = "iden_group_A", description = "This is a identity group."},
        {name = "iden_group_B", description = "This is a identity group."}
    ]
}