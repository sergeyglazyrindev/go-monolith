package core

import "fmt"

type IValidator func(i interface{}, o interface{}) error

type ValidatorRegistry struct {
	Validators map[string]IValidator
}

func NewValidatorRegistry() *ValidatorRegistry {
	return &ValidatorRegistry{
		Validators: make(map[string]IValidator),
	}
}

func (vr *ValidatorRegistry) AddValidator(validatorName string, implementation IValidator) {
	_, exists := vr.Validators[validatorName]
	if exists {
		Trail(WARNING, "You are overriding validator %s", validatorName)
		return
	}
	vr.Validators[validatorName] = implementation
}

func (vr *ValidatorRegistry) GetValidator(validatorName string) (IValidator, error) {
	validator, exists := vr.Validators[validatorName]
	if !exists {
		return nil, fmt.Errorf("no %s validator registered", validatorName)
	}
	return validator, nil
}

func (vr *ValidatorRegistry) GetAllValidators() <-chan IValidator {
	chnl := make(chan IValidator)
	go func() {
		defer close(chnl)
		if vr == nil || vr.Validators == nil {
			return
		}
		for _, validator := range vr.Validators {
			chnl <- validator
		}
	}()
	return chnl
}

var GoMonolithValidatorRegistry *ValidatorRegistry

func init() {
	GoMonolithValidatorRegistry = &ValidatorRegistry{
		Validators: make(map[string]IValidator),
	}
}
