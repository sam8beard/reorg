package rules

import (
	"github.com/sam8beard/reorg/internal/models"
	"log"
	"slices"
)

func Evaluate(ruleSet *models.RuleSet, fileMetadata map[string]models.FileMetadata) (*models.EvaluationResult, error) {
	// Validate ruleset structure
	if err := Validate(ruleSet); err != nil {
		return nil, err
	}

	// Debugging
	log.Printf("RuleSet received:\n%+v\n\n", ruleSet)
	log.Printf("FileMetadata received:\n%+v\n\n", fileMetadata)

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

	// Iterate through files

	for fileUUID := range files {
		// Retrieve metadata for current file using fileUUID
		md := fileMetadata[fileUUID]
		// Evaluate against rules in priority order?
		// Do we need a nested loop for this?
		// Any way to optimize to avoid N^2?
		//
		// Each target is associated with a rule
		for _, target := range targets {
			// Get rule associated with target
			rule := target.Rule

			// Get rule UUID and name
			ruleUUID, ruleName := rule.RuleUUID, rule.RuleName

			// Get matching conditions and actions to execute upon match
			conditions, actions := rule.Conditions, rule.Actions

			// Maybe check if conditions.Extensions/conditions.MimeType exists?
			//
			// If it does, we can use it to check against the file metadata,
			// and if it doesnt match, we can automatically
			//
			// Can now use rule.ActiveConditions to check conditions supplied
			// If file type is supplied by user, we check metadata to see if types are the same
			// If they arent, we can automatically rule out this file -> target binding
			if rule.ActiveConditions["extension"] && rule.ActiveConditions["mime_type"] {
				/*
				 * File type condition has been given by user
				 */

				// In this condition, we can cull every file that does not have the same type in the beginning, reducing sample size to match

				_, mimeTypes := conditions.Extensions, conditions.MimeType

				// Condition not met, continue to next target
				if !slices.Contains(mimeTypes, md.MimeType) {
					continue
				}

				// Condition met, set flag

			}

			if rule.ActiveConditions["name_contains"] {
				/*
				 * Name contains condition has been given by user
				 */
			}

			if rule.ActiveConditions["size"] {
				/*
				 * Size condition has been given by user
				 */
			}

			if rule.ActiveConditions["created"] {
				/*
				 * Created condition has been given by user
				 */
			}

		}

	}
	return &evalResult, nil
}
