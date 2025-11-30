/* Wrapper functions for showing a page based on state and user */
import { Signup, Login, Home } from './pages'
/* Show sign up page */
export function showSignupPage(root) {
	// notifiy server on page visit
	logPageView('signup');
	Signup(root);
}

/* Show log in page */
export function showLoginPage(root) { 
	// notify server on page visit
	logPageView('login');
	Login(root);
}

/* Show homepage */
export function showHomePage(root, user) { 
	// notify server on page visit
	logPageView('home');
	/*
	 *
	 * Change UI based on user and state
	 * 
	 */
	Home(root, user);
}

/* Util for logging page visit */
function logPageView(pageName) { 
	console.log(`Page viewed: ${pageName}`);
	
	/*
	 * Add optional reactions here...
	 *
	 * Send analytics to backend, etc.
	 */
} 
