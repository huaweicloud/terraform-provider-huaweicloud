package repositories

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

type CreateResult struct {
	golangsdk.ErrResult
}

type DeleteResult struct {
	golangsdk.ErrResult
}

type UpdateResult struct {
	golangsdk.ErrResult
}

type GetResult struct {
	golangsdk.Result
}

type ImageRepository struct {
	ID             int    `json:"id"`
	OrganizationID int    `json:"ns_id"`
	Name           string `json:"name"`
	Category       string `json:"category"`
	Description    string `json:"description"`
	CreatorID      string `json:"creator_id"`
	CreatorName    string `json:"creator_name"`
	Size           int    `json:"size"`
	IsPublic       bool   `json:"is_public"`
	NumImages      int    `json:"num_images"`
	NumDownloads   int    `json:"num_downloads"`
	URL            string `json:"url"`
	Path           string `json:"path"`
	InternalPath   string `json:"internal_path"`
	Created        string `json:"created"`
	Updated        string `json:"updated"`
	DomainID       string `json:"domain_id"`
	Prioriry       int    `json:"poriroty"`
}

func (r GetResult) Extract() (*ImageRepository, error) {
	var repo ImageRepository
	err := r.ExtractInto(&repo)
	return &repo, err
}

type RepositoryPage struct {
	pagination.LinkedPageBase
}

func ExtractRepositories(p pagination.Page) ([]ImageRepository, error) {
	var repos []ImageRepository
	err := (p.(RepositoryPage)).ExtractIntoSlicePtr(&repos, "")
	if err != nil {
		return nil, err
	}
	return repos, nil
}
