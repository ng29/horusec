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

package management

import (
	"github.com/ZupIT/horusec/development-kit/pkg/databases/relational/adapter"
	"github.com/ZupIT/horusec/development-kit/pkg/databases/relational/config"
	accountEntities "github.com/ZupIT/horusec/development-kit/pkg/entities/account"
	"github.com/ZupIT/horusec/development-kit/pkg/entities/account/roles"
	horusecEntities "github.com/ZupIT/horusec/development-kit/pkg/entities/horusec"
	rolesEnum "github.com/ZupIT/horusec/development-kit/pkg/enums/account"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestGetAllVulnManagementData(t *testing.T) {
	_ = os.Setenv(config.EnvRelationalDialect, "sqlite3")
	_ = os.Setenv(config.EnvRelationalURI, "tmp.db")
	_ = os.Setenv(config.EnvRelationalLogMode, "false")

	databaseWrite := adapter.NewRepositoryWrite()
	databaseRead := adapter.NewRepositoryRead()

	account := &accountEntities.Account{
		Email:     "test@test.com",
		Username:  "test",
		CreatedAt: time.Now(),
		Password:  "test",
		AccountID: uuid.New(),
	}

	company := &accountEntities.Company{
		CompanyID:   uuid.New(),
		Name:        "test",
		Description: "test",
		CreatedAt:   time.Now(),
	}

	repository := &accountEntities.Repository{
		RepositoryID: uuid.New(),
		CompanyID:    company.CompanyID,
		Name:         "test",
		CreatedAt:    time.Now(),
	}

	accountCompany := &roles.AccountCompany{
		AccountID: account.AccountID,
		CompanyID: company.CompanyID,
		Role:      rolesEnum.Admin,
		CreatedAt: time.Now(),
	}

	accountRepository := &roles.AccountRepository{
		AccountID:    account.AccountID,
		CompanyID:    company.CompanyID,
		RepositoryID: repository.RepositoryID,
		Role:         rolesEnum.Admin,
		CreatedAt:    time.Now(),
	}

	analysis := &horusecEntities.Analysis{
		ID:             uuid.New(),
		RepositoryID:   repository.RepositoryID,
		RepositoryName: "test",
		CompanyID:      company.CompanyID,
		CompanyName:    "test",
		Status:         "success",
		Errors:         "",
		CreatedAt:      time.Now(),
		FinishedAt:     time.Now(),
	}

	vulnerability := &horusecEntities.Vulnerability{
		VulnerabilityID: uuid.New(),
		Line:            "test",
		Column:          "test",
	}

	analysisVulnerabilities := &horusecEntities.AnalysisVulnerabilities{
		VulnerabilityID: vulnerability.VulnerabilityID,
		AnalysisID:      analysis.ID,
	}

	databaseWrite.SetLogMode(true)
	databaseWrite.GetConnection().Table(account.GetTable()).AutoMigrate(account)
	databaseWrite.GetConnection().Table(repository.GetTable()).AutoMigrate(repository)
	databaseWrite.GetConnection().Table(company.GetTable()).AutoMigrate(company)
	databaseWrite.GetConnection().Table(accountRepository.GetTable()).AutoMigrate(accountRepository)
	databaseWrite.GetConnection().Table(accountCompany.GetTable()).AutoMigrate(accountCompany)
	databaseWrite.GetConnection().Table(analysis.GetTable()).AutoMigrate(analysis)
	databaseWrite.GetConnection().Table(vulnerability.GetTable()).AutoMigrate(vulnerability)
	databaseWrite.GetConnection().Table(analysisVulnerabilities.GetTable()).AutoMigrate(analysisVulnerabilities)

	resp := databaseWrite.Create(account, account.GetTable())
	assert.NoError(t, resp.GetError())
	resp = databaseWrite.Create(company, company.GetTable())
	assert.NoError(t, resp.GetError())
	resp = databaseWrite.Create(repository, repository.GetTable())
	assert.NoError(t, resp.GetError())
	resp = databaseWrite.Create(accountRepository, accountRepository.GetTable())
	assert.NoError(t, resp.GetError())
	resp = databaseWrite.Create(accountCompany, accountCompany.GetTable())
	assert.NoError(t, resp.GetError())
	resp = databaseWrite.Create(analysis, analysis.GetTable())
	assert.NoError(t, resp.GetError())
	resp = databaseWrite.Create(vulnerability, vulnerability.GetTable())
	assert.NoError(t, resp.GetError())
	resp = databaseWrite.Create(analysisVulnerabilities, analysisVulnerabilities.GetTable())
	assert.NoError(t, resp.GetError())

	t.Run("should success get vulnerability data with no errors", func(t *testing.T) {
		repo := NewManagementRepository(databaseRead, databaseWrite)

		result, err := repo.GetAllVulnManagementData(repository.RepositoryID, 1, 1, "", "")

		assert.NoError(t, err)
		assert.Len(t, result.Data, 1)
	})

	t.Run("should success get vulnerability data with no errors", func(t *testing.T) {
		repo := NewManagementRepository(databaseRead, databaseWrite)

		result, err := repo.GetAllVulnManagementData(repository.RepositoryID, 1, 1, "test", "")

		assert.NoError(t, err)
		assert.Len(t, result.Data, 0)
	})

	t.Run("should success get vulnerability data with no errors", func(t *testing.T) {
		repo := NewManagementRepository(databaseRead, databaseWrite)

		result, err := repo.GetAllVulnManagementData(repository.RepositoryID, 1, 1, "", "test")

		assert.NoError(t, err)
		assert.Len(t, result.Data, 0)
	})

	t.Run("should success get vulnerability data with no errors", func(t *testing.T) {
		repo := NewManagementRepository(databaseRead, databaseWrite)

		result, err := repo.GetAllVulnManagementData(repository.RepositoryID, 1, 1, "test", "test")

		assert.NoError(t, err)
		assert.Len(t, result.Data, 0)
	})
}