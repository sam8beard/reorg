import axios from 'axios'; 
const DEV_API_BASE = "http://localhost:5173/api";

export async function saveOrgData(orgData) {
	const saveOrgURL = DEV_API_BASE + "/organize/save";
	try {
		const token = sessionStorage.getItem('authToken');
		const response = await axios.post(saveOrgURL, orgData, {
			headers: {
				'Authorization': `Bearer ${token}`
			},
		});
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
