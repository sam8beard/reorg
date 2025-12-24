import { store } from '../state.js';
import { postRuleJson, getPreviewJson } from '../api';
import { showPreview } from '../navigation.js';
import _ from 'lodash';

export async function onRuleSubmit(event, root) { 
	event.preventDefault();

	const form = event.currentTarget;

	// This allows us to get the form data as key:value pairs
	const formData = new FormData(form);

	// Build rule from form data and update preview
	let ruleJson; 
	let ruleSet;
	try { 
		// Build ruleJson and update state
		ruleJson = buildRuleFromForm(formData);

		// Set active rule in state
		store.activeRule = ruleJson.ruleUUID;

		// Add rule to state
		addRule(ruleJson);

		// Add rule binding to state
		addRuleBinding(ruleJson);
		
		// Build ruleSet and update state with preview
		ruleSet = buildRuleSet();

		// Debugging
		console.log(JSON.stringify(ruleSet, null, 2));

		// Get preview object from backend using ruleset object
		store.preview = await getPreviewJson(ruleSet);
	
		// Debugging 
		console.log(JSON.stringify(store.preview, null, 2));
		
		// Show preview page with updated preview
		showPreview();

	} catch(err) {
		alert(err.message);
		return;
	}
	
}

/* Builds ruleset object from current state */
function buildRuleSet() {
	const ruleSet = {
		"uploadUUID": store.upload.uploadUUID,
		"files": {},
		"targets": {}
	};
	
	// Populate files object
	for (let file of store.upload.files) {
		ruleSet.files[file.fileUUID] = file;
	}

	// Populate targets object using rule bindings
	for (let binding of store.ruleBindings) {
		const ruleUUID = binding.ruleUUID;
		const targetUUID = binding.targetUUID;
		console.log("Target UUID: ", targetUUID);
		console.log("ACTIVE target UUID: ", store.activeTarget.targetUUID);
		// Get rule referenced in binding (should be unique)
		let ruleToAdd;
		for (let rule of store.rules) {
			// We have found the matching rule
			if (rule.ruleUUID === ruleUUID) {
				ruleToAdd = rule;
				console.log("Rule to add set: ", ruleToAdd);
			}
		}

		// Get target referenced in binding (should be unique)
		let targetToAdd;
		for (let target of store.targets) {
			console.log("Firing in target loop");
			if (target.targetUUID === targetUUID) {
				targetToAdd = target;
				console.log("Target to add set: ", targetToAdd);
			}
		}
	
		// Make new target member with binded rule and target data
		const newTarget = {
			"targetUUID": targetToAdd.targetUUID,
			"targetName": targetToAdd.targetName,
			"rule": ruleToAdd
		}
		
		// Populate targets list in ruleset
		ruleSet.targets[newTarget.targetUUID] = newTarget;
	}
	
	return ruleSet;
}


/* Adds a new rule to user state */
function addRule(rule) { 
	store.rules = [...store.rules, rule];
}

/* Adds a new rule binding to user state */
function addRuleBinding(rule) {
	store.ruleBindings = [
		...store.ruleBindings,
		{ "ruleUUID": rule.ruleUUID, "targetUUID": store.activeTarget.targetUUID }
	];
}

/* Builds the rule object from the form data on the rule creation page */
function buildRuleFromForm(formData) { 
	// Flags to keep track of input fields
	let nameProvided = false;
	let fileTypeProvided = false;
	let sizeProvided = false;
	let comparatorProvided = false;
	let unitProvided = false;

	const ruleJson = { 
		"uploadUUID": store.upload.uploadUUID,
		"ruleUUID": crypto.randomUUID(),
		"ruleName": null,
		"activeConditions": {
			"extension": false, 
			"mime_type": false, 
			"name_contains": false,
			"size": false, 
			"created": false
		},
		"when": {
			"extension": [],
			"mime_type": [],
			"name_contains": "",
			"size": {
				"comparator": {
					"gt": false,
					"lt": false
				},

				"value": null,

				"unit": {
					"mb": false,
					"gb": false
				},
			},
			"created": {
				"before": null,
				"after": null
			},
		},
		"then": {
			"move_to": store.activeTarget.targetUUID
		}
	};

	// Iterate through form input and populate ruleJson
	for (var entry of formData.entries()) {
		var key = entry[0];
		var val = entry[1];
		switch(key) { 
			// handles the rule building for file type matching
			case 'ruleName':
				if (val.trim() !== '') { 
					nameProvided = true;
					ruleJson.ruleName = val.trim();
				}
				break;
			case 'fileType':
				if (val !== '') {
					fileTypeProvided = true;

					// Add active condition
					ruleJson.activeConditions.extension = true;
					ruleJson.activeConditions.mime_type = true;

					switch(val) {
						case 'image':
							ruleJson.when.extension.push(".jpg",".png",".svg",".gif");
							ruleJson.when.mime_type.push("image/jpeg","image/png","image/gif");
							break;
						case 'pdf':
							ruleJson.when.extension.push(".pdf");
							ruleJson.when.mime_type.push("application/pdf");
							break;
						case 'text':
							ruleJson.when.extension.push(".txt", ".md");
							ruleJson.when.mime_type.push("text/plain", "text/markdown");
							break;
						case 'video':
							ruleJson.when.extension.push(".mp4", ".mov", ".webm");
							ruleJson.when.mime_type.push("video/mp4", "video/webm", "video/quicktime");
							break;
						case 'document':
							ruleJson.when.extension.push(".doc", ".docx");
							ruleJson.when.mime_type.push("application/msword");
							break;
					}
				}
				break;
			case 'nameContains':
				const trimmed = val.trim();
				if (trimmed !== '') { 
					// Add active condition 
					ruleJson.activeConditions.name_contains = true;
					ruleJson.when.name_contains = trimmed;
				}
				
				break;

			case 'comparator':
				console.log("comp: ", val);
				comparatorProvided = true;
				switch(val) { 
					case 'lt':
						ruleJson.when.size.comparator.lt = true;
						break;
					
					case 'gt':
						ruleJson.when.size.comparator.gt = true;
						break;
				}
				break;

			case 'size':
				if (val !== '') {
					console.log("size: ", val);
					sizeProvided = true;
					ruleJson.when.size.value = Number(val);
				}
				
				break;

			case 'unit':
				console.log("unit: ", val);
				unitProvided = true;
				switch(val) { 
					case 'mb':
						ruleJson.when.size.unit.mb = true;
						break;
					
					case 'gb':
						ruleJson.when.size.unit.gb = true;
						break;
				}
				break;

			case 'dateBefore':
				if (val !== '') {
					const timestamp = Date.parse(val);
					ruleJson.when.created.after = timestamp;
					// Add active condition
					ruleJson.activeConditions.created = true;
				}
				break;

			case 'dateAfter':
				if (val !== '') {
					const timestamp = Date.parse(val);
					ruleJson.when.created.after = timestamp;
					// Add active condition
					ruleJson.activeConditions.created = true;
				}
				break;
		} 
	}
	
	// If rule name is not provided 
	if (!nameProvided) { 
		throw new Error('Rule name must be provided');
	}
	
	// If invalid file size input provided
	const anySizeField = sizeProvided || comparatorProvided || unitProvided;
	const allSizeField = sizeProvided && comparatorProvided && unitProvided;
	if (anySizeField && !allSizeField) {
		throw new Error('File size field requires a comparator, value, and unit'); 
	} else if (allSizeField) {
		// Add active condition
		ruleJson.activeConditions.size = true;
	}
	
	// If no conditions are provided
	const anyConditionProvided = ruleJson.activeConditions.extension || ruleJson.activeConditions.mime_type || ruleJson.activeConditions.name_contains || ruleJson.activeConditions.size || ruleJson.activeConditions.created;
	if (nameProvided && !anyConditionProvided) {
		throw new Error(`Please provide at least one condition for ${ruleJson.ruleName}`);
	}

	// If every condition in the submitted rule is already associated with a rule from the current upload
	if (isDuplicate(ruleJson)) {
		throw new Error(`This rule is already associated with another folder for this upload.`);
	}

	console.log(ruleJson);
	for (const key in ruleJson) { 
		console.log("Type of " + key + ": " + typeof(ruleJson[key]));
	}
	return ruleJson;
}

function isDuplicate(ruleJson) {
	// TODO: also check on folder name associated with rule
	// NOTE: will not allow duplicate folder names, but will allow 
	// 	 duplicate rule names AS LONG AS they are associated with different folders
	//
	//	This must be handled on folder creation.
	//
	//
	// Current rule conditions
	const newRuleConditions = ruleJson.when;
	// Upload UUID associated with rule
	const newRuleUploadUUID = ruleJson.uploadUUID;

	// Iterate through all rules
	for (let rule of store.rules) {

		let ruleUploadUUID = rule.uploadUUID;
		let ruleConditions = rule.when;
		
		// Rule is associated with the same upload
		if (newRuleUploadUUID === ruleUploadUUID) {
			// Conditions are deeply equal
			if (_.isEqual(newRuleConditions, ruleConditions)) return true;
		}
	}
	return false;
}
