import cors from 'cors';
import './style.css';
import reorgLogo from '/reorg-logo.png';
import { fetchHealth } from './api';
import { renderHomeView } from './views/Home.js';
import { awaitServer } from './api/health.js'
//import { renderLoadingView } from './views/Loading.js'

// Renders initial view and setup event listeners
async function init() {
	// Check backend server health
	//renderLoadingView();
	await awaitServer();


	// Render homepage
	renderHomeView();


	/*
	 * 
	 * Attach event listeners
	 *
	 */
}
init();
