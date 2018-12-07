package main

import "time"

// Response -  Object to store Response Status and Message
// ================================================================================
type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

/////////////////////////////////////////////////////
// Constant for table names
/////////////////////////////////////////////////////
const (
	EXPLOITATIONREPORT       string = "EXPLOITATIONREPORT"
	HOLDERREPRESENTATION     string = "HOLDERREPRESENTATION"
	ADMINISTRATORAFFILIATION string = "ADMINISTRATORAFFILIATION"
	COPYRIGHTDATAREPORT      string = "COPYRIGHTDATAREPORT"
	ROYALTYSTATEMENT         string = "ROYALTYSTATEMENT"
)

/////////////////////////////////////////////////////
// Constant for the Exploitation Report State field values
/////////////////////////////////////////////////////
const (
	INITIAL                      string = "INITIAL"
	UNKNOWN_RIGHT_HOLDER         string = "UNKNOWN_RIGHT_HOLDER"
	INCONSISTENT_COPYRIGHT_SPLIT string = "INCONSISTENT_COPYRIGHT_SPLIT"
	INCOMPLETE_COPYRIGHT_SPLIT   string = "INCOMPLETE_COPYRIGHT_SPLIT"
)

//ExploitationReport : struct defining data model for Exploitation Reports
type ExploitationReport struct {
	DocType                string  `json:"docType"`
	Source                 string  `json:"source"`
	SongTitle              string  `json:"songTitle"`
	WriterName             string  `json:"writerName"`
	Isrc                   string  `json:"isrc"`
	Units                  int     `json:"units"`
	ExploitationDate       string  `json:"exploitationDate"`
	Amount                 float32 `json:"amount"`
	UsageType              string  `json:"usageType"`
	ExploitationReportUUID string  `json:"exploitationReportUUID"`
	Territory              string  `json:"territory"`
	State                  string  `json:"state"`
}

//RoyaltyStatement : struct defining data model for Royalty Reports
type RoyaltyStatement struct {
	DocType                string `json:"docType"`
	RoyaltyStatementUUID   string `json:"royaltyStatementUUID"`
	ExploitationReportUUID string `json:"exploitationReportUUID"`
	Source                 string `json:"source"`
	Isrc                   string `json:"isrc"`
	SongTitle              string `json:"songTitle"`
	WriterName             string `json:"writerName"`
	Units                  int    `json:"units"`
	ExploitationDate       string `json:"exploitationDate"`
	Amount                 string `json:"amount"`
	RightType              string `json:"rightType"`
	Territory              string `json:"territory"`
	UsageType              string `json:"usageType"`
	RightHolder            string `json:"rightHolder"`
	Administrator          string `json:"administrator"`
	Collector              string `json:"collector"`
	State                  string `json:"state"`
}

//CopyrightDataReport : struct definition
type CopyrightDataReport struct {
	DocType           string        `json:"docType"`
	CopyrightDataUUID string        `json:"copyrightDataReportUUID"`
	Isrc              string        `json:"isrc"`
	SongTitle         string        `json:"songTitle"`
	StartDate         time.Time     `json:"startDate,string"`
	EndDate           time.Time     `json:"endDate,string"`
	RightHolders      []RightHolder `json:"rightHolders"`
}

//RightHolder : struct definition for copyright data report
type RightHolder struct {
	IPI     string `json:"ipi"`
	Percent int    `json:"percent"`
}

//OwnerAdministration : struct defining data model for Owner Administration
type OwnerAdministration struct {
	DocType                 string           `json:"docType"`
	OwnerAdministrationUUID string           `json:"ownerAdministrationUUID"`
	Owner                   string           `json:"owner"`
	OwnerName               string           `json:"ownerName"`
	StartDate               string           `json:"startDate"`
	EndDate                 string           `json:"endDate"`
	Representations         []Representation `json:"representations"`
}

//Representation : struct defining data model for Representation
type Representation struct {
	Selector           string `json:"selector"`
	Representative     string `json:"representative"`
	RepresentativeName string `json:"representativeName"`
}

//AdministratorAffiliation : struct defining data model for Administrator Affiliation
type AdministratorAffiliation struct {
	DocType                      string        `json:"docType"`
	AdministratorAffiliationUUID string        `json:"administratorAffiliationUUID"`
	Administrator                string        `json:"administrator"`
	StartDate                    string        `json:"startDate"`
	EndDate                      string        `json:"endDate"`
	Affiliations                 []Affiliation `json:"affiliations"`
}

//Affiliation : struct defining data model for Affiliation
type Affiliation struct {
	Selector      string `json:"selector"`
	Affiliate     string `json:"affiliate"`
	AffiliateName string `json:"affiliateName"`
}
