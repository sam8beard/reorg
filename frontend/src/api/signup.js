import axios from 'axios';
const DEV_API_BASE = "http://localhost:5173/api";

export async function createAccount(signupRequest) {
	const signupURL = DEV_API_BASE + "/auth/signup";
	try {
		const response = await axios.post(signupURL, signupRequest);
		console.log(response);
		return { data: response.data, err: null };
	} catch (err) {
		if (err.response) {
			let msg = 'Unknown error';
			if (err.response.status === 409) {
				// Account creation conflict
				console.log(err.response.status);
				msg = 'An account with this email or username already exists.';
			} else if (err.response.status === 400) {
				// Invalid signup form input
				msg = 'Invalid account creation details, please try again.';
			}
			return { data: null, err: msg };
		} else { 
			return { data: null, err: 'Network error' };
		} 
	}
}
