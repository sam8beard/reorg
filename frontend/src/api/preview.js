import axios from 'axios'; 
const DEV_API_BASE = "http://localhost:5173/api";

/* Posts rule set object to backend and gets back preview json object */
export async function getPreviewJson(ruleSet) { 
	const previewURL = DEV_API_BASE + "/preview";
	try {
		const response = await axios.post(previewURL, ruleSet);
		return response.data;
	} catch(err) { 
		if (err.response) {
			console.log(err); 
			return { error: err.response.data.error || 'Unknown error' };
		} else { 
			console.log(err);
			return { error: 'Network error' };
		}
	}
}
