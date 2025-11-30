import './style.css';
import { awaitServer } from './api';
import reorgLogo from '/reorg-logo.png';
import { Landing } from './pages';

/*
 * Root of application
 *
 * Loads landing page as initial view.
 */
export default function App(root) {
	// Ensure server is healthy
	awaitServer().then(() => {
		console.log("Server healthy");
	}).finally(() => {
		Landing(root)
	});
}
