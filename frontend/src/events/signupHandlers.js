/* 
 * Handlers and utils for sign up button on landing page
 */
import { showSignupPage } from '../navigation.js';

/*
 * Adds event for sign up button click on landing page.
 *
 * Shows sign up page
 */
export function attachSignupPageHandler(button, root) {
	button.addEventListener('click',() => showSignupPage(root));
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
export function onSignupClick(root) { 

}
