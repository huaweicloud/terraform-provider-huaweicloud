package auth

type CommonResp struct {
	IsSuccess bool   `json:"is_success"`
	Message   string `json:"message"`
}

type DataPermissions struct {
	IsSuccess  bool        `json:"is_success"`
	Message    string      `json:"message"`
	ObjectName string      `json:"object_name"`
	ObjectType string      `json:"object_type"`
	Count      int         `json:"count"`
	Privileges []Privilege `json:"privileges"`
}

type Privilege struct {
	IsAdmin    bool     `json:"is_admin"`
	UserName   string   `json:"user_name"`
	Privileges []string `json:"privileges"`
}

type QueuePermissions struct {
	IsSuccess  bool        `json:"is_success"`
	Message    string      `json:"message"`
	QueueName  string      `json:"queue_name"`
	Privileges []Privilege `json:"privileges"`
}

type DatabasePermissions struct {
	IsSuccess    bool        `json:"is_success"`
	Message      string      `json:"message"`
	DatabaseName string      `json:"database_name"`
	Privileges   []Privilege `json:"privileges"`
}

type TablePermissions struct {
	IsSuccess  bool             `json:"is_success"`
	Message    string           `json:"message"`
	Privileges []TablePrivilege `json:"privileges"`
}

type TablePrivilege struct {
	IsAdmin bool `json:"is_admin"`
	// Objects on which a user has permission.
	// If the object is in the format of databases.Database name.tables.Table name,
	// the user has permission on the database.
	// If the object is in the format of databases.Database name.tables.Table namecolumns.Column name,
	// the user has permission on the table.
	Object     string   `json:"object"`
	Privileges []string `json:"privileges"`
	UserName   string   `json:"user_name"`
}

type TablePermissionsOfUser struct {
	IsSuccess  bool                   `json:"is_success"`
	Message    string                 `json:"message"`
	UserName   string                 `json:"user_name"`
	Privileges []TablePrivilegeOfUser `json:"privileges"`
}

type TablePrivilegeOfUser struct {
	Object     string   `json:"object"`
	Privileges []string `json:"privileges"`
}
