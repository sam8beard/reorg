/*
 * The Home Page
 */
import '../styles/home.css';
import { attachUploadPageHandler } from '../events';
export default function Home(root, identity) { 
	console.log(identity);
	root.innerHTML = `
		<div id='home-root'>
			<!-- Render navbar component here -->
			<div id='home-welcome-banner'>
				<h1 id='home-welcome-msg'>
				</h1>
			</div>
			<div id='home-options-container'>
				<button id='upload-files-btn'>Upload files</button>
			</div>
			<div id='user-info'>
				<div id='dashboard-header-container'>
					<h5>
						Files uploaded
					</h5>
				</div>
				<div id='home-dashboard'>
					<!-- What data/actions are we presenting to the user here? -->
				</div>
			</div>
		</div>
	`;

	// Attach event handlers
	const uploadBtn = root.querySelector('#upload-files-btn');
	attachUploadPageHandler(uploadBtn, root);
}

