package rules

import (
	"github.com/sam8beard/reorg/internal/models"
	"log"
	"slices"
	"strings"
)

type Target struct {
	TargetUUID string
	TargetName string
}

func Evaluate(ruleSet *models.RuleSet, fileMetadata map[string]models.FileMetadata) (*models.EvaluationResult, error) {
	// Validate ruleset structure
	if err := Validate(ruleSet); err != nil {
		return nil, err
	}

	// Debugging
	log.Printf("RuleSet received:\n%+v\n\n", ruleSet)
	log.Printf("FileMetadata received:\n%+v\n\n", fileMetadata)

	// Grab upload UUID, files, and targets from ruleset
	uploadUUID := ruleSet.UploadUUID
	files := ruleSet.Files
	targets := ruleSet.Targets

	// Initialize evaluation result object
	evalResult := models.EvaluationResult{
		UploadUUID: uploadUUID,
		Folders:    make(map[string]*models.Folder, 0),
		Unmatched: models.UnmatchedFolder{
			Name:  "unsorted-files",
			Files: make([]models.File, 0),
		},
	}

	// Debugging
	log.Printf("Upload UUID: \n%s\n\n", uploadUUID)
	log.Printf("Files list: \n%+v\n\n", files)
	log.Printf("Targets list: \n%+v\n\n", targets)

	// Iterate through files
	for fileUUID := range files {
		// Retrieve metadata for current file using fileUUID
		md := fileMetadata[fileUUID]

		// Keep track of matching targets and how many conditions are matched
		matchingTargets := make(map[string]int, 0)

		// Check every target and its associated rule against current file
		for _, target := range targets {

			// Keep track of current target and its conditions matched
			currTargetUUID := target.TargetUUID
			matchingTargets[currTargetUUID] = 0

			// Get rule associated with target
			rule := target.Rule

			// Get matching conditions and actions to execute upon match
			conditions, _ := rule.Conditions, rule.Actions

			/*
				NOTE:

				I just realized we are already keeping track of the folder to move files to
				in the rule.Actions member.

				Can we use this in any capacity? Or do we still need our current logic?
			*/

			if rule.ActiveConditions["extension"] && rule.ActiveConditions["mime_type"] {

				/*
				 * File type condition has been given by user
				 */

				matched := checkFileType(conditions, md)
				if matched {
					// Condition met, increment counter for target
					matchingTargets[currTargetUUID]++
				} else {
					// Condition given and not met
					// Disqualify current target and check next target
					matchingTargets[currTargetUUID] = 0
					continue
				}
			}

			if rule.ActiveConditions["name_contains"] {

				/*
				 * Name contains condition has been given by user
				 */

				matched := checkNameContains(conditions, md)
				if matched {
					matchingTargets[currTargetUUID]++
				} else {
					matchingTargets[currTargetUUID] = 0
					continue
				}
			}

			if rule.ActiveConditions["size"] {

				/*
				 * Size condition has been given by user
				 */

				matched := checkSize(conditions, md)
				if matched {
					matchingTargets[currTargetUUID]++
				} else {
					matchingTargets[currTargetUUID] = 0
					continue
				}
			}

			if rule.ActiveConditions["created"] {

				/*
				 * Created condition has been given by user
				 */

				matched := checkCreated(conditions, md)
				if matched {
					matchingTargets[currTargetUUID]++
				} else {
					matchingTargets[currTargetUUID] = 0
					continue
				}
			}

		}

		// Keep track of most specific and valid target and the greatest amount of conditions matched for any valid target
		bestMatch := make(map[string]Target, 0)
		var bestMatchUUID string
		mostConditions := 0

		// Decide which target receives the current file
		// by iterating through matchingTargets and getting
		// the number of conditions matched
		for targetUUID, numMatched := range matchingTargets {
			targetName := targets[targetUUID].TargetName
			currentTarget := make(map[string]Target, 0)
			currentTarget[targetUUID] = Target{
				TargetUUID: targetUUID,
				TargetName: targetName,
			}

			// If the current target's specificity is greater than the greatest we've seen so far
			if numMatched > mostConditions {
				// Change the greatest we've seen so far to the current target's specificity
				mostConditions = numMatched

				// Track this target
				bestMatch[targetUUID] = currentTarget[targetUUID]
				bestMatchUUID = targetUUID

			} else if numMatched == mostConditions && (numMatched != 0 && mostConditions != 0) && (currentTarget[targetUUID] != bestMatch[targetUUID]) {
				// NOTE: THIS LOGIC IS WRONG (I think?)
				// NOTE: this case is pretty shakey, not really sure how to check for this case
				// This is if there is a tie in the specificity of two targets
				//
				// What should be done?
				// Create two folders and copy file to both? I don't think it will be this easy
				// Create two folders if they dont

				log.Println("firing")
			}
		}

		// No matches were found for this file
		if len(bestMatch) == 0 {
			unmatchedFile := models.File{
				FileUUID: md.FileUUID,
				FileName: md.FileName,
			}
			// Add file to unmatched/unsorted folder
			evalResult.Unmatched.Files = append(evalResult.Unmatched.Files, unmatchedFile)

			// Continue to next file
			continue
		}

		// Check if a folder corresponding to the best match target exists already
		// If not, create new folder
		if _, ok := evalResult.Folders[bestMatchUUID]; !ok {
			// Create new folder
			newFolder := models.Folder{
				TargetUUID: bestMatchUUID,
				TargetName: bestMatch[bestMatchUUID].TargetName,
				Files:      make([]models.File, 0),
			}
			// Add folder to eval result
			evalResult.Folders[bestMatchUUID] = &newFolder

		}

		// Create matched file to add the matching folder
		matchedFile := models.File{
			FileUUID: md.FileUUID,
			FileName: md.FileName,
		}

		// Add file to corresponding folder
		evalResult.Folders[bestMatchUUID].Files = append(evalResult.Folders[bestMatchUUID].Files, matchedFile)

	}

	// Return evaluation result
	return &evalResult, nil
}
func checkFileType(conditions models.Conditions, md models.FileMetadata) bool {

	matches := false

	_, mimeTypes := conditions.Extensions, conditions.MimeType

	// NOTE: this logic is currently wrong, we need to check that EVERY mimetype given is
	// present in md.MimeType, at least I think its wrong, need to see how .Contains works
	if slices.Contains(mimeTypes, md.MimeType) {
		matches = true
	}
	return matches
}

func checkNameContains(conditions models.Conditions, md models.FileMetadata) bool {

	matches := false

	fileName, nameToMatch := md.FileName, conditions.NameContains

	// If the substring specified is in the file name
	if strings.Contains(fileName, nameToMatch) {
		matches = true
	}

	return matches
}

func checkSize(conditions models.Conditions, md models.FileMetadata) bool {

	matches := false

	// Size of file
	fileSize := md.Size

	// File size condition parameters
	sizeParams := conditions.FileSize
	compOptions, unitOptions := sizeParams.Comparator, sizeParams.Unit

	// Comparator choices
	type Comparator int
	const (
		Greater Comparator = iota
		Less
	)
	// Unit choices
	type Unit int
	const (
		MB Unit = iota
		GB
	)
	// Size value for comparison
	valueChoice := sizeParams.Value

	// Get comparator option selected (greater than or less than)
	var compSelected Comparator
	if compOptions.GreaterThan {
		compSelected = Greater
	} else {
		compSelected = Less
	}

	// Get unit option selected (megabytes or gigabytes)
	var unitSelected Unit
	if unitOptions.MB {
		unitSelected = MB
	} else {
		unitSelected = GB
	}

	// Get adjusted file size based on unit of measurement
	const (
		MBConv = 1024 * 1024
		GBConv = MBConv * 1024
	)
	var adjustedSize float64
	switch unitSelected {
	case MB:
		adjustedSize = float64(valueChoice) * MBConv
	case GB:
		adjustedSize = float64(valueChoice) * GBConv
	}

	// Build expression

	// Compare adjusted file size against given value choice
	switch compSelected {
	case Greater:
		matches = float64(fileSize) > adjustedSize
	case Less:
		matches = float64(fileSize) < adjustedSize
	}

	return matches
}

func checkCreated(conditions models.Conditions, md models.FileMetadata) bool {

	matches := false
	return matches
}
