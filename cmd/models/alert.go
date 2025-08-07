package models

type FindMemberSchema struct {
	ID             int      `json:"id"`
	Type           string   `json:"type"`
	Name           string   `json:"name"`
	StudyExplain   string   `json:"studyExplain,omitempty"`
	ProjectExplain string   `json:"projectExplain,omitempty"`
	RecruitNum     int      `json:"recruitNum,omitempty"`
	FrontNum       int      `json:"frontNum,omitempty"`
	BackNum        int      `json:"backNum,omitempty"`
	DataEngNum     int      `json:"dataEngNum,omitempty"`
	DevOpsNum      int      `json:"devOpsNum,omitempty"`
	FullStack      int      `json:"fullStack,omitempty"`
	Leader         []string `json:"leader"`
	Email          []string `json:"email"`
	RecruitExplain string   `json:"recruitExplain"`
	NotionLink     string   `json:"notionLink"`
	Goal           string   `json:"goal,omitempty"`
	Rule           string   `json:"rule,omitempty"`
	Stack          []string `json:"stack,omitempty"`
	Environment    string   `json:"environment,omitempty"`
}

type UserMessageSchema struct {
	TeamID         int    `json:"teamId"`
	TeamName       string `json:"teamName"`
	Type           string `json:"type"`
	LeaderEmail    string `json:"leaderEmail"`
	ApplicantEmail string `json:"applicantEmail"`
	Result         string `json:"result"`
}

type AlertMessageSchema struct {
	Type      string `json:"type"`
	Email     string `json:"email"`
	ChannelID string `json:"channelId"`
	Message   string `json:"message"`
}
