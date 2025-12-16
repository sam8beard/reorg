import { store } from '../state.js';

export async function onRuleSubmit(event, root) { 
	event.preventDefault();

	const form = event.currentTarget;

	// This allows us to get the form data as key:value pairs
	const formData = new FormData(form);
	//for (var pair of formData.entries()) {
	//	console.log(pair[0] + ": " + pair[1]);
	//}
	buildRuleFromForm(formData);

	
}

/* Builds the rule object from the form data on the rule creation page */
function buildRuleFromForm(formData) { 
	const ruleJson = { 
		"when": {
			"extension": [],
			"mime_type": [],
			"name_contains": "",
			"size": {
				"gt": null,
				"lt": null,
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
				break;

			case 'comparator':
				console.log("Firing in comp block: " + entry);
				break;

			case 'size':
				console.log("Firing in size block: " + entry);
				break;

			case 'unit':
				console.log("Firing in unit block: " + entry);
				break;
		} 
	}

	console.log(ruleJson);

}
