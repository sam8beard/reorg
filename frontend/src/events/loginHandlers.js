/* 
 * Handlers and utils for log in button on landing page
 */
import { showLoginPage } from '../navigation.js';

/*
 * Adds event for login button click on landing page.
 *
 * Shows login page
 */
export function attachLoginPageHandler(button, root) {
	button.addEventListener('click', () => showLoginPage(root));
}
/*
 * Adds event for login button on login page
 */
export function attachLoginHandler(button, root) {
	button.addEventListener('click', () => onLoginClick(root));
}

/*
 * Handles login click
 *
 * Fetches user data and uses state to display homepage
 */
export function onLoginClick(root) {
	// grab login field information from root
	//
	// fetch user data (and state?) 
	//
	// pass user data and state to showHomePage()
	//
	// showHomePage(root, userData);
}
