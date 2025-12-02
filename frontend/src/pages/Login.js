/*
 * The Login Page
 */
import { attachLoginHandler } from '../events';
import { store } from '../state.js';

export default function Login(root) {
	root.innerHTML = `
		<div id='login-root'>
			<div id='login-banner'>
				<h1> Log in to your account </h1>
			</div>
			<div id='login-form-container'>
				<form id='login-form' method='post'>
					<div id='username-container'>
						<label class='' for='username'>
							Username
						</label>
						<input type='text' id='username' name='username' placeholder='Username'>
					</div>
					<div id='password-container'>
						<label class='' for='password'>
							Password
						</label>
						<input type='text' id='password' name='password' placeholder='Password'>
					</div>
					<div class='' id='login-error-container'>
						<p id='login-error'></p>
					</div>
					<button type='submit' style='display:none'></button>
				</form>
			</div>
		</div>
	`;

	// If login error is not null, render error message
	const errEl = root.querySelector('#login-error');
	errEl.textContent = (store.loginError) ? "Wrong username or password. Try again." : ""; 


	// Attach handler for login
	const loginForm = root.querySelector('#login-form');
	attachLoginHandler(loginForm, root);
}



