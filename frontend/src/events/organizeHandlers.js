/*
 * Handlers and utils for organize page related actions
 */
import { showOrganize } from '../navigation.js';
import { fetchFiles } from '../api';
import { store } from '../state.js';

export async function onOrganizePageClick(e, root) {
	// Fetch files using upload ID for file preview on organize page
	const uploadId = store.upload.uploadID
	const response = await fetchFiles(uploadId);
	store.upload.files = response
	showOrganize();
}

/* Helper functions for attaching event listeners */
export function attachOrgAIHandler(btn, root) { 
	btn.addEventListener('click', () => onOrgAIClick(root));
}
export function attachOrgRuleHandler(btn, root) { 
	btn.addEventListener('click', () => onOrgRuleClick(root));
}

/* Displays organize with ai option */
async function onOrgAIClick(root) { 
	root.innerHTML = `
		<div id='ai-organize-root'>
			<h1> AI Organize </h1>
		</div>
	`;
}

/* Displays organize with rules option */
async function onOrgRuleClick(root) { 
	root.innerHTML = `
		<div id='rule-organize-root'>
			<h1> Rule organize </h1>
		</div>
	`;
} 
