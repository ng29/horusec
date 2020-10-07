// Copyright 2020 ZUP IT SERVICOS EM TECNOLOGIA E INOVACAO SA
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package horusec

import (
	"encoding/json"
	"time"

	"github.com/ZupIT/horusec/development-kit/pkg/enums/severity"

	"github.com/ZupIT/horusec/development-kit/pkg/enums/horusec"
	"github.com/google/uuid"
)

type Analysis struct {
	ID                      uuid.UUID                 `json:"id" gorm:"Column:analysis_id"`
	RepositoryID            uuid.UUID                 `json:"repositoryID" gorm:"Column:repository_id"`
	RepositoryName          string                    `json:"repositoryName" gorm:"Column:repository_name"`
	CompanyID               uuid.UUID                 `json:"companyID" gorm:"Column:company_id"`
	CompanyName             string                    `json:"companyName" gorm:"Column:company_name"`
	Status                  horusec.Status            `json:"status" gorm:"Column:status"`
	Errors                  string                    `json:"errors" gorm:"Column:errors"`
	CreatedAt               time.Time                 `json:"createdAt" gorm:"Column:created_at"`
	FinishedAt              time.Time                 `json:"finishedAt" gorm:"Column:finished_at"`
	AnalysisVulnerabilities []AnalysisVulnerabilities `json:"analysisVulnerabilities" gorm:"foreignkey:AnalysisID;association_foreignkey:ID"` //nolint:lll gorm usage
}

func (a *Analysis) GetTable() string {
	return "analysis"
}

func (a *Analysis) ToBytes() []byte {
	bytes, _ := json.Marshal(a)
	return bytes
}

func (a *Analysis) GetID() uuid.UUID {
	return a.ID
}

func (a *Analysis) GetIDString() string {
	return a.ID.String()
}

func (a *Analysis) ToString() string {
	return string(a.ToBytes())
}

func (a *Analysis) Map() map[string]interface{} {
	return map[string]interface{}{
		"id":                      a.ID,
		"createdAt":               a.CreatedAt,
		"repositoryID":            a.RepositoryID,
		"repositoryName":          a.RepositoryName,
		"companyName":             a.CompanyName,
		"companyID":               a.CompanyID,
		"status":                  a.Status,
		"errors":                  a.Errors,
		"finishedAt":              a.FinishedAt,
		"analysisVulnerabilities": a.AnalysisVulnerabilities,
	}
}

func (a *Analysis) SetFindOneFilter() map[string]interface{} {
	return map[string]interface{}{"id": a.GetID()}
}

func (a *Analysis) SetAnalysisError(err error) {
	if err != nil {
		toAppend := ""
		if len(a.Errors) > 0 {
			a.Errors += ", " + err.Error()
			return
		}
		a.Errors += toAppend + err.Error()
	}
}

func (a *Analysis) SetupIDInAnalysisContents() *Analysis {
	for key := range a.AnalysisVulnerabilities {
		a.AnalysisVulnerabilities[key].SetCreatedAt()
		a.AnalysisVulnerabilities[key].SetAnalysisID(a.ID)
		a.AnalysisVulnerabilities[key].SetVulnerabilityID(uuid.New())
	}
	return a
}

func (a *Analysis) SetCompanyName(companyName string) *Analysis {
	a.CompanyName = companyName
	return a
}

func (a *Analysis) SetRepositoryName(repositoryName string) *Analysis {
	a.RepositoryName = repositoryName
	return a
}

func (a *Analysis) SetRepositoryID(repositoryID uuid.UUID) *Analysis {
	a.RepositoryID = repositoryID
	return a
}

func (a *Analysis) SetAnalysisFinishedData() *Analysis {
	a.FinishedAt = time.Now()

	if a.HasErrors() {
		a.Status = horusec.Error
		return a
	}

	a.Status = horusec.Success
	return a
}

func (a *Analysis) HasErrors() bool {
	return len(a.Errors) > 0
}

func (a *Analysis) GetTotalVulnerabilities() int {
	return len(a.AnalysisVulnerabilities)
}

func (a *Analysis) GetTotalVulnerabilitiesBySeverity() (total map[severity.Severity]int) {
	total = map[severity.Severity]int{
		severity.High:   0,
		severity.Medium: 0,
		severity.Low:    0,
		severity.Info:   0,
		severity.Audit:  0,
		severity.NoSec:  0,
	}
	for index := range a.AnalysisVulnerabilities {
		total[a.AnalysisVulnerabilities[index].Vulnerability.Severity]++
	}
	return total
}

func (a *Analysis) SortVulnerabilitiesByCriticality() *Analysis {
	analysisVulnerabilities := a.getVulnerabilitiesBySeverity(severity.High)
	analysisVulnerabilities = append(analysisVulnerabilities, a.getVulnerabilitiesBySeverity(severity.Medium)...)
	analysisVulnerabilities = append(analysisVulnerabilities, a.getVulnerabilitiesBySeverity(severity.Low)...)
	analysisVulnerabilities = append(analysisVulnerabilities, a.getVulnerabilitiesBySeverity(severity.Info)...)
	analysisVulnerabilities = append(analysisVulnerabilities, a.getVulnerabilitiesBySeverity(severity.Audit)...)
	analysisVulnerabilities = append(analysisVulnerabilities, a.getVulnerabilitiesBySeverity(severity.NoSec)...)
	a.AnalysisVulnerabilities = analysisVulnerabilities
	return a
}

func (a *Analysis) GetAnalysisWithoutAnalysisVulnerabilities() *Analysis {
	return &Analysis{
		ID:             a.ID,
		RepositoryID:   a.RepositoryID,
		RepositoryName: a.RepositoryName,
		CompanyID:      a.CompanyID,
		CompanyName:    a.CompanyName,
		Status:         a.Status,
		Errors:         a.Errors,
		CreatedAt:      a.CreatedAt,
		FinishedAt:     a.FinishedAt,
	}
}

func (a *Analysis) getVulnerabilitiesBySeverity(search severity.Severity) (response []AnalysisVulnerabilities) {
	for index := range a.AnalysisVulnerabilities {
		if a.AnalysisVulnerabilities[index].Vulnerability.Severity == search {
			response = append(response, a.AnalysisVulnerabilities[index])
		}
	}
	return response
}

func (a *Analysis) SetDefaultVulnerabilityStatusAndType() *Analysis {
	for key := range a.AnalysisVulnerabilities {
		a.AnalysisVulnerabilities[key].Vulnerability.Type = horusec.Vulnerability
	}
	return a
}

func (a *Analysis) SetFalsePositivesAndRiskAcceptInVulnerabilities(
	listFalsePositive, listRiskAccept []string) *Analysis {
	for key := range a.AnalysisVulnerabilities {
		a.setVulnerabilityType(key, listFalsePositive, horusec.FalsePositive)
		a.setVulnerabilityType(key, listRiskAccept, horusec.RiskAccepted)
	}
	return a
}

func (a *Analysis) setVulnerabilityType(keyAnalysisVulnerabilities int,
	listToCheck []string, vulnerabilityType horusec.VulnerabilityType) {
	currentHash := a.AnalysisVulnerabilities[keyAnalysisVulnerabilities].Vulnerability.VulnHash
	for _, flagVulnerabilityHash := range listToCheck {
		if currentHash == flagVulnerabilityHash {
			a.AnalysisVulnerabilities[keyAnalysisVulnerabilities].Vulnerability.Type = vulnerabilityType
			break
		}
	}
}
