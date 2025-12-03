/*
 * Handlers and utils for organize related actions
 */
import { showOrganize } from '../navigation.js';
import { store } from '../state.js';
/*
 * Adds event for guest button click on landing page.
 *
 * Shows organize page
 */
export function attachOrganizePageHandler(button, root) {
	button.addEventListener('click', () => showOrganize());
}

/*
 * Handle uploaded files via drop
 */
export async function dropHandler(e, preview) { 
	e.preventDefault();
	const files = [...e.dataTransfer.items]
		.map((item) => item.getAsFile())
		.filter((file) => file);
	displayFiles(files, preview);
}

/*
 * Handle uploaded files via default input
 */
export async function fileInputHandler(e, preview) {
	displayFiles(e.target.files, preview);
}

/*
 * Display preview for files uploaded
 */
function displayFiles(files, preview) {
	for (const file of files) { 
		const li = document.createElement("li");
		li.appendChild(document.createTextNode(file.name));
		preview.appendChild(li);
	}
} 

/*
 * Handles file upload form submission
 */
export async function onFileSubmit(e) { 
	e.preventDefault();	
	/* Send files to backend endpoint */
	const form = e.target;
	
}

