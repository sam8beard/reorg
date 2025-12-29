import { attachSignupHandler } from '../events';
/*
 * The Signup Page
 */

export default function Signup(root) { 
	root.innerHTML = `
		<div id='signup-root'>
			<div id='signup-banner'>
				<h1> Create an account </h1>
			</div>
			<div id='signup-form-container'>
				<form id='signup-form' method='post'>
					<div id='signup-email-container'>
						<label class='' for='email'>
							Email
						</label>
						<input type='text' id='email' name='email' placeholder='Email'>
					</div>
					<div id='signup-username-container'>
						<label class='' for='username'>
							Username
						</label>
						<input type='text' id='username' name='username' placeholder='Username'>
					</div>
					<div id='signup-password-container'>
						<label class='' for='password'>
							Password
						</label>
						<input type='text' id='password' name='password' placeholder='Password'>
					</div>
					<div class='' id='signup-error-container'>
						<p id='signup-error'></p>
					</div>
					<button type='submit' style='display:none'></button>
				</form>
			</div>
		</div>
	`;

	// Attach handler for signup
	const signupForm = root.querySelector('#signup-form');
	attachSignupHandler(signupForm, root);

}
