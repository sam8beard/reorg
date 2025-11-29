import axios from 'axios';
const DEV_API_BASE = "http://localhost:5173/api";

export async function fetchHealth() {
	var healthURL = DEV_API_BASE + "/health";
	try {
		const response = await axios.get(healthURL);
		console.log(response.data);
		return response
	} catch (error) {
		console.error('Error fetching health:', error);
	}
}


export async function awaitServer() { 
	let healthy = false;
	while (!healthy) { 
		try { 
			const response = await fetchHealth();
			if (response.status === 200) { 
				healthy = true;
			} else { 
				console.log('Bad response from server, retrying...')
				await new Promise(r => setTimeout(r, 2000));
			} 
		} catch (err) { 
			console.log('Could reach server, retrying...')
			await new Promise(r => setTimeout(r, 2000));
		} 
	}
} 

