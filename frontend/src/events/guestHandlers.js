import { createGuest } from '../api';
import { showUpload } from '../navigation.js';
export function attachGuestHandler(button, root) {
	button.addEventListener('click', () => onGuestClick());
}

export async function onGuestClick(root) {
	const response = await createGuest();
	const data = response;
	console.log(data);
	sessionStorage.setItem('authToken', data.token);
	sessionStorage.setItem('guestID', data.guestID);

	showUpload();
}
