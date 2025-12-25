package rules

import (
	"errors"
	"github.com/sam8beard/reorg/internal/models"
	"log"
	"slices"
	"strings"
)

type Target struct {
	TargetUUID string
	TargetName string
}

type BestMatch map[string]Target

type TargetMatches map[string]int

func Evaluate(ruleSet *models.RuleSet, fileMetadata map[string]models.FileMetadata) (*models.EvaluationResult, error) {
	// Debugging
	log.Printf("RuleSet received:\n%+v\n\n", ruleSet)
	log.Printf("FileMetadata received:\n%+v\n\n", fileMetadata)

	// Validate ruleset structure
	// NOTE: Im not even sure we need this, input is validated on the frontend
	// May be good for redundancy's sake
	if err := Validate(ruleSet); err != nil {
		return nil, err
	}

	// Grab upload UUID, files, and targets from ruleset
	uploadUUID := ruleSet.UploadUUID
	files := ruleSet.Files
	targets := ruleSet.Targets

	// Debugging
	log.Printf("Upload UUID: \n%s\n\n", uploadUUID)
	log.Printf("Files list: \n%+v\n\n", files)
	log.Printf("Targets list: \n%+v\n\n", targets)

	// Initialize evaluation result object
	evalResult := models.EvaluationResult{
		UploadUUID: uploadUUID,
		Folders:    make(map[string]*models.Folder, 0),
		Unmatched: models.UnmatchedFolder{
			Name:  "unsorted-files",
			Files: make([]models.File, 0),
		},
	}

	// Iterate through files
	for fileUUID := range files {
		// Retrieve metadata for current file using fileUUID
		md := fileMetadata[fileUUID]

		// Target condition matches for current file
		targetMatches := TargetMatches{}

		// Check every target and its associated rule against the current file
		for _, target := range targets {

			// Get total matches for current target
			if err := targetMatches.getMatches(target, md); err != nil {
				// No matches found for current target
				log.Printf("%v", err)
				// Check next target
				continue
			}
		}

		// Best matched target (greatest number of valid conditions)
		bestMatch := BestMatch{}

		// Get target to add file to
		bestMatchUUID, err := bestMatch.selectTarget(targets, targetMatches)

		// If this file does not match any target's rule
		if err != nil {
			unmatchedFile := models.File{
				FileUUID: md.FileUUID,
				FileName: md.FileName,
			}
			// Add file to unmatched/unsorted folder
			evalResult.Unmatched.Files = append(evalResult.Unmatched.Files, unmatchedFile)

			// Continue to next file
			continue
		}

		// Bind file to folder in evaluation result
		bestMatch.bindFile(&evalResult, bestMatchUUID, md)
	}

	// Return evaluation result
	LogEvalResult(&evalResult)
	return &evalResult, nil
}

/*
Logs an evaluation result
*/
func LogEvalResult(er *models.EvaluationResult) {
	folders := er.Folders
	for _, folder := range folders {
		folderName := folder.TargetName
		files := folder.Files
		log.Printf("Folder Name: %s\n", folderName)
		for _, file := range files {
			log.Printf("%s File Info:\n%+v", folderName, file)
		}

	}

	unmatched := er.Unmatched
	log.Print("Unmatched Folder Files:\n\n")
	for _, file := range unmatched.Files {
		log.Printf("\nFile Info:\n%+v", file)
	}
}

/*
Gets the total number of conditions matched for a given target and file.
Populates the calling object with the number of condition matches for a given target on success.

In order for a target to be considered a possible candidate, all active conditions for a target
must match for the given file.

Returns an error if a the given file fails to match at least one condition for the given target.
*/
func (tm *TargetMatches) getMatches(target models.Target, md models.FileMetadata) error {

	// Keep track of current target and its conditions matched
	currTargetUUID := target.TargetUUID
	(*tm)[currTargetUUID] = 0

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

	// Check conditions, if condition exists, check if there is a match
	for condition, active := range rule.ActiveConditions {
		switch condition {
		case "mime_type":
			if active {
				log.Println("condition detected")
				if checkFileType(conditions, md) {
					log.Println("match on file type")
					// Condition met, increment counter for target
					(*tm)[currTargetUUID]++
				} else {
					// Condition given and not met
					// Disqualify current target and check next target
					(*tm)[currTargetUUID] = 0
					return errors.New("failed match on file type")
				}
			}
		case "name_contains":
			if active {
				log.Println("condition detected")
				if checkNameContains(conditions, md) {
					log.Println("match on name contains")
					// Condition met, increment counter for target
					(*tm)[currTargetUUID]++
				} else {
					// Condition given and not met
					// Disqualify current target and check next target
					(*tm)[currTargetUUID] = 0
					return errors.New("failed match on name contains")
				}
			}
		case "size":
			if active {
				log.Println("condition detected")
				if checkSize(conditions, md) {
					log.Println("match on size")
					// Condition met, increment counter for target
					(*tm)[currTargetUUID]++
				} else {
					// Condition given and not met
					// Disqualify current target and check next target
					(*tm)[currTargetUUID] = 0
					return errors.New("failed match on size")
				}
			}
		case "created":
			if active {
				log.Println("condition detected")
				if checkCreated(conditions, md) {
					log.Println("match on created")
					// Condition met, increment counter for target
					(*tm)[currTargetUUID]++
				} else {
					// Condition given and not met
					// Disqualify current target and check next target
					(*tm)[currTargetUUID] = 0
					return errors.New("failed match on date creation")
				}
			}
		}
	}

	return nil
}

/*
Selects the best fit target using the target with the greatest
number of conditions matched.

Returns an error if the file had no condition matches for any target
*/
func (bm *BestMatch) selectTarget(targets map[string]models.Target, targetMatches map[string]int) (string, error) {
	// Keep track of most specific and valid target and the greatest amount of conditions matched for any valid target
	var bestMatchUUID string
	mostConditions := 0

	// Iterate through target matches
	for targetUUID, numMatched := range targetMatches {
		// Get target name
		targetName := targets[targetUUID].TargetName

		// Create temp target
		currentTarget := make(map[string]Target, 0)
		currentTarget[targetUUID] = Target{
			TargetUUID: targetUUID,
			TargetName: targetName,
		}

		// If the current target's specificity is greater than the greatest we've seen so far
		if numMatched > mostConditions {
			// Change the greatest we've seen so far to the current target's specificity
			mostConditions = numMatched

			// Select this target as the best matched so far
			(*bm)[targetUUID] = currentTarget[targetUUID]
			bestMatchUUID = targetUUID

		} else if numMatched == mostConditions && (numMatched != 0 && mostConditions != 0) && (currentTarget[targetUUID] != (*bm)[targetUUID]) {
			// TODO: WE MUST PREVENT THIS CASE ON RULE CREATION IN THE FRONTEND!!!!!
			//
			// When a rule is created, check if any other target's rule has the EXACT
			// same conditions
			//
			// If they do, notify the user and reject the rule creation

			log.Println("firing")
		}
	}

	// If there are no matches to any target's rule
	if len((*bm)) == 0 {
		return "", errors.New("no matches found")
	}

	return bestMatchUUID, nil

}

/* Binds a file to a folder in the evaluation result using the best fit target */
func (bm *BestMatch) bindFile(evalResult *models.EvaluationResult, bestMatchUUID string, md models.FileMetadata) {
	// Check if a folder corresponding to the best match target exists already
	// If not, create new folder
	if _, ok := evalResult.Folders[bestMatchUUID]; !ok {
		// Create new folder
		newFolder := models.Folder{
			TargetUUID: bestMatchUUID,
			TargetName: (*bm)[bestMatchUUID].TargetName,
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

	// Bind file to corresponding folder
	evalResult.Folders[bestMatchUUID].Files = append(evalResult.Folders[bestMatchUUID].Files, matchedFile)
}

/* Checks if file type condition is met */
func checkFileType(conditions models.Conditions, md models.FileMetadata) bool {

	// Mimetypes for target
	_, mimeTypes := conditions.Extensions, conditions.MimeType

	// NOTE: this logic is currently wrong, we need to check that EVERY mimetype given is
	// present in md.MimeType, at least I think its wrong, need to see how .Contains works
	return slices.Contains(mimeTypes, md.MimeType)
}

/* Checks if name contains condition is met */
func checkNameContains(conditions models.Conditions, md models.FileMetadata) bool {

	// Current file name and substring to match
	fileName, nameToMatch := md.FileName, conditions.NameContains

	// If the substring specified is in the file name
	return strings.Contains(fileName, nameToMatch)

}

/* Checks if size condition is met */
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

	// Get result using comparator, adjusted file size, and value choice
	switch compSelected {
	case Greater:
		matches = float64(fileSize) > adjustedSize
	case Less:
		matches = float64(fileSize) < adjustedSize
	}

	return matches
}

/* Checks if date creation condition is met */
func checkCreated(conditions models.Conditions, md models.FileMetadata) bool {

	matches := false

	// Creation date of file in miliseconds
	creationDate := md.OGTimestamp.UnixMilli()

	// Get data boundary options (default value is nil if omitted)
	beforeDate, afterDate := conditions.Created.Before, conditions.Created.After

	// Check if both a before date and after date were specified
	rangedBoundary := beforeDate != 0 && afterDate != 0

	// Evaluate based on range or single bound
	switch rangedBoundary {

	// Range evaluation
	case true:
		matches = (creationDate < beforeDate) && (creationDate > afterDate)

	// Single bound evaluation
	case false:

		// Comparator choices for single bounded option
		type Comparator int
		const (
			Greater Comparator = iota
			Less
		)
		// Get comparator option selected (greater than or less than)
		var compSelected Comparator
		if beforeDate != 0 {
			compSelected = Less
		} else {
			compSelected = Greater
		}

		// Evaluate based on comparator
		switch compSelected {
		case Less:
			matches = creationDate < beforeDate
		case Greater:
			matches = creationDate > afterDate
		}
	}

	return matches
}
