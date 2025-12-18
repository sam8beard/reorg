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
	} catch(err) {
		alert(err.message);
		return;
	}
	const response = await postRuleJson(ruleJson);
}

/* Builds the rule object from the form data on the rule creation page */
function buildRuleFromForm(formData) { 
	// Flags to keep track of file size input fields
	let sizeProvided = false;
	let comparatorProvided = false;
	let unitProvided = false;

	const ruleJson = { 
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
			"move_to": store.activeTarget.targetName
		}
	}
	for (var entry of formData.entries()) {
		var key = entry[0];
		var val = entry[1];
		switch(key) { 
			// handles the rule building for file type matching
			case 'fileType':
				console.log("Firing in ft block: " + entry);
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
				break;
			case 'nameContains':
				console.log("Firing in nc block: " + entry);
				ruleJson.when.name_contains = val.trim();
				break;

			case 'comparator':
				comparatorProvided = true;
				console.log("Firing in comp block: " + entry);
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
				console.log("Firing in size block: " + entry);
				console.log("Type of size value: " + typeof(val));
				
				break;

			case 'unit':
				unitProvided = true;
				console.log("Firing in unit block: " + entry);
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
