package models

// DeployRequest represents the input parameters for GitHub Actions deployment
type DeployRequest struct {
	ImageTag    string `json:"imageTag"`
	ImageName   string `json:"imageName"`
	Replicas    string `json:"replicas"`
	Environment string `json:"environment"`
}

// ActionsRequestWrapper represents the wrapper for GitHub Actions deployment request
type ActionsRequestWrapper struct {
	Reference string        `json:"ref"`
	Inputs    DeployRequest `json:"inputs"`
}

// ImageDeployRequest represents the input parameters for image deployment
type ImageDeployRequest struct {
	ImageName   string `json:"imageName"`
	ImageTag    string `json:"imageTag"`
	CommitLink  string `json:"commitLink"`
	Environment string `json:"environment"`
}

// StatusRequest represents a deployment status update request
type StatusRequest struct {
	Status      string `json:"status"`
	ImageName   string `json:"imageName"`
	ImageTag    string `json:"imageTag"`
	Environment string `json:"environment"`
	FailedStep  string `json:"failedStep"`
	Logs        string `json:"logs"`
	JobURL      string `json:"jobURL"`
}
