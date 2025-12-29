import '../../src/style.css';
import reorgLogo from '../assets/reorg-logo.png';
import { attachLoginPageHandler, attachSignupPageHandler, attachGuestHandler, attachUploadPageHandler } from '../events';

/*
 * Render landing page
 */
export default function Landing(root) { 
	root.innerHTML = `

		<div class='banner'>
			<div class='welcome-message'>
				<img src="${reorgLogo}" class="reorg-logo" alt="ReOrg logo" />
			</div>
			<div class='actions'>
				<div class='account-options'>
					<button id="login-btn">Log In</button>
					<button id="signup-btn">Sign Up</button>
				</div>
				<div class='guest-option'>
					<button id='guest-btn'>Use without an account</button>
				</div>
			<div>
		</div> 
	`;

	// Attach handlers for navigation options
	const loginBtn = root.querySelector('#login-btn')
	const signupBtn = root.querySelector('#signup-btn')
	const guestBtn = root.querySelector('#guest-btn');
	attachLoginPageHandler(loginBtn, root);
	attachSignupPageHandler(signupBtn, root);
//	attachUploadPageHandler(guestBtn, root);
	attachGuestHandler(guestBtn, root);
}
