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
	"github.com/ZupIT/horusec/development-kit/pkg/databases/relational"
	"github.com/ZupIT/horusec/development-kit/pkg/databases/relational/repository/management"
	"github.com/ZupIT/horusec/development-kit/pkg/entities/api/dto"
	horusecEnums "github.com/ZupIT/horusec/development-kit/pkg/enums/horusec"
	"github.com/google/uuid"
)

type IController interface {
	GetAllVulnManagementData(repositoryID uuid.UUID, page, size int, vulnType horusecEnums.AnalysisVulnerabilitiesType,
		vulnStatus horusecEnums.AnalysisVulnerabilitiesStatus) (vulnManagement dto.VulnManagement, err error)
}

type Controller struct {
	managementRepository management.IManagementRepository
}

func NewManagementController(postgresRead relational.InterfaceRead,
	postgresWrite relational.InterfaceWrite) IController {
	return &Controller{
		managementRepository: management.NewManagementRepository(postgresRead, postgresWrite),
	}
}

func (c *Controller) GetAllVulnManagementData(repositoryID uuid.UUID, page, size int,
	vulnType horusecEnums.AnalysisVulnerabilitiesType,
	vulnStatus horusecEnums.AnalysisVulnerabilitiesStatus) (vulnManagement dto.VulnManagement, err error) {
	return c.managementRepository.GetAllVulnManagementData(repositoryID, page, size, vulnType, vulnStatus)
}