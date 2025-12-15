/*
 * Handlers and utils for organize page related actions
 */
import { showOrganize, showRuleCreation } from '../navigation.js';
import { fetchFiles, postTarget } from '../api';
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
			
			<div id='create-dir-container'>
				<button id='create-dir-btn'>
					Create Folder
				</button>

			</div>

		</div>
	`;
	const createDirBtn = root.querySelector('#create-dir-btn');
	createDirBtn.addEventListener('click', () => onCreateDirClick(root));
}

async function onCreateDirClick(root) { 
	const createDirContainer = root.querySelector('#create-dir-container');	
	createDirContainer.innerHTML = `
		<form id='create-dir-form' method='post'>
			<label for='dir-name'>
				Folder Name: 
			</label>
			<input type='text' id='dir-name' name='dir' required placeholder='Folder Name'>
			<div>
				<button style='margin-top: 0.5em;' type='submit'>Create Folder</button>
			</div>
			<div>
				<p id='create-dir-error'></p>
			</div>
		</form>
	`;
	const createDirForm = createDirContainer.querySelector('#create-dir-form');
	createDirForm.addEventListener('submit', (e) => onCreateDirFormSubmit(e, createDirContainer));
}

async function onCreateDirFormSubmit(event, createDirContainer) { 
	event.preventDefault();
	const form = event.target;

	// generate id and name for new target
	const targetName = form.elements.dir.value.trim();
	if (targetName === "") { 
		const errMsg = createDirContainer.querySelector('#create-dir-error');
		errMsg.innerText = "Must supply folder name";
		return;
	}
	createDirContainer.querySelector('#create-dir-error').innerText = "";	
	const targetUUID = crypto.randomUUID();

	// Do we want to store the ID along with the UUID??
	store.targets.push({targetUUID, targetName});
	const target = { targetUUID: targetUUID, targetName: targetName}
	store.activeTarget = target
	// I don't think we want to post the target here yet...
	//const response = await postTarget(targetUUID, targetName);
	showRuleCreation()
}

