import { store } from '../state.js';
import { postRuleJson } from '../api';

export async function onRuleSubmit(event, root) { 
	event.preventDefault();

	const form = event.currentTarget;

	// This allows us to get the form data as key:value pairs
	const formData = new FormData(form);
	// Build json object from rule input
	const ruleJson = buildRuleFromForm(formData);
	const response = await postRuleJson(ruleJson);
	
}

/* Builds the rule object from the form data on the rule creation page */
function buildRuleFromForm(formData) { 
	const ruleJson = { 
		"when": {
			"extension": [],
			"mime_type": [],
			"name_contains": "",
			"size": {
				"comparator": {
					"gt": false,
					"lt": false,
				},

				"value": null,

				"unit": {
					"mb": false,
					"gb": false,
				},
			},
			"created": {
				"before": null,
				"after": null,
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
				console.log("Firing in size block: " + entry);
				ruleJson.when.size.value = val;
				break;

			case 'unit':

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

	console.log(ruleJson);
	return ruleJson;

}
