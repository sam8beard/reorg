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
 * Recursively read files from dropped directories
 */
function scanFiles(entry) { 
	if (entry.isDirectory) { 
		let dirReader = entry.createReader();
		dirReader.readEntries((entries) => {
			entries.forEach((entry) => { 
				scanFiles(entry);
			});
		});
	}
}

/*
 * Handle uploaded files via drop
 */
export async function dropHandler(e, root, preview) { 
	const items = e.dataTransfer.items;
	e.preventDefault();
	console.log(items);	
	for (const item of items) { 
		const entry = item.webkitGetAsEntry();
		if (entry) { 
			scanFiles(entry);
		}
	}
	// show preview
	displayFiles(items, preview);

	// store dropped files in file input for submission
	const input = root.querySelector("#dir-input");
	//const dt = new DataTransfer();
	//for (const file of items) dt.items.add(file);
	input.files = items.files;
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
	preview.innerText = "";
	for (const file of files) { 
		const li = document.createElement("li");
		li.appendChild(document.createTextNode(file.name));
		preview.appendChild(li);
	}
} 

/*
 * Handles file upload form submission
 */
export async function onFileSubmit(e, root) { 
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

	
	// what should the response contain? 
	const response = await uploadFileForm(formData);
	const preview = root.querySelector("#file-preview");

	// Clear preview and replace with message on submission
	if (response.error) { 
		preview.innerText = "Failed to upload files";
		form.reset();
	} else {
	// what page should we transition to or how should we modify the dom post upload? 
		preview.innerText = "Files uploaded";
		form.reset();
	}
	
	
	console.log(response);
}

