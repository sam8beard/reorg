import axios from 'axios';
const DEV_API_BASE = "http://localhost:5173/api";

export async function postRuleJson(ruleJson) { 
	const ruleURL = DEV_API_BASE + "/rule";
	try { 
		const response = await axios.post(ruleURL, ruleJson);
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
