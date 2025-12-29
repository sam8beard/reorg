import axios from 'axios';
const DEV_API_BASE = "http://localhost:5173/api";

export async function createGuest() {
	const guestURL = DEV_API_BASE + "/auth/guest";
	try {
		const response = await axios.post(guestURL);
		console.log(response);
		sessionStorage.setItem('authToken', data.token);
		sessionStorage.setItem('guestID', data.guestID);
		return response.data;
	} catch (err) {
		if (err.response) { 
			console.log(err); 
			return { error: err.response.data.error || 'Unknown error' };
		} else { 
			console.log(err);
			return { error: 'Network error' };
		} 
	}
}
