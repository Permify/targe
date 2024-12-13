package requirements

import (
	"github.com/Permify/kivo/internal/requirements/aws"
)

type Requirement interface {
	GetFileName() string
	GetName() string
	Install() error
}

var requirements = []Requirement{
	aws.Types{},
	aws.ManagedPolicies{},
}

func GetRequirements() []Requirement {
	return requirements
}
