package domains

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

type CreateResult struct {
	golangsdk.ErrResult
}

type GetResult struct {
	golangsdk.Result
}

type AccessDomain struct {
	Organization string `json:"namespace"`
	Repository   string `json:"repository"`
	// Name
	AccessDomain string `json:"access_domain"`
	Permit       string `json:"permit"`
	Deadline     string `json:"deadline"`
	Description  string `json:"description"`
	CreatorID    string `json:"creator_id"`
	CreatorName  string `json:"creator_name"`
	Created      string `json:"created"`
	Updated      string `json:"updated"`
	// Status. `true`: valid `false`: expired
	Status bool `json:"status"`
}

func (r GetResult) Extract() (*AccessDomain, error) {
	// existance check
	var existingDomain struct {
		Exist bool `json:"exist"`
	}
	if err := r.ExtractIntoStructPtr(&existingDomain, ""); err != nil {
		return nil, err
	}
	if !existingDomain.Exist {
		err404 := golangsdk.ErrDefault404{}
		err404.DefaultErrString = "domain does not exist"
		return nil, err404
	}

	// actual extract
	domain := new(AccessDomain)
	if err := r.ExtractIntoStructPtr(domain, ""); err != nil {
		return nil, err
	}

	return domain, nil
}

type DeleteResult struct {
	golangsdk.ErrResult
}

type UpdateResult struct {
	golangsdk.ErrResult
}

type AccessDomainPage struct {
	pagination.SinglePageBase
}

func ExtractAccessDomains(p pagination.Page) ([]AccessDomain, error) {
	var domains []AccessDomain
	err := p.(AccessDomainPage).ExtractIntoSlicePtr(&domains, "")
	return domains, err
}
