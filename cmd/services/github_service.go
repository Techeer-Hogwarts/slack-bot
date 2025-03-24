package services

type GitHubService interface {
	// Define methods for GitHub service
}

type githubService struct {
	githubURL   string
	githubToken string
}

func NewGitHubService(githubURL, githubToken string) GitHubService {
	return &githubService{
		githubURL:   githubURL,
		githubToken: githubToken,
	}
}
