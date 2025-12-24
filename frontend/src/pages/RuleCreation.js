import { store } from '../state.js';
import { onRuleSubmit } from '../events';
export default function RuleCreation(root, user) {
	// Get active target
	const currTarget = store.activeTarget
	root.innerHTML = `
		<div id='rule-creation-root'>
			<!-- Is this the optimal way to get the most recently added target dir??? -->
			<h1> Define the rule for files in this folder </h1>
			<br>
			<div id='rule-option-container'>
				<button>Use existing rule</button>
				<hr>
				<br>
				<!-- Implement some external module for the multiselect dropdown -->
				<h2>Create new rule</h2>
				<small>Files matching all selected conditions will be moved into <strong>${currTarget.targetName}</strong>.</small>
				<form id='rule-option-form' method='post'>
					<div>
						<label for='ruleName'><strong>Rule name</strong>
							<span style='color:red;'>  *</span>
							<br>
							<input type='text' id='rule-name' name='ruleName' required>
							<br>
							<details>
								<summary>?</summary>
								<small>Give this rule a name so it can be reused for other folders</small>
							</details>
						</label><br>
					</div>
					<div>
						<label for='fileType'><strong>File types include</strong> </label><br>
						<input type='checkbox' id='pdf-choice' name='fileType' value='pdf'>
						<label for='pdf'>PDF files (PDF)</label><br>
						<input type='checkbox' id='image-choice' name='fileType' value='image'>
						<label for='image'>Images (JPG, PNG)</label><br>
						<input type='checkbox' id='text-choice' name='fileType' value='text'>
						<label for='text'>Text files (TXT, MD)</label><br>
						<input type='checkbox' id='video-choice' name='fileType' value='video'>
						<label for='video'>Videos (MP4, MOV, WEBM)</label><br>
						<input type='checkbox' id='document-choice' name='fileType' value='document'>
						<label for='document'>Documents (DOCX, DOC)</label><br>
					</div>
					<div>
						<label for='nameContains'>
							<strong>File name contains text:</strong> 
						</label>
						<input type='text' name='nameContains' id='nameContains'>
						<details style='display:inline;'>
							<summary>
								<small>?</small>
							</summary>
							<small>Matches files whose names include this text</small>
							<small> (e.g. "invoice", "report", "screenshot")</small> 
						</details>
					</div>
					<div> <!-- Need to include greater than/less than options for file size -->
						<label><strong>File size must be:</strong>  </label>
						<select name='comparator' id='comparator'>
							<option value="" selected disabled hidden>Select</option>
							<option value='lt'>Smaller than</option>
							<option value='gt'>Larger than</option>
						</select>
						<input type='number' name='size' id='file-size'>
						<select name='unit' id='unit'>
							<option value="" selected disabled hidden>Select</option>
							<option value='mb'>MB</option>
							<option value='gb'>GB</option>
						</select>
					</div>
					<div> <!-- Need to include before/after options for data created -->
						<label for='dateBefore'>Created <strong>before</strong> this date: </label>
						<input name='dateBefore' type='date' id='date-before'>
					</div>
					<div>
						<label for='dateAfter'>Created <strong>after</strong> this date: </label>
						<input name='dateAfter' type='date' id='date-after'>
					</div>
					<div>
						<button id='rule-btn' type='submit'>Save rule and apply for ${currTarget.targetName}</button>
					</div>
				</form>
			</div>
		</div>
	`;
	const ruleForm = root.querySelector("#rule-option-form");
	ruleForm.addEventListener("submit", (e) => onRuleSubmit(e, root));

}
