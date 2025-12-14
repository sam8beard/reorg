/*
 * Organize page
 */
import { store } from '../state.js';
import { attachOrgAIHandler, attachOrgRuleHandler } from '../events';
export default function Organize(root, userData) { 
	root.innerHTML = `
		<div id='organize-root'>
			<div id='organize-options'>
				<button id='ai-organize-btn' type='click'>
					Organize with AI
				</button>
				<button id='rule-organize-btn' type='click'>
					Organize with Rules
				</button>
			</div>
			<div id='organize-file-preview-container'>
				<h2> Uploaded Files </h2>
				<!-- Fetch all files with the correct upload id and display here -->	
				<ul id='file-list' style='list-style-type: none;'></ul>
			</div>
		</div>
	`;

	// Display preview of uploaded files
	const fileList = root.querySelector('#file-list');
	for (const file of store.upload.files) { 
		const li = document.createElement("li");
		li.appendChild(document.createTextNode(file));
		fileList.appendChild(li);
	}

	// Attach event handlers for organize option buttons
	const orgAIBtn = root.querySelector('#ai-organize-btn');
	const orgRuleBtn = root.querySelector('#rule-organize-btn');
	attachOrgAIHandler(orgAIBtn, root)
	attachOrgRuleHandler(orgRuleBtn, root)
}
