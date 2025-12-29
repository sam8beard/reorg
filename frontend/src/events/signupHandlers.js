/* 
 * Handlers and utils for sign up button on landing page
 */
import { showSignup, showHome } from '../navigation.js';
import { createAccount } from '../api';
import { store } from '../state.js';

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
	button.addEventListener('submit', (e) => onSignupSubmit(root, e));
}
/*
 * Handles sign up submit
 *
 * Create account....etc. 
 */
export async function onSignupSubmit(root, e) {
	// Grab sign up field information from root
	e.preventDefault();
	const form = e.target;
	
	// TODO: VALIDATE USER INPUT
	const email = form.elements.email.value;
	const username = form.elements.username.value;
	const password = form.elements.password.value;

	const signupRequest = {
		"email": email,
		"username": username,
		"password": password,
	}
	
	const response = await createAccount(signupRequest);
	console.log(response);	

	if (response.err) {
		console.log(response);
	} else {
		console.log("Successfully registered user");
		console.log(response);
		sessionStorage.setItem('authToken', response.token);
		sessionStorage.setItem('userID', response.user.userID);
		store.identity.type = 'user';
		store.identity.userID = response.user.userID;
		store.identity.sessionID = response.token;
		console.log(store.identity);
		showHome(store.identity);
	}
}
