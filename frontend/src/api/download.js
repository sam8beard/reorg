import axios from 'axios'; 
const DEV_API_BASE = "http://localhost:5173/api";

/* Downloads zip file of organized file structure */
export async function downloadZip(organizeResult) {
	console.log("firing");
	const downloadURL = DEV_API_BASE + "/download/zip";
	try {
		// Might have to end up sending eval result with get request, serving zip file from the backend, and then downlaod from the frontend
		console.log("firing");
		const response = await axios.post(downloadURL, organizeResult, { responseType: 'blob' });
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
