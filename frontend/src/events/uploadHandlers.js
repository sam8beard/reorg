/*
 * Handlers and utils for upload related actions
 */
import { uploadFileForm } from '../api';
import { showUpload } from '../navigation.js';
import { store } from '../state.js';
import { Tree, Folder, File } from 'https://cdn.jsdelivr.net/npm/@webreflection/file-tree/prod.js';
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
export async function dropHandler(e, root, previewContainer) { 
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
	console.log("DEBUG: ", dt.files);
	console.log("Files: ", input.files);
	displayFiles(dt.files, previewContainer)
}

/*
 * Handle uploaded files via default input
 */
export async function fileInputHandler(e, previewContainer) {
	displayFiles(e.target.files, previewContainer);
}

/*
 * Display preview for files uploaded
 */
function displayFiles(files, previewContainer) {
	const tree = new Tree;
	
	for (const file of files) { 
		tree.append(new File([], file.name));
	}
	previewContainer.appendChild(tree);
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
		alert("Please select at least one file to upload");
		return;
	}

	const formData = new FormData();
	for (const file of files) {
		formData.append(file.name, file, file.lastModified); 
	}
	
	// Get DOM elements to update
	const progressBar = root.querySelector("#upload-progress-bar");
	const progressBarContainer = root.querySelector("#progress-bar-container");
	const preview = root.querySelector("#file-preview");
	const previewContainer = root.querySelector('#upload-preview');
	const uploadBtn = root.querySelector("#upload-btn");
	const organizeBtn = root.querySelector("#organize-page-btn");
	const statusMessage = root.querySelector('#status-container');

	// Upload files
	uploadBtn.disabled = true;
	statusMessage.innerHTML = `
		<h3>Uploading files...</h3>
	`;
	
	progressBarContainer.style = 'width: 20em; background: white; height: 2em; border-radius: 4px;';
	progressBar.style.height = '100%';
	progressBar.style.width = '0%';
	progressBar.style.background = '#409940';
	progressBar.style.borderRadius = '4px';

	const response = await uploadFileForm(formData, progressBar);

	// Clear preview and replace with message on submission
	if (response.error) { 
		statusMessage.innerHTML = `
			<h3>Failed to upload files. Please try again.</h3>	
		`;
		progressBar.style.width = '0%';
		previewContainer.style.display = 'none';
		form.reset();
	} else {
		previewContainer.style.display = 'none';
		statusMessage.innerHTML = `
			<h3> Upload complete </h3>
		`;
		root.querySelector('#drop-zone').style.display = 'none';
		uploadBtn.style.display = 'none';
		uploadBtn.disabled = true;
		organizeBtn.style = 'display: block;';
		organizeBtn.disabled = false;

		// Use upload ID returned from backend response
		console.log("Firing in uploadHandlers: ", response);
		store.upload.uploadUUID = response;
		form.reset();
	}

	
	
}

