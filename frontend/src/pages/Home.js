/*
 * The Home Page
 */

export default function Home(root, userData) { 
	console.log(userData);
	const userID = userData.userID;
	const username = userID.username;
	root.innerHTML = `
		<div id='home-root'>
			<!-- Render navbar component here -->
			<div id='home-welcome-banner'>
				<h1 id='home-welcome-msg'>
					Welcome, ${username}
				</h1>
			</div>
		</div>
	`;
}
