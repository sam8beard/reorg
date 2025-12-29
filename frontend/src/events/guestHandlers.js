import { createGuest } from '../api';
import { showUpload } from '../navigation.js';
import { store } from '../state.js';
export function attachGuestHandler(button, root) {
	button.addEventListener('click', () => onGuestClick(root));
}

export async function onGuestClick(root) {
	const response = await createGuest();
	if (response.err) {
		// handle error here
	} else {
		// Save token and ID in store
		store.identity.type = 'guest';
		store.identity.userID = response.guestID;
		store.identity.sessionID = response.token;
		console.log(store.identity);
		// Persist in session storage for page reloads
		sessionStorage.setItem('authToken', response.token);
		sessionStorage.setItem('guestID', response.guestID);
	}
	

	// Show the upload page
	showUpload();
}
