/* Wrapper functions for showing a page based on state and user */
import { store, subscribe } from './state.js';
import { Landing, Signup, Login, Home } from './pages';

function renderPage() { 
	const root = document.getElementById('root');
	switch (store.currentPage) { 
		case 'landing':
			Landing(root);
			break;
		case 'login':
			Login(root);
			break;
		case 'signup':
			Signup(root);
			break;
		case 'home': 
			Home(root, store.user);
			break;

		/*
		 * Additional pages go here
		 *
		 *
		 *
		 */
	}
}
// Subscribe globally and render landing page
subscribe(renderPage);
renderPage();

/* Show landing page */
export function showLanding() { 
	logPageView('landing');
	store.currentPage = 'landing';
}

/* Show sign up page */
export function showSignup() {
	// notifiy server on page visit
	logPageView('signup');
	store.currentPage = 'signup';
}

/* Show log in page */
export function showLogin() { 
	// notify server on page visit
	logPageView('login');
	store.currentPage = 'login';
}

/* Show homepage */
export function showHome(userData) { 
	// notify server on page visit
	logPageView('home');
	store.user = userData;
	store.isLoggedIn = true;
	store.currentPage = 'home';
	store.loginError = null;
	/*
	 *
	 * Change UI based on user and state
	 * 
	 */
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
