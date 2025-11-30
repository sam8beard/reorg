/*
 * The Login Page
 */
import { attachLoginHandler } from '../events';
import { store, subscribe } from '../state.js';

export default function Login(root) {
	root.innerHTML = `
		<div id='login-root'>
			<div id='login-banner'>
				<h1> Log In to your account </h1>
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
					<div class='' id='login-error'></div>
					<button type='submit' style='display:none'></button>
				</form>
			</div>
		</div>
	`;
	 // Subscribe to loginError changes
	subscribe((prop, value) => {
		if (prop === 'loginError') {
		    const errorEl = root.querySelector('#login-error');
		    if (errorEl) errorEl.textContent = value || '';
		}
	});
	// Attach handler for login
	const loginForm = root.querySelector('#login-form');
	attachLoginHandler(loginForm, root);
}



