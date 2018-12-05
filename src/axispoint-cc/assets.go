package main

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
	ROYALTYREPORT            string = "ROYALTYREPORT"
	EXPLOITATIONREPORT       string = "EXPLOITATIONREPORT"
	HOLDERREPRESENTATION     string = "HOLDERREPRESENTATION"
	ADMINISTRATORAFFILIATION string = "ADMINISTRATORAFFILIATION"
)

//ExploitationReport : struct defining data model for Exploitation Reports
type ExploitationReport struct {
	DocType                string `json:"docType"`
	Source                 string `json:"source"`
	SongTitle              string `json:"songTitle"`
	WriterName             string `json:"writerName"`
	Isrc                   string `json:"isrc"`
	Units                  int    `json:"units"`
	ExploitationDate       string `json:"exploitationDate"`
	Amount                 string `json:"amount"`
	UsageType              string `json:"usageType"`
	ExploitationReportUUID string `json:"exploitationReportUUID"`
	Territory              string `json:"territory"`
}

//RoyaltyReport : struct defining data model for Royalty Reports
type RoyaltyReport struct {
	DocType                string `json:"docType"`
	RoyaltyReportUUID      string `json:"royaltyReportUUID"`
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
	Target                 string `json:"target"`
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
