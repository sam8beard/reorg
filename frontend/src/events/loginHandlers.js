/*
 * Handlers and utils for log in related actions
 */
import { showLogin, showHome } from '../navigation.js';
import { fetchUser } from '../api';
import { store } from '../state.js';

/*
 * Adds event for login button click on landing page.
 *
 * Shows login page
 */
export function attachLoginPageHandler(button, root) {
	button.addEventListener('click', () => showLogin());
}
/*
 * Adds event for login form on login page
 */
export function attachLoginHandler(form, root) {
	form.addEventListener('submit', (e) => {
		onLoginSubmit(root, e);
	});
}

/*
 * Handles login click
 *
 * Fetches user data and uses state to display homepage
 */
async function onLoginClick(root) {
	// Grab login field information from root
	// const username = root.querySelector('#username').value;
	// const password = root.querySelector('#password').value;
	//
	// Fetch user data from backend 
	//
	// pass user data and state to showHomePage()
	//
	// Update store and show page
	// showHome(userData);
}

/*
 * Handles login submit
 *
 * Fetches user data and uses state to display homepage
 */
async function onLoginSubmit(root, e) { 
	e.preventDefault();
	const form = e.target;
	const username = form.elements.username.value;
	const password = form.elements.password.value;
	
	const loginRequest = {
		"usernameOrEmail": username,
		"password": password,
	};
	
	// Fetch user from backend
	const response = await fetchUser(loginRequest);
	console.log(response);	

	// If user does not exist, update state to notify login page and rerender with message
	if (response.error) {
		store.loginError = response.error;
		return;
	// If user does exist, show home page for respective user
	} else {
		console.log("Successfully fetched existing user");
		console.log(response);
		sessionStorage.setItem('authToken', response.token);
		sessionStorage.setItem('userID', response.user.userID);
		store.identity.type = 'user';
		store.identity.userID = response.user.userID;
		store.identity.sessionID = response.token;
		console.log(store.identity);
		showHome(store.identity)
	}
}
