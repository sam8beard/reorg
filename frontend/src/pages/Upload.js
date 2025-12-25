/*
 * Upload page
 */
import { store } from '../state.js';
import { dropHandler, fileInputHandler, onFileSubmit, onOrganizePageClick } from '../events';

export default function Upload(root, userData) { 
	console.log(userData)
	// Grab user data if logged in
	const user = (store.isLoggedIn) ? store.user : null;
	const username = (user) ? store.user.userID.username : "Guest";
	root.innerHTML = `
		<div id='upload-root'>
			<div style='margin-bottom: 5em;' id='upload-banner'>
				<h1> ${username} </h1>
			</div>

			<div  style='margin: 1em;' id='upload-form-container'>
				<!-- Need to decide what user data we're maintaining in database if using with an account 
				
				link to resource: https://developer.mozilla.org/en-US/docs/Web/API/File_API/Using_files_from_web_applications'
				
				drag and drop: https://developer.mozilla.org/en-US/docs/Web/API/HTML_Drag_and_Drop_API/File_drag_and_drop
'
				<form name='upload-form'>
					<div>
						<input id="upload-input" type='file' multiple />
					</div>
				</form>
				-->
				<form id='upload-form' name='upload-form'>
					<div>
						<label  style='border-radius: 0.5em; padding: 5em; margin: 1em; border: 0.05em white dashed' id='drop-zone'>
							Drop files here, or click to upload
							<input style="display: none;" id='dir-input' type='file' webkitdirectory directory multiple />
						</label>
					</div>

					<div style='display: flex; justify-content: center; gap: 1em;'>
						<button style='margin-top: 6em;' id='upload-btn' type='submit'>Upload</button>
						<button style='display:none;' id='organize-page-btn' type='button' disabled='true'>Organize Files</button>
					</div>

				</form> 

				<!-- 
				Here is where we will display the preview for all files 
				uploaded so far
				-->
				<div id='status-container'>
					<div id='status-msg'></div>
				</div>
				<div style='display: flex; justify-content: center; margin: 1em'>
					<div id='progress-bar-container'>
						<div id='upload-progress-bar'>
						</div>
					</div>
				</div>
				<div id='upload-preview'>
				</div>

			</div>
		</div>

	`;

	// Prevent default browser behavior when dragging and dropping files
	document.addEventListener("drop", (e) => { 
		if ([...e.dataTransfer.items].some((item) => item.kind === 'file')) {
			e.preventDefault();
		}
	});
	document.addEventListener("dragover", (e) => {
		if ([...e.dataTransfer.items].some((item) => item.kind === 'file')) {
			e.preventDefault();
		}
	});

	// Here, we'll deal with the behavior of dropping files for the drop target
	const dropZone = root.querySelector("#drop-zone");
	const preview = root.querySelector("#upload-preview");
	const uploadForm = root.querySelector("#upload-form");
	
	// Make sure dropHandler handles multi-file uploads
	dropZone.addEventListener("drop", (e) => dropHandler(e, root, preview));
	// Prevent default browser behavior for dragover event 
	dropZone.addEventListener("dragover", (e) => e.preventDefault());
	// Handle default input 
	dropZone.addEventListener("change", (e) => fileInputHandler(e, preview));
	// Handle form submission 
	uploadForm.addEventListener("submit", (e) => onFileSubmit(e, root));
	
	// Handle transition to Organize page upon successful upload
	const organizeBtn = root.querySelector('#organize-page-btn')
	organizeBtn.addEventListener("click", (e) => onOrganizePageClick(e, root));
}
