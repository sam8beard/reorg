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
export function attachSignupHandler(form) {
	form.addEventListener('submit', onSignupSubmit);

	// Validate input in real time
	form.elements.email.addEventListener('input', () => 
		debounce(() => validateEmail(form), 600)
	);

	form.elements.password.addEventListener('input', () => 
		debounce(() => validatePassword(form), 600)
	);

	form.elements.passwordMatch.addEventListener('input', () => 
		debounce(() => validatePasswordMatch(form), 600)
	);
}
/*
 * Handles sign up submit
 *
 * Create account....etc. 
 */
export async function onSignupSubmit(e) {
	e.preventDefault();

	const form = e.target;

	const emailValid = validateEmail(form);
	const passValid = validatePassword(form);
	const matchValid = validatePasswordMatch(form);

	// Block submit on invalid input 
	if (!emailValid || !passValid || !matchValid) {
		return;
	}
	
	// Build signup request object
	const signupRequest = {
		"email": form.elements.email.value.trim(),
		"username": form.elements.username.value.trim(),
		"password": form.elements.password.value,
	}
	
	const response = await createAccount(signupRequest);

	if (response.err) {
		showFormError(form, response.err);
		return;
	}

	// Debugging 
	console.log("Successfully registered user");

	// Success
	sessionStorage.setItem('authToken', response.data.token);
	sessionStorage.setItem('userID', response.data.user.userID);

	store.identity.type = 'user';
	store.identity.userID = response.data.user.userID;
	store.identity.sessionID = response.data.token;

	showHome(store.identity);
}

function validateEmail(form) {
	const email = form.elements.email;
	const regex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

	const isValid = regex.test(email.value.trim());
	return toggleValidity(email, isValid, 'Email is invalid');
}

function validateUsername(form) {
	const username = form.elements.username;
	// only letters and numbers, at least one letter
	const regex = /^(?=.*[A-Za-z])[A-Za-z0-9]{4,15}$/;
	// greater than length of 3 chars
	const isValid = (username.value.trim().length > 3) && regex.test(username.value.trim());
	return toggleValidity(username, isValid, 'Username must be at least 4 characters, at most 15 characters, only contain letters or numbers, and must have at least 1 letter');
}

function validatePassword(form) {
	const password = form.elements.password;
	// at least 8 chars, one number, one symbol
	const regex = /^(?=.*[\d])(?=.*[\W]).{8,}$/;

	const isValid = regex.test(password.value);

	return toggleValidity(
		password,
		isValid,
		'Password needs 8+ chars, number & symbol'
	);
}

function validatePasswordMatch(form) {
	const password = form.elements.password;
	const match = form.elements.passwordMatch;

	const isValid = password.value === match.value;
	return toggleValidity(match, isValid, 'Passwords do not match');
}

function toggleValidity(input, isValid, msg) {
	let error = input.nextElementSibling;

	if (!error || !error.classList.contains('error')) {
		error = document.createElement('span');
		error.className = 'error';
		input.after(error);
	}

	error.textContent = isValid ? '' : msg;
	input.classList.toggle('is-invalid', !isValid);
	input.setAttribute('aria-invalid', String(!isValid));

	return isValid;
}

// Debounce for live input validation
let debounceTimer;

function debounce(fn, wait) {
	clearTimeout(debounceTimer);
	debounceTimer = setTimeout(fn, wait);
}

// Error on input submission
function showFormError(form, msg) {
	let err = form.querySelector('#form-error');

	if (!err) {
		err = document.createElement('p');
		err.className = 'form-error';
		err.setAttribute('role', 'alert');
		form.prepend(err);
	}

	err.textContent = msg;
}
