import axios from 'axios';
const DEV_API_BASE = "http://localhost:5173/api";

/* Post file data to backend */
export async function uploadFileForm(formData) { 
	const uploadURL = DEV_API_BASE + "/upload";
	try { 
		console.log(formData);
		const response = await axios.post(uploadURL, formData);
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

