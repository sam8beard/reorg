import { store } from '../state.js';
export default function RuleCreation(root, user) {
	// Get most recently pushed target onto target list
	const currTarget = store.targets[store.targets.length - 1]
	root.innerHTML = `
		<div id='rule-creation-root'>
			<!-- Is this the optimal way to get the most recently added target dir??? -->
			<h1> Add rules for ${currTarget.targetName} </h1>
			<div id='rule-option-container'>
				<!-- Implement some external module for the multiselect dropdown -->
				<form id='rule-option-form' method='post'>
					<div>
						<label for='fileTypes'>File Type(s): </label>
						<select name='fileTypes' id='fileTypes' multiple data-multi-select>
							<option value="pdf">PDF (.pdf)</option>
							<option value="image">Images (.jpg, .png)</option>
							<option value="video">Videos (.mp3, .mov)</option>
							<option value="text">Text files (.txt, .md)</option>
						</select>
					</div>
					<div>
						<label for='nameContains'>Name Contains: </label>
						<input type='text' name='nameContains' id='nameContains'>
					</div>
					<div> <!-- Need to include greater than/less than options for file size -->
						<label for='size'>Size (specify MB or GB): </label>
						<input type='text' name='size' id='file-size'>
					</div>
					<div <!-- Need to include before/after options for data created -->
						<label for='date'>Date Created: </label>
						<input type='date' id='date-created'>
					</div>
				</form>
			</div>
		</div>
	`;
}
