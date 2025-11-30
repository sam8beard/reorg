import axios from 'axios';
const DEV_API_BASE = "http://localhost:5173/api";

export async function fetchUser(user) { 
	var userURL = DEV_API_BASE + "/user";
	console.log(userURL);
	try {
		username = user.username
		password = user.password
		const response = await axios.post(userURL, {
				username: username,
				password: password,
		});
		console.log(response.data);
		return response.data;
	} catch (err) { 
		if (err.response) { 
			console.error(err)
			return { error: err.response.data.error || 'Unknown error' };
		} else { 
			console.error(err)
			return { error: 'Network error' };
		} 
	} 
}
