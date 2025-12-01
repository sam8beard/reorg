/*
 * The Home Page
 */
import '../styles/home.css';
import { attachOrganizePageHandler } from '../events';
export default function Home(root, userData) { 
	console.log(userData);
	const userID = userData.userID;
	const username = userID.username;
	root.innerHTML = `
		<div id='home-root'>
			<!-- Render navbar component here -->
			<div id='home-welcome-banner'>
				<h1 id='home-welcome-msg'>
					Welcome, ${username}
				</h1>
			</div>
			<div id='home-options-container'>
				<button id='organize-files-btn'>Organize files</button>
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
	const organizeBtn = root.querySelector('#organize-files-btn');
	attachOrganizePageHandler(organizeBtn, root);
}

