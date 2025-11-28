import './style.css'
import reorgLogo from '/reorg-logo.png'

document.querySelector('#app').innerHTML = `
	<div class='banner'>
		<div class='welcome-message'>
			<img src="${reorgLogo}" class="logo vanilla" alt="ReOrg logo" />
			<h1> welcome </h1>
		</div>
	</div> 
	<div>
	</div>

`

