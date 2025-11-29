import reorgLogo from '../assets/reorg-logo.png';
import '../styles/main.css';
import { fetchHealth } from '../api';
export function renderHomeView() {

	const app = document.querySelector('#app');
	app.innerHTML = `

		<div class='banner'>
			<div class='welcome-message'>
				<img src="${reorgLogo}" class="logo vanilla" alt="ReOrg logo" />
				<h1> welcome </h1>
			</div>
		</div> 
	`;

	/* Test with fetchHealth */
	
//	// Make header element
//	const healthHeader = document.createElement('h2')
//
//	// Fetch server health and display message accordingly
//	fetchHealth().then(value => {
//		console.log(value)
//		var msg = ""	
//		if (value.healthy == true) {
//			msg = 'ready to receive requests'
//		} else {
//			msg = 'not responding'
//		} 
//		healthHeader.textContent = `Server status: ${msg}`
//		app.appendChild(healthHeader)
//
//	}).catch(error => {
//		console.error('Error: ', error)	
//
//		healthHeader.textContent = `Could not fetch server status: ${error}`
//		app.appendChild(healthHeader)
//
//	});

	// example structure for homepage 
	//const app = document.querySelector('#app');

	//  app.innerHTML = `
	//    <h1>Home</h1>
	//    <button id="settings-btn">Settings</button>
	//  `;

	//  document
	//    .querySelector('#settings-btn')
	//    .addEventListener('click', () => {
	//      renderSettingsView();
	//    });
}
