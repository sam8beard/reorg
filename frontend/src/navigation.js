/* Wrapper functions for showing a page based on state and user */
import { store, subscribe } from './state.js';
import { Landing, Signup, Login, Home, Upload, Organize } from './pages';

export function renderPage() { 
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
		case 'upload':
			Upload(root, store.user);
			break;
		case 'organize':
			Organize(root, store.user);
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
	logPageVisit('landing');
	store.currentPage = 'landing';
}

/* Show sign up page */
export function showSignup() {
	// notifiy server on page visit
	logPageVisit('signup');
	store.currentPage = 'signup';
}

/* Show log in page */
export function showLogin() { 
	// notify server on page visit
	logPageVisit('login');
	store.currentPage = 'login';
}

/* Show home page */
export function showHome(userData) { 
	// notify server on page visit
	logPageVisit('home');
	store.user = userData;
	store.isLoggedIn = true;
	store.loginError = null;
	store.currentPage = 'home';
	/*
	 *
	 * Change UI based on user and state
	 * 
	 */
}

/* Show upload page */
export function showUpload() { 
	// notify server on page visit
	logPageVisit('upload');
	store.currentPage = 'upload';
}

/* Show organize page */
export function showOrganize() { 
	// notify server on page visit
	logPageVisit('organize');
	store.currentPage = 'organize';
}


/* Util for logging page visit */
function logPageVisit(pageName) { 
	console.log(`Page visit: ${pageName}`);
	
	/*
	 * Add optional reactions here...
	 *
	 * Send analytics to backend, etc.
	 */
} 
