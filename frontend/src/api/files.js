import axios from 'axios';
const DEV_API_BASE = "http://localhost:5173/api";

/*
 * Fetch files using upload ID
 */
export async function fetchFiles(uploadId) { 
	const filesURL = DEV_API_BASE + "/files";
	try { 
		const response = await axios.post(filesURL, uploadId);
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
