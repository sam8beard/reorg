import './style.css';
import { awaitServer } from './api';
import reorgLogo from '/reorg-logo.png';
import { showLanding } from './navigation.js'
/*
 * Root of application
 *
 * Loads landing page as initial view.
 */
export default function App(root) {
	// Ensure server is healthy
	awaitServer().finally(() => {
		console.log("Server healthy");
		showLanding();
	});
}
