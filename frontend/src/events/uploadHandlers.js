/*
 * Handlers and utils for upload related actions
 */
import { uploadFileForm } from '../api';
import { showUpload } from '../navigation.js';
import { store } from '../state.js';
/*
 * Adds event for guest button click on landing page.
 *
 * Shows upload page
 */
export function attachUploadPageHandler(button, root) {
	button.addEventListener('click', () => showUpload());
}

/*
 * Fetches the file object, wraps it in a promise
 */
function getFilesAsync(entry) { 
	return new Promise((resolve, reject) => { 
		entry.file(resolve, reject);
	});
}

/*
 * Reads the directory, wraps it in a promise
 */
function readEntriesAsync(dirReader) {
	return new Promise((resolve, reject) => {
		dirReader.readEntries((entries) => { 
			resolve(entries), reject;
		});
	});
}

/*
 * Recursively read files from dropped directories
 */
async function scanFiles(entry, dt) { 
	if (entry.isDirectory) { 
		let dirReader = entry.createReader();
		while (true) { 
			const entries = await readEntriesAsync(dirReader);
			if (entries.length === 0) break; 
			for (const ent of entries) { 
				await scanFiles(ent, dt);
			} 
		} 
	} else {
		const file = await getFilesAsync(entry);
		dt.items.add(file);
	}
}

/*
 * Handle uploaded files via drop
 */
export async function dropHandler(e, root, preview) { 
	const items = e.dataTransfer.items;
	e.preventDefault();
	const dt = new DataTransfer();
	for (const item of items) { 
		const entry = item.webkitGetAsEntry();
		if (entry) {
			await scanFiles(entry, dt);
		}
	}
	console.log(dt.files.length);
	console.log(dt.items.length);
	console.log(dt.files);
	// show preview
	// displayFiles(dt.files, preview);
	for (const file of dt.files) { 
		console.log(file.name);
	} 
	// store dropped files in file input for submission
	const input = root.querySelector("#dir-input");
	//const dt = new DataTransfer();
	//for (const file of items) dt.items.add(file);
	input.files = dt.files;
	displayFiles(dt.files, preview)
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
	const uploadBtn = root.querySelector("#upload-btn");
	const organizeBtn = root.querySelector("#organize-page-btn");
	// Clear preview and replace with message on submission
	if (response.error) { 
		preview.innerText = "Failed to upload files";
		form.reset();
	} else {
	// what page should we transition to or how should we modify the dom post upload? 
		preview.innerText = "Files uploaded";
		uploadBtn.innerText = "Upload more files";
		organizeBtn.style.display = 'block';
		organizeBtn.disabled = false;
		form.reset();
	}
	
	
	console.log(response);
}

