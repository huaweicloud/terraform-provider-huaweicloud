package assets

import (
	"github.com/chnsz/golangsdk"
)

type commonResult struct {
	golangsdk.Result
}

type ListAssetsResult struct {
	commonResult
}

func (r ListAssetsResult) Extract() ([]Data, error) {
	var s ListAssetsResponse
	err := r.ExtractInto(&s)
	if err != nil {
		return nil, err
	}
	return s.Data, nil
}

type ListAssetsResponse struct {
	// Error Code
	ErrCode string `json:"error_code"`
	// Error Message
	ErrMsg string `json:"error_msg"`
	// List of datas
	Data []Data `json:"data"`
}

type Data struct {
	// Asset ID
	AssetId string `json:"asset_id"`
	// Deployed Type
	DeployedType string `json:"deployed_type"`
	// Version
	Version string `json:"version"`
	// Version ID
	VersionId string `json:"version_id"`
	// Region
	Region string `json:"region"`
	// Image Deployed Object
	ImgDeployedObj ImageDeployedObj `json:"image_deployed_object"`
	// Software Pkg Deployed Object
	SwPkgDeployedObj SoftwarePkgDeployedObj `json:"software_pkg_deployed_object"`
}

type ImageDeployedObj struct {
	// Image ID
	ImageId string `json:"image_id"`
	// Image Name
	ImageName string `json:"image_name"`
	// OS Type
	OsType string `json:"os_type"`
	// Create Time
	CreateTime string `json:"create_time"`
	// Architecture
	Architecture string `json:"architecture"`
}

type SoftwarePkgDeployedObj struct {
	// Package Name
	PackageName string `json:"package_name"`
	// Internal Path
	InternalPath string `json:"internal_path"`
	// Checksum
	Checksum string `json:"checksum"`
}
