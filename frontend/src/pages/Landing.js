import '../../src/style.css';
import reorgLogo from '../assets/reorg-logo.png';
import { attachLoginPageHandler, attachSignupPageHandler } from '../events';

/*
 * Render landing page
 */
export default function Landing(root) { 
	root.innerHTML = `

		<div class='banner'>
			<div class='welcome-message'>
				<img src="${reorgLogo}" class="logo vanilla" alt="ReOrg logo" />
				<h1> welcome </h1>
			</div>
				<!-- 
				These will be buttons that render the homepage and
				dashboard (will probably have separate component for dashboard) or the signup view. 
				On click, will fetch user data that the homepage is dependent on. 

				For now, we will just render a default homepage with dummy user data.
				-->
			<div class='actions'>
				<button id="login-btn">Log In</button>
				<button id="signup-btn">Sign Up</button>
			<div>
		</div> 
	`;

	// Attach handlers for navigation options
	const loginBtn = root.querySelector('#login-btn')
	const signupBtn = root.querySelector('#signup-btn')
	attachLoginPageHandler(loginBtn, root);
	attachSignupPageHandler(signupBtn, root);
}
