import axios from 'axios';
const DEV_API_BASE = "http://localhost:5173/api";

/* Post target info to backend */
export async function postTarget(targetUUID, targetName) { 
	const targetURL = DEV_API_BASE + "/target";
	try { 
		console.log(targetUUID, targetName);
		const response = await axios.post(targetURL, {
				targetUUID: targetUUID,
				targetName: targetName,
		});
		console.log(response);
		console.log(response.data);
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
