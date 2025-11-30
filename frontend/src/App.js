import './style.css';
import { awaitServer } from './api';
import reorgLogo from '/reorg-logo.png';

export default function App() {
	// Make sure server is up and running
	awaitServer().then(() => {
		console.log("Server healthy");
	});

	/* this will have an option to login or create an account 
	 * 
	 * once either one of these options is clicked, render correct
	 * view.
	 * */
	const root = document.getElementById('root');
	root.innerHTML = `

		<div class='banner'>
			<div class='welcome-message'>
				<img src="${reorgLogo}" class="logo vanilla" alt="ReOrg logo" />
				<h1> welcome </h1>
			</div>

			</div>
				<!-- This will be a button -->
				<p> Log in </p>
			<div>
			
			</div>

				<!-- This will be a button -->
				<p> Sign up </p>
			<div>
		</div> 
	`;
}
