import { store } from '../state.js';
import { onRuleSubmit } from '../events';
export default function RuleCreation(root, user) {
	// Get active target
	const currTarget = store.activeTarget
	root.innerHTML = `
		<div id='rule-creation-root'>
			<!-- Is this the optimal way to get the most recently added target dir??? -->
			<h4> Add rules for <strong>${currTarget.targetName}</strong>: </h4>
			<div id='rule-option-container'>
				<!-- Implement some external module for the multiselect dropdown -->
				<form id='rule-option-form' method='post'>
					<div>
						<label for='fileType'><strong>Move files of type:</strong> </label><br>
						<input type='checkbox' id='pdf-choice' name='fileType' value='pdf'>
						<label for='pdf'>PDF (.pdf)</label><br>
						<input type='checkbox' id='image-choice' name='fileType' value='image'>
						<label for='image'>Images (.jpg, .png)</label><br>
						<input type='checkbox' id='text-choice' name='fileType' value='text'>
						<label for='text'>Text files (.txt, .md)</label><br>
						<input type='checkbox' id='video-choice' name='fileType' value='video'>
						<label for='video'>Videos (.mp4, .mov, .webm)</label><br>
						<input type='checkbox' id='document-choice' name='fileType' value='document'>
						<label for='document'>Documents (.docx, .doc)</label><br>
					</div>
					<div>
						<label for='nameContains'><strong>File name contains:</strong> </label>
						<input type='text' name='nameContains' id='nameContains'>
					</div>
					<div> <!-- Need to include greater than/less than options for file size -->
						<label><strong>File size:</strong>  </label>
						<select name='comparator' id='comparator'>
							<option value="" selected disabled hidden>Select</option>
							<option value='lt'>less than</option>
							<option value='gt'>greater than</option>
						</select>
						<input type='number' name='size' id='file-size'>
						<select name='unit' id='unit'>
							<option value="" selected disabled hidden>Select</option>
							<option value='mb'>MB</option>
							<option value='gb'>GB</option>
						</select>
					</div>
					<div> <!-- Need to include before/after options for data created -->
						<label for='dateBefore'>Created <strong>before</strong>: </label>
						<input name='dateBefore' type='date' id='date-before'>
					</div>
					<div>
						<label for='dateAfter'>Created <strong>after</strong>: </label>
						<input name='dateAfter' type='date' id='date-after'>
					</div>
					<div>
						<button id='rule-btn' type='submit'>Add rules</button>
					</div>
				</form>
			</div>
		</div>
	`;
	const ruleForm = root.querySelector("#rule-option-form");
	ruleForm.addEventListener("submit", (e) => onRuleSubmit(e, root));

}
