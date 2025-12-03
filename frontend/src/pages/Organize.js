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

			<div id='upload-files'>
				<!-- Need to decide what user data we're maintaining in database if using with an account -->
				
			</div>
		</div>

	`;
} 
