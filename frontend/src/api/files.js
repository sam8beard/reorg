import axios from 'axios';
const DEV_API_BASE = "http://localhost:5173/api";

/*
 * Fetch files using upload ID
 */
export async function fetchFiles(uploadID) { 
	const filesURL = DEV_API_BASE + "/files";
	try { 
		const token = sessionStorage.getItem('authToken');
		const response = await axios.post(filesURL, uploadID, {
			headers: {
				'Authorization': `Bearer ${token}`
			}
		});
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
