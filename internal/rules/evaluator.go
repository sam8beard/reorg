package rules

import (
	"github.com/sam8beard/reorg/internal/models"
	"log"
)

func Evaluate(ruleSet *models.RuleSet) (*models.EvaluationResult, error) {
	// Validate ruleset structure
	if err := Validate(ruleSet); err != nil {
		return nil, err
	}

	// Debugging
	log.Printf("Rule set received from frontend server:\n%+v", ruleSet)

	// Initialize evaluation result object
	evalResult := models.EvaluationResult{}
	/*
		Evaluation logic goes here...

		How should we go about building the evaluation result?

		This is the driving force behind the application
	*/

	// Grab upload UUID, files, and targets from ruleset
	uploadUUID := ruleSet.UploadUUID
	files := ruleSet.Files
	targets := ruleSet.Targets

	// Debugging
	log.Printf("Upload UUID: \n%s\n\n", uploadUUID)
	log.Printf("Files list: \n%+v\n\n", files)
	log.Printf("Targets list: \n%+v\n\n", targets)

	return &evalResult, nil
}
