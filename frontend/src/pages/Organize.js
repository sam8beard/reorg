/*
 * Organize page
 */
import { store } from '../state.js';
export default function Organize(root, userData) { 
	console.log(userData)
	// Grab user data if logged in
	const user = (store.isLoggedIn) ? store.user : null;
	const username = (user) ? store.user.userID.username : "Guest";
	root.innerHTML = `
		<div id='organize-root'>
			<h1 id='organize-banner'>
				${username}
			</h1>

			<div  id='upload-form-container'>
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
				<label  style='border: 0.05em white solid' id='drop-zone'>
					Drop files here, or click to upload.
					<input style="display: none;" id='file-input' type='file' multiple />
				</label>
				<!-- 
				Here is where we will display the preview for all files 
				uploaded so far
				-->
				<ul id='file-preview'></ul>

			</div>
		</div>

	`;
} 
