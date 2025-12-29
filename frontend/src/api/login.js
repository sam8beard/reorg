import axios from 'axios';
const DEV_API_BASE = "http://localhost:5173/api";

export async function fetchUser(loginRequest) {
	const loginURL = DEV_API_BASE + "/auth/login";
	try {
		const response = await axios.post(loginURL, loginRequest);
		console.log(response);
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
