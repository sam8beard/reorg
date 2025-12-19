import { store } from '../state.js';
import { postRuleJson } from '../api';

export async function onRuleSubmit(event, root) { 
	event.preventDefault();

	const form = event.currentTarget;

	// This allows us to get the form data as key:value pairs
	const formData = new FormData(form);
	// Build json object from rule input
	let ruleJson; 
	try { 
		ruleJson = buildRuleFromForm(formData);
		store.activeRule = ruleJson.ruleUUID;
		addRule(ruleJson);
		addRuleBinding(ruleJson);
	} catch(err) {
		alert(err.message);
		return;
	}

	// What should we present the user with after adding a rule? 
	//
	// Present preview option?
}

/* Adds a new rule to user state */
function addRule(rule) { 
	store.rules = [...store.rules, rule]
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
		"ruleUUID": crypto.randomUUID(),
		"ruleName": null,
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
				"before": "",
				"after": ""
			},
		},
		"then": {
			"move_to": store.activeTarget.targetUUID
		}
	}
	for (var entry of formData.entries()) {
		var key = entry[0];
		var val = entry[1];
		switch(key) { 
			// handles the rule building for file type matching
			case 'ruleName':
				if (val !== '') { 
					nameProvided = true;
					ruleJson.ruleName = val.trim();
				}
				break;
			case 'fileType':
				if (val !== '') {
					fileTypeProvided = true;
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
				ruleJson.when.name_contains = val.trim();
				break;

			case 'comparator':
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
					sizeProvided = true;
					ruleJson.when.size.value = Number(val);
				}
				
				break;

			case 'unit':
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
				ruleJson.when.created.before = val;
				break;

			case 'dateAfter':
				ruleJson.when.created.after = val;
				break;
		} 
	}
	
	// If rule name is not provided 
	if (!nameProvided) { 
		throw new Error('Rule name must be provided');
	}

	// If at least one file type is not provided
	if (!fileTypeProvided) { 
		throw new Error('At least one file type must be provided');
	}

	// If invalid file size input provided
	const anySizeField = sizeProvided || comparatorProvided || unitProvided;
	const allSizeField = sizeProvided && comparatorProvided && unitProvided 
	if (anySizeField && !allSizeField) {
		throw new Error('File size field requires a comparator, value, and unit'); 
	}
	console.log(ruleJson);
	for (const key in ruleJson) { 
		console.log("Type of " + key + ": " + typeof(ruleJson[key]));
	}
	return ruleJson;

}
