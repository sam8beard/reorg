/*
 * Handlers and utils for organize related actions
 */
import { uploadFileForm } from '../api';
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
export async function dropHandler(e, root, preview) { 
	e.preventDefault();
	const files = [...e.dataTransfer.items]
		.map((item) => item.getAsFile())
		.filter((file) => file);
	// Show preview
	displayFiles(files, preview);

	// Store dropped files in file input for submission
	const input = root.querySelector("#dir-input");
	const dt = new DataTransfer();
	for (const file of files) dt.items.add(file);
	input.files = dt.files;
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
		li.appendChild(document.createTextNode(file.webkitRelativePath || file.name));
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
	const files = [
		...form.querySelector("#dir-input").files
	]
	
	console.log(e.target);
	// If no files were uploaded 
	if (!files || files.length === 0) { 
		// Display some message on DOM 
		// indicating no files were selected
		console.log("No files selected");
		return;
	}

	const formData = new FormData();
	for (const file of files) {
		const path = file.webkitRelativePath || file.name;
		console.log(path);
		formData.append("files", file, path); 
	} 

	
	const response = await uploadFileForm(formData);
	console.log(response);
}

