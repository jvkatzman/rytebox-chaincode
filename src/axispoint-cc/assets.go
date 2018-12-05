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
	ROYALTYREPORT       string = "ROYALTYREPORT"
	EXPLOITATIONREPORT  string = "EXPLOITATIONREPORT"
	COPYRIGHTDATAREPORT string = "COPYRIGHTDATAREPORT"
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
