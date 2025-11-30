/* 
 * Handlers and utils for sign up button on landing page
 */
import { showSignup } from '../navigation.js';

/*
 * Adds event for sign up button click on landing page.
 *
 * Shows sign up page
 */
export function attachSignupPageHandler(button, root) {
	button.addEventListener('click',() => showSignup());
}

/*
 * Adds event for sign up button on sign up page
 */
export function attachSignupHandler(button, root) {
	button.addEventListener('click', () => onSignupClick(root));
}
/*
 * Handles sign up click
 *
 * Create account....etc. 
 */
export async function onSignupClick(root) { 
	// Grab sign up field information from root
	// const username = root.querySelector('#username').value;
	// const password = root.querySelector('#password').value;
	//
	// Validate, check if user exists, create user, etc.
	// 
}
