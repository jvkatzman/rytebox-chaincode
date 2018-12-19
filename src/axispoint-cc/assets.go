package main

// Response -  Object to store Response Status and Message
// ================================================================================
type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

/////////////////////////////////////////////////////
// RoyaltyStatementCreation event name
/////////////////////////////////////////////////////
const EventRoyaltyStatementCreation string = "RoyaltyStatementCreation"

/////////////////////////////////////////////////////
// Constant for table names
/////////////////////////////////////////////////////
const (
	EXPLOITATIONREPORT       string = "EXPLOITATIONREPORT"
	OWNERADMINISTRATION      string = "OWNERADMINISTRATION"
	ADMINISTRATORAFFILIATION string = "ADMINISTRATORAFFILIATION"
	COPYRIGHTDATAREPORT      string = "COPYRIGHTDATAREPORT"
	ROYALTYSTATEMENT         string = "ROYALTYSTATEMENT"
	COLLECTIONRIGHTSREPORT   string = "COLLECTIONRIGHTSREPORT"
)

/////////////////////////////////////////////////////
// Constant for the Exploitation Report State field values
/////////////////////////////////////////////////////
const (
	INITIAL                      string = "INITIAL"
	UNKNOWN_RIGHT_HOLDER         string = "UNKNOWN_RIGHT_HOLDER"
	INCONSISTENT_COPYRIGHT_SPLIT string = "INCONSISTENT_COPYRIGHT_SPLIT"
	INCOMPLETE_COPYRIGHT_SPLIT   string = "INCOMPLETE_COPYRIGHT_SPLIT"
	UNKNOWN                      string = "UNKNOWN"
	MISSING_COPYRIGHT_HOLDER     string = "MISSING_COPYRIGHT_HOLDER"
	MISSING_REPRESENTATIVE       string = "MISSING_REPRESENTATIVE"
	MISSING_AFFILIATE            string = "MISSING_AFFILIATE"
	UNKOWN_ISRC                  string = "UNKOWN_ISRC"
)

/////////////////////////////////////////////////////
// Constant for the Royalty Report Right Type
/////////////////////////////////////////////////////
const (
	OWNERSHIP  string = "OWNERSHIP"
	COLLECTION string = "COLLECTION"
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
	Amount                 float64 `json:"amount"`
	UsageType              string  `json:"usageType"`
	ExploitationReportUUID string  `json:"exploitationReportUUID"`
	Territory              string  `json:"territory"`
	State                  string  `json:"state"`
}

//RoyaltyStatement : struct defining data model for Royalty Reports
type RoyaltyStatement struct {
	DocType                string  `json:"docType"`
	RoyaltyStatementUUID   string  `json:"royaltyStatementUUID"`
	ExploitationReportUUID string  `json:"exploitationReportUUID"`
	Source                 string  `json:"source"`
	Isrc                   string  `json:"isrc"`
	SongTitle              string  `json:"songTitle"`
	WriterName             string  `json:"writerName"`
	Units                  int     `json:"units"`
	ExploitationDate       string  `json:"exploitationDate"`
	Amount                 float64 `json:"amount"`
	RightType              string  `json:"rightType"`
	Territory              string  `json:"territory"`
	UsageType              string  `json:"usageType"`
	RightHolder            string  `json:"rightHolder"`
	Administrator          string  `json:"administrator"`
	Collector              string  `json:"collector"`
	State                  string  `json:"state"`
}

//CopyrightDataReport : struct definition
type CopyrightDataReport struct {
	DocType           string        `json:"docType"`
	CopyrightDataUUID string        `json:"copyrightDataReportUUID"`
	Isrc              string        `json:"isrc"`
	SongTitle         string        `json:"songTitle"`
	StartDate         string        `json:"startDate"`
	EndDate           string        `json:"endDate"`
	RightHolders      []RightHolder `json:"rightHolders"`
}

//RightHolder : struct definition for copyright data report
type RightHolder struct {
	Selector string  `json:"selector"`
	IPI      string  `json:"ipi"`
	Percent  float64 `json:"percent"`
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

//CollectionRights : struct definition
type CollectionRight struct {
	DocType             string        `json:"docType"`
	CollectionRightUUID string        `json:"collectionRightUUID"`
	From                string        `json:"from"`     //EMI, Freddy, owner or admin.. --- also the key
	FromName            string        `json:"fromName"` //for display puposes
	StartDate           string        `json:"startDate"`
	EndDate             string        `json:"endDate"`
	RightHolders        []RightHolder `json:"rightHolders"`
}

//need to coordinate with MATT
//instead of managing
//ownershipAdministration
//AdministratorAffiliation

//rules need to move to collection rights
//they need to describe commission and generation
//generate as many files for 'EMI' to their partners
//depicts privacy at a high level

//RoyaltyStatementCreationEventPayload payload to passed as part of the event.
type RoyaltyStatementCreationEventPayload struct {
	Type                   string `json:"type"`
	TargetOrg              string `json:"targetOrg"`
	TargetIPI              string `json:"targetIPI"`
	ExploitationReportUUID string `json:"exploitationReportUUID"`
}
